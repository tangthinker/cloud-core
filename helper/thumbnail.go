package helper

import (
	"bytes"
	"fmt"
	"image"
	"os/exec"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/heic"
)

// DecodeImage decodes an image from bytes, supporting HEIC via github.com/gen2brain/heic if necessary.
func DecodeImage(data []byte, filename string) (image.Image, error) {
	if strings.HasSuffix(strings.ToLower(filename), ".heic") {
		return DecodeHEIC(data)
	}
	return imaging.Decode(bytes.NewReader(data))
}

// DecodeHEIC decodes a HEIC image using github.com/gen2brain/heic.
func DecodeHEIC(data []byte) (image.Image, error) {
	img, err := heic.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("heic decode failed: %w", err)
	}
	return img, nil
}

// GenerateThumbnail generates a thumbnail from the source image with specified width and height,
func GenerateThumbnail(source image.Image, targetWidth, targetHeight int) (image.Image, error) {

	srcWidth, srcHeight := source.Bounds().Dx(), source.Bounds().Dy()

	aspectRatio := float64(srcWidth) / float64(srcHeight)
	targetAspectRatio := float64(targetWidth) / float64(targetHeight)

	var scaledWidth, scaledHeight int
	if aspectRatio > targetAspectRatio {
		scaledHeight = targetHeight
		scaledWidth = int(float64(scaledHeight) * aspectRatio)
	} else {
		scaledWidth = targetWidth
		scaledHeight = int(float64(scaledWidth) / aspectRatio)
	}

	resized := imaging.Resize(source, scaledWidth, scaledHeight, imaging.Lanczos)

	x0 := (resized.Bounds().Dx() - targetWidth) / 2
	y0 := (resized.Bounds().Dy() - targetHeight) / 2
	cropRect := image.Rect(x0, y0, x0+targetWidth, y0+targetHeight)

	thumbnail := imaging.Crop(resized, cropRect)

	return thumbnail, nil
}

// GenerateVideoThumbnail generates a thumbnail from a video file with specified width and height.
func GenerateVideoThumbnail(filepath string, targetWidth, targetHeight int) (image.Image, error) {
	// Construct FFmpeg command to extract a single frame at 1 second and output to stdout
	cmd := exec.Command(
		"ffmpeg",
		"-noaccurate_seek",
		"-i", filepath,
		"-ss", "1",
		"-vframes", "1",
		"-q:v", "2",
		"-an",
		"-f", "image2",
		"-vcodec", "mjpeg",
		"pipe:",
	)

	// Capture FFmpeg stdout and stderr
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Run FFmpeg command
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %w, stderr: %s", err, stderr.String())
	}

	// Decode the image from ffmpeg stdout
	img, _, err := image.Decode(&out)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Generate thumbnail using the existing function
	thumbnail, err := GenerateThumbnail(img, targetWidth, targetHeight)
	if err != nil {
		return nil, fmt.Errorf("failed to generate thumbnail: %w", err)
	}

	return thumbnail, nil
}
