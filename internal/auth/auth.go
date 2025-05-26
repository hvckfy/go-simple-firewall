package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword создает хеш пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword проверяет пароль против хеша
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSessionToken генерирует случайный токен сессии
func GenerateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// SetSessionCookie устанавливает cookie сессии
func SetSessionCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "firewall_session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Для HTTP, для HTTPS должно быть true
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour), // 24 часа
	}
	http.SetCookie(w, cookie)
}

// GetSessionCookie получает токен сессии из cookie
func GetSessionCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("firewall_session")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// ClearSessionCookie очищает cookie сессии
func ClearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "firewall_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)
}

// SessionManager управляет сессиями
type SessionManager struct {
	sessions map[string]time.Time
}

// NewSessionManager создает новый менеджер сессий
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]time.Time),
	}
}

// CreateSession создает новую сессию
func (sm *SessionManager) CreateSession() (string, error) {
	token, err := GenerateSessionToken()
	if err != nil {
		return "", err
	}

	sm.sessions[token] = time.Now().Add(24 * time.Hour)
	return token, nil
}

// ValidateSession проверяет валидность сессии
func (sm *SessionManager) ValidateSession(token string) bool {
	if token == "" {
		return false
	}

	expiry, exists := sm.sessions[token]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		delete(sm.sessions, token)
		return false
	}

	// Продлеваем сессию
	sm.sessions[token] = time.Now().Add(24 * time.Hour)
	return true
}

// DeleteSession удаляет сессию
func (sm *SessionManager) DeleteSession(token string) {
	delete(sm.sessions, token)
}

// CleanupExpiredSessions очищает истекшие сессии
func (sm *SessionManager) CleanupExpiredSessions() {
	now := time.Now()
	for token, expiry := range sm.sessions {
		if now.After(expiry) {
			delete(sm.sessions, token)
		}
	}
}
