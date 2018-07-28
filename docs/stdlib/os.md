# Package: os

## Function Index

- TerminateSelf() error
- TerminateVM()

## Details

### TerminateSelf

**Author:** ahhh

**Description:** TerminateSelf will kill the current process, this a dangerous function

**Method Signature:**

```
TerminateSelf() error
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `error`      | (optional) function error                  |

**Example Usage:**

```
if (check == true) {
    TerminateSelf();
}
```
