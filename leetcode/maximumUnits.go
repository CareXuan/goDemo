package leetcode

import (
	"sort"
)

func MaximumUnits(boxTypes [][]int, truckSize int) int {
	dict := map[int]int{}
	k := 0
	result := 0
	var keys []int
	for i := 0; i < len(boxTypes); i++ {
		if _, ok := dict[boxTypes[i][1]]; ok {
			dict[boxTypes[i][1]] += boxTypes[i][0]
		} else {
			dict[boxTypes[i][1]] = boxTypes[i][0]
		}
	}
	for j := range dict {
		keys = append(keys, j)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	for m := range keys {
		if k >= truckSize {
			break
		}
		if k+dict[keys[m]] < truckSize {
			result += keys[m] * dict[keys[m]]
			k += dict[keys[m]]
		} else {
			result += keys[m] * (truckSize - k)
			k = truckSize
		}
	}
	return result
}
