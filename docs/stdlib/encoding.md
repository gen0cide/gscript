# Package: encoding

## Function Index

DecodeBase64(data string) (string, error)
EncodeBase64(data string) string
EncodeStringAsBytes(data string) []byte
EncodeBytesAsString(data []byte) string

## Details

### DecodeBase64

**Author:** ahhh

**Description:** decodes a base64 string and returns a string

**Method Signature:**

```
DecodeBase64(data string) string
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `data`    | `string`     | the data to base64 decode                  |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string`     | the decoded string                         |
|-----------|--------------|--------------------------------------------|
| `1`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```

-
### EncodeBase64

**Author:** ahhh

**Description:** encodes a string and returns a base64 string

**Method Signature:**

```
EncodeBase64(data string) string
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `data`    | `string`     | the string to encoded                      |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string`     | the encoded base64 string                  |

**Example Usage:**

```

```

-

### EncodeStringAsBytes

**Author:** ahhh

**Description:** decodes a base64 string and returns a string

**Method Signature:**

```
EncodeStringAsBytes(data string) []bytes
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `data`    | `string`     | the string to bytes                        |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `[]bytes`    | the bytes of the string                    |

**Example Usage:**

```

```

-
### EncodeBytesAsString

**Author:** ahhh

**Description:** decodes a base64 string and returns a string

**Method Signature:**

```
EncodeBytesAsString(data []bytes) string
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `data`    | `[]bytes`    | the data to turn to a string               |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string`     | the string from the bytes                  |

**Example Usage:**

```

```

-