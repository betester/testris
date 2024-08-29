package treecoordinate

import (
	"errors"
	"fmt"
	"math"
)

type TreeNode struct {
	value int
	count int
	left  *TreeNode
	right *TreeNode
	tree  *Tree
}

type Tree struct {
	root *TreeNode
}

type CoordinateTree struct {
	yTree      map[int]*Tree
	yAxisCount map[int]int
}

func add(root *TreeNode, value int) *TreeNode {
	if root == nil {
		return &TreeNode{
			value: value,
			count: 1,
		}
	} else if root.value < value {
		root.right = add(root.right, value)
	} else {
		root.left = add(root.left, value)
	}

	root.count += 1
	return root
}

func find(root *TreeNode, value int) *TreeNode {
	if root == nil {
		return nil
	} else if root.value == value {
		return root
	} else if root.value < value {
		return find(root.right, value)
	}

	return find(root.left, value)
}

func remove(root *TreeNode, value int) *TreeNode {
	if root == nil {
		return nil
	} else if root.value == value {
		if root.left == nil {
			return root.right
		} else if root.right == nil {
			return root.left
		} else {
			successor := findMinimumNode(root.right)
			remove(root, successor.value)
			successor.left = root.left
			successor.right = root.right
			return successor
		}
	} else if root.value < value {
		root.right = remove(root.right, value)
	} else {
		root.left = remove(root.left, value)
	}

	root.count -= 1

	return root
}

func findMinimumNode(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}

	nextMinNode := findMinimumNode(root.left)

	if nextMinNode != nil {
		return nextMinNode
	}

	return root
}

func findMaximumNode(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}

	nextMaxNode := findMaximumNode(root.right)

	if nextMaxNode != nil {
		return nextMaxNode
	}

	return root
}

func upperBound(root *TreeNode, value int) int {
	if root == nil {
		return math.MaxInt32
	}

	if root.value > value {
		return int(math.Min(float64(root.value), float64(upperBound(root.left, value))))
	} else {
		return upperBound(root.right, value)
	}
}

func (c Tree) GetAllValues() []int {
	values := make([]int, 0)
	stack := make([]*TreeNode, 0)

	if c.root != nil {
		stack = append(stack, c.root)
	}

	currIter, maxIter := 0, 1000

	for len(stack) > 0 {
		currentNode := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		values = append(values, currentNode.value)

		if currentNode.left != nil {
			stack = append(stack, currentNode.left)
		}

		if currentNode.right != nil {
			stack = append(stack, currentNode.right)
		}
		if currIter > maxIter {
			fmt.Println("infinite loop for root value of", currentNode.value)
			panic("stopping loop")
		}

		currIter += 1
	}

	return values
}

func (c *Tree) Add(value int) {
	c.root = add(c.root, value)
}

func (c *Tree) Remove(value int) *TreeNode {
	c.root = remove(c.root, value)
	return c.root
}

func (c *Tree) Find(value int) *TreeNode {
	return find(c.root, value)
}

func (c *Tree) UpperBound(value int) *TreeNode {
	upperBoundValue := upperBound(c.root, value)
	return c.Find(upperBoundValue)
}

func (ct *CoordinateTree) Count(y int) int {
	xCount, ok := ct.yAxisCount[y]
	if !ok {
		return 0
	}

	return xCount
}

func (ct *CoordinateTree) TotalCount() int {
	total := 0
	for _, yTree := range ct.yTree {
		total += yTree.root.count
	}

	return total
}

func (ct *CoordinateTree) LocationExist(x, y int) bool {
	yTree, ok := ct.yTree[x]

	if !ok {
		return false
	}

	yCoordinate := yTree.Find(y)

	if yCoordinate != nil {
		return true
	}

	return false
}

func New() CoordinateTree {
	return CoordinateTree{make(map[int]*Tree), make(map[int]int)}
}

func (ct *CoordinateTree) Add(x, y int) {

	_, ok := ct.yTree[x]

	if !ok {
		ct.yTree[x] = &Tree{}
	}

	_, ok = ct.yAxisCount[y]

	if !ok {
		ct.yAxisCount[y] = 0
	}

	ct.yAxisCount[y] += 1

	ct.yTree[x].Add(y)
}

func (ct *CoordinateTree) Remove(x, y int) error {
	yTree, ok := ct.yTree[x]

	if !ok {
		return errors.New(fmt.Sprintf("Could not remove coordinate with axis x:%d y:%d, no x coordinate found", x, y))
	}

	yRoot := yTree.Remove(y)

	if yRoot == nil {
		delete(ct.yTree, x)
	}

	return nil
}

func (ct *CoordinateTree) RemoveAll(y int) error {
	_, ok := ct.yAxisCount[y]

	if !ok {
		return errors.New(fmt.Sprintf("Cannot find coordinate with value y: %d", y))
	}

	delete(ct.yAxisCount, y)
	return nil
}

func (ct *CoordinateTree) UpperBound(x, y int) (int, int, error) {
	yTree, ok := ct.yTree[x]

	if !ok {
		return -1, -1, errors.New(fmt.Sprintf("Coordinate X : %d doesn't exist", x))
	}

	yUpperBound := yTree.UpperBound(y)

	if yUpperBound == nil {
		return -1, -1, errors.New(fmt.Sprintf("Cannot find upper bound for y value of %d", y))
	}

	return x, yUpperBound.value, nil
}

func (ct CoordinateTree) GetAllCoordinate() [][2]float32 {
	coordinates := make([][2]float32, 0)

	for x, yTree := range ct.yTree {
		yLocations := yTree.GetAllValues()

		for _, y := range yLocations {
			coordinates = append(coordinates, [2]float32{float32(x), float32(y)})
		}
	}

	return coordinates
}
