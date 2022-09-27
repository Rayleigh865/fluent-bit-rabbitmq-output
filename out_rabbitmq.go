package main

import (
	"C"
	"encoding/json"
	"log"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/streadway/amqp"
)

var (
	connection   *amqp.Connection
	channel      *amqp.Channel
	exchangeName string
)

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	return output.FLBPluginRegister(def, "rabbitmq", "Output to RabbitMQ")
}

//export FLBPluginInit
func FLBPluginInit(plugin unsafe.Pointer) int {
	// Gets called only once for each instance you have configured.
	var err error

	host := output.FLBPluginConfigKey(plugin, "RabbitHost")
	port := output.FLBPluginConfigKey(plugin, "RabbitPort")
	user := output.FLBPluginConfigKey(plugin, "RabbitUser")
	password := output.FLBPluginConfigKey(plugin, "RabbitPassword")
	exchangeName = output.FLBPluginConfigKey(plugin, "ExchangeName")
	exchangeType := output.FLBPluginConfigKey(plugin, "ExchangeType")

	connection, err = amqp.Dial("amqp://" + user + ":" + password + "@" + host + ":" + port + "/")
	if err != nil {
		LogError("Failed to establish a connection to RabbitMQ: ", err)
		return output.FLB_ERROR
	}

	channel, err = connection.Channel()
	if err != nil {
		LogError("Failed to open a channel: ", err)
		connection.Close()
		return output.FLB_ERROR
	}

	LogInfo("Established successfully a connection to the RabbitMQ-Server")

	err = channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		LogError("Failed to declare an exchange: ", err)
		connection.Close()
		return output.FLB_ERROR
	}

	return output.FLB_OK
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {
	// Gets called with a batch of records to be written to an instance.
	// Create Fluent Bit decoder
	dec := output.NewDecoder(data, int(length))

	// Iterate Records
	for {
		// Extract Record
		ret, _, record := output.GetRecord(dec)
		if ret != 0 {
			break
		}

		parsedRecord := ParseRecord(record)
		jsonString, err := json.Marshal(parsedRecord)
		if err != nil {
			LogError("Couldn't parse record: ", err)
			continue
		}

		err = channel.Publish(
			exchangeName, // exchange
			"",           // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        jsonString,
			})
		if err != nil {
			LogError("Couldn't publish record: ", err)
		}
	}
	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	return output.FLB_OK
}

func ParseRecord(mapInterface map[interface{}]interface{}) map[string]interface{} {
	parsedMap := make(map[string]interface{})
	for k, v := range mapInterface {
		switch t := v.(type) {
		case []byte:
			// prevent encoding to base64
			parsedMap[k.(string)] = string(t)
		case map[interface{}]interface{}:
			parsedMap[k.(string)] = ParseRecord(t)
		case []interface{}:
			parsedMap[k.(string)] = parseSubRecordArray(t)
		default:
			parsedMap[k.(string)] = v
		}
	}
	return parsedMap
}

func parseSubRecordArray(arr []interface{}) *[]interface{} {
	for idx, i := range arr {
		switch t := i.(type) {
		case []byte:
			arr[idx] = string(t)
		case map[interface{}]interface{}:
			arr[idx] = ParseRecord(t)
		case []interface{}:
			arr[idx] = parseSubRecordArray(t)
		default:
			arr[idx] = t
		}
	}
	return &arr
}

func LogInfo(msg string) {
	log.Printf("%s", msg)
}

func LogError(msg string, err error) {
	log.Printf("%s: %s", msg, err)
}

func main() {
}
