package framework

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type JSConfig struct {
	SourceDir string
	OutputDir string
}

type JSCompiler struct {
	config JSConfig
	mu     sync.Mutex
}

func NewJSCompiler(cfg JSConfig) *JSCompiler {
	return &JSCompiler{config: cfg}
}

func (c *JSCompiler) Bundle() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := os.Stat(c.config.OutputDir); os.IsNotExist(err) {
		os.MkdirAll(c.config.OutputDir, 0755)
	}

	inputFile := filepath.Join(c.config.SourceDir, "main.js")
	outputFile := filepath.Join(c.config.OutputDir, "main.js")
	cmd := exec.Command("esbuild", inputFile, "--bundle", "--minify", "--sourcemap", "--outfile="+outputFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("JS Bundle Error:\n%s", string(output))
		return err
	}
	log.Println("JS bundled successfully")
	return nil
}

func (c *JSCompiler) Watch() {
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

	log.Println("JS watcher active (waiting for changes..)")
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
			if ext == ".js" {
				go c.Bundle()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}
