package stats

import (
	"sync"
	"time"
)

// RequestStats структура для хранения статистики запросов
type RequestStats struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int       `json:"count"`
	Blocked   int       `json:"blocked"`
	Allowed   int       `json:"allowed"`
}

// IPStats статистика по IP адресам
type IPStats struct {
	IP       string    `json:"ip"`
	Requests int       `json:"requests"`
	Blocked  int       `json:"blocked"`
	LastSeen time.Time `json:"last_seen"`
}

// UserAgentStats статистика по User-Agent
type UserAgentStats struct {
	UserAgent string `json:"user_agent"`
	Requests  int    `json:"requests"`
	Blocked   int    `json:"blocked"`
}

// Stats основная структура для сбора статистики
type Stats struct {
	mutex          sync.RWMutex
	hourlyStats    []RequestStats
	ipStats        map[string]*IPStats
	userAgentStats map[string]*UserAgentStats
	totalRequests  int64
	totalBlocked   int64
	startTime      time.Time
}

// New создает новый экземпляр статистики
func New() *Stats {
	return &Stats{
		hourlyStats:    make([]RequestStats, 0),
		ipStats:        make(map[string]*IPStats),
		userAgentStats: make(map[string]*UserAgentStats),
		startTime:      time.Now(),
	}
}

// RecordRequest записывает информацию о запросе
func (s *Stats) RecordRequest(ip, userAgent string, blocked bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.totalRequests++
	if blocked {
		s.totalBlocked++
	}

	// Обновляем статистику по IP
	if ipStat, exists := s.ipStats[ip]; exists {
		ipStat.Requests++
		if blocked {
			ipStat.Blocked++
		}
		ipStat.LastSeen = time.Now()
	} else {
		blocked_count := 0
		if blocked {
			blocked_count = 1
		}
		s.ipStats[ip] = &IPStats{
			IP:       ip,
			Requests: 1,
			Blocked:  blocked_count,
			LastSeen: time.Now(),
		}
	}

	// Обновляем статистику по User-Agent
	if uaStat, exists := s.userAgentStats[userAgent]; exists {
		uaStat.Requests++
		if blocked {
			uaStat.Blocked++
		}
	} else {
		blocked_count := 0
		if blocked {
			blocked_count = 1
		}
		s.userAgentStats[userAgent] = &UserAgentStats{
			UserAgent: userAgent,
			Requests:  1,
			Blocked:   blocked_count,
		}
	}

	// Обновляем почасовую статистику
	s.updateHourlyStats(blocked)
}

// updateHourlyStats обновляет почасовую статистику
func (s *Stats) updateHourlyStats(blocked bool) {
	now := time.Now()
	currentHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	// Если это первая запись или новый час
	if len(s.hourlyStats) == 0 || s.hourlyStats[len(s.hourlyStats)-1].Timestamp.Before(currentHour) {
		newStat := RequestStats{
			Timestamp: currentHour,
			Count:     1,
			Blocked:   0,
			Allowed:   1,
		}
		if blocked {
			newStat.Blocked = 1
			newStat.Allowed = 0
		}
		s.hourlyStats = append(s.hourlyStats, newStat)
	} else {
		// Обновляем текущий час
		lastIndex := len(s.hourlyStats) - 1
		s.hourlyStats[lastIndex].Count++
		if blocked {
			s.hourlyStats[lastIndex].Blocked++
		} else {
			s.hourlyStats[lastIndex].Allowed++
		}
	}

	// Оставляем только последние 24 часа
	if len(s.hourlyStats) > 24 {
		s.hourlyStats = s.hourlyStats[len(s.hourlyStats)-24:]
	}
}

// GetHourlyStats возвращает почасовую статистику
func (s *Stats) GetHourlyStats() []RequestStats {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if len(s.hourlyStats) == 0 {
		return []RequestStats{}
	}

	result := make([]RequestStats, len(s.hourlyStats))
	copy(result, s.hourlyStats)
	return result
}

// GetTopIPs возвращает топ IP адресов по количеству запросов
func (s *Stats) GetTopIPs(limit int) []*IPStats {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if len(s.ipStats) == 0 {
		return []*IPStats{}
	}

	var ips []*IPStats
	for _, ipStat := range s.ipStats {
		ips = append(ips, ipStat)
	}

	// Сортируем по количеству запросов
	for i := 0; i < len(ips)-1; i++ {
		for j := i + 1; j < len(ips); j++ {
			if ips[i].Requests < ips[j].Requests {
				ips[i], ips[j] = ips[j], ips[i]
			}
		}
	}

	if len(ips) > limit {
		ips = ips[:limit]
	}

	return ips
}

// GetTopUserAgents возвращает топ User-Agent по количеству запросов
func (s *Stats) GetTopUserAgents(limit int) []*UserAgentStats {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if len(s.userAgentStats) == 0 {
		return []*UserAgentStats{}
	}

	var uas []*UserAgentStats
	for _, uaStat := range s.userAgentStats {
		uas = append(uas, uaStat)
	}

	// Сортируем по количеству запросов
	for i := 0; i < len(uas)-1; i++ {
		for j := i + 1; j < len(uas); j++ {
			if uas[i].Requests < uas[j].Requests {
				uas[i], uas[j] = uas[j], uas[i]
			}
		}
	}

	if len(uas) > limit {
		uas = uas[:limit]
	}

	return uas
}

// GetSummary возвращает общую статистику
func (s *Stats) GetSummary() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	uptime := time.Since(s.startTime)

	return map[string]interface{}{
		"total_requests": s.totalRequests,
		"total_blocked":  s.totalBlocked,
		"total_allowed":  s.totalRequests - s.totalBlocked,
		"unique_ips":     len(s.ipStats),
		"unique_uas":     len(s.userAgentStats),
		"uptime_seconds": int64(uptime.Seconds()),
		"start_time":     s.startTime,
	}
}

// Clear очищает всю статистику
func (s *Stats) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.hourlyStats = make([]RequestStats, 0)
	s.ipStats = make(map[string]*IPStats)
	s.userAgentStats = make(map[string]*UserAgentStats)
	s.totalRequests = 0
	s.totalBlocked = 0
	s.startTime = time.Now()
}
