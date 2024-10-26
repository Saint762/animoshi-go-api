package utils

import (
	"net"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetUserIP(c echo.Context) string {
	// Check for X-Forwarded-For header (in case of proxy)
	xForwardedFor := c.Request().Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// Split the list and find the first IPv4 address
		ips := strings.Split(xForwardedFor, ",")
		for _, ip := range ips {
			trimmedIP := strings.TrimSpace(ip)
			if isIPv4(trimmedIP) {
				return trimmedIP
			}
		}
	}

	// Fallback to RemoteAddr and check if it’s IPv4
	ip, _, err := net.SplitHostPort(c.Request().RemoteAddr)
	if err == nil && isIPv4(ip) {
		return ip
	}

	// If no IPv4 address found, return an empty string or error message
	return ""
}

// Helper function to check if an IP is IPv4
func isIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}
