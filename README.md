# b8

A simple USB keypad with 8 programmable buttons.


## Motivation

I wanted to have a simple keypad I could use to control my computer.

I wanted to be able to write userspace programs in Golang, that would react to the keypress events in the keypad, and execute some programmed action, instead of building long sequences of keypress macros.

I wanted the PCB to be simple (PTH parts only), to have only the buttons and a single indicator LED, and to use the simplest/smallest microcontroller that could handle USB 1.1 and 8 buttons.

I wanted the enclosure to be 3D-printed at home.

I wanted it to be as USB HID compliant as possible, so I could learn more about it.

I wanted the client library to support Linux and Windows.


## What is included

- [Firmware source code](./firmware/)
- [Golang client library](./go/b8/)
- [Printed Circuit Board (Kicad sources)](./pcb/)
- [3D models for a simple enclosure](./3d-models/)
- [`udev` rules for Linux](./share/udev/)


## Pictures

![Front](./share/images/r1.0-front.jpg)
![Back](./share/images/r1.0-back.jpg)
![PCB Front](./share/images/r1.0-pcb-front.jpg)
![PCB Back](./share/images/r1.0-pcb-back.jpg)


## Program examples

### Simple

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rafaelmartins/b8/go/b8"
)

func main() {
	dev, err := b8.GetDevice("")
	if err != nil {
		log.Fatal(err)
	}

	if err := dev.Open(); err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	for i := 0; i < 3; i++ {
		dev.Led(b8.LedFlash)
		time.Sleep(100 * time.Millisecond)
	}

	dev.AddHandler(b8.BUTTON_1, func(b *b8.Button) error {
		fmt.Println("pressed")
		duration := b.WaitForRelease()
		fmt.Printf("released. pressed for %s\n", duration)
		return nil
	})

	if err := dev.Listen(); err != nil {
		log.Fatal(err)
	}
}
```


## F.A.Q.

### How to use this keypad to control `OBS`, similarly to what the `Stream Deck` does?

You can write event handlers that interact with `OBS` by using the `goobs` library: https://github.com/andreykaipov/goobs
