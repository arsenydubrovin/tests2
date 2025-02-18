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

# План тестирования

| Название | Тип | Входные данные | Ожидается |
|---|---|---|---|
| valid log entry | Позитивный | `message = "test message"` | `err = nil`, `entry.message = "test message"` |
| empty message log entry | Негативный | `message = ""` | `err != nil` |
| empty logger JSON | Позитивный | `entries = []` | `json = "[]"`, `err = nil` |
| single log entry JSON | Позитивный | `entries = [{message = "test", level = INFO, timestamp = "2024-01-01T12:00:00Z"}]` | `json = '[{"message":"test","level":"INFO","timestamp":"2024-01-01T12:00:00Z"}]'`, `err = nil` |
| filter WARNING and ERROR | Позитивный | `entries = [INFO, WARNING, ERROR]` | `filtered = [WARNING, ERROR]` |
| logger attestation | Позитивный | `entries = [INFO, WARNING, ERROR], timestamps = [2024-02-17T15:00:00Z, 2024-02-17T16:00:00Z, 2024-02-17T17:00:00Z]` | `filtered = [WARNING, ERROR], json соответствует ожидаемому`, `err = nil` |

| CollapseDuplicates | Входные данные (logs) | from | to | Ожидаемый результат | Ожидаемая ошибка |
|-----------------------------|------------------------|------|----|----------------------|------------------|
| неправильный диапазон | `[{“test”, INFO, fixedTime}]` | `fixedTime + 1ч` | `fixedTime` | `nil` | `"from >= to!"` |
| пустые логи | `[]` | `fixedTime` | `fixedTime + 1ч` | `[]` | `-` |
| нет логов в диапазоне | `[{“test”, INFO, fixedTime - 2ч}, {“test”, INFO, fixedTime - 1ч}]` | `fixedTime` | `fixedTime + 1ч` | `[]` | `-` |
| меньше трёх логов с одинаковым интервалом | `[{“test”, INFO, fixedTime}, {“test”, INFO, fixedTime + 10с}]` | `fixedTime` | `fixedTime + 1ч` | `[{“test”, INFO, fixedTime}, {“test”, INFO, fixedTime + 10с}]` | `-` |
| одинаковые интервалы | `[{“test”, INFO, fixedTime}, {“test”, INFO, fixedTime + 10с}, {“test”, INFO, fixedTime + 20с}, {“test”, INFO, fixedTime + 30с}]` | `fixedTime` | `fixedTime + 1ч` | `[{“test”, INFO, fixedTime}, {“test”, INFO, fixedTime + 30с}]` | `-` |
| разные интервалы | `[{“test”, INFO, fixedTime}, {“test”, INFO, fixedTime + 10с}, {“test”, INFO, fixedTime + 40с}, {“test”, INFO, fixedTime + 50с}]` | `fixedTime` | `fixedTime + 1ч` | `[{“test”, INFO, fixedTime}, {“test”, INFO, fixedTime + 10с}, {“test”, INFO, fixedTime + 40с}, {“test”, INFO, fixedTime + 50с}]` | `-` |
| смешанные логи | `[{“test1”, INFO, fixedTime}, {“test1”, INFO, fixedTime + 10с}, {“test1”, INFO, fixedTime + 20с}, {“test2”, INFO, fixedTime + 5с}, {“test2”, INFO, fixedTime + 15с}]` | `fixedTime` | `fixedTime + 1ч` | `[{“test1”, INFO, fixedTime}, {“test1”, INFO, fixedTime + 20с}, {“test2”, INFO, fixedTime + 5с}, {“test2”, INFO, fixedTime + 15с}]` | `-` |
| одинаковые сообщения, но разные уровни | `[{“test”, INFO, fixedTime}, {“test”, INFO, fixedTime + 10с}, {“test”, INFO, fixedTime + 20с}, {“test”, ERROR, fixedTime + 10с}, {“test”, ERROR, fixedTime + 20с}]` | `fixedTime` | `fixedTime + 1ч` | `[{“test”, INFO, fixedTime}, {“test”, INFO, fixedTime + 20с}, {“test”, ERROR, fixedTime + 10с}, {“test”, ERROR, fixedTime + 20с}]` | `-` |
| логи с нарушенной хронологией | `[{“test”, INFO, fixedTime + 20с}, {“test”, INFO, fixedTime + 10с}, {“test”, INFO, fixedTime}]` | `fixedTime` | `fixedTime + 1ч` | `nil` | `"inconsistent timestamps"` |
| границы диапазона и логов совпадают | `[{“test”, INFO, fixedTime}, {“test”, INFO, fixedTime + 30м}, {“test”, INFO, fixedTime + 1ч}]` | `fixedTime` | `fixedTime + 1ч` | `[{“test”, INFO, fixedTime}, {“test”, INFO, fixedTime + 1ч}]` | `-` |
| интервалы отличаются, но меньше, чем на дельту | `[{“test”, INFO, fixedTime + 10с}, {“test”, INFO, fixedTime + 20с}, {“test”, INFO, fixedTime + 34с}]` | `fixedTime` | `fixedTime + 1ч` | `[{“test”, INFO, fixedTime + 10с}, {“test”, INFO, fixedTime + 34с}]` | `-` |
| интервалы отличаются больше, чем на дельту | `[{“test”, INFO, fixedTime + 10с}, {“test”, INFO, fixedTime + 20с}, {“test”, INFO, fixedTime + 36с}]` | `fixedTime` | `fixedTime + 1ч` | `[{“test”, INFO, fixedTime + 10с}, {“test”, INFO, fixedTime + 20с}, {“test”, INFO, fixedTime + 36с}]` | `-` |
| одинаковые ts у логов | `[{“test”, INFO, fixedTime + 10с}, {“test”, INFO, fixedTime + 10с}, {“test”, INFO, fixedTime + 10с}]` | `fixedTime` | `fixedTime + 1ч` | `[{“test”, INFO, fixedTime + 10с}]` | `-` |
| одинковые логи, один нарушает порядок | `[{“test”, INFO, fixedTime + 10с}, {“test”, INFO, fixedTime + 11с}, {“test”, INFO, fixedTime + 10с}, {“test”, INFO, fixedTime + 10с}]` | `fixedTime` | `fixedTime + 1ч` | `nil` | `"inconsistent timestamps"` |

