package piicodev

import (
	"fmt"
	"testing"
	"time"
)

const (
	// The I2C bus the sensors are connected to (/dev/i2c-...)
	I2CBus = 1

	// Enable different tests
	EnableTestTemperature     = true
	EnableTestPressure        = true
	EnableTestAmbientLight    = true
	EnableTestDistance        = true
	EnableTestMotion          = true
	EnableTestCapacitiveTouch = true
	EnableTestRGBLED          = true
	EnableTestBuzzer          = true
	EnableTestColour          = true
	EnableAirQualitySensor    = true

	// Fun...
	EnableTestColourSensorToRGBLED = false
)

func TestTemperature(t *testing.T) {
	if EnableTestTemperature {
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
}

func TestPressure(t *testing.T) {
	if EnableTestPressure {
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
}

func TestAmbientLight(t *testing.T) {
	if EnableTestAmbientLight {
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
}

func TestDistance(t *testing.T) {
	if EnableTestDistance {
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
}

func TestMotion(t *testing.T) {
	if EnableTestMotion {
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
}

func TestCapacitiveTouch(t *testing.T) {
	if EnableTestCapacitiveTouch {
		var ct *CAP1203
		var err error

		if ct, err = NewCAP1203(CAP1203Address, I2CBus); err != nil {
			t.Fatalf("Error while opening the CAP1203: %v", err)
		}
		defer ct.Close()

		for setSensitivity := 0; setSensitivity <= 7; setSensitivity++ {
			if err = ct.SetSensitivity(setSensitivity); err != nil {
				t.Fatalf("Error setting CAP1203 touch sensitivity: %v", err)
			}

			var sensitivity int
			if sensitivity, err = ct.GetSensitivity(); err != nil {
				t.Fatalf("Error reading CAP1203 touch sensitivity: %v", err)
			}

			if setSensitivity != sensitivity {
				t.Fatalf("Error setting CAP1203 touch sensitivity and reading back: %d != %d", setSensitivity, sensitivity)
			}
		}

		if err = ct.SetSensitivity(6); err != nil {
			t.Fatalf("Error setting CAP1203 touch sensitivity: %v", err)
		}

		var sensitivity int
		if sensitivity, err = ct.GetSensitivity(); err != nil {
			t.Fatalf("Error reading CAP1203 touch sensitivity: %v", err)
		}

		fmt.Printf("Capacitive touch sensor sensitivity: %d\n", sensitivity)

		ct.SetMultipleTouchEnabled(true)

		var multipleTouch bool
		if multipleTouch, err = ct.GetMultipleTouchEnabled(); err != nil {
			t.Fatalf("Error reading CAP1203 multiple touch setting: %v", err)
		}

		fmt.Printf("Capacitive touch sensor multiple touch enabled: %t\n", multipleTouch)

		var touched bool
		if touched, err = ct.IsTouched(); err != nil {
			t.Fatalf("Error reading CAP1203 general status: %v", err)
		}
		fmt.Printf("Capacitive touch sensor was touched: %t\n", touched)

		var count1, count2, count3 int
		if count1, count2, count3, err = ct.ReadDeltaCounts(); err != nil {
			t.Fatalf("Error reading CAP1203 delta counts: %v", err)
		}

		fmt.Printf("Capacitive touch sensor delta counts: %d %d %d\n", count1, count2, count3)

		var status1, status2, status3 bool
		if status1, status2, status3, err = ct.Read(); err != nil {
			t.Fatalf("Error reading CAP1203 touch status: %v", err)
		}

		fmt.Printf("Capacitive touch sensor status: %t %t %t\n", status1, status2, status3)
	}
}

func TestRGBLED(t *testing.T) {
	if EnableTestRGBLED {
		var led *RGBLED
		var err error

		if led, err = NewRGBLED(RGBLEDAddress, I2CBus); err != nil {
			t.Fatalf("Error while opening the RGBLED: %v", err)
		}
		defer led.Close()

		/* Bug in firmware??  Cannot read anything :(
		var deviceID byte
		if deviceID, err = led.GetDeviceID(); err != nil {
			t.Fatalf("Error reading device ID of the RGBLED: %v", err)
		}

		fmt.Printf("RGB LED device ID: %x\n", deviceID)

		var firmwareVersion uint16
		if firmwareVersion, err = led.GetFirmwareVersion(); err != nil {
			t.Fatalf("Error reading firmware version of the RGBLED: %v", err)
		}

		fmt.Printf("RGB LED firmware version: %d\n", firmwareVersion)
		*/

		// Flash each LED red, green and blue in turn at three different brightness levels
		// Flash the power LED as well
		for i := 0; i < 6*3; i++ {
			if i != 0 {
				time.Sleep(150 * time.Millisecond)
			}

			switch (i / 6) % 3 {
			case 0:
				led.SetBrightness(255)
			case 1:
				led.SetBrightness(50)
			case 2:
				led.SetBrightness(5)
			}

			led.ClearPixels()
			switch (i / 2) % 3 {
			case 0:
				led.SetPixel(0, 255, 0, 0)
			case 1:
				led.SetPixel(1, 0, 255, 0)
			case 2:
				led.SetPixel(2, 0, 0, 255)
			}

			state := true
			if i%2 == 0 {
				state = false
			}

			if err = led.EnablePowerLED(state); err != nil {
				t.Fatalf("Error setting RGBLED power LED: %v", err)
			}

			if err = led.Show(); err != nil {
				t.Fatalf("Error setting RGBLED values: %v", err)
			}
		}

		if err = led.Clear(); err != nil {
			t.Fatalf("Error clearing RGBLED values: %v", err)
		}

		if err = led.EnablePowerLED(false); err != nil {
			t.Fatalf("Error setting RGBLED power LED: %v", err)
		}
	}
}

func TestBuzzer(t *testing.T) {
	if EnableTestBuzzer {
		var b *Buzzer
		var err error

		if b, err = NewBuzzer(BuzzerAddress, I2CBus); err != nil {
			t.Fatalf("Error while opening the Buzzer: %v", err)
		}
		defer b.Close()

		var deviceID byte
		if deviceID, err = b.GetDeviceID(); err != nil {
			t.Fatalf("Error reading device ID of the Buzzer: %v", err)
		}

		fmt.Printf("Buzzer device ID: 0x%x\n", deviceID)

		var firmwareVersion [2]byte
		if firmwareVersion, err = b.GetFirmwareVersion(); err != nil {
			t.Fatalf("Error reading firmware version of the Buzzer: %v", err)
		}

		fmt.Printf("Buzzer firmware version: %v\n", firmwareVersion)

		// Flash power led, alternate tones at the three different volume levels
		for i := 0; i < 6*3; i++ {
			if i != 0 {
				time.Sleep(150 * time.Millisecond)
			}

			if err = b.SetVolume(i / 6); err != nil {
				t.Fatalf("Error setting the volume of the Buzzer: %v", err)
			}

			tone := uint16(800)
			state := true
			if i%2 == 0 {
				tone = 500
				state = false
			}

			if err = b.EnablePowerLED(state); err != nil {
				t.Fatalf("Error setting the power LED of the Buzzer: %v", err)
			}

			if err = b.SetTone(tone, 0); err != nil {
				t.Fatalf("Error setting the tone of the Buzzer: %v", err)
			}
		}

		if err = b.NoTone(); err != nil {
			t.Fatalf("Error turning off the Buzzer: %v", err)
		}

		if err = b.EnablePowerLED(false); err != nil {
			t.Fatalf("Error setting the power LED of the Buzzer: %v", err)
		}
	}
}

func TestColour(t *testing.T) {
	if EnableTestColour {
		var c *VEML6040
		var err error

		if c, err = NewVEML6040(VEML6040Address, I2CBus); err != nil {
			t.Fatalf("Error while opening the VEML6040: %v", err)
		}
		defer c.Close()

		for i := 0; i < 10; i++ {
			if i > 0 {
				time.Sleep(100 * time.Millisecond)
			}

			var red, green, blue, white uint16
			if red, green, blue, white, err = c.ReadRGBW(); err != nil {
				t.Fatalf("Error reading RGBW values from VEML6040: %v", err)
			}

			cct := CalculateCCT(red, green, blue)
			hue, saturation, value := CalculateHSV(red, green, blue)
			fmt.Printf("(%d, %d, %d) %d %f (%f, %f, %f)\n", red, green, blue, white, cct, hue, saturation, value)
		}
	}
}

func TestColourSensorToRGBLED(t *testing.T) {
	if EnableTestColourSensorToRGBLED {
		var err error

		var c *VEML6040
		if c, err = NewVEML6040(VEML6040Address, I2CBus); err != nil {
			t.Fatalf("Error while opening the VEML6040: %v", err)
		}
		defer c.Close()

		var led *RGBLED
		if led, err = NewRGBLED(RGBLEDAddress, I2CBus); err != nil {
			t.Fatalf("Error while opening the RGBLED: %v", err)
		}
		defer led.Close()

		led.SetBrightness(120)

		for true {
			var red, green, blue uint16
			if red, green, blue, _, err = c.ReadRGBW(); err != nil {
				t.Fatalf("Error reading RGBW values from VEML6040: %v", err)
			}

			scale := 2.0
			rb := byte((float64(red) * 256.0 * scale) / 65535.0)
			gb := byte((float64(green) * 256.0 * scale) / 65535.0)
			bb := byte((float64(blue) * 256.0 * scale) / 65535.0)
			led.FillPixels(rb, gb, bb)
			led.Show()

			fmt.Printf("(%d, %d, %d)\n", rb, gb, bb)
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func TestEnableAirQualitySensor(t *testing.T) {
	if EnableAirQualitySensor {
		var err error

		var s *ENS160
		if s, err = NewENS160(ENS160Address, I2CBus); err != nil {
			t.Fatalf("Error while opening the ENS160: %v", err)
		}
		defer s.Close()

		if err = s.SetTemperature(22.5); err != nil {
			t.Fatalf("Error setting temperature for ENS160: %v", err)
		}

		time.Sleep(20 * time.Millisecond)

		var temperature float64
		if temperature, err = s.GetTemperature(); err != nil {
			t.Fatalf("Error reading temperature from ENS160: %v", err)
		}
		fmt.Printf("%f\n", temperature)

		if err = s.SetHumidity(40); err != nil {
			t.Fatalf("Error setting humidity for ENS160: %v", err)
		}

		time.Sleep(20 * time.Millisecond)

		var humidity float64
		if humidity, err = s.GetHumidity(); err != nil {
			t.Fatalf("Error reading humidity from ENS160: %v", err)
		}
		fmt.Printf("%f\n", humidity)

		for i := 0; i < 5; i++ {
			var status byte
			if status, err = s.GetStatus(); err != nil {
				t.Fatalf("Error reading status from ENS160: %v", err)
			}

			var operation string
			if operation, err = s.GetOperation(); err != nil {
				t.Fatalf("Error reading operation from ENS160: %v", err)
			}

			var aqi byte
			var aqiRating string
			if aqi, aqiRating, err = s.ReadAQI(); err != nil {
				t.Fatalf("Error reading AQI from ENS160: %v", err)
			}

			var tvoc uint16
			if tvoc, err = s.ReadTVOC(); err != nil {
				t.Fatalf("Error reading TVOC from ENS160: %v", err)
			}

			var eco2 uint16
			var eco2Rating string
			if eco2, eco2Rating, err = s.ReadECO2(); err != nil {
				t.Fatalf("Error reading ECO2 from ENS160: %v", err)
			}

			fmt.Printf("--------------------------------\n    flag: %x\n     AQI: %d [%s]\n    TVOC: %d\n    eCO2: %d ppm [%s]\n  Status: %s\n",
				status, aqi, aqiRating, tvoc, eco2, eco2Rating, operation)
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
