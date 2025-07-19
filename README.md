dev

# infra-lint

Линтер для инфраструктурных файлов (Dockerfile, docker-compose.yml, nginx.conf, Jenkinsfile).

## Установка

### Из релизов GitHub

Скачайте последнюю версию для вашей платформы из [релизов](https://github.com/bakharevd/infra-lint/releases/latest):

```bash
# Linux (AMD64)
curl -L https://github.com/bakharevd/infra-lint/releases/latest/download/infra-lint-linux-amd64.tar.gz | tar xz
sudo mv infra-lint-linux-amd64 /usr/local/bin/infra-lint

# macOS (ARM64/M1)
curl -L https://github.com/bakharevd/infra-lint/releases/latest/download/infra-lint-darwin-arm64.tar.gz | tar xz
sudo mv infra-lint-darwin-arm64 /usr/local/bin/infra-lint

# Windows
# Скачайте infra-lint-windows-amd64.zip и распакуйте в PATH
```

### Из исходников

```bash
go install infra-lint/cmd@latest
```

## Использование

### Сканирование одного файла
```bash
infra-lint --file path/to/Dockerfile
```

### Сканирование директории
```bash
infra-lint --dir path/to/directory
```

### Сканирование git репозитория
```bash
infra-lint --repo https://github.com/user/repo
```

### Фильтрация линтеров

Показать список доступных линтеров:
```bash
infra-lint --list-types
```

Запустить только определенные линтеры:
```bash
infra-lint --dir . --include-types docker,nginx
```

Исключить определенные линтеры:
```bash
infra-lint --dir . --exclude-types jenkins,compose
```

### Дополнительные опции
- `--no-color` - отключить цветной вывод
- `--formatter` - форматировать файлы вместо линтинга

### Форматирование файлов

Форматирование одного файла:
```bash
infra-lint --file path/to/Dockerfile --formatter
```

Форматирование всех файлов в директории:
```bash
infra-lint --dir path/to/directory --formatter
```

Форматтер автоматически:
- Выравнивает отступы (4 пробела)
- Исправляет регистр команд в Dockerfile
- Форматирует YAML в docker-compose файлах
- Выравнивает блоки в nginx.conf
- Форматирует Jenkinsfile с правильными отступами

## Разработка

### Создание релиза

Релизы создаются автоматически при создании git тега:

```bash
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions автоматически:
1. Соберёт бинарные файлы для всех платформ
2. Создаст релиз на GitHub с архивами

### Непрерывная интеграция

При каждом пуше в `main` ветку:
- Запускаются тесты и проверки
- Создаётся "latest" релиз с последними изменениями
- Обновляется Docker образ с тегом `latest`