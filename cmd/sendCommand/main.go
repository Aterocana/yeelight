package main

import (
	"flag"
	"fmt"
	"log"
	"time"
	"yeelight"
)

const helpCommand = `Available commands:
	toggle		- Toggle the device state (no parameters required)
`

func main() {
	var ipAddr string
	var cmd string

	flag.StringVar(&ipAddr, "ip", "192.168.0.20", "specify the IP address on local network of your YeeLight device you'd like to send a command to")
	flag.StringVar(&cmd, "cmd", "toggle", fmt.Sprintf("the command you'd like to send to your YeeLight device\n%s", helpCommand))
	flag.Parse()

	y := &yeelight.YeeLight{
		Location: fmt.Sprintf("%s:55443", ipAddr),
	}
	err := y.Open()
	if err != nil {
		log.Fatalf("Could not reach the device: %v", err)
	}
	defer y.Close()
	errorHandler(err)
	notifications := y.GetNotification()
	errs := y.GetErrors()
	go listenNotifications(notifications)
	go listenErrors(errs)

	var a *yeelight.Answer

	switch cmd {
	case "toggle":
		a, err = y.Toggle()
	default:
		log.Fatalf("command %s not implemented", cmd)
	}

	<-time.After(2 * time.Second)

	// a, err = y.SendRGB(0xff, 0xff, 0xff, yeelight.Effect("smooth"), 500)
	// errorHandler(err)
	// fmt.Println(a)

	// time.Sleep(time.Second)
	a, err = y.SetPower(yeelight.On, yeelight.Smooth, 500, yeelight.RGBMode)
	errorHandler(err)
	fmt.Println(a)

	// for {
	// 	a, err = y.Toggle()
	// 	errorHandler(err)
	// 	fmt.Println("answer:", a)
	// 	time.Sleep(1 * time.Second)

	// 	a, err = y.SetCTAbs(1700, yeelight.Smooth, 500)
	// 	errorHandler(err)
	// 	fmt.Println("answer:", a)
	// 	time.Sleep(1 * time.Second)
	// }
}

func errorHandler(err error) {
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}

func listenNotifications(c <-chan yeelight.Notification) {
	for n := range c {
		log.Printf("Notification: %v\n", n)
	}
}

func listenErrors(c <-chan error) {
	for err := range c {
		log.Printf("Error: %+v\n", err)
	}
}
