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

	var isPowerSave bool
	if isPowerSave, err = light.GetPowerSave(); err != nil {
		t.Fatalf("Failed to read power save mode for VEML6030: %v", err)
	}

	fmt.Println("Ambient light sensor in power save mode:", isPowerSave)

	if err = light.SetGain(VEML6030GainOne); err != nil {
		t.Fatalf("Failed to set gain for VEML6030: %v", err)
	}

	if err = light.SetIntegrationTime(VEML6030IntegrationTime100); err != nil {
		t.Fatalf("Failed to set integration time for VEML6030: %v", err)
	}

	var gain float64
	if gain, err = light.GetGainValue(); err != nil {
		t.Fatalf("Failed to read gain for VEML6030: %v", err)
	}

	fmt.Println("Ambient light sensor gain:", gain)

	var integTime uint16
	if integTime, err = light.GetIntegrationTimeValue(); err != nil {
		t.Fatalf("Failed to read integration time for VEML6030: %v", err)
	}

	fmt.Println("Ambient light integration time(ms):", integTime)

	var lightLux float64
	if lightLux, err = light.Read(); err != nil {
		t.Fatalf("Failed to read ambient light from VEML6030: %v", err)
	}

	fmt.Println("Ambient light (lux):", lightLux)
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

func TestMotion(t *testing.T) {
	var motion *MPU6050
	var err error

	if motion, err = NewMPU6050(MPU6050Address, I2CBus); err != nil {
		t.Fatalf("Error while opening the MPU6050: %v", err)
	}
	defer motion.Close()

	var tempC float64
	if tempC, err = motion.ReadTemperature(); err != nil {
		t.Fatalf("Error reading MPU6050 temperature: %v", err)
	}

	fmt.Println("Motion sensor temperature:", tempC)

	if motion.SetAccelRange(MPU6050AccelRange2G); err != nil {
		t.Fatalf("Error setting MPU6050 accelerometer range: %v", err)
	}

	var accelRange int
	if accelRange, err = motion.GetAccelRangeValue(); err != nil {
		t.Fatalf("Error reading MPU6050 accelerometer range: %v", err)
	}

	fmt.Println("Motion sensor accelerometer range:", accelRange)

	if motion.SetGyroRange(MPU6050GyroRange250Deg); err != nil {
		t.Fatalf("Error setting MPU6050 gyro range: %v", err)
	}

	var gyroRange int
	if gyroRange, err = motion.GetGyroRangeValue(); err != nil {
		t.Fatalf("Error reading MPU6050 gyro range: %v", err)
	}

	fmt.Println("Motion sensor gyro range:", gyroRange)

	var aX, aY, aZ float64
	if aX, aY, aZ, err = motion.ReadAccelData(); err != nil {
		t.Fatalf("Error reading MPU6050 accelerometer data: %v", err)
	}

	fmt.Printf("Motion sensor accelerometer data (%f,%f,%f)\n", aX, aY, aZ)

	var gX, gY, gZ float64
	if gX, gY, gZ, err = motion.ReadGyroData(); err != nil {
		t.Fatalf("Error reading MPU6050 gyro data: %v", err)
	}

	fmt.Printf("Motion sensor gyro data (%f,%f,%f)\n", gX, gY, gZ)
}
