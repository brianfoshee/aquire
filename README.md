## Turn your Raspberry Pi into a scientific grade hydroponic/aquaponic monitoring system!

### Hardware
- Raspberry Pi A+
- Raspberry Pi Supported Wifi Card
- Printed Circuit Board
  - BNC Right Angle PCB Mount
  - 100k SMD Resistor
  - 2x20 Pin 2.54mm Double Row Female Pin Header
  - 1x3 pins 2.54mm Female Header Straight (X2)
  - Atlas PH Circuit (Set to i2c mode)
  - Atlas EC Circuit (Set to i2c mode)
  - PH Probe (Atlas or another that uses BNC connector)
  - EC Probe (Atlas or another that uses BNC connector)
  - DS18B20 Waterproof Digital Temperature Sensor

### Firmware
- Written in Go 
- Communicates with Atlas Sceintic chips via I2C and DS18B20 via 1-wire
- Readings are transmitted via UDP to Statsd/Graphite for data storage and reporting

### Weatherproof Case (3D Printable)
- http://www.thingiverse.com/thing:1463063

### Viewing your data
- Once your device is online, send Brian or Nathan a private message and we'll provide you with a link to view your data
