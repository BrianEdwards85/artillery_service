package main

import (
	"fmt"
	"sort"
)

type GPIOMessage struct{
	addr  uint8
	pin   uint8
	value bool
}

type GPIO struct {
	values map[uint8] byte
	i2c    *AI2C
	c      chan GPIOMessage

	disp chan []string
}

func (g *GPIO) Close(){
	close(g.c)
	g.i2c.Close()
}

func (g *GPIO) poll(){
	for mesg := range g.c {
		value, exists := g.values[mesg.addr]
		if !exists{
			value = 0xff
		}
		var masq byte = 1 << (mesg.pin - 1)
		if mesg.value {
			value = value & ^masq
		} else {
			value = value | masq
		}
		g.values[mesg.addr] = value
		g.i2c.Send(mesg.addr,value)

		p := g.ToString()
		g.disp <- p
	}
}

func byteToString(b byte) []byte{
	r := make([]byte,8)
	var i uint8 = 0
	for ;i < 8; i++{
		var masq byte = 1 << i
		if (b & masq) == 0 {
			r[i] = ('1' + i)
		} else {
			r[i] = '*'
		}
	}
	return r
}

func (g *GPIO) ToString() []string {
	cnt := len(g.values)
	r := make([]string,cnt)
	i := 0
	for k,v := range g.values {
		r[i] = fmt.Sprintf(" %X: %s", k, byteToString(v))
		i += 1
	}
	sort.Strings(r)
	return r
}

func (g *GPIO)Send(addr uint8, pin uint8, value bool){
	msg := GPIOMessage{addr: addr, pin: pin, value: value}
	g.c <- msg
}

func InitGPIO(bus int, disp chan []string) *GPIO {
	c := make(chan GPIOMessage)
	i2c := InitAI2C(bus)
	values := make(map[uint8] byte)
	this := &GPIO{values: values, i2c: i2c, c: c, disp: disp}
	go this.poll()
	return this
}
