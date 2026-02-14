# loglint

`loglint` - модульный плагин для `golangci-lint`, который проверяет текст сообщений в логах `log/slog` и `go.uber.org/zap`.

## Что проверяет

- `start-with-lowercase` (по умолчанию: `true`) - сообщение должно начинаться с маленькой английской буквы.
- `eng-letters` (по умолчанию: `true`) - в сообщении разрешены только английские буквы.
- `no-special-symbols` (по умолчанию: `true`) - запрет специальных символов.
- `allowed-symbols` (по умолчанию: пусто) - символы, которые разрешены дополнительно.
- `sensitive-keywords` (по умолчанию: пусто) - список чувствительных слов (без учета регистра).
- `sensitive-patterns` (по умолчанию: пусто) - список regex-паттернов для поиска чувствительных данных.

## Установка

1. Установите `golangci-lint` v2:

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

2. Создайте файл `.custom-gcl.yml` рядом с проектом:

```yaml
version: v2.9.0
plugins:
  - module: 'github.com/ddovwll/loglint'
    import: 'github.com/ddovwll/loglint/plugin'
    version: latest
```

3. Соберите исполняемый файл:

```bash
golangci-lint custom
```

После этого появится файл `custom-gcl`

## Использование

1. Добавьте `loglint` в `.golangci.yml` вашего Go-проекта.
2. Запускайте линтер через полученный исполняемый файл:

```bash
./custom-gcl run ./...
```

Запуск для применения автоисправлений
```bash
./custom-gcl run ./... --fix
```

## Пример golangci конфига

Пример `.golangci.yml` с подключенным `loglint` в файле `.golangci.example.yml`