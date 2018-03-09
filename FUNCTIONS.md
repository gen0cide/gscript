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




