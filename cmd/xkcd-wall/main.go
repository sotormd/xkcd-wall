package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"xkcd-wall/xkcd"
)

//go:embed assets/default.json
var embeddedDefaultConfig []byte

type Config struct {
	BackgroundColors []string `json:"background-colors"`
	ForegroundColors []string `json:"foreground-colors"`
	Dimensions       string   `json:"dimensions"`
	Target           string   `json:"target"`
	Cache            string   `json:"cache"`
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not determine home directory: %v\n", err)
		os.Exit(1)
	}

	defaultConfigPath := filepath.Join(home, ".config", "xkcd-wall", "config.json")

	configPath := flag.String("c", defaultConfigPath, "Path to config.json")

	defaultComicType := "today"

	comicType := flag.String("t", defaultComicType, "today, random, or <number>")

	flag.Parse()

	configFile := *configPath

	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {

			configDir := filepath.Dir(configFile)

			if err := os.MkdirAll(configDir, 0o755); err != nil {
				fmt.Fprintf(os.Stderr, "Error: could not create config directory: %v\n", err)
				os.Exit(1)
			}

			if err := os.WriteFile(configFile, embeddedDefaultConfig, 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "Error: could not write default config: %v\n", err)
				os.Exit(1)
			}

		} else {
			fmt.Fprintf(os.Stderr, "Error: could not read config file: %v\n", err)
			os.Exit(1)
		}
	}

	var config Config

	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not read config file: %v\n", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not read configuration json: %v\n", err)
		os.Exit(1)
	}

	path, err := xkcd.Get(*comicType, config.Cache)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not fetch comic: %v\n", err)
		os.Exit(1)
	}

	bg := config.BackgroundColors[rand.Intn(len(config.BackgroundColors))]
	fg := config.ForegroundColors[rand.Intn(len(config.ForegroundColors))]

	colored, err := xkcd.Colorize(path, bg, fg, config.Cache)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not colorize comic: %v\n", err)
		os.Exit(1)
	}

	background, err := xkcd.MakeBackground(config.Dimensions, bg, config.Cache)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not create background: %v\n", err)
		os.Exit(1)
	}

	final, err := xkcd.CompositeCenter(colored, background, config.Cache)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not create final image: %v\n", err)
		os.Exit(1)
	}

	if err := xkcd.CopyFile(final, config.Target); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not copy image to target: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(config.Target)
}
