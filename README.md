Dominantcolor
============

[![GoDoc](https://godoc.org/github.com/cenkalti/dominantcolor?status.svg)](http://godoc.org/github.com/cenkalti/dominantcolor)

Find dominant color in images

``` go
import "github.com/cenkalti/dominantcolor"
```

Package dominantcolor provides a function for finding a color that represents the calculated dominant color in the image. This uses a KMean clustering algorithm to find clusters of pixel colors in RGB space.

The algorithm is ported from Chromium source code:

- https://github.com/adobe/chromium/blob/master/ui/gfx/color_analysis.h
- https://github.com/adobe/chromium/blob/master/ui/gfx/color_analysis.cc

See more at: http://godoc.org/github.com/cenkalti/dominantcolor

#### Example

``` go
package main

import (
	"fmt"
	"github.com/cenkalti/dominantcolor"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func FindDomiantColor(fileInput string) (string, error) {
	f, err := os.Open(fileInput)
	defer f.Close()
	if err != nil {
		fmt.Println("File not found:", fileInput)
		return "", err
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}

	return dominantcolor.Hex(dominantcolor.Find(img)), nil
}

func main() {
	fmt.Println(FindDomiantColor("aa.png"))
}

```

#### Output:
```
#CA5527
```
