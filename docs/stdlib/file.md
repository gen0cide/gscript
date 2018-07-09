# Package: file

Handles file operations.

## Function Index

AppendBytesToFile(data []byte, filepath string) error
AppendStringToFile(data string, filepath string) error // will not add line breaks, manage yourself
CopyFile(srcpath string, dstpath string, perms string) (bytesWritten int, err error)

WriteFileFromString(data string, filepath string) error
WriteFileFromBytes(data []byte, filepath string) error
SetPerms(filepath string, perms string) error

ReplaceInFileWithString(match string, new string) error
ReplaceInFileWithRegex(regexString string, replaceWith string) error

ReadFileAsString(filepath string) (string, error)
ReadFileAsBytes(filepath string) ([]byte, error)

## Details

### WriteFileFromBytes

**Author:** ahhh

**Description:** Writes data from a byte array to a file with the parent dirs permissions.

**Method Signature:**

```
WriteFileFromBytes(data []byte, filepath string) error
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `destPath`| `string`     | The location of the file to be written     |
|-----------|--------------|--------------------------------------------|
| `fileData`| `[]byte`     | The data to be written                     |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `error`      | (optional) function error                  |

**Example Usage:**

```
// when an asset was embedded into a binary
var asset = G.asset.GetAssetAsBytes("real.txt");
console.log(asset[0]);
// => "Data in real.txt"
console.log(asset[1]);
// => null
G.WriteFileFromBytes(asset[0], "new_real.txt")
```
-

### WriteFileFromString

**Author:** ahhh

**Description:** Writes data from a string to a file with the parent dirs permissions.

**Method Signature:**

```
WriteFileFromBytes(data string, filepath string) error
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `destPath`| `string`     | The location of the file to be written     |
|-----------|--------------|--------------------------------------------|
| `fileData`| `string`     | The data to be written                     |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `error`      | (optional) function error                  |

**Example Usage:**

```
// when an asset was embedded into a binary
var asset = G.asset.GetAssetAsString("real.txt");
console.log(asset[0]);
// => "Contents in real.txt"
console.log(asset[1]);
// => null
G.WriteFileFromString("new_real.txt")
```
-

### ReadFileAsBytes

**Author:** ahhh

**Description:** Reads data from a file and returns its contents as a byte array.

**Method Signature:**

```
ReadFileAsBytes(readPath string) ([]byte, error)
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `readPath`| `string`     | The location of the file to be written     |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `[]bytes`    | contents of the file as a []byte           |
|-----------|--------------|--------------------------------------------|
| `1`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```
-

### ReadFileAsString

**Author:** ahhh

**Description:** Reads data from a file and returns its contents as a string.

**Method Signature:**

```
ReadFileAsString(readPath string) (string, error)
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `readPath`| `string`     | The location of the file to be written     |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `string `    | contents of the file as a strin            |
|-----------|--------------|--------------------------------------------|
| `1`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```
-