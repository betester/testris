package entity

import (
	"errors"
	"tetris/matrix"
)

var (
	CLOCKWISE      = 1
	ANTI_CLOCKWISE = 2
)

var (
	I = 0
	J = 1
	L = 2
	O = 3
	S = 4
	T = 5
	Z = 6
)

const (
	RED    = 0
	BLUE   = 1
	YELLOW = 2
	GREEN  = 3
)

var BLOCK_OCCUPYING_LOCATION map[int]matrix.Matrix = map[int]matrix.Matrix{
	I: [][]int{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
	J: [][]int{{0, 0}, {0, 1}, {1, 1}, {2, 1}},
	L: [][]int{{0, 0}, {0, 1}, {-1, 1}, {-2, 1}},
	O: [][]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
	S: [][]int{{0, 0}, {-1, 0}, {-1, 1}, {-2, 1}},
	T: [][]int{{0, 0}, {0, 1}, {-1, 1}, {1, 1}},
	Z: [][]int{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
}

var BLOCK_CENTER map[int]int = map[int]int{
	I: 1,
	J: 1,
	L: 1,
	O: -1,
	S: 1,
	T: 1,
	Z: 2,
}

var ORIENTATION_ROTATION map[int]matrix.Matrix = map[int]matrix.Matrix{
	CLOCKWISE:      [][]int{{0, 1}, {-1, 0}},
	ANTI_CLOCKWISE: [][]int{{0, -1}, {1, 0}},
}

type BlockEntity struct {
	EntityType       int
	Color            int
	OccupiedPosition matrix.Matrix
}

func (b *BlockEntity) RotateBlock(orientation int) error {
	matrixRotation, ok := ORIENTATION_ROTATION[orientation]

	if !ok {
		return errors.New("Not a valid orientation")
	}

	center := BLOCK_CENTER[b.EntityType]
	blockCenter := make([]int, 2)
	if center != -1 {
		copy(blockCenter, b.OccupiedPosition[BLOCK_CENTER[b.EntityType]])
	}
	b.OccupiedPosition.Minus(blockCenter)
	transposedPosition := b.OccupiedPosition.Transpose()
	b.OccupiedPosition = matrix.Multiply(matrixRotation, transposedPosition).Transpose()
	b.OccupiedPosition.Add(blockCenter)

	return nil
}

func (b *BlockEntity) MoveBlock(direction [2]int) {
	currentLocation := b.OccupiedPosition

	for _, location := range currentLocation {
		location[0] += direction[0]
		location[1] += direction[1]
	}

	b.OccupiedPosition = currentLocation
}

func New(blockType int, color int, initialLocation [2]int) (BlockEntity, error) {
	var blockEntity BlockEntity
	taken_location, ok := BLOCK_OCCUPYING_LOCATION[blockType]

	occupying_location := make(matrix.Matrix, 4)

	if !ok {
		return blockEntity, errors.New("No block type found")
	}

	for i, position := range taken_location {
		occupying_location[i] = make([]int, 2)
		occupying_location[i][0] = position[0] + initialLocation[0]
		occupying_location[i][1] = position[1] + initialLocation[1]
	}

	blockEntity.Color = color
	blockEntity.EntityType = blockType
	blockEntity.OccupiedPosition = occupying_location
	return blockEntity, nil
}
