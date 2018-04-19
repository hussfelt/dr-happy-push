#include <ESP8266WiFi.h>
#include <WiFiClientSecure.h>

const char* ssid     = "";
const char* password = "";
const char* host = "";
const char* url = "/Prod/";
const char* auth = ""
String currentButton;

// Declare the pins for the Buttons and the LED
int greenButtonPin = D3;
int redButtonPin = D7;
int yellowButtonPin = D5;

// Button state
int lastButtonState = HIGH;
unsigned long lastDebounceTime = 0;
unsigned long debounceDelay = 1000;

// Setup wifi client
WiFiClientSecure client;

void setup() {
  // Setup pins
  pinMode(greenButtonPin, INPUT_PULLUP);
  pinMode(redButtonPin, INPUT);
  pinMode(yellowButtonPin, INPUT);
  pinMode(LED_BUILTIN, OUTPUT);

  // Open serial, with 9600 baud
  Serial.begin(9600);

  // We start by connecting to a WiFi network
  Serial.println();
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  /* Explicitly set the ESP8266 to be a WiFi-client, otherwise, it by default,
     would try to act as both a client and an access-point and could cause
     network-issues with your other WiFi-devices on your WiFi-network. */
  WiFi.mode(WIFI_STA);

  // Begin wifi
  WiFi.begin(ssid, password);

  // While we are not connected wait 500 ms before checking again
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  // If we are connected, output to console
  Serial.println("");
  Serial.println("WiFi connected");
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());
}

void loop(){
  // Read the values of the input from our buttons It can either be 1 or 0  
  int greenButton = digitalRead(greenButtonPin);
  int redButton = digitalRead(redButtonPin);
  int yellowButton = digitalRead(yellowButtonPin);

  // Check if we're connected to wifi
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("Disconnected from wifi. Reconnecting");
    WiFi.mode(WIFI_STA);
    WiFi.begin(ssid, password);
    // Check which button is pressed
    if (greenButton == LOW) {
      currentButton = "Green";
    } else if (yellowButton == LOW) {
      currentButton = "Yellow";
    } else if (redButton == LOW) {
      currentButton = "Red";
    } else {
       // Reset bounce
       lastDebounceTime = 0;
       // Turn the LED off
       digitalWrite(LED_BUILTIN, HIGH);
       return;
    }
    button(currentButton);
  } else {
    // Output failure to terminal and delay 10 seconds
    Serial.println("Could not connect to wifi. We try to connect to ssid:" + ssid + "in 10 seconds.")
    delay(10000);
  }
}

void button(String newCurrentButton) {
  // Setup state-variables
  bool finishedHeaders = false;
  bool currentLineIsBlank = true;
  bool gotResponse = false;

  // Create a variable for storing point in time
  long now;

  // Setup strings to store metadata
  String title = "";
  String headers = "";
  String body = "";

  // Check if we are OK to run actions on button click
  if ((millis() - lastDebounceTime) > debounceDelay) {
    // reset the debouncing timer
    lastDebounceTime = millis();
    // If button pushed, turn LED on
    Serial.println(newCurrentButton + " pushed!");

    // Connect to our host
    if (client.connect(host, 443)) {
      // Output to console
      Serial.println("connected");
      Serial.println(url);

      // Write data to request
      client.println("POST " + url + " HTTP/1.1");
      client.print("Host: "); client.println(host);
      client.println("User-Agent: arduino/1.0");
      client.println("Auth: "); client.println(auth);
      client.println("Content-Type: application/x-www-form-urlencoded");
      client.print("Content-Length: "); client.println(newCurrentButton.length());
      client.println("");
      client.println(newCurrentButton);

      // Store point in time
      now = millis();

      // checking the timeout
      while (millis() - now < 1500) {
        // While request is active
        while (client.available()) {
          // Read response
          char c = client.read();
  
          // Check if response is done with headers
          if (finishedHeaders) {
            body=body+c;
          } else {
            // Else, store that we finished headers
            if (currentLineIsBlank && c == '\n') {
              finishedHeaders = true;
            } else {
              headers = headers + c;
            }
          }

          if (c == '\n') {
            currentLineIsBlank = true;
          }else if (c != '\r') {
            currentLineIsBlank = false;
          }
  
          //marking we got a response
          gotResponse = true;
  
        }

        // We got response, print to console
        if (gotResponse) {
          Serial.println(body);
          break;
        }
      }

      // Light up led, we have a push
      digitalWrite(LED_BUILTIN, LOW);
    }
  }
}