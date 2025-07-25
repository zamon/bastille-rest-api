// Copyright (c) 2025, Zainal Abidin
// SPDX-License-Identifier: BSD-2-Clause

package routes

import (
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {
    RegisterServerRoutes(e)
}