package dominantcolor

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"testing"
)

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
	log.Print("#" + fmt.Sprintf("%.2x", c.R) + fmt.Sprintf("%.2x", c.G) + fmt.Sprintf("%.2x", c.B))
}
