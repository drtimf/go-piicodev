package piicodev

import (
	"fmt"
	"testing"
)

const (
	// The I2C bus the sensors are connected to (/dev/i2c-...)
	I2CBus = 1
)

func TestTemperature(t *testing.T) {
	var temp *TMP117
	var err error
	if temp, err = NewTMP117(TMP117Address, I2CBus); err != nil {
		t.Fatalf("Error while opening the TMP117: %v", err)
	}
	defer temp.Close()

	var tempC float64
	if tempC, err = temp.ReadTempC(); err != nil {
		t.Fatalf("Failed to read temperature from TMP117: %v", err)
	}

	fmt.Println("Current temperature:", tempC)
}


func TestPressure(t *testing.T) {
	var pressure *MS5637
	var err error

	if pressure, err = NewMS5637(MS5637Address, I2CBus); err != nil {
		t.Fatalf("Error while opening the MS5637: %v", err)
	}
	defer pressure.Close()

	var pressureHpa, temperature float64
	if pressureHpa, temperature, err = pressure.Read(); err != nil {
		t.Fatalf("Failed to read pressure and temperature from MS5637: %v", err)
	}

	fmt.Printf("Current pressure: %.2f (%.2f)\n", pressureHpa, temperature)
}


func TestAmbientLight(t *testing.T) {
	var light *VEML6030
	var err error

	if light, err = NewVEML6030(VEML6030Address, I2CBus); err != nil {
		t.Fatalf("Error while opening the VEML6030: %v", err)
	}
	defer light.Close()

	var lightLux float64
	if lightLux, err = light.Read(); err != nil {
		t.Fatalf("Failed to read ambient light from VEML6030: %v", err)
	}

	fmt.Println("Ambient light:", lightLux)
}


func TestDistance(t *testing.T) {
	var dist *VL53L1X
	var err error

	if dist, err = NewVL53L1X(VL53L1XAddress, I2CBus); err != nil {
		t.Fatalf("Error while opening the VL53L1X: %v", err)
	}
	defer dist.Close()

	var rng uint16
	if rng, err = dist.Read(); err != nil {
		t.Fatalf("Failed to read range from VL53L1X: %v", err)
	}

	fmt.Println("Current range:", rng)
}
