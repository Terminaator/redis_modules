#ifndef YEARMODULE_H_
#define YEARMODULE_H_

#include "mainmodule.h"

void is_reset_needed(int *, int *);
int reset_document_doty_counts(RedisModuleCtx *, int *);
int reset(RedisModuleCtx *, int *);

#endif