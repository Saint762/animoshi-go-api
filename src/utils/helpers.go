package utils

import (
	"github.com/labstack/echo/v4"
	"strings"
)

func GetUserIP(c echo.Context) string {
	// Get the IP address from the X-Forwarded-For header if your app is behind a proxy
	xForwardedFor := c.Request().Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// Split the list and return the first IP address
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Fallback to RemoteAddr if X-Forwarded-For is not present
	return c.Request().RemoteAddr
}
