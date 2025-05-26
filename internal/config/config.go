package config

import (
	"encoding/json"
	"os"
	"time"
)

const ConfigFile = "firewall.json"

// SecuritySettings настройки безопасности
type SecuritySettings struct {
	// Защита от запрещенных суффиксов
	EnableSuffixProtection bool     `json:"enable_suffix_protection"`
	ForbiddenSuffixes      []string `json:"forbidden_suffixes"`
	SuffixBanDuration      int      `json:"suffix_ban_duration"` // в часах
	
	// Защита от SQL инъекций
	EnableSQLProtection bool     `json:"enable_sql_protection"`
	SQLKeywords         []string `json:"sql_keywords"`
	
	// Защита от XSS
	EnableXSSProtection bool     `json:"enable_xss_protection"`
	XSSPatterns         []string `json:"xss_patterns"`
	
	// Защита от сканеров
	EnableScannerProtection bool     `json:"enable_scanner_protection"`
	ScannerPaths            []string `json:"scanner_paths"`
	
	// Защита от ботов
	EnableBotProtection bool     `json:"enable_bot_protection"`
	SuspiciousUserAgents []string `json:"suspicious_user_agents"`
	
	// Защита от директорий
	EnableDirectoryProtection bool     `json:"enable_directory_protection"`
	ProtectedDirectories      []string `json:"protected_directories"`
	
	// DDoS защита
	EnableDDoSProtection bool `json:"enable_ddos_protection"`
	DDoSThreshold        int  `json:"ddos_threshold"`        // запросов за период
	DDoSTimeWindow       int  `json:"ddos_time_window"`      // период в секундах
	DDoSBanDuration      int  `json:"ddos_ban_duration"`     // бан в минутах
	
	// Геоблокировка
	EnableGeoBlocking bool     `json:"enable_geo_blocking"`
	BlockedCountries  []string `json:"blocked_countries"`
}

// TemporaryBan временный бан
type TemporaryBan struct {
	IP        string    `json:"ip"`
	Reason    string    `json:"reason"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Config структура конфигурации firewall
type Config struct {
	ListenPort     int               `json:"listen_port"`
	AdminPort      int               `json:"admin_port"`
	TargetPort     int               `json:"target_port"`
	RateLimitRPS   int               `json:"rate_limit_rps"`
	EnableLogging  bool              `json:"enable_logging"`
	EnableFirewall bool              `json:"enable_firewall"`
	BannedIPs      map[string]bool   `json:"banned_ips"`
	AllowedIPs     map[string]bool   `json:"allowed_ips"`
	AllowedUAs     []string          `json:"allowed_user_agents"`
	AutoStart      bool              `json:"auto_start"`
	
	// Аутентификация
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	IsSetup      bool   `json:"is_setup"`
	
	// Настройки безопасности
	Security SecuritySettings `json:"security"`
	
	// Временные баны
	TemporaryBans []TemporaryBan `json:"temporary_bans"`
}

// Load загружает конфигурацию из файла
func Load() (*Config, error) {
	config := GetDefault()

	if data, err := os.ReadFile(ConfigFile); err == nil {
		if err := json.Unmarshal(data, config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// Save сохраняет конфигурацию в файл
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigFile, data, 0644)
}

// GetDefault возвращает конфигурацию по умолчанию
func GetDefault() *Config {
	return &Config{
		ListenPort:     8080,
		AdminPort:      9090,
		TargetPort:     3000,
		RateLimitRPS:   60,
		EnableLogging:  true,
		EnableFirewall: true,
		BannedIPs:      make(map[string]bool),
		AllowedIPs:     make(map[string]bool),
		AllowedUAs:     []string{},
		AutoStart:      false,
		Username:       "",
		PasswordHash:   "",
		IsSetup:        false,
		Security: SecuritySettings{
			EnableSuffixProtection: true,
			ForbiddenSuffixes:      []string{".php", ".asp", ".aspx", ".jsp", ".cgi"},
			SuffixBanDuration:      10,
			
			EnableSQLProtection: true,
			SQLKeywords:         []string{"union", "select", "insert", "delete", "update", "drop", "create", "alter", "exec", "script"},
			
			EnableXSSProtection: true,
			XSSPatterns:         []string{"<script", "javascript:", "onload=", "onerror=", "onclick=", "onmouseover="},
			
			EnableScannerProtection: true,
			ScannerPaths:            []string{"/admin", "/wp-admin", "/phpmyadmin", "/cpanel", "/webmail", "/.env", "/config"},
			
			EnableBotProtection: true,
			SuspiciousUserAgents:    []string{"bot", "crawler", "spider", "scraper", "scanner", "nikto", "sqlmap"},
			
			EnableDirectoryProtection: true,
			ProtectedDirectories:      []string{"/.git", "/.svn", "/backup", "/config", "/logs", "/tmp"},
			
			EnableDDoSProtection: true,
			DDoSThreshold:        100,
			DDoSTimeWindow:       60,
			DDoSBanDuration:      30,
			
			EnableGeoBlocking: false,
			BlockedCountries:  []string{},
		},
		TemporaryBans: []TemporaryBan{},
	}
}

// AddTemporaryBan добавляет временный бан
func (c *Config) AddTemporaryBan(ip, reason string, duration time.Duration) {
	// Удаляем существующий бан для этого IP
	c.RemoveTemporaryBan(ip)
	
	ban := TemporaryBan{
		IP:        ip,
		Reason:    reason,
		ExpiresAt: time.Now().Add(duration),
	}
	
	c.TemporaryBans = append(c.TemporaryBans, ban)
}

// RemoveTemporaryBan удаляет временный бан
func (c *Config) RemoveTemporaryBan(ip string) {
	for i, ban := range c.TemporaryBans {
		if ban.IP == ip {
			c.TemporaryBans = append(c.TemporaryBans[:i], c.TemporaryBans[i+1:]...)
			break
		}
	}
}

// IsTemporarilyBanned проверяет, временно ли заблокирован IP
func (c *Config) IsTemporarilyBanned(ip string) (bool, string) {
	now := time.Now()
	
	for i := len(c.TemporaryBans) - 1; i >= 0; i-- {
		ban := c.TemporaryBans[i]
		
		if ban.IP == ip {
			if now.After(ban.ExpiresAt) {
				// Бан истек, удаляем его
				c.TemporaryBans = append(c.TemporaryBans[:i], c.TemporaryBans[i+1:]...)
				continue
			}
			return true, ban.Reason
		}
	}
	
	return false, ""
}

// CleanupExpiredBans очищает истекшие баны
func (c *Config) CleanupExpiredBans() {
	now := time.Now()
	validBans := []TemporaryBan{}
	
	for _, ban := range c.TemporaryBans {
		if now.Before(ban.ExpiresAt) {
			validBans = append(validBans, ban)
		}
	}
	
	c.TemporaryBans = validBans
}
