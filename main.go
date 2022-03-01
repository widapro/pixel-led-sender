package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MainJson struct {
	MainWeather MainWeather `json:"main"`
}

type MainWeather struct {
	MainTemp float64 `json:"temp"`
}

// ENV VARIABLES
// var NAME = getEnv("VARIABLE NAME", "DEFAULT VALUE")
var broker = getEnv("MQTT_ADDRESS", "nil")                              // Address of MQTT Server
var port = getEnv("MQTT_PORT", "1883")                                  // MQTT port
var mqttUser = getEnv("MQTT_USER", "username")                          // MQTT user
var mqttPassword = getEnv("MQTT_PASSWORD", "change_me")                 // MQTT user password
var weatherToken = getEnv("WEATHER_TOKEN", "change_me")                 // API token from openweathermap.org
var weatherId = getEnv("WEATHER_CITY_ID", "4459467")                    // City ID from openweathermap.org; All existing ID's can be found here: http://bulk.openweathermap.org/sample/city.list.json.gz
var mqttTopic1 = getEnv("MQTT_TOPIC1", "wled/zone0_text")               // MQTT topic to send temperature
var mqttTopic2 = getEnv("MQTT_TOPIC2", "wled/zone1_text")               // MQTT topic to send time
var refreshTime, err = time.ParseDuration(getEnv("REFRESH_TIME", "5s")) // Refresh temperature time in seconds

// MQTT variables
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("MQTT Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("MQTT Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("MQTT Connect lost: %v", err)
}

func publish(client mqtt.Client, message string, mqttTopic string) {
	token := client.Publish(mqttTopic, 0, false, message)
	token.Wait()

}

func GetTempJson() (float64, float64, error) {
	url := "https://api.openweathermap.org/data/2.5/weather?id=" + weatherId + "&units=metric&appid=" + weatherToken

	res, err := http.Get(url)
	if err != nil {
		return 0, 0, fmt.Errorf("error: request to weather provider was failed with error: %s", err)
	}
	if res.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("error: Request to weather provider was failed, Non-OK HTTP status: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("error: response from weather provider was failed with error: %s", err)
	}

	var data MainJson
	json.Unmarshal(body, &data)

	tempCelcium := math.Round(data.MainWeather.MainTemp*10) / 10
	tempFahrenheit := (tempCelcium * 1.8) + 32

	return tempCelcium, tempFahrenheit, nil
}

func postTemp(client mqtt.Client, tempSleepTimer time.Duration) {
	for {
		tempCelcium, tempFahrenheit, err := GetTempJson()
		if err != nil {
			fmt.Println(err)
			fmt.Println("sleeping for the next iteration")
			time.Sleep(tempSleepTimer)
			continue
		}

		tempCelciumString := fmt.Sprintf("%g", tempCelcium)
		tempFahrenheitString := fmt.Sprintf("%.0f", tempFahrenheit)

		for i := 0; i < 6; i++ {
			fmt.Println("MQTT post Celcium: ", tempCelciumString)
			publish(client, tempCelciumString+" C", mqttTopic1)
			time.Sleep(tempSleepTimer)

			fmt.Println("MQTT post Fahrenheit: ", tempFahrenheitString)
			publish(client, tempFahrenheitString+" F", mqttTopic1)
			time.Sleep(tempSleepTimer)
		}
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername(mqttUser)
	opts.SetPassword(mqttPassword)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Refresh time: ", refreshTime)
	go postTemp(client, refreshTime)

	ticker := time.NewTicker(10 * time.Second)
	for _ = range ticker.C {
		timeCurrent := time.Now().Format("15:04")
		publish(client, timeCurrent, mqttTopic2)
		fmt.Println("Current time published at ", timeCurrent)
	}
}
