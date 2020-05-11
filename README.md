# Web Plucker

`webpluck` scrapes a specific values from a web page. It works as a standalone
binary as well as in a API mode.

`webpluck` takes the following as input:
 - URL of the webpage
 - XPATH of the element
 - optional regex to further narrow the selection

and outputs the selected value.


Example `curl` invocation for the server mode:
```
curl 'http://localhost:8080/' --data-urlencode "baseUrl=http://example.com" --data-urlencode "xpath=/html/body/div/p[2]/a/@href" --data-urlencode  "regex=^(?:https?://)?(?:[^@\n]+@)?([^:/\n]+)" -vvv
```
