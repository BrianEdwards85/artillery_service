package main

import (
//	"flag"
	"fmt"
	//"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/d2r2/go-i2c"
)

func createOnMessageReceivedHandler(c chan []string) func(client MQTT.Client, message MQTT.Message) {
	return func (client MQTT.Client, message MQTT.Message) {
		byteArray := message.Payload()
		str := fmt.Sprintf("%s", byteArray)
		i, err := strconv.Atoi(str)
		if(err != nil){
			fmt.Printf("Unknown index [%s]",str)
		} else {
			go func() {
				//fmt.Printf("Received message on topic: %s\nMessage: %d\n", message.Topic(), i)

				l := fmt.Sprintf("Pin[%d] on\n",i)
				m := []string{l}
				c <- m

				fmt.Printf("Pin[%d] on\n",i)
				time.Sleep(1 * time.Second)
				fmt.Printf("Pin[%d] off\n",i)
			}()
		}
	}
}


func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	lcdc := make(chan []string)
	i2c, _ := i2c.NewI2C(0x27, 1)
	lcd, _ := InitALcd(i2c,lcdc)

	go func() {
		<-c
		fmt.Println("signal received, exiting")
		i2c.Close()
		lcd.Close()
		os.Exit(0)
	}()

	fmt.Println("Running")

	opts := MQTT.NewClientOptions().AddBroker("tcp://127.0.0.1:1883").SetClientID("artillery_service")

	handler := createOnMessageReceivedHandler(lcdc)
	
	opts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe("/#", byte(1), handler); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Println("Connected")
	}


	for {
		time.Sleep(1 * time.Second)
	}
}
