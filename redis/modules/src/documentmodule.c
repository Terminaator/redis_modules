#include "documentmodule.h"

int get_document_number(RedisModuleCtx *ctx, RedisModuleString *doty_key, RedisModuleCallReply *year_reply, RedisModuleCallReply *doty_count_reply)
{
    size_t doty_count_string_lenght;
    const char *doty_count_string = RedisModule_StringPtrLen(RedisModule_CreateStringFromCallReply(doty_count_reply), &doty_count_string_lenght);
    const char *year_string = RedisModule_StringPtrLen(RedisModule_CreateStringFromCallReply(year_reply), NULL);
    const char *doty_string = RedisModule_StringPtrLen(doty_key, NULL);

    char *document = malloc(8 + (doty_count_string_lenght < 5 ? 5 : doty_count_string_lenght) + 1);

    strcpy(document, year_string);
    strcat(document, doty_string);
    strcat(document, "/");

    if (doty_count_string_lenght == 1)
        strcat(document, "0000");
    else if (doty_count_string_lenght == 2)
        strcat(document, "000");
    else if (doty_count_string_lenght == 3)
        strcat(document, "00");
    else if (doty_count_string_lenght == 4)
        strcat(document, "0");

    strcat(document, doty_count_string);

    RedisModule_ReplyWithSimpleString(ctx, document);
    return REDISMODULE_OK;
}

int command(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc != 2)
    {
        return RedisModule_WrongArity(ctx);
    }

    RedisModule_AutoMemory(ctx);

    RedisModuleCallReply *reply = RedisModule_Call(ctx, "HEXISTS", "cs", DOCUMENT_KEY, argv[1]);
    if (RedisModule_CallReplyType(reply) != REDISMODULE_REPLY_ERROR && RedisModule_CallReplyInteger(reply) == 1)
    {
        RedisModuleCallReply *year_reply = RedisModule_Call(ctx, "YEAR", "");
        if (RedisModule_CallReplyType(year_reply) != REDISMODULE_REPLY_INTEGER )
        {
            return RedisModule_ReplyWithCallReply(ctx, year_reply);
        }
        
        RedisModuleCallReply *doty_count_reply = RedisModule_Call(ctx, "HINCRBY", "csc", DOCUMENT_KEY, argv[1], "1");

        if (RedisModule_CallReplyType(doty_count_reply) == REDISMODULE_REPLY_INTEGER)
        {
            return get_document_number(ctx, argv[1], year_reply, doty_count_reply);
        }
    }
    return RedisModule_ReplyWithError(ctx, "Error occured when getting document number");
}

int RedisModule_OnLoad(RedisModuleCtx *ctx)
{
    if (RedisModule_Init(ctx, DOCUMENT_MODULE, DOCUMENT_MODULE_VERSION, REDISMODULE_APIVER_1) == REDISMODULE_ERR)
    {
        return REDISMODULE_ERR;
    }
    if (RedisModule_CreateCommand(ctx, DOCUMENT_MODULE_COMMAND, command, "WRITE", 1, 1, 1) == REDISMODULE_ERR)
    {
        return REDISMODULE_ERR;
    }

    return REDISMODULE_OK;
}