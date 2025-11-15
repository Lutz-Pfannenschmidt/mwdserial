package mwdserial

import (
	"go.bug.st/serial"
	"periph.io/x/conn/v3/gpio"
)

type StartBit int

const (
	StartBitLow  StartBit = 0
	StartBitHigh StartBit = 1
	NoStartBit   StartBit = -1
)

// Mode describes a serial port configuration.
type Mode struct {
	TxPin    gpio.PinIO // Pin used for TX (bitbanged). Has to be set for bitbanged writing.
	StartBit StartBit

	SleepState gpio.Level // Pin state to use when sleeping (idle). Typically gpio.High for standard serial communication.
	InvertBits bool       // If true, bits will be inverted when bitbanging

	// Serial port configuration

	BaudRate          int                     // The serial port bitrate (aka Baudrate), also used for bitbanging
	Parity            serial.Parity           // Parity (see Parity type for more info), also used for bitbanging
	StopBits          serial.StopBits         // Stop bits (see StopBits type for more info), use StopBitsBitBanged for bitbanged writing
	InitialStatusBits *serial.ModemOutputBits // Initial output modem bits status (if nil defaults to DTR=true and RTS=true), ignored for bitbanged writing
}

func (m Mode) mode() *serial.Mode {

	return &serial.Mode{
		BaudRate:          m.BaudRate,
		DataBits:          8,
		Parity:            m.Parity,
		StopBits:          m.StopBits,
		InitialStatusBits: m.InitialStatusBits,
	}
}
