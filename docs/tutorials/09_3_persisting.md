# Persisting 

Leveraging GSCRIPT to persist stored assets is a great way to document persitence techniques in code. 
In this way you can write persistence logic in generic terms and persist arbitrary binaries independent of the persistence technique.

## Persistence Examples

### Windows Persitence example

The following is a simple regkey persistence example, showing how the technique can remain independent of the binary

#### Excerpt

`windows.AddRegKeyString("CURRENT_USER", "Software\\Microsoft\\Windows\\CurrentVersion\\Run", "ExampleExe", fullpath);`

#### Full Example:

__https://github.com/ahhh/gscripts/blob/master/attack/windows/runkey_persistence.gs__


### Linux Persistence Example

The following is a simple sshkey persistence example, showing how the technique can remain independent of the binary

### Excerpt:

```js
G.file.WriteFileFromBytes(myUser[0].HomeDir+ "/.ssh/authorized_keys", pubKey[0]);
``` 

#### Full Example:

__https://github.com/ahhh/gscripts/blob/master/attack/linux/sshkey_persistence.gs__


### MacOS Persistence Example

The following is a simple login hook persistence example, showing how the technique can remain independent of the binary

### Exceprt:

```js
G.exec.ExecuteCommand("defaults", ["write", "com.apple.loginwindow", "LoginHook", name]);
```

### Full Example:

__https://github.com/ahhh/gscripts/blob/master/attack/os_x/loginhook_persistence.gs__

