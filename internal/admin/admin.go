package admin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go-simple-firewall/internal/auth"
	"go-simple-firewall/internal/config"
	"go-simple-firewall/internal/logger"
	"go-simple-firewall/internal/stats"
	"go-simple-firewall/pkg/service"
)

// Handler структура для обработки админ панели
type Handler struct {
	config         *config.Config
	logger         *logger.Logger
	stats          *stats.Stats
	template       *template.Template
	loginTemplate  *template.Template
	setupTemplate  *template.Template
	sessionManager *auth.SessionManager
}

// TemplateData структура для передачи данных в шаблон
type TemplateData struct {
	StatusClass      string
	FirewallStatus   string
	LoggingStatus    string
	RateLimitRPS     int
	ListenPort       int
	AdminPort        int
	TargetPort       int
	ServiceStatus    string
	FirewallChecked  string
	LoggingChecked   string
	BannedIPs        string
	AllowedUAs       string
	Logs             string
	Username         string
	
	// Security settings
	SuffixProtectionChecked    string
	SQLProtectionChecked       string
	XSSProtectionChecked       string
	ScannerProtectionChecked   string
	BotProtectionChecked       string
	DirectoryProtectionChecked string
	DDoSProtectionChecked      string
	GeoBlockingChecked         string
	
	ForbiddenSuffixes     string
	SuffixBanDuration     int
	SQLKeywords           string
	XSSPatterns           string
	ScannerPaths          string
	SuspiciousUserAgents  string
	ProtectedDirectories  string
	DDoSThreshold         int
	DDoSTimeWindow        int
	DDoSBanDuration       int
	BlockedCountries      string
	
	TemporaryBans         []config.TemporaryBan
	LogStats              map[string]interface{}
}

// New создает новый admin handler
func New(cfg *config.Config, log *logger.Logger, st *stats.Stats) *Handler {
	tmpl := template.Must(template.New("admin").Parse(getHTMLTemplate()))
	loginTmpl := template.Must(template.New("login").Parse(getLoginTemplate()))
	setupTmpl := template.Must(template.New("setup").Parse(getSetupTemplate()))
	
	sessionManager := auth.NewSessionManager()
	
	// Запускаем очистку истекших сессий каждый час
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				sessionManager.CleanupExpiredSessions()
			}
		}
	}()
	
	return &Handler{
		config:         cfg,
		logger:         log,
		stats:          st,
		template:       tmpl,
		loginTemplate:  loginTmpl,
		setupTemplate:  setupTmpl,
		sessionManager: sessionManager,
	}
}

