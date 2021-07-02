package strings

import (
	"fmt"
	"strconv"
	"strings"
)

// RangeToSlice transform "1-5,7" to [1,2,3,4,5,7]
func RangeToSlice(formatStr string) ([]int, error) {
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
