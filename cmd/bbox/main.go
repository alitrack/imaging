package main

// 计算图像中非零区域的边界框。将边界框作为定义左、上、右和下像素坐标的四元组返回

import (
	"fmt"
	"os"

	ali "github.com/alitrack/imaging"
	imaging "github.com/disintegration/imaging"
)

func help() {
	fmt.Println("Usage: bbox <image path>")
	fmt.Println("Return: The bounding box is returned as a 4-tuple defining the left, upper, right, and lower pixel coordinate")
	fmt.Println("Image format support: jpg(or jpeg), png, gif, tif (or tiff) and bmp.")
	fmt.Println("Author: Steven Lee")
	fmt.Println("Website: alitrack.com")
}

func main() {

	if l := len(os.Args); l != 2 {
		help()
		return
	}

	// Open a test image.
	src, err := imaging.Open(os.Args[1])
	if err != nil {
		fmt.Printf("failed to open image: %v\n", err.Error())
		return
	}

	x1, y1, x2, y2 := ali.GetBBox(src)
	fmt.Printf("(%v,%v,%v,%v)\n", x1, y1, x2, y2)
}
