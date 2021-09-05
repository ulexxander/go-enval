# enval - simple package with small footprint for reading env variables

## Motivation

I used https://github.com/spf13/viper library in some backend projects that required configuration.
But I was always had feeling that this package is huge with very large API and some dependencies I never needed.

In fact it was missing some basic features I wanted -
return error if required variable is missing or value is invalid (not just return zero value and try to use it).

Common problem I had in most projects:

- New developer clones repository and runs executable
- Program connects to postgres with host, username, password empty strings and port=0 (because required variables were not set)
- So why not just validate all variables on program start and print all configuration errors :)

This little package does one thing and does it well (reading variables).
But it may not cover all your project configuration requirements, for example - load variables from `.env` file.

For that I can recommend using https://github.com/joho/godotenv or any similar packages.

## Features

Lookuper type

- Getting `string, int, bool` variables values calling `os.LookupEnv` under the hood
- Support any custom type using `Custom*` methods and implementing ParseFunc - `func(val string) (interface{}, error)`
- Provide default values when retrieving variables
- Support errors like variable is missing or value is in invalid format
- Defer error check - no need to check error on every variable read
- Use `Err() error` method to check all read varialbes at once
- Use `ErrByVariable` map to check if specific variable has error
- Override `LookupFunc` to lookup values from anywhere (see examples and unit tests)

## Examples

Reading variables

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

Custom type

```go
l := enval.NewLookuper()
// no need to override LookupFunc in your code
// NewLookuper sets it to os.LookupEnv
l.LookupFunc = exampleVariablesLookupFunc

type logLevel string

const (
  INFO  logLevel = "INFO"
  WARN  logLevel = "WARN"
  ERROR logLevel = "ERROR"
)

pf := func(val string) (interface{}, error) {
  switch logLevel(val) {
  case INFO:
    return INFO, nil
  case WARN:
    return WARN, nil
  case ERROR:
    return ERROR, nil
  default:
		return logLevel(""), fmt.Errorf("unknown log level: %s", val)
  }
}

level1 := l.Custom("LOG_LEVEL", pf).(logLevel)
level2 := l.CustomWithDefault("LOG_LEVEL_WITH_DEFAULT", INFO, pf).(logLevel)
_ = l.CustomWithDefault("LOG_LEVEL_INVALID_ONE", INFO, pf).(logLevel)
fmt.Println(level1, level2)

if err := l.Err(); err != nil {
  fmt.Println(err)
}
// Output:
// WARN INFO
// LOG_LEVEL_INVALID_ONE: unknown log level: DEBUGGIN
```
