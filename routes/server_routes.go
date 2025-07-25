// Copyright (c) 2025, Zainal Abidin
// SPDX-License-Identifier: BSD-2-Clause

package routes

import (
	"bastille-rest-api/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterServerRoutes(e *echo.Echo) {
	e.GET("/", controllers.Ping)
	e.GET("/ping", controllers.Ping)
	e.GET("/get-server-spec", controllers.GetServerSpec)
	e.GET("/activate-website", controllers.ActivateWebsite)
	e.GET("/deactivate-website", controllers.DeactivateWebsite)
	e.GET("/suspend-website", controllers.SuspendWebsite)
	e.POST("/bastille-pkg-list", controllers.BastillePkgList)
	e.GET("/bastille-list-all", controllers.BastilleListAll)
	e.GET("/bootstrap-list", controllers.BootstrapList)
	e.POST("/create-jail", controllers.CreateJail)
	e.POST("/set-jail-quota", controllers.SetJailQuota)
	e.POST("/stop-jail", controllers.StopJail)
	e.POST("/destroy-jail", controllers.DestroyJail)
	e.POST("/clone-jail", controllers.CloneJail)
}