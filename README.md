# go-piicodev
A Go implementation of the Core Electronics PiicoDev drivers

The company, Core Electronics (https://core-electronics.com.au/), produce a set of products called PiicoDev.  They supply a set of drivers written in Python (https://github.com/CoreElectronics).  This is a port of those drivers to Go for the Raspberry Pi.

Currently supported are the Core Electronics PiicoDev:
* Pressure Seneor MS5637
* Temperature Sensor TMP117
* Ambient Light Sensor VEML6030
* Distance Sensor VL53L1X

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

