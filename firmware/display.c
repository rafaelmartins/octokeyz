// SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: BSD-3-Clause

#include <assert.h>
#include <stdlib.h>
#include <string.h>

#include <stm32f0xx.h>

#include "display.h"

// ensure that our assumptions are correct
static_assert(display_chars_per_line == 21);
static_assert(display_lines          == 8);
static_assert(display_font_width     == 5);
static_assert(display_font_height    == 7);

static const uint8_t init_commands[] = {
    /* send only commands */
    0x00,
    /* set segment remap (reverse direction) */
    0xA1,
    /* set COM output scan direction (COM[N-1] to COM0) */
    0xC8,
    /* charge pump setting (enable during display on) */
    0x8D, 0x14,
    /* display on */
    0xAF,
};

#define _line_data_init(line)               \
{                                           \
    /* set line address */                  \
    0x80, 0xB0 | line,                      \
    /* set column address 4 lower bits */   \
    0x80, 0x02,                             \
    /* set column address 4 higher bits */  \
    0x80, 0x10,                             \
    /* send only data */                    \
    0x40,                                   \
}
static uint8_t line_data_init[] = _line_data_init(0);

#define _line(line)             \
{                               \
    .toggle = false,            \
    .pending = {false, false},  \
    .data = {                   \
        _line_data_init(line),  \
        _line_data_init(line),  \
    },                          \
}

