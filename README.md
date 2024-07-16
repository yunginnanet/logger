# logger

Convenience wrapper for [zerolog](https://github.com/rs/zerolog) and [zwrap](https://github.com/yunginnanet/zwrap).

`import "git.tcp.direct/kayos/logger"`

---

```go
const MaximumSyncErrors = 10000
```

#### func  CreateDatedLogFile

```go
func CreateDatedLogFile(directory, prefix string) (*os.File, error)
```

#### func  CreateDatedLogFileCtx

```go
func CreateDatedLogFileCtx(ctx context.Context, directory, prefix string) (*os.File, error)
```

#### func  CreateDatedLogFileFormatted

```go
func CreateDatedLogFileFormatted(directory, prefix string, format string) (*os.File, error)
```

#### func  CreateDatedLogFileFormattedCtx

```go
func CreateDatedLogFileFormattedCtx(ctx context.Context, directory, prefix string, format string) (*os.File, error)
```

#### func  DisableSyncErrorAccounting

```go
func DisableSyncErrorAccounting()
```
DisableSyncErrorAccounting disables storing and panicking upon exceeding the
maximum sync error constant with regard to log files created with
[CreateDatedLogfile] and managed by [StartPeriodicSync]. By default, this
behavior is enabled.

#### func  EnableSyncErrorAccounting

```go
func EnableSyncErrorAccounting()
```
EnableSyncErrorAccounting enables storing and panicking upon exceeding the
maximum sync error constant with regard to log files created with
[CreateDatedLogfile] and managed by [StartPeriodicSync]. By default, this
behavior is enabled.

#### func  StartPeriodicSync

```go
func StartPeriodicSync(ctx context.Context, f *os.File, dur time.Duration)
```

#### type Log

```go
type Log struct {
}
```


#### func  Global

```go
func Global() *Log
```
Global acquires the assigned global logger.

IMPORTANT: you MUST make your instance of [Log] globally accecible by calling
[WithGlobalPackageAccess].

#### func  NewLogger

```go
func NewLogger(writers ...io.Writer) *Log
```
NewLogger creates a logger that writes to the given writers, as well as pretty
prints to stdout.

#### func  NewLoggerNoColor

```go
func NewLoggerNoColor(writers ...io.Writer) *Log
```

#### func  NewQuietLogger

```go
func NewQuietLogger(writers ...io.Writer) *Log
```
NewQuietLogger creates a logger that writes to the given writers with no console
writer added.

#### func (*Log) AddWriter

```go
func (l *Log) AddWriter(w io.Writer)
```
AddWriter adds a writer to the logger.

Note: this may have unintended consequences if certain [zwrap.Logger]
configuration values have been set via [Log.C]. That said, the [zwrap.Logger]
prefix value will be preserved.

#### func (*Log) C

```go
func (l *Log) C() zwrap.ZWrapLogger
```
C returns a [zwrap.ZWrapLogger] which is a highly compattible interface to fit
many other log intrefaces.

#### func (*Log) WithGlobalPackageAccess

```go
func (l *Log) WithGlobalPackageAccess()
```

#### func (*Log) Z

```go
func (l *Log) Z() *zerolog.Logger
```
Z rerturns a pointer to the underlying [zerolog.Logger].

---
