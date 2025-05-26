# 🤝 Руководство для разработчиков

Добро пожаловать в проект Go Simple Firewall! Мы рады вашему желанию внести вклад в развитие проекта.

## 📋 Содержание

- [Как внести вклад](#как-внести-вклад)
- [Настройка окружения разработки](#настройка-окружения-разработки)
- [Структура проекта](#структура-проекта)
- [Стандарты кодирования](#стандарты-кодирования)
- [Отправка изменений](#отправка-изменений)
- [Сообщение об ошибках](#сообщение-об-ошибках)
- [Предложение новых функций](#предложение-новых-функций)

## 🚀 Как внести вклад

Есть множество способов помочь проекту:

- 🐛 **Сообщайте об ошибках** - помогите нам найти и исправить баги
- 💡 **Предлагайте новые функции** - поделитесь идеями по улучшению
- 📝 **Улучшайте документацию** - помогите сделать документацию лучше
- 🔧 **Пишите код** - исправляйте баги и добавляйте новые функции
- 🌍 **Переводите** - добавьте поддержку новых языков

## 🛠️ Настройка окружения разработки

### Предварительные требования

- **Go 1.21+** - [Установка Go](https://golang.org/doc/install)
- **Git** - для работы с репозиторием

### Клонирование репозитория

```bash
# Форкните репозиторий на GitHub, затем клонируйте свой форк
git clone https://github.com/YOUR_USERNAME/go-simple-firewall.git
cd go-simple-firewall

# Добавьте оригинальный репозиторий как upstream
git remote add upstream https://github.com/hvckfy/go-simple-firewall.git
```

### Установка зависимостей

```bash
# Загрузите зависимости
go mod tidy

# Проверьте, что все работает
go build -o firewall cmd/firewall/main.go
./firewall --help
```

## 📁 Структура проекта

```
go-simple-firewall/
├── cmd/
│   └── firewall/
│       └── main.go              # Точка входа приложения
├── internal/                    # Внутренние пакеты
│   ├── admin/                   # Админ-панель
│   │   ├── admin.go
│   │   ├── template.go
│   │   └── templates.go
│   ├── auth/                    # Аутентификация
│   │   └── auth.go
│   ├── config/                  # Конфигурация
│   │   └── config.go
│   ├── ddos/                    # DDoS защита
│   │   └── ddos.go
│   ├── firewall/                # Основная логика firewall
│   │   └── firewall.go
│   ├── logger/                  # Логирование
│   │   └── logger.go
│   ├── ratelimit/               # Rate limiting
│   │   └── ratelimit.go
│   ├── security/                # Модули безопасности
│   │   └── security.go
│   └── stats/                   # Статистика
│       └── stats.go
├── pkg/                         # Публичные пакеты
│   ├── service/                 # Управление сервисами
│   │   └── service.go
│   └── utils/                   # Утилиты
│       └── utils.go
├── docs/                        # Документация
├── go.mod                       # Go модуль
├── go.sum                       # Контрольные суммы зависимостей
└── README.md                    # Основная документация
```

### Описание пакетов

| Пакет | Описание |
|-------|----------|
| `cmd/firewall` | Точка входа, инициализация приложения |
| `internal/admin` | Веб-интерфейс администратора |
| `internal/auth` | Система аутентификации и авторизации |
| `internal/config` | Управление конфигурацией |
| `internal/ddos` | Защита от DDoS атак |
| `internal/firewall` | Основная логика firewall |
| `internal/logger` | Система логирования |
| `internal/ratelimit` | Ограничение скорости запросов |
| `internal/security` | Модули безопасности |
| `internal/stats` | Сбор и анализ статистики |
| `pkg/service` | Управление системными сервисами |
| `pkg/utils` | Общие утилиты |

## 📝 Стандарты кодирования

### Go Code Style

Мы следуем стандартным соглашениям Go:

1. **Форматирование**: используйте `gofmt` или `goimports`
2. **Именование**: следуйте [Go naming conventions](https://golang.org/doc/effective_go.html#names)
3. **Комментарии**: документируйте публичные функции и типы
4. **Обработка ошибок**: всегда обрабатывайте ошибки явно

### Примеры хорошего кода

```go
// Package security provides various security modules for the firewall.
package security

import (
    "fmt"
    "net/http"
    "strings"
)

// SecurityChecker checks incoming requests for various security threats.
type SecurityChecker struct {
    config *config.Config
}

// New creates a new SecurityChecker instance.
func New(cfg *config.Config) *SecurityChecker {
    return &SecurityChecker{
        config: cfg,
    }
}

// CheckRequest analyzes the request for security threats.
// Returns true if the request should be blocked, along with the reason.
func (sc *SecurityChecker) CheckRequest(r *http.Request) (blocked bool, reason string) {
    if sc.config.Security.EnableSQLProtection {
        if blocked, reason := sc.checkSQLInjection(r); blocked {
            return true, reason
        }
    }
    
    return false, ""
}

// checkSQLInjection checks for SQL injection attempts in the request.
func (sc *SecurityChecker) checkSQLInjection(r *http.Request) (bool, string) {
    // Implementation details...
    return false, ""
}
```

## 📤 Отправка изменений

### Workflow

1. **Создайте ветку для вашей функции:**
```bash
git checkout -b feature/new-security-module
```

2. **Внесите изменения и зафиксируйте их:**
```bash
git add .
git commit -m "feat: add new security module for detecting XYZ attacks"
```

3. **Синхронизируйтесь с upstream:**
```bash
git fetch upstream
git rebase upstream/main
```

4. **Отправьте изменения в свой форк:**
```bash
git push origin feature/new-security-module
```

5. **Создайте Pull Request на GitHub**

### Формат коммитов

Мы используем [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Типы коммитов:**
- `feat:` - новая функция
- `fix:` - исправление бага
- `docs:` - изменения в документации
- `style:` - форматирование кода
- `refactor:` - рефакторинг кода
- `chore:` - обновление зависимостей, настройки

**Примеры:**
```
feat(security): add XSS protection module
fix(ratelimit): fix memory leak in rate limiter
docs(api): update API documentation for new endpoints
test(security): add unit tests for SQL injection detection
```

### Требования к Pull Request

1. **Описание**: четко опишите, что делает ваш PR
2. **Документация**: обновите документацию при необходимости
3. **Линтинг**: убедитесь, что код проходит все проверки
4. **Размер**: старайтесь делать PR небольшими и сфокусированными

### Шаблон Pull Request

```markdown
## Описание

Краткое описание изменений.

## Тип изменения

- [ ] Исправление бага
- [ ] Новая функция
- [ ] Критическое изменение
- [ ] Обновление документации

## Как протестировано

Опишите тесты, которые вы провели.

## Чеклист

- [ ] Код следует стандартам проекта
- [ ] Проведен self-review кода
- [ ] Код прокомментирован в сложных местах
- [ ] Внесены соответствующие изменения в документацию
- [ ] Изменения не генерируют новых предупреждений
- [ ] Добавлены тесты, которые доказывают, что исправление эффективно или функция работает
- [ ] Новые и существующие unit тесты проходят локально
```

## 🐛 Сообщение об ошибках

### Перед созданием Issue

1. Проверьте [существующие Issues](https://github.com/hvckfy/go-simple-firewall/issues)
2. Убедитесь, что используете последнюю версию
3. Попробуйте воспроизвести ошибку на чистой установке

### Шаблон Bug Report

```markdown
**Описание бага**
Четкое и краткое описание того, что является багом.

**Шаги для воспроизведения**
1. Перейти к '...'
2. Нажать на '....'
3. Прокрутить вниз до '....'
4. Увидеть ошибку

**Ожидаемое поведение**
Четкое и краткое описание того, что вы ожидали.

**Скриншоты**
Если применимо, добавьте скриншоты для объяснения проблемы.

**Окружение:**
- ОС: [например, Ubuntu 20.04]
- Версия Go: [например, 1.21.5]
- Версия Firewall: [например, 1.0.0]

**Дополнительный контекст**
Добавьте любой другой контекст о проблеме здесь.

**Логи**
```
Вставьте соответствующие логи здесь
```
```

## 💡 Предложение новых функций

### Шаблон Feature Request

```markdown
**Связана ли ваша функция с проблемой? Опишите.**
Четкое и краткое описание проблемы. Например: Я всегда расстраиваюсь, когда [...]

**Опишите решение, которое вы хотели бы**
Четкое и краткое описание того, что вы хотите.

**Опишите альтернативы, которые вы рассматривали**
Четкое и краткое описание любых альтернативных решений или функций, которые вы рассматривали.

**Дополнительный контекст**
Добавьте любой другой контекст или скриншоты о запросе функции здесь.
```

## 🏷️ Релизы

### Версионирование

Мы используем [Semantic Versioning](https://semver.org/):

- **MAJOR** версия для несовместимых изменений API
- **MINOR** версия для новой функциональности с обратной совместимостью
- **PATCH** версия для исправлений багов с обратной совместимостью

### Процесс релиза

1. Обновление версии в коде
2. Создание changelog
3. Создание Git тега
4. Сборка бинарных файлов
5. Публикация релиза на GitHub

## 📞 Связь с сообществом

- **GitHub Issues** - для багов и предложений функций
- **GitHub Discussions** - для общих вопросов и обсуждений
- **Email** - для приватных вопросов

## 🙏 Благодарности

Спасибо всем, кто вносит вклад в проект! Ваши усилия делают Go Simple Firewall лучше для всех.

### Список контрибьюторов

Все контрибьюторы будут добавлены в этот список:

- [@hvckfy](https://github.com/hvckfy) - создатель и основной разработчик

---

**Помните**: каждый вклад ценен, независимо от размера. Даже исправление опечатки в документации помогает проекту!
