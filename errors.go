package yeelight

import "errors"

// ErrPartialDiscovery is the error raised during the discovery
// when the search message sent is not fully delivered.
// A smaller chunk of bytes is sent instead the whole packet.
var ErrPartialDiscovery = errors.New("UDP Request: sent partial search message")

// ErrConnNotInitialized is the error raised when an operation is done on underlying
// UDP Multi-cast connection, which is not initialized yet,
// or on yeelight device tcp connection.
var ErrConnNotInitialized = errors.New("connection not initialized")

// ErrWrongAdvertisement is the error raised when an arrived advertisement message
// is not properly formatted.
var ErrWrongAdvertisement = errors.New("Wrong advertisemnt format")

// ErrWrongNotification is the error raised when an arrived notification message
// is not properly formatted.
// var ErrWrongNotification = errors.New("Wrong notification message")

// ErrInvalidRange is the error raised when a specified value is out of range
// for a specific property: for example power differs from "on" or "off",
// bright is not in the range 1-100, color_mode is different from 1,2 or 3...
var ErrInvalidRange = errors.New("Invalid range value")

// ErrInvalidType is the error raised when a command parameter value is
// of the wrong type.
var ErrInvalidType = errors.New("Invalid parameter type")

// ErrTimedOut is the error raised when a TCP communication doesn't arrive in time
// from YeeLight device.
var ErrTimedOut = errors.New("TCP communication timed out")

// ErrFailedCmd is the error raised when a command got an error answer
// or it couldn't be sent.
var ErrFailedCmd = errors.New("Failed command")

// ErrConnDrop is the error raised when a TCP connection is dropped.
var ErrConnDrop = errors.New("TCP communication is dropped")

// ErrUnknownCommand is the erorr raised when an answer for an external command
// (sent from another master) is received.
var ErrUnknownCommand = errors.New("Answer received for an unknown command")
