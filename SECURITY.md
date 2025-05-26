# 🔒 Настройки безопасности Go Simple Firewall

Подробное руководство по настройке всех модулей безопасности и защиты от различных типов атак.

## 📋 Содержание

- [Обзор модулей безопасности](#обзор-модулей-безопасности)
- [SQL Injection Protection](#sql-injection-protection)
- [XSS Protection](#xss-protection)
- [DDoS Protection](#ddos-protection)
- [Scanner Protection](#scanner-protection)
- [Bot Protection](#bot-protection)
- [Directory Protection](#directory-protection)
- [Suffix Protection](#suffix-protection)
- [Geo-blocking](#geo-blocking)
- [Rate Limiting](#rate-limiting)
- [IP Management](#ip-management)
- [Мониторинг безопасности](#мониторинг-безопасности)
- [Лучшие практики](#лучшие-практики)

## 🛡️ Обзор модулей безопасности

Go Simple Firewall предоставляет многоуровневую защиту от различных типов атак:

| Модуль | Описание | Тип блокировки | По умолчанию |
|--------|----------|----------------|--------------|
| SQL Injection | Защита от SQL инъекций | Мгновенная | ✅ Включен |
| XSS Protection | Защита от XSS атак | Мгновенная | ✅ Включен |
| DDoS Protection | Защита от DDoS атак | Временный бан | ✅ Включен |
| Scanner Protection | Защита от сканеров | Мгновенная | ✅ Включен |
| Bot Protection | Защита от ботов | Мгновенная | ✅ Включен |
| Directory Protection | Защита директорий | Мгновенная | ✅ Включен |
| Suffix Protection | Защита от опасных файлов | Временный бан | ✅ Включен |
| Geo-blocking | Блокировка по странам | Мгновенная | ❌ Выключен |

## 💉 SQL Injection Protection

### Описание

Модуль анализирует все входящие запросы (GET и POST параметры) на наличие SQL ключевых слов и паттернов, характерных для SQL инъекций.

### Конфигурация

```json
{
  "security": {
    "enable_sql_protection": true,
    "sql_keywords": [
      "union", "select", "insert", "delete", "update",
      "drop", "create", "alter", "exec", "script",
      "declare", "cast", "convert", "having", "where",
      "order", "group", "by", "from", "into", "values",
      "information_schema", "sysobjects", "syscolumns"
    ]
  }
}
```

### Примеры блокируемых запросов

```
GET /search?q=1' UNION SELECT * FROM users--
POST /login with data: username=admin' OR '1'='1
GET /product?id=1; DROP TABLE products;
```

### Настройка через админ-панель

1. Откройте раздел "Security Protection"
2. Найдите "SQL Injection Protection"
3. Включите/выключите модуль
4. Отредактируйте список SQL ключевых слов
5. Нажмите "Update Security Settings"

### Рекомендуемые ключевые слова

**Базовые SQL команды:**
```
union, select, insert, delete, update, drop, create, alter
```

**Функции и операторы:**
```
exec, execute, cast, convert, concat, substring, ascii, char
```

**Системные таблицы:**
```
information_schema, sysobjects, syscolumns, pg_tables, sqlite_master
```

**Комментарии и завершители:**
```
--, /*, */, #, ;
```

## 🎭 XSS Protection

### Описание

Защищает от Cross-Site Scripting атак, анализируя параметры запросов на наличие JavaScript кода и HTML тегов.

### Конфигурация

```json
{
  "security": {
    "enable_xss_protection": true,
    "xss_patterns": [
      "<script", "javascript:", "onload=", "onerror=",
      "onclick=", "onmouseover=", "onfocus=", "onblur=",
      "eval(", "alert(", "confirm(", "prompt(",
      "document.cookie", "window.location", "innerHTML"
    ]
  }
}
```

### Примеры блокируемых запросов

```
GET /search?q=<script>alert('XSS')</script>
POST /comment with data: text=<img src=x onerror=alert(1)>
GET /redirect?url=javascript:alert(document.cookie)
```

### Расширенные XSS паттерны

**HTML события:**
```
onload, onerror, onclick, onmouseover, onfocus, onblur, onchange, onsubmit
```

**JavaScript функции:**
```
eval, setTimeout, setInterval, Function, alert, confirm, prompt
```

**DOM манипуляции:**
```
document.write, innerHTML, outerHTML, insertAdjacentHTML
```

**Протоколы:**
```
javascript:, data:, vbscript:, livescript:
```

## 🌊 DDoS Protection

### Описание

Интеллектуальная защита от DDoS атак, отслеживающая количество запросов с каждого IP адреса в определенном временном окне.

### Конфигурация

```json
{
  "security": {
    "enable_ddos_protection": true,
    "ddos_threshold": 100,
    "ddos_time_window": 60,
    "ddos_ban_duration": 30
  }
}
```

### Параметры

| Параметр | Описание | Рекомендуемые значения |
|----------|----------|------------------------|
| `ddos_threshold` | Максимум запросов за период | 50-200 |
| `ddos_time_window` | Временное окно (секунды) | 30-120 |
| `ddos_ban_duration` | Длительность бана (минуты) | 15-60 |

### Рекомендации по настройке

**Для статических сайтов:**
```json
{
  "ddos_threshold": 50,
  "ddos_time_window": 60,
  "ddos_ban_duration": 30
}
```

**Для API сервисов:**
```json
{
  "ddos_threshold": 200,
  "ddos_time_window": 30,
  "ddos_ban_duration": 15
}
```

**Для высоконагруженных приложений:**
```json
{
  "ddos_threshold": 500,
  "ddos_time_window": 60,
  "ddos_ban_duration": 10
}
```

### Алгоритм работы

1. Firewall отслеживает все запросы от каждого IP
2. Подсчитывает количество запросов в скользящем временном окне
3. При превышении порога создает временный бан
4. Автоматически снимает бан по истечении времени

## 🔍 Scanner Protection

### Описание

Защищает от автоматических сканеров уязвимостей, блокируя доступ к часто сканируемым путям.

### Конфигурация

```json
{
  "security": {
    "enable_scanner_protection": true,
    "scanner_paths": [
      "/admin", "/wp-admin", "/phpmyadmin", "/cpanel",
      "/webmail", "/.env", "/config", "/backup",
      "/test", "/demo", "/staging", "/dev",
      "/api/v1", "/api/v2", "/rest", "/graphql"
    ]
  }
}
```

### Популярные пути сканеров

**Административные панели:**
```
/admin, /administrator, /wp-admin, /cpanel, /plesk, /webmin
```

**Файлы конфигурации:**
```
/.env, /config, /.git, /.svn, /web.config, /.htaccess
```

**Базы данных:**
```
/phpmyadmin, /adminer, /mysql, /postgresql, /mongodb
```

**Резервные копии:**
```
/backup, /backups, /dump, /sql, /database
```

**Тестовые окружения:**
```
/test, /testing, /demo, /staging, /dev, /development
```

### Настройка для конкретных CMS

**WordPress:**
```
/wp-admin, /wp-login.php, /wp-config.php, /wp-content/uploads,
/xmlrpc.php, /wp-json, /readme.html, /license.txt
```

**Drupal:**
```
/admin, /user, /node, /sites/default, /modules, /themes,
/install.php, /update.php, /cron.php
```

**Joomla:**
```
/administrator, /installation, /configuration.php, /htaccess.txt,
/web.config.txt, /joomla.xml
```

## 🤖 Bot Protection

### Описание

Фильтрует запросы от вредоносных ботов и автоматических инструментов на основе анализа User-Agent заголовков.

### Конфигурация

```json
{
  "security": {
    "enable_bot_protection": true,
    "suspicious_user_agents": [
      "bot", "crawler", "spider", "scraper", "scanner",
      "nikto", "sqlmap", "nmap", "masscan", "zap",
      "curl", "wget", "python-requests", "go-http-client"
    ]
  }
}
```

### Категории ботов

**Сканеры уязвимостей:**
```
nikto, nessus, openvas, acunetix, burpsuite, zap, w3af
```

**SQL инъекция инструменты:**
```
sqlmap, havij, pangolin, safe3si, bsqlbf
```

**Сетевые сканеры:**
```
nmap, masscan, zmap, unicornscan, hping
```

**Автоматические инструменты:**
```
curl, wget, python-requests, go-http-client, java, perl
```

**Вредоносные боты:**
```
semrushbot, ahrefsbot, mj12bot, dotbot, blexbot
```

### Исключения для легитимных ботов

Если вам нужно разрешить определенных ботов (например, поисковых), используйте whitelist:

```json
{
  "allowed_user_agents": [
    "Googlebot", "Bingbot", "Slurp", "DuckDuckBot",
    "Baiduspider", "YandexBot", "facebookexternalhit"
  ]
}
```

## 📁 Directory Protection

### Описание

Блокирует доступ к системным и конфиденциальным директориям, которые не должны быть доступны через веб.

### Конфигурация

```json
{
  "security": {
    "enable_directory_protection": true,
    "protected_directories": [
      "/.git", "/.svn", "/backup", "/config",
      "/logs", "/tmp", "/.env", "/node_modules",
      "/.vscode", "/.idea", "/vendor", "/storage"
    ]
  }
}
```

### Системы контроля версий

```
/.git, /.svn, /.hg, /.bzr, /CVS, /.gitignore, /.gitmodules
```

### Конфигурационные файлы

```
/.env, /config, /.config, /settings, /conf, /etc
```

### Временные и служебные директории

```
/tmp, /temp, /cache, /logs, /log, /var, /storage
```

### IDE и редакторы

```
/.vscode, /.idea, /.sublime-project, /.atom, /.brackets
```

### Зависимости и библиотеки

```
/node_modules, /vendor, /bower_components, /packages
```

## 📄 Suffix Protection

### Описание

Блокирует запросы к файлам с потенциально опасными расширениями и создает временные баны для нарушителей.

### Конфигурация

```json
{
  "security": {
    "enable_suffix_protection": true,
    "forbidden_suffixes": [
      ".php", ".asp", ".aspx", ".jsp", ".cgi",
      ".pl", ".py", ".rb", ".sh", ".bat"
    ],
    "suffix_ban_duration": 10
  }
}
```

### Категории опасных файлов

**Серверные скрипты:**
```
.php, .asp, .aspx, .jsp, .cgi, .pl, .py, .rb, .lua
```

**Исполняемые файлы:**
```
.exe, .bat, .cmd, .sh, .ps1, .vbs, .jar
```

**Конфигурационные файлы:**
```
.conf, .config, .ini, .cfg, .properties, .xml
```

**Базы данных:**
```
.sql, .db, .sqlite, .mdb, .accdb
```

### Настройка длительности бана

| Уровень безопасности | Длительность | Описание |
|---------------------|--------------|----------|
| Низкий | 1-2 часа | Для тестовых окружений |
| Средний | 6-12 часов | Для обычных сайтов |
| Высокий | 24-48 часов | Для критичных систем |
| Максимальный | 168 часов (неделя) | Для высокозащищенных систем |

## 🌍 Geo-blocking

### Описание

Блокирует запросы из определенных стран на основе IP адресов (требует внешний сервис геолокации).

### Конфигурация

```json
{
  "security": {
    "enable_geo_blocking": true,
    "blocked_countries": ["CN", "RU", "KP", "IR", "SY"]
  }
}
```

### ISO коды стран

**Часто блокируемые страны:**
```
CN - Китай
RU - Россия  
KP - Северная Корея
IR - Иран
SY - Сирия
VE - Венесуэла
CU - Куба
```

**Европейские страны:**
```
DE - Германия
FR - Франция
GB - Великобритания
IT - Италия
ES - Испания
```

### Реализация геоблокировки

Для полноценной работы геоблокировки рекомендуется интеграция с:

1. **MaxMind GeoIP2** - коммерческая база данных
2. **IP2Location** - альтернативный сервис
3. **GeoJS** - бесплатный API сервис

## ⏱️ Rate Limiting

### Описание

Ограничивает количество запросов с одного IP адреса в единицу времени.

### Конфигурация

```json
{
  "rate_limit_rps": 60
}
```

### Рекомендуемые значения

| Тип приложения | RPS | Обоснование |
|----------------|-----|-------------|
| Статический сайт | 30-60 | Низкая интерактивность |
| Блог/CMS | 60-120 | Средняя активность пользователей |
| E-commerce | 120-200 | Высокая активность покупателей |
| API сервис | 200-500 | Множественные запросы от клиентов |
| Микросервисы | 500+ | Межсервисное взаимодействие |

### Алгоритм работы

1. **Sliding Window** - скользящее временное окно
2. **Token Bucket** - ведро токенов для сглаживания пиков
3. **Fixed Window** - фиксированные временные интервалы

## 🏠 IP Management

### Whitelist (Белый список)

Всегда разрешенные IP адреса, которые обходят все проверки:

```json
{
  "allowed_ips": {
    "127.0.0.1": true,
    "::1": true,
    "192.168.1.0/24": true,
    "10.0.0.0/8": true
  }
}
```

### Blacklist (Черный список)

Постоянно заблокированные IP адреса:

```json
{
  "banned_ips": {
    "192.168.1.100": true,
    "203.0.113.0/24": true,
    "2001:db8::/32": true
  }
}
```

### Поддерживаемые форматы

- **Одиночный IP:** `192.168.1.1`
- **CIDR подсеть:** `192.168.1.0/24`
- **IPv6:** `2001:db8::1`
- **IPv6 подсеть:** `2001:db8::/32`

## 📊 Мониторинг безопасности

### Метрики безопасности

Firewall собирает следующие метрики:

```json
{
  "security_metrics": {
    "total_blocked": 1250,
    "sql_injection_attempts": 45,
    "xss_attempts": 23,
    "ddos_attacks": 8,
    "scanner_attempts": 156,
    "bot_requests": 89,
    "directory_access_attempts": 34,
    "malicious_file_requests": 12
  }
}
```

### Алерты и уведомления

Настройка уведомлений о критических событиях:

```bash
# Мониторинг логов в реальном времени
tail -f firewall.log | grep "ATTACK"

# Подсчет атак за последний час
grep "$(date '+%Y-%m-%d %H')" firewall.log | grep "ATTACK" | wc -l

# Топ атакующих IP
grep "ATTACK" firewall.log | awk '{print $3}' | sort | uniq -c | sort -nr | head -10
```

### Интеграция с системами мониторинга

**Prometheus метрики:**
```
firewall_requests_total{status="blocked"}
firewall_attacks_total{type="sql_injection"}
firewall_banned_ips_total
```

**Grafana дашборд:**
- График заблокированных запросов
- Топ атакующих IP
- Распределение типов атак
- Статус модулей безопасности

## 🛡️ Лучшие практики

### 1. Многоуровневая защита

Используйте несколько модулей одновременно:

```json
{
  "security": {
    "enable_sql_protection": true,
    "enable_xss_protection": true,
    "enable_ddos_protection": true,
    "enable_scanner_protection": true,
    "enable_bot_protection": true
  }
}
```

### 2. Настройка под тип приложения

**Для API сервисов:**
- Отключите suffix protection
- Увеличьте rate limit
- Включите строгую bot protection

**Для статических сайтов:**
- Включите все модули защиты
- Уменьшите rate limit
- Настройте агрессивную DDoS защиту

### 3. Регулярное обновление правил

```bash
# Еженедельное обновление списков
# Добавляйте новые SQL ключевые слова
# Обновляйте списки вредоносных User-Agent'ов
# Анализируйте логи на предмет новых угроз
```

### 4. Мониторинг и анализ

```bash
# Ежедневный анализ логов
grep "$(date '+%Y-%m-%d')" firewall.log | grep "BLOCKED" > daily_blocks.log

# Еженедельный отчет по безопасности
awk '/ATTACK/ {attacks[$4]++} END {for (type in attacks) print type, attacks[type]}' firewall.log
```

### 5. Тестирование защиты

```bash
# Тест SQL injection защиты
curl "http://localhost:8080/search?q=1' UNION SELECT * FROM users--"

# Тест XSS защиты  
curl "http://localhost:8080/comment" -d "text=<script>alert('xss')</script>"

# Тест DDoS защиты
for i in {1..200}; do curl http://localhost:8080/ & done
```

### 6. Резервное копирование конфигурации

```bash
# Ежедневное резервное копирование
cp firewall.json /backup/firewall-$(date +%Y%m%d).json

# Версионирование в Git
git add firewall.json
git commit -m "Security config update $(date)"
```

### 7. Документирование изменений

Ведите журнал изменений безопасности:

```
2024-01-15: Добавлены новые SQL ключевые слова
2024-01-14: Увеличен DDoS threshold до 150
2024-01-13: Заблокированы IP из диапазона 203.0.113.0/24
```

### 8. Обучение команды

- Регулярно обучайте команду новым угрозам
- Проводите учения по реагированию на инциденты
- Документируйте процедуры реагирования

## 🚨 Реагирование на инциденты

### При обнаружении атаки

1. **Немедленные действия:**
   - Заблокируйте атакующий IP
   - Увеличьте уровень логирования
   - Уведомите команду безопасности

2. **Анализ:**
   - Изучите логи атаки
   - Определите тип и масштаб угрозы
   - Оцените потенциальный ущерб

3. **Устранение:**
   - Обновите правила защиты
   - Закройте обнаруженные уязвимости
   - Усильте мониторинг

### Автоматическое реагирование

```bash
#!/bin/bash
# Скрипт автоматического реагирования

# Проверка количества атак за последний час
ATTACKS=$(grep "$(date '+%Y-%m-%d %H')" firewall.log | grep "ATTACK" | wc -l)

if [ $ATTACKS -gt 50 ]; then
    # Увеличиваем защиту
    curl -X POST http://localhost:9090/admin \
         -d "action=update_security&ddos_threshold=50"
    
    # Отправляем уведомление
    echo "High attack volume detected: $ATTACKS attacks" | \
         mail -s "Security Alert" admin@example.com
fi
```

---

Эта документация поможет вам настроить максимально эффективную защиту для вашего приложения. Регулярно обновляйте правила безопасности и следите за новыми угрозами!

