package piicodev

import (
	"fmt"
	"testing"
	"time"
)

const (
	EnableTestQwiicPIR = true
	EnableTestAHT10    = true
	EnableTestLM75A    = true
)

func TestQwiicPIR(t *testing.T) {
	if EnableTestQwiicPIR {
		var err error

		var p *QwiicPIR
		if p, err = NewQwiicPIR(QwiicPIRAddress, I2CBus); err != nil {
			t.Fatalf("Error while opening the QwiicPIR: %v", err)
		}
		defer p.Close()

		var deviceID byte
		if deviceID, err = p.GetDeviceID(); err != nil {
			t.Fatalf("Error reading device ID of the QwiicPIR: %v", err)
		}

		fmt.Printf("QwiicPIR device ID: 0x%x\n", deviceID)

		var firmwareVersion [2]byte
		if firmwareVersion, err = p.GetFirmwareVersion(); err != nil {
			t.Fatalf("Error reading firmware version of the QwiicPIR: %v", err)
		}

		fmt.Printf("QwiicPIR firmware version: %v\n", firmwareVersion)

		var debounceTime uint16
		if debounceTime, err = p.GetDebounceTime(); err != nil {
			t.Fatalf("Error reading debounce time of the QwiicPIR: %v", err)
		}

		fmt.Printf("QwiicPIR debounce time: %d\n", debounceTime)

		for i := 0; i < 10; i++ {
			if i > 0 {
				time.Sleep(200 * time.Millisecond)
			}

			var detected bool
			if detected, err = p.GetRawReading(); err != nil {
				t.Fatalf("Error reading raw detection status of the QwiicPIR: %v", err)
			}

			var objDetected bool
			if objDetected, err = p.IsObjectDetected(); err != nil {
				t.Fatalf("Error reading object detected status of the QwiicPIR: %v", err)
			}

			var objRemoved bool
			if objRemoved, err = p.IstObjectRemoved(); err != nil {
				t.Fatalf("Error reading object removed status of the QwiicPIR: %v", err)
			}

			var available bool
			if available, err = p.IsAvailable(); err != nil {
				t.Fatalf("Error reading available status of the QwiicPIR: %v", err)
			}

			if err = p.ClearEventBits(); err != nil {
				t.Fatalf("Error clearing the event status bits of the QwiicPIR: %v", err)
			}

			fmt.Printf("Raw: %t, Available: %t, Detected: %t, Removed: %t\n", detected, available, objDetected, objRemoved)
		}

		for i := 0; i < 10; i++ {
			if i > 0 {
				time.Sleep(200 * time.Millisecond)
			}
			var objDetected, objRemoved, available bool

			if available, objDetected, objRemoved, err = p.GetDebounceEvents(); err != nil {
				t.Fatalf("Error reading event status of the QwiicPIR: %v", err)
			}

			if available {
				if objDetected {
					fmt.Println("Detected")
				}

				if objRemoved {
					fmt.Println("Removed")
				}
			}
		}
	}
}

func TestAHT10(t *testing.T) {
	if EnableTestAHT10 {
		var err error

		var s *AHT10
		if s, err = NewAHT10(AHT10Address, I2CBus); err != nil {
			t.Fatalf("Error while opening the AHT10: %v", err)
		}

		for i := 0; i < 10; i++ {
			var humidity, temperature float64
			if temperature, humidity, err = s.ReadSensor(); err != nil {
				t.Fatalf("Error reading sensor valyes from a AHT10: %v", err)
			}

			fmt.Printf("%.4f %%, %.4f C\n", humidity, temperature)
			time.Sleep(100 * time.Millisecond)
		}

		defer s.Close()
	}
}

func TestLM75A(t *testing.T) {
	if EnableTestLM75A {
		var err error

		var s *LM75A
		if s, err = NewLM75A(LM75AAddress, I2CBus); err != nil {
			t.Fatalf("Error while opening the LM75A: %v", err)
		}

		for i := 0; i < 10; i++ {
			var temperature float64
			if temperature, err = s.ReadTemperature(); err != nil {
				t.Fatalf("Error reading sensor valyes from a LM75A: %v", err)
			}

			fmt.Printf("%.4f C\n", temperature)
			time.Sleep(100 * time.Millisecond)
		}

		defer s.Close()
	}
}
