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