// ServeHTTP обрабатывает запросы к админ панели
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Если система не настроена, показываем форму первоначальной настройки
	if !h.config.IsSetup {
		h.handleSetup(w, r)
		return
	}
	
	// Обрабатываем логин/логаут
	if r.URL.Path == "/admin/login" {
		h.handleLogin(w, r)
		return
	}
	
	if r.URL.Path == "/admin/logout" {
		h.handleLogout(w, r)
		return
	}
	
	// Проверяем аутентификацию для всех остальных запросов
	if !h.isAuthenticated(r) {
		h.redirectToLogin(w, r)
		return
	}
	
	// API endpoints
	if strings.HasPrefix(r.URL.Path, "/admin/api/") {
		h.handleStatsAPI(w, r)
		return
	}
	
	if r.Method == "GET" {
		h.showAdminPanel(w)
		return
	}

	if r.Method == "POST" {
		h.handleAdminAction(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleSetup обрабатывает первоначальную настройку
func (h *Handler) handleSetup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := h.setupTemplate.Execute(w, nil); err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	
	if r.Method == "POST" {
		r.ParseForm()
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		
		// Валидация
		if username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}
		
		if len(password) < 6 {
			http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
			return
		}
		
		if password != confirmPassword {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}
		
		// Хешируем пароль
		passwordHash, err := auth.HashPassword(password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		
		// Сохраняем в конфиг
		h.config.Username = username
		h.config.PasswordHash = passwordHash
		h.config.IsSetup = true
		
		if err := h.config.Save(); err != nil {
			http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
			return
		}
		
		h.logger.LogInfo("Admin user created: " + username)
		
		// Создаем сессию и перенаправляем в админ панель
		token, err := h.sessionManager.CreateSession()
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}
		
		auth.SetSessionCookie(w, token)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleLogin обрабатывает логин
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := h.loginTemplate.Execute(w, nil); err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	
	if r.Method == "POST" {
		r.ParseForm()
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")
		
		// Проверяем учетные данные
		if username != h.config.Username || !auth.CheckPassword(password, h.config.PasswordHash) {
			time.Sleep(1 * time.Second) // Защита от брутфорса
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		
		// Создаем сессию
		token, err := h.sessionManager.CreateSession()
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}
		
		auth.SetSessionCookie(w, token)
		h.logger.LogInfo("Admin logged in: " + username)
		
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleLogout обрабатывает логаут
func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	if token, err := auth.GetSessionCookie(r); err == nil {
		h.sessionManager.DeleteSession(token)
	}
	
	auth.ClearSessionCookie(w)
	h.logger.LogInfo("Admin logged out")
	
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// isAuthenticated проверяет аутентификацию пользователя
func (h *Handler) isAuthenticated(r *http.Request) bool {
	token, err := auth.GetSessionCookie(r)
	if err != nil {
		return false
	}
	
	return h.sessionManager.ValidateSession(token)
}

// redirectToLogin перенаправляет на страницу логина
func (h *Handler) redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// showAdminPanel отображает админ панель
func (h *Handler) showAdminPanel(w http.ResponseWriter) {
	data := h.prepareTemplateData()
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.template.Execute(w, data); err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// prepareTemplateData подготавливает данные для шаблона
func (h *Handler) prepareTemplateData() *TemplateData {
	statusClass := "enabled"
	if !h.config.EnableFirewall {
		statusClass = "disabled"
	}

	firewallStatus := "Enabled"
	if !h.config.EnableFirewall {
		firewallStatus = "Disabled"
	}

	loggingStatus := "Enabled"
	if !h.config.EnableLogging {
		loggingStatus = "Disabled"
	}

	firewallChecked := ""
	if h.config.EnableFirewall {
		firewallChecked = "checked"
	}

	loggingChecked := ""
	if h.config.EnableLogging {
		loggingChecked = "checked"
	}

	bannedIPs := []string{}
	for ip := range h.config.BannedIPs {
		bannedIPs = append(bannedIPs, ip)
	}
	bannedIPsStr := strings.Join(bannedIPs, ", ")
	if bannedIPsStr == "" {
		bannedIPsStr = "None"
	}

	allowedUAsStr := strings.Join(h.config.AllowedUAs, ",")
	if allowedUAsStr == "" {
		allowedUAsStr = ""
	}
	
	logs := h.getRecentLogs()

	// Проверяем статус сервиса
	serviceStatus := "Not installed"
	if isActive, err := service.GetServiceStatus(); err == nil {
		if isActive {
			serviceStatus = "Running"
		} else {
			serviceStatus = "Installed but not running"
		}
	}

	// Security settings checkboxes
	suffixProtectionChecked := ""
	if h.config.Security.EnableSuffixProtection {
		suffixProtectionChecked = "checked"
	}
	
	sqlProtectionChecked := ""
	if h.config.Security.EnableSQLProtection {
		sqlProtectionChecked = "checked"
	}
	
	xssProtectionChecked := ""
	if h.config.Security.EnableXSSProtection {
		xssProtectionChecked = "checked"
	}
	
	scannerProtectionChecked := ""
	if h.config.Security.EnableScannerProtection {
		scannerProtectionChecked = "checked"
	}
	
	botProtectionChecked := ""
	if h.config.Security.EnableBotProtection {
		botProtectionChecked = "checked"
	}
	
	directoryProtectionChecked := ""
	if h.config.Security.EnableDirectoryProtection {
		directoryProtectionChecked = "checked"
	}
	
	ddosProtectionChecked := ""
	if h.config.Security.EnableDDoSProtection {
		ddosProtectionChecked = "checked"
	}
	
	geoBlockingChecked := ""
	if h.config.Security.EnableGeoBlocking {
		geoBlockingChecked = "checked"
	}

	return &TemplateData{
		StatusClass:     statusClass,
		FirewallStatus:  firewallStatus,
		LoggingStatus:   loggingStatus,
		RateLimitRPS:    h.config.RateLimitRPS,
		ListenPort:      h.config.ListenPort,
		AdminPort:       h.config.AdminPort,
		TargetPort:      h.config.TargetPort,
		ServiceStatus:   serviceStatus,
		FirewallChecked: firewallChecked,
		LoggingChecked:  loggingChecked,
		BannedIPs:       bannedIPsStr,
		AllowedUAs:      allowedUAsStr,
		Logs:            logs,
		Username:        h.config.Username,
		
		// Security settings
		SuffixProtectionChecked:    suffixProtectionChecked,
		SQLProtectionChecked:       sqlProtectionChecked,
		XSSProtectionChecked:       xssProtectionChecked,
		ScannerProtectionChecked:   scannerProtectionChecked,
		BotProtectionChecked:       botProtectionChecked,
		DirectoryProtectionChecked: directoryProtectionChecked,
		DDoSProtectionChecked:      ddosProtectionChecked,
		GeoBlockingChecked:         geoBlockingChecked,
		
		ForbiddenSuffixes:     strings.Join(h.config.Security.ForbiddenSuffixes, ","),
		SuffixBanDuration:     h.config.Security.SuffixBanDuration,
		SQLKeywords:           strings.Join(h.config.Security.SQLKeywords, ","),
		XSSPatterns:           strings.Join(h.config.Security.XSSPatterns, ","),
		ScannerPaths:          strings.Join(h.config.Security.ScannerPaths, ","),
		SuspiciousUserAgents:  strings.Join(h.config.Security.SuspiciousUserAgents, ","),
		ProtectedDirectories:  strings.Join(h.config.Security.ProtectedDirectories, ","),
		DDoSThreshold:         h.config.Security.DDoSThreshold,
		DDoSTimeWindow:        h.config.Security.DDoSTimeWindow,
		DDoSBanDuration:       h.config.Security.DDoSBanDuration,
		BlockedCountries:      strings.Join(h.config.Security.BlockedCountries, ","),
		
		TemporaryBans:         h.config.TemporaryBans,
		LogStats:              h.logger.GetLogStats(),
	}
}

// handleAdminAction обрабатывает действия админ панели
func (h *Handler) handleAdminAction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	action := r.FormValue("action")

	switch action {
	case "update_settings":
		h.updateSettings(r)
	case "update_security":
		h.updateSecuritySettings(r)
	case "unban_temp_ip":
		h.unbanTempIP(r.FormValue("ip"))
	case "clear_temp_bans":
		h.clearTempBans()
	case "ban_ip":
		h.banIP(r.FormValue("ip"))
	case "unban_ip":
		h.unbanIP(r.FormValue("ip"))
	case "update_uas":
		h.updateUserAgents(r.FormValue("allowed_uas"))
	case "clear_logs":
		h.logger.Clear()
	case "restart":
		h.config.Save()
		os.Exit(0)
	case "install_service":
		err := service.InstallService()
		if err != nil {
			h.logger.LogError("Failed to install service", err)
			http.Error(w, "Failed to install service: "+err.Error(), http.StatusInternalServerError)
			return
		} else {
			h.logger.LogInfo("Service installed successfully")
		}
	case "uninstall_service":
		err := service.UninstallService()
		if err != nil {
			h.logger.LogError("Failed to uninstall service", err)
			http.Error(w, "Failed to uninstall service: "+err.Error(), http.StatusInternalServerError)
			return
		} else {
			h.logger.LogInfo("Service uninstalled successfully")
		}
	case "start_service":
		err := service.StartService()
		if err != nil {
			h.logger.LogError("Failed to start service", err)
			http.Error(w, "Failed to start service: "+err.Error(), http.StatusInternalServerError)
			return
		} else {
			h.logger.LogInfo("Service started successfully")
		}
	case "stop_service":
		err := service.StopService()
		if err != nil {
			h.logger.LogError("Failed to stop service", err)
			http.Error(w, "Failed to stop service: "+err.Error(), http.StatusInternalServerError)
			return
		} else {
			h.logger.LogInfo("Service stopped successfully")
		}
	case "clear_stats":
		h.stats.Clear()
		h.logger.LogInfo("Statistics cleared")
	}

	h.config.Save()
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// updateSettings обновляет настройки
func (h *Handler) updateSettings(r *http.Request) {
	h.config.EnableFirewall = r.FormValue("enable_firewall") == "on"
	h.config.EnableLogging = r.FormValue("enable_logging") == "on"

	if rateLimit, err := strconv.Atoi(r.FormValue("rate_limit")); err == nil {
		h.config.RateLimitRPS = rateLimit
	}

	if listenPort, err := strconv.Atoi(r.FormValue("listen_port")); err == nil {
		h.config.ListenPort = listenPort
	}

	if targetPort, err := strconv.Atoi(r.FormValue("target_port")); err == nil {
		h.config.TargetPort = targetPort
	}

	if adminPort, err := strconv.Atoi(r.FormValue("admin_port")); err == nil {
		h.config.AdminPort = adminPort
	}

	h.logger.SetEnabled(h.config.EnableLogging)
}

// banIP добавляет IP в черный список
func (h *Handler) banIP(ip string) {
	if ip != "" {
		h.config.BannedIPs[ip] = true
		h.logger.LogInfo("IP banned: " + ip)
	}
}

// unbanIP удаляет IP из черного списка
func (h *Handler) unbanIP(ip string) {
	if ip != "" {
		delete(h.config.BannedIPs, ip)
		h.logger.LogInfo("IP unbanned: " + ip)
	}
}

// updateUserAgents обновляет список разрешенных User-Agent'ов
func (h *Handler) updateUserAgents(uasStr string) {
	if uasStr == "" {
		h.config.AllowedUAs = []string{}
	} else {
		h.config.AllowedUAs = strings.Split(uasStr, ",")
		for i, ua := range h.config.AllowedUAs {
			h.config.AllowedUAs[i] = strings.TrimSpace(ua)
		}
	}
}

// getRecentLogs получает последние логи
func (h *Handler) getRecentLogs() string {
	file, err := os.Open(logger.LogFile)
	if err != nil {
		return "No logs available"
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Показываем последние 20 строк
	start := len(lines) - 20
	if start < 0 {
		start = 0
	}

	result := strings.Join(lines[start:], "\n")
	if result == "" {
		return "No logs available"
	}

	return result
}

// handleStatsAPI обрабатывает API запросы для статистики
func (h *Handler) handleStatsAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.URL.Path {
	case "/admin/api/hourly-stats":
		json.NewEncoder(w).Encode(h.stats.GetHourlyStats())
	case "/admin/api/top-ips":
		json.NewEncoder(w).Encode(h.stats.GetTopIPs(10))
	case "/admin/api/top-uas":
		json.NewEncoder(w).Encode(h.stats.GetTopUserAgents(10))
	case "/admin/api/summary":
		json.NewEncoder(w).Encode(h.stats.GetSummary())
	case "/admin/api/service-status":
		isActive, _ := service.GetServiceStatus()
		json.NewEncoder(w).Encode(map[string]bool{"active": isActive})
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

// updateSecuritySettings обновляет настройки безопасности
func (h *Handler) updateSecuritySettings(r *http.Request) {
	h.config.Security.EnableSuffixProtection = r.FormValue("enable_suffix_protection") == "on"
	h.config.Security.EnableSQLProtection = r.FormValue("enable_sql_protection") == "on"
	h.config.Security.EnableXSSProtection = r.FormValue("enable_xss_protection") == "on"
	h.config.Security.EnableScannerProtection = r.FormValue("enable_scanner_protection") == "on"
	h.config.Security.EnableBotProtection = r.FormValue("enable_bot_protection") == "on"
	h.config.Security.EnableDirectoryProtection = r.FormValue("enable_directory_protection") == "on"
	h.config.Security.EnableDDoSProtection = r.FormValue("enable_ddos_protection") == "on"
	h.config.Security.EnableGeoBlocking = r.FormValue("enable_geo_blocking") == "on"
	
	// Update lists
	if suffixes := r.FormValue("forbidden_suffixes"); suffixes != "" {
		h.config.Security.ForbiddenSuffixes = strings.Split(suffixes, ",")
		for i, suffix := range h.config.Security.ForbiddenSuffixes {
			h.config.Security.ForbiddenSuffixes[i] = strings.TrimSpace(suffix)
		}
	}
	
	if keywords := r.FormValue("sql_keywords"); keywords != "" {
		h.config.Security.SQLKeywords = strings.Split(keywords, ",")
		for i, keyword := range h.config.Security.SQLKeywords {
			h.config.Security.SQLKeywords[i] = strings.TrimSpace(keyword)
		}
	}
	
	if patterns := r.FormValue("xss_patterns"); patterns != "" {
		h.config.Security.XSSPatterns = strings.Split(patterns, ",")
		for i, pattern := range h.config.Security.XSSPatterns {
			h.config.Security.XSSPatterns[i] = strings.TrimSpace(pattern)
		}
	}
	
	if paths := r.FormValue("scanner_paths"); paths != "" {
		h.config.Security.ScannerPaths = strings.Split(paths, ",")
		for i, path := range h.config.Security.ScannerPaths {
			h.config.Security.ScannerPaths[i] = strings.TrimSpace(path)
		}
	}
	
	if agents := r.FormValue("suspicious_user_agents"); agents != "" {
		h.config.Security.SuspiciousUserAgents = strings.Split(agents, ",")
		for i, agent := range h.config.Security.SuspiciousUserAgents {
			h.config.Security.SuspiciousUserAgents[i] = strings.TrimSpace(agent)
		}
	}
	
	if dirs := r.FormValue("protected_directories"); dirs != "" {
		h.config.Security.ProtectedDirectories = strings.Split(dirs, ",")
		for i, dir := range h.config.Security.ProtectedDirectories {
			h.config.Security.ProtectedDirectories[i] = strings.TrimSpace(dir)
		}
	}
	
	if countries := r.FormValue("blocked_countries"); countries != "" {
		h.config.Security.BlockedCountries = strings.Split(countries, ",")
		for i, country := range h.config.Security.BlockedCountries {
			h.config.Security.BlockedCountries[i] = strings.TrimSpace(country)
		}
	}
	
	// Update numeric values
	if duration, err := strconv.Atoi(r.FormValue("suffix_ban_duration")); err == nil {
		h.config.Security.SuffixBanDuration = duration
	}
	
	if threshold, err := strconv.Atoi(r.FormValue("ddos_threshold")); err == nil {
		h.config.Security.DDoSThreshold = threshold
	}
	
	if window, err := strconv.Atoi(r.FormValue("ddos_time_window")); err == nil {
		h.config.Security.DDoSTimeWindow = window
	}
	
	if banDuration, err := strconv.Atoi(r.FormValue("ddos_ban_duration")); err == nil {
		h.config.Security.DDoSBanDuration = banDuration
	}
	
	h.logger.LogAdmin(h.config.Username, "UPDATE_SECURITY", "Security settings updated")
}

// unbanTempIP удаляет временный бан IP
func (h *Handler) unbanTempIP(ip string) {
	if ip != "" {
		h.config.RemoveTemporaryBan(ip)
		h.logger.LogAdmin(h.config.Username, "UNBAN_TEMP_IP", "Temporary ban removed for IP: "+ip)
	}
}

// clearTempBans очищает все временные баны
func (h *Handler) clearTempBans() {
	count := len(h.config.TemporaryBans)
	h.config.TemporaryBans = []config.TemporaryBan{}
	h.logger.LogAdmin(h.config.Username, "CLEAR_TEMP_BANS", fmt.Sprintf("Cleared %d temporary bans", count))
}
