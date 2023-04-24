/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/spf13/cobra"
	"github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/clients/vault"
	"github.com/tupyy/tinyedge-controller/internal/configuration"
	"github.com/tupyy/tinyedge-controller/internal/interceptors"
	"github.com/tupyy/tinyedge-controller/internal/repo"
	"github.com/tupyy/tinyedge-controller/internal/repo/git"
	"github.com/tupyy/tinyedge-controller/internal/servers"
	"github.com/tupyy/tinyedge-controller/internal/services"
	"github.com/tupyy/tinyedge-controller/internal/services/auth"
	"github.com/tupyy/tinyedge-controller/internal/services/repository"
	"github.com/tupyy/tinyedge-controller/internal/workers"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	edgePb "github.com/tupyy/tinyedge-controller/pkg/grpc/edge"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		logger := setupLogger()
		defer logger.Sync()

		undo := zap.ReplaceGlobals(logger)
		defer undo()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		conf := configuration.GetConfiguration()
		vaultClient, err := vault.NewVaultAppRoleClient(ctx, vault.VaultParameters{
			Address:         conf.VaultAddress,
			ApproleRoleID:   conf.VaultApproleRoleID,
			ApproleSecretID: conf.VaultAppRoleSecretID,
		})
		if err != nil {
			zap.S().Fatal(err)
		}
		certRepo := repo.NewCertificate(vaultClient, "pki_int", "home.net", "tinyedge-role")
		secretRepo := repo.NewSecret(vaultClient, "tinyedge")
		zap.S().Info("vault repositories created")

		pgClient, err := pg.New(pg.ClientParams{
			Host:     "localhost",
			Port:     5433,
			DBName:   "tinyedge",
			User:     "postgres",
			Password: "postgres",
		})
		if err != nil {
			zap.S().Fatal(err)
		}

		deviceRepo, err := repo.NewDeviceRepo(pgClient)
		if err != nil {
			zap.S().Fatal(err)
		}
		repoRepo, err := repo.NewRepository(pgClient)
		if err != nil {
			zap.S().Fatal(err)
		}
		manifestRepo, err := repo.NewManifest(pgClient)
		if err != nil {
			zap.S().Fatal(err)
		}
		// cacheRepo := cache.NewCacheRepo()

		// git repo
		gitRepo := git.New("/home/cosmin/tmp/git")

		// create services
		zap.S().Info("create services")
		certService := services.NewCertificate(certRepo)
		manifestService := services.NewManifest(deviceRepo, manifestRepo, gitRepo)
		deviceService := services.NewDevice(deviceRepo)
		configurationService := services.NewConfiguration(deviceService)
		edgeService := services.NewEdge(deviceRepo, configurationService, certService)
		authService := auth.New(certService, deviceRepo)
		repoService := repository.NewRepositoryService(repoRepo, gitRepo, secretRepo)

		scheduler := workers.New(5 * time.Second)
		scheduler.AddWorker(workers.NewGitOpsWorker(repoService, configurationService))
		go scheduler.Start(ctx)

		tlsConfig, err := certService.TlsConfig(ctx, conf.DefaultCertificateTTL)
		if err != nil {
			zap.S().Fatal(err)
		}
		zap.S().Info("tls config created")

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
		if err != nil {
			zap.S().Fatalf("failed to listen: %v", err)
		}

		connAdmin, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8081))
		if err != nil {
			zap.S().Fatalf("failed to listen: %v", err)
		}

		grpcEdgeServer := createEdgeServer(tlsConfig, authService, logger)
		edgeServer := servers.NewEdgeServer(edgeService)
		edgePb.RegisterEdgeServiceServer(grpcEdgeServer, edgeServer)
		go grpcEdgeServer.Serve(lis)

		grpcAdminServer := createAdminServer(logger)
		adminServer := servers.NewAdminServer(repoService, manifestService, deviceService, configurationService)
		admin.RegisterAdminServiceServer(grpcAdminServer, adminServer)
		grpcAdminServer.Serve(connAdmin)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func setupLogger() *zap.Logger {
	loggerCfg := &zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "severity",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder, EncodeCaller: zapcore.ShortCallerEncoder},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	plain, err := loggerCfg.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		panic(err)
	}

	return plain
}

func createEdgeServer(tlsConfig *tls.Config, auth *auth.Service, logger *zap.Logger) *grpc.Server {
	// start edge server
	creds := credentials.NewTLS(tlsConfig)
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Float64("grpc.time_s", duration.Seconds())
		}),
	}
	altOpts := []grpc_ctxtags.Option{
		grpc_ctxtags.WithFieldExtractor(func(fullMethod string, req interface{}) map[string]interface{} {
			type idStruct struct {
				DeviceID string `json:"device_id"`
			}
			var id idStruct
			m := make(map[string]interface{})
			d, err := json.Marshal(req)
			if err != nil {
				return m
			}
			if err := json.Unmarshal(d, &id); err != nil {
				return m
			}
			m["device_id"] = id.DeviceID
			return m
		}),
	}
	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		interceptors.AuthInterceptor(auth),
		grpc_ctxtags.UnaryServerInterceptor(altOpts...),
		grpc_zap.UnaryServerInterceptor(logger, zapOpts...),
	))

	grpc_zap.ReplaceGrpcLoggerV2(logger)
	return grpc.NewServer(opts...)
}

func createAdminServer(logger *zap.Logger) *grpc.Server {
	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Float64("grpc.time_s", duration.Seconds())
		}),
	}
	opts := []grpc.ServerOption{}
	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		grpc_zap.UnaryServerInterceptor(logger, zapOpts...),
	))

	grpc_zap.ReplaceGrpcLoggerV2(logger)
	return grpc.NewServer(opts...)
}
