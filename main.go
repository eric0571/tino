package main

import (
	"net/http"

	"tino/util/system"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	system.JSONPrint("Ts", "22")
	e.Logger.Fatal(e.Start(":1323"))
}
