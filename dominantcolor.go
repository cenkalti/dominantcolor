// Copyright (c) 2011 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dominantcolor

import (
	"image"
	"image/color"
	"math/rand"
	"sort"

	"github.com/nfnt/resize"
)

const (
	maxSize       = 64
	nCluster      = 4
	maxSample     = 10
	nIterations   = 50
	maxBrightness = 665
	minDarkness   = 100
)

func Find(img image.Image) color.RGBA {
	img = resize.Thumbnail(maxSize, maxSize, img, resize.NearestNeighbor)
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	rnd := rand.New(rand.NewSource(42))
	randomPoint := func() (x, y int) {
		x = bounds.Min.X + int(rnd.Int31n(int32(width)))
		y = bounds.Min.Y + int(rnd.Int31n(int32(height)))
		return
	}
	// Pick a starting point for each cluster
	clusters := make(kMeanClusters, 0, nCluster)
	for i := 0; i < nCluster; i++ {
		// Try up to 10 times to find a unique color. If no unique color can be
		// found, destroy this cluster.
		colorUnique := false
		for j := 0; j < maxSample; j++ {
			ri, gi, bi, _ := img.At(randomPoint()).RGBA()
			r, g, b := uint8(ri), uint8(gi), uint8(bi)
			// Check to see if we have seen this color before.
			colorUnique := !clusters.ContainsCentroid(r, g, b)
			// If we have a unique color set the center of the cluster to
			// that color.
			if colorUnique {
				c := new(kMeanCluster)
				c.SetCentroid(r, g, b)
				clusters = append(clusters, c)
				break
			}
		}
		if !colorUnique {
			break
		}
	}
	convergence := false
	for i := 0; i < nIterations && !convergence && len(clusters) != 0; i++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				ri, gi, bi, _ := img.At(x, y).RGBA()
				r, g, b := uint8(ri), uint8(gi), uint8(bi)
				// Figure out which cluster this color is closest to in RGB space.
				closest := clusters.Closest(r, g, b)
				closest.AddPoint(r, g, b)
			}
		}
		// Calculate the new cluster centers and see if we've converged or not.
		convergence = true
		for _, c := range clusters {
			convergence = convergence && c.CompareCentroidWithAggregate()
			c.RecomputeCentroid()
		}
	}
	// Sort the clusters by population so we can tell what the most popular
	// color is.
	sort.Sort(byWeight(clusters))
	// Loop through the clusters to figure out which cluster has an appropriate
	// color. Skip any that are too bright/dark and go in order of weight.
	var col color.RGBA
	for i, c := range clusters {
		var r, g, b uint8
		c.GetCentroid(&r, &g, &b)
		// Sum the RGB components to determine if the color is too bright or too dark.
		var summedColor uint32 = uint32(r) + uint32(g) + uint32(b)

		if summedColor < maxBrightness && summedColor > minDarkness {
			// If we found a valid color just set it and break. We don't want to
			// check the other ones.
			col.R = r
			col.G = g
			col.B = b
			col.A = 0xFF
			break
		} else if i == 0 {
			// We haven't found a valid color, but we are at the first color so
			// set the color anyway to make sure we at least have a value here.
			col.R = r
			col.G = g
			col.B = b
			col.A = 0xFF
		}
	}
	return col
}
