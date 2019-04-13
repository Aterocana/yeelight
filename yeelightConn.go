package yeelight

import (
	"bufio"
	"encoding/json"
	"net"

	"github.com/pkg/errors"
)

// Close closes the TCP connection to the yeelight device
func (y *YeeLight) Close() error {
	y.connMutex.Lock()
	defer y.connMutex.Unlock()
	if y.tcpSocket != nil {
		return y.tcpSocket.Close()
	}
	return errors.WithStack(ErrConnNotInitialized)
}

// Open opens the TCP connection to the yeelight device
func (y *YeeLight) Open() error {
	if y.errs == nil {
		y.errs = make(chan error)
	}
	if y.events == nil {
		y.events = make(chan Notification)
	}
	y.Close()
	var err error
	y.connMutex.Lock()
	y.tcpSocket, err = net.Dial("tcp", y.Location)
	y.connMutex.Unlock()
	go y.readTCP()
	return errors.Wrap(err, "couldn't open TCP connection")
}

// readTCP is a loop which listens for TCP messages.
// If a generic event (state change) arrives, it is signaled
// through event chan.
// If it's a command answer, the answer is sent back to the
// command caller (if known, otherwise it is assumed comes from
// another application and so forwarded to error chan and discarded).
func (y *YeeLight) readTCP() {
	for y.tcpSocket != nil {
		msg, err := bufio.NewReader(y.tcpSocket).ReadBytes('\n')
		if err != nil {
			y.errs <- errors.Wrapf(err, "failed to parsing msg from yeelight %s", y.Location)
			continue
		}
		go func(msg []byte) {
			var a Answer
			err = json.Unmarshal(msg, &a)
			if err != nil {
				y.errs <- errors.Wrapf(err, "failed to parsing msg from yeelight %s", y.Location)
				return
			}
			if a.ID == 0 {
				go func(msg []byte) {
					ns := y.parseNotifications(msg)
					if len(ns) == 0 {
						y.errs <- errors.Wrap(ErrUnknownCommand, "empty notification")
					}
					for _, n := range ns {
						y.events <- n
					}
				}(msg)
				return
			}
			if a.ID < 0 {
				y.errs <- errors.Wrapf(ErrFailedCmd, "yeelight %s: command failed", y.Location)
				return
			}

			y.releaseAnswerChan(a.ID, &a)
		}(msg)
	}
	y.errs <- errors.WithStack(ErrConnDrop)
	y.Close()
	y.Open()
}
