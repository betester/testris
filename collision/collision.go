package collision

import (
	"errors"
	"fmt"
	treecoordinate "tetris/tree_coordinate"
)

type Collision struct {
	MaxWitdh       int
	MaxHeight      int
	OccupiedBlocks treecoordinate.CoordinateTree
}

func (c Collision) ValidLocation(x, y int) bool {
	return x >= 0 && y >= 0 && x <= c.MaxWitdh && y <= c.MaxHeight
}

// Will find the possible position location for the block to be located near the blocked position
func (c Collision) GetNonBlockingPosition(x, y int) (int, int, error) {
	lx, ly, err := c.OccupiedBlocks.UpperBound(x, y)

	if err != nil {
		return lx, ly, err
	}

	if !c.ValidLocation(lx, ly) {
		return -1, -1, errors.New(fmt.Sprintf("Not a valid location for non blocking position x: %d y: %d", x, y))
	}

	return lx, ly, nil
}

func (c Collision) Collide(x, y int) bool {
	return c.OccupiedBlocks.LocationExist(x, y)
}

func (c Collision) GetYCount(y int) int {
	return c.OccupiedBlocks.Count(y)
}

func (c Collision) GetTotalCount() int {
	return c.OccupiedBlocks.TotalCount()
}

func (c Collision) RemoveBlock(y int) error {
	// probably will be really slow, but hey, what can you do about it?
	for i := range c.MaxWitdh + 1 {
		c.RemoveOccupiedBlocks(i, y)
	}

	removalError := c.OccupiedBlocks.RemoveAll(y)

	j := y - 1
	removedBlocks := make([][2]int, 0)

	for j >= 0 {
		for i := range c.MaxWitdh + 1 {
			if c.Collide(i, j) {
				c.RemoveOccupiedBlocks(i, j)
				removedBlocks = append(removedBlocks, [2]int{i, j})
			}
		}
		c.OccupiedBlocks.RemoveAll(j)
		j -= 1
	}

	fmt.Println(len(removedBlocks))
	for _, val := range removedBlocks {
		c.AddOccupiedBlocks(val[0], val[1]+1)
	}

	return removalError
}

func (c Collision) RemoveOccupiedBlocks(x, y int) error {

	if !c.ValidLocation(x, y) && !c.OccupiedBlocks.LocationExist(x, y) {
		return errors.New(fmt.Sprintf("Not a valid location for removing position x: %d y: %d", x, y))
	}

	err := c.OccupiedBlocks.Remove(x, y)

	return err
}

func (c *Collision) AddOccupiedBlocks(x, y int) error {
	if !c.ValidLocation(x, y) && c.OccupiedBlocks.LocationExist(x, y) {
		return errors.New(fmt.Sprintf("Not a valid location for adding position x: %d y: %d", x, y))
	}

	c.OccupiedBlocks.Add(x, y)

	return nil
}

func (c Collision) GetAllBlocks() [][2]float32 {
	return c.OccupiedBlocks.GetAllCoordinate()
}
