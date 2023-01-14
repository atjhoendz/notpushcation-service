package httpsvc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/atjhoendz/notpushcation-service/internal/model"

	"github.com/kumparan/go-utils"
	"github.com/labstack/echo/v4"
)

type service struct {
	group     *echo.Group
	sseBroker *model.SSEBroker
}

type sseClient struct {
	name   string
	userID string
	events chan *sseEvent
}

type sseEvent struct {
	UserID int64 `json:"user_id"`
}

func RouteService(group *echo.Group) {
	srv := &service{
		group: group,
		//sseBroker: sseBroker,
	}
	srv.initRoutes()
}

func (s *service) initRoutes() {
	s.group.GET("/test-sse/", s.handleSSECobs())
	s.group.GET("/test-sse2/", s.handleThreadSSECobs())
	//s.group.GET("/sse", s.handleSSE())
	s.group.GET("/request/", func(c echo.Context) error {
		req := c.Request()
		format := `
				<code>
				  Protocol: %s<br>
				  Host: %s<br>
				  Remote Address: %s<br>
				  Method: %s<br>
				  Path: %s<br>
				</code>
			  `
		return c.HTML(http.StatusOK, fmt.Sprintf(format, req.Proto, req.Host, req.RemoteAddr, req.Method, req.URL.Path))
	})
}

func (s *service) handleSSECobs() echo.HandlerFunc {
	return func(c echo.Context) error {
		client := &sseClient{name: c.Request().RemoteAddr, userID: c.QueryParam("user_id"), events: make(chan *sseEvent, 10)}
		fmt.Printf("%+v", client)

		go s.mutateUser(client)

		c.Response().Header().Add("Access-Control-Allow-Origin", "*")
		c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
		c.Response().Header().Add("Content-Type", "text/event-stream")
		c.Response().Header().Add("Cache-Control", "no-cache")
		c.Response().Header().Add("Connection", "keep-alive")

		timeout := time.After(1 * time.Minute)
		select {
		case ev := <-client.events:
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			_ = enc.Encode(ev)
			_, _ = fmt.Fprintf(c.Response().Writer, "data: %v\n\n", buf.String())
			fmt.Printf("data: %v\n", buf.String())
		case <-timeout:
			_, _ = fmt.Fprintf(c.Response().Writer, "msg: nothing to sent\n\n")
		}

		if f, ok := c.(http.Flusher); ok {
			fmt.Println("flush")
			f.Flush()
		}

		return nil
	}
}

func (s *service) mutateUser(client *sseClient) {
	for {
		u := &sseEvent{UserID: utils.GenerateID()}
		client.events <- u
	}
}
