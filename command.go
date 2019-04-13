package yeelight

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Notification is the struct describing an event occurred on YeeLight device.
type Notification struct {
	Property string
	Status   string
}

// Effect is the effect parameter in a command.
type Effect string

const (
	// Sudden is the sudden effect in transition.
	Sudden Effect = "sudden"
	// Smooth is the smooth effect in transition.
	Smooth Effect = "smooth"
)

func (e Effect) isValid() bool {
	if e == Sudden || e == Smooth {
		return true
	}
	return false
}

func isValidDuration(d int) bool {
	return d >= 30
}

func (y *YeeLight) parseNotifications(msg []byte) []Notification {
	parsed := struct {
		Method string            `json:"method"`
		Params map[string]string `json:"params"`
	}{}
	s := strings.TrimSuffix(string(msg), "\r\n")

	if err := errors.WithStack(json.Unmarshal([]byte(s), &parsed)); err != nil {
		y.errs <- errors.Wrapf(err, "failed to parse notification: %s", s)
		return nil
	}

	var res []Notification
	for k, v := range parsed.Params {
		n := Notification{
			Property: k,
			Status:   v,
		}
		go y.updateProperty(n)
		res = append(res, n)
	}

	return res
}

func (y *YeeLight) updateProperty(n Notification) {
	switch n.Property {
	case "power":
		y.setPower(n.Status)
	case "bright":
		y.setBright(n.Status)
	case "color_mode":
		y.setColorMode(n.Status)
	case "ct":
		y.setColorTemperature(n.Status)
	case "rgb":
		y.setRGB(n.Status)
	case "hue":
		y.setHue(n.Status)
	case "sat":
		y.setSaturation(n.Status)
	case "name":
		y.setName(n.Status)
	}
}

