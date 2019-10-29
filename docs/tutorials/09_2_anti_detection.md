# Anti-Detection
Good anti-detection should run with a low priority seeking for artifacts that are only available on virtual machines, detection sandboxes, or analyst machines. Once these items are detected the individual gscript can run KillSelf() which looks up its own PID and terminates the process. Through this method you can use multiple gscripts (some of which are anti-detection based) to modularly add these capabilities to a genesis binary and other payloads.

## TerminateSelf()
Calling this function will get the PID of the gscript master process and terminate this process. This itself is a fairly detectable event as opposed to invoking a crash.

Call it like so: 

```sh
G.os.TerminateSelf();
```

### Example
https://github.com/ahhh/gscripts/blob/master/anti-re/sandbox_hostname.gs
