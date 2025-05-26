package ratelimit

import (
	"sync"
	"time"
)

// Limiter структура для ограничения запросов
type Limiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
}

// New создает новый rate limiter
func New(limit int) *Limiter {
	return &Limiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
	}
}

// IsAllowed проверяет, разрешен ли запрос для данного IP
func (rl *Limiter) IsAllowed(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-time.Minute)

	// Очищаем старые запросы
	if requests, exists := rl.requests[ip]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[ip] = validRequests
	}

	// Проверяем лимит
	if len(rl.requests[ip]) >= rl.limit {
		return false
	}

	// Добавляем текущий запрос
	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}

// UpdateLimit обновляет лимит запросов
func (rl *Limiter) UpdateLimit(limit int) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.limit = limit
}

// GetStats возвращает статистику по IP
func (rl *Limiter) GetStats() map[string]int {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	stats := make(map[string]int)
	for ip, requests := range rl.requests {
		stats[ip] = len(requests)
	}
	return stats
}

// Clear очищает все записи
func (rl *Limiter) Clear() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.requests = make(map[string][]time.Time)
}
