package main

import (
	"net/http"

	"github.com/eric0571/tino/util/config"
	"github.com/eric0571/tino/util/system"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	system.JSONPrint("Ts", "22")
	config.HelloWorld()
	e.Logger.Fatal(e.Start(":1323"))
}
