package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Root(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, e.Reverse("list-vampires"))
	}).Name = "root"
}
