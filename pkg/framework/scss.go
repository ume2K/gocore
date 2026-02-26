package framework

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type SCSSConfig struct {
	SourceDir string
	OutputDir string
	Debug     bool
}

type SCSSCompiler struct {
	config SCSSConfig
	mu     sync.Mutex
}

func NewSCSSCompiler(cfg SCSSConfig) *SCSSCompiler {
	return &SCSSCompiler{config: cfg}
}

func (c *SCSSCompiler) Compile() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := os.Stat(c.config.OutputDir); os.IsNotExist(err) {
		os.MkdirAll(c.config.OutputDir, 0755)
	}

	inputFile := filepath.Join(c.config.SourceDir, "main.scss")
	outputFile := filepath.Join(c.config.OutputDir, "main.css")

	cmd := exec.Command("sass", inputFile, outputFile, "--no-source-map", "--style=compressed")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("SCSS Error:\n%s", string(output))
		return err
	}
	log.Println("SCSS compiled successfully")
	return nil
}

func (c *SCSSCompiler) Watch() {
	absPath, err := filepath.Abs(c.config.SourceDir)
	if err != nil {
		log.Printf("Could not resolve path: %v", err)
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error scanning directories: %v", err)
	}

	log.Println("SCSS watcher active (waiting for changes..)")
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				continue
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					watcher.Add(event.Name)
				}
			}
			ext := filepath.Ext(event.Name)
			if ext == ".scss" || ext == ".css" {
				go c.Compile()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}
