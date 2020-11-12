#include "Arduino.h"
#include "Ticker.h"
#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <WiFiClientSecure.h>
#include <Adafruit_Sensor.h>
#include <Adafruit_BME280.h>

#define AGENT_HOST    "http://192.168.0.1071"     /* REQUIRED. Target host to send data */
#define AGENT_PORT    443                         /* REQUIRED. Target port to send data */
#define AGENT_TOKEN   "DEVELOP"                   /* REQUIRED. Host token (data group) to save records */
#define AGENT_NAME    "Agent-1"                   /* REQUIRED. Unique arduino board name to identify this agent */

#define WIFI_SSID     "WIFI_NAME" /* REQUIRED. WIFI hotspot name */
#define WIFI_PASS     "WIFI_PASS" /* REQUIRED. WIFI hotspot password */
#define WIFI_TIMEOUT  10          /* Number of seconds for connection establish */

#define SENSORS_SLEEP 60          /* Sensors read interval in seconds */
#define BUFFER        10          /* Number of items to collect before send */

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
      }

      /* Reset position */
      clear();
      return "{\"token\":\"" + String(AGENT_TOKEN) + "\","
             "\"name\":\""+ String(AGENT_NAME) + "\","
             "\"data\":[" + result + "]}";
    }

    void clear() {
      for (int i = BUFFER - 1; i >= 0; i--) {
        if (!data[i]) continue;
        delete data[i];
      }

      position = 0;
    }

    void add(Sensor *sensor) {
      if (position == BUFFER) {
        return;
      }

      data[position] = sensor;
      position++;
    }

    bool isFull() {
      return position == BUFFER;
    }

  private:
    Sensor *data[BUFFER];
    int position;
};

const char* ssid = WIFI_SSID;
const char* password = WIFI_PASS;
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
  WiFi.setAutoReconnect(true);
  WiFi.mode(WIFI_STA);
  WiFi.disconnect();
  
  /* Setup Adafruit BME280 */
  bme.begin(&Wire);
  bme.setSampling(
    Adafruit_BME280::MODE_FORCED,
    Adafruit_BME280::SAMPLING_X1,   // temperature
    Adafruit_BME280::SAMPLING_NONE, // no pressure
    Adafruit_BME280::SAMPLING_X1,   // humidity
    Adafruit_BME280::FILTER_OFF);

  Serial.println(F("Agent started."));
}

int now = 0;
void loop() {
  if (millis() - now <= SENSORS_SLEEP * 1000) {
    return;
  }

  // Can't use Ticker because NodeMCU has internal timers
  readSensors();
  now = millis();
}

void readSensors() {
  digitalWrite(LED_BUILTIN, LOW); /* Single LED flash for read sensors operation */

  bme.takeForcedMeasurement();
  Sensor *current = new Sensor(bme.readTemperature(), bme.readHumidity());
  measures.add(current);

  Serial.println("Sensor values,"
                 " t: '" + String(current->temperature) + " C'"
                 " , h: '" + String(current->humidity) + " %'");

  if (measures.isFull()) {
    sendData();
  } else {
    delay(50);
  }
  
  digitalWrite(LED_BUILTIN, HIGH);
}

void sendData() {
  if (WiFi.status() != WL_CONNECTED) {
    /* Connect to WI-FI */
    Serial.println(F("Starting WIFI connection ..."));
    WiFi.begin(ssid, password);

    int timeLeft = WIFI_TIMEOUT * 1000;
    while (WiFi.status() != WL_CONNECTED && timeLeft >= 0) {
      /* Blink LED while it's connecting */
      digitalWrite(LED_BUILTIN, LOW);
      delay(100);
      digitalWrite(LED_BUILTIN, HIGH);
      delay(100);

      timeLeft -= 200;
    }

    if (WiFi.status() != WL_CONNECTED) {
      WiFi.disconnect();
      measures.clear();
      Serial.println(F("Failed to establish WIFI connection"));
      return;
    } else {
      /* WI-FI connected, print connection details */
      Serial.print(F("WIFI connected to "));
      Serial.println(ssid);
      Serial.println("DHCP gives " + WiFi.localIP());
    }
  }

  /* Send data */
  Serial.print(F("Sending data to "));
  Serial.println(String(AGENT_HOST));

  String postData = measures.toJsonAndClear();

  WiFiClientSecure client;
  client.setInsecure(); // ignore HTTPS fingerprint, insecure way to connect
  client.connect(AGENT_HOST, AGENT_PORT);
  
  HTTPClient http;
  http.begin(client, AGENT_HOST); 
  http.addHeader("Content-Type", "application/json"); 
  auto httpCode = http.POST(postData); 
  String payload = http.getString();
  Serial.print(F("Server responded with "));
  Serial.print(httpCode);
  Serial.println(payload);
  http.end();
  
  Serial.println(F("Data sent."));
}
