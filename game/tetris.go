package game

import (
	"math"
	"tetris/collision"
	"tetris/entity"
	eventhandler "tetris/event_handler"
	"tetris/spawner"
	renderer "tetris/ui"
	"time"
)

const (
	PAUSE = 1
	PLAY  = 2
	LOSE  = 4
)

const (
	CHANGE_LEVEL_DURATION_SECOND = 60
)

const (
	MOVING_BLOCK   = 0
	SPAWNING_BLOCK = 1
	BLOCK_STOPS    = 2
)

var DIRECTION_MAP map[int][2]int = map[int][2]int{
	eventhandler.RIGHT: {1, 1},
	eventhandler.LEFT:  {-1, 1},
}

var LEVEL_SPEED map[int]float64 = map[int]float64{
	0: 0.025,
	1: 0.030,
	2: 0.045,
	3: 0.07,
	4: 0.1,
}

type TetrisGame struct {
	MaxWitdh           int
	MaxHeight          int
	State              int
	BlockState         int
	Score              int
	Level              int
	speedUpMultiplier  int
	gainedScore        int
	CurrentBlock       *entity.BlockEntity
	Spawner            spawner.BlockSpawner
	CollisionDetector  collision.Collision
	Renderer           renderer.Renderer
	currentSpeed       float64 // could also probably use time, but to lazy for now
	blockColors        [][]int
	blockProjectionPos [][2]float32
	startTime          time.Time
}

func (tg *TetrisGame) Play() {
	tg.State = PLAY
	tg.BlockState = SPAWNING_BLOCK
	tg.blockColors = make([][]int, tg.MaxWitdh+1)
	tg.blockProjectionPos = make([][2]float32, 4)
	tg.startTime = time.Now()

	for i := range len(tg.blockColors) {
		tg.blockColors[i] = make([]int, tg.MaxHeight+1)
	}

	tg.Renderer.Init("Tetris")
	defer tg.Renderer.Close()

	for !tg.Renderer.ShouldClose() {
		tg.Update(tg.ReceiveEvent())
		tg.Render()
	}
}

func (tg *TetrisGame) Continue() {
	tg.State = PLAY
}

