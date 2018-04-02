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



## `AppendFileBytes(path, fileData)`

Addes a given byte array to the end of a file

### Argument List

 * **path** *string*
 * **fileData** *[]byte*

### Returned Object Fields

 * **fileError** *error*

---



## `AppendFileString(path, addString)`

Addes a given string to the end of a file

### Argument List

 * **path** *string*
 * **addString** *string*

### Returned Object Fields

 * **fileError** *error*

---



## `Asset(assetName)`

Retrieves a packed asset from the VM embedded file store.

### Argument List

 * **assetName** *string*

### Returned Object Fields

 * **fileData** *[]byte*
 * **err** *error*

---



## `B64Decode(data)`

Base64 decodes a string to a byte array

### Argument List

 * **data** *string*

### Returned Object Fields

 * **value** *[]byte*
 * **execError** *error*

---



## `B64Encode(data)`

Base64 encodes a byte array

### Argument List

 * **data** *[]byte*

### Returned Object Fields

 * **value** *string*

---



## `CheckIfCPUCountIsHigherThanOne()`

Checks to see if the system we are on has at least 2 CPUs (or cores) (sandbox check).

### Argument List


### Returned Object Fields

 * **areWeInASandbox** *bool*

---



## `CheckIfRAMAmountIsBelow1GB()`

Checks to see if the system we are on has at least 1GB of ram (sandbox check).

### Argument List


### Returned Object Fields

 * **areWeInASandbox** *bool*

---



## `CheckIfWineGetUnixFileNameExists()`

Checks for the existance of a stdLib function called "unix_file_name_exists". This is a sandbox check to discover if the current proccess is running within a Wine enviornment

### Argument List


### Returned Object Fields

 * **areWeRunningInWine** *bool*
 * **runtimeError** *error*

---



## `CheckSandboxUsernames()`

Checks to see if the current username contains a series of typical sandbox strings.

### Argument List


### Returned Object Fields

 * **areWeInASandbox** *bool*
 * **runtimeError** *error*

---



## `Chmod(path, perms)`

Change the permissions on a path.

### Argument List

 * **path** *string*
 * **perms** *int64*

### Returned Object Fields

 * **osError** *error*

---



## `CopyFile(srcPath, dstPath, perms)`

Reads the contents of one file and copies it to another with the given permissions.

### Argument List

 * **srcPath** *string*
 * **dstPath** *string*
 * **perms** *int64*

### Returned Object Fields

 * **fileError** *error*

---



## `CreateDir(path)`

Creates a directory at a given path or return an error

### Argument List

 * **path** *string*

### Returned Object Fields

 * **fileError** *error*

---



## `DNSQuestion(target, request)`

Issues a DNS query and returns it's response

### Argument List

 * **target** *string*
 * **request** *string*

### Returned Object Fields

 * **answer** *string*
 * **runtimeError** *error*

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



## `DeleteFile(path)`

Deletes a file at a given path or returns an error

### Argument List

 * **path** *string*

### Returned Object Fields

 * **fileError** *error*

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



## `FileExists(path)`

Checks if a file exists and returns a bool

### Argument List

 * **path** *string*

### Returned Object Fields

 * **fileExists** *bool*
 * **fileError** *error*

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



## `GetLocalIPs()`

Gets an array of Ip addresses for the host

### Argument List


### Returned Object Fields

 * **addresses** *[]string*

---



## `GetMACAddress()`

Gets the MAC address of the interface with an IPv4 address

### Argument List


### Returned Object Fields

 * **address** *string*
 * **runtimeError** *error*

---



## `GetProcName(pid)`

Returns the name of a target proccess

### Argument List

 * **pid** *int*

### Returned Object Fields

 * **procName** *string*
 * **runtimeError** *error*

---



## `HTTPGetFile(url)`

Retrives a file from an HTTP(s) endpoint

### Argument List

 * **url** *string*

### Returned Object Fields

 * **statusCode** *int*
 * **file** *[]byte*
 * **runtimeError** *error*

---



## `Halt()`

Stop the current VM from continuing execution.

### Argument List


### Returned Object Fields

 * **value** *bool*

---



## `InjectIntoProc(shellcode, proccessID)`

Inject shellcode into a provided PID

### Argument List

 * **shellcode** *string*
 * **proccessID** *int64*

### Returned Object Fields

 * **runtimeError** *error*

---



## `InjectIntoSelf(shellcode)`

Inject shellcode into the parrent proccess

### Argument List

 * **shellcode** *string*

### Returned Object Fields

 * **runtimeError** *error*

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



## `IsDebuggerPresent()`

