package main

import "time"


type TriggerMessage struct{
	Addr  uint8
	Pin   uint8
}

type Trigger struct {
	gpio    *GPIO
	c       chan TriggerMessage
}

func (t *Trigger) Close(){
	close(t.c)
	t.gpio.Close()
}

func (t *Trigger) TriggerPin(mesg TriggerMessage) {
	t.gpio.Send(mesg.Addr, mesg.Pin, true)

	time.Sleep(1 * time.Second)

	t.gpio.Send(mesg.Addr, mesg.Pin, false)
}

func (t *Trigger)Send(addr uint8, pin uint8){
	mesg := TriggerMessage{Addr: addr, Pin: pin}
	t.c <- mesg
}

func (t *Trigger)poll(){
	for mesg := range t.c {
		go t.TriggerPin(mesg)
	}
}

func InitTrigger(bus int, disp chan []string) *Trigger{
	c := make(chan TriggerMessage)
	gpio := InitGPIO(bus, disp)
	this := &Trigger{gpio: gpio, c: c}
	go this.poll()
	return this
}

