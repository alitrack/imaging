# imaging

Imaging is a simple image processing package for Go

## GetBBox 

return the bounding box as a 4-tuple defining the left, upper, right, and lower pixel coordinate.

If the image is completely empty, this method returns (0,0,0,0).

other solution,
* GostScript
```
$ gs -sDevice=bbox tempCmykPdfFile.pdf|grep BoundingBox
>> return 
 %%BoundingBox: 13 48 199 100
```

* Python PIL
```
import Image
im=Image.open("myfile.png")
print im.getbbox()
```

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

