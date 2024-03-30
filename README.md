<div align=center>

    # serve
    
    ###### serve -port 8080

</div>

## CLI

### Installation

```
go install github.com/nazhard/web-server/cmd/serve@latest
```

### Usage

```sh
$ serve -log=false
# serve current directory and no cute log

$ serve -h
# print help

$ serve -dir dist -port 8080
# serve files inside "dir" and use 8080 as port
```

## Use in your own projects

### Installation

```sh
go get github.com/nazhard/web-server@latest
```

### Usage

```go
package main

import (
    "github.com/nazhard/web-server"
)

func main() {
    server := &webserver.New{
        Dir:  "path/to/static/assets",
        Port: "8080",
        Log:  true,
    }

    server.Start()
}
```

## Why?

Just for test purposes.
It's faster than `pnpm preview` you know.

My use case:
- pnpm build
- serve -dir dist

Yeah, just to see if everything is okay.
