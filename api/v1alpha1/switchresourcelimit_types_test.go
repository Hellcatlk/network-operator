package v1alpha1

import "testing"

func TestSwitchResourceLimitSpecVerify(t *testing.T) {
	untaggedVLAN := 20
	cases := []struct {
		name          string
		spec          *SwitchResourceLimitSpec
		configuration *SwitchPortConfiguration
		expectedError bool
	}{
		{
			name: "no limit",
			spec: nil,
			configuration: &SwitchPortConfiguration{
				Spec: SwitchPortConfigurationSpec{
					VLANs:        "1-10",
					UntaggedVLAN: &untaggedVLAN,
				},
			},
			expectedError: false,
		},
		{
			name: "no limit",
			spec: &SwitchResourceLimitSpec{},
			configuration: &SwitchPortConfiguration{
				Spec: SwitchPortConfigurationSpec{
					VLANs:        "1-10",
					UntaggedVLAN: &untaggedVLAN,
				},
			},
			expectedError: false,
		},
		{
			name: "vlan range limit: 1-20",
			spec: &SwitchResourceLimitSpec{
				VLANRange: "1-20",
			},
			configuration: &SwitchPortConfiguration{
				Spec: SwitchPortConfigurationSpec{
					VLANs:        "1-10",
					UntaggedVLAN: &untaggedVLAN,
				},
			},
			expectedError: false,
		},
		{
			name: "vlan range limit: 1-10",
			spec: &SwitchResourceLimitSpec{
				VLANRange: "1-10",
			},
			configuration: &SwitchPortConfiguration{
				Spec: SwitchPortConfigurationSpec{
					VLANs:        "1-10",
					UntaggedVLAN: &untaggedVLAN,
				},
			},
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.spec.Verify(c.configuration)
			if (err != nil) != c.expectedError {
				t.Errorf("Got unexpected error: %v", err)
			}
		})
	}
}
