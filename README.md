# Web Plucker

`webpluck` scrapes a specific values from a web page. It works as a standalone
binary as well as in a API mode.

## Inputs
`webpluck` takes the following input:
 - URL of the webpage
 - XPATH of the element
 - optional regex to further narrow the selection

and outputs the selected value.

### URL of the webpage
This is the link of the webpage which has the desired information we want to extract.

For example, if we want to scrape the founders of the StackOverflow website from its company page, the URL is:
https://stackoverflow.com/company. The desired value we want to extract is: **Joel Spolsky and Jeff Atwood**

<img width="507" alt="baseUrl" src="https://user-images.githubusercontent.com/14211134/81618604-5335bf00-9405-11ea-8b8c-ddb75e194983.png">

### XPATH of the element
This is the link of the **xpath** of the element on the page that contains the required information. A good way to get this information is to (on a chrome browser):
 - Right click on the place were the information is present
 - Click "Inspect" to open the Chrome developer tools window with the element highligted
 - On the highlighed value in the HTML source code, `Right click -> Copy -> Copy xpath`
 - The copied value is the xpath we need


Get xpath Step 1             |  Get xpath Step 2
:-------------------------:|:-------------------------:
<img width="352" alt="Screen Shot 2020-05-12 at 4 06 13 AM" src="https://user-images.githubusercontent.com/14211134/81619156-8d539080-9406-11ea-99bf-17e9e4da7e87.png" > | <img width="355" alt="Screen Shot 2020-05-12 at 4 08 02 AM" src="https://user-images.githubusercontent.com/14211134/81619157-8e84bd80-9406-11ea-8941-b6c6e0dfab46.png">

The xpath in example above comes out to be: ```//*[@id="content"]/section[3]/ol/li[1]/ol/li[2]/text()```

### `regex` to pluck the right value

Note the the xpath above leads us to the value: *Joel Spolsky and Jeff Atwood launch Stack Overflow*

Since we want to trim that down further, we'll provide a regex value to extract just the names. 

This regex will fetch just the names (the value in parenthesis): 
``` ^(*.) launch .* ```

## Sample hosted invocation

## Sample API invocation

Armed with the knowledge of `baseUrl`, `xpath` and `regex`, we can now call the API endpoint by POSTing these three params:
Example `curl` invocation for the server mode:
```
curl 'https://api.code.express/webpluck/' \
      --data-urlencode 'baseUrl=https://stackoverflow.com/company' \
      --data-urlencode 'xpath=//*[@id="content"]/section[3]/ol/li[1]/ol/li[2]/text()' \
      --data-urlencode 'regex=^(.*) launch .*' -g
```

The result from the API is as follows. The `pluckedData` field returns the value extracted:
```
{
  "baseUrl": "https://stackoverflow.com/company",
  "pluckedData": "Joel Spolsky and Jeff Atwood",
  "regex": "^(.*) launch .*",
  "xpath": "//*[@id=\"content\"]/section[3]/ol/li[1]/ol/li[2]/text()"
```
