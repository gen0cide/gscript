# Package: file

Handles file operations.

## Function Index
- WriteFileFromBytes(filepath string, data []byte) error
- WriteFileFromString(filepath string, data string) error
- ReadFileAsBytes(filepath string) ([]byte, error)
- ReadFileAsString(filepath string) (string, error)
- AppendBytesToFile(filepath string, data []byte) error
- AppendStringToFile(filepath string, data string) error
- CopyFile(srcpath string, dstpath string, perms string) (bytesWritten int, err error)
- ReplaceInFileWithString(file string, match string, new string) error
- ReplaceInFileWithRegex(file string, regexString string, replaceWith string) error
- SetPerms(filepath string, perms string) error
- CheckExists(targetPath string) bool

## Details

### WriteFileFromBytes

**Author:** ahhh

**Description:** Writes data from a byte array to a file with the parent dirs permissions.

**Method Signature:**

```
WriteFileFromBytes(filepath string, data []byte) error
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `filepath`| `string`     | The location of the file to be written     |
| `data`    | `[]byte`     | The data to be written                     |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `error`      | (optional) function error                  |

**Example Usage:**

```
var myBin = GetAssetAsBytes("example.bin");
errors = G.file.WriteFileFromBytes("example_test", myBin[0]);
console.log("errors: "+errors);
```

### WriteFileFromString

**Author:** ahhh

**Description:** Writes data from a string to a file with the parent dirs permissions.

**Method Signature:**

```
WriteFileFromString(filepath string, data string) error
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `filepath`| `string`     | The location of the file to be written     |
| `data`    | `string`     | The data to be written                     |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `error`      | (optional) function error                  |

**Example Usage:**

```
var writeStringErrors = G.file.WriteFileFromString("example_test", "Example test\n");
console.log("errors: "+ writeStringErrors);
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
var fileBytes = G.file.ReadFileAsBytes("example_test");
console.log("errors: "+ fileBytes[1]);
console.log("file bytes: "+ fileBytes[0]);
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
var readFile = G.file.ReadFileAsString("example_test3");
console.log("errors: "+ readFile[1]);
console.log("file contents: "+readFile[0]);
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
var copyErrors = G.file.CopyFile("example_test", "example_test2");
console.log("errors: " +copyErrors);
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
var myBin = GetAssetAsBytes("example.bin");
var appendedFileError = G.file.AppendFileBytes("example_test", myBin[0]);
console.log("errors: "+ appendedFileError);
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
var appendedFileError = G.file.AppendFileString("example_test", "Appended String\n");
console.log("errors: "+ appendedFileError);
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
var replaceError = G.file.ReplaceInFileWithString("example_test", "test", "replace");
console.log("errors: "+ replaceError);
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
var replaceError = G.file.ReplaceInFileWithRegex("example_test", "(Test)", "replaced");
console.log("errors: "+ replaceError);
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
var permErrors = G.file.SetPerms("example_test", 0777);
console.log("errors: "+permErrors);
```

### CheckExists

**Author:** ahhh

**Description:** Takes a file or directory and checks to see if it exists on the file system


**Method Signature:**

```
CheckExists(targetPath string) bool
```

**Arguments:**

| Label         | Type         | Description                                |
|---------------|--------------|--------------------------------------------|
| `targetPAth`  | `string`     | The location of the file or dir to check   |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `bool`       | if the file or dir exists                  |

**Example Usage:**

```
var exists = G.file.CheckExists("example_test");
console.log("Does it: "+exists);
```
