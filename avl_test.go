// avl_test.go - AVL tree tests.
//
// To the extent possible under law, Yawning Angel has waived all copyright
// and related or neighboring rights to avl, using the Creative
// Commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package avl

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAVLTree(t *testing.T) {
	require := require.New(t)

	cmpInt := func(a, b interface{}) int {
		aInt, bInt := a.(int), b.(int)
		switch {
		case aInt < bInt:
			return -1
		case aInt > bInt:
			return 1
		default:
			return 0
		}
	}

	tree := New(cmpInt)
	require.Equal(0, tree.Len(), "Len(): empty")
	require.Nil(tree.First(), "First(): empty")
	require.Nil(tree.Last(), "Last(): empty")

	iter := tree.Iterator(Forward)
	require.Nil(iter.First(), "Iterator: First(), empty")
	require.Nil(iter.Next(), "Iterator: Next(), empty")

	// Test insertion.
	const nrEntries = 1024
	insertedMap := make(map[int]*Node)
	for len(insertedMap) != nrEntries {
		v := rand.Int()
		if insertedMap[v] != nil {
			continue
		}
		insertedMap[v] = tree.Insert(v)
	}
	require.Equal(nrEntries, tree.Len(), "Len(): After insertion")

	// Ensure that all entries can be found.
	for k, v := range insertedMap {
		require.Equal(v, tree.Find(k), "Find(): %v", k)
		require.Equal(k, v.Value, "Find(): %v Value", k)
	}

	// Test the forward/backward iterators.
	inOrder := make([]int, 0, nrEntries)
	for k := range insertedMap {
		inOrder = append(inOrder, k)
	}
	sort.Ints(inOrder)

	iter = tree.Iterator(Forward)
	visited := 0
	for node := iter.First(); node != nil; node = iter.Next() {
		v, idx := node.Value.(int), visited
		require.Equal(inOrder[visited], v, "Iterator: Forward[%v]", idx)
		require.Equal(node, iter.Get(), "Iterator: Forward[%v]: Get()", idx)
		visited++
	}
	require.Equal(nrEntries, visited, "Iterator: Forward: Visited")

	iter = tree.Iterator(Backward)
	visited = 0
	for node := iter.First(); node != nil; node = iter.Next() {
		v, idx := node.Value.(int), nrEntries-1-visited
		require.Equal(inOrder[idx], v, "Iterator: Backward[%v]", idx)
		require.Equal(node, iter.Get(), "Iterator: Backward[%v]: Get()", idx)
		visited++
	}
	require.Equal(nrEntries, visited, "Iterator: Backward: Visited")

	// Test removal.
	for i, idx := range rand.Perm(nrEntries) {
		v := inOrder[idx]
		n := tree.Find(v)
		require.Equal(v, n.Value, "Find(): %v (Pre-remove)", v)
		tree.Remove(n)
		require.Equal(nrEntries-(i+1), tree.Len(), "Len(): %v (Post-remove)", v)

		n = tree.Find(v)
		require.Nil(n, "Find(): %v (Post-remove)", v)
	}
	require.Equal(0, tree.Len(), "Len(): After removal")
}
