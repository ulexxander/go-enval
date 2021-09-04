# enval - simple package with small footprint for reading env variables

## Features

Lookuper type

- Getting `string, int, bool` variables values calling `os.LookupEnv` under the hood
- Support any custom type using `Custom*` methods and implementing ParseFunc - `func(val string) (interface{}, error)`
- Provide default values when retrieving variables
- Support errors like variable is missing or value is in invalid format
- Defer error check - no need to check error on every variable read
- Use `Err() error` method to check all read varialbes at once
- Use `ErrByVariable` map to check if specific variable has error
- Override `LookupFunc` to lookup values from anywhere (see unit tests)

## Example

```go
l := enval.NewLookuper()
// no need to override LookupFunc in your code
// NewLookuper sets it to os.LookupEnv
l.LookupFunc = exampleVariablesLookupFunc

type config struct {
  PostgresHost          string
  PostgresPort          int
  PostgresUsername      string
  PostgresPassword      string
  MaintenanceMode       bool
  MaxConcurrentRequests int
}

_ = config{
  PostgresHost:          l.String("POSTGRES_HOST"),
  PostgresPort:          l.IntWithDefault("POSTGRES_PORT", 5432),
  PostgresUsername:      l.String("POSTGRES_USERNAME"),
  PostgresPassword:      l.String("POSTGRES_PASSWORD"),
  MaintenanceMode:       l.Bool("MAINTENANCE_MODE"),
  MaxConcurrentRequests: l.IntWithDefault("MAX_CONCURRENT_REQUESTS", 64),
}

if err := l.Err(); err != nil {
  fmt.Println(err)
}
// Output: POSTGRES_PASSWORD: variable missing, MAX_CONCURRENT_REQUESTS: unparsable int: strconv.ParseInt: parsing "q123": invalid syntax
```
