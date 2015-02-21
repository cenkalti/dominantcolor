package dominantcolor

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"
	"testing"
)

// https://www.mozilla.org/en-US/styleguide/identity/firefox/color/
var firefoxOrange = color.RGBA{R: 230, G: 96}

func TestFind(t *testing.T) {
	f, err := os.Open("firefox.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	c := Find(img)
	d := distance(c, firefoxOrange)
	if d > 55 {
		t.Errorf("Found color is not close. Distance: %f", d)
		t.Log("Found dominant color:", Hex(c))
		t.Log("Firefox orange:      ", Hex(firefoxOrange))
	}
}

func distance(a, b color.RGBA) float64 {
	dr := uint32(a.R) - uint32(b.R)
	dg := uint32(a.G) - uint32(b.G)
	db := uint32(a.B) - uint32(b.B)
	return math.Sqrt(float64(dr*dr + dg*dg + db*db))
}
