package leetcode

func CountPairs(deliciousness []int) int {
	count := 0
	for i := 0; i < len(deliciousness); i++ {
		for j := i + 1; j < len(deliciousness); j++ {
			var tmp int
			tmp = deliciousness[i] + deliciousness[j]
			if tmp&(tmp-1) == 0 && tmp != 0{
				count++
			}
		}
	}
	return count
}
