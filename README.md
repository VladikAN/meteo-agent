Meteo station for remote *temperature* and *humidity* measurements.

![version - 1 - deployment](/pics/deployment.png)

Demo results are observed by grafana:

![version - 1 - grafana](/pics/grafana.png)

# Arduino

Sketch is based on NodeMCU with wi-fi onboard and Adafruit BME280.

Board official page: [NodeMcu - an open-source firmware and development kit](http://www.nodemcu.com/index_en.html).

## Install

For the NodeMCU Amica board install [board support tools](https://create.arduino.cc/projecthub/electropeak/getting-started-w-nodemcu-esp8266-on-arduino-ide-28184f).

Then download libraries by calling `Sketch` > `Include Library` > `Manage Libraries` and searching for:

* Adafruit Unified Sensor.

* Adafruit_BME280.

### Compile

Update `#define` variables to desired host settings and wi-fi credentials.

* **AGENT_HOST** and **AGENT_PORT** is a target host and port to send JSON measurements.

* **AGENT_TOKEN** is a unique device group name. Basically any string variable, like `MyHome`, `Garage` or else.

* **AGENT_NAME** is a unique arduino board name. Basically any string value to identify masurements spot, like `Kitchen`, `Bathroom` or else.

* **WIFI_SSID** is a WIFI network name to connect.

* **WIFI_PASS** is a WIFI network password.

## To do

* NodeMCU has [deep-sleep feature](https://randomnerdtutorials.com/esp8266-deep-sleep-with-arduino-ide/) for the energy saving.

## References

* [Getting started with board](https://create.arduino.cc/projecthub/electropeak/getting-started-w-nodemcu-esp8266-on-arduino-ide-28184f).

# Backend

Backend is implemented with Golang. By default application is starting at 8081 port.

Sample request for data save:

`curl -X POST http://localhost:8081 -d "{\"token\":\"group-token\",\"name\":\"agent-1\",\"data\":[{\"o\":0,\"t\":25,\"h\":40}]}"`