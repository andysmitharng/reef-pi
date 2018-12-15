package mockpca9685

import (
	"fmt"
	"log"
	"sort"

	"github.com/reef-pi/rpi/i2c"

	"github.com/reef-pi/reef-pi/controller/settings"
	"github.com/reef-pi/types/driver"
)

type mockPwmChannel struct {
	name string
}

func (m *mockPwmChannel) Name() string    { return m.name }
func (m *mockPwmChannel) Close() error    { return nil }
func (m *mockPwmChannel) LastState() bool { return false }

func (m *mockPwmChannel) Set(value float64) error {
	log.Printf("mockpca9685: setting pwm %s to %f", m.name, value)
	return nil
}

func (m *mockPwmChannel) Write(on bool) error {
	if on {
		return m.Set(100)
	}
	return m.Set(0)
}

type mockDriver struct {
	channels map[string]*mockPwmChannel
}

func (m *mockDriver) Close() error { return nil }

func (m *mockDriver) Metadata() driver.Metadata {
	return driver.Metadata{
		Name:        "pca9685",
		Description: "Mock driver - no actual hardware",
		Capabilities: driver.Capabilities{
			PWM:    true,
			Output: true,
		},
	}
}

func (m *mockDriver) PWMChannels() []driver.PWMChannel {
	var chs []driver.PWMChannel
	for _, ch := range m.channels {
		chs = append(chs, ch)
	}
	sort.Slice(chs, func(i, j int) bool { return chs[i].Name() < chs[j].Name() })
	return chs
}

func (m *mockDriver) GetPWMChannel(name string) (driver.PWMChannel, error) {
	ch, ok := m.channels[name]
	if !ok {
		return nil, fmt.Errorf("unknown pwm channel %s", name)
	}
	return ch, nil
}

func (m *mockDriver) OutputPins() []driver.OutputPin {
	var chs []driver.OutputPin
	for _, ch := range m.channels {
		chs = append(chs, ch)
	}
	sort.Slice(chs, func(i, j int) bool { return chs[i].Name() < chs[j].Name() })
	return chs
}

func (m *mockDriver) GetOutputPin(name string) (driver.OutputPin, error) {
	pin, ok := m.channels[name]
	if !ok {
		return nil, fmt.Errorf("unknown pin with name %s", name)
	}
	return pin, nil
}

func NewMockDriver(s settings.Settings, b i2c.Bus) (driver.Driver, error) {
	d := &mockDriver{
		channels: make(map[string]*mockPwmChannel),
	}
	for i := 0; i < 8; i++ {
		d.channels[fmt.Sprintf("%d", i)] = &mockPwmChannel{name: fmt.Sprintf("%d", i)}
	}
	return d, nil
}