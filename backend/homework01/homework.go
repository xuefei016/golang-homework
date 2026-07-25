package homework01

import "sort"

// 1. 只出现一次的数字
// 给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
func SingleNumber(nums []int) int {
	// 方法1
	// counts := make(map[int]int)
	// for _, num := range nums {
	// 	counts[num]++
	// }

	// for num, count := range counts {
	// 	if count == 1 {
	// 		return num
	// 	}
	// }
	// return 0

	// 方法2
	result := 0
	for _, num := range nums {
		result ^= num
	}
	return result
}

// 2. 回文数
// 判断一个整数是否是回文数
func IsPalindrome(x int) bool {
	if x < 0 {
		return false
	}

	// 方法1
	// str := strconv.Itoa(x)
	// i, j := 0, len(str)-1
	// for(i < j) {
	// 	if(str[i] != str[j]) {
	// 		return false
	// 	}
	// 	i++
	// 	j--
	// }
	// return true

	// 方法2
	if x != 0 && x%10 == 0 {
		return false
	}
	reverted := 0
	for x > reverted {
		reverted = reverted*10 + x%10
		x /= 10
	}
	return x == reverted || x == reverted/10
}

// 3. 有效的括号
// 给定一个只包括 '(', ')', '{', '}', '[', ']' 的字符串，判断字符串是否有效
func IsValid(s string) bool {
	pairs := map[byte]byte{
		')': '(',
		'}': '{',
		']': '[',
	}

	stack := []byte{}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if pair, ok := pairs[c]; ok {
			if len(stack) == 0 || pair != stack[len(stack)-1] {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, c)
		}
	}
	return len(stack) == 0
}

// 4. 最长公共前缀
// 查找字符串数组中的最长公共前缀
func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	index := 0
	pause := false
	for index < len(strs[0]) {
		for i := 1; i < len(strs); i++ {
			if index >= len(strs[i]) || strs[0][index] != strs[i][index] {
				pause = true
				break
			}
			if i == len(strs)-1 {
				index++
			}
		}
		if pause {
			break
		}
	}

	if index == 0 {
		return ""
	}
	return strs[0][:index]
}

// 5. 加一
// 给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func PlusOne(digits []int) []int {
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}
	digits = make([]int, len(digits)+1)
	digits[0] = 1
	return digits
}

// 6. 删除有序数组中的重复项
// 给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
// 不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。
func RemoveDuplicates(nums []int) int {
	slow := 0
	fast := 1
	if len(nums) == 0 {
		return 0
	}
	for fast <= len(nums)-1 {
		if nums[slow] == nums[fast] {
			fast++
		} else {
			slow++
			nums[slow] = nums[fast]
			fast++
		}
	}
	return slow + 1
}

// 7. 合并区间
// 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。
func Merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return nil
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	result := [][]int{intervals[0]}
	last_right := intervals[0][1]
	for i := 1; i < len(intervals); i++ {
		left := intervals[i][0]
		right := intervals[i][1]

		if last_right >= left {
			if right > last_right {
				result[len(result)-1][1] = right
				last_right = right
			}
		} else {
			result = append(result, intervals[i])
			last_right = right
		}
	}
	return result

}

// 8. 两数之和
// 给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
func TwoSum(nums []int, target int) []int {
	seen := make(map[int]int)
	for i, num := range nums {
		if seenIndex, ok := seen[target-num]; ok {
			return []int{seenIndex, i}
		}
		seen[num] = i
	}
	return nil
}
