package logger

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const LogFile = "firewall.log"

// Logger структура для логирования
type Logger struct {
	file    *os.File
	enabled bool
	mutex   sync.Mutex
}

// RequestInfo детальная информация о запросе
type RequestInfo struct {
	Timestamp    time.Time
	IP           string
	Method       string
	URL          string
	UserAgent    string
	Referer      string
	Status       string
	StatusCode   int
	Reason       string
	QueryParams  string
	PostData     string
	ContentType  string
	ContentLength int64
}

// New создает новый logger
func New(enabled bool) (*Logger, error) {
	logger := &Logger{
		enabled: enabled,
	}

	if enabled {
		var err error
		logger.file, err = os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
	}

	return logger, nil
}

// LogRequest записывает HTTP запрос в лог с детальной информацией
func (l *Logger) LogRequest(r *http.Request, status string, statusCode int) {
	l.LogRequestWithReason(r, status, statusCode, "")
}

// LogRequestWithReason записывает HTTP запрос в лог с указанием причины блокировки
func (l *Logger) LogRequestWithReason(r *http.Request, status string, statusCode int, reason string) {
	if !l.enabled || l.file == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	info := l.extractRequestInfo(r, status, statusCode, reason)
	logEntry := l.formatLogEntry(info)

	l.file.WriteString(logEntry)
}

// extractRequestInfo извлекает детальную информацию из запроса
func (l *Logger) extractRequestInfo(r *http.Request, status string, statusCode int, reason string) RequestInfo {
	// Получаем полный URL с параметрами
	fullURL := r.URL.String()
	if r.URL.RawQuery != "" {
		fullURL = r.URL.Path + "?" + r.URL.RawQuery
	} else {
		fullURL = r.URL.Path
	}

	// Извлекаем параметры запроса
	queryParams := ""
	if r.URL.RawQuery != "" {
		// Декодируем параметры для лучшей читаемости
		if decoded, err := url.QueryUnescape(r.URL.RawQuery); err == nil {
			queryParams = decoded
		} else {
			queryParams = r.URL.RawQuery
		}
	}

	// Извлекаем POST данные (только для небольших запросов)
	postData := ""
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		if r.ContentLength > 0 && r.ContentLength < 1024 { // Ограничиваем размер
			r.ParseForm()
			if len(r.PostForm) > 0 {
				postData = l.formatPostData(r.PostForm)
			}
		}
	}

	return RequestInfo{
		Timestamp:     time.Now(),
		IP:            getClientIP(r),
		Method:        r.Method,
		URL:           fullURL,
		UserAgent:     r.UserAgent(),
		Referer:       r.Referer(),
		Status:        status,
		StatusCode:    statusCode,
		Reason:        reason,
		QueryParams:   queryParams,
		PostData:      postData,
		ContentType:   r.Header.Get("Content-Type"),
		ContentLength: r.ContentLength,
	}
}

// formatPostData форматирует POST данные, скрывая пароли
func (l *Logger) formatPostData(postForm url.Values) string {
	var parts []string
	
	for key, values := range postForm {
		for _, value := range values {
			// Скрываем пароли и другие чувствительные данные
			if l.isSensitiveField(key) {
				parts = append(parts, fmt.Sprintf("%s=[HIDDEN]", key))
			} else {
				// Ограничиваем длину значения
				if len(value) > 100 {
					value = value[:100] + "..."
				}
				parts = append(parts, fmt.Sprintf("%s=%s", key, value))
			}
		}
	}
	
	return strings.Join(parts, "&")
}

// isSensitiveField проверяет, является ли поле чувствительным
func (l *Logger) isSensitiveField(fieldName string) bool {
	sensitiveFields := []string{
		"password", "passwd", "pwd", "pass",
		"token", "secret", "key", "auth",
		"credit", "card", "cvv", "ssn",
		"api_key", "access_token", "refresh_token",
	}
	
	fieldLower := strings.ToLower(fieldName)
	for _, sensitive := range sensitiveFields {
		if strings.Contains(fieldLower, sensitive) {
			return true
		}
	}
	
	return false
}