// command is the struct describing a command to be sent over device's TCP connection.
type command struct {
	ID     int           `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

// json() is used to convert the command into bytes which can be send
// over device's TCP connection.
func (cmd *command) json() []byte {
	j, _ := json.Marshal(cmd)
	return append(j, '\r', '\n')
}

// Answer is the struct describing a command result received by device's TCP connection.
type Answer struct {
	ID     int           `json:"id"`
	Result []interface{} `json:"result"`
}

// nextCommand prepare the map entry with Answer, returning the index
// for the command.
func (y *YeeLight) nextCommand() int {
	y.idMutex.Lock()
	defer y.idMutex.Unlock()

	if y.pendingCmds == nil {
		y.pendingCmds = make(map[int]chan Answer)
	}

	// increment before in order to have a idCommand meaningful zero value.
	// If cmd.ID is zero, than an error occurs.
	y.idCommand++
	y.pendingCmds[y.idCommand] = make(chan Answer, 1)

	return y.idCommand
}

// newCommand is used to build a command.
func (y *YeeLight) newCommand(method string, params []interface{}) (*command, error) {
	for _, param := range params {
		switch p := param.(type) {
		case Effect:
			if !p.isValid() {
				return nil, errors.Wrapf(ErrInvalidType, "invalid effect: %s", p)
			}
		case PowerValue:
			if !p.isValid() {
				return nil, errors.Wrapf(ErrInvalidType, "invalid power value: %s", p)
			}
		case TurnOnValue:
			if !p.isValid() {
				return nil, errors.Wrapf(ErrInvalidType, "invalid turn on value: %d", p)
			}
		case int:
		default:
			return nil, errors.Wrapf(ErrInvalidType, "invalid parameter: %v", p)

		}
	}
	id := y.nextCommand()
	return &command{
		ID:     id,
		Method: method,
		Params: params,
	}, nil
}

// releaseAnswerChan frees Answer chan in pendingCmds map.
// if a is specified, then a is sent into the chan before closing it.
func (y *YeeLight) releaseAnswerChan(id int, a *Answer) {
	y.idMutex.Lock()
	defer y.idMutex.Unlock()
	c, ok := y.pendingCmds[id] // retrieving the chan of the open "transaction"
	if !ok {
		y.errs <- errors.Wrapf(ErrUnknownCommand, "unknown %d command", id)
		return
	}
	if a != nil {
		c <- *a
	}
	close(c) // transaction is successfully ended, so close its chan and delete it from pendingCmds map
	delete(y.pendingCmds, id)
}

// sendCommand sends a command to YeeLight device through its
// TCP connection.
func (y *YeeLight) sendCommand(cmd *command) (*Answer, error) {
	if y.tcpSocket == nil {
		y.releaseAnswerChan(cmd.ID, nil)
		return nil, errors.WithStack(ErrConnNotInitialized)
	}
	y.idMutex.RLock()
	respChan, ok := y.pendingCmds[cmd.ID]
	y.idMutex.RUnlock()
	if !ok {
		return nil, errors.WithStack(ErrFailedCmd)
	}
	y.tcpSocket.Write(cmd.json())
	select {
	case a := <-respChan:
		return &a, nil
	case <-time.After(commandTimeout):
		y.releaseAnswerChan(cmd.ID, nil)
		return nil, errors.Wrapf(ErrTimedOut, "failed command %v", cmd)
	}
}

// SendRGB is used to send a set_rgb command
func (y *YeeLight) SendRGB(r, g, b uint8, effect Effect, duration int) (*Answer, error) {
	if !isValidDuration(duration) {
		return nil, errors.Wrapf(ErrInvalidType, "invalid duration value: %d", duration)
	}
	val := RGBValue{r, g, b}

	cmd, err := y.newCommand("set_rgb", []interface{}{val.Get(), effect, duration})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return y.sendCommand(cmd)
}

// SetCTAbs is used to send a set_ct_abx command (set color temperature).
func (y *YeeLight) SetCTAbs(ct int, effect Effect, duration int) (*Answer, error) {
	if ct < 1700 || ct > 6500 {
		return nil, errors.WithStack(ErrInvalidRange)
	}
	if !isValidDuration(duration) {
		return nil, errors.Wrapf(ErrInvalidType, "invalid duration value: %d", duration)
	}
	cmd, err := y.newCommand("set_ct_abx", []interface{}{ct, effect, duration})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return y.sendCommand(cmd)
}

// SetHSV is used to change the color of YeeLight device.
func (y *YeeLight) SetHSV(hue, sat int, effect Effect, duration int) (*Answer, error) {
	if !isValidDuration(duration) {
		return nil, errors.Wrapf(ErrInvalidType, "invalid duration value: %d", duration)
	}
	cmd, err := y.newCommand("set_hsv", []interface{}{hue, sat, effect, duration})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return y.sendCommand(cmd)
}

// SetBright is used to change the brightness of YeeLight device.
func (y *YeeLight) SetBright(bright int, effect Effect, duration int) (*Answer, error) {
	if bright < 1 || bright > 100 {
		return nil, errors.Wrapf(ErrInvalidRange, "invalid bright value: %d", bright)
	}
	if !isValidDuration(duration) {
		return nil, errors.Wrapf(ErrInvalidType, "invalid duration value: %d", duration)
	}
	cmd, err := y.newCommand("set_bright", []interface{}{bright, effect, duration})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return y.sendCommand(cmd)
}

// SetPower is used to switch ON or OFF the YeeLight device.
func (y *YeeLight) SetPower(power PowerValue, effect Effect, duration int, mode TurnOnValue) (*Answer, error) {
	if !power.isValid() {
		return nil, errors.Wrapf(ErrInvalidType, "invalid power value: %v", power)
	}
	if !isValidDuration(duration) {
		return nil, errors.Wrapf(ErrInvalidType, "invalid duration value: %d", duration)
	}
	cmd, err := y.newCommand("set_power", []interface{}{power, effect, duration, mode})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return y.sendCommand(cmd)
}

// Toggle is used to send a toogle command
func (y *YeeLight) Toggle() (*Answer, error) {
	cmd, err := y.newCommand("toggle", []interface{}{})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return y.sendCommand(cmd)
}
