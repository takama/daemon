Go Daemon
=========

A daemon package for use with Go (golang) services

[![GoDoc](https://godoc.org/github.com/takama/daemon?status.svg)](https://godoc.org/github.com/takama/daemon)

### Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/takama/daemon"
)

func main() {
	service, err := daemon.New("name", "description")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	status, err := service.Install()
	if err != nil {
		log.Fatal(status, "\nError: ", err)
	}
	fmt.Println(status)
}
```

## Author

[Igor Dolzhikov](https://github.com/takama)


## License

[MIT Public License](https://github.com/takama/daemon/blob/master/LICENSE)
