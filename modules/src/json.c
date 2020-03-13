#include "json.h"

const char* convert(char *res)
{
    struct json_object *parsed_json = json_tokener_parse(res);
    struct json_object *Response;
    json_object_object_get_ex(parsed_json, "Response", &Response);
    return json_object_get_string(Response);
}