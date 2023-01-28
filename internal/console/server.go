package console

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kumparan/go-connect"
	connectMiddleware "github.com/kumparan/go-connect/middleware"
	"github.com/ulule/limiter/v3"

	"github.com/nats-io/nats.go"

	"github.com/atjhoendz/notpushcation-service/internal/subscriber"
	"github.com/kumparan/ferstream"

	"github.com/r3labs/sse/v2"

	"github.com/atjhoendz/notpushcation-service/internal/db"
	"github.com/atjhoendz/notpushcation-service/internal/repository"

	"github.com/atjhoendz/notpushcation-service/internal/delivery/httpsvc"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/atjhoendz/notpushcation-service/internal/config"
	"github.com/atjhoendz/notpushcation-service/internal/delivery/graphqlsvc/graph"
	"github.com/atjhoendz/notpushcation-service/internal/delivery/graphqlsvc/graph/generated"
	"github.com/atjhoendz/notpushcation-service/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Run:   run,
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {
	sigCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	quitCh := make(chan bool, 1)

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	db.InitializeCockroachConn()

	httpServer := echo.New()

	onSubscribeCB := func(streamID string, sub *sse.Subscriber) {
		log.Infof("client connected, streamID: %s, url: %s", streamID, sub.URL)
	}
	onUnsubscribeCB := func(streamID string, sub *sse.Subscriber) {
		log.Infof("client disconnected, streamID: %s, url: %s", streamID, sub.URL)
	}

	//sseServer := sse.New()
	sseServer := sse.NewWithCallback(onSubscribeCB, onUnsubscribeCB)
	sseServer.AutoStream = true
	sseServer.AutoReplay = false

	threadRepository := repository.NewThreadRepository(db.CockroachDB)
	threadUsecase := usecase.NewThreadUsecase(threadRepository, sseServer)
	liveBlogPostUsecase := usecase.NewLiveBlogPostUsecase(sseServer)

	natsOpts := []nats.Option{
		nats.RetryOnFailedConnect(config.NATSJSRetryOnFailedConnect()),
		nats.MaxReconnects(config.NATSJSMaxReconnect()),
		nats.ReconnectWait(config.NATSJSReconnectWait()),
	}

	subs := subscriber.NewSubscriber(liveBlogPostUsecase)
	js, err := ferstream.NewNATSConnection(
		config.NATSJSHost(),
		[]ferstream.JetStreamRegistrar{
			liveBlogPostUsecase,
			subs,
		},
		natsOpts...,
	)
	continueOrFatal(err)
	defer ferstream.SafeClose(js)

	go func() {
		for {
			select {
			case <-sigCh:
				sseServer.Close()
				gracefulShutdown(httpServer)
				quitCh <- true
			case e := <-errCh:
				log.Error(e)
				gracefulShutdown(httpServer)
				quitCh <- true
			}
		}
	}()

	redisOpts := &connect.RedisConnectionPoolOptions{
		DialTimeout:     config.RedisDialTimeout(),
		ReadTimeout:     config.RedisReadTimeout(),
		WriteTimeout:    config.RedisWriteTimeout(),
		IdleCount:       config.RedisMaxIdleConn(),
		PoolSize:        config.RedisMaxActiveConn(),
		IdleTimeout:     config.RedisIdleTimeout(),
		MaxConnLifetime: config.RedisMaxConnLifetime(),
	}

	go func() {
		httpServer.Pre(middleware.AddTrailingSlash())
		httpServer.Use(middleware.Logger())
		httpServer.Use(middleware.Recover())

		apiGroup := httpServer.Group("")
		httpsvc.RouteService(apiGroup)

		c := generated.Config{Resolvers: &graph.Resolver{
			ThreadUsecase:       threadUsecase,
			LiveBlogPostUsecase: liveBlogPostUsecase,
		}}

		graphqlHandler := handler.NewDefaultServer(
			generated.NewExecutableSchema(c),
		)

		if config.EnableIntrospection() {
			graphqlHandler.Use(extension.Introspection{})
		}

		gqlGroup := httpServer.Group("/query")
		gqlGroup.Any("/", func(c echo.Context) error {
			graphqlHandler.ServeHTTP(c.Response(), c.Request())
			return nil
		})

		sseGroup := httpServer.Group("/sse")
		if config.EnableRateLimiter() {
			redisClient, err := connect.NewGoRedisConnectionPool(config.RedisRateLimiterHost(), redisOpts)
			continueOrFatal(err)

			sseIPRateLimiter, err := connectMiddleware.NewRedisIPRateLimiter(redisClient, limiter.Rate{
				Period: config.RateLimiterPeriod(),
				Limit:  config.RateLimiterRequestLimit(),
			}, config.RateLimiterExcludeIPs())
			continueOrFatal(err)

			sseGroup.Use(sseIPRateLimiter.Limit())
		}
		sseGroup.GET("/", func(c echo.Context) error {
			c.Response().Header().Add("Access-Control-Allow-Origin", "*")
			c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
			sseServer.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		})

		// start graphql server
		//err := httpServer.Start(fmt.Sprintf(":%s", config.HTTPPort()))
		err := httpServer.StartTLS(fmt.Sprintf(":%s", config.HTTPPort()), "cert.pem", "key.pem")
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	<-quitCh
	log.Info("exiting")
}

func gracefulShutdown(graphqlSvr *echo.Echo) {
	if graphqlSvr != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := graphqlSvr.Shutdown(ctx); err != nil {
			graphqlSvr.Logger.Fatal(err)
		}
	}
}

func continueOrFatal(err error) {
	if err != nil {
		log.Fatal()
	}
}
