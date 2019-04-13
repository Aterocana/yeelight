package yeelight

import (
	"net"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewDiscoveryService(t *testing.T) {
	t.Run("NewDiscoveryService()", func(t *testing.T) {
		ds := NewDiscoveryService()
		if _, ok := ds.(*discoveryService); !ok {
			t.Errorf("NewDiscoveryService() error: underlying type is not *discoveryService")
		}
	})
}

func Test_getMyIPs(t *testing.T) {
	t.Run("getMyIPs()", func(t *testing.T) {
		got, err := getMyIPs()
		if err != nil {
			t.Errorf("getMyIPs() error = %v", err)
			return
		}
		for _, val := range got {
			if val == "::1" {
				continue
			}
			if len(strings.Split(val, ".")) == 4 {
				continue
			}
			ipv6 := strings.Split(val, "::")
			if len(ipv6) == 2 {
				if len(strings.Split(ipv6[1], ":")) == 4 {
					continue
				}
			}

			t.Errorf("getMyIPs(): %s has not a x.x.x.x format", val)
		}
	})
}

type wrappedIPv4Addr struct {
	net.IP
}

func (a wrappedIPv4Addr) Network() string {
	return ""
}

func Test_addrIsIn(t *testing.T) {
	type args struct {
		addr net.Addr
		ips  []string
	}
	type test struct {
		name string
		args args
		want bool
	}

	tests := []test{
		test{
			name: "ip found",
			args: args{
				addr: wrappedIPv4Addr{net.IPv4(192, 168, 0, 1)},
				ips:  []string{"192.168.0.2", "192.168.0.1"},
			},
			want: true,
		},
		test{
			name: "missing ip",
			args: args{
				addr: wrappedIPv4Addr{net.IPv4(192, 168, 0, 1)},
				ips:  []string{"192.168.0.2", "192.168.0.3"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addrIsIn(tt.args.addr, tt.args.ips); got != tt.want {
				t.Errorf("addrIsIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoveryService_Open(t *testing.T) {
	t.Run("correct Open()", func(t *testing.T) {
		service := NewDiscoveryService()
		if err := service.Open(); err != nil {
			t.Errorf("discoveryService.Open() error = %v", err)
		}
	})

	t.Run("wrong Open()", func(t *testing.T) {
		service := NewDiscoveryService()
		if err := service.Open(); err != nil {
			t.Errorf("discoveryService.Open() error = %v", err)
		}
		errs := service.GetErrors()
		service.Open()
		select {
		case err := <-errs:
			if !strings.Contains(err.Error(), "bind: address already in use") {
				t.Errorf("discoveryService.Open() error = %v, expected \"address already in use\"", err)
			}
		case <-time.After(time.Second):
			t.Errorf("discoveryService.Open() expecting error, timed out")
		}
	})
}

func Test_discoveryService_GetDiscoveredDevices(t *testing.T) {
	service := NewDiscoveryService()
	c := service.GetDiscoveredDevices()
	underlyingService, ok := service.(*discoveryService)
	if !ok {
		t.Errorf("GetDiscoveredDevices() error: underlying type is not *discoveryService")
	}
	if reflect.DeepEqual(underlyingService.discoveredDevices, c) {
		t.Errorf("discoveryService.GetDiscoveredDevices() = %v, want %v", c, underlyingService.discoveredDevices)
	}
}

func Test_discoveryService_GetErrors(t *testing.T) {
	service := NewDiscoveryService()
	c := service.GetErrors()
	underlyingService, ok := service.(*discoveryService)
	if !ok {
		t.Errorf("GetErrors() error: underlying type is not *discoveryService")
	}
	if reflect.DeepEqual(underlyingService.errorsChan, c) {
		t.Errorf("discoveryService.GetErrors() = %v, want %v", c, underlyingService.errorsChan)
	}
}

func Test_discoveryService_DiscoveryRequest(t *testing.T) {
	service := NewDiscoveryService()
	t.Run("Calling DiscoveryRequest() before Open()", func(t *testing.T) {
		err := service.DiscoveryRequest()
		if err != nil {
			t.Errorf("discoveryService.DiscoveryRequest() expecting no errors, got: %+v", err)
			return
		}
		errs := service.GetErrors()
		devs := service.GetDiscoveredDevices()
		// timeout := time.After(time.Second)
		select {
		case <-time.After(200 * time.Millisecond):
			return
		case err := <-errs:
			t.Errorf("discoveryService.DiscoveryRequest() expecting no errors, got: %+v", err)
		case dev := <-devs:
			t.Errorf("discoveryService.DiscoveryRequest() expecting no devices, got: %v", dev)
		}
	})
}
