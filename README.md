# go-piicodev

A Go implementation of the Core Electronics PiicoDev drivers

The company, Core Electronics (https://core-electronics.com.au/), produce a set of products called PiicoDev. They supply a set of drivers written in Python (https://github.com/CoreElectronics). This is a port of those drivers to Go for the Raspberry Pi.

Currently supported are the Core Electronics PiicoDev:

- Pressure Seneor MS5637
- Temperature Sensor TMP117
- Ambient Light Sensor VEML6030
- Colour Sensor VEML6040
- Distance Sensor VL53L1X
- Motion Sensor MPU-6050
- Capacitive Touch Sensor CAP1203
- Air Quality Sensor ENS160
- 3 x RGB LED
- Buzzer
- Potentiometer
- Switch

Now adding other I2C devices

- Qwiic PIR Sensor
- AHT10 Temperature and humidity sensor
- LM75a Temperature sensor and thermal watchdog

For example:

```
package main

import (
	"fmt"
	"github.com/drtimf/go-piicodev"
)

func main() {
	var err error
	var t *piicodev.TMP117

	if t, err = piicodev.NewTMP117(piicodev.TMP117Address, 1); err != nil {
		fmt.Println(err)
		return
	}

	defer t.Close()

	var tempC float64
	if tempC, err = t.ReadTempC(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Current temperature:", tempC)
}
```
