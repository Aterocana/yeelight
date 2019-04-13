package yeelight

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// SupportedFeatures is the list of supported as bool values.
type SupportedFeatures struct {
	GetProp    bool `json:"get_prop"`
	SetDefault bool `json:"set_default"`
	SetPower   bool `json:"set_power"`
	Toggle     bool `json:"toggle"`
	SetBright  bool `json:"set_bright"`
	StartCF    bool `json:"start_cf"`
	StopCF     bool `json:"stop_cf"`
	SetScene   bool `json:"set_scene"`
	CronAdd    bool `json:"cron_add"`
	CronGet    bool `json:"cron_get"`
	CronDel    bool `json:"cron_del"`
	SetCtAbx   bool `json:"set_ct_abx"`
	SetRGB     bool `json:"set_rgb"`
}

// PowerValue is the YeeLight device's power value.
type PowerValue string

const (
	// On is when the bulb is powered.
	On PowerValue = "on"
	// Off is when the bulb is not powered.
	Off PowerValue = "off"
)

func (p PowerValue) isValid() bool {
	return p == "on" || p == "off"
}

// TurnOnValue is the optional SetPower parameter
type TurnOnValue int

const (
	// NormalMode is the default value: turn the
	// device ON with the previous mode.
	NormalMode TurnOnValue = iota

	// CTMode turns the device ON with ColorTemperature Mode.
	CTMode

	// RGBMode turns the device ON with RGB Mode.
	RGBMode

	// HSVMode turns the device ON with HSV Mode.
	HSVMode

	// ColorFlowMode turns the device ON with ColorFlow Mode.
	ColorFlowMode

	// NightLightMode turns the device ON with NightLight Mode.
	NightLightMode
)

func (t TurnOnValue) isValid() bool {
	return t > 0 && t < 6
}

// ColorModeValue is the YeeLight device's color mode.
type ColorModeValue int

const (
	_ = iota
	// ColorMode is when device is in RGB mode.
	ColorMode //1

	// ColorTemperature is when device is in CT .
	ColorTemperature //2

	// HSV is when device is in Hue Saturation Value mode.
	HSV //3
)

func (c ColorModeValue) String() string {
	switch c {
	case ColorMode:
		return "rgb"
	case ColorTemperature:
		return "temperature"
	case HSV:
		return "hsv"
	default:
		return "unknown mode"
	}
}

// MarshalJSON convert a ColorModeValue in json value.
func (c ColorModeValue) MarshalJSON() ([]byte, error) {
	str := c.String()
	return json.Marshal(str)
}

// RGBValue is the 3 byte RGB representation.
type RGBValue struct {
	red   uint8
	green uint8
	blue  uint8
}

// NewRGB instantiate a RGBValue from its int value.
func NewRGB(val int) (RGBValue, error) {
	if val < 0 || val > 0xffffff {
		return RGBValue{}, errors.Wrapf(ErrInvalidRange, "invalid RGB value: %d", val)
	}
	return RGBValue{
		red:   uint8(0x00FF0000 & val >> 16),
		green: uint8(0x0000FF00 & val >> 8),
		blue:  uint8(0x000000FF & val),
	}, nil
}

// Get returns rgb as int value
func (rgb RGBValue) Get() int {
	return int(rgb.red)<<16 + int(rgb.green)<<8 + int(rgb.blue)
}

// MarshalJSON serialize RGBValue struct in json format.
func (rgb RGBValue) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(`{"r":%d,"g":%d,"b":%d}`, rgb.red, rgb.green, rgb.blue)
	return []byte(s), nil
}

// setSupport creates the support map[string]bool from a
// parsed string.
func (y *YeeLight) setSupport(support string) {
	features := strings.Split(support, " ")
	for _, feature := range features {
		switch strings.ToLower(feature) {
		case "get_prop":
			y.Support.GetProp = true
		case "set_default":
			y.Support.SetDefault = true
		case "set_power":
			y.Support.SetPower = true
		case "toggle":
			y.Support.Toggle = true
		case "set_bright":
			y.Support.SetBright = true
		case "start_cf":
			y.Support.StartCF = true
		case "stop_cf":
			y.Support.StopCF = true
		case "set_scene":
			y.Support.SetScene = true
		case "cron_add":
			y.Support.CronAdd = true
		case "cron_get":
			y.Support.CronGet = true
		case "cron_del":
			y.Support.CronDel = true
		case "set_ct_abx":
			y.Support.SetCtAbx = true
		case "set_rgb":
			y.Support.SetRGB = true
		}
	}
}

func (y *YeeLight) setPower(val string) error {
	v := PowerValue(val)
	if v != On && v != Off {
		return errors.Wrapf(ErrInvalidRange, "invalid power value: %v", val)
	}
	y.propMutex.Lock()
	y.Power = v
	y.propMutex.Unlock()
	return nil
}

func (y *YeeLight) setBright(val string) error {
	v, err := strconv.Atoi(val)
	if err != nil {
		return errors.Wrapf(err, "could not convert %s to a bright value", val)
	}
	if v < 1 || v > 100 {
		return errors.Wrapf(ErrInvalidRange, "invalid bright value: %d", v)
	}
	y.propMutex.Lock()
	y.Brightness = v
	y.propMutex.Unlock()
	return nil
}

func (y *YeeLight) setColorMode(val string) error {
	v, err := strconv.Atoi(val)
	if err != nil {
		return errors.Wrapf(err, "could not convert %s to a color_mode value", val)
	}
	if v < 1 || v > 3 {
		return errors.Wrapf(ErrInvalidRange, "invalid color_mode value: %d", v)
	}
	y.propMutex.Lock()
	y.ColorMode = ColorModeValue(v)
	y.propMutex.Unlock()
	return nil
}

func (y *YeeLight) setColorTemperature(val string) error {
	v, err := strconv.Atoi(val)
	if err != nil {
		return errors.Wrapf(err, "could not convert %s to a ct value", val)
	}
	if v < 1700 || v > 6500 {
		return errors.Wrapf(ErrInvalidRange, "invalid ct value: %d", v)
	}
	y.propMutex.Lock()
	y.ColorTemperature = v
	y.propMutex.Unlock()
	return nil
}

func (y *YeeLight) setRGB(val string) error {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return errors.Wrapf(err, "could not convert %s to a rgb value", val)
	}
	v, err := NewRGB(intVal)
	if err != nil {
		return errors.Wrapf(ErrInvalidRange, "invalid rgb value: %d", v)
	}
	y.propMutex.Lock()
	y.RGB = v
	y.propMutex.Unlock()
	return nil
}

func (y *YeeLight) setHue(val string) error {
	v, err := strconv.Atoi(val)
	if err != nil {
		return errors.Wrapf(err, "could not convert %s to a hue value", val)
	}
	if v < 0 || v > 359 {
		return errors.Wrapf(ErrInvalidRange, "invalid hue value: %d", v)
	}
	y.propMutex.Lock()
	y.Hue = v
	y.propMutex.Unlock()
	return nil
}

func (y *YeeLight) setSaturation(val string) error {
	v, err := strconv.Atoi(val)
	if err != nil {
		return errors.Wrapf(err, "could not convert %s to a sat value", val)
	}
	if v < 0 || v > 100 {
		return errors.Wrapf(ErrInvalidRange, "invalid sat value: %d", v)
	}
	y.propMutex.Lock()
	y.Saturation = v
	y.propMutex.Unlock()
	return nil
}

func (y *YeeLight) setName(val string) {
	y.propMutex.Lock()
	y.Name = val
	y.propMutex.Unlock()
}
