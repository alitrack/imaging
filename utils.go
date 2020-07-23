package imaging

import (
	"image"
	"image/color"
	"image/draw"
)

var backgroundColor color.Color = color.White

//MergeGrids return merged image
func MergeGrids(images []image.Image, imageCountDX, imageCountDY int) (*image.RGBA, error) {
	var canvas *image.RGBA
	imageBoundX := 0
	imageBoundY := 0

	imageBoundX = images[0].Bounds().Dx()
	imageBoundY = images[0].Bounds().Dy()
	canvasBoundX := imageCountDX * imageBoundX
	canvasBoundY := imageCountDY * imageBoundY

	canvasMaxPoint := image.Point{canvasBoundX, canvasBoundY}
	canvasRect := image.Rectangle{image.Point{0, 0}, canvasMaxPoint}
	canvas = image.NewRGBA(canvasRect)

	// draw grids one by one
	for i, img := range images {
		x := i % imageCountDX
		y := i / imageCountDX
		minPoint := image.Point{x * imageBoundX, y * imageBoundY}
		maxPoint := minPoint.Add(image.Point{imageBoundX, imageBoundY})
		nextGridRect := image.Rectangle{minPoint, maxPoint}

		draw.Draw(canvas, nextGridRect, &image.Uniform{backgroundColor}, image.Point{}, draw.Src)
		draw.Draw(canvas, nextGridRect, img, image.Point{}, draw.Over)
	}
	return canvas, nil
}

//from https://github.com/oliamb/cutter/blob/master/cutter.go
// An interface that is
// image.Image + SubImage method.
type subImageSupported interface {
	SubImage(r image.Rectangle) image.Image
}

func cropWithCopy(img image.Image, cr image.Rectangle) (image.Image, error) {
	result := image.NewRGBA(cr)
	draw.Draw(result, cr, img, cr.Min, draw.Src)
	return result, nil
}

//Crop crop image with given rect
func Crop(src image.Image, rect image.Rectangle) (dst image.Image, err error) {
	if dst, ok := src.(subImageSupported); ok {
		// fmt.Println(ok)
		return dst.SubImage(rect), nil
	}
	return cropWithCopy(src, rect)
}
