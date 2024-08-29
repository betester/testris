package entity

import (
	"fmt"
	"testing"
	"tetris/matrix"
)

func TestNew(t *testing.T) {
	block, err := New(I, [2]int{1, 0})
	expectedPosition := matrix.Matrix{{1, 0}, {2, 0}, {3, 0}, {4, 0}}

	if err != nil {
		t.Error("Fail creating block")
		t.Fail()
	}

	if !block.OccupiedPosition.Equal(expectedPosition) {
		t.Error("Created block is not equal with the expected position")
		t.Fail()
	}
}

func TestMoveBlock(t *testing.T) {
	block, err := New(I, [2]int{0, 0})
	expectedPosition := matrix.Matrix{{1, 0}, {2, 0}, {3, 0}, {4, 0}}

	if err != nil {
		t.Error("Fail creating block")
		t.Fail()
	}

	block.MoveBlock([2]float64{1, 0})

	if !block.OccupiedPosition.Equal(expectedPosition) {
		t.Error("Moved block is not equal with the expected position")
		t.Fail()
	}
}

func TestRotateBlock(t *testing.T) {
	block, err := New(Z, [2]int{0, 0})
	currentPosition := matrix.Copy(block.OccupiedPosition)

	if err != nil {
		t.Error("Fail creating block")
		t.Fail()
	}

	t.Log(block.OccupiedPosition.ToString())
	block.RotateBlock(ANTI_CLOCKWISE)
	t.Log(block.OccupiedPosition.ToString())
	block.RotateBlock(CLOCKWISE)

	if !block.OccupiedPosition.Equal(currentPosition) {
		t.Error(
			fmt.Printf(
				"Rotated block is not equal with the expected position\nresult\n%s\nexpected\n%s",
				block.OccupiedPosition.ToString(),
				currentPosition.ToString(),
			),
		)
		t.Fail()
	}
}
