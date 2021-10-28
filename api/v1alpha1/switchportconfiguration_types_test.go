package v1alpha1

import "testing"

func TestIsEqual(t *testing.T) {
	cases := []struct {
		target   *SwitchPortConfigurationSpec
		actual   *SwitchPortConfigurationSpec
		expected bool
	}{
		{
			target:   nil,
			actual:   nil,
			expected: true,
		},
		{
			target:   &SwitchPortConfigurationSpec{},
			actual:   nil,
			expected: true,
		},
		{
			target:   nil,
			actual:   &SwitchPortConfigurationSpec{},
			expected: true,
		},
		{
			target:   &SwitchPortConfigurationSpec{},
			actual:   &SwitchPortConfigurationSpec{},
			expected: true,
		},
		{
			target: &SwitchPortConfigurationSpec{
				TaggedVLANRange: "1-10,11",
			},
			actual: &SwitchPortConfigurationSpec{
				TaggedVLANRange: "1-12",
			},
			expected: false,
		},
		{
			target: &SwitchPortConfigurationSpec{
				TaggedVLANRange: "1-10,11,12",
			},
			actual: &SwitchPortConfigurationSpec{
				TaggedVLANRange: "1-12",
			},
			expected: true,
		},
		{
			target: &SwitchPortConfigurationSpec{
				UntaggedVLAN:    new(int),
				TaggedVLANRange: "1-10,11,12",
			},
			actual: &SwitchPortConfigurationSpec{
				UntaggedVLAN:    new(int),
				TaggedVLANRange: "1-12",
			},
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			got := c.target.IsEqual(c.actual)
			if c.expected != got {
				t.Errorf("Expected: %v, got: %v", c.expected, got)
			}
		})
	}
}
