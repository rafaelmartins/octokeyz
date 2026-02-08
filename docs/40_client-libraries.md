# Client Libraries

The octokeyz firmware communicates over a custom USB HID protocol using vendor-specific usage pages (`0xFF00`-`0xFF03`). Interacting with the device from userspace requires a client library that understands this protocol. Currently, a Go library is available.

## Go

**Package**: `rafaelmartins.com/p/octokeyz`

The Go client library provides device discovery, button event handling, LED control, and display output for both hardware variants.

- [Source code](https://github.com/rafaelmartins/go-octokeyz)
- [API documentation](https://pkg.go.dev/rafaelmartins.com/p/octokeyz)
- License: BSD-3-Clause

### Installation

```
go get rafaelmartins.com/p/octokeyz
```

### Device Lifecycle

`Enumerate()` returns a list of all connected octokeyz devices. `GetDevice()` returns a specific device by serial number, or auto-detects if exactly one device is connected when called with an empty string.

Once you have a `*Device`, call `Open()` before interacting with it and `Close()` when done.

```go
dev, err := octokeyz.GetDevice("")
if err != nil {
    log.Fatal(err)
}

if err := dev.Open(); err != nil {
    log.Fatal(err)
}
defer dev.Close()

fmt.Println("connected to:", dev.SerialNumber())
```

If multiple devices are connected, use `Enumerate()` to list them, or `GetDevice()` with a specific serial number:

```go
devices, err := octokeyz.Enumerate()
if err != nil {
    log.Fatal(err)
}
for _, d := range devices {
    fmt.Println(d.SerialNumber())
}

dev, err := octokeyz.GetDevice("200014000A43304D45363820")
```

### Handling Button Events

Register callbacks for individual buttons with `AddHandler()`, then call `Listen()` to start the blocking event loop.

Handlers receive a `*Button` argument. Inside a handler, `WaitForRelease()` blocks until the button is released and returns the press duration. `GetID()` returns the `ButtonID`.

```go
dev.AddHandler(octokeyz.BUTTON_1, func(b *octokeyz.Button) error {
    fmt.Printf("button %s pressed\n", b)
    duration := b.WaitForRelease()
    fmt.Printf("released after %s\n", duration)
    return nil
})

errCh := make(chan error, 1)
if err := dev.Listen(errCh); err != nil {
    log.Fatal(err)
}
```

Button constants are `BUTTON_1` through `BUTTON_8`.

`Listen()` blocks indefinitely, dispatching events to registered handlers. If a handler returns an error, it is wrapped in a `ButtonHandlerError` and sent to the `errCh` channel without stopping the event loop. Pass `nil` for `errCh` to discard handler errors.

### Modifier Buttons

The `Modifier` type implements shift/modifier-like functionality. Register its `Handler` method on a button, then check `Pressed()` from other handlers to branch on modifier state.

```go
modifier := &octokeyz.Modifier{}

dev.AddHandler(octokeyz.BUTTON_8, modifier.Handler)

dev.AddHandler(octokeyz.BUTTON_1, func(b *octokeyz.Button) error {
    b.WaitForRelease()
    if modifier.Pressed() {
        fmt.Println("button 1 + modifier")
    } else {
        fmt.Println("button 1")
    }
    return nil
})
```

### LED Control

`Led()` sets the indicator LED state. Five states are available:

| State | Behavior |
|-------|----------|
| `LedOn` | Steady on |
| `LedFlash` | Single 50ms pulse |
| `LedSlowBlink` | 250ms period |
| `LedFastBlink` | 100ms period |
| `LedOff` | Off |

`LedFlash` triggers a single short pulse on the firmware side. To create a visible flash pattern, call it multiple times with sleep intervals:

```go
for i := 0; i < 3; i++ {
    dev.Led(octokeyz.LedFlash)
    time.Sleep(100 * time.Millisecond)
}
```

### Display Control

Display functions are only operational on [octokeyz-mega](20_octokeyz-mega.md). On the basic variant, the firmware reports no display capability by returning an `ErrDeviceDisplayNotSupported` error whenever a function that requires a display is called.

The display has 8 lines (`DisplayLine1` through `DisplayLine8`) with 21 characters per line. `GetDisplayCharsPerLine()` returns this value at runtime, or `0` if the device has no display.

`DisplayLine()` writes a string to a specific line with alignment:

```go
dev.DisplayLine(octokeyz.DisplayLine1, "octokeyz", octokeyz.DisplayLineAlignCenter)
dev.DisplayLine(octokeyz.DisplayLine3, "left",     octokeyz.DisplayLineAlignLeft)
dev.DisplayLine(octokeyz.DisplayLine4, "right",    octokeyz.DisplayLineAlignRight)
```

Alignment options are `DisplayLineAlignLeft`, `DisplayLineAlignRight`, and `DisplayLineAlignCenter`.

`DisplayClearLine()` clears a single line. `DisplayClear()` clears the entire display immediately. `DisplayClearWithDelay()` clears the display after a firmware-side delay (millisecond resolution, up to 65535ms) -- useful for showing transient information without needing a client-side timer:

```go
dev.DisplayLine(octokeyz.DisplayLine1, "done!", octokeyz.DisplayLineAlignCenter)
dev.DisplayClearWithDelay(2 * time.Second)
```

### Complete Example

```go
package main

import (
    "fmt"
    "log"
    "time"

    "rafaelmartins.com/p/octokeyz"
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

### API Reference

For the full API surface, type details, and additional examples, see the [package documentation on pkg.go.dev](https://pkg.go.dev/rafaelmartins.com/p/octokeyz).

Source code is available at [github.com/rafaelmartins/go-octokeyz](https://github.com/rafaelmartins/go-octokeyz).
