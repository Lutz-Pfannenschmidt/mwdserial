package mwdserial

import (
	"fmt"
	"time"

	"go.bug.st/serial"
	"periph.io/x/conn/v3/gpio"
)

type MWDSerial struct {
	port  serial.Port
	mode  Mode
	delay time.Duration // 1 / BaudRate
}

func Open(portName string, mode Mode) (*MWDSerial, error) {
	serialPort, err := serial.Open(portName, mode.mode())
	if err != nil {
		return nil, err
	}

	if mode.TxPin != nil {
		err := mode.TxPin.Out(gpio.High)
		if err != nil {
			serialPort.Close()
			return nil, fmt.Errorf("failed to set TxPin high: %w", err)
		}
	} else {
		return nil, fmt.Errorf("TxPin must be set for bitbanged writing")
	}

	if mode.SleepState != gpio.Low && mode.SleepState != gpio.High {
		serialPort.Close()
		return nil, fmt.Errorf("SleepState must be either gpio.Low or gpio.High")
	}

	mode.TxPin.Out(mode.SleepState)

	return &MWDSerial{port: serialPort, mode: mode, delay: time.Millisecond * time.Duration(1000/mode.BaudRate)}, nil
}

func (m *MWDSerial) SetMode(mode Mode) error {
	err := m.port.SetMode(mode.mode())
	if err != nil {
		return err
	}
	m.mode = mode
	return nil
}

// Read
func (m *MWDSerial) Read(p []byte) (n int, err error) {
	return m.port.Read(p)
}

// Write
func (m *MWDSerial) Write(p []byte) (n int, err error) {
	bits := make([]bool, 0, len(p)*8)
	for _, b := range p {
		for i := range 8 {
			bit := (b & (1 << i)) != 0
			bits = append(bits, bit)
		}
	}
	return m.WriteBits(bits)
}

// Drain
func (m *MWDSerial) Drain() error {
	return m.port.Drain()
}

// ResetInputBuffer
func (m *MWDSerial) ResetInputBuffer() error {
	return m.port.ResetInputBuffer()
}

// ResetOutputBuffer
func (m *MWDSerial) ResetOutputBuffer() error {
	return m.port.ResetOutputBuffer()
}

// SetDTR
func (m *MWDSerial) SetDTR(dtr bool) error {
	return m.port.SetDTR(dtr)
}

// SetRTS
func (m *MWDSerial) SetRTS(rts bool) error {
	return m.port.SetRTS(rts)
}

// GetModemStatusBits
func (m *MWDSerial) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	return m.port.GetModemStatusBits()
}

// SetReadTimeout
func (m *MWDSerial) SetReadTimeout(t time.Duration) error {
	return m.port.SetReadTimeout(t)
}

// Close
func (m *MWDSerial) Close() error {
	return m.port.Close()
}

// Break
func (m *MWDSerial) Break(duration time.Duration) error {
	return m.port.Break(duration)
}

func (m *MWDSerial) WriteBits(bits []bool) (n int, err error) {
	err = m.write(bits)
	if err != nil {
		return 0, err
	}
	return len(bits), nil
}

func invertBits(bits *[]bool) {
	for i, bit := range *bits {
		(*bits)[i] = !bit
	}
}

func (m *MWDSerial) getParityBit(bits []bool) int {
	switch m.mode.Parity {
	case serial.NoParity:
		return -1
	case serial.MarkParity:
		return 1
	case serial.SpaceParity:
		return 0
	}

	// Even or Odd Parity

	oneBitsCount := 0
	for _, bit := range bits {
		if bit {
			oneBitsCount++
		}
	}

	if m.mode.Parity == serial.EvenParity {
		if oneBitsCount%2 == 0 {
			return 0
		}
		return 1
	}

	// Odd Parity
	if oneBitsCount%2 == 0 {
		return 1
	}
	return 0
}

func (m *MWDSerial) write(bits []bool) error {
	if m.mode.InvertBits {
		invertBits(&bits)

		switch m.mode.StartBit {
		case StartBitLow:
			m.mode.TxPin.Out(gpio.High)
		case StartBitHigh:
			m.mode.TxPin.Out(gpio.Low)
		}
	} else {
		switch m.mode.StartBit {
		case StartBitLow:
			m.mode.TxPin.Out(gpio.Low)
		case StartBitHigh:
			m.mode.TxPin.Out(gpio.High)
		}
	}

	time.Sleep(m.delay)

	for _, bit := range bits {
		m.mode.TxPin.Out(gpio.Level(bit))
		time.Sleep(m.delay)
	}

	// Parity bit
	parityBit := m.getParityBit(bits)
	if parityBit != -1 {
		if m.mode.InvertBits {
			parityBit = 1 - parityBit
		}
		m.mode.TxPin.Out(gpio.Level(parityBit == 1))
		time.Sleep(m.delay)
	}

	m.mode.TxPin.Out(m.mode.SleepState)

	return nil
}
