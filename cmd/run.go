/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
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
	"github.com/tupyy/tinyedge-controller/internal/repo/cache"
	"github.com/tupyy/tinyedge-controller/internal/repo/certificate"
	deviceRepo "github.com/tupyy/tinyedge-controller/internal/repo/postgres"
	"github.com/tupyy/tinyedge-controller/internal/servers"
	"github.com/tupyy/tinyedge-controller/internal/services/edge"
	edgePb "github.com/tupyy/tinyedge-controller/pkg/grpc/edge"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
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

		certificateManager := certificate.New(vaultClient, "pki", "home.net")
		deviceRepo, err := deviceRepo.New(pgClient)
		if err != nil {
			zap.S().Fatal(err)
		}
		configurationRepo := cache.New()

		edgeService := edge.New(deviceRepo, configurationRepo, certificateManager)
		edgeServer := servers.NewEdgeServer(edgeService)

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
		if err != nil {
			zap.S().Fatalf("failed to listen: %v", err)
		}

		// start edge server
		var opts []grpc.ServerOption

		zapOpts := []grpc_zap.Option{
			grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
				return zap.Int64("grpc.time_ns", duration.Nanoseconds())
			}),
		}
		opts = append(opts, grpc_middleware.WithUnaryServerChain(
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
