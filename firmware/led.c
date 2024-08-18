// SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: BSD-3-Clause

#include <stm32f0xx.h>

#include "led.h"


void
led_init(void)
{
    RCC->AHBENR |= RCC_AHBENR_GPIOBEN;
    RCC->APB1ENR |= RCC_APB1ENR_TIM3EN;

    GPIOB->MODER &= ~GPIO_MODER_MODER0;
    GPIOB->MODER |= GPIO_MODER_MODER0_1;
    GPIOB->AFR[0] |= (1 << GPIO_AFRL_AFSEL0_Pos);

    TIM3->PSC = SystemCoreClock / 1000 - 1;
}


void
led_set_state(led_state_t state)
{
    TIM3->CR1 = 0;
    TIM3->ARR = 0;
    TIM3->CCR3 = 0;
    TIM3->CCMR2 = 0;
    TIM3->CCER = TIM_CCER_CC3E;
    TIM3->CNT = 0;

    switch (state) {
    case LED_ON:
        TIM3->CCMR2 |= TIM_CCMR2_OC3M_0;
        // fallthrough

    case LED_OFF:
        TIM3->CCMR2 |= TIM_CCMR2_OC3M_2;
        break;

    case LED_FLASH:
        TIM3->ARR = 50;
        TIM3->CCR3 = 1;
        TIM3->CCMR2 = TIM_CCMR2_OC3M_2 | TIM_CCMR2_OC3M_1 | TIM_CCMR2_OC3M_0;
        TIM3->EGR = TIM_EGR_UG;
        TIM3->CR1 = TIM_CR1_OPM | TIM_CR1_CEN;
        break;

    case LED_SLOW_BLINK:
        TIM3->ARR += 150;
        // fallthrough

    case LED_FAST_BLINK:
        TIM3->ARR += 99;
        TIM3->CCR3 = TIM3->ARR;
        TIM3->CCMR2 = TIM_CCMR2_OC3M_1 | TIM_CCMR2_OC3M_0;
        TIM3->CR1 = TIM_CR1_CEN;
        break;
    }
}
