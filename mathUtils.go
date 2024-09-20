package main

import (
	"math"
)

const ratio = math.Pi / 180

func RadToDeg(rad float64) (deg float64) {
	deg = rad / (ratio)
	return
}

func DegToRad(deg float32) (rad float32) {
	rad = deg * (ratio)
	return
}
