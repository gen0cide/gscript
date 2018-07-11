# Package: rand

## Function Index

- GetInt(min int, max int) int
- GetAlphaNumericString(len int) string
- GetAlphaString(len int) string
- GetAlphaNumericSpecialString(len int) string
- GetBool() bool

## Details

### GetInt

**Author:** ahhh

**Description:**  RandomInt generates a random number between min and max

**Method Signature:**

```
GetInt(min int, max int) int
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `min`     | `int`        | the min range of int to will return        |
| `max`     | `int`        | the max range of int to will return        |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `int`        | a random int between min and max           |

**Example Usage:**

```

```

### GetAlphaNumericString

**Author:** ahhh

**Description:** GetAlphaNumericString Generates a random alpha numeric string of a specified length, all mixed case

**Method Signature:**

```
GetAlphaNumericString(len int) string
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `len`     | `int`        | length of random string to return          |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string`     | a random Alphanumeric string of the length |

**Example Usage:**

```

```

### GetAlphaString

**Author:** ahhh

**Description:** GetAlphaString generates a random alpha string of a specified length, all mixed case

**Method Signature:**

```
GetAlphaString(len int) string
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `len`     | `int`        | length of random string to return          |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string`     | a random alpha string of the length |

**Example Usage:**

```

```

### GetAlphaNumericSpecialString

**Author:** ahhh

**Description:** generates a random alphanumeric and special char string of a specified length

**Method Signature:**

```
GetAlphaNumericSpecialString(len int) string
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `len`     | `int`        | length of random string to return          |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string`     | a random string of the length              |

**Example Usage:**

```

```

### GetBool

**Author:** ahhh

**Description:** generates a random bool, true or false

**Method Signature:**

```
GetBool() bool
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `bool`       | a random bool, true or false               |

**Example Usage:**

```

```