// formatLogEntry форматирует запись лога
func (l *Logger) formatLogEntry(info RequestInfo) string {
	timestamp := info.Timestamp.Format("2006-01-02 15:04:05")
	
	// Базовая информация
	logEntry := fmt.Sprintf("[%s] %s %s %s \"%s\"",
		timestamp,
		info.IP,
		info.Method,
		info.URL,
		info.UserAgent,
	)
	
	// Добавляем статус
	if info.Reason != "" {
		logEntry += fmt.Sprintf(" - %s (%d) - REASON: %s", info.Status, info.StatusCode, info.Reason)
	} else {
		logEntry += fmt.Sprintf(" - %s (%d)", info.Status, info.StatusCode)
	}
	
	// Добавляем Referer если есть
	if info.Referer != "" {
		logEntry += fmt.Sprintf(" - REFERER: %s", info.Referer)
	}
	
	// Добавляем параметры запроса если есть
	if info.QueryParams != "" {
		logEntry += fmt.Sprintf(" - PARAMS: %s", info.QueryParams)
	}
	
	// Добавляем POST данные если есть
	if info.PostData != "" {
		logEntry += fmt.Sprintf(" - POST: %s", info.PostData)
	}
	
	// Добавляем Content-Type если есть
	if info.ContentType != "" {
		logEntry += fmt.Sprintf(" - CONTENT-TYPE: %s", info.ContentType)
	}
	
	// Добавляем размер контента если больше 0
	if info.ContentLength > 0 {
		logEntry += fmt.Sprintf(" - SIZE: %d bytes", info.ContentLength)
	}
	
	logEntry += "\n"
	
	return logEntry
}

// LogSecurity записывает событие безопасности
func (l *Logger) LogSecurity(ip, event, details string) {
	if !l.enabled || l.file == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	logEntry := fmt.Sprintf("[%s] SECURITY: %s - %s - %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		ip,
		event,
		details,
	)

	l.file.WriteString(logEntry)
}

// LogTemporaryBan записывает информацию о временном бане
func (l *Logger) LogTemporaryBan(ip, reason string, duration time.Duration) {
	if !l.enabled || l.file == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	logEntry := fmt.Sprintf("[%s] TEMP_BAN: %s - %s - Duration: %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		ip,
		reason,
		duration,
	)

	l.file.WriteString(logEntry)
}

// LogAttack записывает информацию об атаке
func (l *Logger) LogAttack(ip, attackType, details, userAgent, url string) {
	if !l.enabled || l.file == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	logEntry := fmt.Sprintf("[%s] ATTACK: %s - Type: %s - URL: %s - UA: \"%s\" - Details: %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		ip,
		attackType,
		url,
		userAgent,
		details,
	)

	l.file.WriteString(logEntry)
}

// LogInfo записывает информационное сообщение
func (l *Logger) LogInfo(message string) {
	if !l.enabled || l.file == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	logEntry := fmt.Sprintf("[%s] INFO: %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		message,
	)

	l.file.WriteString(logEntry)
}

// LogError записывает ошибку
func (l *Logger) LogError(message string, err error) {
	if !l.enabled || l.file == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	logEntry := fmt.Sprintf("[%s] ERROR: %s - %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		message,
		err,
	)

	l.file.WriteString(logEntry)
}

// LogAdmin записывает действия администратора
func (l *Logger) LogAdmin(username, action, details string) {
	if !l.enabled || l.file == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	logEntry := fmt.Sprintf("[%s] ADMIN: %s - %s - %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		username,
		action,
		details,
	)

	l.file.WriteString(logEntry)
}

// Clear очищает лог файл
func (l *Logger) Clear() error {
	if l.file == nil {
		return nil
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	return l.file.Truncate(0)
}

// Close закрывает лог файл
func (l *Logger) Close() error {
	if l.file == nil {
		return nil
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	return l.file.Close()
}

// SetEnabled включает/выключает логирование
func (l *Logger) SetEnabled(enabled bool) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if enabled && !l.enabled {
		var err error
		l.file, err = os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
	} else if !enabled && l.enabled && l.file != nil {
		l.file.Close()
		l.file = nil
	}

	l.enabled = enabled
	return nil
}

// GetLogStats возвращает статистику логов
func (l *Logger) GetLogStats() map[string]interface{} {
	if l.file == nil {
		return map[string]interface{}{
			"enabled": false,
			"size":    0,
		}
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	fileInfo, err := l.file.Stat()
	if err != nil {
		return map[string]interface{}{
			"enabled": l.enabled,
			"size":    0,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"enabled":      l.enabled,
		"size":         fileInfo.Size(),
		"last_modified": fileInfo.ModTime(),
	}
}

// getClientIP извлекает IP клиента из запроса
func getClientIP(r *http.Request) string {
	// Проверяем заголовки прокси
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Берем первый IP из списка
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Проверяем Cloudflare
	cfIP := r.Header.Get("CF-Connecting-IP")
	if cfIP != "" {
		return cfIP
	}

	// Используем RemoteAddr как последний вариант
	ip := r.RemoteAddr
	if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
		ip = ip[:colonIndex]
	}

	return ip
}
