#include "mainmodule.h"

int command(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{

    if (argc != 1)
    {
        return RedisModule_WrongArity(ctx);
    }

    RedisModule_AutoMemory(ctx);

    RedisModuleCallReply *reply = RedisModule_Call(ctx, "EXISTS", "c", PROCEDURE_KEY);
    if (RedisModule_CallReplyType(reply) == REDISMODULE_REPLY_INTEGER && RedisModule_CallReplyInteger(reply) == 1)
    {
        reply = RedisModule_Call(ctx, "INCR", "c", PROCEDURE_KEY);
        if (RedisModule_CallReplyType(reply) == REDISMODULE_REPLY_INTEGER)
        {
            RedisModule_ReplyWithCallReply(ctx, reply);
            return REDISMODULE_OK;
        }
    }
    else if (RedisModule_CallReplyType(reply) == REDISMODULE_REPLY_INTEGER && RedisModule_CallReplyInteger(reply) == 0)
    {
        return RedisModule_ReplyWithError(ctx, "Error occurred when getting value");
    }

    return RedisModule_ReplyWithError(ctx, "Error occurred when getting value");
}

int RedisModule_OnLoad(RedisModuleCtx *ctx)
{
    if (RedisModule_Init(ctx, PROCEDURE_MODULE, PROCEDURE_MODULE_VERSION, REDISMODULE_APIVER_1) == REDISMODULE_ERR)
    {
        return REDISMODULE_ERR;
    }
    if (RedisModule_CreateCommand(ctx, PROCEDURE_MODULE_VERSION_COMMAND, command, "WRITE", 0, 0, 0) == REDISMODULE_ERR)
    {
        return REDISMODULE_ERR;
    }

    return REDISMODULE_OK;
}