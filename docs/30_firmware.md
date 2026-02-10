# Firmware

Single bare-metal C firmware for the STM32F042K6/K4 (ARM Cortex-M0). No HAL, no RTOS -- direct register access throughout. A single binary runs on both hardware variants: the firmware probes for an SSD1306 display at startup and silently disables display features if none is found. See the hardware pages for [octokeyz](10_octokeyz.md) and [octokeyz-mega](20_octokeyz-mega.md).

The firmware implements a USB HID device with a custom vendor protocol. Host-side interaction is handled by the [Go client library](40_client-libraries.md) or any custom HID implementation that speaks the protocol described below.

## Building from Source

### Prerequisites

- CMake 3.22+
- Ninja build system (optional, but recommended)
- ARM embedded GCC cross-compiler (`arm-none-eabi-gcc`)
- Git

The two library dependencies -- [cmake-cmsis-stm32](https://github.com/rafaelmartins/cmake-cmsis-stm32) (build framework) and [usbd-fs-stm32](https://github.com/rafaelmartins/usbd-fs-stm32) (USB device stack) -- are fetched automatically via CMake FetchContent during configuration.

### Configure and Build

```bash
cmake -B build -DCMAKE_BUILD_TYPE=Release -G Ninja
cmake --build build
```

### Output Artifacts

The build produces the following in the `build/firmware/` directory:

| File | Format |
|------|--------|
| `octokeyz.elf` | ELF binary |
| `octokeyz.bin` | Raw binary |
| `octokeyz.hex` | Intel HEX |
| `octokeyz.dfu` | DFU with suffix (for `dfu-util`) |
| `octokeyz.map` | Linker map |

The firmware version is derived from git tags matching the `v[0-9]*` pattern. Between tagged releases, the version includes the commit count and abbreviated hash (e.g. `0.0.72-1b14`).

### Memory Layout

The linker script (`firmware/STM32F042KxTx_FLASH.ld`) targets the smaller K4 memory to ensure compatibility with both STM32F042K4 and K6 variants:

- Flash: 16 KB at `0x08000000`
- RAM: 6 KB at `0x20000000`

## Flashing

### Using ST-Link

For a freshly built firmware, with `st-flash` (from https://github.com/stlink-org/stlink) installed, run:

```bash
cmake --build build --target octokeyz-stlink-write
```

### Using USB DFU

Firmware flashing can be done over USB using the STM32's built-in DFU bootloader. There are several ways to enter DFU mode:

**Empty Microcontroller:** An empty microcontroller boots directly into the DFU bootloader.

**Button Combo:** With an octokeyz firmware already running on the microcontroller, hold buttons 1 and 5 simultaneously while plugging in the USB cable. The microcontroller boots directly into the DFU bootloader.

**Bootloader Pin Header:** Connect a jumper to JP1 (Boot to DFU) pin header and plug the USB cable. The microcontroller boots directly into the DFU bootloader.

Once in DFU mode, follow the instructions from my generic [Hardware Build Manual's "STM32 (USB DFU)" section](@@/hardware/build-manual/#stm32-usb-dfu).

## Linux udev Rules

Linux users may have issues connecting to the macro pad as a normal user.

The file `share/udev/60-octokeyz.rules` grants device access to users in the `plugdev` group and to logged-in users via the `uaccess` tag:

```
SUBSYSTEMS=="usb", DRIVERS=="usb", ATTRS{idVendor}=="1d50", ATTRS{idProduct}=="6184", GROUP="plugdev", MODE="0660", TAG+="uaccess"
```

Install it:

```bash
sudo cp share/udev/60-octokeyz.rules /etc/udev/rules.d/
sudo udevadm control --reload-rules && sudo udevadm trigger
```

## Architecture Overview

The firmware runs on HSI48 at 48 MHz (internal oscillator, no external crystal). Flash wait state is set to 1 cycle as required at this frequency. AHB and APB buses run undivided.

### Peripheral Map

| Peripheral | Function |
|------------|----------|
| GPIOA PA0-PA7 | Button inputs (internal pull-up) |
| GPIOB PB0 | LED output (TIM3 CH3, alternate function) |
| GPIOA PA9 / PA10 | I2C1 SCL / SDA (display, open-drain with pull-up) |
| USB FS | USB device |
| I2C1 | SSD1306 display communication (address `0x3c`) |
| DMA1 Channel 2 | I2C1 TX -- display data transfer |
| TIM3 | LED PWM generation |
| TIM16 | Display clear delay (one-pulse mode) |
| TIM17 | HID idle rate timer (one-pulse mode) |
| RTC BKP0R | Bootloader magic value handling |

### Source Files

| File | Responsibility |
|------|----------------|
| `main.c` | Clock init, GPIO setup, USB and display init, main loop (`usbd_task()` + `display_task()`), USB callback implementations |
| `descriptors.c` | USB device/config/HID descriptors, string descriptors, HID report descriptor |
| `display.c` | SSD1306 driver: I2C + DMA init, double-buffered line rendering, font rasterization, delayed clear |
| `led.c` | TIM3 PWM setup, 5-state LED control (on, flash, slow blink, fast blink, off) |
| `idle.c` | HID idle rate tracking via TIM17, SET_IDLE/GET_IDLE request handling |
| `bootloader.c` | DFU entry detection (RTC backup register check), system memory remap and jump, reset trigger |

## USB HID Protocol

### Device Identity

| Field | Value |
|-------|-------|
| USB version | 2.0 Full-Speed |
| Device class | HID (interface-level) |
| VID | `0x1d50` (generously provided by OpenMoko) |
| PID | `0x6184` |
| Manufacturer string | `rgm.io` |
| Product string | `octokeyz` |
| Serial number | STM32 unique device ID |
| Max power | 100 mA (bus-powered) |
| HID version | 1.11 |

### Endpoints

The firmware only defines endpoint 1 in both directions:

| Direction | Type | Max Packet Size | Interval |
|-----------|------|-----------------|----------|
| IN | Interrupt | 64 bytes | 10 ms |
| OUT | Interrupt | 64 bytes | 10 ms |

### HID Reports

| ID | Kind | Size (bytes) | Description |
|-----------|------|--------------|-------------|
| 1 | Input | 1 | Button states -- bits 0-7 map to buttons 1-8, `1` = pressed |
| 1 | Output | 1 | LED control -- bits 0-2 encode state: 1=on, 2=flash, 3=slow blink, 4=fast blink, 5=off |
| 1 | Feature | 1 | Capabilities -- bit 0: device has a display |
| 2 | Output | 23 | Display line data -- byte 0: line number (5 bits), byte 1: alignment (2 bits, 1=left, 2=right, 3=center), bytes 2-22: ASCII text (21 chars max) |
| 2 | Feature | 3 | Display capabilities -- byte 0: number of lines, byte 1: characters per line, byte 2 bit 0: supports clear |
| 3 | Output | 2 | Display clear with delay -- 16-bit little-endian milliseconds |

Report sizes listed above do not include the report ID byte.

### Vendor Usage Pages

| Usage Page | Name | Contains |
|------------|------|----------|
| `0xFF00` | octokeyz | Application collection, capabilities feature |
| `0xFF01` | octokeyz Key | Button states (keys 1-8) |
| `0xFF02` | octokeyz LED | LED control (states 1-5) |
| `0xFF03` | octokeyz Display | Display capabilities, line data, clear command |

> [!NOTE]
> This is a custom vendor HID protocol, not a standard keyboard or consumer device. Interacting with it requires the [Go client library](40_client-libraries.md) or a custom implementation that understands the report structure and vendor usage pages described above.

### HID Idle Rate

The default idle rate is 500 ms (value 125 in 4 ms units), configurable via standard HID SET_IDLE and GET_IDLE requests. When idle rate is non-zero, button state reports are sent periodically even if no state change has occurred. Setting idle rate to 0 disables periodic reports -- the device only sends a report when button state actually changes.
