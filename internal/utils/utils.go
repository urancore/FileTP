package utils

import (
	"fmt"
	"math"
	"net"
	"net/http"
	"strings"
)

func GetUserIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// Берем первый IP из списка (клиентский адрес)
		ips := strings.Split(ip, ", ")
		if len(ips) > 0 {
			ip = ips[0]
		}
		return ip
	}

	ip = r.Header.Get("X-Real-IP")

	if ip != "" {
		return ip
	}

	// Если заголовков нет, используем RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if ip == "::1" {
		return "localhost"
	}

	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func FormatFileSize(size int64) string {
	if size == 0 {
		return "0B"
	}

	sizes := []string{"B", "K", "M", "G", "T", "P", "E"}
	base := 1000.0

	order := int(math.Log(float64(size)) / math.Log(base))
	if order > len(sizes)-1 {
		order = len(sizes) - 1
	} else if order < 0 {
		order = 0
	}

	value := float64(size) / math.Pow(base, float64(order))
	return fmt.Sprintf("%.1f%s", value, sizes[order])
}
