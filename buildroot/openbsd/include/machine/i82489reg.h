#ifndef _BSD_MACHINE_I82489REG_H_
#define _BSD_MACHINE_I82489REG_H_

#if defined (__i386__)
#include "i386/i82489reg.h"
#elif defined (__amd64__) || defined (__x86_64__)
#include "amd64/i82489reg.h"
#elif defined (__arm__)
#include "arm/i82489reg.h"
#elif defined (__arm64__) || defined (__aarch64__)
#include "arm64/i82489reg.h"
#else
#error architecture not supported
#endif

#endif /* _BSD_MACHINE_I82489REG_H_ */
