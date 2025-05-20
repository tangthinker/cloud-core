package helper

import (
	"github.com/disintegration/imaging"
	"image"
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
