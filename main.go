package main

import (
	"math/rand"
	"tetris/collision"
	"tetris/game"
	"tetris/spawner"
	treecoordinate "tetris/tree_coordinate"
	renderer "tetris/ui"
	"time"
)

// TODO: fix bug agane on projection, it should decide based on the distance probably
// TODO: add sound

func main() {
	totalBlockHorizontal, totalVertical := 10, 20
	width := 600
	height := 800
	blockXSize, blockYSize := 30, 30
	speedUpMultiplier := 4
	coordinateTree := treecoordinate.New()
	collisionDetector := collision.Collision{MaxWitdh: totalBlockHorizontal, MaxHeight: totalVertical, OccupiedBlocks: coordinateTree}
	spawnerBlock := spawner.BlockSpawner{MaxWidth: totalBlockHorizontal, Randomizer: *rand.New(rand.NewSource(time.Now().Unix()))}
	raylibRenderer := renderer.Renderer{
		Height:               int32(height),
		Width:                int32(width),
		BlockXSize:           int32(blockXSize),
		BlockYSize:           int32(blockYSize),
		TotalHorizontalBlock: totalBlockHorizontal,
		TotalVerticalBlock:   totalVertical,
		TargetFps:            60,
	}
	tetrisGame := game.New(totalBlockHorizontal, totalVertical, collisionDetector, spawnerBlock, raylibRenderer, speedUpMultiplier, 1)

	tetrisGame.Play()
}
