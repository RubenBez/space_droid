package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type GameState struct {
	Player *PlayerShip

	Bullets   []*Bullet
	Asteroids []*Asteroid

	Camera rl.Camera2D

	GameOver bool

	Paused bool
}

const shouldDrawBoundingBox = false
const shouldDrawAsteroidInfo = true

const screenWidth float32 = 800
const screenHeight float32 = 450

const screenCenterX = screenWidth / 2
const screenCenterY = screenHeight / 2

func main() {
	rl.InitWindow(int32(screenWidth), int32(screenHeight), "RayLib In Go")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var state = &GameState{
		Player:    nil,
		Bullets:   []*Bullet{},
		Asteroids: []*Asteroid{},
		Camera:    rl.NewCamera2D(rl.Vector2Zero(), rl.Vector2Zero(), 0, 1),
		GameOver:  false,
		Paused:    false,
	}

	RestartGame(state)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.BeginMode2D(state.Camera)

		rl.ClearBackground(rl.Black)

		if state.Paused {
			var size = rl.MeasureTextEx(rl.GetFontDefault(), "PAUSE", 20, 0)
			rl.DrawText("PAUSE", 800/2-int32(size.X)/2, 450/2-int32(size.Y)/2, 20, rl.Red)
		} else {
			if !state.GameOver {
				ProcessPlayer(state)
				ProcessBullets(state)
				ProcessAsteroids(state)
			} else {
				var size = rl.MeasureTextEx(rl.GetFontDefault(), "GAME OVER", 20, 0)
				rl.DrawText("GAME OVER", 800/2-int32(size.X)/2, 450/2-int32(size.Y)/2, 20, rl.Red)
			}
		}

		if rl.IsKeyPressed(rl.KeyP) {
			state.Paused = !state.Paused
		}
		if rl.IsKeyPressed(rl.KeyR) {
			RestartGame(state)
		}

		DrawPlayer(state.Player)
		for i := range state.Bullets {
			DrawBullet(state.Bullets[i])
		}
		for i := range state.Asteroids {
			DrawAsteroid(state.Asteroids[i])
			if shouldDrawAsteroidInfo {
				state.Asteroids[i].DrawInfo()
			}
		}
		ProcessCollision(state)
		DrawStats(state)
		rl.EndMode2D()
		rl.EndDrawing()
	}
}

func RestartGame(state *GameState) {
	state.GameOver = false
	state.Paused = false
	for i := range state.Bullets {
		state.Bullets[i] = nil
	}
	state.Bullets = []*Bullet{}

	for i := range state.Asteroids {
		state.Asteroids[i] = nil
	}
	state.Asteroids = []*Asteroid{}

	state.Player = nil
	state.Player = NewPlayerShip(rl.NewVector2(screenCenterX, screenCenterY), 0, 20.0, 4)

	//SpawnAsteroid(state, rl.NewVector2(150, 150), float32(0), 0, Small)
	//SpawnAsteroid(state, rl.NewVector2(150, 250), float32(0), 0, Medium)
	//SpawnAsteroid(state, rl.NewVector2(150, 350), float32(0), 0, Large)
	for range 10 {
		var size = AsteroidSize(rl.GetRandomValue(int32(Small), int32(Large)))
		var x = rl.GetRandomValue(0, int32(screenWidth)/2)
		var y = rl.GetRandomValue(0, int32(screenHeight)/2)
		var rotation = rl.GetRandomValue(0, 360)
		SpawnAsteroid(state, rl.NewVector2(float32(x), float32(y)), float32(rotation), 1, size)
	}
}

func GetRandomValueF(min int32, max int32) float32 {
	return float32(rl.GetRandomValue(min, max))
}

func GetRandomAngle() float32 {
	return GetRandomValueF(-360, 360)
}

