# Native Functions



## `AddRegKeyBinary(registryString, path, name, value)`

Add a binary registry key

### Argument List

 * **registryString** *string*
 * **path** *string*
 * **name** *string*
 * **value** *[]byte*

### Returned Object Fields

 * **runtimeError** *error*

---



## `AddRegKeyDWORD(registryString, path, name, value)`

Add a DWORD registry key

### Argument List

 * **registryString** *string*
 * **path** *string*
 * **name** *string*
 * **value** *uint32*

### Returned Object Fields

 * **runtimeError** *error*

---



## `AddRegKeyExpandedString(registryString, path, name, value)`

Add an expanded string registry key

### Argument List

 * **registryString** *string*
 * **path** *string*
 * **name** *string*
 * **value** *string*

### Returned Object Fields

 * **runtimeError** *error*

---



## `AddRegKeyQWORD(registryString, path, name, value)`

Add a QWORD registry key

### Argument List

 * **registryString** *string*
 * **path** *string*
 * **name** *string*
 * **value** *uint64*

### Returned Object Fields

 * **runtimeError** *error*

---



## `AddRegKeyString(registryString, path, name, value)`

Add a string registry key

### Argument List

 * **registryString** *string*
 * **path** *string*
 * **name** *string*
 * **value** *string*

### Returned Object Fields

 * **runtimeError** *error*

---



## `AddRegKeyStrings(registryString, path, name, value)`

Add a registry key of type string(s)

### Argument List

 * **registryString** *string*
 * **path** *string*
 * **name** *string*
 * **value** *[]string*

### Returned Object Fields

 * **runtimeError** *error*

---



## `Asset(assetName)`

Retrieves a packed asset from the VM embedded file store.

### Argument List

 * **assetName** *string*

### Returned Object Fields

 * **fileData** *[]byte*
 * **err** *error*

---



## `DelRegKey(registryString, path)`

Delete a registry key

### Argument List

 * **registryString** *string*
 * **path** *string*

### Returned Object Fields

 * **runtimeError** *error*

---



## `DelRegKeyValue(registryString, path, value)`

Delete a registry key value

### Argument List

 * **registryString** *string*
 * **path** *string*
 * **value** *string*

### Returned Object Fields

 * **runtimeError** *error*

---



## `DeobfuscateString(str)`

Basic string deobfuscator function.

### Argument List

 * **str** *string*

### Returned Object Fields

 * **value** *string*

---



## `EnvVars()`

Returns a map of enviornment variable names to their corrisponding values.

### Argument List


### Returned Object Fields

 * **vars** *map[string]string*

---



## `ExecuteCommand(baseCmd, cmdArgs)`

Executes system commands.

### Argument List

 * **baseCmd** *string*
 * **cmdArgs** *[]string*

### Returned Object Fields

 * **retObject** *VMExecResponse*

---



## `FindProcByName(procName)`

Returns the Pid of a given proccess, if the proccess can not be found, an error is returned

### Argument List

 * **procName** *string*

### Returned Object Fields

 * **pid** *int*
 * **procError** *error*

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



## `GetEnvVar(vars)`

Returns the value of a given enviornment variable

### Argument List

 * **vars** *string*

### Returned Object Fields

 * **value** *string*

---



## `GetProcName(pid)`

Returns the name of a target proccess

### Argument List

 * **pid** *int*

### Returned Object Fields

 * **procName** *string*
 * **runtimeError** *error*

---



## `Halt()`

Stop the current VM from continuing execution.

### Argument List


### Returned Object Fields

 * **value** *bool*

---



## `InstallSystemService(path, name, displayName, description)`

Installs a target binary as a system service

### Argument List

 * **path** *string*
 * **name** *string*
 * **displayName** *string*
 * **description** *string*

### Returned Object Fields

 * **installError** *error*

---



## `MD5(data)`

Perform an MD5() hash on data.

### Argument List

 * **data** *[]byte*

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



## `QueryRegKey(registryString, path)`

Retrive a registry key

### Argument List

 * **registryString** *string*
 * **path** *string*

### Returned Object Fields

 * **keyObj** *RegistryRetValue*
 * **runtimeError** *error*

---



## `RandomInt(min, max)`

Generates a random number between min and max.

### Argument List

 * **min** *int64*
 * **max** *int64*

### Returned Object Fields

 * **value** *int*

---



## `RandomMixedCaseString(strlen)`

Generates a random mixed case alpha numeric string of a specified length.

### Argument List

 * **strlen** *int64*

### Returned Object Fields

 * **value** *string*

---



## `RandomString(strlen)`

Generates a random alpha numeric string of a specified length.

### Argument List

 * **strlen** *int64*

### Returned Object Fields

 * **value** *string*

---



## `RemoveServiceByName(name)`

Uninstalls a system service

### Argument List

 * **name** *string*

### Returned Object Fields

 * **removealError** *error*

---



## `RunningProcs()`

Returns an array of int's representing active PIDs currently running

### Argument List


### Returned Object Fields

 * **pids** *[]int*
 * **runtimeError** *error*

---



## `Signal(signal, pid)`

Sends a signal to a target proccess

### Argument List

 * **signal** *int*
 * **pid** *int*

### Returned Object Fields

 * **runtimeError** *error*

---



## `StartServiceByName(name)`

Starts a system service

### Argument List

 * **name** *string*

### Returned Object Fields

 * **installError** *error*

---



## `StopServiceByName(name)`

Stops a system service

### Argument List

 * **name** *string*

### Returned Object Fields

 * **installError** *error*

---



## `StripSpaces(str)`

Strip any unicode characters out of a string.

### Argument List

 * **str** *string*

### Returned Object Fields

 * **value** *string*

---



## `Timestamp()`

Get the system's current timestamp in epoch format.

### Argument List


### Returned Object Fields

 * **value** *int64*

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



## `XorBytes(aByteArray, bByteArray)`

XOR two byte arrays together.

### Argument List

 * **aByteArray** *[]byte*
 * **bByteArray** *[]byte*

### Returned Object Fields

 * **value** *[]byte*

---




