package bambu

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PrintState struct {
	mu       sync.Mutex
	previous string
}

type MQTTConfig struct {
	IP         string
	AccessCode string
	Serial     string
}

type reportMessage struct {
	Print *printStatus `json:"print"`
}

type printStatus struct {
	GcodeState string `json:"gcode_state"`
}

func Listen(cfg MQTTConfig, onFinish func()) error {
	broker := fmt.Sprintf("ssl://%s:8883", cfg.IP)
	topic := fmt.Sprintf("device/%s/report", cfg.Serial)

	state := &PrintState{}

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetUsername("bblp").
		SetPassword(cfg.AccessCode).
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true}).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(c mqtt.Client) {
			log.Printf("connected to printer MQTT at %s", cfg.IP)
			tok := c.Subscribe(topic, 0, func(_ mqtt.Client, msg mqtt.Message) {
				handleMessage(msg.Payload(), state, onFinish)
			})
			tok.Wait()
			if tok.Error() != nil {
				log.Printf("subscribe error: %v", tok.Error())
			}
		}).
		SetConnectionLostHandler(func(_ mqtt.Client, err error) {
			log.Printf("MQTT connection lost: %v", err)
		})

	client := mqtt.NewClient(opts)
	tok := client.Connect()
	tok.Wait()
	if tok.Error() != nil {
		return fmt.Errorf("MQTT connect: %w", tok.Error())
	}

	log.Printf("subscribed to %s", topic)
	return nil
}

func handleMessage(payload []byte, state *PrintState, onFinish func()) {
	var msg reportMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return
	}
	if msg.Print == nil || msg.Print.GcodeState == "" {
		return
	}

	current := msg.Print.GcodeState

	state.mu.Lock()
	prev := state.previous
	state.previous = current
	state.mu.Unlock()

	if prev != "" && prev != "FINISH" && current == "FINISH" {
		log.Printf("print finished (was %s → FINISH)", prev)
		onFinish()
	}
}
