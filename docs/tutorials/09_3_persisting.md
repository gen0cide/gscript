# Persisting 

Leveraging GSCRIPT to persist stored assets is a great way to document persitence techniques in code. 
In this way you can write persistence logic in generic terms and persist arbitrary binaries independant of the persistence technique.

## Persistence Examples

The following are some simple persistence examples on the main three operating systems GSCRIPT targets.

### Windows Persitence example

The following is a simple regkey persistence example, showing how the technique can remain indpendant of the binary

#### Full Example:

https://github.com/ahhh/gscripts/blob/master/attack/windows/runkey_persistence.gs

#### Excerpt

`windows.AddRegKeyString("CURRENT_USER", "Software\\Microsoft\\Windows\\CurrentVersion\\Run", "ExampleExe", fullpath);`


### Linux Persistence Example

The following is a simple sshkey persistence example, showing how the technique can remain independant of the binary

#### Full Example:

https://github.com/ahhh/gscripts/blob/master/attack/linux/sshkey_persistence.gs

### Excerpt:

```js
G.file.WriteFileFromBytes(myUser[0].HomeDir+ "/.ssh/authorized_keys", pubKey[0]);
``` 


### MacOS Persistence Example

The following is a simple loginhook persistence example, showing how the technique can remain independant of the binary

### Full Example:

https://github.com/ahhh/gscripts/blob/master/attack/os_x/loginhook_persistence.gs

### Exceprt:

```js
G.exec.ExecuteCommand("defaults", ["write", "com.apple.loginwindow", "LoginHook", name]);
```
