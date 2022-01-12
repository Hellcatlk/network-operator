package strings

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// RangeToSlice transform "1-5,7" to [1,2,3,4,5,7]
func RangeToSlice(formatStr string) ([]int, error) {
	if formatStr == "" {
		return nil, nil
	}
	nums := []int{}
	for _, str := range strings.Split(formatStr, ",") {
		r := strings.Split(str, "-")
		switch len(r) {
		case 1:
			num, err := strconv.Atoi(r[0])
			if err != nil {
				return nil, err
			}
			nums = append(nums, num)
		case 2:
			begin, err := strconv.Atoi(r[0])
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(r[1])
			if err != nil {
				return nil, err
			}
			for ; begin <= end; begin = begin + 1 {
				nums = append(nums, begin)
			}
		default:
			return nil, fmt.Errorf("invalid format")
		}
	}

	return nums, nil
}

// SliceToRange transform [1,2,3,4,5,7] to "1-5,7"
func SliceToRange(nums []int) string {
	if len(nums) == 0 {
		return ""
	}
	if len(nums) == 1 {
		return strconv.Itoa(nums[0])
	}
	sort.Ints(nums)
	length := len(nums)
	j := 0
	for i := 1; i < length; i++ {
		if nums[i] != nums[j] {
			j++
			if j < i {
				nums[i], nums[j] = nums[j], nums[i]
			}
		}
	}
	nums = nums[:j+1]
	length = len(nums)

	formatStr := ""
	var begin, end int
	begin = nums[0]
	end = nums[0]
	for i := 1; i < length; i++ {
		if end == nums[i]-1 {
			end = nums[i]
			if i == (length - 1) {
				formatStr = formatStr + strconv.Itoa(begin) + "-" + strconv.Itoa(end)
			}
		} else {
			if begin == end {
				formatStr = formatStr + strconv.Itoa(begin) + ","
			} else {
				formatStr = formatStr + strconv.Itoa(begin) + "-" + strconv.Itoa(end) + ","
			}
			begin = nums[i]
			end = nums[i]
			if i == (length - 1) {
				formatStr = formatStr + strconv.Itoa(end)
			}
		}
	}

	return formatStr
}

// RangeContains if allowed contains target return nil, else return err.
func RangeContains(allowed string, target string) error {
	// Get allowed vlan range
	allowedSlice, err := RangeToSlice(allowed)
	if err != nil {
		return err
	}
	allowedMap := make(map[int]struct{})
	for _, value := range allowedSlice {
		allowedMap[value] = struct{}{}
	}

	// Get target vlan range
	targetSlice, err := RangeToSlice(target)
	if err != nil {
		return err
	}

	// Check vlan range
	for _, value := range targetSlice {
		_, existed := allowedMap[value]
		if !existed {
			return fmt.Errorf("%d is out of allowed range: %s", value, allowed)
		}
	}

	return nil
}

// LastJSON find last json from string
func LastJSON(str string) ([]byte, error) {
	lastJSON := ""
	count := 0
	for _, char := range str {
		switch char {
		case '{':
			count = count + 1
			if count == 1 {
				lastJSON = ""
			}
		case '}':
			count = count - 1
			if count == 0 {
				lastJSON = lastJSON + "}"
			}
		}
		if count > 0 {
			lastJSON = lastJSON + string(char)
		} else if count < 0 {
			return nil, fmt.Errorf("invalid string: %s", str)
		}
	}
	if count != 0 || len(lastJSON) == 0 {
		return nil, fmt.Errorf("invalid string: %s", str)
	}

	return []byte(lastJSON), nil
}

// SliceContains check the slice contains str or not, if exist return true, else return false.
func SliceContains(slice []string, str string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	return false
}

// SliceDelete str from slice, if success return true, else return false.
func SliceDelete(slice *[]string, str string) bool {
	for i := 0; i < len(*slice); i++ {
		if (*slice)[i] == str {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			return true
		}
	}
	return false
}

// Expansion combine the contents of A and B
func Expansion(a, b string) (string, error) {
	if a == "" {
		return b, nil
	}
	if b == "" {
		return a, nil
	}

	arr1, err := RangeToSlice(a)
	if err != nil {
		return "", err
	}

	arr2, err := RangeToSlice(b)
	if err != nil {
		return "", err
	}

	arr1 = append(arr1, arr2...)

	return SliceToRange(arr1), nil
}

// Shrink remove B from the content of A
func Shrink(a, b string) (string, error) {
	var k int

	if a == "" {
		return "", nil
	}
	if b == "" {
		return a, nil
	}

	arr1, err := RangeToSlice(a)
	if err != nil {
		return "", err
	}

	arr2, err := RangeToSlice(b)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(arr1); i++ {
		k = arr1[i]
		for _, v := range arr2 {
			if k == v {
				arr1 = append(arr1[:i], arr1[i+1:]...)
				i--
			}
		}
	}

	return SliceToRange(arr1), nil
}
