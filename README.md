# fatal

fatal by [LitFill](https://github.com/LitFill).
library for fatal assignment or logging (error management)

## Example

```go
Log(http.ServeAndListen(port),
    "Can not serve and listen",
    "port", port
)

file := Assign(os.Create("log.txt"),
    "Can not create file",
    "file", "log.txt"
)
```