Discovers weather or not the current proccess is being debugged

### Argument List


### Returned Object Fields

 * **isBeingDebugged** *bool*
 * **runtimeError** *error*

---



## `IsTCPPortInUse(port)`

States whether or not a given TCP port is avalible for use

### Argument List

 * **port** *string*

### Returned Object Fields

 * **state** *bool*

---



## `IsUDPPortInUse(port)`

States whether or not a given UDP port is avalible for use

### Argument List

 * **port** *string*

### Returned Object Fields

 * **state** *bool*

---



## `MD5(data)`

Perform an MD5() hash on data.

### Argument List

 * **data** *[]byte*

### Returned Object Fields

 * **value** *string*

---



## `MakeDebuggable()`

Reverse the debugging state set by MakeUnDebuggable by detaching the debugger

### Argument List


### Returned Object Fields

 * **worked** *bool*
 * **runtimeError** *error*

---



## `MakeUnDebuggable()`

Make the current proccess un debuggable by attaching a debugger to itself

### Argument List


### Returned Object Fields

 * **worked** *bool*
 * **runtimeError** *error*

---



## `ModTime(path)`

Retrieves the last modified time of a path.

### Argument List

 * **path** *string*

### Returned Object Fields

 * **modTime** *int64*
 * **fileError** *error*

---



## `ModifyTimestamp(path, accessTime, modifyTime)`

Change the access and modified time of a file.

### Argument List

 * **path** *string*
 * **accessTime** *int64*
 * **modifyTime** *int64*

### Returned Object Fields

 * **fileError** *error*

---



## `ObfuscateString(str)`

Basic string obfuscator function.

### Argument List

 * **str** *string*

### Returned Object Fields

 * **value** *string*

---



## `PostJSON(url, json)`

Transmits a JSON object to a URL and retruns the HTTP status code and response

### Argument List

 * **url** *string*
 * **json** *string*

### Returned Object Fields

 * **statusCode** *int*
 * **response** *[]byte*
 * **runtimeError** *error*

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



## `ReadFile(path)`

Reads a file path and returns a byte array

### Argument List

 * **path** *string*

### Returned Object Fields

 * **fileBytes** *[]byte*
 * **fileError** *error*

---



## `RemoveServiceByName(name)`

Uninstalls a system service

### Argument List

 * **name** *string*

### Returned Object Fields

 * **removalError** *error*

---



## `ReplaceFileString(file, match, replacement)`

Searches a file for a string and replaces each instance found of that string. Returns the amount of strings replaced

### Argument List

 * **file** *string*
 * **match** *string*
 * **replacement** *string*

### Returned Object Fields

 * **stringsReplaced** *int*
 * **fileError** *error*

---



## `RetrievePolymorphicData()`

Retrive data stored within uninitalized space at the end of the gscript binary

### Argument List


### Returned Object Fields

 * **data** *[]byte*
 * **runtimeError** *error*

---



## `RunningProcs()`

Returns an array of int's representing active PIDs currently running

### Argument List


### Returned Object Fields

 * **pids** *[]int*
 * **runtimeError** *error*

---



## `SHA1(data)`

Perform an SHA1() hash on data.

### Argument List

 * **data** *[]byte*

### Returned Object Fields

 * **value** *string*

---



## `SSHCmd(hostAndPort, cmd, username, password, key)`

Runs a command on a target host via SSH and returns the result (stdOut only). Uses both password and key authentication options

### Argument List

 * **hostAndPort** *string*
 * **cmd** *string*
 * **username** *string*
 * **password** *string*
 * **key** *[]byte*

### Returned Object Fields

 * **response** *string*
 * **runtimeError** *error*

---



## `SelfPath()`

Retrieves the path to the currently running executable.

### Argument List


### Returned Object Fields

 * **path** *string*
 * **osError** *error*

---



## `ServePathOverHTTPS(port, path, timeout)`

Starts an HTTPS webserver on a given port (default 443) for $X (default 30) number of seconds that acts as a file server rooted in a given path

### Argument List

 * **port** *string*
 * **path** *string*
 * **timeout** *int64*

### Returned Object Fields

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

 * **startError** *error*

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



## `WritePolymorphicData(data)`

Write data to the uninitalized space at the end of the gscript binary

### Argument List

 * **data** *[]byte*

### Returned Object Fields

 * **runtimeError** *error*

---



## `WriteTempFile(name, fileData)`

Writes data from a byte array to a temporary file and returns the full temp file path and name.

### Argument List

 * **name** *string*
 * **fileData** *[]byte*

### Returned Object Fields

 * **fullpath** *string*
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




