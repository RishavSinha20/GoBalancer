package config

import (
	"log"
	"os"
	"time"
)

func WatchConfig(path string, reload func(*Config)) {
	var lastModified time.Time

	for {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if info.ModTime().After(lastModified) {
			cfg, err := LoadConfig(path)

			if err == nil {
				log.Println("Config reloaded")
				reload(cfg)
				lastModified = info.ModTime()
			}
		}
		time.Sleep(5 * time.Second)
	}
}
