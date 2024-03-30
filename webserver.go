package webserver

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

type New struct {
	Dir  string
	Port string
	Log  bool
}

func (s *New) Start() {
	if !hasHTMLFilesInDir(s.Dir) {
		log.Println("No .html files found in the directory")
		return
	}

	go s.startHTTPServer()

	watcher := setupFileWatcher(s.Dir)
	defer watcher.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down...")
}

func (s *New) startHTTPServer() {
	http.HandleFunc("/", s.handler)

	log.Printf("Hello Onii-chan! Server running on localhost:%s\n\n", s.Port)
	log.Fatal(http.ListenAndServe(":"+s.Port, nil))
}

func (s *New) handler(w http.ResponseWriter, r *http.Request) {
	if s.Log {
		log.Printf("%s %s", r.Method, r.URL.Path)
	}

	path := r.URL.Path

	if path == "" || path == "/" || strings.HasSuffix(path, "/") {
		path = filepath.Join(s.Dir, path, "index.html")
		http.ServeFile(w, r, path)
		return
	}

	if !strings.Contains(path, ".") {
		path += ".html"
	}

	http.ServeFile(w, r, filepath.Join(s.Dir, path))
}

func setupFileWatcher(dir string) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
