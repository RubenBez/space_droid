package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type State int32

const (
	Menu State = iota
	Instructions
	Game
)

type GameData struct {
	Player *PlayerShip

	Bullets   []*Bullet
	Asteroids []*Asteroid

	Camera rl.Camera2D

	GameRunning bool
	GameOver    bool
	Paused      bool
	Win         bool

	GameState State

	MenuIndex int32

	FxShoot           rl.Sound
	FxAsteroidDestroy rl.Sound
	FxSpaceShipDead   rl.Sound
	FxWin             rl.Sound
}

const shouldDrawBoundingBox = false
const shouldDrawAsteroidInfo = false
const shouldDrawStats = false

const screenWidth float32 = 800
const screenHeight float32 = 450

const screenCenterX = screenWidth / 2
const screenCenterY = screenHeight / 2

func main() {
	rl.InitWindow(int32(screenWidth), int32(screenHeight), "RayLib In Go")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyNull)

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	var data = &GameData{
		Player:            nil,
		Bullets:           []*Bullet{},
		Asteroids:         []*Asteroid{},
		Camera:            rl.NewCamera2D(rl.Vector2Zero(), rl.Vector2Zero(), 0, 1),
		GameRunning:       true,
		GameOver:          false,
		Paused:            false,
		Win:               false,
		GameState:         Menu,
		MenuIndex:         0,
		FxShoot:           rl.LoadSound("assets/audio/shoot.wav"),
		FxAsteroidDestroy: rl.LoadSound("assets/audio/asteroid_destroy.wav"),
		FxSpaceShipDead:   rl.LoadSound("assets/audio/space_ship_dead.wav"),
		FxWin:             rl.LoadSound("assets/audio/win.wav"),
	}
	defer rl.UnloadSound(data.FxShoot)
	defer rl.UnloadSound(data.FxAsteroidDestroy)
	defer rl.UnloadSound(data.FxSpaceShipDead)
	defer rl.UnloadSound(data.FxWin)

	RestartGame(data)

	for data.GameRunning {
		rl.BeginDrawing()
		rl.BeginMode2D(data.Camera)
		rl.ClearBackground(rl.Black)

		switch data.GameState {
		case Menu:
			ProcessMenuState(data)
		case Instructions:
			ProcessInstructionsState(data)
		case Game:
			ProcessGameState(data)
		}

		if rl.WindowShouldClose() {
			data.GameRunning = false
		}

		rl.EndMode2D()
		rl.EndDrawing()
	}
}

func ProcessInstructionsState(data *GameData) {
	DrawTextCenter("Instructions", 70, 42, rl.Green)

	var y float32 = 200
	DrawTextCenter("Use W,A,S,D or ARROW KEYS to move and 'SPACE' to shoot", y, 18, rl.White)
	y += 20
	DrawTextCenter("Destroy all the asteroids to win and don't get hit by one", y, 18, rl.White)

	y = 400
	DrawMenuItem("Back", y, true)

	if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) {
		data.GameState = Menu
	}
}

func ProcessMenuState(data *GameData) {
	DrawTextCenter("SPACE DROID", 70, 42, rl.Green)

	var y float32 = 150

	DrawMenuItem("Play", y, data.MenuIndex == 0)
	y += 50
	DrawMenuItem("Instructions", y, data.MenuIndex == 1)
	y += 50
	DrawMenuItem("Quit", y, data.MenuIndex == 2)

	y = 400
	DrawTextCenter("Made by Ruben Bezuidenhout", y, 18, rl.Blue)

	if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) {
		data.MenuIndex++
		data.MenuIndex %= 3
	}

	if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) {
		data.MenuIndex--
		if data.MenuIndex < 0 {
			data.MenuIndex = 2
		}
		data.MenuIndex %= 3
	}

	if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) {
		if data.MenuIndex == 0 {
			RestartGame(data)
			data.GameState = Game
		}

		if data.MenuIndex == 1 {
			data.GameState = Instructions
		}

		if data.MenuIndex == 2 {
			data.GameRunning = false
		}
	}
}

func DrawMenuItem(text string, y float32, selected bool) {
	var size = rl.MeasureTextEx(rl.GetFontDefault(), text, 16, 0)
	var px = screenWidth/2 - size.X/2
	var py = y
	rl.DrawText(text, int32(px), int32(py), 16, rl.RayWhite)

	if selected {
		rl.DrawLineEx(rl.NewVector2(px-10, py+size.Y/2), rl.NewVector2(px-20, (py+size.Y/2)-10), 2, rl.RayWhite)
		rl.DrawLineEx(rl.NewVector2(px-10, py+size.Y/2), rl.NewVector2(px-20, (py+size.Y/2)+10), 2, rl.RayWhite)
	}

}

