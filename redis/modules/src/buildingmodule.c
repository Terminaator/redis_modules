#include "ehrcodemodule.h"

int command(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{

    if (argc != 1)
    {
        return RedisModule_WrongArity(ctx);
    }

    RedisModule_AutoMemory(ctx);

    RedisModuleCallReply *reply = RedisModule_Call(ctx, "HEXISTS", "cc", EHR_CODE_SET_KEY, EHR_CODE_SET_BUILDING_FIELD);

    if (RedisModule_CallReplyType(reply) == REDISMODULE_REPLY_INTEGER)
    {
        if (RedisModule_CallReplyInteger(reply) == 1)
        {
            reply = RedisModule_Call(ctx, "HINCRBY", "ccl", EHR_CODE_SET_KEY, EHR_CODE_SET_BUILDING_FIELD, 1);
            if (RedisModule_CallReplyType(reply) == REDISMODULE_REPLY_INTEGER)
            {
                long long building_code = RedisModule_CallReplyInteger(reply);
                if (BUILDING_MIN_VALUE <= building_code && building_code < BUILDING_MAX_VALUE)
                {
                    RedisModule_ReplyWithLongLong(ctx, building_code);
                    return REDISMODULE_OK;
                }
            }
        }
        else
        {
            return RedisModule_ReplyWithError(ctx, "Value not initialized");
        }
    }

    return RedisModule_ReplyWithError(ctx, "Error occurred when getting value");
}

int RedisModule_OnLoad(RedisModuleCtx *ctx)
{
    if (RedisModule_Init(ctx, BUILDING_MODULE, BUILDING_MODULE_VERSION, REDISMODULE_APIVER_1) == REDISMODULE_ERR)
    {
        return REDISMODULE_ERR;
    }
    if (RedisModule_CreateCommand(ctx, BUILDING_MODULE_COMMAND, command, "WRITE", 0, 0, 0) == REDISMODULE_ERR)
    {
        return REDISMODULE_ERR;
    }

    return REDISMODULE_OK;
}