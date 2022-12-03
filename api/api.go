package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mct-joken/jkojs-agent/types"
)

func StartServer() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/run", handler)
	e.Logger.Fatal(e.Start(":3030"))
}

func handler(c echo.Context) error {
	req := types.StartExecRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	res := types.StartExecResponse{}
	return c.JSON(http.StatusAccepted, res)
}
