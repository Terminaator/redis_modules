#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <curl/curl.h>
#include <json-c/json.h>

struct string
{
  char *ptr;
  size_t len;
};

void init_string(struct string *s)
{
  s->len = 0;
  s->ptr = malloc(s->len + 1);
  if (s->ptr == NULL)
  {
    fprintf(stderr, "malloc() failed\n");
    exit(EXIT_FAILURE);
  }
  s->ptr[0] = '\0';
}

size_t writefunc(void *ptr, size_t size, size_t nmemb, struct string *s)
{
  size_t new_len = s->len + size * nmemb;
  s->ptr = realloc(s->ptr, new_len + 1);
  if (s->ptr == NULL)
  {
    fprintf(stderr, "realloc() failed\n");
    exit(EXIT_FAILURE);
  }
  memcpy(s->ptr + s->len, ptr, size * nmemb);
  s->ptr[new_len] = '\0';
  s->len = new_len;

  return size * nmemb;
}

const char* get()
{
  CURL *curl;
  CURLcode res;
  struct json_object *parsed_json;
  struct json_object *Respond;

  curl = curl_easy_init();
  struct string s;
  init_string(&s);
  if (curl)
  {
    struct curl_slist *chunk = NULL;
    chunk = curl_slist_append(chunk, "X-Session-Token: eeee");

    curl_easy_setopt(curl, CURLOPT_URL, "127.0.0.1:8080/building");
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, chunk);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, writefunc);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &s);
    res = curl_easy_perform(curl);

    printf("%s\n", s.ptr);
    parsed_json = json_tokener_parse(s.ptr);
    json_object_object_get_ex(parsed_json, "Respond", &Respond);
    free(s.ptr);

    /* always cleanup */
    curl_easy_cleanup(curl);
    return json_object_get_string(Respond);
  }
  return s.ptr;
}

int main()
{
  const char* a = get();
  printf("%s\n", a);
}