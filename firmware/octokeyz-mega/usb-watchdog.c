// SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: BSD-3-Clause

#include <stm32f0xx.h>

#include "usb-watchdog.h"

#define timeout_ms 100


void
usb_watchdog_init(void)
{
    RCC->APB1ENR |= RCC_APB1ENR_TIM2EN;

    TIM2->PSC = SystemCoreClock / 1000 - 1;
    TIM2->ARR = timeout_ms - 1;
    TIM2->CR1 = TIM_CR1_URS;
    TIM2->DIER = TIM_DIER_UIE;
}


void
usb_watchdog_reset(void)
{
    TIM2->EGR |= TIM_EGR_UG;
    if ((TIM2->CR1 & TIM_CR1_CEN) != TIM_CR1_CEN)
        TIM2->CR1 |= TIM_CR1_CEN;
}


void
usb_watchdog_task(void)
{
    if ((TIM2->SR & TIM_SR_UIF) == TIM_SR_UIF) {
        TIM2->SR &= ~TIM_SR_UIF;
        TIM2->CR1 &= ~TIM_CR1_CEN;
        usb_watchdog_cb();
    }
}
