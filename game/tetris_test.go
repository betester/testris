package game

import (
	"fmt"
	"math/rand"
	"testing"
	"tetris/collision"
	eventhandler "tetris/event_handler"
	"tetris/matrix"
	"tetris/spawner"
	treecoordinate "tetris/tree_coordinate"
)

func TestUpdateSpawningBlock(t *testing.T) {

	colisionDetector := collision.Collision{
		MaxWitdh:       100,
		MaxHeight:      100,
		OccupiedBlocks: treecoordinate.New(),
	}
	spawnerBlock := spawner.BlockSpawner{MaxWidth: 100, Randomizer: *rand.New(rand.NewSource(42069))}
	game := TetrisGame{
		MaxWitdh:          100,
		MaxHeight:         100,
		CollisionDetector: colisionDetector,
		Spawner:           spawnerBlock,
		BlockState:        SPAWNING_BLOCK,
		BlockSpeed:        1,
		State:             PAUSE,
	}
	game.Continue()

	game.Update(eventhandler.UpdateEvent{})

	if game.CurrentBlock == nil {
		t.Error("Current block should not be nil")
		t.Fail()
	}

	// each current block should be on collision detector
	if game.BlockState != MOVING_BLOCK {
		t.Error("The block state should be moving state after spawning the block")
		t.Fail()
	}

}

func TestUpdateMovingBlock(t *testing.T) {
	colisionDetector := collision.Collision{
		MaxWitdh:       100,
		MaxHeight:      100,
		OccupiedBlocks: treecoordinate.New(),
	}
	spawnerBlock := spawner.BlockSpawner{MaxWidth: 100, Randomizer: *rand.New(rand.NewSource(42069))}
	game := TetrisGame{
		MaxWitdh:          100,
		MaxHeight:         100,
		CollisionDetector: colisionDetector,
		Spawner:           spawnerBlock,
		BlockState:        SPAWNING_BLOCK,
		BlockSpeed:        1,
		State:             PAUSE,
	}

	game.Continue()

	game.Update(eventhandler.UpdateEvent{})
	previousState := matrix.Copy(game.CurrentBlock.OccupiedPosition)
	game.Update(eventhandler.UpdateEvent{})

	allTheSame := true
	currentState := game.CurrentBlock.OccupiedPosition

	for i := range previousState {
		for j := range previousState[i] {
			allTheSame = allTheSame && previousState[i][j] == currentState[i][j]
		}
	}

	if allTheSame {
		t.Error("The block should have moved if there are no blocking path")
		t.Fail()
	}
}

func TestUpdateMovingBlockUntilDownShouldStop(t *testing.T) {

	colisionDetector := collision.Collision{
		MaxWitdh:       100,
		MaxHeight:      100,
		OccupiedBlocks: treecoordinate.New(),
	}
	spawnerBlock := spawner.BlockSpawner{MaxWidth: 100, Randomizer: *rand.New(rand.NewSource(42069))}
	game := TetrisGame{
		MaxWitdh:          100,
		MaxHeight:         100,
		CollisionDetector: colisionDetector,
		Spawner:           spawnerBlock,
		BlockState:        SPAWNING_BLOCK,
		BlockSpeed:        1,
		State:             PAUSE,
	}

	game.Continue()

	for range 100 {
		game.Update(eventhandler.UpdateEvent{})
	}

	if game.BlockState != BLOCK_STOPS {
		t.Errorf("Game state should stop found %d", game.BlockState)
		t.Fail()
	}

	game.Update(eventhandler.UpdateEvent{})

	totalBlockCount := game.CollisionDetector.GetTotalCount()
	if totalBlockCount != 4 {
		t.Errorf("Total count should be 4 instead found %d", totalBlockCount)
		t.Fail()
	}

}

