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

var dir = flag.String("dir", ".", "serve specific folder.")
var port = flag.String("port", "9000", "specific port")

func main() {
	flag.Parse()

	if flag.NFlag() == 0 && (flag.Arg(0) == "-h" || flag.Arg(0) == "--help") {
		flag.Usage()
		os.Exit(0)
	}

	if !hasHTMLFilesInDir(*dir) {
		fmt.Println("No .html files found in the directory")
		return
	}

	go startHTTPServer()

	watcher := setupFileWatcher()
	defer watcher.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down...")
}

func startHTTPServer() {
	http.HandleFunc("/", handler)

	fmt.Printf("Hello onii-chan! Your server running on localhost:%s nyaa~\n\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	path := r.URL.Path

	if path == "" || path == "/" || strings.HasSuffix(path, "/") {
		path = filepath.Join(*dir, path, "index.html")
		http.ServeFile(w, r, path)
		return
	}

	if !strings.Contains(path, ".") {
		path += ".html"
	}

	http.ServeFile(w, r, *dir+path)
}

func setupFileWatcher() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
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

	return watcher
}

func hasHTMLFilesInDir(dir string) bool {
	files, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		log.Println(err)
		return false
	}
	return len(files) > 0
}
