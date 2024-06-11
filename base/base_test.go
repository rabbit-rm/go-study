package main

import (
	"testing"
)

func TestRange(t *testing.T) {
	var nums = [3]int{1, 2, 3}
	// range 会发生值拷贝，所以修改 v 不会影响原数组
	// nums 也会发生拷贝，v 的值来自拷贝之后的 nums
	// ha := a, 如果是切片，底层数据是共享的，如果是数组，是值拷贝
	for i, v := range nums {
		if i == 0 {
			nums[0], nums[1] = 100, 200
			t.Log(nums)
		}
		nums[i] = v + 100
	}
	t.Log(nums)
}
