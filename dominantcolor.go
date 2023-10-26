// Copyright (c) 2011 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package dominantcolor provides a function for finding
// a color that represents the calculated dominant color in the
// image. This uses a KMean clustering algorithm to find clusters of pixel
// colors in RGB space.
//
// The algorithm is ported from Chromium source code:
//
//	https://src.chromium.org/svn/trunk/src/ui/gfx/color_analysis.h
//	https://src.chromium.org/svn/trunk/src/ui/gfx/color_analysis.cc
//
// RGB KMean Algorithm (N clusters, M iterations):
//
// 1. Pick N starting colors by randomly sampling the pixels. If you see a
// color you already saw keep sampling. After a certain number of tries
// just remove the cluster and continue with N = N-1 clusters (for an image
// with just one color this should devolve to N=1). These colors are the
// centers of your N clusters.
//
// 2. For each pixel in the image find the cluster that it is closest to in RGB
// space. Add that pixel's color to that cluster (we keep a sum and a count
// of all of the pixels added to the space, so just add it to the sum and
// increment count).
//
// 3. Calculate the new cluster centroids by getting the average color of all of
// the pixels in each cluster (dividing the sum by the count).
//
// 4. See if the new centroids are the same as the old centroids.
//
// a) If this is the case for all N clusters than we have converged and can move on.
//
// b) If any centroid moved, repeat step 2 with the new centroids for up to M iterations.
//
// 5. Once the clusters have converged or M iterations have been tried, sort
// the clusters by weight (where weight is the number of pixels that make up
// this cluster).
//
// 6. Going through the sorted list of clusters, pick the first cluster with the
// largest weight that's centroid falls between |lower_bound| and
// |upper_bound|. Return that color.
// If no color fulfills that requirement return the color with the largest
// weight regardless of whether or not it fulfills the equation above.
package dominantcolor

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"sort"

	"github.com/disintegration/imaging"
)

const (
	resizeTo      = 256
	maxSample     = 10
	nIterations   = 50
	maxBrightness = 665
	minDarkness   = 100

	nClustersDefault = 4
)

type Color struct {
	color.RGBA
	Weight float64
}

func findClusters(img image.Image, nCluster int) (kMeanClusterGroup, float64) {
	// Shrink image for faster processing.
	img = imaging.Fit(img, resizeTo, resizeTo, imaging.NearestNeighbor)

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	rnd := rand.New(rand.NewSource(0))
	randomPoint := func() (x, y int) {
		x = bounds.Min.X + rnd.Intn(width)
		y = bounds.Min.Y + rnd.Intn(height)
		return
	}
	// Pick a starting point for each cluster.
	clusters := make(kMeanClusterGroup, 0, nCluster)
	for i := 0; i < nCluster; i++ {
		// Try up to 10 times to find a unique color. If no unique color can be
		// found, destroy this cluster.
		colorUnique := false
		for j := 0; j < maxSample; j++ {
			ri, gi, bi, a := img.At(randomPoint()).RGBA()
			// Ignore transparent pixels.
			if a == 0 {
				continue
			}
			r, g, b := uint8(ri/0x101), uint8(gi/0x101), uint8(bi/0x101)
			// Check to see if we have seen this color before.
			colorUnique = !clusters.ContainsCentroid(r, g, b)
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
				ri, gi, bi, a := img.At(x, y).RGBA()
				// Ignore transparent pixels.
				if a == 0 {
					continue
				}
				r, g, b := uint8(ri/0x101), uint8(gi/0x101), uint8(bi/0x101)
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
	return clusters, float64(width) * float64(height)
}

// Find returns the dominant color in img.
func Find(img image.Image) color.RGBA {
	colors := FindN(img, nClustersDefault)
	if len(colors) == 0 {
		return color.RGBA{0, 0, 0, 0}
	}
	// Loop through the clusters to figure out which cluster has an appropriate
	// color. Skip any that are too bright/dark and go in order of weight.
	for _, c := range colors {
		// Sum the RGB components to determine if the color is too bright or too dark.
		summedColor := uint16(c.R) + uint16(c.G) + uint16(c.B)

		if summedColor < maxBrightness && summedColor > minDarkness {
			// If we found a valid color just set it and break. We don't want to
			// check the other ones.
			return c
		}
	}
	// We haven't found a valid color, return the first one.
	return colors[0]
}

// FindN returns the first-N dominant colors in an image.
// If nClusters is less than or equal to 0, the value defaults to 4.
// Clusters are returned in their order of dominance.
func FindN(img image.Image, nClusters int) []color.RGBA {
	colors := FindWeight(img, nClusters)
	cols := []color.RGBA{}
	for _, c := range colors {
		cols = append(cols, c.RGBA)
	}
	return cols
}

func FindWeight(img image.Image, nClusters int) []Color {
	if nClusters <= 0 {
		nClusters = nClustersDefault
	}

	clusters, totalWeight := findClusters(img, nClusters)

	colors := []Color{}
	for _, c := range clusters {
		r, g, b := c.Centroid()
		colors = append(colors, Color{
			RGBA:   color.RGBA{R: r, G: g, B: b, A: 0xff},
			Weight: float64(c.weight) / totalWeight,
		})
	}
	return colors
}

// Hex returns a string representing the color in "#AABBCC" format.
func Hex(c color.RGBA) string {
	return "#" + fmt.Sprintf("%.2X%.2X%.2X", c.R, c.G, c.B)
}
