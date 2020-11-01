#include "Arduino.h"
#include "Ticker.h"
#include <ESP8266WiFi.h>
#include <WiFiClientSecure.h>
#include <Adafruit_Sensor.h>
#include <Adafruit_BME280.h>

#define SENSORS_TIMEOUT   10          /* Sensors read interval in seconds */
#define SENSORS_AVG_TIMES 5           /* Number of measures to make AVG value in SENSORS_TIMEOUT */
#define WIFI_SSID         "WIFI_SSID" /* WIFI hotspot name */
#define WIFI_PASS         "WIFI_PASS" /* WIFI hotspot password */
#define WIFI_TIMEOUT      1           /* Timeout to collect and send data in seconds */

const char* ssid = WIFI_SSID;
const char* password = WIFI_PASS;

Ticker sensorsTicker;
Ticker transferTicker;

Adafruit_BME280 bme;

class Sensor {
  public:
    Sensor(float t, int p, int h) {
      temperature = t;
      pressure = p;
      humidity = h;
    }
  
    float temperature;
    int pressure;
    int humidity;
};

void setup() {
  Serial.begin(9600);
  Serial.println("Starting agent ...");

  /* Init pins and modules */
  pinMode(LED_BUILTIN, OUTPUT);  
  digitalWrite(LED_BUILTIN, HIGH);

  /* Setup WIFI */
  WiFi.persistent(false);
  WiFi.mode(WIFI_STA);

  /* Setup Adafruit BME */
  bme.begin(&Wire);
  bme.setSampling(Adafruit_BME280::MODE_FORCED,
                  Adafruit_BME280::SAMPLING_X1, // temperature
                  Adafruit_BME280::SAMPLING_X1, // pressure
                  Adafruit_BME280::SAMPLING_X1, // humidity
                  Adafruit_BME280::FILTER_OFF);

  /* Start timers */
  sensorsTicker.attach_ms(SENSORS_TIMEOUT * 1000 / SENSORS_AVG_TIMES, readSensors);
  //transferTicker.attach_ms(1200, sendD);

  Serial.println("Agent started.");
}

void loop() {
}

void readSensors() {
  digitalWrite(LED_BUILTIN, LOW);

  bme.takeForcedMeasurement();

  Sensor current(
    bme.readTemperature(),
    bme.readPressure() / 100.0F,
    bme.readHumidity());
    
  Serial.println("Sensor values, t: '" + String(current.temperature) + " C' , h: '" + String(current.humidity) + " %' , p: '" + String(current.pressure) + " hPA'.");

  delay(50);
  digitalWrite(LED_BUILTIN, HIGH);
}

void sendData() {
  /* Connect to WI-FI */
  Serial.println("Starting WIFI connection ...");
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    /* Blink LED while it's connecting */
    digitalWrite(LED_BUILTIN, LOW);
    delay(100);
    digitalWrite(LED_BUILTIN, HIGH);
    delay(100);
  }

  /* WI-FI connected, print connection details */
  Serial.print("WIFI connected to ");
  Serial.println(ssid);
  Serial.println("DHCP gives " + WiFi.localIP());

  /* Send data */
  Serial.println("Sending data ...");

  Serial.println("Data sent.");

  /* Close connection */
  //WiFi.stop();
  Serial.println("WIFI connection closed.");
}
