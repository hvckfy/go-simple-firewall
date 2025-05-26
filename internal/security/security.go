package security

import (
	"net/http"
	"strings"
	"time"

	"go-simple-firewall/internal/config"
	"go-simple-firewall/pkg/utils"
)

// SecurityChecker проверяет различные угрозы безопасности
type SecurityChecker struct {
	config *config.Config
}

// New создает новый SecurityChecker
func New(cfg *config.Config) *SecurityChecker {
	return &SecurityChecker{
		config: cfg,
	}
}

// CheckRequest проверяет запрос на различные угрозы
func (sc *SecurityChecker) CheckRequest(r *http.Request) (blocked bool, reason string) {
	clientIP := utils.GetClientIP(r)

	// Проверяем временные баны
	if banned, banReason := sc.config.IsTemporarilyBanned(clientIP); banned {
		return true, "Temporarily banned: " + banReason
	}

	// Проверка запрещенных суффиксов
	if sc.config.Security.EnableSuffixProtection {
		if blocked, reason := sc.checkForbiddenSuffixes(r); blocked {
			// Добавляем временный бан
			duration := time.Duration(sc.config.Security.SuffixBanDuration) * time.Hour
			sc.config.AddTemporaryBan(clientIP, reason, duration)
			return true, reason
		}
	}

	// Проверка SQL инъекций
	if sc.config.Security.EnableSQLProtection {
		if blocked, reason := sc.checkSQLInjection(r); blocked {
			return true, reason
		}
	}

	// Проверка XSS
	if sc.config.Security.EnableXSSProtection {
		if blocked, reason := sc.checkXSS(r); blocked {
			return true, reason
		}
	}

	// Проверка сканеров
	if sc.config.Security.EnableScannerProtection {
		if blocked, reason := sc.checkScanner(r); blocked {
			return true, reason
		}
	}

	// Проверка ботов
	if sc.config.Security.EnableBotProtection {
		if blocked, reason := sc.checkBot(r); blocked {
			return true, reason
		}
	}

	// Проверка директорий
	if sc.config.Security.EnableDirectoryProtection {
		if blocked, reason := sc.checkDirectory(r); blocked {
			return true, reason
		}
	}

	return false, ""
}

// checkForbiddenSuffixes проверяет запрещенные суффиксы
func (sc *SecurityChecker) checkForbiddenSuffixes(r *http.Request) (bool, string) {
	path := strings.ToLower(r.URL.Path)

	for _, suffix := range sc.config.Security.ForbiddenSuffixes {
		if strings.Contains(path, strings.ToLower(suffix)) {
			return true, "Forbidden suffix detected: " + suffix
		}
	}

	return false, ""
}

// checkSQLInjection проверяет SQL инъекции
func (sc *SecurityChecker) checkSQLInjection(r *http.Request) (bool, string) {
	// Проверяем URL параметры
	queryParams := r.URL.Query()
	for _, values := range queryParams {
		for _, value := range values {
			if sc.containsSQLKeywords(value) {
				return true, "SQL injection attempt detected"
			}
		}
	}

	// Проверяем POST данные
	if r.Method == "POST" {
		r.ParseForm()
		for _, values := range r.PostForm {
			for _, value := range values {
				if sc.containsSQLKeywords(value) {
					return true, "SQL injection attempt detected"
				}
			}
		}
	}

	return false, ""
}

// containsSQLKeywords проверяет наличие SQL ключевых слов
func (sc *SecurityChecker) containsSQLKeywords(input string) bool {
	input = strings.ToLower(input)

	for _, keyword := range sc.config.Security.SQLKeywords {
		if strings.Contains(input, strings.ToLower(keyword)) {
			return true
		}
	}

	return false
}

// checkXSS проверяет XSS атаки
func (sc *SecurityChecker) checkXSS(r *http.Request) (bool, string) {
	// Проверяем URL параметры
	queryParams := r.URL.Query()
	for _, values := range queryParams {
		for _, value := range values {
			if sc.containsXSSPatterns(value) {
				return true, "XSS attempt detected"
			}
		}
	}

	// Проверяем POST данные
	if r.Method == "POST" {
		r.ParseForm()
		for _, values := range r.PostForm {
			for _, value := range values {
				if sc.containsXSSPatterns(value) {
					return true, "XSS attempt detected"
				}
			}
		}
	}

	return false, ""
}

// containsXSSPatterns проверяет наличие XSS паттернов
func (sc *SecurityChecker) containsXSSPatterns(input string) bool {
	input = strings.ToLower(input)

	for _, pattern := range sc.config.Security.XSSPatterns {
		if strings.Contains(input, strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}

// checkScanner проверяет попытки сканирования
func (sc *SecurityChecker) checkScanner(r *http.Request) (bool, string) {
	path := strings.ToLower(r.URL.Path)

	for _, scannerPath := range sc.config.Security.ScannerPaths {
		if strings.HasPrefix(path, strings.ToLower(scannerPath)) {
			return true, "Scanner attempt detected: " + scannerPath
		}
	}

	return false, ""
}

// checkBot проверяет подозрительных ботов
func (sc *SecurityChecker) checkBot(r *http.Request) (bool, string) {
	userAgent := strings.ToLower(r.UserAgent())

	for _, suspiciousUA := range sc.config.Security.SuspiciousUserAgents {
		if strings.Contains(userAgent, strings.ToLower(suspiciousUA)) {
			return true, "Suspicious bot detected: " + suspiciousUA
		}
	}

	return false, ""
}

// checkDirectory проверяет доступ к защищенным директориям
func (sc *SecurityChecker) checkDirectory(r *http.Request) (bool, string) {
	path := strings.ToLower(r.URL.Path)

	for _, protectedDir := range sc.config.Security.ProtectedDirectories {
		if strings.HasPrefix(path, strings.ToLower(protectedDir)) {
			return true, "Access to protected directory: " + protectedDir
		}
	}

	return false, ""
}

// UpdateConfig обновляет конфигурацию
func (sc *SecurityChecker) UpdateConfig(cfg *config.Config) {
	sc.config = cfg
}
