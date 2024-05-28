# Juke It Out

[![Lint](https://github.com/AlexMeuer/juke/actions/workflows/lint.yml/badge.svg)](https://github.com/AlexMeuer/juke/actions/workflows/lint.yml)
[![Test](https://github.com/AlexMeuer/juke/actions/workflows/test.yml/badge.svg)](https://github.com/AlexMeuer/juke/actions/workflows/test.yml)

## Library assumptions

### Logging

We use [zerolog](https://github.com/rs/zerolog) instead of the standard `log` package.

```go
import (
	"github.com/rs/zerolog/log"
)
```

### Json

We use [jsoniter](https://github.com/json-iterator/go) instead of the standard `json` packages.

```go
import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary
json.Unmarshal(input, &data)
```
