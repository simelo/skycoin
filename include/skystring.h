#ifndef LIBSKY_STRING_H
#define LIBSKY_STRING_H

#include <stdio.h>
#include <stdlib.h>
#include "libskycoin.h"

extern void randBytes(GoSlice *bytes, size_t n);

extern void strnhex(unsigned char* buf, char *str, int n);

extern void strhex(unsigned char* buf, char *str);

extern int hexnstr(const char* hex, unsigned char* str, int n);


#endif //LIBSKY_STRING_H
