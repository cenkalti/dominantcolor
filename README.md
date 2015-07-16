Dominantcolor
============

[![GoDoc](https://godoc.org/github.com/cenkalti/dominantcolor?status.svg)](http://godoc.org/github.com/cenkalti/dominantcolor)

Find dominant color in images

```
import "github.com/cenkalti/dominantcolor"
```

Package dominantcolor provides a function for finding a color that represents the calculated dominant color in the image. This uses a KMean clustering algorithm to find clusters of pixel colors in RGB space.

The algorithm is ported from Chromium source code:

https://src.chromium.org/svn/trunk/src/ui/gfx/color_analysis.h
https://src.chromium.org/svn/trunk/src/ui/gfx/color_analysis.cc

See more at: http://godoc.org/github.com/cenkalti/dominantcolor

####Example

```
f, _ := os.Open("firefox.png")
img, _, _ := image.Decode(f)
f.Close()
fmt.Println(dominantcolor.Hex(dominantcolor.Find(img)))
```

####Output:
```
#CA5527
```
