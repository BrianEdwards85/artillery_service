package main

import {
	"time"
//	"sync"
}

type GPIOMessage struct{
	addr  int
	pin   uint8
	value bool
//	mux   sync.Mutex
}

type GPIO struct {
	values [127] byte
	i2c    *AI2C
	c      chan GPIOMessage
}

func (g *GPIO) Close(){
	close(g.c)
}

func (g *GPIO) poll(){
	for mesg := range g.c {
		value := g.values[mesg.addr]
		masq := 1 << (msg.pin - 1)
		if msg.value {
			value = value & ^masq
		} else {
			value = value | masq
		}
		g.values[mesg.addr] = value
		g.i2c.Send(msg.addr,value)

	}
}

func InitGPIO(bus int) *GPIO {
	c := make(chan GPIOMessage)
}

