package cluster

import (
	"fmt"
)

// Coordinates is a slice of float64
type Coordinates []float64

// Coordinates implements the Observation interface for a plain set of float64
// coordinates
func (c Coordinates) Coordinates() Coordinates {
	return Coordinates(c)
}

// Distance returns the euclidean distance between two coordinates
func (c Coordinates) Distance(p2 Coordinates) float64 {
	return elu(c, p2)
}

type Observation interface {
	Coordinates() Coordinates
	Distance(p2 Coordinates) float64
	GetID() string
}

type Observations []Observation

// Center returns the center coordinates of a set of Cluster
func (c Observations) Center() (Coordinates, error) {
	var l = len(c)
	if l == 0 {
		return nil, fmt.Errorf("there is no mean for an empty set of points")
	}

	cc := make([]float64, len(c[0].Coordinates()))
	for _, point := range c {
		for j, v := range point.Coordinates() {
			cc[j] += v
		}
	}

	var mean Coordinates
	for _, v := range cc {
		mean = append(mean, v/float64(l))
	}
	return mean, nil
}

// AverageDistance returns the average distance between o and all observations
func AverageDistance(o Observation, observations Observations) float64 {
	var d float64
	var l int

	for _, observation := range observations {
		dist := o.Distance(observation.Coordinates())
		if dist == 0 {
			continue
		}

		l++
		d += dist
	}

	if l == 0 {
		return 0
	}
	return d / float64(l)
}
