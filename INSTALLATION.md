# 📦 Установка Go Simple Firewall

Это подробное руководство по установке и настройке Go Simple Firewall на различных операционных системах.

## 📋 Содержание

- [Системные требования](#системные-требования)
- [Установка из исходного кода](#установка-из-исходного-кода)
- [Установка на Linux](#установка-на-linux)
- [Установка на Windows](#установка-на-windows)
- [Установка на macOS](#установка-на-macos)
- [Docker установка](#docker-установка)
- [Настройка как сервис](#настройка-как-сервис)
- [Первоначальная настройка](#первоначальная-настройка)
- [Устранение неполадок](#устранение-неполадок)

## 🖥️ Системные требования

### Минимальные требования:
- **CPU**: 1 ядро, 1 GHz
- **RAM**: 512 MB
- **Диск**: 100 MB свободного места
- **ОС**: Linux, Windows 10+, macOS 10.14+

### Рекомендуемые требования:
- **CPU**: 2+ ядра, 2+ GHz
- **RAM**: 2+ GB
- **Диск**: 1+ GB свободного места
- **Сеть**: Стабильное интернет-соединение

### Программное обеспечение:
- **Go**: версия 1.21 или выше
- **Git**: для клонирования репозитория
- **Права администратора**: для установки как системный сервис

## 🔧 Установка из исходного кода

### 1. Установка Go

#### Linux (Ubuntu/Debian):
\`\`\`bash
# Обновляем систему
sudo apt update && sudo apt upgrade -y

# Устанавливаем Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Добавляем в PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Проверяем установку
go version
\`\`\`

#### Linux (CentOS/RHEL):
\`\`\`bash
# Устанавливаем Go
sudo dnf install golang -y
# или для старых версий: sudo yum install golang -y

# Проверяем установку
go version
\`\`\`

#### Windows:
1. Скачайте Go с [официального сайта](https://golang.org/dl/)
2. Запустите установщик и следуйте инструкциям
3. Откройте Command Prompt и проверьте: `go version`

#### macOS:
\`\`\`bash
# Используя Homebrew
brew install go

# Или скачайте с официального сайта
# https://golang.org/dl/

# Проверяем установку
go version
\`\`\`

### 2. Клонирование и сборка

\`\`\`bash
# Клонируем репозиторий
git clone https://github.com/hvckfy/go-simple-firewall.git
cd go-simple-firewall

# Загружаем зависимости
go mod tidy

# Собираем проект
go build -o firewall cmd/firewall/main.go

# Проверяем сборку
./firewall --help
\`\`\`

## 🐧 Установка на Linux

### Ubuntu/Debian

\`\`\`bash
# 1. Обновляем систему
sudo apt update && sudo apt upgrade -y

# 2. Устанавливаем необходимые пакеты
sudo apt install git golang-go build-essential -y

# 3. Клонируем и собираем
git clone https://github.com/hvckfy/go-simple-firewall.git
cd go-simple-firewall
go mod tidy
go build -o firewall cmd/firewall/main.go

# 4. Создаем директорию для установки
sudo mkdir -p /opt/go-simple-firewall
sudo cp firewall /opt/go-simple-firewall/
sudo chmod +x /opt/go-simple-firewall/firewall

# 5. Создаем символическую ссылку
sudo ln -sf /opt/go-simple-firewall/firewall /usr/local/bin/firewall

# 6. Проверяем установку
firewall --help
\`\`\`

### CentOS/RHEL/Fedora

\`\`\`bash
# 1. Обновляем систему
sudo dnf update -y  # или yum update -y для старых версий

# 2. Устанавливаем необходимые пакеты
sudo dnf install git golang gcc -y

# 3. Клонируем и собираем
git clone https://github.com/hvckfy/go-simple-firewall.git
cd go-simple-firewall
go mod tidy
go build -o firewall cmd/firewall/main.go

# 4. Устанавливаем
sudo mkdir -p /opt/go-simple-firewall
sudo cp firewall /opt/go-simple-firewall/
sudo chmod +x /opt/go-simple-firewall/firewall
sudo ln -sf /opt/go-simple-firewall/firewall /usr/local/bin/firewall
\`\`\`

### Arch Linux

\`\`\`bash
# 1. Обновляем систему
sudo pacman -Syu

# 2. Устанавливаем необходимые пакеты
sudo pacman -S git go gcc

# 3. Клонируем и собираем
git clone https://github.com/hvckfy/go-simple-firewall.git
cd go-simple-firewall
go mod tidy
go build -o firewall cmd/firewall/main.go

# 4. Устанавливаем
sudo mkdir -p /opt/go-simple-firewall
sudo cp firewall /opt/go-simple-firewall/
sudo chmod +x /opt/go-simple-firewall/firewall
sudo ln -sf /opt/go-simple-firewall/firewall /usr/local/bin/firewall
\`\`\`

## 🪟 Установка на Windows

### Метод 1: Установка через Git Bash

\`\`\`bash
# 1. Установите Git и Go с официальных сайтов
# Git: https://git-scm.com/download/win
# Go: https://golang.org/dl/

# 2. Откройте Git Bash и клонируйте репозиторий
git clone https://github.com/hvckfy/go-simple-firewall.git
cd go-simple-firewall

# 3. Соберите проект
go mod tidy
go build -o firewall.exe cmd/firewall/main.go

# 4. Проверьте сборку
./firewall.exe --help
\`\`\`

### Метод 2: Установка через PowerShell

\`\`\`powershell
# 1. Клонируем репозиторий
git clone https://github.com/hvckfy/go-simple-firewall.git
Set-Location go-simple-firewall

# 2. Собираем проект
go mod tidy
go build -o firewall.exe cmd/firewall/main.go

# 3. Создаем директорию для установки
New-Item -ItemType Directory -Force -Path "C:\Program Files\GoSimpleFirewall"
Copy-Item firewall.exe "C:\Program Files\GoSimpleFirewall\"

# 4. Добавляем в PATH (требует прав администратора)
$env:PATH += ";C:\Program Files\GoSimpleFirewall"
[Environment]::SetEnvironmentVariable("PATH", $env:PATH, [EnvironmentVariableTarget]::Machine)
\`\`\`

## 🍎 Установка на macOS

### Метод 1: Используя Homebrew

\`\`\`bash
# 1. Устанавливаем Homebrew (если не установлен)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 2. Устанавливаем Go
brew install go git

# 3. Клонируем и собираем
git clone https://github.com/hvckfy/go-simple-firewall.git
cd go-simple-firewall
go mod tidy
go build -o firewall cmd/firewall/main.go

# 4. Устанавливаем
sudo mkdir -p /usr/local/bin
sudo cp firewall /usr/local/bin/
sudo chmod +x /usr/local/bin/firewall
\`\`\`

### Метод 2: Ручная установка

\`\`\`bash
# 1. Скачайте и установите Go с https://golang.org/dl/
# 2. Установите Xcode Command Line Tools
xcode-select --install

# 3. Клонируем и собираем
git clone https://github.com/hvckfy/go-simple-firewall.git
cd go-simple-firewall
go mod tidy
go build -o firewall cmd/firewall/main.go

# 4. Устанавливаем
sudo cp firewall /usr/local/bin/
sudo chmod +x /usr/local/bin/firewall
\`\`\`

## 🐳 Docker установка

### Создание Dockerfile

\`\`\`dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o firewall cmd/firewall/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/firewall .
EXPOSE 8080 9090
CMD ["./firewall"]
\`\`\`

### Сборка и запуск

\`\`\`bash
# Сборка образа
docker build -t go-simple-firewall .

# Запуск контейнера
docker run -d \
  --name firewall \
  -p 8080:8080 \
  -p 9090:9090 \
  -v $(pwd)/firewall.json:/root/firewall.json \
  -v $(pwd)/firewall.log:/root/firewall.log \
  go-simple-firewall

# Проверка статуса
docker ps
docker logs firewall
\`\`\`

### Docker Compose

\`\`\`yaml
version: '3.8'

services:
  firewall:
    build: .
    ports:
      - "8080:8080"
      - "9090:9090"
    volumes:
      - ./firewall.json:/root/firewall.json
      - ./firewall.log:/root/firewall.log
    restart: unless-stopped
    environment:
      - GO_ENV=production
\`\`\`

## ⚙️ Настройка как сервис

### Linux (systemd)

\`\`\`bash
# 1. Запустите firewall с правами администратора
sudo ./firewall

# 2. Откройте админ-панель: http://localhost:9090/admin
# 3. Создайте учетную запись администратора
# 4. В разделе "System Management" нажмите "Install as Service"

# 5. Проверьте статус сервиса
sudo systemctl status go-simple-firewall

# 6. Управление сервисом
sudo systemctl start go-simple-firewall
sudo systemctl stop go-simple-firewall
sudo systemctl restart go-simple-firewall
sudo systemctl enable go-simple-firewall  # автозапуск

# 7. Просмотр логов
sudo journalctl -u go-simple-firewall -f
\`\`\`

### Windows Service

\`\`\`cmd
# 1. Запустите Command Prompt от имени администратора
# 2. Перейдите в директорию с firewall.exe
# 3. Запустите firewall
firewall.exe

# 4. Откройте админ-панель: http://localhost:9090/admin
# 5. Создайте учетную запись администратора
# 6. В разделе "System Management" нажмите "Install as Service"

# 7. Управление сервисом через Services.msc или командную строку:
sc start GoSimpleFirewall
sc stop GoSimpleFirewall
sc query GoSimpleFirewall
\`\`\`

### macOS LaunchDaemon

\`\`\`bash
# 1. Запустите firewall с правами администратора
sudo ./firewall

# 2. Откройте админ-панель: http://localhost:9090/admin
# 3. Создайте учетную запись администратора
# 4. В разделе "System Management" нажмите "Install as Service"

# 5. Управление сервисом
sudo launchctl start com.gosimplefirewall.daemon
sudo launchctl stop com.gosimplefirewall.daemon
sudo launchctl list | grep gosimplefirewall
\`\`\`

## 🎯 Первоначальная настройка

### 1. Первый запуск

\`\`\`bash
# Запустите firewall
./firewall

# Вы увидите сообщения:
# 🔥 Firewall started on port 8080, proxying to 3000
# 📊 Admin panel: http://localhost:9090/admin
\`\`\`

### 2. Создание учетной записи администратора

1. Откройте браузер и перейдите по адресу: `http://localhost:9090/admin`
2. Вы увидите страницу первоначальной настройки
3. Введите желаемое имя пользователя и пароль
4. Нажмите "Create Admin Account"

### 3. Базовая конфигурация

После входа в админ-панель:

1. **Настройте порты**:
   - Listen Port: порт, на котором работает firewall (по умолчанию 8080)
   - Target Port: порт вашего приложения (по умолчанию 3000)
   - Admin Port: порт админ-панели (по умолчанию 9090)

2. **Настройте Rate Limiting**:
   - Установите лимит запросов в минуту (по умолчанию 60)

3. **Включите нужные модули защиты**:
   - SQL Injection Protection
   - XSS Protection
   - DDoS Protection
   - Scanner Protection
   - Bot Protection

### 4. Проверка работы

\`\`\`bash
# Проверьте, что firewall работает
curl http://localhost:8080

# Проверьте админ-панель
curl http://localhost:9090/admin
\`\`\`

## 🔧 Настройка конфигурационного файла

Firewall создает файл `firewall.json` с настройками:

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
  "security": {
    "enable_suffix_protection": true,
    "forbidden_suffixes": [".php", ".asp", ".aspx", ".jsp", ".cgi"],
    "suffix_ban_duration": 10,
    "enable_sql_protection": true,
    "sql_keywords": ["union", "select", "insert", "delete", "update"],
    "enable_xss_protection": true,
    "xss_patterns": ["<script", "javascript:", "onload="],
    "enable_scanner_protection": true,
    "scanner_paths": ["/admin", "/wp-admin", "/phpmyadmin"],
    "enable_bot_protection": true,
    "suspicious_user_agents": ["bot", "crawler", "spider"],
    "enable_directory_protection": true,
    "protected_directories": ["/.git", "/.svn", "/backup"],
    "enable_ddos_protection": true,
    "ddos_threshold": 100,
    "ddos_time_window": 60,
    "ddos_ban_duration": 30,
    "enable_geo_blocking": false,
    "blocked_countries": []
  },
  "temporary_bans": []
}
\`\`\`

## 🚨 Устранение неполадок

### Проблема: Порт уже используется

\`\`\`bash
# Найдите процесс, использующий порт
sudo netstat -tlnp | grep :8080
# или
sudo lsof -i :8080

# Завершите процесс
sudo kill -9 <PID>

# Или измените порт в конфигурации
\`\`\`

### Проблема: Нет прав доступа

\`\`\`bash
# Linux: запустите с sudo
sudo ./firewall

# Или измените владельца файла
sudo chown $USER:$USER firewall
chmod +x firewall
\`\`\`

### Проблема: Не удается установить как сервис

\`\`\`bash
# Linux: проверьте права и systemd
sudo systemctl --version
sudo systemctl daemon-reload

# Windows: запустите от имени администратора
# macOS: проверьте права sudo
\`\`\`

### Проблема: Firewall не блокирует атаки

1. Проверьте, что firewall включен в админ-панели
2. Убедитесь, что нужные модули защиты активированы
3. Проверьте логи: `tail -f firewall.log`
4. Убедитесь, что запросы проходят через firewall

### Проблема: Высокое потребление ресурсов

1. Уменьшите лимит rate limiting
2. Отключите ненужные модули защиты
3. Очистите старые логи
4. Проверьте количество временных банов

### Логи и диагностика

\`\`\`bash
# Просмотр логов firewall
tail -f firewall.log

# Просмотр системных логов (Linux)
sudo journalctl -u go-simple-firewall -f

# Проверка статуса процесса
ps aux | grep firewall

# Проверка сетевых соединений
sudo netstat -tlnp | grep firewall
\`\`\`

## 📞 Получение помощи

Если у вас возникли проблемы:

1. Проверьте [FAQ](https://github.com/hvckfy/go-simple-firewall/wiki/FAQ)
2. Поищите решение в [Issues](https://github.com/hvckfy/go-simple-firewall/issues)
3. Создайте новый [Issue](https://github.com/hvckfy/go-simple-firewall/issues/new) с подробным описанием проблемы
4. Присоединяйтесь к [Discussions](https://github.com/hvckfy/go-simple-firewall/discussions)

При создании Issue приложите:
- Версию ОС
- Версию Go
- Конфигурационный файл (без паролей)
- Логи с ошибкой
- Шаги для воспроизведения проблемы

