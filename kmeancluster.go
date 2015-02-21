// Copyright (c) 2011 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dominantcolor

import (
	"math"
)

type kMeanCluster struct {
	centroid [3]uint8

	// Holds the sum of all the points that make up this cluster. Used to
	// generate the next centroid as well as to check for convergence.
	aggregate [3]uint32
	counter   uint32

	// The weight of the cluster, determined by how many points were used
	// to generate the previous centroid.
	weight uint32
}

func (k *kMeanCluster) SetCentroid(r uint8, g uint8, b uint8) {
	k.centroid[0] = r
	k.centroid[1] = g
	k.centroid[2] = b
}

func (k *kMeanCluster) GetCentroid(r *uint8, g *uint8, b *uint8) {
	*r = k.centroid[0]
	*g = k.centroid[1]
	*b = k.centroid[2]
}

func (k *kMeanCluster) IsAtCentroid(r uint8, g uint8, b uint8) bool {
	return r == k.centroid[0] && g == k.centroid[1] && b == k.centroid[2]
}

// Recomputes the centroid of the cluster based on the aggregate data. The
// number of points used to calculate this center is stored for weighting
// purposes. The aggregate and counter are then cleared to be ready for the
// next iteration.
func (k *kMeanCluster) RecomputeCentroid() {
	if k.counter > 0 {
		k.centroid[0] = uint8(k.aggregate[0] / k.counter)
		k.centroid[1] = uint8(k.aggregate[1] / k.counter)
		k.centroid[2] = uint8(k.aggregate[2] / k.counter)

		k.aggregate[0] = 0
		k.aggregate[1] = 0
		k.aggregate[2] = 0
		k.weight = k.counter
		k.counter = 0
	}
}

func (k *kMeanCluster) AddPoint(r uint8, g uint8, b uint8) {
	k.aggregate[0] += uint32(r)
	k.aggregate[1] += uint32(g)
	k.aggregate[2] += uint32(b)
	k.counter++
}

// Just returns the distance^2. Since we are comparing relative distances
// there is no need to perform the expensive sqrt() operation.
func (k *kMeanCluster) GetDistanceSqr(r uint8, g uint8, b uint8) uint32 {
	return (uint32(r)-uint32(k.centroid[0]))*(uint32(r)-uint32(k.centroid[0])) +
		(uint32(g)-uint32(k.centroid[1]))*(uint32(g)-uint32(k.centroid[1])) +
		(uint32(b)-uint32(k.centroid[2]))*(uint32(b)-uint32(k.centroid[2]))
}

// In order to determine if we have hit convergence or not we need to see
// if the centroid of the cluster has moved. This determines whether or
// not the centroid is the same as the aggregate sum of points that will be
// used to generate the next centroid.
func (k *kMeanCluster) CompareCentroidWithAggregate() bool {
	if k.counter == 0 {
		return false
	}
	return uint8(k.aggregate[0]/uint32(k.counter)) == k.centroid[0] &&
		uint8(k.aggregate[1]/uint32(k.counter)) == k.centroid[1] &&
		uint8(k.aggregate[2]/uint32(k.counter)) == k.centroid[2]
}

type kMeanClusters []*kMeanCluster

func (a kMeanClusters) ContainsCentroid(r, g, b uint8) bool {
	for _, c := range a {
		if c.IsAtCentroid(r, g, b) {
			return true
		}
	}
	return false
}

func (a kMeanClusters) Closest(r, g, b uint8) *kMeanCluster {
	var closest *kMeanCluster
	var distanceToClosest uint32 = math.MaxUint32
	for _, c := range a {
		d := c.GetDistanceSqr(r, g, b)
		if d < distanceToClosest {
			distanceToClosest = d
			closest = c
		}
	}
	return closest
}

type byWeight kMeanClusters

func (a byWeight) Len() int           { return len(a) }
func (a byWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byWeight) Less(i, j int) bool { return a[i].weight > a[j].weight }
