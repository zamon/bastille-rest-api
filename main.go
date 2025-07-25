package main

import (
	"bastille-rest-api/routes"
	"fmt"
	"log"
	"os"
	"time"
	"strings"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("APP_PORT")

	// fallback default
	if port == "" {
		port = "8080"
	}

	e := echo.New()
	allowedIP := os.Getenv("IP_WHITELIST")
	e.Use(IPWhitelistMiddleware(allowedIP))

	routes.Init(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

func IPWhitelistMiddleware(allowedIP string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			if strings.Contains(ip, ":") {
				ip = strings.Split(ip, ":")[0]
			}

			if ip != allowedIP {
				time.Sleep(2 * time.Second)
				return nil
			}

			return next(c)
		}
	}
}