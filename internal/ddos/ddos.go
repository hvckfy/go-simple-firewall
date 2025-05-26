package ddos

import (
	"sync"
	"time"

	"go-simple-firewall/internal/config"
)

// Request информация о запросе
type Request struct {
	Timestamp time.Time
	IP        string
}

// DDoSProtection защита от DDoS атак
type DDoSProtection struct {
	config   *config.Config
	requests []Request
	mutex    sync.RWMutex
}

// New создает новую DDoS защиту
func New(cfg *config.Config) *DDoSProtection {
	ddos := &DDoSProtection{
		config:   cfg,
		requests: make([]Request, 0),
	}
	
	// Запускаем очистку старых запросов каждые 30 секунд
	go ddos.cleanupOldRequests()
	
	return ddos
}

// CheckRequest проверяет запрос на DDoS
func (d *DDoSProtection) CheckRequest(ip string) (blocked bool, reason string) {
	if !d.config.Security.EnableDDoSProtection {
		return false, ""
	}
	
	d.mutex.Lock()
	defer d.mutex.Unlock()
	
	now := time.Now()
	timeWindow := time.Duration(d.config.Security.DDoSTimeWindow) * time.Second
	threshold := d.config.Security.DDoSThreshold
	
	// Добавляем текущий запрос
	d.requests = append(d.requests, Request{
		Timestamp: now,
		IP:        ip,
	})
	
	// Считаем запросы от этого IP за временное окно
	count := 0
	cutoff := now.Add(-timeWindow)
	
	for _, req := range d.requests {
		if req.IP == ip && req.Timestamp.After(cutoff) {
			count++
		}
	}
	
	if count > threshold {
		// Добавляем временный бан
		banDuration := time.Duration(d.config.Security.DDoSBanDuration) * time.Minute
		d.config.AddTemporaryBan(ip, "DDoS attack detected", banDuration)
		return true, "DDoS attack detected"
	}
	
	return false, ""
}

// cleanupOldRequests очищает старые запросы
func (d *DDoSProtection) cleanupOldRequests() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			d.mutex.Lock()
			
			now := time.Now()
			timeWindow := time.Duration(d.config.Security.DDoSTimeWindow) * time.Second
			cutoff := now.Add(-timeWindow)
			
			validRequests := make([]Request, 0)
			for _, req := range d.requests {
				if req.Timestamp.After(cutoff) {
					validRequests = append(validRequests, req)
				}
			}
			
			d.requests = validRequests
			d.mutex.Unlock()
		}
	}
}

// UpdateConfig обновляет конфигурацию
func (d *DDoSProtection) UpdateConfig(cfg *config.Config) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.config = cfg
}
