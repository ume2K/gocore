package framework

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func LoadEnv(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables.")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Kommentare und leere Zeilen ignorieren
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split am ERSTEN Gleichheitszeichen (wichtig für Base64 Strings oder Connection Strings)
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Optional: Anführungszeichen entfernen, falls vorhanden (DX)
		value = strings.Trim(value, `"'`)

		os.Setenv(key, value)
	}
}

// GetEnv holt einen Wert oder einen Fallback
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
