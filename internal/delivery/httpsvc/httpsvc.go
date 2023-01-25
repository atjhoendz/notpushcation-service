package httpsvc

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type service struct {
	group *echo.Group
}

// RouteService :nodoc:
func RouteService(group *echo.Group) {
	srv := &service{
		group: group,
	}
	srv.initRoutes()
}

func (s *service) initRoutes() {
	s.group.GET("/info", func(c echo.Context) error {
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
