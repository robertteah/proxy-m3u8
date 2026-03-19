package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/dovakiin0/proxy-m3u8/config"
	"github.com/dovakiin0/proxy-m3u8/internal/handler"
	mdlware "github.com/dovakiin0/proxy-m3u8/internal/middleware"
)

func init() {
	godotenv.Load()
	config.InitConfig()
	config.RedisConnect()
}

func main() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Response().Header()
			h.Set("Access-Control-Allow-Origin", "*")
			h.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			h.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Range")
			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}
			return next(c)
		}
	})

	customCacheConfig := mdlware.CacheControlConfig{
		MaxAge:         3600, // 1 hour
		Public:         true,
		MustRevalidate: true,
	}
	e.Use(mdlware.CacheControlWithConfig(customCacheConfig))
	e.GET("/m3u8-proxy", handler.M3U8ProxyHandler)

	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	port := config.Env.Port

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
