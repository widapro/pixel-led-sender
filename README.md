pixel-led-sender
------------------
Sender for pixel led mqtt panel

How to run
----------
The recommended way to build and run
> do not forget to set all required environment variables
```go
go build .
./pixel-led-sender
```

 Requiered environment variables
--------------------------------
```bash
export MQTT_ADDRESS="192.168.XXX.XXX"  # Requierd! Address of MQTT Server or domain name
export MQTT_PORT="1883"              # Optional. Default value 1883
export MQTT_USER="mqtt_user"         # MQTT user. Default value username
export MQTT_PASSWORD="change_me"     # Requierd! MQTT PASSWORD
export WEATHER_TOKEN="change_me"     # Requierd! API token from openweathermap.org
export WEATHER_CITY_ID="4459467"     # Requierd! City ID from openweathermap.org; All existing ID's can be found here: http://bulk.openweathermap.org/sample/city.list.json.gz
export MQTT_TOPIC1="wled/zone0_text" # Optional. MQTT topic to send temperature. Default value "wled/zone0_text"
export MQTT_TOPIC2="wled/zone1_text" # Optional. MQTT topic to send time. Default value "wled/zone1_text"
```
