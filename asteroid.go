package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type AsteroidSize int32

func (as AsteroidSize) Name() string {
	switch as {
	case Small:
		return "Small"
	case Medium:
		return "Medium"
	case Large:
		return "Large"
	}

	return "Unknown"
}

const (
	Small AsteroidSize = iota
	Medium
	Large
)

type Asteroid struct {
	Position     rl.Vector2
	Rotation     float32
	Scale        float32
	Speed        float32
	Size         AsteroidSize
	ShouldDelete bool
	RenderPoints []rl.Vector2
}

func (a *Asteroid) DrawInfo() {
	var text = "Size:" + a.Size.Name()
	var textSize = rl.MeasureTextEx(rl.GetFontDefault(), text, 10, 0)
	var boundingBox = a.GetBoundingBox()
	var position = rl.NewVector2(boundingBox.X+boundingBox.Width/2, (boundingBox.Y-boundingBox.Height/2)-textSize.Y)
	rl.DrawText(text, int32(position.X), int32(position.Y), 10, rl.RayWhite)
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
	var transform = func(point rl.Vector2) rl.Vector2 {
		return rl.Vector2Add(rl.Vector2Scale(rl.Vector2Rotate(point, DegToRad(a.Rotation)), a.Scale), a.Position)
	}
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
	var pos = transform(rl.NewVector2(left+height/2, up+width/2))
	return rl.NewRectangle(pos.X, pos.Y, width*a.Scale, height*a.Scale)
}

func (a *Asteroid) GenerateAsteroid() {
	var numPoints = rl.GetRandomValue(6, 10)
	var points = make([]rl.Vector2, numPoints)
	var sections = float32(360 / numPoints)
	var r float32 = 0
	for i := range numPoints {
		var pos = rl.NewVector2(float32(math.Cos(float64(DegToRad(r)))), float32(math.Sin(float64(DegToRad(r)))))
		pos = rl.Vector2Normalize(pos)
		var offset1 = float32(rl.GetRandomValue(-1, int32(a.Scale))) / a.Scale
		var offset2 = float32(rl.GetRandomValue(-1, int32(a.Scale))) / a.Scale
		pos = rl.Vector2Add(pos, rl.NewVector2(offset1, offset2))
		r += sections

		points[i] = rl.NewVector2(pos.X, pos.Y)
	}
	a.RenderPoints = points
}

func (a *Asteroid) GetScaleForSize() float32 {
	switch a.Size {
	case Small:
		return 8
	case Medium:
		return 15
	case Large:
		return 24
	}

	return 8
}

func NewAsteroid(position rl.Vector2, rotation float32, size AsteroidSize, speed float32) *Asteroid {
	var a = &Asteroid{Position: position, Rotation: rotation, Size: size, Speed: speed}
	a.Scale = a.GetScaleForSize()
	a.GenerateAsteroid()
	return a
}
