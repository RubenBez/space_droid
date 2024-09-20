package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type Asteroid struct {
	Position     rl.Vector2
	Rotation     float32
	Scale        float32
	Speed        float32
	ShouldDelete bool
	RenderPoints []rl.Vector2
}

func (a *Asteroid) GetScaledRenderPoints() []rl.Vector2 {
	var transform = func(point rl.Vector2) rl.Vector2 {
		return rl.Vector2Add(rl.Vector2Scale(rl.Vector2Rotate(point, DegToRad(a.Rotation)), a.Scale), a.Position)
	}
	var newPoints = make([]rl.Vector2, len(a.RenderPoints))
	for i, point := range a.RenderPoints {
		newPoints[i] = transform(point)
	}

	if shouldDrawBoundingBox {
		for i := range newPoints {
			rl.DrawLineEx(newPoints[i], newPoints[(i+1)%len(newPoints)], 2, rl.Pink)
		}
	}

	return newPoints
}

func (a *Asteroid) GetBoundingBox() rl.Rectangle {
	var left float32 = 0
	var right float32 = 0
	var up float32 = 0
	var down float32 = 0
	for _, point := range a.RenderPoints {
		if point.X < left {
			left = point.X
		}
		if point.X > right {
			right = point.X
		}

		if point.Y > up {
			up = point.Y
		}
		if point.Y < down {
			down = point.Y
		}
	}

	var width = left - right
	var height = up - down
	return rl.NewRectangle(a.Position.X, a.Position.Y, width*a.Scale, height*a.Scale)
}

func (a *Asteroid) GenerateAsteroid() {
	var numPoints = rl.GetRandomValue(6, 10)
	var points = make([]rl.Vector2, numPoints)
	var sections = float32(360 / numPoints)
	var r float32 = 0
	for i := range numPoints {
		var pos = rl.NewVector2(float32(math.Cos(float64(DegToRad(r)))), float32(math.Sin(float64(DegToRad(r)))))
		pos = rl.Vector2Normalize(pos)
		var offset1 = float32(rl.GetRandomValue(-int32(a.Scale), int32(a.Scale))) / 100
		var offset2 = float32(rl.GetRandomValue(-int32(a.Scale), int32(a.Scale))) / 100
		pos = rl.Vector2Add(pos, rl.NewVector2(offset1, offset2))
		r += sections

		points[i] = rl.NewVector2(pos.X, pos.Y)
	}
	a.RenderPoints = points
}

func NewAsteroid(position rl.Vector2, rotation float32, scale float32, speed float32) *Asteroid {
	var a = &Asteroid{Position: position, Rotation: rotation, Scale: scale, Speed: speed}
	a.GenerateAsteroid()
	return a
}
