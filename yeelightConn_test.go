package yeelight

import (
	"bufio"
	"io"
	"net"
	"testing"

	"github.com/pkg/errors"
)

var mockTCP struct {
	net.Listener
	listening bool
	stop      chan struct{}
	buf       io.Reader
}

func startMockTCP() error {
	if mockTCP.listening {
		return nil
	}
	mockTCP.stop = make(chan struct{})
	var err error
	mockTCP.Listener, err = net.Listen("tcp4", ":0")
	if err != nil {
		return err
	}
	mockTCP.listening = true

	go func() {
		c, _ := mockTCP.Accept()
		mockTCP.buf = bufio.NewReader(c)
	}()

	return nil
}

func stopMockTCP() error {
	mockTCP.listening = false
	return mockTCP.Listener.Close()
}

func TestYeeLight_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		errType error
	}{
		{
			"Close a nil connection",
			true,
			ErrConnNotInitialized,
		},
		{
			"Close an established connection",
			false,
			ErrConnNotInitialized,
		},
	}
	err := startMockTCP()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer stopMockTCP()
	for _, tt := range tests {
		y := &YeeLight{}
		if !tt.wantErr {
			y.Location = mockTCP.Addr().String()
			y.Open()
		}
		t.Run(tt.name, func(t *testing.T) {
			err := y.Close()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Close() error: expected no errors, got %+v", err)
					return
				}
				if tt.errType == nil {
					return
				}
				if errors.Cause(err) != tt.errType {
					t.Errorf("Close() error: expected %v, got %+v", tt.errType, err)
					return
				}
			}
		})
	}
}