static struct {
    bool initialized;
    bool reset_line;
    uint8_t current_line;
    struct {
        bool toggle;
        bool pending[2];
        uint8_t data[2][display_screen_width + sizeof(line_data_init)];
    } lines[display_lines];
} display = {
    .initialized = false,
    .reset_line = false,
    .current_line = 0,
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
#undef _line_data_init


static inline bool
check_availability(void)
{
    for (uint8_t i = 10; i; i--) {
        I2C1->CR1 = I2C_CR1_PE;
        I2C1->CR2 = (display_address << 1) | I2C_CR2_START | I2C_CR2_AUTOEND;

        for (uint16_t cnt = 0xffff; cnt; cnt--) {
            // display is (explicitly) not ready
            if ((I2C1->ISR & (I2C_ISR_STOPF | I2C_ISR_NACKF)) == (I2C_ISR_STOPF | I2C_ISR_NACKF)) {
                TIM16->ARR = 49;  // 50ms
                TIM16->EGR = TIM_EGR_UG;
                TIM16->CR1 |= TIM_CR1_CEN;
                while ((TIM16->SR & TIM_SR_UIF) != TIM_SR_UIF);
                TIM16->SR &= ~TIM_SR_UIF;
                break;
            }

            // display is ready
            if ((I2C1->ISR & (I2C_ISR_STOPF | I2C_ISR_NACKF)) == I2C_ISR_STOPF) {
                I2C1->CR1 &= ~I2C_CR1_PE;
                return true;
            }
        }
        I2C1->CR1 &= ~I2C_CR1_PE;
    }
    return false;
}


bool
display_init(void)
{
    RCC->CFGR3 |= RCC_CFGR3_I2C1SW_SYSCLK;

    RCC->AHBENR |=  RCC_AHBENR_DMA1EN | RCC_AHBENR_GPIOAEN;
    RCC->APB1ENR |= RCC_APB1ENR_I2C1EN;
    RCC->APB2ENR |= RCC_APB2ENR_TIM16EN;

    GPIOA->OTYPER |= GPIO_OTYPER_OT_9 | GPIO_OTYPER_OT_10;
    GPIOA->PUPDR |= GPIO_PUPDR_PUPDR9_0 | GPIO_PUPDR_PUPDR10_0;
    GPIOA->AFR[1] |= (4 << GPIO_AFRH_AFSEL9_Pos) | (4 << GPIO_AFRH_AFSEL10_Pos);
    GPIOA->MODER |= (GPIO_MODER_MODER9_1 | GPIO_MODER_MODER10_1);

    TIM16->CR1 = TIM_CR1_OPM | TIM_CR1_URS;
    TIM16->DIER = TIM_DIER_UIE;
    TIM16->PSC = SystemCoreClock / 1000 - 1;

    I2C1->TIMINGR = 0x00200C1E;  // ~1MHz

    if (!check_availability())
        return false;

    I2C1->CR1 = I2C_CR1_TXDMAEN | I2C_CR1_PE;

    DMA1_Channel2->CPAR = (uint32_t) &(I2C1->TXDR);
    DMA1_Channel2->CCR = DMA_CCR_DIR | DMA_CCR_PL | DMA_CCR_MINC | DMA_CCR_TCIE;

    display_clear();
    return true;
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

    if ((I2C1->CR1 & I2C_CR1_TXDMAEN) != I2C_CR1_TXDMAEN)
        return true;

    if ((TIM16->CR1 & TIM_CR1_CEN) == TIM_CR1_CEN) {
        TIM16->CR1 &= ~TIM_CR1_CEN;
        display_clear();
    }

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

    bool toggle = display.lines[line].toggle;
    uint8_t *buf = display.lines[line].data[toggle] + sizeof(line_data_init);

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

    display.lines[line].pending[toggle] = true;
    if (line == 0)
        display.reset_line = true;
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
display_clear_line(uint8_t line)
{
    bool toggle = display.lines[line].toggle;
    memset(display.lines[line].data[toggle] + sizeof(line_data_init), 0, display_screen_width);
    display.lines[line].pending[toggle] = true;
}


void
display_clear(void)
{
    for (uint8_t i = 0; i < display_lines; i++)
        display_clear_line(i);
}


void
display_clear_with_delay(uint16_t ms)
{
    if (ms <= 1) {
        display_clear();
        return;
    }

    TIM16->ARR = ms - 1;
    TIM16->EGR = TIM_EGR_UG;
    TIM16->CR1 |= TIM_CR1_CEN;
}


static inline void
dma_send(const uint8_t *data, uint32_t data_len)
{
    I2C1->CR2 = (display_address << 1) | (data_len << I2C_CR2_NBYTES_Pos) |
        I2C_CR2_AUTOEND | I2C_CR2_START;

    DMA1_Channel2->CMAR = (uint32_t) data;
    DMA1_Channel2->CNDTR = data_len;
    DMA1_Channel2->CCR |= DMA_CCR_EN;
}


void
display_task(void)
{
    if ((I2C1->CR1 & I2C_CR1_TXDMAEN) != I2C_CR1_TXDMAEN)
        return;

    if ((TIM16->SR & TIM_SR_UIF) == TIM_SR_UIF) {
        TIM16->SR &= ~TIM_SR_UIF;
        display_clear();
        return;
    }

    bool toggle = display.lines[display.current_line].toggle;

    if (((DMA1->ISR & DMA_ISR_TCIF2) == DMA_ISR_TCIF2) && ((I2C1->ISR & I2C_ISR_STOPF) == I2C_ISR_STOPF)) {
        I2C1->ICR = I2C_ICR_STOPCF;
        DMA1->IFCR = DMA_IFCR_CTCIF2;
        DMA1_Channel2->CCR &= ~DMA_CCR_EN;

        if (!display.initialized)
            display.initialized = true;

        display.lines[display.current_line].pending[!toggle] = false;
        return;
    }

    if ((DMA1_Channel2->CCR & DMA_CCR_EN) == DMA_CCR_EN)
        return;

    if (!display.initialized) {
        dma_send(init_commands, sizeof(init_commands));
        return;
    }

    if (display.reset_line) {
        display.current_line = 0;
        display.reset_line = false;
    }

    if (display.lines[display.current_line].pending[toggle]) {
        display.lines[display.current_line].toggle = !toggle;
        dma_send(display.lines[display.current_line].data[toggle], sizeof(display.lines[0].data[0]));
        return;
    }

    if (++display.current_line == display_lines)
        display.current_line = 0;
}
