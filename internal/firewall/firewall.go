package firewall

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"go-simple-firewall/internal/admin"
	"go-simple-firewall/internal/config"
	"go-simple-firewall/internal/ddos"
	"go-simple-firewall/internal/logger"
	"go-simple-firewall/internal/ratelimit"
	"go-simple-firewall/internal/security"
	"go-simple-firewall/internal/stats"
	"go-simple-firewall/pkg/utils"
)

// Firewall основная структура firewall
type Firewall struct {
	config           *config.Config
	rateLimiter      *ratelimit.Limiter
	proxy            *httputil.ReverseProxy
	logger           *logger.Logger
	adminPanel       *admin.Handler
	firewallServer   *http.Server
	adminServer      *http.Server
	stats            *stats.Stats
	securityChecker  *security.SecurityChecker
	ddosProtection   *ddos.DDoSProtection
	mutex            sync.RWMutex
}

// New создает новый firewall
func New(cfg *config.Config) (*Firewall, error) {
	// Создаем logger
	log, err := logger.New(cfg.EnableLogging)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %v", err)
	}

	// Создаем rate limiter
	rateLimiter := ratelimit.New(cfg.RateLimitRPS)

	// Создаем proxy
	targetURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", cfg.TargetPort))
	if err != nil {
		return nil, fmt.Errorf("invalid target URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Создаем статистику
	statistics := stats.New()

	// Создаем security checker
	securityChecker := security.New(cfg)

	// Создаем DDoS защиту
	ddosProtection := ddos.New(cfg)

	// Создаем admin panel
	adminPanel := admin.New(cfg, log, statistics)

	fw := &Firewall{
		config:          cfg,
		rateLimiter:     rateLimiter,
		proxy:           proxy,
		logger:          log,
		adminPanel:      adminPanel,
		stats:           statistics,
		securityChecker: securityChecker,
		ddosProtection:  ddosProtection,
	}

	// Запускаем очистку истекших банов каждые 10 минут
	go fw.cleanupExpiredBans()

	return fw, nil
}

// cleanupExpiredBans очищает истекшие баны
func (fw *Firewall) cleanupExpiredBans() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			fw.mutex.Lock()
			oldCount := len(fw.config.TemporaryBans)
			fw.config.CleanupExpiredBans()
			newCount := len(fw.config.TemporaryBans)
			
			if oldCount != newCount {
				fw.logger.LogInfo(fmt.Sprintf("Cleaned up %d expired temporary bans", oldCount-newCount))
				fw.config.Save()
			}
			fw.mutex.Unlock()
		}
	}
}

// Start запускает firewall и админ-панель на разных портах
func (fw *Firewall) Start() error {
	// Запускаем админ-панель в отдельной горутине
	go fw.startAdminServer()

	// Запускаем основной firewall сервер
	return fw.startFirewallServer()
}

// startFirewallServer запускает основной firewall сервер
func (fw *Firewall) startFirewallServer() error {
	fw.firewallServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", fw.config.ListenPort),
		Handler: http.HandlerFunc(fw.ServeHTTP),
	}

	fw.logger.LogInfo(fmt.Sprintf("Firewall started on port %d, proxying to %d", fw.config.ListenPort, fw.config.TargetPort))
	fmt.Printf("🔥 Firewall started on port %d, proxying to %d\n", fw.config.ListenPort, fw.config.TargetPort)
	fmt.Printf("📊 Admin panel: http://localhost:%d/admin\n", fw.config.AdminPort)

	return fw.firewallServer.ListenAndServe()
}

// startAdminServer запускает админ-панель
func (fw *Firewall) startAdminServer() error {
	mux := http.NewServeMux()
	
	// Админ панель и API
	mux.HandleFunc("/admin", fw.adminPanel.ServeHTTP)
	mux.HandleFunc("/admin/", fw.adminPanel.ServeHTTP)
	mux.HandleFunc("/admin/api/", fw.adminPanel.ServeHTTP)

	fw.adminServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", fw.config.AdminPort),
		Handler: mux,
	}

	fw.logger.LogInfo(fmt.Sprintf("Admin panel started on port %d", fw.config.AdminPort))
	fmt.Printf("📊 Admin panel started on port %d\n", fw.config.AdminPort)

	return fw.adminServer.ListenAndServe()
}

// ServeHTTP обрабатывает HTTP запросы firewall
func (fw *Firewall) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientIP := utils.GetClientIP(r)
	userAgent := r.UserAgent()
	
	// Проверка на блокировку
	if blocked, reason := fw.isBlocked(r); blocked {
		// Логируем заблокированный запрос с детальной информацией
		fw.logger.LogRequestWithReason(r, "BLOCKED", 403, reason)
		
		// Логируем как атаку если это серьезная угроза
		if fw.isSecurityThreat(reason) {
			fw.logger.LogAttack(clientIP, fw.getAttackType(reason), reason, userAgent, r.URL.String())
		}
		
		fw.stats.RecordRequest(clientIP, userAgent, true)
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// Логируем разрешенный запрос
	fw.logger.LogRequest(r, "ALLOWED", 200)
	fw.stats.RecordRequest(clientIP, userAgent, false)
	
	// Проксирование запроса
	fw.proxy.ServeHTTP(w, r)
}

