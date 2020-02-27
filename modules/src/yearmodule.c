#include "yearmodule.h"

void is_reset_needed(int *year, int *update)
{
    time_t current_time = time(NULL);
    int *current_year = &(localtime(&current_time)->tm_year);

    *current_year %= 100;

    if (*current_year != *year)
    {
        *year = *current_year;
        *update = 1;
    }
}
int reset_document_doty_counts(RedisModuleCtx *ctx, int *year)
{
    RedisModuleCallReply *reply = RedisModule_Call(ctx, "HKEYS", "c", DOCUMENT_KEY);

    if (RedisModule_CallReplyType(reply) == REDISMODULE_REPLY_ARRAY)
    {
        size_t reply_len = RedisModule_CallReplyLength(reply);
        if (reply_len > 0)
        {
            for (size_t x = 0; x < reply_len; x++)
            {
                RedisModuleCallReply *element = RedisModule_CallReplyArrayElement(reply, x);
                RedisModuleCallReply *elememt_reply = RedisModule_Call(ctx, "HSET", "csc", DOCUMENT_KEY, RedisModule_CreateStringFromCallReply(element), "0");
                if (RedisModule_CallReplyType(elememt_reply) == REDISMODULE_REPLY_ERROR)
                {
                    return RedisModule_ReplyWithError(ctx, "Error occurred when getting value");
                }
            }
        }
        RedisModule_ReplyWithLongLong(ctx, *year);
        return REDISMODULE_OK;
    }

    return RedisModule_ReplyWithError(ctx, "Error occurred when getting value");
}

int reset(RedisModuleCtx *ctx, int *year)
{
    RedisModuleCallReply *reply = RedisModule_Call(ctx, "SET", "cl", YEAR_KEY, *year);
    if (RedisModule_CallReplyType(reply) == REDISMODULE_REPLY_ERROR)
    {
        return RedisModule_ReplyWithError(ctx, "Error occurred when getting value");
    }
    return reset_document_doty_counts(ctx, year);
}

int command(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc != 1)
    {
        return RedisModule_WrongArity(ctx);
    }

    RedisModule_AutoMemory(ctx);

    RedisModuleCallReply *reply = RedisModule_Call(ctx, "GET", "c", YEAR_KEY);

    if (RedisModule_CallReplyType(reply) == REDISMODULE_REPLY_STRING)
    {
        int year = atoi(RedisModule_CallReplyStringPtr(reply, NULL)), update = 0;
        is_reset_needed(&year, &update);

        if (update == 1)
        {
            return reset(ctx, &year);
        }
        else
        {
            RedisModule_ReplyWithLongLong(ctx, year);
            return REDISMODULE_OK;
        }
    }

    return RedisModule_ReplyWithError(ctx, "Error occurred when getting value");
}

int RedisModule_OnLoad(RedisModuleCtx *ctx)
{
    if (RedisModule_Init(ctx, YEAR_MODULE, YEAR_MODULE_VERSION, REDISMODULE_APIVER_1) == REDISMODULE_ERR)
    {
        return REDISMODULE_ERR;
    }
    if (RedisModule_CreateCommand(ctx, YEAR_MODULE_COMMAND, command, "WRITE", 0, 0, 0) == REDISMODULE_ERR)
    {
        return REDISMODULE_ERR;
    }
    return REDISMODULE_OK;
}