package utils

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

const imageMagick = "convert"
const tempFileName = "imgConvertTmp"

var mu sync.Mutex

func IsImageMagickInstalled() bool {
	_, err := exec.LookPath(imageMagick)
	return err == nil
}

// ConvertJP2ToPNG uses ImageMagick to convert JP2 images to PNG
// It will use ImageMagick installed on the system to avoid the cgo requirement for this project
func ConvertJP2ToPNG(input []byte) ([]byte, error) {
	mu.Lock()
	defer mu.Unlock()

	jp2 := fmt.Sprintf("%s.jp2", tempFileName)
	png := fmt.Sprintf("%s.png", tempFileName)

	err := os.WriteFile(jp2, input, 0644)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(imageMagick, jp2, png)
	err = cmd.Run()
	if err != nil {
		// Revert
		os.Remove(jp2)
		return nil, err
	}

	// Successfully converted jp2 to png
	os.Remove(jp2)
	data, err := os.ReadFile(png)
	if err != nil {
		return nil, err
	}
	os.Remove(png)
	return data, nil
}