func (tg *TetrisGame) Update(event eventhandler.UpdateEvent) {

	tg.Level = int(math.Min(4, time.Now().Sub(tg.startTime).Seconds()/float64(CHANGE_LEVEL_DURATION_SECOND)))

	if tg.State == PAUSE || tg.State == LOSE {
		return
	}

	if event.GameState == PAUSE {
		tg.State = PAUSE
	} else if tg.BlockState == SPAWNING_BLOCK {
		block, err := tg.Spawner.Spawn()
		if err != nil {
			panic(err.Error())
		}
		tg.CurrentBlock = &block
		tg.BlockState = MOVING_BLOCK

		for i, location := range block.OccupiedPosition {
			tg.blockProjectionPos[i][0] = float32(location[0])
			tg.blockProjectionPos[i][1] = float32(tg.MaxHeight)
		}

	} else if tg.BlockState == MOVING_BLOCK {

		tg.currentSpeed += LEVEL_SPEED[tg.Level]
		if event.MovingDirection == eventhandler.DOWN {
			tg.currentSpeed += float64(tg.speedUpMultiplier)*LEVEL_SPEED[tg.Level] - LEVEL_SPEED[tg.Level] // cancels out the previous addition
		}
		baseDirection := [2]int{0, 1}
		direction, ok := DIRECTION_MAP[event.MovingDirection]
		if ok {
			baseDirection[0] = direction[0]
			baseDirection[1] = direction[1]
		}

		collisionOnSpawnPoint, collide, reachedBottom, outOfBounds := false, false, false, false
		collideVertically := false

		if tg.currentSpeed < 1 {
			baseDirection[1] = 0
		} else {
			tg.currentSpeed = 0
		}

		tg.CurrentBlock.RotateBlock(event.RotateDirection)

		for _, location := range tg.CurrentBlock.OccupiedPosition {
			// TODO: handle case for going down immediately
			x, y := location[0]+baseDirection[0], location[1]+baseDirection[1]
			reachedBottom = y >= tg.MaxHeight || reachedBottom
			_, uy, _ := tg.CollisionDetector.GetNonBlockingPosition(location[0], location[1])
			collideVertically = collideVertically || uy == y && y != -1
			collide = (tg.CollisionDetector.Collide(x, y)) || collide
			collisionOnSpawnPoint = (location[1] == 0 && collide) || collisionOnSpawnPoint
			outOfBounds = outOfBounds || !tg.CollisionDetector.ValidLocation(x, y)
		}

		if collisionOnSpawnPoint {
			tg.State = LOSE
			return
		}

		if (collide || outOfBounds) && !(collideVertically) {
			if event.RotateDirection == entity.CLOCKWISE {
				tg.CurrentBlock.RotateBlock(entity.ANTI_CLOCKWISE)
			} else if event.RotateDirection == entity.ANTI_CLOCKWISE {
				tg.CurrentBlock.RotateBlock(entity.CLOCKWISE)
			}
			return
		}

		if !outOfBounds && !collide {
			tg.CurrentBlock.MoveBlock(baseDirection)

			maxUpperBoundY, defaultUpperYMax, yUpperMax := 1000000, -1, -1
			for _, location := range tg.CurrentBlock.OccupiedPosition {
				_, uy, _ := tg.CollisionDetector.GetNonBlockingPosition(location[0], location[1])

				if maxUpperBoundY > uy-1 && uy != -1 {
					maxUpperBoundY = uy - 1
					yUpperMax = location[1]
				} else if maxUpperBoundY == uy-1 && yUpperMax < location[1] {
					yUpperMax = location[1]
				}

				if defaultUpperYMax <= location[1] {
					defaultUpperYMax = location[1]
				}
			}

			if maxUpperBoundY == 1000000 {
				maxUpperBoundY = tg.MaxHeight
				yUpperMax = defaultUpperYMax
			}

			for i, location := range tg.CurrentBlock.OccupiedPosition {
				tg.blockProjectionPos[i][0] = float32(location[0])
				tg.blockProjectionPos[i][1] = float32(maxUpperBoundY - (yUpperMax - location[1]))
			}

		}

		if reachedBottom || collideVertically {
			tg.BlockState = BLOCK_STOPS
		}

	} else if tg.BlockState == BLOCK_STOPS {

		totalRemoveBlock := 0

		for _, location := range tg.CurrentBlock.OccupiedPosition {
			tg.CollisionDetector.AddOccupiedBlocks(location[0], location[1])
			tg.blockColors[location[0]][location[1]] = tg.CurrentBlock.Color
		}

		for _, location := range tg.CurrentBlock.OccupiedPosition {
			y := location[1]
			// TODO: handle node deletion properly
			totalCoordinate := tg.CollisionDetector.GetYCount(y)

			if totalCoordinate == tg.MaxWitdh+1 {
				// pop all of the blocks from the tree
				error := tg.CollisionDetector.RemoveBlock(y)
				if error == nil {
					totalRemoveBlock += 1
				}
			}
		}

		tg.gainedScore = totalRemoveBlock * tg.MaxWitdh
		tg.Score += totalRemoveBlock * tg.MaxWitdh
		tg.BlockState = SPAWNING_BLOCK
		tg.CurrentBlock = nil
	}
}

func (tg *TetrisGame) Render() {

	if tg.State == PLAY {
		blocks := tg.CollisionDetector.GetAllBlocks()
		if tg.CurrentBlock != nil {
			for _, location := range tg.CurrentBlock.OccupiedPosition {
				blocks = append(blocks, [2]float32{float32(location[0]), float32(location[1])})
				tg.blockColors[location[0]][location[1]] = tg.CurrentBlock.Color
			}
		}
		projectionColor := -1
		if tg.CurrentBlock != nil {
			projectionColor = tg.CurrentBlock.Color
		}
		tg.Renderer.RenderPlay(blocks, tg.blockColors, tg.blockProjectionPos, projectionColor, tg.gainedScore, tg.Level, tg.Score, time.Now().Sub(tg.startTime))
		tg.gainedScore = 0
	} else if tg.State == LOSE {
		tg.Renderer.RenderLose(tg.Score)
	}
}

func (tg TetrisGame) ReceiveEvent() eventhandler.UpdateEvent {
	return eventhandler.HandleEvent()
}

func New(MaxWidth, MaxHeight int,
	CollisionDetector collision.Collision,
	Spawner spawner.BlockSpawner,
	Renderer renderer.Renderer,
	speedUpMultiplier int,
	level int) TetrisGame {
	return TetrisGame{
		MaxWitdh:          MaxWidth,
		MaxHeight:         MaxHeight,
		CollisionDetector: CollisionDetector,
		Spawner:           Spawner,
		BlockState:        SPAWNING_BLOCK,
		Level:             level,
		State:             PAUSE,
		Renderer:          Renderer,
		speedUpMultiplier: speedUpMultiplier,
	}
}
