#include <ESP8266WiFi.h>
#include <WiFiClientSecure.h>

#define SENSORS_TIMEOUT   1 /* Sensors read interval in minutes */
#define SENSORS_AVG_TIMES 5 /* Number of measures to make AVG value in SENSORS_TIMEOUT */

#define WIFI_SSID     "Secret WIFI name"
#define WIFI_PASS     "Secret WIFI password"
#define WIFI_TIMEOUT  1

const char* ssid = WIFI_SSID;
const char* password = WIFI_PASS;

void setup() {
  Serial.begin(9600);
  Serial.println("Starting agent ...");

  /* Init pins and modules */
  pinMode(LED_BUILTIN, OUTPUT);  
  digitalWrite(LED_BUILTIN, LOW);

  Serial.println("Agent started.");
}

void loop() {
  
}

void readSensors() {
  
}

void sendData() {
  /* Connect to WI-FI */
  WiFi.mode(WIFI_STA);

  /* Connect to WI-FI */
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    /* Blink LED while it's connecting */
    digitalWrite(LED_BUILTIN, HIGH);
    delay(100);
    digitalWrite(LED_BUILTIN, LOW);
    delay(100);
  }

  /* WI-FI connected, print connection details */
  Serial.println("WIFI connected to " + ssid);
  Serial.println("DHCP gives " + WiFi.localIP());

  /* Send data */

  /* Close connection */
  WiFi.stop();
}
