package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
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

func GetRandomValueF(min int32, max int32) float32 {
	return float32(rl.GetRandomValue(min, max))
}

func GetRandomAngle() float32 {
	return GetRandomValueF(-360, 360)
}

func CheckCollisionRotatedRect(rectA rl.Rectangle, rotA float32, rectB rl.Rectangle, rotB float32) bool {
	var a1, a2, a3, a4 = GetPointsFromRect(rectA)
	var b1, b2, b3, b4 = GetPointsFromRect(rectB)

	var transform = func(point rl.Vector2, rotation float32, pos rl.Vector2) rl.Vector2 {
		return rl.Vector2Add(rl.Vector2Scale(rl.Vector2Rotate(rl.Vector2Subtract(point, pos), DegToRad(rotation)), 1), pos)
	}

	var aPos = rl.NewVector2(rectA.X, rectA.Y)
	var aPoints = []rl.Vector2{
		transform(a1, rotA, aPos),
		transform(a2, rotA, aPos),
		transform(a3, rotA, aPos),
		transform(a4, rotA, aPos),
	}
	var bPos = rl.NewVector2(rectB.X, rectB.Y)
	var bPoints = []rl.Vector2{
		transform(b1, rotB, bPos),
		transform(b2, rotB, bPos),
		transform(b3, rotB, bPos),
		transform(b4, rotB, bPos),
	}

	return CheckCollisionPoly(aPoints, bPoints)
}

func CheckCollisionPoly(pointsA []rl.Vector2, pointsB []rl.Vector2) bool {
	// go through each of the vertices, plus the next
	// vertex in the list
	var next = 0
	for current := 0; current < len(pointsA); current++ {
		// get next vertex in list
		// if we've hit the end, wrap around to 0
		next = current + 1
		if next == len(pointsA) {
			next = 0
		}

		// get the PVectors at our current position
		// this makes our if statement a little cleaner
		var vc = pointsA[current] // c for "current"
		var vn = pointsA[next]    // n for "next"

		var collision bool = CollisionPolyLine(pointsB, vc, vn)
		if collision {
			return true
		}

		// optional: check if the 2nd polygon is INSIDE the first
		collision = rl.CheckCollisionPointPoly(pointsB[0], pointsA)
		if collision {
			return true
		}
	}

	return false
}

func CollisionPolyLine(points []rl.Vector2, lineStart rl.Vector2, lineEnd rl.Vector2) bool {
	// go through each of the vertices, plus the next
	// vertex in the list
	var next = 0
	for current := 0; current < len(points); current++ {
		// get next vertex in list
		// if we've hit the end, wrap around to 0
		next = current + 1
		if next == len(points) {
			next = 0
		}

		var lineStart2 = points[current]
		var lineEnd2 = points[next]

		var hitPoint = rl.NewVector2(0, 0)
		var hit bool = rl.CheckCollisionLines(lineStart, lineEnd, lineStart2, lineEnd2, &hitPoint)
		if hit {
			return true
		}
	}

	// never got a hit
	return false
}

func GetPointsFromRect(rectangle rl.Rectangle) (rl.Vector2, rl.Vector2, rl.Vector2, rl.Vector2) {
	w := rectangle.Width / 2
	h := rectangle.Height / 2
	return rl.NewVector2(rectangle.X-w, rectangle.Y+h),
		rl.NewVector2(rectangle.X+w, rectangle.Y+h),
		rl.NewVector2(rectangle.X+w, rectangle.Y-h),
		rl.NewVector2(rectangle.X-w, rectangle.Y-h)
}

func GetPointsFromRectSlice(rectangle rl.Rectangle) []rl.Vector2 {
	var a1, a2, a3, a4 = GetPointsFromRect(rectangle)
	return []rl.Vector2{a1, a2, a3, a4}
}

func WrapCoordinates(position rl.Vector2) (newPos rl.Vector2) {
	newPos.X = position.X
	newPos.Y = position.Y
	if position.X < 0.0 {
		newPos.X = position.X + screenWidth
	}
	if position.X >= screenWidth {
		newPos.X = position.X - screenWidth
	}

	if position.Y < 0.0 {
		newPos.Y = position.Y + screenHeight
	}
	if position.Y >= screenHeight {
		newPos.Y = position.Y - screenHeight
	}
	return newPos
}
