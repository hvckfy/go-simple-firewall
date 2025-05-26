# ⚙️ Конфигурация Go Simple Firewall

Подробное руководство по настройке и конфигурации всех компонентов Go Simple Firewall.

## 📋 Содержание

- [Структура конфигурации](#структура-конфигурации)
- [Основные настройки](#основные-настройки)
- [Настройки безопасности](#настройки-безопасности)
- [Управление IP адресами](#управление-ip-адресами)
- [Rate Limiting](#rate-limiting)
- [Логирование](#логирование)
- [Временные баны](#временные-баны)
- [Примеры конфигураций](#примеры-конфигураций)
- [Переменные окружения](#переменные-окружения)

## 📁 Структура конфигурации

Firewall использует JSON файл `firewall.json` для хранения всех настроек. Файл создается автоматически при первом запуске.

\`\`\`json
{
  "listen_port": 8080,
  "admin_port": 9090,
  "target_port": 3000,
  "rate_limit_rps": 60,
  "enable_logging": true,
  "enable_firewall": true,
  "banned_ips": {},
  "allowed_ips": {},
  "allowed_user_agents": [],
  "auto_start": false,
  "username": "admin",
  "password_hash": "$2a$10$...",
  "is_setup": true,
  "security": { ... },
  "temporary_bans": []
}
\`\`\`

## 🔧 Основные настройки

### Сетевые порты

| Параметр | Описание | По умолчанию | Диапазон |
|----------|----------|--------------|----------|
| `listen_port` | Порт для входящих запросов | 8080 | 1-65535 |
| `admin_port` | Порт админ-панели | 9090 | 1-65535 |
| `target_port` | Порт целевого приложения | 3000 | 1-65535 |

\`\`\`json
{
  "listen_port": 8080,
  "admin_port": 9090,
  "target_port": 3000
}
\`\`\`

### Основные переключатели

| Параметр | Описание | По умолчанию |
|----------|----------|--------------|
| `enable_firewall` | Включить/выключить firewall | true |
| `enable_logging` | Включить/выключить логирование | true |
| `auto_start` | Автозапуск при старте системы | false |

\`\`\`json
{
  "enable_firewall": true,
  "enable_logging": true,
  "auto_start": false
}
\`\`\`

## 🛡️ Настройки безопасности

### Защита от запрещенных суффиксов

Блокирует запросы к файлам с определенными расширениями:

\`\`\`json
{
  "security": {
    "enable_suffix_protection": true,
    "forbidden_suffixes": [".php", ".asp", ".aspx", ".jsp", ".cgi", ".pl"],
    "suffix_ban_duration": 10
  }
}
\`\`\`

**Параметры:**
- `enable_suffix_protection` — включить защиту
- `forbidden_suffixes` — список запрещенных расширений
- `suffix_ban_duration` — длительность бана в часах

### Защита от SQL инъекций

Обнаруживает и блокирует попытки SQL инъекций:

\`\`\`json
{
  "security": {
    "enable_sql_protection": true,
    "sql_keywords": [
      "union", "select", "insert", "delete", "update", 
      "drop", "create", "alter", "exec", "script",
      "declare", "cast", "convert", "having", "where"
    ]
  }
}
\`\`\`

### Защита от XSS атак

Блокирует попытки Cross-Site Scripting:

\`\`\`json
{
  "security": {
    "enable_xss_protection": true,
    "xss_patterns": [
      "<script", "javascript:", "onload=", "onerror=", 
      "onclick=", "onmouseover=", "onfocus=", "onblur=",
      "eval(", "alert(", "confirm(", "prompt("
    ]
  }
}
\`\`\`

### Защита от сканеров

Блокирует доступ к часто сканируемым путям:

\`\`\`json
{
  "security": {
    "enable_scanner_protection": true,
    "scanner_paths": [
      "/admin", "/wp-admin", "/phpmyadmin", "/cpanel",
      "/webmail", "/.env", "/config", "/backup",
      "/test", "/demo", "/staging", "/dev"
    ]
  }
}
\`\`\`

### Защита от ботов

Фильтрует подозрительные User-Agent'ы:

\`\`\`json
{
  "security": {
    "enable_bot_protection": true,
    "suspicious_user_agents": [
      "bot", "crawler", "spider", "scraper", "scanner",
      "nikto", "sqlmap", "nmap", "masscan", "zap"
    ]
  }
}
\`\`\`

### Защита директорий

Блокирует доступ к системным директориям:

\`\`\`json
{
  "security": {
    "enable_directory_protection": true,
    "protected_directories": [
      "/.git", "/.svn", "/backup", "/config", 
      "/logs", "/tmp", "/.env", "/node_modules"
    ]
  }
}
\`\`\`

### DDoS защита

Защищает от атак типа "отказ в обслуживании":

\`\`\`json
{
  "security": {
    "enable_ddos_protection": true,
    "ddos_threshold": 100,
    "ddos_time_window": 60,
    "ddos_ban_duration": 30
  }
}
\`\`\`

**Параметры:**
- `ddos_threshold` — максимум запросов за период
- `ddos_time_window` — временное окно в секундах
- `ddos_ban_duration` — длительность бана в минутах

### Геоблокировка

Блокирует запросы из определенных стран:

\`\`\`json
{
  "security": {
    "enable_geo_blocking": true,
    "blocked_countries": ["CN", "RU", "KP", "IR"]
  }
}
\`\`\`

## 🌐 Управление IP адресами

### Черный список IP

Постоянно заблокированные IP адреса:

\`\`\`json
{
  "banned_ips": {
    "192.168.1.100": true,
    "10.0.0.50": true,
    "203.0.113.0/24": true
  }
}
\`\`\`

### Белый список IP

Всегда разрешенные IP адреса (обходят все проверки):

\`\`\`json
{
  "allowed_ips": {
    "127.0.0.1": true,
    "192.168.1.1": true,
    "10.0.0.1": true
  }
}
\`\`\`

### Управление через админ-панель

1. Откройте админ-панель: `http://localhost:9090/admin`
2. Перейдите в раздел "IP Management"
3. Добавьте IP в соответствующий список

### Управление через API

\`\`\`bash
# Заблокировать IP
curl -X POST http://localhost:9090/admin \
  -d "action=ban_ip&ip=192.168.1.100"

# Разблокировать IP
curl -X POST http://localhost:9090/admin \
  -d "action=unban_ip&ip=192.168.1.100"
\`\`\`

## ⏱️ Rate Limiting

### Базовая конфигурация

\`\`\`json
{
  "rate_limit_rps": 60
}
\`\`\`

Ограничивает количество запросов с одного IP до 60 в минуту.

### Рекомендуемые значения

| Тип сайта | RPS | Описание |
|-----------|-----|----------|
| Статический сайт | 30-60 | Низкая нагрузка |
| Блог/CMS | 60-120 | Средняя нагрузка |
| API сервис | 120-300 | Высокая нагрузка |
| Высоконагруженный API | 300+ | Очень высокая нагрузка |

### Настройка через админ-панель

1. Откройте раздел "Settings"
2. Измените значение "Rate Limit (req/min)"
3. Нажмите "Save Settings"

## 📝 Логирование

### Конфигурация логирования

\`\`\`json
{
  "enable_logging": true
}
\`\`\`

### Типы логов

Firewall записывает различные типы событий:

1. **HTTP запросы** — все входящие запросы
2. **События безопасности** — обнаруженные атаки
3. **Административные действия** — действия в админ-панели
4. **Системные события** — запуск, остановка, ошибки

### Формат логов

\`\`\`
[2024-01-15 14:30:25] 192.168.1.100 GET /api/users "Mozilla/5.0..." - ALLOWED (200)
[2024-01-15 14:30:26] 192.168.1.101 POST /login "curl/7.68.0" - BLOCKED (403) - REASON: Rate limit exceeded
[2024-01-15 14:30:27] SECURITY: 192.168.1.102 - SQL injection attempt - union select * from users
[2024-01-15 14:30:28] ADMIN: admin - BAN_IP - IP banned: 192.168.1.102
\`\`\`

### Ротация логов

Для предотвращения переполнения диска настройте ротацию логов:

#### Linux (logrotate)

Создайте файл `/etc/logrotate.d/go-simple-firewall`:

\`\`\`
/path/to/firewall.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 root root
    postrotate
        systemctl reload go-simple-firewall
    endscript
}
\`\`\`

## ⏰ Временные баны

### Структура временного бана

\`\`\`json
{
  "temporary_bans": [
    {
      "ip": "192.168.1.100",
      "reason": "DDoS attack detected",
      "expires_at": "2024-01-15T15:30:00Z"
    }
  ]
}
\`\`\`

### Автоматические баны

Firewall автоматически создает временные баны при:

- Превышении DDoS лимитов
- Обращении к запрещенным файлам
- Множественных попытках атак

### Управление временными банами

Через админ-панель в разделе "Temporary Bans":

- Просмотр активных банов
- Ручное снятие банов
- Очистка всех банов

## 📚 Примеры конфигураций

### Конфигурация для блога

\`\`\`json
{
  "listen_port": 80,
  "admin_port": 9090,
  "target_port": 3000,
  "rate_limit_rps": 60,
  "enable_logging": true,
  "enable_firewall": true,
  "security": {
    "enable_suffix_protection": true,
    "forbidden_suffixes": [".php", ".asp"],
    "enable_sql_protection": true,
    "enable_xss_protection": true,
    "enable_scanner_protection": true,
    "enable_bot_protection": false,
    "enable_ddos_protection": true,
    "ddos_threshold": 50,
    "ddos_time_window": 60
  }
}
\`\`\`

### Конфигурация для API

\`\`\`json
{
  "listen_port": 8080,
  "admin_port": 9090,
  "target_port": 3000,
  "rate_limit_rps": 300,
  "enable_logging": true,
  "enable_firewall": true,
  "security": {
    "enable_suffix_protection": false,
    "enable_sql_protection": true,
    "enable_xss_protection": true,
    "enable_scanner_protection": true,
    "enable_bot_protection": true,
    "enable_ddos_protection": true,
    "ddos_threshold": 200,
    "ddos_time_window": 30
  }
}
\`\`\`

### Конфигурация для высокой безопасности

\`\`\`json
{
  "listen_port": 443,
  "admin_port": 9090,
  "target_port": 3000,
  "rate_limit_rps": 30,
  "enable_logging": true,
  "enable_firewall": true,
  "security": {
    "enable_suffix_protection": true,
    "forbidden_suffixes": [".php", ".asp", ".jsp", ".cgi", ".pl", ".py"],
    "suffix_ban_duration": 24,
    "enable_sql_protection": true,
    "enable_xss_protection": true,
    "enable_scanner_protection": true,
    "enable_bot_protection": true,
    "enable_directory_protection": true,
    "enable_ddos_protection": true,
    "ddos_threshold": 20,
    "ddos_time_window": 60,
    "ddos_ban_duration": 60,
    "enable_geo_blocking": true,
    "blocked_countries": ["CN", "RU", "KP"]
  }
}
\`\`\`

## 🌍 Переменные окружения

Некоторые настройки можно переопределить через переменные окружения:

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `FIREWALL_CONFIG` | Путь к конфигурационному файлу | `./firewall.json` |
| `FIREWALL_LOG` | Путь к файлу логов | `./firewall.log` |
| `FIREWALL_LISTEN_PORT` | Порт для входящих запросов | `8080` |
| `FIREWALL_ADMIN_PORT` | Порт админ-панели | `9090` |
| `FIREWALL_TARGET_PORT` | Порт целевого приложения | `3000` |

### Пример использования

\`\`\`bash
# Linux/macOS
export FIREWALL_LISTEN_PORT=80
export FIREWALL_ADMIN_PORT=8080
./firewall

# Windows
set FIREWALL_LISTEN_PORT=80
set FIREWALL_ADMIN_PORT=8080
firewall.exe
\`\`\`

## 🔄 Применение изменений

### Через админ-панель

Большинство изменений применяются автоматически при сохранении в админ-панели.

### Через конфигурационный файл

1. Отредактируйте `firewall.json`
2. Перезапустите firewall:

\`\`\`bash
# Если запущен как процесс
pkill firewall
./firewall

# Если запущен как сервис
sudo systemctl restart go-simple-firewall
\`\`\`

### Горячая перезагрузка

Некоторые настройки применяются без перезапуска:

- Rate limiting
- IP списки
- Настройки безопасности
- Включение/выключение модулей

## ✅ Валидация конфигурации

Firewall автоматически проверяет корректность конфигурации при запуске:

- Валидность JSON
- Корректность портов (1-65535)
- Валидность IP адресов
- Корректность временных интервалов

При обнаружении ошибок firewall:
1. Выводит сообщение об ошибке
2. Использует значения по умолчанию
3. Создает корректный конфигурационный файл

## 🔍 Мониторинг конфигурации

### Проверка текущих настроек

\`\`\`bash
# Просмотр конфигурации
cat firewall.json | jq '.'

# Проверка статуса через API
curl http://localhost:9090/admin/api/summary
\`\`\`

### Резервное копирование

\`\`\`bash
# Создание резервной копии
cp firewall.json firewall.json.backup.$(date +%Y%m%d_%H%M%S)

# Автоматическое резервное копирование (cron)
0 2 * * * cp /path/to/firewall.json /backup/firewall.json.$(date +\%Y\%m\%d)
\`\`\`

## 🚨 Устранение проблем конфигурации

### Firewall не запускается

1. Проверьте синтаксис JSON:
\`\`\`bash
cat firewall.json | jq '.'
\`\`\`

2. Проверьте права доступа:
\`\`\`bash
ls -la firewall.json
chmod 644 firewall.json
\`\`\`

3. Проверьте доступность портов:
\`\`\`bash
netstat -tlnp | grep :8080
\`\`\`

### Настройки не применяются

1. Убедитесь, что файл сохранен
2. Проверьте права на запись
3. Перезапустите firewall
4. Проверьте логи на ошибки

### Высокое потребление ресурсов

1. Уменьшите rate limit
2. Отключите ненужные модули защиты
3. Настройте ротацию логов
4. Очистите временные баны