// isSecurityThreat определяет, является ли причина блокировки угрозой безопасности
func (fw *Firewall) isSecurityThreat(reason string) bool {
	threats := []string{
		"SQL injection", "XSS attempt", "Scanner attempt", 
		"Forbidden suffix", "DDoS attack", "Suspicious bot",
		"Protected directory",
	}
	
	reasonLower := strings.ToLower(reason)
	for _, threat := range threats {
		if strings.Contains(reasonLower, strings.ToLower(threat)) {
			return true
		}
	}
	
	return false
}

// getAttackType определяет тип атаки по причине блокировки
func (fw *Firewall) getAttackType(reason string) string {
	reasonLower := strings.ToLower(reason)
	
	if strings.Contains(reasonLower, "sql") {
		return "SQL_INJECTION"
	}
	if strings.Contains(reasonLower, "xss") {
		return "XSS"
	}
	if strings.Contains(reasonLower, "scanner") {
		return "SCANNER"
	}
	if strings.Contains(reasonLower, "ddos") {
		return "DDOS"
	}
	if strings.Contains(reasonLower, "bot") {
		return "BOT"
	}
	if strings.Contains(reasonLower, "suffix") {
		return "MALICIOUS_FILE"
	}
	if strings.Contains(reasonLower, "directory") {
		return "DIRECTORY_TRAVERSAL"
	}
	
	return "UNKNOWN"
}

// isBlocked проверяет, заблокирован ли запрос
func (fw *Firewall) isBlocked(r *http.Request) (bool, string) {
	fw.mutex.RLock()
	defer fw.mutex.RUnlock()

	if !fw.config.EnableFirewall {
		return false, ""
	}

	clientIP := utils.GetClientIP(r)
	userAgent := strings.ToLower(r.UserAgent())

	// Проверка забаненных IP
	if fw.config.BannedIPs[clientIP] {
		return true, "IP permanently banned"
	}

	// Проверка разрешенных IP (если список не пуст)
	if len(fw.config.AllowedIPs) > 0 && !fw.config.AllowedIPs[clientIP] {
		return true, "IP not in whitelist"
	}

	// Проверка User-Agent whitelist (если список не пуст, то проверяем whitelist)
	if len(fw.config.AllowedUAs) > 0 {
		allowed := false
		for _, allowedUA := range fw.config.AllowedUAs {
			if strings.Contains(userAgent, strings.ToLower(allowedUA)) {
				allowed = true
				break
			}
		}
		if !allowed {
			return true, "User-Agent not in whitelist"
		}
	}

	// Проверка rate limit
	if !fw.rateLimiter.IsAllowed(clientIP) {
		return true, "Rate limit exceeded"
	}

	// DDoS защита
	if blocked, reason := fw.ddosProtection.CheckRequest(clientIP); blocked {
		// Логируем временный бан
		fw.logger.LogTemporaryBan(clientIP, reason, time.Duration(fw.config.Security.DDoSBanDuration)*time.Minute)
		return true, reason
	}

	// Проверки безопасности
	if blocked, reason := fw.securityChecker.CheckRequest(r); blocked {
		// Если это суффикс, логируем временный бан
		if strings.Contains(reason, "Forbidden suffix") {
			duration := time.Duration(fw.config.Security.SuffixBanDuration) * time.Hour
			fw.logger.LogTemporaryBan(clientIP, reason, duration)
		}
		return true, reason
	}

	return false, ""
}

// UpdateConfig обновляет конфигурацию
func (fw *Firewall) UpdateConfig(newConfig *config.Config) error {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	// Обновляем rate limiter если изменился лимит
	if fw.config.RateLimitRPS != newConfig.RateLimitRPS {
		fw.rateLimiter.UpdateLimit(newConfig.RateLimitRPS)
	}

	// Обновляем proxy если изменился target port
	if fw.config.TargetPort != newConfig.TargetPort {
		targetURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", newConfig.TargetPort))
		if err != nil {
			return fmt.Errorf("invalid target URL: %v", err)
		}
		fw.proxy = httputil.NewSingleHostReverseProxy(targetURL)
	}

	// Обновляем logger если изменились настройки логирования
	if fw.config.EnableLogging != newConfig.EnableLogging {
		fw.logger.SetEnabled(newConfig.EnableLogging)
	}

	// Обновляем security checker и ddos protection
	fw.securityChecker.UpdateConfig(newConfig)
	fw.ddosProtection.UpdateConfig(newConfig)

	fw.config = newConfig
	return fw.config.Save()
}

// Shutdown корректно завершает работу firewall
func (fw *Firewall) Shutdown() {
	fw.logger.LogInfo("Firewall shutting down")
	
	if fw.firewallServer != nil {
		fw.firewallServer.Close()
	}

	if fw.adminServer != nil {
		fw.adminServer.Close()
	}

	if fw.logger != nil {
		fw.logger.Close()
	}

	if fw.config != nil {
		fw.config.Save()
	}
}

// GetStats возвращает статистику firewall
func (fw *Firewall) GetStats() map[string]interface{} {
	fw.mutex.RLock()
	defer fw.mutex.RUnlock()

	return map[string]interface{}{
		"firewall_enabled":   fw.config.EnableFirewall,
		"logging_enabled":    fw.config.EnableLogging,
		"rate_limit":         fw.config.RateLimitRPS,
		"banned_ips":         len(fw.config.BannedIPs),
		"allowed_ips":        len(fw.config.AllowedIPs),
		"allowed_uas":        len(fw.config.AllowedUAs),
		"temporary_bans":     len(fw.config.TemporaryBans),
		"rate_stats":         fw.rateLimiter.GetStats(),
		"log_stats":          fw.logger.GetLogStats(),
	}
}
