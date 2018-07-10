# Package: file

Handles file operations.

## Function Index
- WriteFileFromBytes(data []byte, filepath string) error
- WriteFileFromString(data string, filepath string) error
- ReadFileAsBytes(filepath string) ([]byte, error)
- ReadFileAsString(filepath string) (string, error)
- AppendBytesToFile(data []byte, filepath string) error
- AppendStringToFile(data string, filepath string) error // will not add line breaks, manage yourself
- CopyFile(srcpath string, dstpath string, perms string) (bytesWritten int, err error)
- ReplaceInFileWithString(match string, new string) error
- ReplaceInFileWithRegex(regexString string, replaceWith string) error
- SetPerms(filepath string, perms string) error

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
| `1`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```

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
| `1`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```

### CopyFile

**Author:** ahhh

**Description:** Reads data from a file location and copies it to a new location with the original files perms.
Returns the number of bytes copied as an int and an error

**Method Signature:**

```
CopyFile(readPath, destPath string)  (int, error) 
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `readPath`| `string`     | The location of the file to be read        |
| `destPath`| `string`     | The location of the file to be written     |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `int`        | Number of bytes copied                     |
| `1`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```

### AppendFileBytes

**Author:** ahhh

**Description:** Adds data from a byte array to the end of a file, does not add new lines so handle these in your string.

**Method Signature:**

```
AppendFileBytes(targetFile string, data []byte]) error
```

**Arguments:**

| Label        | Type         | Description                                |
|--------------|--------------|--------------------------------------------|
| `targetFile` | `string`     | The location of the file to be written     |
| `data`       | `[]byte`     | The data to be written                     |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```

### AppendFileString

**Author:** ahhh

**Description:** Adds data from a string to the end of a file, does not add new lines so handle these in your string.

**Method Signature:**

```
AppendFileString(targetFile, data string) error
```

**Arguments:**

| Label        | Type         | Description                                |
|--------------|--------------|--------------------------------------------|
| `targetFile` | `string`     | The location of the file to be written     |
| `data`       | `string`     | The data to be written                     |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```

### ReplaceInFileWithString

**Author:** ahhh

**Description:** ReplaceInFileWithString searches a file for a string and replaces each instance found of that string. Returns the amount of strings replaced

**Method Signature:**

```
ReplaceInFileWithString(file, match, replacement string) (int, error)
```

**Arguments:**

| Label         | Type         | Description                                |
|---------------|--------------|--------------------------------------------|
| `file`        | `string`     | The location of the file to be replaced    |
| `match`       | `string`     | The strigns to search and replace          |
| `replacement` | `string`     | The string to replace it with              |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `int`        | Number of strings replaced                 |
| `1`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```

### ReplaceInFileWithRegex

**Author:** ahhh

**Description:** ReplaceInFileWithRegex searches a file for a regex and replaces each instance of that, with your repalce string. Returns the amount of strings replaced

**Method Signature:**

```
ReplaceInFileWithRegex(file, match, replacement string) (int, error)
```

**Arguments:**

| Label         | Type         | Description                                |
|---------------|--------------|--------------------------------------------|
| `file`        | `string`     | The location of the file to be replaced    |
| `match`       | `string`     | the regex to match on strings on           |
| `replacement` | `string`     | The string to replace it with              |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `int`        | Number of strings replaced                 |
| `1`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```

### SetPerms

**Author:** ahhh

**Description:** ReplaceInFileWithRegex searches a file for a regex and replaces each instance of that, with your repalce string. Returns the amount of strings replaced

**Method Signature:**

```
SetPerms(file string, unixPerms int64) (error)
```

**Arguments:**

| Label         | Type         | Description                                |
|---------------|--------------|--------------------------------------------|
| `file`        | `string`     | The location of the file to be replaced    |
| `unixPerms`   | `int64`      | The file perms to set, unix / chmod style  |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `error`      | (optional) function error                  |

**Example Usage:**

```

```
