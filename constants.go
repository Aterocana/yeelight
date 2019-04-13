package yeelight

import (
	"fmt"
	"net"
	"time"
)

const (
	// manHeader is the MAN header value in discovery message
	manHeader = `"ssdp:discover"`

	// udpAddress is the UDP Multicast address where discovery messages are sent and received
	udpAddress = "239.255.255.250"

	// udpPort is the port where discovery messages are sent
	udpPort = 1982

	// startLine is the first line in sent discovery messages
	startLine = "M-SEARCH"
)

// searchMessage is used to send a discovery message in UDP multicast group where
// YeeLight devices listen to.
var searchMessage = []byte(fmt.Sprintf("M-SEARCH * HTTP/1.1\r\nHOST:%s:%d\r\nMAN:\"ssdp:discover\"\r\nST:wifi_bulb", udpAddress, udpPort))

// groupAddr is the UDP multicast group where YeeLight devices listen to.
var groupAddr = net.IPv4(239, 255, 255, 250)

// discoveryAnswerHeader is the header sent in discovery answers by
// YeeLight devices.
var discoveryAnswerHeader = []byte(("HTTP/1.1 200 OK\r\n"))

// advertisementHeader is the header sent periodically in advertisement
// messages by YeeLight devices.
var advertisementHeader = []byte("NOTIFY * HTTP/1.1\r\n")

// commandTimeout is the time waited for a command answer before raising an error
// and release the connection mutex.
var commandTimeout = time.Second
