package yeelight

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func Test_parseFromMap(t *testing.T) {
	type test struct {
		name    string
		lines   map[string]string
		want    *YeeLight
		wantErr bool
		errType error // expected error, leaving it nil if a generic error is expected.
	}
	tests := []test{
		test{
			name: "full map (power on)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			want: &YeeLight{
				CacheControl:    "CacheControlValue",
				Location:        "LocationValue",
				ID:              "IDValue",
				Name:            "NameValue",
				Model:           "ModelValue",
				FirmwareVersion: "FWVerValue",
				Support: SupportedFeatures{
					GetProp:    true,
					SetDefault: true,
					SetPower:   true,
				},
				Power:            On,
				Brightness:       1,
				ColorMode:        ColorTemperature,
				ColorTemperature: 1700,
				RGB: RGBValue{
					red:   0,
					green: 0,
					blue:  4,
				},
				Hue:        5,
				Saturation: 6,
			},
			wantErr: false,
		},
		test{
			name: "full map (power off)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "off",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			want: &YeeLight{
				CacheControl:    "CacheControlValue",
				Location:        "LocationValue",
				ID:              "IDValue",
				Name:            "NameValue",
				Model:           "ModelValue",
				FirmwareVersion: "FWVerValue",
				Support: SupportedFeatures{
					GetProp:    true,
					SetDefault: true,
					SetPower:   true,
				},
				Power:            Off,
				Brightness:       1,
				ColorMode:        ColorTemperature,
				ColorTemperature: 1700,
				RGB: RGBValue{
					red:   0,
					green: 0,
					blue:  4,
				},
				Hue:        5,
				Saturation: 6,
			},
			wantErr: false,
		},
		test{
			name: "missing field cache_control",
			lines: map[string]string{
				"location":   "LocationValue",
				"id":         "IDValue",
				"name":       "NameValue",
				"model":      "ModelValue",
				"fw_ver":     "FWVerValue",
				"support":    "get_prop set_default set_power",
				"power":      "on",
				"bright":     "1",
				"color_mode": "2",
				"ct":         "1700",
				"rgb":        "4",
				"hue":        "5",
				"sat":        "6",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "missing field location",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "missing field id",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "missing field name",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "missing field model",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "missing field fw_ver",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "missing field support",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "missing field power",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "missing field bright",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "missing field color_mode",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "missing field ct",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "missing field rgb",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "missing field hue",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "missing field sat",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
			},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
		test{
			name: "invalid power",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "INVALID_POWER",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrInvalidRange,
		},
		test{
			name: "invalid bright (wrong format: Atoi fails)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "off",
				"bright":        "bright",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
		},
		test{
			name: "invalid bright (less than 1)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "off",
				"bright":        "0",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrInvalidRange,
		},
		test{
			name: "invalid bright (more than 100)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "off",
				"bright":        "101",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrInvalidRange,
		},
		test{
			name: "invalid color_mode (wrong format: Atoi fails)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "off",
				"bright":        "1",
				"color_mode":    "color_mode",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
		},
		test{
			name: "invalid color_mode (less than 1)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "off",
				"bright":        "1",
				"color_mode":    "0",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrInvalidRange,
		},
		test{
			name: "invalid color_mode (more than 3)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "4",
				"ct":            "1700",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrInvalidRange,
		},
		test{
			name: "invalid ct (ColorTemperature) (wrong format: Atoi fails)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "3",
				"ct":            "ct",
				"rgb":           "4",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
		},
		test{
			name: "invalid RGB (less than 0)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "-1",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrInvalidRange,
		},
		test{
			name: "invalid RGB (more than 0xffffff)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "16777216",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrInvalidRange,
		},
		test{
			name: "invalid RGB (wrong format: Atoi fails)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "rgb",
				"hue":           "5",
				"sat":           "6",
			},
			wantErr: true,
		},
		test{
			name: "invalid Hue (wrong format: Atoi fails)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "400",
				"hue":           "hue",
				"sat":           "6",
			},
			wantErr: true,
		},
		test{
			name: "invalid Hue (less than 0)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "400",
				"hue":           "-1",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrInvalidRange,
		},
		test{
			name: "invalid Hue (more than 359)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "400",
				"hue":           "360",
				"sat":           "6",
			},
			wantErr: true,
			errType: ErrInvalidRange,
		},
		test{
			name: "invalid Saturation (wrong format: Atoi fails)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "400",
				"hue":           "1",
				"sat":           "sat",
			},
			wantErr: true,
		},
		test{
			name: "invalid Saturation (less than 0)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "400",
				"hue":           "1",
				"sat":           "-1",
			},
			wantErr: true,
			errType: ErrInvalidRange,
		},
		test{
			name: "invalid Saturation (more than 100)",
			lines: map[string]string{
				"cache-control": "CacheControlValue",
				"location":      "LocationValue",
				"id":            "IDValue",
				"name":          "NameValue",
				"model":         "ModelValue",
				"fw_ver":        "FWVerValue",
				"support":       "get_prop set_default set_power",
				"power":         "on",
				"bright":        "1",
				"color_mode":    "2",
				"ct":            "1700",
				"rgb":           "400",
				"hue":           "360",
				"sat":           "101",
			},
			wantErr: true,
			errType: ErrInvalidRange,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFromMap(tt.lines)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("parseFromMap() error = %+v, wantErr %v", err, tt.wantErr)
					return
				}
				if errors.Cause(err) != tt.errType && tt.errType != nil {
					t.Errorf("parseFromMap() error = %v, expected error %v", err, tt.errType)
					return
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newFromAdvertisement(t *testing.T) {
	type test struct {
		name    string
		msg     []byte
		want    *YeeLight
		wantErr bool
		errType error // expected error, leaving it nil if a generic error is expected.
	}
	tests := []test{
		test{
			name: "discovery answer",
			msg:  []byte("HTTP/1.1 200 OK\r\nCache-Control: max-age-3600\r\nDate: \r\nExt: \r\nLocation: yeelight://192.168.0.20:55443\r\nNTS: ssdp:alive\r\nServer: POSIX, UPnP/1.0 YGLC/1\r\nid: 0x000000000458bdfa\r\nmodel: color\r\nfw_ver: 70\r\nsupport: get_prop set_default set_power toggle set_bright start_cf stop_cf set_scene cron_add cron_get cron_del set_ct_abx set_rgb set_hsv set_adjust adjust_bright adjust_ct adjust_color set_music set_name\r\npower: off\r\nbright: 53\r\ncolor_mode: 2\r\nct: 2634\r\nrgb: 16711680\r\nhue: 359\r\nsat: 100\r\nname: my-bulb\r\n"),
			want: &YeeLight{
				CacheControl:    "max-age-3600",
				Location:        "192.168.0.20:55443",
				ID:              "0x000000000458bdfa",
				Model:           "color",
				FirmwareVersion: "70",
				Support: SupportedFeatures{
					GetProp:    true,
					SetDefault: true,
					SetPower:   true,
					Toggle:     true,
					SetBright:  true,
					StartCF:    true,
					StopCF:     true,
					SetScene:   true,
					CronAdd:    true,
					CronGet:    true,
					CronDel:    true,
					SetCtAbx:   true,
					SetRGB:     true,
				},
				Power:            Off,
				Brightness:       53,
				ColorMode:        ColorTemperature,
				ColorTemperature: 2634,
				RGB: RGBValue{
					red:   0xff,
					green: 0,
					blue:  0,
				},
				Hue:        359,
				Saturation: 100,
				Name:       "my-bulb",
			},
		},
		test{
			name: "advertisement message",
			msg:  []byte("NOTIFY * HTTP/1.1\r\nCache-Control: max-age-3600\r\nDate: \r\nExt: \r\nLocation: yeelight://192.168.0.20:55443\r\nNTS: ssdp:alive\r\nServer: POSIX, UPnP/1.0 YGLC/1\r\nid: 0x000000000458bdfa\r\nmodel: color\r\nfw_ver: 70\r\nsupport: get_prop set_default set_power toggle set_bright start_cf stop_cf set_scene cron_add cron_get cron_del set_ct_abx set_rgb set_hsv set_adjust adjust_bright adjust_ct adjust_color set_music set_name\r\npower: off\r\nbright: 53\r\ncolor_mode: 2\r\nct: 2634\r\nrgb: 16711680\r\nhue: 359\r\nsat: 100\r\nname: my-bulb\r\n"),
			want: &YeeLight{
				CacheControl:    "max-age-3600",
				Location:        "192.168.0.20:55443",
				ID:              "0x000000000458bdfa",
				Model:           "color",
				FirmwareVersion: "70",
				Support: SupportedFeatures{
					GetProp:    true,
					SetDefault: true,
					SetPower:   true,
					Toggle:     true,
					SetBright:  true,
					StartCF:    true,
					StopCF:     true,
					SetScene:   true,
					CronAdd:    true,
					CronGet:    true,
					CronDel:    true,
					SetCtAbx:   true,
					SetRGB:     true,
				},
				Power:            Off,
				Brightness:       53,
				ColorMode:        ColorTemperature,
				ColorTemperature: 2634,
				RGB: RGBValue{
					red:   0xff,
					green: 0,
					blue:  0,
				},
				Hue:        359,
				Saturation: 100,
				Name:       "my-bulb",
			},
		},
		test{
			name:    "empty message",
			msg:     []byte{},
			want:    &YeeLight{},
			wantErr: true,
			errType: io.EOF,
		},
		test{
			name:    "wrong start line",
			msg:     []byte("WRONG\r\nCache-Control: max-age-3600\r\nDate: \r\nExt: \r\nLocation: yeelight://192.168.0.20:55443\r\nNTS: ssdp:alive\r\nServer: POSIX, UPnP/1.0 YGLC/1\r\nid: 0x000000000458bdfa\r\nmodel: color\r\nfw_ver: 70\r\nsupport: get_prop set_default set_power toggle set_bright start_cf stop_cf set_scene cron_add cron_get cron_del set_ct_abx set_rgb set_hsv set_adjust adjust_bright adjust_ct adjust_color set_music set_name\r\npower: off\r\nbright: 53\r\ncolor_mode: 2\r\nct: 2634\r\nrgb: 16711680\r\nhue: 359\r\nsat: 100\r\nname: my-bulb\r\n"),
			want:    &YeeLight{},
			wantErr: true,
			errType: ErrWrongAdvertisement,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newFromAdvertisement(tt.msg)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("newFromAdvertisement() error = %+v, wantErr %v", err, tt.wantErr)
					return
				}
				if errors.Cause(err) != tt.errType && tt.errType != nil {
					t.Errorf("parseFromMap() error = %v, expected error %v", err, tt.errType)
					return
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newFromAdvertisement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_MarshalJSON(t *testing.T) {
	type test struct {
		name    string
		arg     *YeeLight
		want    []byte
		wantErr bool
	}
	tests := []test{
		test{
			name: "serializing JSON format",
			arg: &YeeLight{
				CacheControl:    "max-age-3600",
				Location:        "192.168.0.20",
				ID:              "0x000000000458bdfa",
				Model:           "color",
				FirmwareVersion: "70",
				Support: SupportedFeatures{
					GetProp:    true,
					SetDefault: true,
					SetPower:   true,
					Toggle:     true,
					SetBright:  true,
					StartCF:    true,
					StopCF:     true,
					SetScene:   true,
					CronAdd:    true,
					CronGet:    true,
					CronDel:    true,
					SetCtAbx:   true,
					SetRGB:     true,
				},
				Power:            Off,
				Brightness:       53,
				ColorMode:        ColorTemperature,
				ColorTemperature: 2634,
				RGB: RGBValue{
					red:   0xff,
					green: 0,
					blue:  0,
				},
				Hue:        359,
				Saturation: 100,
				Name:       "my-bulb",
			},
			want: []byte(`{"cache_control":"max-age-3600","location":"192.168.0.20","id":"0x000000000458bdfa","model":"color","fw_ver":"70","support":{"get_prop":true,"set_default":true,"set_power":true,"toggle":true,"set_bright":true,"start_cf":true,"stop_cf":true,"set_scene":true,"cron_add":true,"cron_get":true,"cron_del":true,"set_ct_abx":true,"set_rgb":true},"power":"off","brightness":53,"color_mode":"temperature","color_temperature":2634,"rgb":{"r":255,"g":0,"b":0},"hue":359,"saturation":100,"name":"my-bulb"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aux, err := json.Marshal(tt.arg)
			if err != nil {
				if tt.wantErr {
					return
				}
				return
			}
			if bytes.Compare(aux, tt.want) != 0 {
				t.Errorf("MarshalJSON(): expected %v be equal to %v\n", aux, tt.want)
			}
		})
	}
}
func Test_UnmarshalJSON(t *testing.T) {
	type test struct {
		name    string
		arg     []byte
		want    *YeeLight
		wantErr bool
	}
	tests := []test{
		test{
			name: "serializing JSON format",
			want: &YeeLight{
				CacheControl:    "max-age-3600",
				Location:        "192.168.0.20",
				ID:              "0x000000000458bdfa",
				Model:           "color",
				FirmwareVersion: "70",
				Support: SupportedFeatures{
					GetProp:    true,
					SetDefault: true,
					SetPower:   true,
					Toggle:     true,
					SetBright:  true,
					StartCF:    true,
					StopCF:     true,
					SetScene:   true,
					CronAdd:    true,
					CronGet:    true,
					CronDel:    true,
					SetCtAbx:   true,
					SetRGB:     true,
				},
				Power:            Off,
				Brightness:       53,
				ColorMode:        ColorTemperature,
				ColorTemperature: 2634,
				RGB: RGBValue{
					red:   0xff,
					green: 0,
					blue:  0,
				},
				Hue:        359,
				Saturation: 100,
				Name:       "my-bulb",
			},
			arg: []byte(`{"cache_control":"max-age-3600","location":"192.168.0.20","id":"0x000000000458bdfa","model":"color","fw_ver":"70","support":{"get_prop":true,"set_default":true,"set_power":true,"toggle":true,"set_bright":true,"start_cf":true,"stop_cf":true,"set_scene":true,"cron_add":true,"cron_get":true,"cron_del":true,"set_ct_abx":true,"set_rgb":true},"power":"off","brightness":53,"color_mode":"temperature","color_temperature":2634,"rgb":{"r":255,"g":0,"b":0},"hue":359,"saturation":100,"name":"my-bulb"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aux := YeeLight{}
			err := json.Unmarshal(tt.arg, &aux)
			if err != nil {
				if tt.wantErr {
					return
				}
				return
			}
			if !reflect.DeepEqual(&aux, tt.want) {
				t.Errorf("MarshalJSON(): expected %v be equal to %v\n", &aux, tt.want)
			}
		})
	}
}

func Test_String(t *testing.T) {
	type test struct {
		name    string
		arg     *YeeLight
		want    string
		wantErr bool
	}
	tests := []test{
		test{
			name: "serializing JSON format",
			arg: &YeeLight{
				CacheControl:    "max-age-3600",
				Location:        "192.168.0.20",
				ID:              "0x000000000458bdfa",
				Model:           "color",
				FirmwareVersion: "70",
				Support: SupportedFeatures{
					GetProp:    true,
					SetDefault: true,
					SetPower:   true,
					Toggle:     true,
					SetBright:  true,
					StartCF:    true,
					StopCF:     true,
					SetScene:   true,
					CronAdd:    true,
					CronGet:    true,
					CronDel:    true,
					SetCtAbx:   true,
					SetRGB:     true,
				},
				Power:            Off,
				Brightness:       53,
				ColorMode:        ColorTemperature,
				ColorTemperature: 2634,
				RGB: RGBValue{
					red:   0xff,
					green: 0,
					blue:  0,
				},
				Hue:        359,
				Saturation: 100,
				Name:       "my-bulb",
			},
			want: `{"cache_control":"max-age-3600","location":"192.168.0.20","id":"0x000000000458bdfa","model":"color","fw_ver":"70","support":{"get_prop":true,"set_default":true,"set_power":true,"toggle":true,"set_bright":true,"start_cf":true,"stop_cf":true,"set_scene":true,"cron_add":true,"cron_get":true,"cron_del":true,"set_ct_abx":true,"set_rgb":true},"power":"off","brightness":53,"color_mode":"temperature","color_temperature":2634,"rgb":{"r":255,"g":0,"b":0},"hue":359,"saturation":100,"name":"my-bulb"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.arg.String() != tt.want {
				t.Errorf("String(): expected %s be equal to %s\n", tt.arg.String(), tt.want)
			}
		})
	}
}

func TestYeeLight_GetErrors(t *testing.T) {
	type test struct {
		name    string
		errChan chan error
	}
	tests := []test{
		test{
			name:    "GetErrors",
			errChan: make(chan error),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &YeeLight{
				errs: tt.errChan,
			}
			got := y.GetErrors()
			err := errors.New("test error")
			go func() {
				tt.errChan <- err
			}()
			select {
			case errRcv := <-got:
				if !reflect.DeepEqual(err, errRcv) {
					t.Errorf("YeeLight.GetErrors(): expecting %v, got %v", err, errRcv)
				}
				return
			case <-time.After(100 * time.Millisecond):
				t.Errorf("YeeLight.GetErrors(): timed out")
			}
		})
	}
}

func TestYeeLight_GetNotification(t *testing.T) {
	type test struct {
		name   string
		events chan Notification
	}
	tests := []test{
		test{
			name:   "GetNotification",
			events: make(chan Notification),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &YeeLight{
				events: tt.events,
			}
			got := y.GetNotification()
			n := Notification{"prop", "status"}
			go func() {
				tt.events <- n
			}()
			select {
			case nRcv := <-got:
				if !reflect.DeepEqual(n, nRcv) {
					t.Errorf("YeeLight.GetNotification(): expecting %v, got %v", n, nRcv)
				}
				return
			case <-time.After(100 * time.Millisecond):
				t.Errorf("YeeLight.GetNotification(): timed out")
			}
		})
	}
}
