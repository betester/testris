package treecoordinate

import (
	"fmt"
	"testing"
)

func TestFindAndAdd(t *testing.T) {
	tree := Tree{}

	tree.Add(2)
	tree.Add(1)
	tree.Add(3)

	node1 := tree.Find(1)

	if node1 == nil {
		t.Error("Cannot find node with value 1")
		t.Fail()
	}

	if tree.root == nil || tree.root.value != 2 {
		t.Error("Wrong insert on tree, 2 should be the root")
		t.Fail()
	}

	if tree.root.left.value != 1 {
		t.Error("Wrong insert on tree, 1 should be on the left of the root")
		t.Fail()
	}

	if tree.root.right.value != 3 {
		t.Error("Wrong insert on tree, 3 should be on the right of the root")
		t.Fail()
	}
}

func TestRemove(t *testing.T) {

	tree := Tree{}

	tree.Add(7)
	tree.Add(1)
	tree.Add(12)
	tree.Add(11)
	tree.Add(15)

	tree.Remove(7)

	node1 := tree.Find(7)

	if node1 != nil {
		t.Error("Wrong on remove, node with value 7 still exist")
		t.Fail()
	}

	tree.Remove(11)

	node2 := tree.Find(11)

	if node2 != nil {
		t.Error("Wrong on remove, node with value 11 still exist")
		t.Fail()
	}

	if tree.root.value != 12 {
		t.Error("Wrong on remove, 12 is not the root")
		t.Fail()
	}
}

func TestRemoveLeftSide(t *testing.T) {

	tree := Tree{}

	tree.Add(3)
	tree.Add(2)
	tree.Add(1)

	tree.Remove(3)

	node1 := tree.Find(3)

	if node1 != nil {
		t.Error("Wrong on remove, node with value 3 still exist")
		t.Fail()
	}

	if tree.root.value != 2 {
		t.Error("Wrong on remove, 2 is not the root")
		t.Fail()
	}
}

func TestUpperBound(t *testing.T) {

	tree := Tree{}

	tree.Add(10)
	tree.Add(7)
	tree.Add(9)
	tree.Add(6)
	tree.Add(4)
	tree.Add(15)

	result := tree.UpperBound(6)

	if result.value != 7 {
		t.Error(fmt.Sprintf("The upper bound value should be 6, found %d instead", result.value))
		t.Fail()
	}
}

func TestAddCoordinateTree(t *testing.T) {
	cTree := New()

	cTree.Add(1, 2)
	result := cTree.LocationExist(1, 2)

	if !result {
		t.Error("Position 1 and 2 could not be found on the coordinate tree")
	}
}

func TestRemoveCoordinateTree(t *testing.T) {
	cTree := New()

	cTree.Add(1, 2)
	cTree.Remove(1, 2)
	result := cTree.LocationExist(1, 2)

	if result {
		t.Error("Position 1 and 2 should not be found on the coordinate tree")
	}
}

func TestUpperBoundCoordinateTree(t *testing.T) {
	cTree := New()

	cTree.Add(1, 2)
	cTree.Add(1, 3)
	cTree.Add(1, 4)

	ux, uy, err := cTree.UpperBound(1, 3)

	if err != nil {
		t.Error(err)
	}

	if ux != 1 && uy != 4 {
		t.Error(fmt.Sprintf("Wrong upper bound, should return 1 and 4 instead fodund %d and %d", ux, uy))
	}
}

func TestCountForCoordinateTree(t *testing.T) {
	cTree := New()

	cTree.Add(1, 2)
	cTree.Add(1, 3)
	cTree.Add(1, 4)

	count := cTree.TotalCount()

	if count != 3 {
		t.Error(fmt.Sprintf("Total node count should be 3 found %d", count))
	}
}

func TestCountAddingElement(t *testing.T) {
	tree := Tree{}

	tree.Add(1)
	tree.Add(2)
	tree.Add(3)

	if tree.root.count != 3 {
		t.Error(fmt.Sprintf("Total node count should be 3 found %d", tree.root.count))
	}

	tree.Remove(2)

	if tree.root.count != 2 {
		t.Error(fmt.Sprintf("Total node count should be 2 found %d", tree.root.count))
	}
}
