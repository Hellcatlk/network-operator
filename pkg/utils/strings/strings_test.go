package strings

import (
	"reflect"
	"testing"
)

func TestRangeToSlice(t *testing.T) {
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
		{
			name:          ",7",
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := RangeToSlice(c.name)
			if !reflect.DeepEqual(c.expected, got) {
				t.Errorf("expected: %v, got: %v", c.expected, got)
			}
			if (err != nil) != c.expectedError {
				t.Errorf("got unexpected error: %v", err)
			}
		})
	}
}

func TestLastJSON(t *testing.T) {
	cases := []struct {
		data          string
		expected      []byte
		expectedError bool
	}{
		{
			data:          "invalid",
			expectedError: true,
		},
		{
			data:     "test:{a:123}",
			expected: []byte("{a:123}"),
		},
		{
			data:     "test:{a:123}, {b:123}",
			expected: []byte("{b:123}"),
		},
		{
			data:          "test:{a:123}}",
			expectedError: true,
		},
		{
			data:          "test:{{a:123}",
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.data, func(t *testing.T) {
			got, err := LastJSON(c.data)
			if !reflect.DeepEqual(c.expected, got) {
				t.Errorf("expected: %v, got: %v", c.expected, got)
			}
			if (err != nil) != c.expectedError {
				t.Errorf("got unexpected error: %v", err)
			}
		})
	}
}

func TestContains(t *testing.T) {
	cases := []struct {
		slice    []string
		str      string
		expected bool
	}{
		{
			slice:    []string{"test1"},
			str:      "",
			expected: false,
		},
		{
			slice:    []string{""},
			str:      "test1",
			expected: false,
		},
		{
			slice:    []string{"test1"},
			str:      "test1",
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(t.Name(), func(t *testing.T) {
			got := SliceContains(c.slice, c.str)
			if c.expected != got {
				t.Fatalf("expected: %v, got: %v", c.expected, got)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		slice         []string
		str           string
		expectedSlice []string
		expected      bool
	}{
		{
			slice:         []string{"test1"},
			str:           "",
			expectedSlice: []string{"test1"},
			expected:      false,
		},
		{
			slice:         []string{""},
			str:           "test1",
			expectedSlice: []string{""},
			expected:      false,
		},
		{
			slice:         []string{"test1"},
			str:           "test1",
			expectedSlice: []string{},
			expected:      true,
		},
		{
			slice:         []string{"test1", "test2", "test3"},
			str:           "test1",
			expectedSlice: []string{"test2", "test3"},
			expected:      true,
		},
		{
			slice:         []string{"test1", "test2", "test3"},
			str:           "test2",
			expectedSlice: []string{"test1", "test3"},
			expected:      true,
		},
		{
			slice:         []string{"test1", "test2", "test3"},
			str:           "test3",
			expectedSlice: []string{"test1", "test2"},
			expected:      true,
		},
	}

	for _, c := range cases {
		t.Run(t.Name(), func(t *testing.T) {
			got := SliceDelete(&c.slice, c.str)
			if c.expected != got || !reflect.DeepEqual(c.expectedSlice, c.slice) {
				t.Errorf("expectedSlice: %v(%d), slice: %v(%d)", c.expectedSlice, len(c.expectedSlice), c.slice, len(c.slice))
				t.Errorf("expected: %v, got: %v", c.expected, got)
			}
		})
	}
}

func TestSliceToRange(t *testing.T) {
	cases := []struct {
		arr      []int
		expected string
	}{
		{
			arr:      []int{1, 2, 2, 3, 4, 5, 7},
			expected: "1-5,7",
		},
		{
			arr:      []int{},
			expected: "",
		},
		{
			arr:      []int{1, 2, 3, 4, 5, 6, 7},
			expected: "1-7",
		},
		{
			arr:      []int{1, 5, 7},
			expected: "1,5,7",
		},
	}

	for _, c := range cases {
		t.Run(t.Name(), func(t *testing.T) {
			got := SliceToRange(c.arr)
			if c.expected != got {
				t.Errorf("expected: %v, got: %v", c.expected, got)
			}
		})
	}
}

func TestExpansion(t *testing.T) {
	cases := []struct {
		str1          string
		str2          string
		expected      string
		expectedError bool
	}{
		{
			str1:     "1-5,7",
			str2:     "6,8",
			expected: "1-8",
		},
		{
			str1:     "1-5,7",
			str2:     "10-15",
			expected: "1-5,7,10-15",
		},
		{
			str1:          "1--5,7",
			str2:          "6,8",
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(t.Name(), func(t *testing.T) {
			got, err := Expansion(c.str1, c.str2)
			if c.expected != got {
				t.Errorf("expected: %v, got: %v", c.expected, got)
			}
			if (err != nil) != c.expectedError {
				t.Errorf("got unexpected error: %v", err)
			}
		})
	}
}

func TestShrink(t *testing.T) {
	cases := []struct {
		str1          string
		str2          string
		expected      string
		expectedError bool
	}{
		{
			str1:     "1-5,7",
			str2:     "",
			expected: "1-5,7",
		},
		{
			str1:     "1-5,7",
			str2:     "1-5,7",
			expected: "",
		},
		{
			str1:          "1--5,7",
			str2:          "6,8",
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(t.Name(), func(t *testing.T) {
			got, err := Shrink(c.str1, c.str2)
			if c.expected != got {
				t.Errorf("expected: %v, got: %v", c.expected, got)
			}
			if (err != nil) != c.expectedError {
				t.Errorf("got unexpected error: %v", err)
			}
		})
	}
}
