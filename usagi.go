package usagi

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const usagiLogo = `
 ___  ___  ________  ________  ________  ___     
|\  \|\  \|\   ____\|\   __  \|\   ____\|\  \    
\ \  \\\  \ \  \___|\ \  \|\  \ \  \___|\ \  \   
 \ \  \\\  \ \_____  \ \   __  \ \  \  __\ \  \  
  \ \  \\\  \|____|\  \ \  \ \  \ \  \|\  \ \  \ 
   \ \_______\____\_\  \ \__\ \__\ \_______\ \__\
    \|_______|\_________\|__|\|__|\|_______|\|__|
             \|_________|                                                                      
  🐰 Usagi: The Cute Static File Serve 🐇
`

type Jump struct {
	Dir  string
	Port string
	Log  bool
}

func (u *Jump) Start() {
	fmt.Print(usagiLogo)

	fs := http.FileServer(http.Dir(filepath.Join(u.Dir, "assets")))

	http.Handle("/assets", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", u.handler)

	fmt.Printf("\n  Hello Onii-chan! Server running on http://localhost:%s\n\n", u.Port)
	log.Fatal(http.ListenAndServe(":"+u.Port, nil))

	u.startFileWatcher()

	select {}
}

func (u *Jump) handler(w http.ResponseWriter, r *http.Request) {
	if u.Log {
		log.Printf("%s %s", r.Method, r.URL.Path)
	}

	path := r.URL.Path
	if path == "/" || strings.HasSuffix(path, "/") {
		path = filepath.Join(u.Dir, path, "index.html")
	} else {
		path = filepath.Join(u.Dir, path)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		u.handleNotFound(w, r)
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(path))
	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}

	http.ServeFile(w, r, path)
}

func (u *Jump) handleNotFound(w http.ResponseWriter, r *http.Request) {
	errorPage := filepath.Join(u.Dir, "404.html")
	if _, err := os.Stat(errorPage); err == nil {
		http.ServeFile(w, r, errorPage)
	} else {
		u.listDirectory(w, r)
	}
}

func (u *Jump) listDirectory(w http.ResponseWriter, r *http.Request) {
	dirPath := filepath.Join(u.Dir, r.URL.Path)
	files, err := os.ReadDir(dirPath)
	if err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html><body><h1>Directory Listing</h1><ul>"))

	for _, file := range files {
		fileName := file.Name()
		if file.IsDir() {
			fileName += "/"
		}
		w.Write([]byte("<li><a href=\"" + fileName + "\">" + fileName + "</a></li>"))
	}

	w.Write([]byte("</ul></body></html>"))
}

func (u *Jump) startFileWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("failed to create watcher: %v", err)
	}

	err = filepath.Walk(u.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		return watcher.Add(path)
	})

	if err != nil {
		log.Fatalf("failed to watch directory: %v", err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					log.Printf("File modified: %s", event.Name)
				}
			case err := <-watcher.Errors:
				log.Printf("watcher error: %v", err)
			}
		}
	}()
}
