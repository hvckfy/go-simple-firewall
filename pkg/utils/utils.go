package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetClientIP извлекает IP клиента из HTTP запроса
func GetClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	return strings.Split(r.RemoteAddr, ":")[0]
}

// InstallService устанавливает firewall как системный сервис
func InstallService() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	serviceName := "go-firewall"
	serviceContent := fmt.Sprintf(`[Unit]
Description=Go Simple Firewall
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=%s
ExecStart=%s
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
`, filepath.Dir(execPath), execPath)

	serviceFile := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)

	if err := os.WriteFile(serviceFile, []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("failed to write service file: %v", err)
	}

	// Перезагружаем systemd и включаем сервис
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		log.Printf("Warning: failed to reload systemd: %v", err)
	}

	if err := exec.Command("systemctl", "enable", serviceName).Run(); err != nil {
		log.Printf("Warning: failed to enable service: %v", err)
	}

	log.Printf("Service installed: %s", serviceName)
	return nil
}

// IsValidIP проверяет, является ли строка валидным IP адресом
func IsValidIP(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}

		num := 0
		for _, char := range part {
			if char < '0' || char > '9' {
				return false
			}
			num = num*10 + int(char-'0')
		}

		if num > 255 {
			return false
		}
	}

	return true
}

// Contains проверяет, содержит ли слайс строку
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
