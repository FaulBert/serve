package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

var dir = flag.String("dir", ".", "serve spesific folder.")

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	urlPath := r.URL.Path

	if urlPath == "" || urlPath == "/" {
		http.ServeFile(w, r, *dir+"/index.html")
		return
	}

	if !strings.Contains(urlPath, ".") {
		urlPath += ".html"
	}

	http.ServeFile(w, r, *dir+urlPath)
}

func main() {
	port := flag.String("port", "9000", "spesific port")

	flag.Parse()

	if flag.NFlag() == 0 && (flag.Arg(0) == "-h" || flag.Arg(0) == "--help") {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	fmt.Printf("Hello onii-chan! Your connection at localhost:%s nyaa~ \n", *port)
	go func() {
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(":"+*port, nil))
	}()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					log.Println("File modified:", event.Name)
				}
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down...")
}
