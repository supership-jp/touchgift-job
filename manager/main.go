package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"touchgift-job-manager/config"
	"touchgift-job-manager/infra"
	"touchgift-job-manager/injector"

	"github.com/gin-gonic/gin"
)

func SignalContext(ctx context.Context, logger *infra.Logger) (context.Context, context.CancelFunc) {
	parent, cancelParent := context.WithCancel(ctx)
	go func() {
		defer cancelParent()
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sig)
		select {
		case <-parent.Done():
			logger.Info().Msg("context done")
			return
		case s := <-sig:
			switch s {
			case syscall.SIGINT, syscall.SIGTERM:
				logger.Info().Msg("signal stop")
				return
			}
		}
	}()
	return parent, cancelParent
}

func listenAndServe(ctx context.Context, port string, router *gin.Engine, logger *infra.Logger) chan error {
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 20 * time.Second,
	}

	go func() {
		logger.Info().Msgf("Listening and serving HTTP on :%s", port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to ListenAndServe")
		}
	}()
	ch := make(chan error)
	go func() {
		for {
			<-ctx.Done()
			logger.Info().Msgf("Server shutting down...:%s", port)
			timeout, tCancel := context.WithTimeout(context.Background(), config.Env.Server.ShutdownTimeout)
			defer tCancel()
			err := srv.Shutdown(timeout)
			logger.Info().Msgf("Server shutdown :%s", port)
			ch <- err
			return
		}
	}()
	return ch
}

func run(sctx context.Context, cancel context.CancelFunc, logger *infra.Logger) {
	var wg sync.WaitGroup
	done := make(chan bool, 1)
	defer close(done)

	adminRouter, initializeAdmin, terminateAdmin := injector.AdminRoute(sctx, infra.NewRouter(logger))
	adminServerCh := listenAndServe(sctx, config.Env.Server.AdminPort, adminRouter, logger)
	wg.Add(1)
	logger.Info().Msg("Start to initialize for admin.")
	if err := initializeAdmin(); err != nil {
		logger.Error().Err(err).Msg("Failed to initialize for admin.")
		cancel()
		return
	}

	router := injector.Route(infra.NewRouter(logger))
	serverCh := listenAndServe(sctx, config.Env.Server.Port, router, logger)
	wg.Add(1)

	go func() {
		wg.Wait()
		done <- true
	}()
	for {
		select {
		case err := <-serverCh:
			logger.Info().Msg("Server stopped")
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to server stop: ")
			}
			wg.Done()
		case err := <-adminServerCh:
			logger.Info().Msg("Terminating for admin server.")
			if err := terminateAdmin(); err != nil {
				logger.Error().Err(err).Msg("Failed to terminate for admin.")
			}
			logger.Info().Msg("Admin server stopped")
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to admin server stop: ")
			}
			wg.Done()
		case <-done:
			logger.Info().Msg("Stopped")
			return
		}
	}
}

func main() {
	logger := infra.GetLogger()
	ctx, cancel := SignalContext(context.Background(), logger)
	run(ctx, cancel, logger)
}
