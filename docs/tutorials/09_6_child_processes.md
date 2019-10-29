# Child Processes

## Command Exec 

When calling `G.exec.ExecuteCommand` a child process is spawned that waits for execution to finish before proceding. 

## Async Exec

When calling `G.exec.ExecuteComandAsync` all Signals for SigHup will be caught by default. 

### Async Example

```js
var running = G.exec.ExecuteCommandAsync(naming, [""]);
if (running[1] != null) {
    console.log("errors: "+running[1].Error());
} else {
    console.log("pid: "+running[0].Process.Pid);
}
```
