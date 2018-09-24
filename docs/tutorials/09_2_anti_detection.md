# Anti-Detection
Good antidetection should run with a low priority seeking for artifacts that are only avaialble on virtual machines, detection sandboxes, or analyst machines. Once these items are detected the individual gscript can run KillSelf() which looks up its own PID and terminates the process. In this way you can use multiple gscripts, some of which are anti-detection based, to add these capabilities in a modular way to a genesis binary and other payloads.

## TerminateSelf()
Calling this function will get the PID of the gscript master process and terminate this process. This itself is a fairly detectable event as opposed to invoking a crash.
You call it like so: G.os.TerminateSelf();

### Example
https://github.com/ahhh/gscripts/blob/master/anti-re/sandbox_hostname.gs



