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

// Compile nutzt esbuild zum Bundeln und Minifizieren
func (c *JSCompiler) Bundle() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := os.Stat(c.config.OutputDir); os.IsNotExist(err) {
		os.MkdirAll(c.config.OutputDir, 0755)
	}

	inputFile := filepath.Join(c.config.SourceDir, "main.js")
	outputFile := filepath.Join(c.config.OutputDir, "main.js")

	// Der Befehl: esbuild assets/js/main.js --bundle --minify --outfile=public/js/app.js
	cmd := exec.Command("esbuild", inputFile, "--bundle", "--minify", "--sourcemap", "--outfile="+outputFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("❌ JS Bundle Error:\n%s", string(output))
		return err
	}

	log.Println("✅ JS bundled successfully")
	return nil
}

func (c *JSCompiler) Watch() {
	// 1. Absoluten Pfad ermitteln (hilft gegen Pfad-Probleme unter Windows)
	absPath, err := filepath.Abs(c.config.SourceDir)
	if err != nil {
		log.Printf("❌ Konnte absoluten Pfad nicht ermitteln: %v", err)
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// 2. Rekursiv alle Unterordner zum Watcher hinzufügen
	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Nur Verzeichnisse watchen
		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				log.Printf("❌ Fehler beim Watchen von %s: %v", path, err)
			} else {
				// Debug-Output, damit du siehst, dass es klappt
				log.Printf("👀 Watching Directory: %s", path)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("❌ Fehler beim Scannen der Ordner: %v", err)
	}

	log.Println("✅ JS Watcher läuft und wartet auf Änderungen...")

	// 3. Die Endlos-Schleife (Blockierend)
	// Da du in main.go "go scss.Watch()" aufrufst, ist blockieren hier okay.
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Ignoriere CHMOD (passiert oft bei Editoren beim Speichern)
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				continue
			}

			// Wenn ein NEUER Ordner erstellt wird, müssen wir ihn auch watchen
			if event.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					log.Printf("📂 Neuer Ordner erkannt: %s", event.Name)
					watcher.Add(event.Name)
				}
			}

			// Nur bei .scss oder .css Dateien kompilieren
			ext := filepath.Ext(event.Name)
			if ext == ".js" {
				log.Printf("⚡ Änderung erkannt in: %s", filepath.Base(event.Name))

				// Optional: Kleines Debouncing (damit er nicht 2x feuert),
				// aber meistens reicht es so.
				go c.Bundle() // Compile in goroutine, damit der Watcher nicht blockiert
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("❌ Watcher error:", err)
		}
	}
}
