# Medium Well Done Serial

A simple Golang library for half bit banged serial communication.
(Hardware Reading, bitbanged writing)
"Wrapper" around [go.bug.st/serial](https://pkg.go.dev/go.bug.st/serial) and [periph.io/x/periph](https://periph.io/).

MWDSerial implements the serial.Port interface from go.bug.st/serial,
so it can be used as a near drop-in replacement for any serial port in Go programs.

## Disclaimer

I know that bitbanging serial communication in software is not the best idea.
Especially not in a non-realtime OS like Linux.
Especially especially not half bitbanged serial communication.

Plese don't use this library in production systems or anything critical.
It sucks.
