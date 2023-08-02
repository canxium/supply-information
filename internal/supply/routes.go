package supply

import (
	"github.com/labstack/echo/v4"
)

func MapRoutes(group *echo.Group, h *handlers) {
	group.GET("/info/cau", h.GetSupplyInfo)
}
