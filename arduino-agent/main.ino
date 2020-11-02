#include "Arduino.h"
#include "Ticker.h"
#include <ESP8266WiFi.h>
#include <WiFiClientSecure.h>
#include <Adafruit_Sensor.h>
#include <Adafruit_BME280.h>

#define SENSORS_SLEEP 2           /* Sensors read interval in seconds */
#define WIFI_SSID     "WIFI_SSID" /* WIFI hotspot name */
#define WIFI_PASS     "WIFI_PASS" /* WIFI hotspot password */
#define WIFI_SLEEP    10          /* Timeout to collect and send data in seconds */
#define BUFFER        WIFI_SLEEP / SENSORS_SLEEP  /* Circle buffer size */

class Sensor {
  public:
    Sensor() {}
    Sensor(float t, int h) {
      temperature = t;
      humidity = h;
    }
    String toJson(int offset) {
      return "{\"o\":" + String(offset) + ","
             "\"t\":" + String(temperature) + ","
             "\"h\":" + String(humidity) + "}";
    }

    float temperature;
    int humidity;
};

class Measures {
  public:
    Measures() {}

    String toJsonAndClear() {
      int offset = 0;
      String result = "";

      /* Read values from end */
      for (int i = position - 1; i >= 0; i--) {
        if (!data[i]) continue;
        String json = data[i]->toJson(offset);
        result.concat(result.length() == 0 ? json : ("," + json));
        offset += SENSORS_SLEEP;
        delete data[i];
      }

      for (int i = BUFFER - 1; i >= position; i--) {
        if (!data[i]) continue;
        String json = data[i]->toJson(offset);
        result.concat(result.length() == 0 ? json : ("," + json));
        offset += SENSORS_SLEEP;
        delete data[i];
      }

      /* Reset position */
      position = 0;
      return "[" + result + "]";
    }

    void add(Sensor *sensor) {
      if (position == BUFFER) {
        position = 0;
      }

      data[position] = sensor;
      position++;
    }

  private:
    Sensor *data[BUFFER];
    int position;
};

const char* ssid = WIFI_SSID;
const char* password = WIFI_PASS;
Ticker sensorsTicker;
Ticker transferTicker;
Adafruit_BME280 bme;
Measures measures = Measures();

void setup() {
  Serial.begin(9600);
  Serial.println(F("Starting agent ..."));

  /* Init pins */
  pinMode(LED_BUILTIN, OUTPUT);
  digitalWrite(LED_BUILTIN, HIGH);

  /* Setup WIFI */
  WiFi.persistent(false);
  WiFi.mode(WIFI_STA);

  /* Setup Adafruit BME280 */
  bme.begin(&Wire);
  bme.setSampling(
    Adafruit_BME280::MODE_FORCED,
    Adafruit_BME280::SAMPLING_X1,   // temperature
    Adafruit_BME280::SAMPLING_NONE, // pressure
    Adafruit_BME280::SAMPLING_X1,   // humidity
    Adafruit_BME280::FILTER_OFF);

  /* Start timers */
  sensorsTicker.attach_ms(SENSORS_SLEEP * 1000, readSensors);
  transferTicker.attach_ms(WIFI_SLEEP * 1000, sendData);

  Serial.println(F("Agent started."));
}

void loop() {}

void readSensors() {
  digitalWrite(LED_BUILTIN, LOW); /* Single LED flash for read sensors operation */

  bme.takeForcedMeasurement();
  Sensor *current = new Sensor(
    bme.readTemperature(),
    bme.readHumidity());
  measures.add(current);

  Serial.println("Sensor values,"
                 " t: '" + String(current->temperature) + " C'"
                 " , h: '" + String(current->humidity) + " %'");

  delay(50);
  digitalWrite(LED_BUILTIN, HIGH);
}

void sendData() {
  if (WiFi.status() != WL_CONNECTED) {
    /* Connect to WI-FI */
    Serial.println(F("Starting WIFI connection ..."));
    WiFi.begin(ssid, password);
    while (WiFi.status() != WL_CONNECTED) {
      /* Blink LED while it's connecting */
      digitalWrite(LED_BUILTIN, LOW);
      delay(100);
      digitalWrite(LED_BUILTIN, HIGH);
      delay(100);
    }

    /* WI-FI connected, print connection details */
    Serial.print(F("WIFI connected to "));
    Serial.println(ssid);
    Serial.println("DHCP gives " + WiFi.localIP());
  }

  /* Send data */
  Serial.println(F("Sending data ..."));
  Serial.println(measures.toJsonAndClear());
  Serial.println(F("Data sent."));
}
