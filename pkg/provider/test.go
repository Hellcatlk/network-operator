package provider

import corev1 "k8s.io/api/core/v1"

// TestSwitch just use for test
type TestSwitch struct {
}

// GetOS  return switch's os
func (s *TestSwitch) GetOS() string {
	return "test"
}

// GetProtocol return switch's protocol
func (s *TestSwitch) GetProtocol() string {
	return "test"
}

// GetHost return switch's host
func (s *TestSwitch) GetHost() string {
	return ""
}

// GetSecret return switch's certificate secret reference
func (s *TestSwitch) GetSecret() *corev1.SecretReference {
	return &corev1.SecretReference{}
}

// GetOptions return switch's options
func (s *TestSwitch) GetOptions() map[string]string {
	return nil
}
