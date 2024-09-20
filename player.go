package main

import rl "github.com/gen2brain/raylib-go/raylib"

type PlayerShip struct {
	Position rl.Vector2
	Rotation float32
	Scale    float32
	Speed    float32
}

func (p PlayerShip) GetBoundingBox() rl.Rectangle {
	return rl.NewRectangle(p.Position.X, p.Position.Y, p.Scale, p.Scale)
}

func NewPlayerShip(position rl.Vector2, rotation float32, scale float32, speed float32) *PlayerShip {
	return &PlayerShip{Position: position, Rotation: rotation, Scale: scale, Speed: speed}
}
