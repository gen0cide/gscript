# Native Functions



## `Asset(assetName)`

Retrieves a packed asset from the VM embedded file store.

### Argument List

 * **assetName** *string*

### Returned Object Fields

 * **fileData** *[]byte*
 * **err** *error*

---



## `RandomString(strlen)`

Generates a random alpha numeric string of a specified length.

### Argument List

 * **strlen** *int64*

### Returned Object Fields

 * **value** *string*

---



## `RandomMixedCaseString(strlen)`

Generates a random mixed case alpha numeric string of a specified length.

### Argument List

 * **strlen** *int64*

### Returned Object Fields

 * **value** *string*

---



## `RandomInt(min, max)`

Generates a random number between min and max.

### Argument List

 * **min** *int64*
 * **max** *int64*

### Returned Object Fields

 * **value** *int*

---



## `XorBytes(aByteArray, bByteArray)`

XOR two byte arrays together.

### Argument List

 * **aByteArray** *[]byte*
 * **bByteArray** *[]byte*

### Returned Object Fields

 * **value** *[]byte*

---



## `StripSpaces(str)`

Strip any unicode characters out of a string.

### Argument List

 * **str** *string*

### Returned Object Fields

 * **value** *string*

---



## `ObfuscateString(str)`

Basic string obfuscator function.

### Argument List

 * **str** *string*

### Returned Object Fields

 * **value** *string*

---



## `DeobfuscateString(str)`

Basic string deobfuscator function.

### Argument List

 * **str** *string*

### Returned Object Fields

 * **value** *string*

---



## `MD5(data)`

Perform an MD5() hash on data.

### Argument List

 * **data** *[]byte*

### Returned Object Fields

 * **value** *string*

---



## `Timestamp()`

Get the system's current timestamp in epoch format.

### Argument List


### Returned Object Fields

 * **value** *int64*

---



## `Halt()`

Stop the current VM from continuing execution.

### Argument List


### Returned Object Fields

 * **value** *bool*

---



## `ExecuteCommand(baseCmd, cmdArgs)`

Executes system commands.

### Argument List

 * **baseCmd** *string*
 * **cmdArgs** *[]string*

### Returned Object Fields

 * **retObject** *VMExecResponse*

---



## `ForkExecuteCommand(baseCmd, cmdArgs)`

Executes system commands via a forked call.

### Argument List

 * **baseCmd** *string*
 * **cmdArgs** *[]string*

### Returned Object Fields

 * **pid** *int*
 * **execError** *error*

---



## `WriteFile(path, fileData, perms)`

Writes data from a byte array to a file with the given permissions.

### Argument List

 * **path** *string*
 * **fileData** *[]byte*
 * **perms** *int64*

### Returned Object Fields

 * **bytesWritten** *int*
 * **fileError** *error*

---



