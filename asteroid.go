package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Asteroid struct {
	Position     rl.Vector2
	Rotation     float32
	Scale        float32
	Speed        float32
	ShouldDelete bool
}

func (a Asteroid) GetBoundingBox() rl.Rectangle {
	return rl.NewRectangle(a.Position.X, a.Position.Y, a.Scale, a.Scale)
}

func NewAsteroid(position rl.Vector2, rotation float32, scale float32, speed float32) *Asteroid {
	return &Asteroid{Position: position, Rotation: rotation, Scale: scale, Speed: speed}
}
