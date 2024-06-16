# octokeyz

A simple USB macropad with 8 programmable buttons.


## Motivation / Project requirements

- [x] I want to have a simple macropad I can use to control my computer.
- [x] I want to be able to write userspace programs in Golang, that will react to the keypress events in the macropad and execute some Golang code, instead of building long sequences of keypress macros.
- [x] I want the PCB to be simple (PTH parts only), to have only the buttons and a single indicator LED, and to use the simplest/smallest microcontroller that can handle USB 1.1 and 8 buttons, like the [`ATtiny4313`](https://www.microchip.com/en-us/product/attiny4313).
- [x] I want the enclosure to be 3D-printable at home.
- [x] I want the firmware to be as USB HID compliant as possible, so I can learn more about the USB stack and specifications.
- [x] I want the client library to support at least Linux and Windows.

> [!NOTE]
> After using the original [`octokeyz`](#octokeyz) macropad for a few months I realized that having a small OLED screen added to the macropad could be very useful. This new addition required using a more powerful microcontroller (I picked the [`STM32F042K6`](https://www.st.com/en/microcontrollers-microprocessors/stm32f042k6.html)/[`STM32F042K4`](https://www.st.com/en/microcontrollers-microprocessors/stm32f042k4.html), which is not PTH, but still quite easy to hand-solder).
>
> These additions resulted in a new `octokeyz` macropad variant named [`octokeyz-mega`](#octokeyz-mega). This new variant also includes support for 5-pin mechanical keyboard switches instead of the simpler 12mm SPST push-buttons used in the original variant.
>
> Given the low price of the STM32 microcontrollers nowadays, and to simplify the project maintenance, I ended up converting the original `octokeyz` variant to also use the [`STM32F042K6`](https://www.st.com/en/microcontrollers-microprocessors/stm32f042k6.html)/[`STM32F042K4`](https://www.st.com/en/microcontrollers-microprocessors/stm32f042k4.html) parts. This way we can deploy the same DFU-capable firmware to any of the board variants.


## Variants

> [!TIP]
> The following resources are common to all variants:
>
> - [Firmware source code](./firmware/)
> - [Golang client library](./go/octokeyz/)
> - [`udev` rules for Linux](./share/udev/)


### octokeyz-mega

- [Schematics](./pcb/octokeyz-mega/octokeyz-mega.pdf)
- [Interactive Bill of Materials](https://rafaelmartins.github.io/octokeyz/ibom/octokeyz-mega.html)
- [Kicad sources](./pcb/octokeyz-mega/)
- [Enclosure 3D models](./3d-models/octokeyz-mega/)

![octokeyz-mega Front](./share/images/octokeyz-mega/front.jpg)
![octokeyz-mega Top](./share/images/octokeyz-mega/top.jpg)
![octokeyz-mega Side](./share/images/octokeyz-mega/side.jpg)


### octokeyz

- [Schematics](./pcb/octokeyz-mega/octokeyz.pdf)
- [Interactive Bill of Materials](https://rafaelmartins.github.io/octokeyz/ibom/octokeyz.html)
- [Kicad sources](./pcb/octokeyz/)
- [Enclosure 3D models](./3d-models/octokeyz/)

![octokeyz Front](./share/images/octokeyz/front.jpg)
![octokeyz Top](./share/images/octokeyz/top.jpg)


## Program examples

### Simple

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rafaelmartins/octokeyz/go/octokeyz"
)

func main() {
	dev, err := octokeyz.GetDevice("")
	if err != nil {
		log.Fatal(err)
	}

	if err := dev.Open(); err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	for i := 0; i < 3; i++ {
		dev.Led(octokeyz.LedFlash)
		time.Sleep(100 * time.Millisecond)
	}

	dev.AddHandler(octokeyz.BUTTON_1, func(b *octokeyz.Button) error {
		fmt.Println("pressed")
		duration := b.WaitForRelease()
		fmt.Printf("released. pressed for %s\n", duration)
		return nil
	})

	if err := dev.Listen(nil); err != nil {
		log.Fatal(err)
	}
}
```


## F.A.Q.

### How to implement a more complex client software?

Please check the Golang [API documentation](https://pkg.go.dev/github.com/rafaelmartins/octokeyz/go/octokeyz).

### How to use this macropad to control `OBS`, similarly to what the `Stream Deck` does?

It is possible to write Golang code that interacts with `OBS` by using the `goobs` library: https://github.com/andreykaipov/goobs. This library could be easily integrated with our [Golang client library](./go/octokeyz/).
