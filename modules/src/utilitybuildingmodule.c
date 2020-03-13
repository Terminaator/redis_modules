#include "ehrcodemodule.h"

int command(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc != 1)
    {
        return RedisModule_WrongArity(ctx);
    }

    RedisModule_AutoMemory(ctx);

    RedisModuleCallReply *reply = RedisModule_Call(ctx, "HEXISTS", "cc", EHR_CODE_SET_KEY, EHR_CODE_SET_UTILITY_BUILDING_FIELD);

    if (RedisModule_CallReplyType(reply) == REDISMODULE_REPLY_INTEGER)
    {
        if (RedisModule_CallReplyInteger(reply) == 1)
        {
            reply = RedisModule_Call(ctx, "HINCRBY", "ccl", EHR_CODE_SET_KEY, EHR_CODE_SET_UTILITY_BUILDING_FIELD, 1);
            if (RedisModule_CallReplyType(reply) == REDISMODULE_REPLY_INTEGER)
            {
                long long utiliy_building_code = RedisModule_CallReplyInteger(reply);
                if (UTILITY_BUILDING_MIN_VALUE <= utiliy_building_code && utiliy_building_code < UTILITY_BUILDING_MAX_VALUE)
                {
                    RedisModule_ReplyWithLongLong(ctx, utiliy_building_code);
                    return REDISMODULE_OK;
                }
            }
        }
        else
        {
            return RedisModule_ReplyWithError(ctx, "Error occurred when getting value");
        }
    }

    return RedisModule_ReplyWithError(ctx, "Error occurred when getting value");
}

int RedisModule_OnLoad(RedisModuleCtx *ctx)
{
    if (RedisModule_Init(ctx, UTILITY_BUILDING_MODULE, UTILITY_BUILDING_MODULE_VERSION, REDISMODULE_APIVER_1) == REDISMODULE_ERR)
    {
        return REDISMODULE_ERR;
    }

    if (RedisModule_CreateCommand(ctx, UTILITY_BUILDING_MODULE_VERSION_COMMAND, command, "WRITE", 0, 0, 0) == REDISMODULE_ERR)
    {
        return REDISMODULE_ERR;
    }

    return REDISMODULE_OK;
}