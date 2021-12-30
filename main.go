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

func main() {
	const defaultWindowWidth = 800
	const defaultWindowHeight = 800

	rl.SetConfigFlags(rl.FlagWindowResizable)
	var windowWidth int32 = defaultWindowWidth
	var windowHeight int32 = defaultWindowHeight

	rl.SetConfigFlags(rl.FlagMsaa4xHint)

	rl.InitWindow(windowWidth, windowHeight, "Mandelbrot and Julia set visualizer")
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		// Update on resize
		windowWidth = int32(rl.GetScreenWidth())
		windowHeight = int32(rl.GetScreenHeight())

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		var y int32
		var x int32
		for y = 0; y < windowHeight; y++ {
			for x = 0; x < windowWidth; x++ {
				pos := rl.Vector2{
					X: float32(x),
					Y: float32(y)}
				mPixel := mandelbrotPixel(pos,
					float32(windowWidth),
					float32(windowHeight),
					30)
				jPixel := juliaPixel(pos,
					float32(windowWidth),
					float32(windowHeight),
					30)

				pixelColor := rl.Black
				if mPixel == rl.Black {
					pixelColor = jPixel
				} else {
					pixelColor = mPixel
				}

				rl.DrawPixelV(pos, pixelColor)
			}
		}

		rl.EndDrawing()

	}

	rl.CloseWindow()
}
