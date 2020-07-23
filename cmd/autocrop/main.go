package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	ali "github.com/alitrack/imaging"
	"github.com/disintegration/imaging"
)

func help() {
	fmt.Println("Usage: autocrop <input image path> <output image path>")
	// fmt.Println("Return: The bounding box is returned as a 4-tuple defining the left, upper, right, and lower pixel coordinate")
	fmt.Println("Image format support: jpg(or jpeg), png, gif, tif (or tiff) and bmp.")
	fmt.Println("Author: Steven Lee")
	fmt.Println("Website: alitrack.com")
}

func main() {

	if l := len(os.Args); l != 3 {
		help()
		return
	}

	// Open a test image.
	//imaging.AutoOrientation(true)
	src, err := imaging.Open(os.Args[1])
	if err != nil {
		fmt.Printf("failed to open image: %v\n", err.Error())
		return
	}

	x1, y1, x2, y2 := ali.GetBBox(src)

	fmt.Printf("bbox:(%v,%v,%v,%v)\n", x1, y1, x2, y2)

	rect := image.Rect(x1, y1, x2, y2)

	dst, err := crop(src, rect)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = imaging.Save(dst, os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

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

func crop(src image.Image, rect image.Rectangle) (dst image.Image, err error) {
	if dst, ok := src.(subImageSupported); ok {
		fmt.Println(ok)
		return dst.SubImage(rect), nil
	}

	return cropWithCopy(src, rect)

}
