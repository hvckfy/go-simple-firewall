# 🔌 API документация Go Simple Firewall

Подробная документация по REST API для управления firewall и получения статистики.

## 📋 Содержание

- [Обзор API](#обзор-api)
- [Аутентификация](#аутентификация)
- [Статистика](#статистика)
- [Управление IP](#управление-ip)
- [Конфигурация](#конфигурация)
- [Логи](#логи)
- [Системное управление](#системное-управление)
- [WebSocket API](#websocket-api)
- [Примеры использования](#примеры-использования)

## 🌐 Обзор API

Base URL: `http://localhost:9090/admin/api/`

Все API endpoints требуют аутентификации через сессионные cookies или API ключи.

### Формат ответов

Все ответы возвращаются в формате JSON:

```json
{
  "success": true,
  "data": { ... },
  "message": "Operation completed successfully",
  "timestamp": "2024-01-15T14:30:25Z"
}
```

### Коды ошибок

| Код | Описание |
|-----|----------|
| 200 | Успешно |
| 400 | Неверный запрос |
| 401 | Не авторизован |
| 403 | Доступ запрещен |
| 404 | Не найдено |
| 500 | Внутренняя ошибка сервера |

## 🔐 Аутентификация

### Вход в систему

```http
POST /admin/login
Content-Type: application/x-www-form-urlencoded

username=admin&password=your_password
```

**Ответ:**
```json
{
  "success": true,
  "message": "Login successful",
  "session_expires": "2024-01-16T14:30:25Z"
}
```

### Выход из системы

```http
POST /admin/logout
```

### Проверка сессии

```http
GET /admin/api/session
```

**Ответ:**
```json
{
  "success": true,
  "data": {
    "authenticated": true,
    "username": "admin",
    "expires_at": "2024-01-16T14:30:25Z"
  }
}
```

## 📊 Статистика

### Общая статистика

```http
GET /admin/api/summary
```

**Ответ:**
```json
{
  "success": true,
  "data": {
    "total_requests": 15420,
    "total_blocked": 1250,
    "total_allowed": 14170,
    "unique_ips": 342,
    "unique_uas": 89,
    "uptime_seconds": 86400,
    "start_time": "2024-01-14T14:30:25Z"
  }
}
```

### Почасовая статистика

```http
GET /admin/api/hourly-stats
```

**Ответ:**
```json
{
  "success": true,
  "data": [
    {
      "timestamp": "2024-01-15T14:00:00Z",
      "count": 450,
      "blocked": 23,
      "allowed": 427
    },
    {
      "timestamp": "2024-01-15T15:00:00Z", 
      "count": 523,
      "blocked": 31,
      "allowed": 492
    }
  ]
}
```

### Топ IP адресов

```http
GET /admin/api/top-ips?limit=10
```

**Параметры:**
- `limit` (optional) - количество записей (по умолчанию 10)

**Ответ:**
```json
{
  "success": true,
  "data": [
    {
      "ip": "192.168.1.100",
      "requests": 1250,
      "blocked": 45,
      "last_seen": "2024-01-15T14:30:25Z"
    },
    {
      "ip": "10.0.0.50",
      "requests": 890,
      "blocked": 12,
      "last_seen": "2024-01-15T14:25:10Z"
    }
  ]
}
```

### Топ User-Agent'ов

```http
GET /admin/api/top-uas?limit=10
```

**Ответ:**
```json
{
  "success": true,
  "data": [
    {
      "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
      "requests": 2340,
      "blocked": 15
    },
    {
      "user_agent": "curl/7.68.0",
      "requests": 156,
      "blocked": 89
    }
  ]
}
```

### Статистика атак

```http
GET /admin/api/attack-stats
```

**Ответ:**
```json
{
  "success": true,
  "data": {
    "sql_injection": 45,
    "xss_attempts": 23,
    "ddos_attacks": 8,
    "scanner_attempts": 156,
    "bot_requests": 89,
    "directory_traversal": 34,
    "malicious_files": 12
  }
}
```

## 🌐 Управление IP

### Получить списки IP

```http
GET /admin/api/ip-lists
```

**Ответ:**
```json
{
  "success": true,
  "data": {
    "banned_ips": [
      "192.168.1.100",
      "203.0.113.0/24"
    ],
    "allowed_ips": [
      "127.0.0.1",
      "192.168.1.0/24"
    ],
    "temporary_bans": [
      {
        "ip": "10.0.0.50",
        "reason": "DDoS attack detected",
        "expires_at": "2024-01-15T15:30:25Z"
      }
    ]
  }
}
```

### Заблокировать IP

```http
POST /admin/api/ban-ip
Content-Type: application/json

{
  "ip": "192.168.1.100",
  "reason": "Manual ban"
}
```

**Ответ:**
```json
{
  "success": true,
  "message": "IP 192.168.1.100 has been banned"
}
```

### Разблокировать IP

```http
POST /admin/api/unban-ip
Content-Type: application/json

{
  "ip": "192.168.1.100"
}
```

### Добавить в whitelist

```http
POST /admin/api/whitelist-ip
Content-Type: application/json

{
  "ip": "192.168.1.1",
  "comment": "Office IP"
}
```

### Удалить из whitelist

```http
DELETE /admin/api/whitelist-ip
Content-Type: application/json

{
  "ip": "192.168.1.1"
}
```

### Создать временный бан

```http
POST /admin/api/temp-ban
Content-Type: application/json

{
  "ip": "10.0.0.50",
  "reason": "Suspicious activity",
  "duration_minutes": 60
}
```

### Снять временный бан

```http
DELETE /admin/api/temp-ban
Content-Type: application/json

{
  "ip": "10.0.0.50"
}
```

## ⚙️ Конфигурация

### Получить конфигурацию

```http
GET /admin/api/config
```

**Ответ:**
```json
{
  "success": true,
  "data": {
    "listen_port": 8080,
    "admin_port": 9090,
    "target_port": 3000,
    "rate_limit_rps": 60,
    "enable_logging": true,
    "enable_firewall": true,
    "security": {
      "enable_sql_protection": true,
      "enable_xss_protection": true,
      "enable_ddos_protection": true
    }
  }
}
```

### Обновить основные настройки

```http
PUT /admin/api/config/basic
Content-Type: application/json

{
  "rate_limit_rps": 120,
  "enable_logging": true,
  "enable_firewall": true
}
```

### Обновить настройки безопасности

```http
PUT /admin/api/config/security
Content-Type: application/json

{
  "enable_sql_protection": true,
  "sql_keywords": ["union", "select", "insert"],
  "enable_ddos_protection": true,
  "ddos_threshold": 150,
  "ddos_time_window": 60
}
```

### Обновить порты

```http
PUT /admin/api/config/ports
Content-Type: application/json

{
  "listen_port": 80,
  "target_port": 3000,
  "admin_port": 9090
}
```

## 📝 Логи

### Получить последние логи

```http
GET /admin/api/logs?lines=100&filter=BLOCKED
```

**Параметры:**
- `lines` (optional) - количество строк (по умолчанию 50)
- `filter` (optional) - фильтр по содержимому
- `since` (optional) - с определенного времени (ISO 8601)

**Ответ:**
```json
{
  "success": true,
  "data": {
    "logs": [
      "[2024-01-15 14:30:25] 192.168.1.100 GET /api/users - BLOCKED (403) - REASON: Rate limit exceeded",
      "[2024-01-15 14:30:24] 10.0.0.50 POST /login - ALLOWED (200)"
    ],
    "total_lines": 15420,
    "filtered_lines": 156
  }
}
```

### Очистить логи

```http
DELETE /admin/api/logs
```

### Экспорт логов

```http
GET /admin/api/logs/export?format=csv&since=2024-01-15T00:00:00Z
```

**Параметры:**
- `format` - формат экспорта (csv, json, txt)
- `since` - с определенного времени
- `until` - до определенного времени

## 🔧 Системное управление

### Статус сервиса

```http
GET /admin/api/service-status
```

**Ответ:**
```json
{
  "success": true,
  "data": {
    "active": true,
    "status": "running",
    "uptime": "2 days, 14 hours",
    "memory_usage": "45.2 MB",
    "cpu_usage": "2.1%"
  }
}
```

### Перезапуск системы

```http
POST /admin/api/restart
```

### Установка как сервис

```http
POST /admin/api/install-service
```

### Удаление сервиса

```http
DELETE /admin/api/install-service
```

### Запуск сервиса

```http
POST /admin/api/service/start
```

### Остановка сервиса

```http
POST /admin/api/service/stop
```

### Системная информация

```http
GET /admin/api/system-info
```

**Ответ:**
```json
{
  "success": true,
  "data": {
    "version": "1.0.0",
    "go_version": "go1.21.5",
    "os": "linux",
    "arch": "amd64",
    "hostname": "firewall-server",
    "pid": 12345,
    "start_time": "2024-01-14T14:30:25Z"
  }
}
```

## 🔄 WebSocket API

### Подключение к WebSocket

```javascript
const ws = new WebSocket('ws://localhost:9090/admin/api/ws');

ws.onopen = function() {
    console.log('Connected to firewall WebSocket');
};

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('Received:', data);
};
```

### Типы сообщений

**Статистика в реальном времени:**
```json
{
  "type": "stats_update",
  "data": {
    "requests_per_second": 15,
    "blocked_per_second": 2,
    "active_connections": 45
  }
}
```

**Новая атака:**
```json
{
  "type": "attack_detected",
  "data": {
    "ip": "192.168.1.100",
    "type": "sql_injection",
    "url": "/api/users?id=1' UNION SELECT",
    "timestamp": "2024-01-15T14:30:25Z"
  }
}
```

**Новый бан:**
```json
{
  "type": "ip_banned",
  "data": {
    "ip": "10.0.0.50",
    "reason": "DDoS attack detected",
    "duration": "30 minutes"
  }
}
```

### Подписка на события

```javascript
// Подписка на определенные типы событий
ws.send(JSON.stringify({
    "action": "subscribe",
    "events": ["attack_detected", "ip_banned", "stats_update"]
}));
```

## 💡 Примеры использования

### Python клиент

```python
import requests
import json

class FirewallAPI:
    def __init__(self, base_url, username, password):
        self.base_url = base_url
        self.session = requests.Session()
        self.login(username, password)
    
    def login(self, username, password):
        response = self.session.post(
            f"{self.base_url}/admin/login",
            data={"username": username, "password": password}
        )
        response.raise_for_status()
    
    def get_stats(self):
        response = self.session.get(f"{self.base_url}/admin/api/summary")
        return response.json()
    
    def ban_ip(self, ip, reason="Manual ban"):
        response = self.session.post(
            f"{self.base_url}/admin/api/ban-ip",
            json={"ip": ip, "reason": reason}
        )
        return response.json()
    
    def get_top_ips(self, limit=10):
        response = self.session.get(
            f"{self.base_url}/admin/api/top-ips",
            params={"limit": limit}
        )
        return response.json()

# Использование
api = FirewallAPI("http://localhost:9090", "admin", "password")

# Получить статистику
stats = api.get_stats()
print(f"Total requests: {stats['data']['total_requests']}")

# Заблокировать IP
api.ban_ip("192.168.1.100", "Suspicious activity")

# Получить топ IP
top_ips = api.get_top_ips(5)
for ip_info in top_ips['data']:
    print(f"IP: {ip_info['ip']}, Requests: {ip_info['requests']}")
```

### JavaScript клиент

```javascript
class FirewallAPI {
    constructor(baseUrl) {
        this.baseUrl = baseUrl;
    }
    
    async login(username, password) {
        const response = await fetch(`${this.baseUrl}/admin/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `username=${username}&password=${password}`,
            credentials: 'include'
        });
        
        if (!response.ok) {
            throw new Error('Login failed');
        }
        
        return response.json();
    }
    
    async getStats() {
        const response = await fetch(`${this.baseUrl}/admin/api/summary`, {
            credentials: 'include'
        });
        return response.json();
    }
    
    async banIP(ip, reason = 'Manual ban') {
        const response = await fetch(`${this.baseUrl}/admin/api/ban-ip`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ ip, reason }),
            credentials: 'include'
        });
        return response.json();
    }
    
    async getTopIPs(limit = 10) {
        const response = await fetch(
            `${this.baseUrl}/admin/api/top-ips?limit=${limit}`,
            { credentials: 'include' }
        );
        return response.json();
    }
}

// Использование
const api = new FirewallAPI('http://localhost:9090');

// Авторизация
await api.login('admin', 'password');

// Получение статистики
const stats = await api.getStats();
console.log('Total requests:', stats.data.total_requests);

// Блокировка IP
await api.banIP('192.168.1.100', 'Suspicious activity');
```

### Bash скрипты

```bash
#!/bin/bash

# Конфигурация
BASE_URL="http://localhost:9090"
USERNAME="admin"
PASSWORD="password"

# Функция для авторизации
login() {
    curl -c cookies.txt -X POST "$BASE_URL/admin/login" \
         -d "username=$USERNAME&password=$PASSWORD"
}

# Функция для получения статистики
get_stats() {
    curl -b cookies.txt "$BASE_URL/admin/api/summary" | jq '.'
}

# Функция для блокировки IP
ban_ip() {
    local ip=$1
    local reason=${2:-"Manual ban"}
    
    curl -b cookies.txt -X POST "$BASE_URL/admin/api/ban-ip" \
         -H "Content-Type: application/json" \
         -d "{\"ip\":\"$ip\",\"reason\":\"$reason\"}"
}

# Функция для получения топ IP
get_top_ips() {
    local limit=${1:-10}
    curl -b cookies.txt "$BASE_URL/admin/api/top-ips?limit=$limit" | jq '.data'
}

# Использование
login
get_stats
ban_ip "192.168.1.100" "Automated ban"
get_top_ips 5
```

### Мониторинг скрипт

```bash
#!/bin/bash

# Скрипт мониторинга firewall
FIREWALL_API="http://localhost:9090/admin/api"
ALERT_EMAIL="admin@example.com"
THRESHOLD_ATTACKS=50

# Авторизация
curl -c /tmp/fw_cookies.txt -X POST "http://localhost:9090/admin/login" \
     -d "username=admin&password=password" > /dev/null 2>&1

# Получение статистики
STATS=$(curl -s -b /tmp/fw_cookies.txt "$FIREWALL_API/summary")
ATTACKS=$(echo "$STATS" | jq -r '.data.total_blocked')

# Проверка превышения порога
if [ "$ATTACKS" -gt "$THRESHOLD_ATTACKS" ]; then
    echo "High attack volume detected: $ATTACKS attacks" | \
         mail -s "Firewall Alert" "$ALERT_EMAIL"
fi

# Получение топ атакующих IP
TOP_IPS=$(curl -s -b /tmp/fw_cookies.txt "$FIREWALL_API/top-ips?limit=5")
echo "$TOP_IPS" | jq -r '.data[] | "\(.ip): \(.requests) requests, \(.blocked) blocked"'

# Очистка
rm -f /tmp/fw_cookies.txt
```

## 🔒 Безопасность API

### Rate Limiting

API endpoints имеют собственные лимиты:

- Аутентификация: 5 попыток в минуту
- Статистика: 60 запросов в минуту  
- Управление: 30 запросов в минуту

### CORS

API поддерживает CORS для веб-приложений:

```
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Allow-Methods: GET, POST, PUT, DELETE
Access-Control-Allow-Headers: Content-Type, Authorization
```

### API ключи (планируется)

В будущих версиях будет добавлена поддержка API ключей:

```http
GET /admin/api/summary
Authorization: Bearer your-api-key
```

---

Эта документация покрывает все доступные API endpoints. Для получения дополнительной информации обращайтесь к исходному коду или создавайте Issues в репозитории.

