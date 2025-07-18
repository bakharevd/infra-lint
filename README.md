TO-DO:
- Исправление ковычек, скобок, отступов


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