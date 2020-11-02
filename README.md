Meteo station for remore *temperature* and *humidity* measurements.

# Arduino

Sketch is based on NodeMCU with wi-fi onboard and Adafruit BME280.

Board official page: [NodeMcu - an open-source firmware and development kit](http://www.nodemcu.com/index_en.html).

## Install

For the NodeMCU Amica board [download driver CP2102](https://www.silabs.com/products/development-tools/software/usb-to-uart-bridge-vcp-drivers) firstly.

Then dowload libraries by calling `Sketch` > `Include Library` > `Manage Libraries` and searching for:

* Adafruit Unified Sensor.

* Adafruit_BME280.

## Extention points

* NodeMCU has [deep-sleep feature](https://randomnerdtutorials.com/esp8266-deep-sleep-with-arduino-ide/) for the energy saving.

## References

* [Getting started with board](https://create.arduino.cc/projecthub/electropeak/getting-started-w-nodemcu-esp8266-on-arduino-ide-28184f).

# Backend

Backend is implemented with Golang. By default application is starting at 8081 port.

Sample request:

`curl -X POST http://localhost:8081 -d "{\"token\":\"group-token\",\"name\":\"agent-1\",\"data\":[{\"o\":0,\"t\":25,\"h\":40}]}"`