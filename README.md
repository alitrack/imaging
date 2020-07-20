# imaging

Imaging is a simple image processing package for Go

## GetBBox 

return the bounding box as a 4-tuple defining the left, upper, right, and lower pixel coordinate.

If the image is completely empty, this method returns (0,0,0,0).

## Installation

```
$ go get github.com/alitrack/imaging
```

## Usage

```
var img image.Image
//your decode 
fmt.Println(imaging.GetBBox(img))
//if img==nil will return 0.0.0.0 

```



## License
MIT

## Author 
Steven Lee(steven#alitrack.com)

