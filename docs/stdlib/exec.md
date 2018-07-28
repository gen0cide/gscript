# Package: exec

## Function Index

- ExecuteCommand(progname string, args []string) (pid int, stdout string, stderr string, exitCode int, err error)
- ExecuteCommandAsync(progname string, args []string) (proc *exec.Cmd, err error)

## Details

### ExecuteCommand

**Author:** ahhh

**Description:** Executes a single command and waits for it to complete

**Method Signature:**

```
ExecuteCommand(progname string, args []string) (pid int, stdout string, stderr string, exitCode int, err error)
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `command` | `string`     | The target command to run                  |
|-----------|--------------|--------------------------------------------|
| `args`    | `[]string`   | The arguments to command                   |

**Returns:**

| Position  | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `0`       | `int`        | the pid of the executed program            |
| `1`       | `string`     | the stdout of the execute progam           |
| `2`       | `string`     | the stderr of the execute progam           |
| `3`       | `int`        | the exitcode of the executed program       |
| `4`       | `error`      | (optional) function error                  |

**Example Usage:**

```
var exec = G.exec.ExecuteCommand("launchctl", ["submit", "-l", label, "--", name]);
if (exec[4] == null) {
    console.log("Pid: "+exec[0]);;
    console.log("stdout: "+exec[1]);
    console.log("stderr: "+exec[2]);
    console.log("exit code: "+exec[3]);
} else {
    console.log("go errors: "+exec[4].Error())  ;
}   
```

### ExecuteCommandAsync

**Author:** ahhh

**Description:** Executes a single command but does not wait for it to complete

**Method Signature:**

```
ExecuteCommandAsync(progname string, args []string) (proc *exec.Cmd, err error)
```

**Arguments:**

| Label     | Type         | Description                                |
|-----------|--------------|--------------------------------------------|
| `command` | `string`     | The target command to run                  |
| `args`    | `[]string`   | The arguments to command                   |

**Returns:**

| Position  | Type             | Description                                |
|-----------|------------------|--------------------------------------------|
| `0`       | `proc *exec.Cmd` | an  exec/cmd object                        |
| `1`       | `error`          | (optional) function error                  |

**Example Usage:**

```
var fullpath = "/bin/cp"
var running = G.exec.ExecuteCommandAsync(fullpath, ["arg1", "arg2"]);
if (running[1] != null) {
    console.log("errors: "+running[1].Error());
} else {
    console.log("pid: "+running[0].Process.Pid);
}
```