func DrawTextCenter(text string, y float32, fontSize int32, color rl.Color) {
	var size = rl.MeasureTextEx(rl.GetFontDefault(), text, float32(fontSize), 0)
	rl.DrawText(text, int32(screenWidth/2-size.X/2), int32(y), fontSize, color)
}

func ProcessGameState(data *GameData) {
	if data.Paused {
		DrawTextCenter("PAUSE", screenHeight/2, 20, rl.Red)
	} else {
		if !data.Win {
			if !data.GameOver {
				ProcessPlayer(data)
				ProcessBullets(data)
				ProcessAsteroids(data)
				ProcessCollision(data)
			} else {
				DrawTextCenter("GAME OVER", screenHeight/2, 20, rl.Red)
				DrawTextCenter("PRESS 'R' TO TRY AGAIN", (screenHeight+40)/2, 20, rl.Red)
			}
		} else {
			DrawTextCenter("YOU WON!!", screenHeight/2, 20, rl.Gold)
			DrawTextCenter("PRESS 'R' TO TRY AGAIN", (screenHeight+40)/2, 20, rl.Gold)
		}
	}

	if rl.IsKeyPressed(rl.KeyP) {
		data.Paused = !data.Paused
	}

	if rl.IsKeyPressed(rl.KeyEscape) {
		data.GameState = Menu
	}

	if rl.IsKeyPressed(rl.KeyR) {
		RestartGame(data)
	}

	DrawPlayer(data.Player)
	for i := range data.Bullets {
		DrawBullet(data.Bullets[i])
	}
	for i := range data.Asteroids {
		DrawAsteroid(data.Asteroids[i])
		if shouldDrawAsteroidInfo {
			data.Asteroids[i].DrawInfo()
		}
	}

	if shouldDrawStats {
		DrawStats(data)
	}
}

func RestartGame(data *GameData) {
	data.GameOver = false
	data.Paused = false
	data.Win = false
	for i := range data.Bullets {
		data.Bullets[i] = nil
	}
	data.Bullets = []*Bullet{}

	for i := range data.Asteroids {
		data.Asteroids[i] = nil
	}
	data.Asteroids = []*Asteroid{}

	data.Player = nil
	data.Player = NewPlayerShip(rl.NewVector2(screenCenterX, screenCenterY), 0, 20.0, 2)

	//SpawnAsteroid(data, rl.NewVector2(150, 150), float32(0), 0, Small)
	//SpawnAsteroid(data, rl.NewVector2(150, 250), float32(0), 0, Medium)
	//SpawnAsteroid(data, rl.NewVector2(150, 350), float32(0), 0, Large)
	for range 10 {
		var size = AsteroidSize(rl.GetRandomValue(int32(Small), int32(Large)))
		var x = rl.GetRandomValue(0, int32(screenWidth)/2)
		var y = rl.GetRandomValue(0, int32(screenHeight)/2)
		var rotation = rl.GetRandomValue(0, 360)
		SpawnAsteroid(data, rl.NewVector2(float32(x), float32(y)), float32(rotation), 1, size)
	}
}

func ProcessCollision(data *GameData) {
	for _, b := range data.Bullets {
		for _, a := range data.Asteroids {
			if CheckCollisionPoly(a.GetScaledRenderPoints(), GetPointsFromRectSlice(b.GetBoundingBox())) {
				if a.Size == Large {
					SpawnAsteroid(data, a.Position, a.Rotation+GetRandomAngle(), a.Speed+GetRandomValueF(0, 5)/5, a.Size-1)
					SpawnAsteroid(data, a.Position, a.Rotation+GetRandomAngle(), a.Speed+GetRandomValueF(0, 5)/5, a.Size-1)
					SpawnAsteroid(data, a.Position, a.Rotation+GetRandomAngle(), a.Speed+GetRandomValueF(0, 5)/5, a.Size-1)
					SpawnAsteroid(data, a.Position, a.Rotation+GetRandomAngle(), a.Speed+GetRandomValueF(0, 5)/5, a.Size-1)
				}
				if a.Size == Medium {
					SpawnAsteroid(data, a.Position, a.Rotation+GetRandomAngle(), a.Speed+GetRandomValueF(0, 3)/3, a.Size-1)
					SpawnAsteroid(data, a.Position, a.Rotation+GetRandomAngle(), a.Speed+GetRandomValueF(0, 3)/3, a.Size-1)
				}
				rl.PlaySound(data.FxAsteroidDestroy)
				a.ShouldDelete = true
				b.ShouldDelete = true
			}
		}
	}

	for _, a := range data.Asteroids {
		if CheckCollisionPoly(data.Player.GetScaledRenderPoints(), a.GetScaledRenderPoints()) {
			rl.PlaySound(data.FxSpaceShipDead)
			data.GameOver = true
		}
	}
}

