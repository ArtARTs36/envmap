# envmap

```shell
go get github.com/artarts36/envmap
```

This library converts environment configuration structure to `map[string]string`.

**envmap** does the reverse of [caarlos0/env](https://github.com/caarlos0/env), which converts environment variables into a Go structure.

Usage example:

```go
package main

import (
	"fmt"
	"time"

	"github.com/artarts36/envmap"
)

type Config struct {
	Mode string        `env:"MODE"`
	Timeout time.Duration `env:"TIMEOUT"`
}

func main() {
	result, _ := envmap.Convert(Config{
		Mode: "prod",
		Timeout: 30 * time.Second,
	}, envmap.WithPrefix("APP_"))

	fmt.Println(result)
	// map[string]string
	// APP_MODE: prod
	// APP_TIMEOUT: 30s
}
```
