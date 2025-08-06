# envmap

```shell
go get github.com/artarts36/envmap
```

This library converts configuration structure to map[string]string.

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
