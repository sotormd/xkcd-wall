package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type XKCDResponse struct {
	Num int    `json:"num"`
	Img string `json:"img"`
}

func fetchJSON(url string, out any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func resolveComicNumber(kind string) (int, error) {
	switch kind {
	case "today":
		var latest XKCDResponse
		if err := fetchJSON("https://xkcd.com/info.0.json", &latest); err != nil {
			return 0, err
		}
		return latest.Num, nil

	case "random":
		var latest XKCDResponse
		if err := fetchJSON("https://xkcd.com/info.0.json", &latest); err != nil {
			return 0, err
		}

		rand.Seed(time.Now().UnixNano())
		return rand.Intn(latest.Num) + 1, nil

	default:
		n, err := strconv.Atoi(kind)
		if err != nil || n <= 0 {
			return 0, fmt.Errorf("invalid comic type: %s", kind)
		}
		return n, nil
	}
}

func downloadImage(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error: %s", resp.Status)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func Get(kind string, cacheDir string) (string, error) {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", err
	}

	num, err := resolveComicNumber(kind)
	if err != nil {
		return "", err
	}

	var comic XKCDResponse
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", num)

	if err := fetchJSON(url, &comic); err != nil {
		return "", err
	}

	if comic.Img == "" {
		return "", fmt.Errorf("no image found for comic %d", num)
	}

	outPath := filepath.Join(cacheDir, "comic.png")

	if err := downloadImage(comic.Img, outPath); err != nil {
		return "", err
	}

	return outPath, nil
}
