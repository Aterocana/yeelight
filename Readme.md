# YeeLight Go Library

This is a simple Go library which integrates YeeLights.
The code is working so far, but not all functions are implemented yet.
Discovery part is fully implemented, like most of the commands.

## Documentation

Documentation can be generated with `go doc`.

## Examples

In the `cmd` folder you have a couple of simple examples using the library.

* `discover` finds all device in your local network, printing their IPs. It closes after 30 seconds.
* `sendCommand` sends a command to a specified device. So fat just `toggle` is implemented. Run it with `--help` option to have a detailed description.