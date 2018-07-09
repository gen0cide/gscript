# Package: net

## Function Index

CheckForInUseTCP(port int) (bool, error)
CheckForInUseUDP(port int) (bool, error)

## Details

### CheckForInUseTCP

**Author:** ahhh

**Description:** 

**Method Signature:**

```
CheckForInUseTCP(port int) (bool, error)
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `port`    | `int`        | The target port to see if it's in use      |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `bool`       | is the port open or closed                 |
|-----------|--------------|--------------------------------------------|
| `1`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```
-

### CheckForInUseUDP

**Author:** ahhh

**Description:** 

**Method Signature:**

```
CheckForInUseUDP(port int) (bool, error)
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `port`    | `int`        | The target port to see if it's in use      |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `bool`       | is the port open or closed                 |
|-----------|--------------|--------------------------------------------|
| `1`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```
-