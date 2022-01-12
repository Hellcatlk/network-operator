package v1alpha1

import "testing"

func TestPortVerify(t *testing.T) {
	untaggedVLAN := 20
	cases := []struct {
		name          string
		port          *Port
		configuration *SwitchPortConfiguration
		expectedError bool
	}{
		{
			name: "disabled",
			port: &Port{
				PhysicalPortName: "test",
				Disabled:         true,
			},
			configuration: &SwitchPortConfiguration{
				Spec: SwitchPortConfigurationSpec{
					TaggedVLANRange: "1-10",
					UntaggedVLAN:    &untaggedVLAN,
				},
			},
			expectedError: true,
		},
		{
			name: "vlan is out of range",
			port: &Port{
				PhysicalPortName: "test",
				VLANRange:        "1-5",
			},
			configuration: &SwitchPortConfiguration{
				Spec: SwitchPortConfigurationSpec{
					TaggedVLANRange: "1-10",
					UntaggedVLAN:    &untaggedVLAN,
				},
			},
			expectedError: true,
		},
		{
			name: "trunk disabled",
			port: &Port{
				PhysicalPortName: "test",
				TrunkDisabled:    true,
			},
			configuration: &SwitchPortConfiguration{
				Spec: SwitchPortConfigurationSpec{
					TaggedVLANRange: "1-10",
					UntaggedVLAN:    &untaggedVLAN,
				},
			},
			expectedError: true,
		},
		{
			name: "vlan in range",
			port: &Port{
				PhysicalPortName: "test",
				VLANRange:        "1-20",
			},
			configuration: &SwitchPortConfiguration{
				Spec: SwitchPortConfigurationSpec{
					TaggedVLANRange: "1-10",
					UntaggedVLAN:    &untaggedVLAN,
				},
			},
			expectedError: false,
		},
		{
			name: "vlan in range",
			port: &Port{
				PhysicalPortName: "test",
				VLANRange:        "1-20",
			},
			configuration: &SwitchPortConfiguration{
				Spec: SwitchPortConfigurationSpec{
					UntaggedVLAN: &untaggedVLAN,
				},
			},
			expectedError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.port.Verify(c.configuration)
			if (err != nil) != c.expectedError {
				t.Errorf("Got unexpected error: %v", err)
			}
		})
	}
}
