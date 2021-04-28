package test

import (
	"context"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
	"github.com/metal3-io/networkconfiguration-operator/pkg/device"
)

func NewTT(ctx context.Context, address string) (sw device.Switch, err error) {
	return &TestTest{}, nil
}

// TestTest just for test
type TestTest struct {
}

func (tt *TestTest) PowerOn(ctx context.Context) (err error) {
	return
}

func (tt *TestTest) PowerOff(ctx context.Context) (err error) {
	return
}

func (tt *TestTest) CreateVlan(ctx context.Context, vlans []v1alpha1.VLAN) (err error) {
	return
}

func (tt *TestTest) DeleteVlan(ctx context.Context, vlans []v1alpha1.VLAN) (err error) {
	return
}

func (tt *TestTest) GetPortAttr(ctx context.Context, portID string) (configuration *v1alpha1.SwitchPortConfiguration, err error) {
	return &v1alpha1.SwitchPortConfiguration{}, nil
}

func (tt *TestTest) SetPortAttr(ctx context.Context, portID string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	return
}

func (tt *TestTest) ResetPort(ctx context.Context, portID string) (err error) {
	return
}
