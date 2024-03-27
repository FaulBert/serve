package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var serve = flag.String("dir", ".", "serve spesific folder.")

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	urlPath := r.URL.Path

	if urlPath == "" || urlPath == "/" {
		http.ServeFile(w, r, *serve+"/index.html")
		return
	}

	if !strings.Contains(urlPath, ".") {
		urlPath += ".html"
	}

	http.ServeFile(w, r, *serve+urlPath)
}

func main() {
	port := flag.String("port", "9000", "spesific port")

	flag.Parse()

	if flag.NFlag() == 0 && (flag.Arg(0) == "-h" || flag.Arg(0) == "--help") {
		printUsage()
		os.Exit(0)
	}

	http.HandleFunc("/", handler)

	fmt.Printf("Hello onii-chan! Your connection at localhost:%s nyaa~ \n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}
