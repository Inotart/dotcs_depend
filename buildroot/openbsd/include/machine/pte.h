#ifndef _BSD_MACHINE_PTE_H_
#define _BSD_MACHINE_PTE_H_

#if defined (__i386__)
#include "i386/pte.h"
#elif defined (__amd64__) || defined (__x86_64__)
#include "amd64/pte.h"
#elif defined (__arm__)
#include "arm/pte.h"
#elif defined (__arm64__) || defined (__aarch64__)
#include "arm64/pte.h"
#else
#error architecture not supported
#endif

#endif /* _BSD_MACHINE_PTE_H_ */
