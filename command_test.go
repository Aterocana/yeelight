package yeelight

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
)

type orderedNotifications []Notification

func (n orderedNotifications) Len() int {
	return len(n)
}

func (n orderedNotifications) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n orderedNotifications) Less(i, j int) bool {
	return strings.Compare(n[i].Status, n[j].Status) < 0
}

func Test_newCommand(t *testing.T) {
	type args struct {
		id     int
		method string
		params []interface{}
	}
	type test struct {
		name    string
		args    args
		want    *command
		wantErr bool
		errType error
	}
	tests := []test{
		test{
			name: "good parameters",
			args: args{
				method: "method",
				params: []interface{}{
					Effect("smooth"),
					1,
				},
			},
			want: &command{
				Method: "method",
				Params: []interface{}{
					Effect("smooth"),
					1,
				},
			},
			wantErr: false,
		},
		test{
			name: "wrong parameters",
			args: args{
				method: "method",
				params: []interface{}{
					Effect("smooth"),
					1,
					0.2,
				},
			},
			wantErr: true,
			errType: ErrInvalidType,
		},
		test{
			name: "wrong effect parameters",
			args: args{
				method: "method",
				params: []interface{}{
					Effect("smoother"),
					1,
				},
			},
			wantErr: true,
			errType: ErrInvalidType,
		},
		test{
			name: "good parameters after wrong cmds (to test IDs)",
			args: args{
				method: "method",
				params: []interface{}{
					Effect("smooth"),
					1,
				},
			},
			want: &command{
				Method: "method",
				Params: []interface{}{
					Effect("smooth"),
					1,
				},
			},
			wantErr: false,
		},
	}
	y := &YeeLight{}
	index := 0
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := y.newCommand(tt.args.method, tt.args.params)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("newCommand() expected no error, got %v", err)
					return
				}
				if tt.errType == nil {
					return
				}
				if tt.errType != errors.Cause(err) {
					t.Errorf("newCommand() error = %v, want %v", err, tt.errType)
				}
				return
			}
			if tt.wantErr {
				t.Errorf("newCommand() expected errors, got no errors")
				return
			}
			// increment index to matching cmd ID, since a new ID is released only on
			// correct cmds, index is incremented after knowing tt.wantErr is false.
			index++
			tt.want.ID = index
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_command_JSON(t *testing.T) {
	type fields struct {
		ID     int
		Method string
		Params []interface{}
	}
	type test struct {
		name   string
		fields fields
		want   []byte //without trailing \r\n
	}
	tests := []test{
		test{
			name:   "json",
			fields: fields{0, "set_power", []interface{}{"on", "smooth", 500}},
			want:   []byte(`{"id":0,"method":"set_power","params":["on","smooth",500]}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &command{
				ID:     tt.fields.ID,
				Method: tt.fields.Method,
				Params: tt.fields.Params,
			}
			got := cmd.json()
			tt.want = append(tt.want, '\r', '\n')
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("command.JSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYeeLight_updateProperty(t *testing.T) {
	tests := []struct {
		name string
		n    Notification
	}{
		{
			"power ON",
			Notification{"power", "on"},
		},
		{
			"power OFF",
			Notification{"power", "off"},
		},
		{
			"bright",
			Notification{"bright", "90"},
		},
		{
			"color mode: ColorMode",
			Notification{"color_mode", "1"},
		},
		{
			"color mode: ColorTemperature",
			Notification{"color_mode", "2"},
		},
		{
			"color mode: HSV",
			Notification{"color_mode", "3"},
		},
		{
			"color temperature",
			Notification{"ct", "4000"},
		},
		{
			"rgb",
			Notification{"rgb", "16711680"},
		},
		{
			"hue",
			Notification{"hue", "100"},
		},
		{
			"saturation",
			Notification{"sat", "35"},
		},
		{
			"name",
			Notification{"name", "my-bulb"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &YeeLight{}
			y.updateProperty(tt.n)
			switch tt.n.Property {
			case "power":
				if y.Power != PowerValue(tt.n.Status) {
					t.Errorf("updateProperty(), expected power to be %v, instead is %v", PowerValue(tt.n.Status), y.Power)
				}
			case "bright":
				v, _ := strconv.Atoi(tt.n.Status)
				if y.Brightness != v {
					t.Errorf("updateProperty(), expected brightness to be %v, instead is %v", v, y.Brightness)
				}
			case "color_mode":
				v, _ := strconv.Atoi(tt.n.Status)
				if y.ColorMode != ColorModeValue(v) {
					t.Errorf("updateProperty(), expected color mode to be %v, instead is %v", ColorModeValue(v), y.ColorMode)
				}
			case "ct":
				v, _ := strconv.Atoi(tt.n.Status)
				if y.ColorTemperature != v {
					t.Errorf("updateProperty(), expected color temperature to be %v, instead is %v", v, y.ColorTemperature)
				}
			case "rgb":
				v, _ := strconv.Atoi(tt.n.Status)
				rgb, _ := NewRGB(v)
				if y.RGB != rgb {
					t.Errorf("updateProperty(), expected RGB to be %v, instead is %v", rgb, y.RGB)
				}
			case "hue":
				v, _ := strconv.Atoi(tt.n.Status)
				if y.Hue != v {
					t.Errorf("updateProperty(), expected hue to be %v, instead is %v", v, y.Hue)
				}
			case "sat":
				v, _ := strconv.Atoi(tt.n.Status)
				if y.Saturation != v {
					t.Errorf("updateProperty(), expected saturation to be %v, instead is %v", v, y.Saturation)
				}
			case "name":
				if y.Name != tt.n.Status {
					t.Errorf("updateProperty(), expected name to be %v, instead is %v", tt.n.Status, y.Name)
				}
			}
		})
	}
}

func TestYeeLight_parseNotifications(t *testing.T) {
	tests := []struct {
		name    string
		msg     []byte
		want    []Notification
		wantErr bool
	}{
		{
			"simple notification",
			[]byte("{\"method\":\"props\",\"params\":{\"power\":\"on\"}}\r\n"),
			[]Notification{
				Notification{"power", "on"},
			},
			false,
		},
		{
			"complex notification",
			[]byte("{\"method\":\"props\",\"params\":{\"power\":\"on\",\"bright\":\"20\"}}\r\n"),
			[]Notification{
				Notification{"power", "on"},
				Notification{"bright", "20"},
			},
			false,
		},
		{
			"malformed notification (malformed JSON)",
			[]byte("{}\"method\":\"props\",\"params\":{\"power\":\"on\",\"bright\":\"20\"}}\r\n"),
			[]Notification{
				Notification{"power", "on"},
				Notification{"bright", "20"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &YeeLight{
				errs: make(chan error),
			}
			if tt.wantErr {
				go y.parseNotifications(tt.msg)
				select {
				case <-y.errs:
				case <-time.After(200 * time.Millisecond):
					t.Errorf("parseNotification(): expected error, timed out")
				}
				return
			}
			got := orderedNotifications(y.parseNotifications(tt.msg))
			sort.Sort(got)

			wanted := orderedNotifications(tt.want)
			sort.Sort(wanted)
			if !reflect.DeepEqual(got, wanted) {
				t.Errorf("YeeLight.parseNotifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEffect_isValid(t *testing.T) {
	type test struct {
		name string
		e    Effect
		want bool
	}
	tests := []test{
		test{
			name: "smooth",
			e:    Effect("smooth"),
			want: true,
		},
		test{
			name: "sudden",
			e:    Effect("sudden"),
			want: true,
		},
		test{
			name: "any othe value",
			e:    Effect("any other value"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.isValid(); got != tt.want {
				t.Errorf("Effect.isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidDuration(t *testing.T) {
	tests := []struct {
		name string
		d    int
		want bool
	}{
		{"correct duration", 30, true},
		{"wrong duration", 29, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidDuration(tt.d); got != tt.want {
				t.Errorf("isValidDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYeeLight_nextCommand(t *testing.T) {
	t.Run("nextCommand()", func(t *testing.T) {
		count := 100
		y := &YeeLight{}
		wg := sync.WaitGroup{}
		wg.Add(count)
		mux := sync.Mutex{}
		res := make(map[int]bool)
		for index := 0; index < count; index++ {
			go func(id int) {
				id = y.nextCommand()
				mux.Lock()
				res[id] = true
				mux.Unlock()
				wg.Done()
			}(index)
		}
		wg.Wait()
		if len(res) != count {
			t.Errorf("nextCommand() failed: got %d index instead of %d", len(res), count)
		}
	})
}
