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
				Disabled: true,
			},
			configuration: &SwitchPortConfiguration{
				Spec: SwitchPortConfigurationSpec{
					VLANs:        "1-10",
					UntaggedVLAN: &untaggedVLAN,
				},
			},
			expectedError: true,
		},
		{
			name: "trunk disabled",
			port: &Port{
				TrunkDisabled: true,
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
			err := c.port.Verify(c.configuration)
			if (err != nil) != c.expectedError {
				t.Errorf("Got unexpected error: %v", err)
			}
		})
	}
}
