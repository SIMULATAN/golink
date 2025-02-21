package main

import (
	"context"
	_ "embed"
	"github.com/labstack/echo/v4"
	"golink/config"
	"golink/frontend/home"
	"golink/route"
	"golink/service"
	"golink/service/link"
	"log"
	"time"
)

var (
	//go:embed sql/schema.sql
	DatabaseSchema string
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error while loading config!", err)
	}

	app := echo.New()
	app.Debug = true

	app.HTTPErrorHandler = func(err error, c echo.Context) {
		route.HandleError(app, err, c)
	}

	var linkService service.LinkService
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	linkService, err = link.NewPostgresService(ctx, cfg.Postgres, DatabaseSchema)
	if err != nil {
		log.Fatalln("Error while connecting to postgres!", err)
	}
	linkService.Init()

	app.GET("/", func(c echo.Context) error {
		return route.Render(c, home.Home(cfg.BaseUrl))
	})

	linkHandler := route.LinkHandler{
		LinkService: linkService,
		BaseUrl:     cfg.BaseUrl,
	}

	app.POST("/link", linkHandler.CreateLink)
	app.GET("/to/:code", linkHandler.RedirectLink)

	app.Static("/", "frontend/assets")

	app.Logger.Fatal(app.Start(":8080"))
}
