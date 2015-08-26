Get

```
{
    "method": "GET",
    "url": "http:192.168.100.122:8010/intersect?type=1&uid=2&targets=3",
    "httpVersion": "HTTP/1.1",
    "queryString": [
        {
            "name": "type",
            "value": "1"
        },
        {
            "name": "uid",
            "value": "2"
        },
        {
            "name": "targets",
            "value": "3"
        }
    ],
    "headers": [
        {
            "name": "Accept",
            "value": "*/*"
        }
    ],
    "cookies": []
}
```

Post

```
{
    "method": "POST",
    "headers": [
        {
            "name": "Content-type",
            "value": "application/json"
        },
        {
            "name": "Accept",
            "value": "*/*"
        }
    ],
    "cookies": [],
    "url": "http:192.168.100.122:8001/user/save",
    "httpVersion": "HTTP/1.1",
    "queryString": [],
    "postData": {
        "mimeType": "application/json",
        "text": "{\"name\":\"sumory\",\"sex\":\"男\",\"age\":12}"
    }
}
```

```
{
    "method": "POST",
    "url": "http:192.168.100.122:8001/user/save",
    "httpVersion": "HTTP/1.1",
    "queryString": [],
    "headers": [
        {
            "name": "Content-type",
            "value": "application/x-www-form-urlencoded"
        },
        {
            "name": "Accept",
            "value": "*/*"
        }
    ],
    "cookies": [],
    "postData": {
        "mimeType": "application/x-www-form-urlencoded",
        "params": [
            {
                "name": "name",
                "value": "ss"
            },
            {
                "name": "age",
                "value": "45"
            }
        ]
    }
}
```