func DrawAsteroid(asteroid *Asteroid) {
	DrawLines(asteroid.Position, asteroid.Rotation, asteroid.Scale, asteroid.RenderPoints)

	DrawBoundingBox(asteroid.GetBoundingBox(), asteroid.Rotation)
}

func ProcessAsteroids(data *GameData) {
	for i := len(data.Asteroids) - 1; i >= 0; i-- {
		if data.Asteroids[i].ShouldDelete {
			data.Asteroids[i] = nil
			data.Asteroids = append(data.Asteroids[:i], data.Asteroids[i+1:]...)
		}
	}

	for i := range data.Asteroids {
		var theta = float64(DegToRad(data.Asteroids[i].Rotation))
		var direction = rl.NewVector2(float32(math.Cos(theta)), float32(math.Sin(theta)))
		data.Asteroids[i].Position = rl.Vector2Add(data.Asteroids[i].Position, rl.Vector2Multiply(direction, rl.NewVector2(data.Asteroids[i].Speed, data.Asteroids[i].Speed)))

		data.Asteroids[i].Position = WrapCoordinates(data.Asteroids[i].Position)
	}
}

func ProcessBullets(data *GameData) {
	// Cleanup before processing again
	for i := len(data.Bullets) - 1; i >= 0; i-- {
		if data.Bullets[i].ShouldDelete {
			data.Bullets[i] = nil
			data.Bullets = append(data.Bullets[:i], data.Bullets[i+1:]...)
		}
	}

	for i := range data.Bullets {
		data.Bullets[i].Lifetime -= rl.GetFrameTime()

		if data.Bullets[i].Lifetime <= 0 {
			data.Bullets[i].ShouldDelete = true
		}
		var theta = float64(DegToRad(data.Bullets[i].Rotation))
		var direction = rl.NewVector2(float32(math.Cos(theta)), float32(math.Sin(theta)))
		data.Bullets[i].Position = rl.Vector2Add(data.Bullets[i].Position, rl.Vector2Multiply(direction, rl.NewVector2(data.Bullets[i].Speed, data.Bullets[i].Speed)))
	}
}

func ProcessPlayer(data *GameData) {
	var player = data.Player
	const drag = 0.015

	var theta = float64(DegToRad(player.Rotation))
	var lookDirection = rl.NewVector2(float32(math.Cos(theta)), float32(math.Sin(theta)))
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		player.Velocity = rl.Vector2Add(player.Velocity, rl.Vector2Scale(lookDirection, player.Speed*rl.GetFrameTime()))
	}

	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		player.Rotation -= 3
	}

	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		player.Rotation += 3
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		rl.PlaySound(data.FxShoot)
		SpawnBullet(data, player.Position, player.Rotation, 8)
	}

	player.Velocity = rl.Vector2Scale(player.Velocity, 1-drag)
	player.Position = rl.Vector2Add(player.Position, player.Velocity)

	data.Player.Position = WrapCoordinates(data.Player.Position)

	if len(data.Asteroids) == 0 {
		rl.PlaySound(data.FxWin)
		data.Win = true
	}
}

func DrawStats(data *GameData) {
	var y int32 = 10
	rl.DrawText(fmt.Sprintf("Number of Bullets: %d", len(data.Bullets)), 2, y, 10, rl.RayWhite)
	y += 10
	rl.DrawText(fmt.Sprintf("Number of astroids: %d", len(data.Asteroids)), 2, y, 10, rl.RayWhite)
	y += 10
	rl.DrawText(fmt.Sprintf("Player pos: %f.0, %f.0", data.Player.Position.X, data.Player.Position.Y), 2, y, 10, rl.RayWhite)
}

func SpawnBullet(data *GameData, spawnPosition rl.Vector2, rotation float32, speed float32) {
	var bullet = NewBullet(spawnPosition, 10, rotation, speed, 1)
	data.Bullets = append(data.Bullets, bullet)
}

func SpawnAsteroid(data *GameData, spawnPosition rl.Vector2, rotation float32, speed float32, size AsteroidSize) {
	var asteroid = NewAsteroid(spawnPosition, rotation, size, speed)
	data.Asteroids = append(data.Asteroids, asteroid)
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
