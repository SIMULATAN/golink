package main

import (
	"github.com/labstack/echo/v4"
	"golink/frontend/home"
	"golink/route"
	"golink/service"
)

func main() {
	app := echo.New()
	app.Debug = true

	app.HTTPErrorHandler = func(err error, c echo.Context) {
		route.HandleError(app, err, c)
	}

	linkService := service.LinkService{}
	linkService.Init()

	app.GET("/", func(c echo.Context) error {
		return route.Render(c, home.Home())
	})

	linkHandler := route.LinkHandler{
		LinkService: linkService,
		BaseUrl:     "localhost:8080",
	}

	app.POST("/link", linkHandler.CreateLink)
	app.GET("/to/:code", linkHandler.RedirectLink)

	app.Logger.Fatal(app.Start(":8080"))
}
