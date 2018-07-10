# Package: crypto

## Function Index

- GetMD5FromString(data string) string
- GetMD5FromBytes(data []byte) string
- GetSHA1FromString(data string) string
- GetSHA1FromBytes(data []byte) string
- GetSHA256FromString(data string) string
- GetSHA256FromBytes(data []byte) string
- GenerateRSASSHKeyPair(size int) (pubkey string, privkey string, error)

## Details

### GetMD5FromString

**Author:** ahhh

**Description:** GetMD5FromString takes a string and returns the md5 string of it


**Method Signature:**

```
GetMD5FromString(data string) string
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `data`    | `string`     | the input string to hash                   |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string`     | MD5 value of the input string              |

**Example Usage:**

```

```

### GetMD5FromBytes

**Author:** ahhh

**Description:** GetMD5FromBytes takes a byte array and returns the md5 string of it


**Method Signature:**

```
GetMD5FromBytes(data []byte) string
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `data`    | `[]byte`     | the input data to hash                     |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string`     | MD5 value of the input string              |

**Example Usage:**

```

```

### GetSHA1FromString

**Author:** ahhh

**Description:** GetSHA1FromString takes a string and returns the sha1 string of it


**Method Signature:**

```
GetSHA1FromString(data string) string
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `data`    | `string`     | the input string to hash                   |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string`     | sha1 value of the input string             |

**Example Usage:**

```

```

### GetSHA1FromBytes

**Author:** ahhh

**Description:** GetSHA1FromBytes takes a byte array and returns the sha1 string of it


**Method Signature:**

```
GetSHA1FromBytes(data []byte) string
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `data`    | `[]byte`     | the input data to hash                     |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string`     | SHA1 value of the input string              |

**Example Usage:**

```

```

### GetSHA256FromString

**Author:** ahhh

**Description:** GetSHA256FromString takes a string and returns the sha256 string of it


**Method Signature:**

```
GetSHA256FromString(data string) string
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `data`    | `string`     | the input string to hash                   |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string`     | sha256 value of the input string             |

**Example Usage:**

```

```

### GetSHA256FromBytes

**Author:** ahhh

**Description:** GetSHA256FromBytes takes a byte array and returns the sha256 string of it


**Method Signature:**

```
GetSHA256FromBytes(data []byte) string
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `data`    | `[]byte`     | the input data to hash                     |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string`     | SHA256 value of the input string              |

**Example Usage:**

```

```

### GenerateRSASSHKeyPair

**Author:** ahhh

**Description:** Generates a new rsa key pair the size of the arg, and returns the pub and priv key as strings

**Method Signature:**

```
GenerateRSASSHKeyPair(keySize int) (pubkey string, privkey string, error)
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `keySize` | `int`        | The size of the key to gen                 |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string`     | the generated private key as a string      |
|-----------|--------------|--------------------------------------------|
| `1`       | `string`     | the generated public key as a string       |
|-----------|--------------|--------------------------------------------|
| `2`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```