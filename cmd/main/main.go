package main

import (
	"context"
	"net"
	"net/http"

	"github.com/spf13/viper"
	"github.com/tvanriel/db-gui/app/connectionpool"
	"github.com/tvanriel/db-gui/app/domain"
	"github.com/tvanriel/db-gui/app/frontend"
	"github.com/tvanriel/db-gui/app/httpinterface"
	"github.com/tvanriel/db-gui/app/mysqlconnection"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Configuration struct {
	ListenAddress string
	DevMode       bool
}

func NewConfiguration() (*Configuration, error) {

	viper.SetDefault("Listen", "127.0.0.1:8080")
	viper.AddConfigPath("/etc/dbg/config")
	viper.AddConfigPath("/opt/dbg/config")
	viper.AddConfigPath("$HOME/.config/dbg")
	viper.AddConfigPath(".")

	viper.BindEnv("Listen", "LISTEN")
	viper.BindEnv("Dev", "DEV")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()

	if err != nil {
		viper.SafeWriteConfig()
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {

			return nil, err
		}
	}

	return &Configuration{
		ListenAddress: viper.GetString("Listen"),
		DevMode:       viper.GetBool("Dev"),
	}, nil

}

func NewDatabaseGuiHandlerMux(dbg httpinterface.DbGuiHandler) http.Handler {
	return dbg.Handler()
}

func NewMemoryConnectionPool() domain.ConnectionPool {
	return connectionpool.NewMemoryConnectionPool(mysqlconnection.NewMySQLConnection)
}

func NewEmbedBackendFrontend() frontend.Frontend {
	return frontend.NewEmbedBackendFrontend()
}

func NewFilesystemBackedFrontend() frontend.Frontend {
	return frontend.NewFilesystemBackedFrontend()
}

func NewFrontend(config *Configuration, log *zap.Logger) frontend.Frontend {
	if config.DevMode {
		log.Named("Frontend").Info("Starting frontend in Filesystem mode")
		return NewFilesystemBackedFrontend()
	}
	return NewEmbedBackendFrontend()
}

func NewLogger(config *Configuration) *zap.Logger {
	if config.DevMode {
		logger, _ := zap.NewDevelopment()
		return logger
	}
	logger, _ := zap.NewProduction()
	return logger
}

func NewDatabaseGuiHandler(pool domain.ConnectionPool, frontend frontend.Frontend) httpinterface.DbGuiHandler {
	return *httpinterface.NewDbGuiHandler(pool, frontend)
}

func NewHttpServe(lc fx.Lifecycle, config *Configuration, mux http.Handler) *http.Server {
	srv := &http.Server{
		Addr:    config.ListenAddress,
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

func main() {
	app := fx.New(
		fx.Provide(
			NewConfiguration,
			NewLogger,
			NewHttpServe,
			NewDatabaseGuiHandler,
			NewDatabaseGuiHandlerMux,
			NewMemoryConnectionPool,
			NewFrontend,
		),
		fx.Invoke(func(*http.Server) {}),
	)

	app.Run()

}
