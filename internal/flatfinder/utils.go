package flatfinder

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

// storeConfig - Write current config to disk
func (c *LocalConfig) storeConfig() {
	configFilePath := getConfigFilePath()

	json, err := json.Marshal(c)
	if err != nil {
		log.Fatal("Failed to JSONify config")
	}

	err = os.WriteFile(configFilePath, json, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// loadConfig - Pull existing config (if exists)
func (c *LocalConfig) loadConfig() {
	configFilePath := getConfigFilePath()
	if fileExists(configFilePath) {
		data, err := os.ReadFile(configFilePath)
		if err != nil {
			log.Fatal(err)
		}

		// Load it into global
		err = json.Unmarshal(data, c)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Loaded %d previously posted property IDs", len(c.PostedProperties))
	} else {
		// Create empty map for first run
		maps := make(map[int64]bool)
		c.PostedProperties = maps
	}
}

// getConfigFilePath - Returns a string of the config file pathg
func getConfigFilePath() string {
	// path := ""
	// switch runtime.GOOS {
	// case "linux":
	// 	if os.Getenv("XDG_CONFIG_HOME") != "" {
	// 		path = os.Getenv("XDG_CONFIG_HOME")
	// 	} else {
	// 		path = filepath.Join(os.Getenv("HOME"), ".config")
	// 	}
	// case "windows":
	// 	path = os.Getenv("APPDATA")
	// case "darwin":
	// 	path = os.Getenv("HOME") + "/Library/Application Support"
	// default:
	// 	log.Fatalf("Unsupported platform? %s", runtime.GOOS)
	// }

	// path = path + fmt.Sprintf("%c", os.PathSeparator) + "flatfinder"
	// err := os.MkdirAll(path, os.ModePerm)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// return path + fmt.Sprintf("%c", os.PathSeparator) + "flatfinder.json"
	return "flatfinder.json"
}

// fileExists - Check if a file exists
func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		log.Fatal(err)
	}

	return false
}
