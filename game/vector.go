package game

import (
	"fmt"
	"math"
)

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

func (v Vec3) String() string {
	return fmt.Sprintf("(%v, %v, %v)", v.X, v.Y, v.Z)
}

func (v Vec3) AngleTo(other Vec3) float32 {
	f := float32(math.Atan2(float64(other.Y-v.Y), float64(other.X-v.X)))
	if f < 0.0 {
		f += math.Pi * 2.0
	}

	if f > math.Pi*2.0 {
		f -= math.Pi * 2.0
	}

	return f
}
