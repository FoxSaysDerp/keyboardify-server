package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"

	"keyboardify-server/auth"
	"keyboardify-server/controllers"
	"keyboardify-server/db"
	mdlwr "keyboardify-server/middleware"
	"keyboardify-server/route"
)

func main() {
	godotenv.Load()

	e := echo.New()
	firebaseAuth := auth.InitFirebase()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", controllers.Db)
			c.Set("firebaseAuth", firebaseAuth)
			return next(c)
		}
	})

	logger := zerolog.New(os.Stdout)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	db.Init()

	e.Static("/api/public", "public")
	route.Init(e.Group("/api", mdlwr.Auth))

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
