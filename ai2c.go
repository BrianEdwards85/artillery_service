package main

import "github.com/d2r2/go-i2c"

type AI2CMessage struct {
	Addr uint8
	Value byte
}

type AI2C struct{
	i2cMap    map [uint8] *i2c.I2C
	bus       int
	c         chan AI2CMessage
}

func (a *AI2C) getI2CPort(addr uint8) (*i2c.I2C, error){
	i2cPort, exists := a.i2cMap[addr]
	if exists {
		return i2cPort, nil
	} else {
		i2cPort, err := i2c.NewI2C(addr, a.bus)
		if err != nil {
			return nil, err
		}

		i2cPort.WriteRegU8(9,0xff)
		i2cPort.WriteRegU8(0,0)

		a.i2cMap[addr] = i2cPort
		return i2cPort, nil
	}
}

func (a *AI2C) poll(){
	for mesg := range a.c {
		i2cPort, err := a.getI2CPort(mesg.Addr)
		if err == nil{
			i2cPort.WriteRegU8(9,mesg.Value)
		}
	}
}

func (a *AI2C) Close(){
	close(a.c)
	for _, i2cPort := range a.i2cMap {
		if i2cPort != nil {
			i2cPort.Close()
		}
	}
}

func (a *AI2C) Send(addr uint8, value byte){
	msg := AI2CMessage{Addr: addr, Value: value}
	a.c <- msg
}

func InitAI2C(bus int) *AI2C {
	c := make(chan AI2CMessage)
	i2cMap := make(map[uint8] *i2c.I2C)
	this := &AI2C{i2cMap: i2cMap, bus: bus, c: c}
	go this.poll()
	return this
}

/*
func trigger(i2c *i2c.I2C, pins []uint8){
	var masq uint8 = 0
	for _, pin := range pins {
		masq = masq | (1 << (pin - 1))
	}
	_ = i2c.WriteRegU8(9,^masq)

	time.Sleep(1 * time.Second)

	_ = i2c.WriteRegU8(9,255)
}

func main() {
	ai2c := InitAI2C(1)

	ai2c.Send(0x20,0xfe)
	time.Sleep(1 * time.Second)

	ai2c.Send(0x20,0xfd)
	time.Sleep(1 * time.Second)

	ai2c.Send(0x20,0xfb)
	time.Sleep(1 * time.Second)

	ai2c.Send(0x20,0xff)
	time.Sleep(1 * time.Second)

	ai2c.Close()
}
        i2c, err := i2c.NewI2C(0x20, 1)
        checkError(err)
        defer i2c.Close()

	b, err := i2c.ReadRegU8(9)
	fmt.Printf("S: %d\n", b);

	err = i2c.WriteRegU8(9,240)

	b, err = i2c.ReadRegU8(9)
	fmt.Printf("S: %d\n", b);

	time.Sleep(1 * time.Second)

	err = i2c.WriteRegU8(9,255)
	b, err = i2c.ReadRegU8(9)
	fmt.Printf("S: %d\n", b);


	arr := []uint8 {1, 3}
	trigger(i2c, arr)
	}
*/
