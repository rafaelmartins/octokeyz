// SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: BSD-3-Clause

#include <stdbool.h>

#include <stm32f0xx.h>

#define BOOT_TO_DFU 0xdeadbeef


void
bootloader_entry(void)
{
    if ((RCC->CSR & RCC_CSR_SFTRSTF) == 0)
        return;

    RCC->CSR &= ~RCC_CSR_SFTRSTF;

    if (RTC->BKP0R != BOOT_TO_DFU)
        return;

    RCC->APB1ENR |= RCC_APB1ENR_PWREN;

    PWR->CR = PWR_CR_DBP;
    RTC->BKP0R = 0;
    PWR->CR = 0;

    RCC->APB1ENR &= ~RCC_APB1ENR_PWREN;

    if (FLASH->OBR & FLASH_OBR_BOOT_SEL) {
        RCC->AHBENR |= RCC_AHBENR_GPIOBEN;
        GPIOB->MODER &= ~GPIO_MODER_MODER8;
        if ((GPIOB->IDR & GPIO_IDR_8) != 0)
            return;
        GPIOB->MODER |= GPIO_MODER_MODER8_0;
        GPIOB->BSRR = GPIO_BSRR_BS_8;
    }

    RCC->APB2ENR |= RCC_APB2ENR_SYSCFGEN;
    SYSCFG->CFGR1 = SYSCFG_CFGR1_MEM_MODE_0;

    __IO uint32_t bootloader_stack_pointer = *(uint32_t*) 0x1FFFC400UL;
    __IO uint32_t bootloader_address = *(uint32_t*) 0x1FFFC404UL;

    // avoid using __setMSP() and jump to bootloader in C because gcc tends to
    // generate code that stores the bootloader address in the stack pointer
    // instead of a normal register when building with -O3. gcc can't identify
    // that we changed the stack pointer and loads an invalid address when
    // trying to jump to bootloader.
    __ASM volatile (
        "MSR msp, %0\r\n"
        "BLX %1\r\n"
        :: "r" (bootloader_stack_pointer), "r" (bootloader_address)
    );

    while (true)
        __NOP();
}


void
bootloader_reset(void)
{
    RCC->APB1ENR |= RCC_APB1ENR_PWREN;

    PWR->CR = PWR_CR_DBP;
    RTC->BKP0R = BOOT_TO_DFU;
    PWR->CR = 0;

    NVIC_SystemReset();
}
