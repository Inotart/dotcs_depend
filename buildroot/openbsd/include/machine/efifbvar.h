#ifndef _BSD_MACHINE_EFIFBVAR_H_
#define _BSD_MACHINE_EFIFBVAR_H_

#if defined (__i386__)
#include "i386/efifbvar.h"
#elif defined (__amd64__) || defined (__x86_64__)
#include "amd64/efifbvar.h"
#elif defined (__arm__)
#include "arm/efifbvar.h"
#elif defined (__arm64__) || defined (__aarch64__)
#include "arm64/efifbvar.h"
#else
#error architecture not supported
#endif

#endif /* _BSD_MACHINE_EFIFBVAR_H_ */
