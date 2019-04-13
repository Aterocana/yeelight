package yeelight

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestNewRGB(t *testing.T) {
	type args struct {
		val int
	}
	type test struct {
		name    string
		arg     int
		want    RGBValue
		wantErr bool
	}
	tests := []test{
		test{
			name: "0xffffff",
			arg:  0xffffff,
			want: RGBValue{
				red:   0xff,
				green: 0xff,
				blue:  0xff,
			},
		},
		test{
			name: "0x000000",
			arg:  0x000000,
			want: RGBValue{
				red:   0x00,
				green: 0x00,
				blue:  0x00,
			},
		},
		test{
			name: "0x112233",
			arg:  0x112233,
			want: RGBValue{
				red:   0x11,
				green: 0x22,
				blue:  0x33,
			},
		},
		test{
			name:    "out of range (more than 0xffffff)",
			arg:     0xffffff1,
			wantErr: true,
		},
		test{
			name:    "out of range (less than 0)",
			arg:     -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRGB(tt.arg)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewRGB() error: wanted no errors, got %v", err)
					return
				}
				if errors.Cause(err) != ErrInvalidRange {
					t.Errorf("NewRGB() error: %+v", err)
					return
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRGB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYeeLight_setSupport(t *testing.T) {
	type fields struct {
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
	type test struct {
		name    string
		support string
		args    fields
		wantErr bool
	}

	tests := []test{
		test{
			name:    "setSupport: correct empty conversion",
			support: "",
			args:    fields{},
			wantErr: false,
		},
		test{
			name:    "setSupport: wrong empty conversion",
			support: "get_prop",
			args:    fields{},
			wantErr: true,
		},
		test{
			name:    "setSupport: correct single fields conversion",
			support: "get_prop",
			args: fields{
				GetProp: true,
			},
			wantErr: false,
		},
		test{
			name:    "setSupport: wrong single fields conversion",
			support: "get_prop",
			args: fields{
				GetProp: false,
			},
			wantErr: true,
		},
		test{
			name:    "setSupport: correct ignoring unknown fields",
			support: "get_prop gets_props",
			args: fields{
				GetProp: true,
			},
			wantErr: false,
		},
		test{
			name:    "setSupport: correct multiple fields conversion",
			support: "get_prop set_rgb",
			args: fields{
				GetProp: true,
				SetRGB:  true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &YeeLight{
				Support: SupportedFeatures{
					GetProp:    tt.args.GetProp,
					SetDefault: tt.args.SetDefault,
					SetPower:   tt.args.SetPower,
					Toggle:     tt.args.Toggle,
					SetBright:  tt.args.SetBright,
					StartCF:    tt.args.StartCF,
					StopCF:     tt.args.StopCF,
					SetScene:   tt.args.SetScene,
					CronAdd:    tt.args.CronAdd,
					CronGet:    tt.args.CronGet,
					CronDel:    tt.args.CronDel,
					SetCtAbx:   tt.args.SetCtAbx,
					SetRGB:     tt.args.SetRGB,
				},
			}
			aux := &YeeLight{}
			aux.setSupport(tt.support)
			if reflect.DeepEqual(y.Support, aux.Support) {
				if tt.wantErr {
					t.Errorf("setSupport error: wanted %v being different from %v", aux.Support, y.Support)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("setSupport error: wanted %v being equal to %v", aux.Support, y.Support)
				}
			}
		})
	}
}

func TestColorModeValue_String(t *testing.T) {
	type test struct {
		name string
		c    ColorModeValue
		want string
	}
	tests := []test{
		test{
			name: "colorMode: rgb",
			c:    ColorModeValue(1),
			want: "rgb",
		},
		test{
			name: "colorMode: temperature",
			c:    ColorModeValue(2),
			want: "temperature",
		},
		test{
			name: "colorMode: hsv",
			c:    ColorModeValue(3),
			want: "hsv",
		},
		test{
			name: "colorMode: less than 1",
			c:    ColorModeValue(0),
			want: "unknown mode",
		},
		test{
			name: "colorMode: more than 3",
			c:    ColorModeValue(4),
			want: "unknown mode",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("ColorModeValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRGBValue_Get(t *testing.T) {
	type fields struct {
		red   uint8
		green uint8
		blue  uint8
	}
	type test struct {
		name   string
		fields fields
		want   int
	}
	tests := []test{
		test{
			name: "red",
			fields: fields{
				red:   0xff,
				green: 0,
				blue:  0,
			},
			want: 16711680,
		},
		test{
			name: "green",
			fields: fields{
				red:   0,
				green: 0xff,
				blue:  0,
			},
			want: 65280,
		},
		test{
			name: "blue",
			fields: fields{
				red:   0,
				green: 0,
				blue:  0xff,
			},
			want: 255,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgb := RGBValue{
				red:   tt.fields.red,
				green: tt.fields.green,
				blue:  tt.fields.blue,
			}
			if got := rgb.Get(); got != tt.want {
				t.Errorf("RGBValue.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPowerValue_isValid(t *testing.T) {
	tests := []struct {
		name string
		p    PowerValue
		want bool
	}{
		{"ON value", PowerValue("on"), true},
		{"OFF value", PowerValue("off"), true},
		{"any other value", PowerValue("any other value"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.isValid(); got != tt.want {
				t.Errorf("PowerValue.isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTurnOnValue_isValid(t *testing.T) {
	tests := []struct {
		name string
		t    TurnOnValue
		want bool
	}{
		// TODO: completare
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.isValid(); got != tt.want {
				t.Errorf("TurnOnValue.isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
