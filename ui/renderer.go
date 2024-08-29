package renderer

import (
	"fmt"
	"tetris/entity"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	TEXT_SCORE_DURATION_SECOND = 2
)

var BLOCK_COLORS map[int]rl.Color = map[int]rl.Color{
	entity.RED:    rl.Red,
	entity.BLUE:   rl.Blue,
	entity.YELLOW: rl.Yellow,
	entity.GREEN:  rl.Green,
}

type Renderer struct {
	Height               int32
	Width                int32
	BlockXSize           int32
	BlockYSize           int32
	TotalHorizontalBlock int
	TotalVerticalBlock   int
	TargetFps            int32
	xOffset              int32
	yOffset              int32
	currentGainedScore   int
	timeGainedScore      time.Time
}

func (r *Renderer) Init(gameName string) {

	fmt.Println("Initializing game")

	r.xOffset = r.Width/2 - r.BlockXSize*int32(r.TotalHorizontalBlock)/2
	r.yOffset = r.Height/2 - r.BlockYSize*int32(r.TotalVerticalBlock)/2
	rl.InitWindow(r.Width, r.Height, gameName)
	rl.SetTargetFPS(r.TargetFps)
}

func (r *Renderer) RenderPlay(
	blockPositions [][2]float32,
	color [][]int,
	blockProjectionPos [][2]float32,
	currentBlockColor,
	gainedScore,
	level,
	currentScore int,
	elapsedTime time.Duration) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	for i := range r.TotalHorizontalBlock + 1 {
		for j := range r.TotalVerticalBlock + 1 {
			rl.DrawRectangleLines(
				r.BlockXSize*int32(i)+r.xOffset,
				r.BlockYSize*int32(j)+r.yOffset,
				r.BlockXSize,
				r.BlockYSize,
				rl.Gray,
			)
		}
	}

	for i := range blockProjectionPos {
		xPosition := float32(r.BlockXSize)*blockProjectionPos[i][0] + float32(r.xOffset)
		yPosition := float32(r.BlockYSize)*blockProjectionPos[i][1] + float32(r.yOffset)

		rl.DrawRectangleLines(
			int32(xPosition),
			int32(yPosition),
			r.BlockXSize,
			r.BlockYSize,
			BLOCK_COLORS[currentBlockColor],
		)
	}

	for i := range blockPositions {

		xPosition := float32(r.BlockXSize)*blockPositions[i][0] + float32(r.xOffset)
		yPosition := float32(r.BlockYSize)*blockPositions[i][1] + float32(r.yOffset)
		blockColor := color[int(blockPositions[i][0])][int(blockPositions[i][1])]

		rl.DrawRectangleV(
			rl.Vector2{X: float32(xPosition), Y: float32(yPosition)},
			rl.Vector2{X: float32(r.BlockXSize), Y: float32(r.BlockYSize)},
			BLOCK_COLORS[blockColor],
		)
	}
	if gainedScore > 0 {
		r.currentGainedScore = gainedScore
		r.RenderGainedScore(gainedScore)
		r.timeGainedScore = time.Now()
	} else if time.Now().Sub(r.timeGainedScore).Seconds() < TEXT_SCORE_DURATION_SECOND {
		r.RenderGainedScore(r.currentGainedScore)
	}

	r.RenderTimeElapsed(elapsedTime)
	r.RenderScore(currentScore)
	r.RenderLevel(level)
	rl.EndDrawing()
}

func (r Renderer) RenderGainedScore(gainedScore int) {
	rl.DrawText(fmt.Sprintf("+%d", gainedScore), r.Width/2-2, r.Height/4, 30, rl.White)
}

func (r Renderer) RenderLevel(level int) {
	rl.DrawText(fmt.Sprintf("Level: %d", level), r.Width/2-r.xOffset-30, r.Height/12, 20, rl.White)
}

func (r Renderer) RenderScore(score int) {
	rl.DrawText(fmt.Sprintf("Current Score: %d", score), r.Width/2-r.xOffset-30, r.Height/18, 20, rl.White)
}

func (r Renderer) RenderTimeElapsed(elapsedTime time.Duration) {
	rl.DrawText(fmt.Sprintf("Elapsed Time: %.0f:%.2f", elapsedTime.Minutes(), elapsedTime.Seconds()), r.Width/2+r.xOffset-100, r.Height/18, 20, rl.White)
}

func (r Renderer) RenderLose(score int) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.White)
	rl.DrawText(fmt.Sprintf("You lose with score %d", score), r.xOffset+30, r.Height/2, 20, rl.LightGray)
	rl.EndDrawing()
}

func (r Renderer) ShouldClose() bool {
	return rl.WindowShouldClose()
}

func (r Renderer) Close() {
	rl.CloseWindow()
}
