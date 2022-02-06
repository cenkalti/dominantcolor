package dominantcolor_test

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
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

func testImage(t *testing.T) image.Image {
	f, err := os.Open("firefox.png")
	if err != nil {
		t.Fatal(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return img
}

func TestFind(t *testing.T) {
	img := testImage(t)
	c := dominantcolor.Find(img)
	d := distance(c, firefoxOrange)
	t.Log("Found dominant color:", dominantcolor.Hex(c))
	t.Log("Firefox orange:      ", dominantcolor.Hex(firefoxOrange))
	t.Logf("Distance:             %.2f", d)
	if d > 50 {
		t.Errorf("Found color is not close.")
	}
}

func TestFindN(t *testing.T) {
	img := testImage(t)
	colors := dominantcolor.FindN(img, 4)

	if len(colors) != 4 {
		t.Error("Did not find 4 colors. Got:", len(colors))
	}

	for i, c := range colors {
		t.Logf("%d/%d Found dominant color: %s", i, len(colors), dominantcolor.Hex(c))

		paletted := image.NewPaletted(image.Rect(0, 0, 64, 64), []color.Color{c})
		f, err := os.OpenFile(fmt.Sprintf("_test_palette_%d.png", i), os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		_ = png.Encode(f, paletted)
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
