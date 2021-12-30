package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func Normalize(c complex64) float64 {
	x := float64(real(c))
	y := float64(imag(c))
	return x*x + y*y
}

func mandelbrotPixel(pos rl.Vector2, xdim float32, ydim float32, maxIter int) rl.Color {
	c := complex(pos.Y, pos.X)
	c = c * complex(2.4/ydim, 0)
	c = c - complex(1.2*xdim/ydim+0.5, 1.2)
	z := c

	iter := 0
	for iter = 0; (iter < maxIter) && (Normalize(z) <= 4.0); iter++ {
		z = z*z + c
	}
	if iter == maxIter {
		return rl.Black
	}
	ci := 512 * iter / maxIter
	if iter < (maxIter / 2) {
		return rl.NewColor(uint8(ci), 0, 0, 255)
	} else {
		return rl.NewColor(255, uint8(ci)-255, uint8(ci)-255, 255)
	}
}

func juliaPixel(pos rl.Vector2, xdim float32, ydim float32, maxIter int) rl.Color {
	// Cantor dust
	const k = -0.6 + 0.6i

	c := complex(pos.Y, pos.X)
	c = c * complex(2.4/ydim, 0)
	c = c - complex(1.2*xdim/ydim+0.5, 1.2)
	z := c

	iter := 0
	for iter = 0; (iter < maxIter) && (Normalize(z) <= 4.0); iter++ {
		z = z*z + k
	}
	if iter == maxIter {
		return rl.Black
	}
	ci := 512 * iter / maxIter
	if iter < (maxIter / 2) {
		return rl.NewColor(0, 0, uint8(ci), 255)
	} else {
		return rl.NewColor(uint8(ci)-255, uint8(ci)-255, 255, 255)
	}

}

func generatePixels(width, height int) []byte {

	result := make([]byte, width*height*4)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pos := rl.NewVector2(float32(x), float32(y))
			mPixel := mandelbrotPixel(pos,
				float32(width),
				float32(height),
				30)
			jPixel := juliaPixel(pos,
				float32(width),
				float32(height),
				30)

			pixelColor := rl.Black
			if mPixel == rl.Black {
				pixelColor = jPixel
			} else {
				pixelColor = mPixel
			}

			index := (y*width + x) * 4
			result[index] = byte(pixelColor.R)
			result[index+1] = byte(pixelColor.G)
			result[index+2] = byte(pixelColor.B)
			result[index+3] = byte(pixelColor.A)

		}
	}
	return result
}

func cameraUpdate(camera *rl.Camera2D, cameraTarget rl.Vector2, windowWidth int32, windowHeight int32, previousMousePosition *rl.Vector2) rl.Vector2 {
	const cameraMoveSpeed = 8

	// Move camera with keyboard
	if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		cameraTarget.X += cameraMoveSpeed
	} else if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		cameraTarget.X -= cameraMoveSpeed
	} else if rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
		cameraTarget.Y -= cameraMoveSpeed
	} else if rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
		cameraTarget.Y += cameraMoveSpeed
	}

	// Move camera with mouse
	mousePos := rl.GetMousePosition()
	delta := rl.Vector2Subtract(*previousMousePosition, mousePos)
	*previousMousePosition = mousePos
	if rl.IsMouseButtonDown(0) {
		cameraTarget = rl.GetScreenToWorld2D(rl.Vector2Add(camera.Offset, delta), *camera)
	}

	return cameraTarget
}

func main() {

	const defaultWindowWidth = 800
	const defaultWindowHeight = 800

	rl.SetConfigFlags(rl.FlagWindowResizable)
	var windowWidth int32 = defaultWindowWidth
	var windowHeight int32 = defaultWindowHeight

	rl.SetConfigFlags(rl.FlagMsaa4xHint)

	rl.InitWindow(windowWidth, windowHeight, "Mandelbrot and Julia set visualizer")
	rl.SetTargetFPS(60)

	// Camera setup\
	cameraTarget := rl.NewVector2(float32(windowWidth/2.0), float32(windowHeight/2.0))
	camera := rl.Camera2D{}
	camera.Target = cameraTarget
	camera.Offset = rl.NewVector2(cameraTarget.X, cameraTarget.Y)
	camera.Rotation = 0.0
	camera.Zoom = 1.0

	mousePos := rl.GetMousePosition()

	// generate patterns

	const drawAreaWidth = 1000
	const drawAreaHeight = 1000

	pixels := generatePixels(drawAreaWidth, drawAreaHeight)
	image := rl.NewImage(pixels, drawAreaWidth, drawAreaHeight, 1, rl.UncompressedR8g8b8a8)
	texture := rl.LoadTextureFromImage(image)

	for !rl.WindowShouldClose() {

		// Update on resize
		windowWidth = int32(rl.GetScreenWidth())
		windowHeight = int32(rl.GetScreenHeight())

		cameraTarget = cameraUpdate(&camera, cameraTarget, windowWidth, windowHeight, &mousePos)
		camera.Target = cameraTarget

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginMode2D(camera)

		rl.DrawTexture(
			texture,
			windowWidth/2-texture.Width/2,
			windowHeight/2-texture.Height/2,
			rl.White)

		rl.EndMode2D()
		rl.EndDrawing()

	}

	rl.UnloadImage(image)
	rl.UnloadTexture(texture)
	rl.CloseWindow()
}
