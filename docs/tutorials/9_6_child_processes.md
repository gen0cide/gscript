# Child Processes

## Async Exec
When calling G.exec.ExecuteComandAsync it will specifically catch all Signals for SigHup by default. 

### Async Example
var running = G.exec.ExecuteCommandAsync(naming, [""]);
if (running[1] != null) {
    console.log("errors: "+running[1].Error());
} else {
    console.log("pid: "+running[0].Process.Pid);
}
