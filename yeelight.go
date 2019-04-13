package yeelight

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// YeeLight is a struct representing a YeeLight Smart LED Bulb.
type YeeLight struct {
	CacheControl    string            `json:"cache_control,omitempty"`
	Location        string            `json:"location,omitempty"`
	ID              string            `json:"id,omitempty"`
	Model           string            `json:"model,omitempty"`
	FirmwareVersion string            `json:"fw_ver,omitempty"`
	Support         SupportedFeatures `json:"support"`

	Power            PowerValue     `json:"power,omitempty"`
	Brightness       int            `json:"brightness,omitempty"`
	ColorMode        ColorModeValue `json:"color_mode,omitempty"`
	ColorTemperature int            `json:"color_temperature,omitempty"`
	RGB              RGBValue       `json:"rgb,omitempty"`
	Hue              int            `json:"hue,omitempty"`
	Saturation       int            `json:"saturation,omitempty"`

	Name string `json:"name"`

	propMutex sync.RWMutex

	tcpSocket net.Conn
	connMutex sync.RWMutex

	idMutex sync.RWMutex
	// idCommand is the command ID used to identify correspondant Answer
	idCommand int
	// pendingCmds is the map where are stored chan for answers of sent commands.
	// Once the "transaction" is done, the chan is closed and the map entry deleted.
	pendingCmds map[int]chan Answer

	errs   chan error
	events chan Notification
}

func (y *YeeLight) String() string {
	b, _ := json.Marshal(y)
	return string(b)
	// return fmt.Sprintf("<id: %s, fw: %s, IP: %s> support :%s", y.ID, y.FirmwareVersion, y.Location, y.Support)
}

// GetErrors returns error chan where YeeLight device errors are sent.
func (y *YeeLight) GetErrors() <-chan error {
	return y.errs
}

// GetNotification returns events chan where YeeLight device events are sent.
func (y *YeeLight) GetNotification() <-chan Notification {
	return y.events
}

// newFromAdvertisement build a YeeLight struct from an advertisement
// message.
func newFromAdvertisement(msg []byte) (*YeeLight, error) {
	buf := bytes.NewBuffer(msg)
	chunk, err := buf.ReadBytes('\n')
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !(bytes.Compare(chunk, discoveryAnswerHeader) == 0 || bytes.Compare(chunk, advertisementHeader) == 0) {
		return nil, errors.Wrapf(ErrWrongAdvertisement, "wrong advertisement header: %s", string(chunk))
	}

	lines := make(map[string]string)

	for {
		chunk, err = buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrap(err, "newFromAdventisement failed")
		}
		header := strings.SplitN(string(chunk), ": ", 2)
		if len(header) >= 2 {
			lines[strings.ToLower(header[0])] = strings.TrimSuffix(header[1], "\r\n")
		}
	}

	return parseFromMap(lines)
}

// parseFromMap build a YeeLight struct from a map key-value.
func parseFromMap(lines map[string]string) (*YeeLight, error) {
	y := &YeeLight{}
	var err error
	var found bool

	y.CacheControl, found = lines["cache-control"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing cache control header")
	}

	y.Location, found = lines["location"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing location header")
	}
	y.Location = strings.TrimPrefix(y.Location, "yeelight://")

	y.ID, found = lines["id"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing id header")
	}

	y.Name, found = lines["name"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing name header")
	}

	y.Model, found = lines["model"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing model header")
	}

	y.FirmwareVersion, found = lines["fw_ver"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing fw_ver header")
	}

	supportedFeatures, found := lines["support"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing support header")
	}
	y.setSupport(supportedFeatures)

	val, found := lines["power"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing power header")
	}
	err = y.setPower(val)
	if err != nil {
		return nil, errors.Wrap(err, "wrong power value")
	}

	val, found = lines["bright"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing bright header")
	}
	err = y.setBright(val)
	if err != nil {
		return nil, errors.Wrap(err, "wrong bright value")
	}

	val, found = lines["color_mode"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing color_mode header")
	}
	err = y.setColorMode(val)
	if err != nil {
		return nil, errors.Wrap(err, "wrong color_mode value")
	}

	val, found = lines["ct"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing ct header")
	}
	err = y.setColorTemperature(val)
	if err != nil {
		return nil, errors.Wrap(err, "wrong ct value")
	}

	val, found = lines["rgb"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing rgb header")
	}
	err = y.setRGB(val)
	if err != nil {
		return nil, errors.Wrap(err, "wrong rgb value")
	}

	val, found = lines["hue"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing hue header")
	}
	err = y.setHue(val)
	if err != nil {
		return nil, errors.Wrap(err, "wrong hue value")
	}

	val, found = lines["sat"]
	if !found {
		return nil, errors.Wrap(ErrWrongAdvertisement, "missing sat header")
	}
	err = y.setSaturation(val)
	if err != nil {
		return nil, errors.Wrap(err, "wrong sat value")
	}

	return y, nil
}