| inRange | ts | from | to | Ожидаемый результат |
|-------------------|----|------|----|----------------------|
| одинаковые границы | `fixedTime` | `fixedTime` | `fixedTime + 1ч` | `true` |
| ts равен верхней границе | `fixedTime + 1ч` | `fixedTime` | `fixedTime + 1ч` | `true` |
| ts равен нижней границе | `fixedTime` | `fixedTime` | `fixedTime + 1ч` | `true` |
| ts в диапазоне | `fixedTime + 30м` | `fixedTime` | `fixedTime + 1ч` | `true` |
| ts меньше нижней границы | `fixedTime - 1ч` | `fixedTime` | `fixedTime + 1ч` | `false` |
| ts больше верхней границы | `fixedTime + 2ч` | `fixedTime` | `fixedTime + 1ч` | `false` |

| hasPeriodicity | intervals | Ожидаемый результат |
|---------------------------|-----------|----------------------|
| пустой список | `[]` | `true` |
| один интервал | `[10с]` | `true` |
| одинаковый интервал | `[10с, 10с, 10с]` | `true` |
| интервалы отличаются, но меньше дельты | `[10с, 12с, 11с]` | `true` |
| интервалы отличаются больше дельты | `[10с, 16с, 10с, 16с]` | `false` |

| getIntervals | logs | Ожидаемый результат | Ожидаемая ошибка |
|-------------------------|------|----------------------|------------------|
| правильные интервалы | `[{fixedTime}, {fixedTime + 10с}, {fixedTime + 20с}]` | `[10с, 10с]` | `-` |
| нарушение хронологии | `[{fixedTime + 20с}, {fixedTime + 10с}]` | `nil` | `"inconsistent timestamps"` |
| один лог | `[{fixedTime}]` | `[]` | `-` |
| пустой список | `[]` | `[]` | `-` |
