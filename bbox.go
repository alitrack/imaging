package imaging

import (
	"image"
	"sync"
)

//GetBBox return the bounding box as a 4-tuple defining the left, upper, right, and lower pixel coordinate.
//If the image is completely empty, this method returns (0,0,0,0).
func GetBBox(src image.Image) (x1, y1, x2, y2 int) {
	if src == nil {
		return 0, 0, 0, 0
	}

	c := src.At(0, 0)
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	var wg sync.WaitGroup
	wg.Add(4)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for y1 = 0; y1 < height; y1++ {
			for j := 0; j < width; j++ {
				if src.At(j, y1) != c {
					return
				}
			}
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for x1 = 0; x1 < width; x1++ {
			for j := 0; j < height; j++ {
				if src.At(x1, j) != c {
					return
				}
			}
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for x2 = width - 1; x2 >= 0; x2-- {
			for j := height - 1; j >= 0; j-- {
				if src.At(x2, j) != c {
					x2++
					return
				}
			}
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for y2 = height - 1; y2 >= 0; y2-- {
			for j := width - 1; j >= 0; j-- {
				if src.At(j, y2) != c {
					y2++
					return
				}
			}
		}
	}(&wg)

	wg.Wait()
	return
}
