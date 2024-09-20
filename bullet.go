package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Bullet struct {
	Position     rl.Vector2
	Scale        float32
	Rotation     float32
	Speed        float32
	Lifetime     float32
	ShouldDelete bool
}

func (b Bullet) GetBoundingBox() rl.Rectangle {
	return rl.NewRectangle(b.Position.X, b.Position.Y, b.Scale, b.Scale)
}

func NewBullet(position rl.Vector2, scale float32, rotation float32, speed float32, lifetime float32) *Bullet {
	return &Bullet{Position: position, Scale: scale, Rotation: rotation, Speed: speed, Lifetime: lifetime}
}
