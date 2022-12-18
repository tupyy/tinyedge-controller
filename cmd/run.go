/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/spf13/cobra"
	"github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/clients/vault"
	"github.com/tupyy/tinyedge-controller/internal/configuration"
	"github.com/tupyy/tinyedge-controller/internal/interceptors"
	"github.com/tupyy/tinyedge-controller/internal/repo/cache"
	certRepo "github.com/tupyy/tinyedge-controller/internal/repo/certificate"
	deviceRepo "github.com/tupyy/tinyedge-controller/internal/repo/postgres"
	"github.com/tupyy/tinyedge-controller/internal/servers"
	"github.com/tupyy/tinyedge-controller/internal/services/auth"
	"github.com/tupyy/tinyedge-controller/internal/services/certificate"
	"github.com/tupyy/tinyedge-controller/internal/services/edge"
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

		pgClient, err := pg.New(pg.ClientParams{
			Host:     "localhost",
			Port:     5432,
			DBName:   "tinyedge",
			User:     "postgres",
			Password: "postgres",
		})

		certRepo := certRepo.New(vaultClient, "pki_int", "home.net", "tinyedge-role")
		deviceRepo, err := deviceRepo.New(pgClient)
		if err != nil {
			zap.S().Fatal(err)
		}
		configurationRepo := cache.New()

		certService := certificate.New(certRepo)
		edgeService := edge.New(deviceRepo, configurationRepo, certService)
		authService := auth.New(certService, deviceRepo)
		edgeServer := servers.NewEdgeServer(edgeService)

		tlsConfig, err := createTlsConfig(
			"/home/cosmin/projects/tinyedge-controller/resources/certificates/ca.pem",
			"/home/cosmin/projects/tinyedge-controller/resources/certificates/cert.pem",
			"/home/cosmin/projects/tinyedge-controller/resources/certificates/key.pem",
		)

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
		if err != nil {
			zap.S().Fatalf("failed to listen: %v", err)
		}

		// start edge server
		creds := credentials.NewTLS(tlsConfig)
		opts := []grpc.ServerOption{grpc.Creds(creds)}

		zapOpts := []grpc_zap.Option{
			grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
				return zap.Int64("grpc.time_ns", duration.Nanoseconds())
			}),
		}
		opts = append(opts, grpc_middleware.WithUnaryServerChain(
			interceptors.AuthInterceptor(authService),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger, zapOpts...),
		))
		// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
		grpc_zap.ReplaceGrpcLoggerV2(logger)
		grpcServer := grpc.NewServer(opts...)
		edgePb.RegisterEdgeServiceServer(grpcServer, edgeServer)
		grpcServer.Serve(lis)
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

func createTlsConfig(caroot, certFile, keyFile string) (*tls.Config, error) {
	// read certificates
	caRoot, err := os.ReadFile(caroot)
	if err != nil {
		return nil, err
	}

	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("cannot copy system certificate pool: %w", err)
	}

	pool.AppendCertsFromPEM(caRoot)
	config := tls.Config{
		RootCAs:    pool,
		ClientCAs:  pool,
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	cert, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	privateKey, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	cc, err := tls.X509KeyPair(cert, privateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create x509 key pair: %w", err)
	}

	config.Certificates = []tls.Certificate{cc}

	return &config, nil
}
