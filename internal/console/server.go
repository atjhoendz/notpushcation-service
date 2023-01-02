package console

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/atjhoendz/notpushcation-service/internal/config"
	"github.com/atjhoendz/notpushcation-service/internal/delivery/graphqlsvc/graph"
	"github.com/atjhoendz/notpushcation-service/internal/delivery/graphqlsvc/graph/generated"
	"github.com/atjhoendz/notpushcation-service/internal/onesignal"
	"github.com/atjhoendz/notpushcation-service/internal/usecase"
	"github.com/kumparan/go-connect"
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

	graphqlServer := echo.New()

	httpClient := connect.NewHTTPConnection(config.DefaultHTTPOptions())

	notificationComposer := usecase.NewNotificationComposer()
	onesignalClient := onesignal.NewClient(httpClient, config.OnesignalAppID())
	messageProcessorUsecase := usecase.NewMessageProcessorUsecase(usecase.Composers{NotificationComposer: notificationComposer}, onesignalClient)

	go func() {
		for {
			select {
			case <-sigCh:
				gracefulShutdown(graphqlServer)
				quitCh <- true
			case e := <-errCh:
				log.Error(e)
				gracefulShutdown(graphqlServer)
				quitCh <- true
			}
		}
	}()

	go func() {
		graphqlServer.Pre(middleware.AddTrailingSlash())
		graphqlServer.Use(middleware.Logger())
		graphqlServer.Use(middleware.Recover())

		c := generated.Config{Resolvers: &graph.Resolver{MessageProcessorUsecase: messageProcessorUsecase}}

		graphqlHandler := handler.NewDefaultServer(
			generated.NewExecutableSchema(c),
		)

		if config.EnableIntrospection() {
			graphqlHandler.Use(extension.Introspection{})
		}

		gqlGroup := graphqlServer.Group("/query")
		gqlGroup.Any("/", func(c echo.Context) error {
			graphqlHandler.ServeHTTP(c.Response(), c.Request())
			return nil
		})

		// start graphql server
		err := graphqlServer.Start(fmt.Sprintf(":%s", config.HTTPPort()))
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
