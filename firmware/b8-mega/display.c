/*
 * b8: A simple USB keypad with 8 programmable buttons.
 *
 * SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
 * SPDX-License-Identifier: BSD-3-Clause
 */

#include <assert.h>
#include <stdlib.h>
#include <string.h>

#include <stm32f0xx.h>

#include "display.h"
#include "display-font.h"

#define display_address 0x3c
#define display_screen_width 128
#define display_screen_height 64
#define display_chars_per_line (display_screen_width / (display_font_width + 1))
#define display_lines (display_screen_height / 8)  // info from ssd1306 datasheet

// ensure that our assumptions are correct
static_assert(display_chars_per_line == 21);
static_assert(display_lines          == 8);
static_assert(display_font_width     == 5);
static_assert(display_font_height    == 7);

#define _line(line)                                 \
{                                                   \
    /* render empty lines to clean screen */        \
    .state = DISPLAY_LINE_STATE_PENDING_COMMANDS,   \
    .commands = {                                   \
        /* Co = 0; D/C# = 0 */                      \
        0x00,                                       \
        /* set line address */                      \
        0xB0 | line,                                \
        /* set column address 4 lower bits */       \
        0x02,                                       \
        /* set column address 4 higher bits */      \
        0x10,                                       \
    },                                              \
    .data = {                                       \
        /* Co = 0; D/C# = 1 */                      \
        0x40,                                       \
    },                                              \
    .with_backlog = false,                          \
}

static struct {
    bool initialized;
    uint8_t current_line;
    const uint8_t commands[6];
    struct {
        enum {
            DISPLAY_LINE_STATE_FREE,
            DISPLAY_LINE_STATE_PENDING_COMMANDS,
            DISPLAY_LINE_STATE_SENDING_COMMANDS,
            DISPLAY_LINE_STATE_PENDING_DATA,
            DISPLAY_LINE_STATE_SENDING_DATA,
        } state;
        const uint8_t commands[4];
        uint8_t data[display_screen_width + 1];
        uint8_t backlog[display_screen_width];
        bool with_backlog;
    } lines[display_lines];
} display = {
    .initialized = false,
    .current_line = 0,
    .commands = {
        /* Co = 0; D/C# = 0 */
        0x00,
        /* set segment remap (reverse direction) */
        0xA1,
        /* set COM output scan direction (COM[N-1] to COM0) */
        0xC8,
        /* charge pump setting (enable during display on) */
        0x8D, 0x14,
        /* display on */
        0xAF,
    },
    .lines = {
        _line(0),
        _line(1),
        _line(2),
        _line(3),
        _line(4),
        _line(5),
        _line(6),
        _line(7),
    },
};

#undef _line


void
display_init(void)
{
    RCC->AHBENR |=  RCC_AHBENR_DMA1EN | RCC_AHBENR_GPIOAEN;
    RCC->APB1ENR |= RCC_APB1ENR_I2C1EN;

    GPIOA->OTYPER |= GPIO_OTYPER_OT_9 | GPIO_OTYPER_OT_10;
    GPIOA->PUPDR |= GPIO_PUPDR_PUPDR9_0 | GPIO_PUPDR_PUPDR10_0;
    GPIOA->AFR[1] |= (4 << GPIO_AFRH_AFSEL9_Pos) | (4 << GPIO_AFRH_AFSEL10_Pos);
    GPIOA->MODER |= (GPIO_MODER_MODER9_1 | GPIO_MODER_MODER10_1);

    I2C1->TIMINGR = 0x00200C1E;  // ~1MHz
    I2C1->CR1 = I2C_CR1_TXDMAEN | I2C_CR1_PE;

    DMA1_Channel2->CPAR = (uint32_t) &(I2C1->TXDR);
    DMA1_Channel2->CCR = DMA_CCR_DIR | DMA_CCR_PL | DMA_CCR_MINC | DMA_CCR_TCIE;
}


static uint8_t*
get_data_ptr(uint8_t line)
{
    switch (display.lines[line].state) {
    case DISPLAY_LINE_STATE_FREE:
        display.lines[line].state = DISPLAY_LINE_STATE_PENDING_COMMANDS;
        // fall through

    case DISPLAY_LINE_STATE_PENDING_COMMANDS:
        return display.lines[line].data + 1;

    default:
    }

    display.lines[line].with_backlog = true;
    return display.lines[line].backlog;
}


static uint8_t
safe_strlen(const char *str)
{
    uint8_t i = 0;
    while (i < display_chars_per_line && str[i] != 0)
        i++;
    return i;
}


