package app

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/kuzin57/shad-networks/internal/config"
	"github.com/kuzin57/shad-networks/internal/generated"
	graphcleaner "github.com/kuzin57/shad-networks/internal/pkg/graph_cleaner"
	graphgenerator "github.com/kuzin57/shad-networks/internal/pkg/graph_generator"
	"github.com/kuzin57/shad-networks/internal/pkg/visualizer"
	"github.com/kuzin57/shad-networks/internal/repositories"
	graphrepo "github.com/kuzin57/shad-networks/internal/repositories/graph"
	"github.com/kuzin57/shad-networks/internal/server"
	"github.com/kuzin57/shad-networks/internal/services/graph"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v3"

	slicesutils "github.com/kuzin57/shad-networks/internal/utils/slices"
)

func Create(confPath string) fx.Option {
	conf := mustLoadConfig(confPath)

	return fx.Options(
		fx.Supply(conf),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			NewGRPCServer,
			zap.NewProduction,
			server.NewServer,
			fx.Annotate(graph.NewService, fx.As(new(server.GraphService))),
			fx.Annotate(graphgenerator.NewGenerator, fx.As(new(graph.GraphGenerator))),
			fx.Annotate(repositories.NewNeo4jDriver, fx.As(new(repositories.Driver))),
			fx.Annotate(graphrepo.NewRepository, fx.As(new(graph.GraphRepository)), fx.As(new(graphcleaner.GraphRepository))),
			graphcleaner.NewCleaner,
			fx.Annotate(visualizer.NewVisualizer, fx.As(new(server.Visualizer))),
		),
		fx.Invoke(func(server *grpc.Server) {}),
	)
}

func NewGRPCServer(lc fx.Lifecycle, serverAPI *server.Server, log *zap.Logger, config *config.Config) *grpc.Server {
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("Recovered from panic", zap.Any("argument", p))

			return status.Errorf(codes.Internal, "internal error")
		})),
		logging.UnaryServerInterceptor(
			logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
				switch level {
				case logging.LevelDebug:
					log.Debug(msg, slicesutils.Map(fields, func(i int, field any) zap.Field {
						return zap.Any(fmt.Sprintf("arg_%d", i), field)
					})...)
				case logging.LevelError:
					log.Error(msg, slicesutils.Map(fields, func(i int, field any) zap.Field {
						return zap.Any(fmt.Sprintf("arg_%d", i), field)
					})...)
				case logging.LevelInfo:
					log.Info(msg, slicesutils.Map(fields, func(i int, field any) zap.Field {
						return zap.Any(fmt.Sprintf("arg_%d", i), field)
					})...)
				case logging.LevelWarn:
					log.Warn(msg, slicesutils.Map(fields, func(i int, field any) zap.Field {
						return zap.Any(fmt.Sprintf("arg_%d", i), field)
					})...)
				}
			}),
			logging.WithLogOnEvents(
				logging.PayloadReceived, logging.PayloadSent,
			),
		),
	))

	generated.RegisterGraphServer(srv, serverAPI)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l, err := net.Listen("tcp", fmt.Sprintf(":%d", config.App.Port))
			if err != nil {
				return fmt.Errorf("listen: %w", err)
			}

			log.Sugar().Infof("starting server on port %d", config.App.Port)

			go srv.Serve(l)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("stopping grpc server...")

			srv.GracefulStop()

			return nil
		},
	})

	return srv
}

func mustLoadConfig(path string) *config.Config {
	confContent, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	conf := &config.Config{}

	if err := yaml.Unmarshal(confContent, conf); err != nil {
		panic(err)
	}

	return conf
}
