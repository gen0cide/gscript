# Package: requests

## Function Index

PostURL(url string, data string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)
PostJSON(url string, jsondata string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)
PostFile(url string, filepath string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)
GetURLAsString(url string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)
GetURLAsBytes(url string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)

## Details

### PostURL

**Author:** ahhh

**Description:** 

**Method Signature:**

```
PostURL(url string, data string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)
```

**Arguments:**

| Label             | Type                | Description                                |
|-------------------|---------------------|--------------------------------------------|
| `url`             | `string`            | The url to query                           |
|-------------------|---------------------|--------------------------------------------|
| `data`            | `string`            | the data to post to target                 |
|-------------------|---------------------|--------------------------------------------|
| `headers`         | `map[string]string` | the headers to be sent to the target       |
|-------------------|---------------------|--------------------------------------------|
| `ignoresslerrors` | `bool`              | bool to ignore invalid ssl certificates    |

**Returns:**

| Position  | Type             | Description                                |
|-----------|------------------|--------------------------------------------|
| `0`       | `*http.Response` | is the port open or closed                 |
|-----------|------------------|--------------------------------------------|
| `1`       | `string`         | the body of the response as a string       |
|-----------|------------------|--------------------------------------------|
| `2`       | `error`          | (optional) function error                  |

**Example Usage:**

```

```
-

### PostJSON

**Author:** ahhh

**Description:** 

**Method Signature:**

```
PostJSON(url string, jsondata string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)
```

**Arguments:**

| Label             | Type                | Description                                |
|-------------------|---------------------|--------------------------------------------|
| `url`             | `string`            | The url to query                           |
|-------------------|---------------------|--------------------------------------------|
| `json`            | `string`            | The JSON string to send to the target      |
|-------------------|---------------------|--------------------------------------------|
| `headers`         | `map[string]string` | the headers to be sent to the target       |
|-------------------|---------------------|--------------------------------------------|
| `ignoresslerrors` | `bool`              | bool to ignore invalid ssl certificates    |


**Returns:**

| Position  | Type             | Description                                |
|-----------|------------------|--------------------------------------------|
| `0`       | `*http.Response` | is the port open or closed                 |
|-----------|------------------|--------------------------------------------|
| `1`       | `string`         | the body of the response as a string       |
|-----------|------------------|--------------------------------------------|
| `2`       | `error`          | (optional) function error                  |

**Example Usage:**

```

```
-
### PostFile

**Author:** ahhh

**Description:** 

**Method Signature:**

```
PostURL(url string, data string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)
```

**Arguments:**

| Label             | Type                | Description                                |
|-------------------|---------------------|--------------------------------------------|
| `url`             | `string`            | The url to query                           |
|-------------------|---------------------|--------------------------------------------|
| `filepath`        | `string`            | path to the binary to post                 |
|-------------------|---------------------|--------------------------------------------|
| `headers`         | `map[string]string` | the headers to be sent to the target       |
|-------------------|---------------------|--------------------------------------------|
| `ignoresslerrors` | `bool`              | bool to ignore invalid ssl certificates    |

**Returns:**

| Position  | Type             | Description                                |
|-----------|------------------|--------------------------------------------|
| `0`       | `*http.Response` | is the port open or closed                 |
|-----------|------------------|--------------------------------------------|
| `1`       | `string`         | the body of the response as a string       |
|-----------|------------------|--------------------------------------------|
| `2`       | `error`          | (optional) function error                  |

**Example Usage:**

```

```
-
### GetURLAsString

**Author:** ahhh

**Description:** 

**Method Signature:**

```
GetURLAsString(url string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)
```

**Arguments:**

| Label             | Type                | Description                                |
|-------------------|---------------------|--------------------------------------------|
| `url`             | `string`            | The url to query                           |
|-------------------|---------------------|--------------------------------------------|
| `headers`         | `map[string]string` | the headers to be sent to the target       |
|-------------------|---------------------|--------------------------------------------|
| `ignoresslerrors` | `bool`              | bool to ignore invalid ssl certificates    |

**Returns:**

| Position  | Type             | Description                                |
|-----------|------------------|--------------------------------------------|
| `0`       | `*http.Response` | is the port open or closed                 |
|-----------|------------------|--------------------------------------------|
| `1`       | `string`         | the body of the response as a string       |
|-----------|------------------|--------------------------------------------|
| `2`       | `error`          | (optional) function error                  |

**Example Usage:**

```

```
-
### GetURLAsBytes

**Author:** ahhh

**Description:** 

**Method Signature:**

```
GetURLAsBytes(url string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)
```

**Arguments:**

| Label             | Type                | Description                                |
|-------------------|---------------------|--------------------------------------------|
| `url`             | `string`            | The url to query                           |
|-------------------|---------------------|--------------------------------------------|
| `headers`         | `map[string]string` | the headers to be sent to the target       |
|-------------------|---------------------|--------------------------------------------|
| `ignoresslerrors` | `bool`              | bool to ignore invalid ssl certificates    |

**Returns:**

| Position  | Type             | Description                                |
|-----------|------------------|--------------------------------------------|
| `0`       | `*http.Response` | is the port open or closed                 |
|-----------|------------------|--------------------------------------------|
| `1`       | `[]byte`         | the body of the response as a []byte       |
|-----------|------------------|--------------------------------------------|
| `2`       | `error`          | (optional) function error                  |

**Example Usage:**

```

```
-