bool
display_line(uint8_t line, const char *str, display_halign_t align)
{
    if (str == NULL || line >= display_lines)
        return false;

    uint8_t start = 0;
    uint8_t len = safe_strlen(str);

    switch (align) {
    case DISPLAY_HALIGN_LEFT:
        break;

    case DISPLAY_HALIGN_RIGHT:
        start = (display_chars_per_line - len) * (display_font_width + 1);
        break;

    case DISPLAY_HALIGN_CENTER:
        start = ((display_chars_per_line - len) * (display_font_width + 1)) / 2;
        break;
    }

    uint8_t *buf = get_data_ptr(line);

    uint8_t i = 0;
    uint8_t j = 0;
    uint8_t k = 0;

    while (i < display_screen_width) {
        if (i < start || j >= len) {
            buf[i++] = 0;
            continue;
        }

        buf[i++] = display_font[(uint8_t) str[j]][k++];
        if (k == display_font_width) {
            buf[i++] = 0;
            j++;
            k = 0;
        }
    }

    return true;
}


bool
display_line_from_report(uint8_t *buf, uint8_t len)
{
    if (len != display_chars_per_line + 2)
        return false;

    return display_line(buf[0], (const char*) &(buf[2]), buf[1]);
}


void
display_clear(void)
{
    for (uint8_t i = 0; i < display_lines; i++)
        display_line(i, "", DISPLAY_HALIGN_LEFT);
}


static bool
dma_send(const uint8_t *data, uint32_t data_len)
{
    if ((DMA1_Channel2->CCR & DMA_CCR_EN) == DMA_CCR_EN)
        return false;

    I2C1->CR2 = (display_address << 1) | (data_len << I2C_CR2_NBYTES_Pos) |
        I2C_CR2_AUTOEND | I2C_CR2_START;

    DMA1_Channel2->CMAR = (uint32_t) data;
    DMA1_Channel2->CNDTR = data_len;
    DMA1_Channel2->CCR |= DMA_CCR_EN;

    return true;
}


static bool
task_start(void)
{
    if (!display.initialized) {
        if (dma_send(display.commands, sizeof(display.commands))) {
            display.initialized = true;
            return true;
        }
        return false;
    }

    switch (display.lines[display.current_line].state) {
    case DISPLAY_LINE_STATE_FREE:
        if (++display.current_line == display_lines)
            display.current_line = 0;
        break;

    case DISPLAY_LINE_STATE_SENDING_COMMANDS:
    case DISPLAY_LINE_STATE_SENDING_DATA:
        break;

    case DISPLAY_LINE_STATE_PENDING_COMMANDS:
        if (dma_send(display.lines[display.current_line].commands,
            sizeof(display.lines[display.current_line].commands)))
        {
            display.lines[display.current_line].state = DISPLAY_LINE_STATE_SENDING_COMMANDS;
            return true;
        }
        break;

    case DISPLAY_LINE_STATE_PENDING_DATA:
        if (dma_send(display.lines[display.current_line].data,
            sizeof(display.lines[display.current_line].data)))
        {
            display.lines[display.current_line].state = DISPLAY_LINE_STATE_SENDING_DATA;
            return true;
        }
        break;
    }
    return false;
}


static void
task_stop(void)
{
    switch (display.lines[display.current_line].state) {
    case DISPLAY_LINE_STATE_FREE:
    case DISPLAY_LINE_STATE_PENDING_COMMANDS:
    case DISPLAY_LINE_STATE_PENDING_DATA:
        break;

    case DISPLAY_LINE_STATE_SENDING_COMMANDS:
        display.lines[display.current_line].state = DISPLAY_LINE_STATE_PENDING_DATA;
        break;

    case DISPLAY_LINE_STATE_SENDING_DATA:
        display.lines[display.current_line].state = DISPLAY_LINE_STATE_FREE;
        if (display.lines[display.current_line].with_backlog) {
            memcpy(display.lines[display.current_line].data + 1,
                display.lines[display.current_line].backlog, display_screen_width);
            display.lines[display.current_line].with_backlog = false;
        }
        if (++display.current_line == display_lines)
            display.current_line = 0;
        break;
    }
}


void
display_task(void)
{
    static bool busy = false;

    if ((DMA1->ISR & DMA_ISR_TCIF2) == (DMA_ISR_TCIF2) && (I2C1->ISR & I2C_ISR_STOPF) == (I2C_ISR_STOPF)) {
        I2C1->ICR |= I2C_ICR_STOPCF;
        DMA1->IFCR = DMA_IFCR_CTCIF2;
        DMA1_Channel2->CCR &= ~DMA_CCR_EN;

        task_stop();
        busy = false;
        return;
    }

    if ((!busy) && task_start())
        busy = true;
}
