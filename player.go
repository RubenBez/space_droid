package main

import rl "github.com/gen2brain/raylib-go/raylib"

type PlayerShip struct {
	Position     rl.Vector2
	Rotation     float32
	Scale        float32
	Speed        float32
	RenderPoints []rl.Vector2
}

func (p PlayerShip) GetBoundingBox() rl.Rectangle {
	return rl.NewRectangle(p.Position.X, p.Position.Y, p.Scale, p.Scale)
}

func (p PlayerShip) GetScaledRenderPoints() []rl.Vector2 {
	var transform = func(point rl.Vector2) rl.Vector2 {
		return rl.Vector2Add(rl.Vector2Scale(rl.Vector2Rotate(point, DegToRad(p.Rotation-90)), p.Scale), p.Position)
	}
	var newPoints = make([]rl.Vector2, len(p.RenderPoints))
	for i, point := range p.RenderPoints {
		newPoints[i] = transform(point)
	}

	if shouldDrawBoundingBox {
		for i := range newPoints {
			rl.DrawLineEx(newPoints[i], newPoints[(i+1)%len(newPoints)], 2, rl.Pink)
		}
	}

	return newPoints
}

func NewPlayerShip(position rl.Vector2, rotation float32, scale float32, speed float32) *PlayerShip {
	var p = &PlayerShip{Position: position, Rotation: rotation, Scale: scale, Speed: speed}
	p.RenderPoints = []rl.Vector2{
		rl.NewVector2(0.0, 0.5),
		rl.NewVector2(-0.5, -0.5),
		rl.NewVector2(-0.3, -0.2),
		rl.NewVector2(0.3, -0.2),
		rl.NewVector2(0.5, -0.5),
	}
	return p
}
