package heatmap

import (
	"image"
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"os"
)

// saveBytesToFile writes bytes to a file
func saveBytesToFile(filepath string, data []byte) error {
	return os.WriteFile(filepath, data, 0644)
}

// LoadImage loads an image from a file.
// Supports common formats: PNG, JPEG, GIF, etc.
func LoadImage(filepath string) (image.Image, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}
