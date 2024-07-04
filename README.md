# fatal

fatal by [LitFill](https://github.com/LitFill).
Library for fatal assignment or logging (error management).
Using [log/slog package](https://pkg.go.dev/log/slog).

## Example

```go
Log(http.ServeAndListen(port),
    logger,
    "Can not serve and listen",
    "port", port,
)

file := Assign(os.Create(filename),
    logger,
    "Can not create file",
    "file name", filename,
)
```

## Usage notes

> [!NOTE]
> I don't know how to optimize this any further.

how this typically used:

```go
func main() {
  logFile := fatal.CreateLogFile("log.json")
  defer logFile.Close()
  logger := fatal.CreateLogger(io.MultiWriter(logFile, os.Stderr), slog.LevelInfo)
}
```

and then use it like normal.
