package main

import (
        "log"

        "github.com/d2r2/go-hd44780"
        "github.com/d2r2/go-i2c"
)

var displayLines = [...]hd44780.ShowOptions{hd44780.SHOW_LINE_1, hd44780.SHOW_LINE_2, hd44780.SHOW_LINE_3, hd44780.SHOW_LINE_4}

func checkError(err error) {
        if err != nil {
                log.Fatal(err)
        }
}

type ALcd struct {
	i2c       *i2c.I2C
	lcd       *hd44780.Lcd
	c         chan []string
}

func (a *ALcd) poll(){
	for t := range a.c {
		a.lcd.Clear()
		for i, l := range t {
			if i < 4 {
				truncLine := l
				if len(l) > 20{
					truncLine = l[:20]
				}
				a.lcd.ShowMessage(truncLine,displayLines[i])
			}}}}

func (a *ALcd) Close(){
	close(a.c)
}

func InitALcd(i2c *i2c.I2C, c chan []string) (*ALcd, error){
        lcd, err := hd44780.NewLcd(i2c, hd44780.LCD_20x4)
	if(err != nil){
		return nil, err
	}
	err = lcd.BacklightOn()
	if(err != nil){
		return nil, err
	}
	this := &ALcd{i2c: i2c, lcd: lcd, c: c}
	go this.poll()
	return this, nil
}

