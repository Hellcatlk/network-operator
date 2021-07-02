package strings

import (
	"reflect"
	"testing"
)

func TestToSlice(t *testing.T) {
	cases := []struct {
		name          string
		expected      []int
		expectedError bool
	}{
		{
			name:     "1-5,7",
			expected: []int{1, 2, 3, 4, 5, 7},
		},
		{
			name:          "1-5,,7",
			expectedError: true,
		},
		{
			name:          "1--5,7",
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ToSlice(c.name)
			if !reflect.DeepEqual(c.expected, got) {
				t.Errorf("expected: %v, got: %v", c.expected, got)
			}
			if (err != nil) != c.expectedError {
				t.Errorf("got unexpected error: %v", err)
			}
		})
	}
}
