package helper

import (
	"bytes"
	"fmt"
	"image"
	"os/exec"

	"github.com/disintegration/imaging"
)

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
		"-ss", "1", // Seek to 1 second (放在 -i 前面更快)
		"-i", filepath, // Input video file
		"-vframes", "1", // Extract 1 frame
		"-q:v", "2", // Set quality (2 is high quality for JPEG)
		"-an",          // Disable audio
		"-f", "image2", // Output format
		"-vcodec", "mjpeg", // Output as JPEG
		"pipe:1", // Output to stdout
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
