package game

import "fmt"

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

func (v Vec3) String() string {
	return fmt.Sprintf("(%v, %v, %v)", v.X, v.Y, v.Z)
}
