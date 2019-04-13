package yeelight

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"golang.org/x/net/ipv4"

	"github.com/pkg/errors"
)

var (
	// discoveryStarted is locked when a discovery is already started, unlock otherwise.
	discoveryStarted sync.Mutex
)

// DiscoveryService is a service which listen to UDP Multi-cast group for devices discovery.
type DiscoveryService interface {
	// DiscoveryRequest should perform a discovery/search request to find YeeLight devices.
	DiscoveryRequest() error

	// GetDiscoveredDevices returns a chan when a YeeLight pointer is sent when a
	// YeeLight device sends an advertisement or answers to a discovery request
	GetDiscoveredDevices() <-chan *YeeLight

	// GetErrors returns a chan where errors during discovery are sent.
	GetErrors() <-chan error

	Open() error
}

// discoveryService is a DiscoveryService implementation
type discoveryService struct {
	// joinedMulticast is unlocked when DiscoveryService is listen on UDP Multi-cast group.
	joinedMulticast sync.Mutex

	// udpConn is the Multicast UDP packet connection.
	udpConn net.PacketConn

	// discoveredDevices is the chan where YeeLight discovered devices are notified.
	discoveredDevices chan *YeeLight

	errorsChan chan error
}

// NewDiscoveryService instantiate a DiscoveryService,
func NewDiscoveryService() DiscoveryService {
	service := discoveryService{
		discoveredDevices: make(chan *YeeLight),
		errorsChan:        make(chan error),
	}
	service.joinedMulticast.Lock()
	return &service
}

func getMyIPs() ([]string, error) {
	result := make([]string, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return result, errors.WithStack(err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return result, errors.WithStack(err)
		}
		for _, addr := range addrs {
			result = append(result, strings.Split(addr.String(), "/")[0])
		}
	}
	return result, nil
}

func addrIsIn(addr net.Addr, ips []string) bool {
	for _, ip := range ips {
		addrString := strings.Split(addr.String(), ":")[0]
		if addrString == ip {
			return true
		}
	}
	return false
}

func (service *discoveryService) Open() error {
	ssdp, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", udpAddress, udpPort))
	if err != nil {
		return errors.WithStack(err)
	}
	go func(ssdp *net.UDPAddr) {
		service.udpConn, err = net.ListenPacket("udp4", "0.0.0.0:1982")
		if err != nil {
			service.errorsChan <- errors.WithStack(err)
			return
		}
		defer service.udpConn.Close()

		service.joinedMulticast.Unlock()
		multicastConn := ipv4.NewPacketConn(service.udpConn)
		if err := errors.WithStack(multicastConn.JoinGroup(nil, ssdp)); err != nil {
			service.errorsChan <- err
			return
		}
		if err := errors.WithStack(multicastConn.SetControlMessage(ipv4.FlagDst, true)); err != nil {
			service.errorsChan <- err
			return
		}
		defer multicastConn.Close()

		myIPs, err := getMyIPs()
		if err != nil {
			service.errorsChan <- errors.WithStack(err)
			return
		}

		for {
			buf := make([]byte, 2048)
			n, _, addr, err := multicastConn.ReadFrom(buf)
			if err != nil {
				service.errorsChan <- err
				return
			}
			go func(length int, yeelightAddr net.Addr, rawMsg []byte) {
				// check if source address is my IP
				if addrIsIn(yeelightAddr, myIPs) {
					return
				}
				buf[n] = '\r'
				buf[n+1] = '\n'
				n = n + 2

				y, err := newFromAdvertisement(buf[:n])
				if err != nil {
					service.errorsChan <- errors.WithStack(err)
					return
				}
				service.discoveredDevices <- y
			}(n, addr, buf)
		}
	}(ssdp)
	return nil
}

// GetDiscoveredDevices returns a chan where every device discovered is sent.
func (service *discoveryService) GetDiscoveredDevices() <-chan *YeeLight {
	return service.discoveredDevices
}

// GetErrors returns a chan where every raised error is sent.
func (service *discoveryService) GetErrors() <-chan error {
	return service.errorsChan
}

// DiscoveryRequest sent a discovery request on UDP multi-cast group.
func (service *discoveryService) DiscoveryRequest() error {
	destAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", udpAddress, udpPort))
	if err != nil {
		return errors.WithStack(err)
	}
	go func(service *discoveryService, destAddr *net.UDPAddr) {
		service.joinedMulticast.Lock()
		if service.udpConn == nil {
			service.errorsChan <- errors.WithStack(ErrConnNotInitialized)
		}

		n, err := service.udpConn.WriteTo(searchMessage, destAddr)
		if err != nil {
			service.errorsChan <- errors.WithStack(err)
		}
		if n < len(searchMessage) {
			service.errorsChan <- errors.WithStack(ErrPartialDiscovery)
		}
	}(service, destAddr)

	return nil
}
