package dominantcolor_test

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"
	"testing"

	"github.com/cenkalti/dominantcolor"
)

// https://www.mozilla.org/en-US/styleguide/identity/firefox/color/
var firefoxOrange = color.RGBA{R: 230, G: 96}

func Example() {
	f, _ := os.Open("firefox.png")
	img, _, _ := image.Decode(f)
	f.Close()
	fmt.Println(dominantcolor.Hex(dominantcolor.Find(img)))
	// Output: #CB5A27
}

func TestFind(t *testing.T) {
	f, err := os.Open("firefox.png")
	if err != nil {
		t.Fatal(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	c := dominantcolor.Find(img)
	d := distance(c, firefoxOrange)
	t.Log("Found dominant color:", dominantcolor.Hex(c))
	t.Log("Firefox orange:      ", dominantcolor.Hex(firefoxOrange))
	t.Logf("Distance:             %.2f", d)
	if d > 50 {
		t.Errorf("Found color is not close.")
	}
}

func distance(a, b color.RGBA) float64 {
	dr := uint32(a.R) - uint32(b.R)
	dg := uint32(a.G) - uint32(b.G)
	db := uint32(a.B) - uint32(b.B)
	return math.Sqrt(float64(dr*dr + dg*dg + db*db))
}

func BenchmarkFind(b *testing.B) {
	f, err := os.Open("firefox.png")
	if err != nil {
		b.Fatal(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		b.Fatal(err)
	}
	f.Close()
	for i := 0; i < b.N; i++ {
		dominantcolor.Find(img)
	}
}
