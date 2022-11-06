package main

import (
	"fmt"
	"math/rand"
	"time"
	_ "embed"

	"github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	ScreenWidth  int32
	ScreenHeight int32
	Cols         int32
	Rows         int32
	Pause        bool
	Cells        [][]*Cell
}

type Cell struct {
	Position rl.Vector2
	Size     rl.Vector2
	Next     bool
	Alive    bool
}

const scl int32 = 10

var fps = int32(10)

//go:embed icon.png
var iconC []byte

func main() {
	rand.Seed(time.Now().Unix())

	game := &Game{}
	game.Init(false)

	rl.InitWindow(game.ScreenWidth, game.ScreenHeight, "Conway's Life Game using Raylib-Go")
	rl.SetTargetFPS(fps)

	icon := rl.LoadImageFromMemory(".png", iconC, int32(len(iconC)))

	rl.SetWindowIcon(*icon)

	rl.UnloadImage(icon)

	for !rl.WindowShouldClose() {
		if !game.Pause {
			game.Update()
		}

		game.Input()

		game.Render()
	}

	rl.CloseWindow()
}

func (g *Game) Init(clear bool) {
	g.ScreenWidth = 800
	g.ScreenHeight = 600
	g.Cols = g.ScreenWidth / scl
	g.Rows = g.ScreenHeight / scl

	g.Cells = make([][]*Cell, g.Cols+1)
	for i := int32(0); i <= g.Cols; i++ {
		g.Cells[i] = make([]*Cell, g.Rows+1)
	}

	for x := int32(0); x <= g.Cols; x++ {
		for y := int32(0); y <= g.Rows; y++ {
			g.Cells[x][y] = &Cell{}
			g.Cells[x][y].Position = rl.NewVector2(float32(x*scl), float32(y*scl))
			g.Cells[x][y].Size = rl.NewVector2(float32(scl), float32(scl))
			if rand.Float64() < 0.1 && !clear {
				g.Cells[x][y].Alive = true
			}
		}
	}
}

func (g *Game) Input() {
	if rl.IsKeyPressed(rl.KeyR) {
		g.Init(false)
	}
	if rl.IsKeyPressed(rl.KeyC) {
		g.Init(true)
	}
	if rl.IsKeyPressed(rl.KeyRight) && g.Pause {
		g.Update()
	}
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		g.Click(rl.GetMouseX(), rl.GetMouseY())
		println("click")
	}
	if rl.IsKeyPressed(rl.KeySpace) {
		g.Pause = !g.Pause
	}
	if rl.IsKeyPressed(rl.KeyEqual) {
		fps += 10
		rl.SetTargetFPS(fps)
	}
	if rl.IsKeyPressed(rl.KeyMinus) {
		fps -= 10
		rl.SetTargetFPS(fps)
	}
}

func (g *Game) Update() {
	for i := int32(0); i <= g.Cols; i++ {
		for j := int32(0); j <= g.Rows; j++ {
			nc := g.CountNeighbors(i, j)
			if g.Cells[i][j].Alive {
				if nc < 2 {
					g.Cells[i][j].Next = false
				} else if nc > 3 {
					g.Cells[i][j].Next = false
				} else {
					g.Cells[i][j].Next = true
				}
			} else {
				if nc == 3 {
					g.Cells[i][j].Next = true
				}
			}
		}
	}

	for i := int32(0); i <= g.Cols; i++ {
		for j := int32(0); j < g.Rows; j++ {
			g.Cells[i][j].Alive = g.Cells[i][j].Next
		}
	}
}

func (g *Game) Render() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.Black)

	if g.Pause {
		v := rl.MeasureTextEx(rl.GetFontDefault(), "Paused", 80, 2)
		rl.DrawText("Paused", (g.ScreenWidth/2) - (int32(v.X)/2), (g.ScreenHeight/2) - (int32(v.Y)/2), 60, rl.DarkGray)
	} else {
		v := rl.MeasureTextEx(rl.GetFontDefault(), fmt.Sprintf("FPS: %v", fps), 80, 2)
		rl.DrawText(fmt.Sprintf("FPS: %v", fps), (g.ScreenWidth/2) - (int32(v.X)/2), (g.ScreenHeight/2) - (int32(v.Y)/2), 60, rl.DarkGray)
	}

	for i := range g.Cells {
		for j := range g.Cells[i] {
			if g.Cells[i][j].Alive {
				rl.DrawRectangleV(g.Cells[i][j].Position, g.Cells[i][j].Size, rl.Yellow)
			}
		}
	}

	rl.EndDrawing()
}

func (g *Game) Click(x, y int32) {
	for i := int32(0); i <= g.Cols; i++ {
		for j := int32(0); j <= g.Rows; j++ {
			cell := g.Cells[i][j].Position
			if int32(cell.X) < x && int32(cell.X)+scl > x && int32(cell.Y) < y && int32(cell.Y)+scl > y {
				g.Cells[i][j].Alive = !g.Cells[i][j].Alive
				println("update")
			}
		}
	}
}

func (g *Game) CountNeighbors(x, y int32) int {
	var count int

	for i := int32(-1); i < 2; i++ {
		for j := int32(-1); j < 2; j++ {
			if g.Cells[(x+i+g.Cols)%g.Cols][(y+j+g.Rows)%g.Rows].Alive {
				count++
			}
		}
	}

	if g.Cells[x][y].Alive {
		count--
	}

	return count
}
