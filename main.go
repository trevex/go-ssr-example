package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/trevex/go-ssr-example/views"
)

type Map map[string]interface{}

func DefaultMap(c *gin.Context) Map {
	req := c.Request
	m := Map{
		"HX": map[string]interface{}{
			"Partial": req.Header.Get("HX-Request") != "",
			"Target":  req.Header.Get("HX-Target"),
		},
		"Req": map[string]interface{}{
			"Path": req.URL.Path,
		},
	}
	return m
}

func (m Map) Field(key string, value interface{}) Map {
	m[key] = value
	return m
}

func NewServerCmd(log *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "server [flags]",
		SilenceErrors: true,
	}

	var (
		isDev bool
	)

	flags := cmd.Flags()
	flags.BoolVar(&isDev, "dev", false, "TODO")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		gin.SetMode(gin.ReleaseMode)

		if isDev { // Ok, we want to start the development setup
			gin.SetMode(gin.DebugMode)

			// TODO: do not use embed file systems
		}

		// Let's set up our gin router
		r := gin.New()
		r.HTMLRender = views.MustRenderer()
		r.Use(gin.Recovery())
		// TODO: add gzip, limit, logger middlewares

		// We use a group without any particular path to add our protected
		// routes and use our authentication middleware
		g := r.Group("/")
		// g.Use(p.OnlyAuthenticated)
		{
			g.StaticFS("/public/", views.PublicFileSystem())

			g.GET("/", func(c *gin.Context) {
				c.HTML(http.StatusOK, "index.html.tmpl", DefaultMap(c))
			})
			g.GET("/devices/", func(c *gin.Context) {
				c.HTML(http.StatusOK, "devices.html.tmpl", DefaultMap(c))
			})
		}

		// Start HTTP server at :4000.
		log.Info("Starting HTTP server on http://localhost:4000...")
		srv := &http.Server{
			Addr:    ":4000",
			Handler: r,
		}
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Error("error listening on server", zap.Error(err))
			}
		}()

		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctxShutdown); err != nil {
			log.Error("error shutting down server", zap.Error(err))
			return err
		}
		return nil
	}

	return cmd
}

func main() {
	log, _ := zap.NewProduction()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rootCmd := &cobra.Command{
		Use:           "example action [flags]",
		SilenceErrors: true,
	}

	// Add all sub-commands
	rootCmd.AddCommand(NewServerCmd(log))

	// Make sure to cancel the context if a signal was received
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Warn("received signal", zap.String("signal", sig.String()))
		cancel()
	}()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Error("command failed", zap.Error(err))
		os.Exit(1)
	}
}