func ProcessCollision(state *GameState) {
	for _, b := range state.Bullets {
		for _, a := range state.Asteroids {
			if CheckCollisionPoly(a.GetScaledRenderPoints(), GetPointsFromRectSlice(b.GetBoundingBox())) {

				if a.Size == Large {
					SpawnAsteroid(state, a.Position, a.Rotation+GetRandomAngle(), a.Speed, a.Size-1)
					SpawnAsteroid(state, a.Position, a.Rotation+GetRandomAngle(), a.Speed, a.Size-1)
					SpawnAsteroid(state, a.Position, a.Rotation+GetRandomAngle(), a.Speed, a.Size-1)
					SpawnAsteroid(state, a.Position, a.Rotation+GetRandomAngle(), a.Speed, a.Size-1)
				}
				if a.Size == Medium {
					SpawnAsteroid(state, a.Position, a.Rotation+GetRandomAngle(), a.Speed, a.Size-1)
					SpawnAsteroid(state, a.Position, a.Rotation+GetRandomAngle(), a.Speed, a.Size-1)
				}
				a.ShouldDelete = true
				b.ShouldDelete = true
			}
		}
	}

	for _, a := range state.Asteroids {
		if CheckCollisionPoly(state.Player.GetScaledRenderPoints(), a.GetScaledRenderPoints()) {
			state.GameOver = true
		}
	}
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

func DrawAsteroid(asteroid *Asteroid) {
	DrawLines(asteroid.Position, asteroid.Rotation, asteroid.Scale, asteroid.RenderPoints)

	DrawBoundingBox(asteroid.GetBoundingBox(), asteroid.Rotation)
}

func ProcessAsteroids(state *GameState) {
	for i := len(state.Asteroids) - 1; i >= 0; i-- {
		if state.Asteroids[i].ShouldDelete {
			state.Asteroids[i] = nil
			state.Asteroids = append(state.Asteroids[:i], state.Asteroids[i+1:]...)
		}
	}

	for i := range state.Asteroids {
		var theta = float64(DegToRad(state.Asteroids[i].Rotation))
		var direction = rl.NewVector2(float32(math.Cos(theta)), float32(math.Sin(theta)))
		state.Asteroids[i].Position = rl.Vector2Add(state.Asteroids[i].Position, rl.Vector2Multiply(direction, rl.NewVector2(state.Asteroids[i].Speed, state.Asteroids[i].Speed)))

		state.Asteroids[i].Position = WrapCoordinates(state.Asteroids[i].Position)
	}
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

func ProcessBullets(state *GameState) {
	// Cleanup before processing again
	for i := len(state.Bullets) - 1; i >= 0; i-- {
		if state.Bullets[i].ShouldDelete {
			state.Bullets[i] = nil
			state.Bullets = append(state.Bullets[:i], state.Bullets[i+1:]...)
		}
	}

	for i := range state.Bullets {
		state.Bullets[i].Lifetime -= rl.GetFrameTime()

		if state.Bullets[i].Lifetime <= 0 {
			state.Bullets[i].ShouldDelete = true
		}
		var theta = float64(DegToRad(state.Bullets[i].Rotation))
		var direction = rl.NewVector2(float32(math.Cos(theta)), float32(math.Sin(theta)))
		state.Bullets[i].Position = rl.Vector2Add(state.Bullets[i].Position, rl.Vector2Multiply(direction, rl.NewVector2(state.Bullets[i].Speed, state.Bullets[i].Speed)))
	}
}

func ProcessPlayer(state *GameState) {
	var player = state.Player

	var theta = float64(DegToRad(player.Rotation))
	var lookDirection = rl.NewVector2(float32(math.Cos(theta)), float32(math.Sin(theta)))
	if rl.IsKeyDown(rl.KeyW) {
		player.Position = rl.Vector2Add(player.Position, rl.Vector2Multiply(rl.Vector2Normalize(lookDirection), rl.NewVector2(player.Speed, player.Speed)))
	}

	if rl.IsKeyDown(rl.KeyA) {
		player.Rotation -= 3
	}

	if rl.IsKeyDown(rl.KeyD) {
		player.Rotation += 3
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		SpawnBullet(state, player.Position, player.Rotation, 8)
	}

	state.Player.Position = WrapCoordinates(state.Player.Position)
}

func DrawStats(state *GameState) {
	var y int32 = 10
	rl.DrawText(fmt.Sprintf("Number of Bullets: %d", len(state.Bullets)), 2, y, 10, rl.RayWhite)
	y += 10
	rl.DrawText(fmt.Sprintf("Number of astroids: %d", len(state.Asteroids)), 2, y, 10, rl.RayWhite)
	y += 10
	rl.DrawText(fmt.Sprintf("Player pos: %f.0, %f.0", state.Player.Position.X, state.Player.Position.Y), 2, y, 10, rl.RayWhite)
}

func SpawnBullet(state *GameState, spawnPosition rl.Vector2, rotation float32, speed float32) {
	var bullet = NewBullet(spawnPosition, 10, rotation, speed, 1)
	state.Bullets = append(state.Bullets, bullet)
}

func SpawnAsteroid(state *GameState, spawnPosition rl.Vector2, rotation float32, speed float32, size AsteroidSize) {
	var asteroid = NewAsteroid(spawnPosition, rotation, size, speed)
	state.Asteroids = append(state.Asteroids, asteroid)
}

func DrawBullet(bullet *Bullet) {
	DrawLines(bullet.Position, bullet.Rotation+45, bullet.Scale, []rl.Vector2{
		rl.NewVector2(0.5, 0.5),
		rl.NewVector2(0.5, -0.5),
		rl.NewVector2(-0.5, -0.5),
	})

	DrawBoundingBox(bullet.GetBoundingBox(), bullet.Rotation)
}

func DrawPlayer(player *PlayerShip) {
	DrawLines(player.Position, player.Rotation-90, player.Scale, player.RenderPoints)

	DrawBoundingBox(player.GetBoundingBox(), player.Rotation)
}

func DrawLines(position rl.Vector2, rotation float32, scale float32, points []rl.Vector2) {
	var transform = func(point rl.Vector2) rl.Vector2 {
		return rl.Vector2Add(rl.Vector2Scale(rl.Vector2Rotate(point, DegToRad(rotation)), scale), position)
	}
	for i := range points {
		rl.DrawLineEx(transform(points[i]), transform(points[(i+1)%len(points)]), 2, rl.White)
	}
}

func DrawBoundingBox(boundingBox rl.Rectangle, rotation float32) {
	var a1, a2, a3, a4 = GetPointsFromRect(boundingBox)

	var pos = rl.NewVector2(boundingBox.X, boundingBox.Y)

	var transform = func(point rl.Vector2) rl.Vector2 {
		return rl.Vector2Add(rl.Vector2Scale(rl.Vector2Rotate(rl.Vector2Subtract(point, pos), DegToRad(rotation)), 1), pos)
	}

	if shouldDrawBoundingBox {
		rl.DrawLineEx(transform(a1), transform(a2), 2, rl.Pink)
		rl.DrawLineEx(transform(a2), transform(a3), 2, rl.Pink)
		rl.DrawLineEx(transform(a3), transform(a4), 2, rl.Pink)
		rl.DrawLineEx(transform(a4), transform(a1), 2, rl.Pink)
	}
}
