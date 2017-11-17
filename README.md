z# gscript

Genesis Scripting Engine

![Genesis Logo](http://www.city-church.org.uk/sites/default/files/styles/sidebar_series_image/public/series/image/genesis-in-the-beginning.jpg?itok=EJFz0LWt)

**WARNING**: This library is under active development. API is ***NOT*** stable and will have breaking changes for the foreseeable future.

## Description

GENESIS Scripting (gscript for short) is a technology I've developed to allow dynamic runtime execution of malware installation based on parameters determined at runtime.

Inspiration for this comes from my old AutoRune™ days and from the need for malware to basically become self aware without a bunch of duplicate overhead code.

GScript uses a JS V8 Virtual Machine to interpret your genesis script and allow it to hook into the malware initialization.

The Engine itself is referred commonly as "GSE" - Genesis Scripting Engine.

## What is GENESIS btw?

GENESIS was created by @vyrus, @gen0cide, @emperorcow, and @ahhh for dynamically bundling multiple payloads into one dropper for faster deployment of implants for the CCDC Red Team.

For more information on this work we do every year, see my blog post outlining our toolbox:

 * <https://alexlevinson.wordpress.com/2017/05/09/know-your-opponent-my-ccdc-toolbox/>

GSE's goal is to allow intelligent deployment of those payloads.

## Variables

### User Defined (You define/overwrite as needed)

| Variable Name  | Type     | Default              | Purpose                                                                                                   |
|----------------|----------|----------------------|-----------------------------------------------------------------------------------------------------------|
| `source_url`   | `string` | `null`               | Location GSE should download implant from.                                                                |
| `packed_file`  | `string` | `null`               | A packed GSE File (created with the GSE compiler)                                                         |
| `source_bytes` | `array`  | Generated At Runtime | The contents to be written to `file_dest` during `Deploy()`                                               |
| `file_dest`    | `string` | `null`               | The location where `source_bytes` get's written to, and subsequently executed.                            |
| `exec_args`    | `array`  | `[]`                 | Array of strings that will be passed to the `dest_file` Exec() call.                                      |
| `timeout`      | `int`    | `180000`             | The global timeout for the entire GSE VM. Default = 3 minutes. Can be overwritten to shorten or lengthen. |

### Read Only (Defined at Runtime by GSE)

| Variable Name | Type     | Example                                                 | Purpose                                                                                                      |
|---------------|----------|---------------------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| `user_info`   | `object` | `{uid: 0, gid: 0, username: "root", home_dir: "/root"}` | Information about the User. Will be basically whatever is returned with https://golang.org/pkg/os/user/#User |
| `hostname`    | `string` | `example01`                                             | The hostname of the machine.                                                                                 |
| `ip_addrs`    | `array`  | `["127.0.0.1","192.168.1.5"]`                           | The IP addresses of the machine.                                                                             |
| `os`          | `string` | `linux`                                                 | The operating system (basically `runtime.GOOS`)                                                              |
| `arch`        | `string` | `amd64`                                                 | The CPU architecture (basically `runtime.GOARCH`)                                                             |

## Builtin Functions
These functions are available to you automatically within the GSE scripting context.



#### Halt()

Terminates the current GSE VM gracefully.

##### Argument List

None

##### Return Type

`boolean` (true = success, false = error)

---

#### DeleteFile(path)

Delete the file located at `path`.

##### Argument List

 * `path` (String) - Path to file you wish to delete.

##### Return Type

`boolean` (true = success, false = error)

---

#### CopyFile(srcPath, dstPath)

Copy file from `srcPath` to `dstPath`.

##### Argument List

 * `srcPath` (String) - Path to source file.
 * `dstFile` (String) - Path to destination file.

##### Return Type

`boolean` (true = success, false = error)

---

#### WriteFile(path, bytes, perms)

Write `bytes` to `path` and set perms to `perms`.

##### Argument List

 * `path` (String) - Path to file you wish to write.
 * `bytes` (Array) - Array of bytes you wish to write to the `path` location.
 * `perms` (String) - Octal unix permissions represented as a string. ie: `0777`.

##### Return Type

`boolean` (true = success, false = error)

---

#### ExecuteFile(path, args)

Execute a file located at `path` with `args` as arguments.

##### Argument List

 * `path` (String) - Path to file you wish to execute.
 * `args` (Array) - Arguments to pass to during file execution.

##### Return Type

`boolean` (true = success, false = error)

---

#### AppendFile(path, bytes)

Append `bytes` to the file located at `path`.

##### Argument List

 * `path` (String) - Path to file you wish to append.
 * `bytes` (Array) - Array of bytes you wish to append.

##### Return Type

`boolean` (true = success, false = error)

---

#### ReplaceInFile(path, target, replace)

Replace any instances of `target` with `replace` in the file located at `path`.

##### Argument List

 * `path` (String) - Path to file you wish to modify.
 * `target` (String) - String value you wish to replace in the file.
 * `replace` (String) - String value you wish to substitute `target` with.

##### Return Type

`boolean` (true = success, false = error)

---

#### Signal(pid, signal)

Send a signal to another process.

##### Argument List

 * `pid` (String) - Process ID you wish to signal
 * `signal` (Integer) - Type of signal you wish to send (9, 15, etc.)

##### Return Type

`boolean` (true = success, false = error)

---

#### RetrieveFileFromURL(url)

Retrieve a file via `GET` for a given `url`.

##### Argument List

 * `url` (String) - Full URL of location you wish to retrieve.

##### Return Type

`[]bytes` - Byte array of body response.

---

#### DNSQuery(question, type)

Perform a DNS lookup.

##### Argument List

 * `question` (String) - DNS query question (eg: "twitter.com")
 * `type` (String) - DNS question type (A, CNAME, MX, etc.)

##### Return Type

`Object` - Reference `VMDNSQueryResponse` in `response_objects.go` for object details.

---

#### HTTPRequest(method, url, body, headers)

Perform an HTTP/S request.

##### Argument List

 * `method` (String) - HTTP Method (GET, POST, PUT, DELETE, HEAD, etc.)
 * `url` (String) - Full URL (including https://) you wish to make a request to.
 * `body` (String) - Any body you wish to include (nil if none).
 * `headers` (Object) - A key/value object that will be set as HTTP Request Headers.

##### Return Type

`Object` - Reference `VMHTTPRequestResponse` in `response_objects.go` for object details.

---

#### Exec(cmd, args)

Execute the given command and arguments.

##### Argument List

 * `cmd` (String) - Base command you wish to run
 * `args` (Array) - Arguments as an array of strings.

##### Return Type

`Object` - Reference `VMExecResponse` in `response_objects.go` for object details.

---

#### MD5(bytes)

Create a MD5 hash of the given bytes.

##### Argument List

 * `path` (Array) - Array of bytes.

##### Return Type

`string` - Hex encoded MD5 hash.

---

#### SHA1(bytes)

Create a SHA1 hash of the given bytes.

##### Argument List

 * `bytes` (Array) - Array of bytes.

##### Return Type

`string` - Hex encoded SHA1 hash.

---

#### B64Encode(bytes)

Perform a Base64 encode on `bytes`.

##### Argument List

 * `bytes` (Array) - Array of bytes you wish to base64 encode.

##### Return Type

`string` - Base64 encoded string representation.

---

#### B64Decode(string)

Perform a Base64 decode on `string`.

##### Argument List

 * `string` (String) - Base64 encoded string

##### Return Type

`[]bytes` - Byte array of the deserialized b64 string.

---

#### Timestamp()

Get current time in Epoch.

##### Argument List

None

##### Return Type

`integer` - Current time in Epoch.

---

#### CPUStats()

Retreive specs about the machine's CPU.

##### Argument List

None

##### Return Type

`Object` - Reference `VMCPUStatsResponse` in `response_objects.go` for object details.

---

#### MemStats()

Retreive specs about the machine's memory.

##### Argument List

None

##### Return Type

`Object` - Reference `VMMemStatsResponse` in `response_objects.go` for object details.

---

#### SSHExec(host, port, creds, cmds)

Executes SSH commands on the given host.

##### Argument List

 * `host` (String) - Host you wish to connect to.
 * `port` (String) - Port you wish to connect to.
 * `creds` (Object) - Credential Object: `{ username: "", password: "", privateKey: "" }`
 * `cmds` (Array) - Commands you wish to run as an array of strings.

##### Return Type

`Object` - Reference `VMSSHExecResponse` in `response_objects.go` for object details.

---

#### Sleep(seconds)

Sleep for `seconds` number of seconds.

##### Argument List

 * `seconds` (Int) - Sleep Duration

##### Return Type

`boolean` (true = success, false = error)

---

#### GetDirsInPath()

Get a list of all the directories currently in our PATH.

##### Argument List

None

##### Return Type

`[]string` - Array of directories in the current PATH as strings.

---

#### EnvVars()

Retrieve an array of all environment variables in the current execution.

##### Argument List

None

##### Return Type

`Object` - Reference `VMEnvVarsResponse` in `response_objects.go` for object details.

---

#### GetEnv(varname)

Retrieve the value for Environment Variable `varname`.

##### Argument List

 * `varname` (String) - Environment variable name.

##### Return Type

`string` - Value, empty if undefined.

---

#### FileCreateTime(path)

Lookup the creation time for file located at `path`.

##### Argument List

 * `path` (String) - Path to target file.

##### Return Type

`int` - Last modified time in Epoch format.

---

#### FileModifyTime(path)

Lookup the last modified time for file located at `path`.

##### Argument List

 * `path` (String) - Path to target file.

##### Return Type

`int` - Last modified time in Epoch format.

---

#### LoggedInUsers()

Gets an array of unique users currently logged in.

##### Argument List

None

##### Return Type

`[]string` - Array of usernames as strings.

---

#### UsersRunningProcs()

Gets an array of unique users currently running processes.

##### Argument List

None

##### Return Type

`[]string` - Array of usernames as strings.

---

#### ServeDataOverHTTP(data, port, timeout)

Starts an HTTPServer that will respond to `GET /` with the `data` provided on port `port`.

##### Argument List

 * `data` (String) - Data you wish to serve.
 * `port` (Int) - What port should we listen on?
 * `timeout` (Int) - How many seconds should we listen? (Cannot be > global `timeout` variable!)

##### Return Type

`boolean` (true = success, false = error)

---

## Notes

This is just my design chicken scratch. I'll slowly migrate this stuff over to more formal documentation as I implement.

```

// Genesis Hooks (Can be user defined)
function BeforeDeploy() {};
function Deploy()       {};
function AfterDeploy()  {};
function OnError()      {};

// Research Functions (Should *not* be overridden!)
// These functions allow you to get information about
// a given system in order to make decisions based off
// the context the runtime is executing in.
function LocalUserExists(username)      { return boolean; };
function ProcExistsWithName(name)       { return boolean; };
function CanReadFile(path)              { return boolean; };
function CanWriteFile(path)             { return boolean; };
function CanExecFile(path)              { return boolean; };
function FileExists(path)               { return boolean; };
function DirExists(path)                { return boolean; };
function FileContains(path, match)      { return boolean; };
function IsVM()                         { return boolean; };
function IsAWS()                        { return boolean; };
function HasPublicIP()                  { return boolean; };
function CanMakeTCPConn(dst, port)      { return boolean; };
function ExpectedDNS(query, type, resp) { return boolean; };
function CanMakeHTTPConn(url)           { return boolean; };
function DetectSSLMITM(url, cert_fp)    { return boolean; };
function CmdSuccessful(cmd)             { return boolean; };
function CanPing(dst)                   { return boolean; };
function TCPPortInUse(port)             { return boolean; };
function UDPPortInUse(port)             { return boolean; };
function ExistsInPath(progname)         { return boolean; };
function CanSudo()                      { return boolean; };
function Matches(string, match)         { return boolean; };
function CanSSHLogin(ip, port, u, p)    { return boolean; };

```

## TODO

 * Implement All Functions
 * Implement Global Variables
 * Implement Runtime Variable Loading
 * Implement Timeout
 * Implement Hook Callers
 * Implement CLI Framework (cmd/gscript)
 * Implement Compiler / Crypter

## Credits

Shoutouts to the homies:

 * vyrus
 * ahhh
 * cmccsec
 * carnal0wnage
 * indi303
 * emperorcow
 * rossja
