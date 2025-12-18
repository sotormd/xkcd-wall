package xkcd

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func parseHexColor(s string) (color.RGBA, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid color: %s", s)
	}

	r, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	g, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	b, err := strconv.ParseUint(s[4:6], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}, nil
}

func lerp(a, b uint8, t float64) uint8 {
	return uint8(float64(a)*(1-t) + float64(b)*t)
}

func Colorize(
	inputPath string,
	backgroundHex string,
	foregroundHex string,
	cacheDir string,
) (string, error) {

	bg, err := parseHexColor(backgroundHex)
	if err != nil {
		return "", err
	}

	fg, err := parseHexColor(foregroundHex)
	if err != nil {
		return "", err
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	out := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			if a == 0 {
				out.Set(x, y, bg)
				continue
			}

			// Convert to 8-bit
			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)

			// Perceptual luminance
			lum := (0.299*r8 + 0.587*g8 + 0.114*b8) / 255.0

			// Optional gamma correction for smoother edges
			lum = math.Pow(lum, 1.8)

			out.Set(x, y, color.RGBA{
				R: lerp(fg.R, bg.R, lum),
				G: lerp(fg.G, bg.G, lum),
				B: lerp(fg.B, bg.B, lum),
				A: 255,
			})
		}
	}

	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", err
	}

	outPath := filepath.Join(cacheDir, "colored.png")
	outFile, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if err := png.Encode(outFile, out); err != nil {
		return "", err
	}

	return outPath, nil
}

func parseDimension(s string) (int, int, error) {
	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid dimension: %s", s)
	}

	w, err := strconv.Atoi(parts[0])
	if err != nil || w <= 0 {
		return 0, 0, fmt.Errorf("invalid width: %s", parts[0])
	}

	h, err := strconv.Atoi(parts[1])
	if err != nil || h <= 0 {
		return 0, 0, fmt.Errorf("invalid height: %s", parts[1])
	}

	return w, h, nil
}

func MakeBackground(
	dimension string,
	bgHex string,
	cacheDir string,
) (string, error) {

	width, height, err := parseDimension(dimension)
	if err != nil {
		return "", err
	}

	bg, err := parseHexColor(bgHex)
	if err != nil {
		return "", err
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fast fill
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bg)
		}
	}

	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", err
	}

	outPath := filepath.Join(cacheDir, "background.png")
	outFile, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		return "", err
	}

	return outPath, nil
}

func CompositeCenter(
	fgPath string,
	bgPath string,
	cacheDir string,
) (string, error) {

	// --- Load background ---
	bgFile, err := os.Open(bgPath)
	if err != nil {
		return "", err
	}
	defer bgFile.Close()

	bgImg, _, err := image.Decode(bgFile)
	if err != nil {
		return "", err
	}

	// --- Load foreground ---
	fgFile, err := os.Open(fgPath)
	if err != nil {
		return "", err
	}
	defer fgFile.Close()

	fgImg, _, err := image.Decode(fgFile)
	if err != nil {
		return "", err
	}

	bgBounds := bgImg.Bounds()
	fgBounds := fgImg.Bounds()

	// --- Create output ---
	out := image.NewRGBA(bgBounds)

	// Draw background
	draw.Draw(out, bgBounds, bgImg, bgBounds.Min, draw.Src)

	// Compute centered position
	offsetX := (bgBounds.Dx() - fgBounds.Dx()) / 2
	offsetY := (bgBounds.Dy() - fgBounds.Dy()) / 2

	if offsetX < 0 || offsetY < 0 {
		return "", fmt.Errorf("foreground larger than background")
	}

	fgPoint := image.Point{
		X: bgBounds.Min.X + offsetX,
		Y: bgBounds.Min.Y + offsetY,
	}

	fgRect := image.Rectangle{
		Min: fgPoint,
		Max: fgPoint.Add(fgBounds.Size()),
	}

	// Draw foreground with alpha
	draw.Draw(out, fgRect, fgImg, fgBounds.Min, draw.Over)

	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", err
	}

	outPath := filepath.Join(cacheDir, "final.png")
	outFile, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if err := png.Encode(outFile, out); err != nil {
		return "", err
	}

	return outPath, nil
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	// Make sure data is flushed to disk
	return out.Sync()
}
