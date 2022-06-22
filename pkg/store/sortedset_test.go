package store

import (
	"fmt"
	"testing"
)

func TestSortedSet(t *testing.T) {
	sortedset := New()

	sortedset.AddOrUpdate("a", 89, "Kelly")
	sortedset.AddOrUpdate("b", 100, "Staley")
	sortedset.AddOrUpdate("c", 100, "Jordon")
	sortedset.AddOrUpdate("d", -321, "Park")
	sortedset.AddOrUpdate("e", 101, "Albert")
	sortedset.AddOrUpdate("f", 99, "Lyman")
	sortedset.AddOrUpdate("g", 99, "Singleton")
	sortedset.AddOrUpdate("h", 70, "Audrey")

	sortedset.AddOrUpdate("e", 99, "ntr")

	sortedset.IterFuncByRankRange(0, -1, func(key string, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})

	nodes := sortedset.GetByScoreRange(100, 90, &GetByScoreRangeOptions{Limit: 5})
	for _, node := range nodes {
		fmt.Println("node:", node.Key(), node.Value, node.Score())
	}

}
