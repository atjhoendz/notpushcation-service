package httpsvc

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/atjhoendz/notpushcation-service/internal/model"
	"github.com/labstack/echo/v4"
)

func (s *service) handleThreadSSECobs() echo.HandlerFunc {
	return func(c echo.Context) error {
		c = addSSEHeader(c)

		messageChan := make(chan *model.SSEMessage)
		s.sseBroker.AddNewClient(messageChan)

		//defer func() {
		//	s.sseBroker.CloseClient(messageChan)
		//}()

		timeout := time.After(1 * time.Minute)
		select {
		case msg := <-messageChan:
			_, _ = fmt.Fprintf(c.Response().Writer, "event: %s\n", msg.Event)
			_, _ = fmt.Fprintf(c.Response().Writer, "data: %s\n\n", msg.Data)
			log.Info("message written to buffer")
		case <-timeout:
			_, _ = fmt.Fprintf(c.Response().Writer, "data: nothing to sent")
		}

		if f, ok := c.(http.Flusher); ok {
			fmt.Println("flush")
			f.Flush()
		}

		//f, ok := c.(http.Flusher)
		//if !ok {
		//	log.Error("streaming unsupported")
		//	//return c.JSON(http.StatusInternalServerError, "streaming unsupported")
		//	return errors.New("streaming unsupported")
		//}
		//f.Flush()

		return nil
	}
}

func addSSEHeader(c echo.Context) echo.Context {
	c.Response().Header().Add("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Add("Content-Type", "text/event-stream")
	c.Response().Header().Add("Cache-Control", "no-cache")
	c.Response().Header().Add("Connection", "keep-alive")

	return c
}
