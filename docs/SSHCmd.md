# `SSHCmd(hostAndPort, cmd, username, password, key)`

Run a command on a target host via SSH.

## Argument List

 * `hostAndPort` (string) - Hostname and port running sshd.
 * `cmd` (string) - Command to run.
 * `username` (string) - User name.
 * `password` (string) - Password.
 * `key` ([]byte) - SSH Key.

## Return Type

 * `obj.response` (string) - Output of command (stdout only!)
 * `obj.runtimeError` (error) - Error.

## Example

```js
var obj = SSHCmd(hostAndPort, cmd, username, password, key)
// obj.response
// obj.runtimeError
```

