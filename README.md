# 🔥 Go Simple Firewall

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Windows%20%7C%20macOS-lightgrey?style=for-the-badge)](https://github.com/hvckfy/go-simple-firewall)

**Go Simple Firewall** — это мощный, легкий и простой в использовании веб-firewall с продвинутыми функциями безопасности, написанный на Go. Предоставляет комплексную защиту от различных типов атак с удобным веб-интерфейсом для управления.

![Firewall Dashboard](https://via.placeholder.com/800x400/2563eb/ffffff?text=Go+Simple+Firewall+Dashboard)

## ✨ Основные возможности

### 🛡️ Защита от атак
- **SQL Injection Protection** — Обнаружение и блокировка SQL-инъекций
- **XSS Protection** — Защита от Cross-Site Scripting атак
- **DDoS Protection** — Интеллектуальная защита от DDoS атак
- **Scanner Protection** — Блокировка сканеров уязвимостей
- **Bot Protection** — Фильтрация вредоносных ботов
- **Directory Traversal Protection** — Защита от обхода директорий
- **Geo-blocking** — Блокировка по географическому признаку

### 🚀 Производительность и управление
- **Rate Limiting** — Ограничение количества запросов
- **IP Whitelisting/Blacklisting** — Управление доступом по IP
- **User-Agent Filtering** — Фильтрация по User-Agent
- **Temporary Bans** — Система временных блокировок
- **Real-time Statistics** — Статистика в реальном времени
- **Detailed Logging** — Подробное логирование всех событий

### 🎛️ Удобное управление
- **Web Admin Panel** — Современный веб-интерфейс
- **Authentication System** — Безопасная система аутентификации
- **Service Integration** — Установка как системный сервис
- **JSON Configuration** — Простая конфигурация через JSON
- **Auto-restart** — Автоматический перезапуск при сбоях

## 🏗️ Архитектура

\`\`\`
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client        │───▶│  Go Firewall    │───▶│  Target Server  │
│                 │    │  (Port 8080)    │    │  (Port 3000)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │  Admin Panel    │
                       │  (Port 9090)    │
                       └─────────────────┘
\`\`\`

Firewall работает как reverse proxy, перехватывая все входящие запросы, анализируя их на предмет угроз и пропуская только безопасные запросы к целевому серверу.

## 🚀 Быстрый старт

### Предварительные требования

- Go 1.21 или выше
- Linux/Windows/macOS
- Права администратора (для установки как сервис)

### Установка

1. **Клонируйте репозиторий:**
\`\`\`bash
git clone https://github.com/hvckfy/go-simple-firewall.git
cd go-simple-firewall
\`\`\`

2. **Соберите проект:**
\`\`\`bash
go mod tidy
go build -o firewall cmd/firewall/main.go
\`\`\`

3. **Запустите firewall:**
\`\`\`bash
./firewall
\`\`\`

4. **Откройте админ-панель:**
   - Перейдите по адресу: `http://localhost:9090/admin`
   - При первом запуске создайте учетную запись администратора

## 📋 Подробная документация

- 📖 [Установка и настройка](INSTALLATION.md)
- ⚙️ [Конфигурация](CONFIGURATION.md)
- 🔒 [Настройки безопасности](SECURITY.md)
- 🔌 [API документация](API.md)
- 🤝 [Руководство для разработчиков](CONTRIBUTING.md)

## 🖥️ Скриншоты

### Главная панель управления
![Dashboard](https://via.placeholder.com/600x400/f8fafc/1e293b?text=System+Status+Dashboard)

### Статистика в реальном времени
![Statistics](https://via.placeholder.com/600x400/f8fafc/1e293b?text=Real-time+Statistics)

### Настройки безопасности
![Security Settings](https://via.placeholder.com/600x400/f8fafc/1e293b?text=Security+Configuration)

## 🔧 Конфигурация по умолчанию

\`\`\`json
{
  "listen_port": 8080,
  "admin_port": 9090,
  "target_port": 3000,
  "rate_limit_rps": 60,
  "enable_logging": true,
  "enable_firewall": true,
  "security": {
    "enable_suffix_protection": true,
    "enable_sql_protection": true,
    "enable_xss_protection": true,
    "enable_scanner_protection": true,
    "enable_bot_protection": true,
    "enable_directory_protection": true,
    "enable_ddos_protection": true,
    "ddos_threshold": 100,
    "ddos_time_window": 60,
    "ddos_ban_duration": 30
  }
}
\`\`\`

## 📊 Мониторинг и логирование

### Типы логов
- **Request Logs** — Все HTTP запросы с детальной информацией
- **Security Events** — События безопасности и атаки
- **Admin Actions** — Действия администратора
- **System Events** — Системные события и ошибки

### Статистика
- Почасовая статистика запросов
- Топ IP-адресов по активности
- Статистика по User-Agent
- Информация о заблокированных запросах

## 🛠️ Использование как сервис

### Linux (systemd)
\`\`\`bash
# Установка как сервис
sudo ./firewall
# В админ-панели нажмите "Install as Service"

# Управление сервисом
sudo systemctl start go-simple-firewall
sudo systemctl stop go-simple-firewall
sudo systemctl status go-simple-firewall
\`\`\`

### Windows
\`\`\`cmd
# Запуск от имени администратора
firewall.exe
# В админ-панели нажмите "Install as Service"

# Управление через Services.msc или sc
sc start GoSimpleFirewall
sc stop GoSimpleFirewall
\`\`\`

## 🔐 Безопасность

### Аутентификация
- Безопасная система сессий
- Хеширование паролей с bcrypt
- Защита от брутфорс атак
- Автоматическое истечение сессий

### Защищенные данные
- Пароли скрываются в логах
- Безопасные cookies
- CSRF защита
- Валидация всех входных данных

## 🚨 Типы атак, от которых защищает firewall

| Тип атаки | Описание | Действие |
|-----------|----------|----------|
| SQL Injection | Попытки внедрения SQL кода | Блокировка + лог |
| XSS | Cross-Site Scripting атаки | Блокировка + лог |
| DDoS | Превышение лимита запросов | Временный бан |
| Scanner | Сканирование уязвимостей | Блокировка + лог |
| Bot | Вредоносные боты | Блокировка + лог |
| Directory Traversal | Попытки обхода директорий | Блокировка + лог |
| Malicious Files | Запросы к опасным файлам | Временный бан |

## 📈 Производительность

- **Низкая задержка** — Минимальное влияние на время отклика
- **Высокая пропускная способность** — Обработка тысяч запросов в секунду
- **Эффективное использование памяти** — Оптимизированные структуры данных
- **Масштабируемость** — Легко адаптируется под нагрузку

## 🤝 Вклад в проект

Мы приветствуем вклад в развитие проекта! Пожалуйста, ознакомьтесь с [руководством для разработчиков](CONTRIBUTING.md).

### Как помочь проекту:
- 🐛 Сообщайте об ошибках
- 💡 Предлагайте новые функции
- 📝 Улучшайте документацию
- 🔧 Отправляйте Pull Request'ы

## 📄 Лицензия

Этот проект распространяется под лицензией MIT. Подробности в файле [LICENSE](LICENSE).

## 🙏 Благодарности

- Команде Go за отличный язык программирования
- Сообществу разработчиков за вдохновение
- Всем, кто тестирует и улучшает проект

## 📞 Поддержка

- 📧 **Email**: [your-email@example.com](mailto:your-email@example.com)
- 🐛 **Issues**: [GitHub Issues](https://github.com/hvckfy/go-simple-firewall/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/hvckfy/go-simple-firewall/discussions)

---

<div align="center">

**⭐ Если проект оказался полезным, поставьте звездочку! ⭐**

Made with ❤️ by [hvckfy](https://github.com/hvckfy)

</div>