func TestUpdateMovingBlockOutOfBoundsShouldStayTheSame(t *testing.T) {

	colisionDetector := collision.Collision{
		MaxWitdh:       100,
		MaxHeight:      100,
		OccupiedBlocks: treecoordinate.New(),
	}
	spawnerBlock := spawner.BlockSpawner{MaxWidth: 100, Randomizer: *rand.New(rand.NewSource(42069))}
	game := TetrisGame{
		MaxWitdh:          100,
		MaxHeight:         100,
		CollisionDetector: colisionDetector,
		Spawner:           spawnerBlock,
		BlockState:        SPAWNING_BLOCK,
		BlockSpeed:        1,
		State:             PAUSE,
	}
	game.Continue()

	for range 34 {
		game.Update(eventhandler.UpdateEvent{MovingDirection: eventhandler.RIGHT})
	}

	for _, location := range game.CurrentBlock.OccupiedPosition {
		if !colisionDetector.ValidLocation(location[0], location[1]) {
			t.Errorf("Invalid block location x: %d, y: %d should not exceed the bound limit", location[0], location[1])
			t.Fail()
		}
	}
}

func TestUpdateMovingCollidingWithExistingBlockShouldStop(t *testing.T) {

	colisionDetector := collision.Collision{
		MaxWitdh:       100,
		MaxHeight:      100,
		OccupiedBlocks: treecoordinate.New(),
	}
	spawnerBlock := spawner.BlockSpawner{MaxWidth: 100, Randomizer: *rand.New(rand.NewSource(42069))}
	game := TetrisGame{
		MaxWitdh:          100,
		MaxHeight:         100,
		CollisionDetector: colisionDetector,
		Spawner:           spawnerBlock,
		BlockState:        SPAWNING_BLOCK,
		BlockSpeed:        1,
		State:             PAUSE,
	}
	game.Continue()

	// spawns the current block
	game.Update(eventhandler.UpdateEvent{})

	previousState := matrix.Copy(game.CurrentBlock.OccupiedPosition)

	for range 99 {
		game.Update(eventhandler.UpdateEvent{})
	}

	if game.BlockState != BLOCK_STOPS {
		t.Errorf("Block state should have stop after 100 iteration")
		t.Fail()
	}
	game.Update(eventhandler.UpdateEvent{})

	if game.BlockState != SPAWNING_BLOCK {
		t.Errorf("Game state after stopping the block should be spawning a new block")
		t.Fail()
	}

	game.Update(eventhandler.UpdateEvent{})
	game.CurrentBlock.OccupiedPosition = previousState

	for range 98 {
		game.Update(eventhandler.UpdateEvent{})
	}

	if game.BlockState != BLOCK_STOPS {
		t.Errorf("Block state should have changed to stop due to colliding with other block")
		t.Fail()
	}
}

func TestStoppingBlockShouldPopBlockIfAllOccupied(t *testing.T) {

	colisionDetector := collision.Collision{
		MaxWitdh:       16,
		MaxHeight:      100,
		OccupiedBlocks: treecoordinate.New(),
	}
	spawnerBlock := spawner.BlockSpawner{MaxWidth: 16, Randomizer: *rand.New(rand.NewSource(42069))}
	game := TetrisGame{
		MaxWitdh:          16,
		MaxHeight:         100,
		CollisionDetector: colisionDetector,
		Spawner:           spawnerBlock,
		BlockState:        SPAWNING_BLOCK,
		BlockSpeed:        1,
		State:             PAUSE,
	}
	game.Continue()

	i := 0
	for i <= 12 {
		game.Update(eventhandler.UpdateEvent{})
		game.CurrentBlock.OccupiedPosition[0][0] = i
		game.CurrentBlock.OccupiedPosition[1][0] = i + 1
		game.CurrentBlock.OccupiedPosition[2][0] = i + 2
		game.CurrentBlock.OccupiedPosition[3][0] = i + 3

		game.CurrentBlock.OccupiedPosition[0][1] = 0
		game.CurrentBlock.OccupiedPosition[1][1] = 0
		game.CurrentBlock.OccupiedPosition[2][1] = 0
		game.CurrentBlock.OccupiedPosition[3][1] = 0

		t.Log(game.CurrentBlock.OccupiedPosition.ToString())

		for range 101 {
			game.Update(eventhandler.UpdateEvent{})
		}

		game.Update(eventhandler.UpdateEvent{})

		t.Log(game.CurrentBlock.OccupiedPosition.ToString())
		i += 4
	}

	if game.Score != 16 {
		t.Errorf("Current score should be 16, found %d instead", game.Score)
		t.Fail()
	}

	for i := range 16 {
		if colisionDetector.Collide(i, 100) {
			t.Errorf("Popped position of %d 100 should not exist as it already removed", i)
			t.Fail()
		}
	}
}
