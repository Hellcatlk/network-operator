package strings

import (
	"fmt"
	"strconv"
	"strings"
)

// ToSlice transform "1-5,7" to [1,2,3,4,5,7]
func ToSlice(formatStr string) ([]int, error) {
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
