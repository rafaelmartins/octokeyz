// SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: BSD-3-Clause

#include <stdbool.h>
#include <stdint.h>

#include <stm32f0xx.h>

#include "idle.h"

static uint8_t idle_val = 125;  // 500ms


void
idle_init(void)
{
    RCC->APB2ENR |= RCC_APB2ENR_TIM17EN;

    TIM17->CR1 = TIM_CR1_OPM | TIM_CR1_URS;
    TIM17->DIER = TIM_DIER_UIE;
    TIM17->PSC = SystemCoreClock / 1000 - 1;
}


void
idle_set(uint8_t val)
{
    idle_val = val;
}


uint8_t
idle_get(void)
{
    return idle_val;
}


bool
idle_request(void)
{
    if (idle_val == 0) {
        if ((TIM17->CR1 & TIM_CR1_CEN) == TIM_CR1_CEN)
            TIM17->CR1 &= ~TIM_CR1_CEN;
        return false;
    }

    if ((TIM17->SR & TIM_SR_UIF) == TIM_SR_UIF) {
        TIM17->SR &= ~TIM_SR_UIF;
        return true;
    }

    if ((TIM17->CR1 & TIM_CR1_CEN) != TIM_CR1_CEN) {
        TIM17->ARR = (((uint16_t) idle_val) << 2) - 1;
        TIM17->EGR = TIM_EGR_UG;
        TIM17->CR1 |= TIM_CR1_CEN;
    }
    return false;
}
