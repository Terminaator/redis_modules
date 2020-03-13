#include "request.h"

struct string
{
  char *ptr;
  size_t len;
};

struct curl_slist *headers()
{
  struct curl_slist *chunk = NULL;
  chunk = curl_slist_append(chunk, "X-Session-Token: eeee");
  return chunk;
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

const char* get(char *URL)
{
  CURL *curl;
  CURLcode res;

  curl = curl_easy_init();
  if (curl)
  {
    struct string r;
    init_string(&r);

    curl_easy_setopt(curl, CURLOPT_URL, "127.0.0.1:8080/building");
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers());
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, writefunc);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &r);

    res = curl_easy_perform(curl);
    curl_easy_cleanup(curl);

    if (res == CURLE_OK)
    {
      return convert(r.ptr);
    }
  }
  return "error";
}