[![Build Status](https://github.com/arsenydubrovin/tests2/actions/workflows/ci.yml/badge.svg)](https://github.com/arsenydubrovin/tests2/actions)

## Типы и Константы

### `type Level int`
Тип, представляющий уровень логирования.

#### Константы:
- `INFO` – Информационное сообщение.
- `WARNING` – Предупреждение.
- `ERROR` – Ошибка.

#### `func (l Level) String() string`
Возвращает строковое представление уровня логирования.

## Структуры

### `type LogEntry struct`
Запись лога.

#### Поля:
- `message string` – Текст сообщения.
- `level Level` – Уровень логирования.
- `timestamp time.Time` – Временная метка.

#### `func NewLogEntry(message string) (*LogEntry, error)`
Создает новую запись лога с уровнем `INFO`.
Возвращает ошибку, если сообщение пустое.

#### `func (l *LogEntry) SetLevel(level Level)`
Устанавливает уровень логирования для записи.

#### `func (l *LogEntry) SetTimestamp(timestamp time.Time)`
Устанавливает временную метку для записи.

#### `func (l *LogEntry) String() string`
Возвращает строковое представление записи в формате: `[<timestamp>] <level>: <message>`

## Структуры

### `type Logger struct`
Логгер, управляющий записями лога.

#### Поля:
- `entries []*LogEntry` – Список записей.

#### `func NewLogger() *Logger`
Создает новый логгер.

#### `func (l *Logger) AddEntry(entry *LogEntry)`
Добавляет запись в логгер.
Если временная метка не установлена, устанавливает текущую.

#### `func (l *Logger) GetEntries(level ...Level) []*LogEntry`
Возвращает список записей с указанными уровнями.
Если уровни не указаны, возвращает все записи.

#### `func (l *Logger) MarshalJSON() ([]byte, error)`
Сериализует записи лога в формат JSON.
