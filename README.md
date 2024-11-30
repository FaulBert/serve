<span align=center>

# Usagi Serve
    
###### usagi -port 8080

</span>

---

## CLI

### Installation

```
go install github.com/nazhard/usagi-serve/cmd/usagi@latest
```

### Usage

```sh
$ usagi -log=false
# serve current directory and no cute log

$ usagi -h
# print help

$ usagi -dir dist -port 8080
# serve files inside "dir" and use 8080 as port
```

## Use in your own projects

### Installation

```sh
go get github.com/nazhard/usagi-serve@latest
```

### Usage

```go
package main

import (
    "github.com/nazhard/usagi-serve"
)

func main() {
    server := &usagi.Jump{
        Dir:  "path/to/static/assets",
        Port: "8080",
        Log:  true,
    }

    server.Start()
}
```

## Why?

I don't know. Why not?

old reason :
> Just for test purposes.
> It's faster than `pnpm preview` you know.
> 
> My use case:
> - pnpm build
> - serve -dir dist
>
> Yeah, just to see if everything is okay.
