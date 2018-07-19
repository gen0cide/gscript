# Anti-Detection
Good antidetection should run with a low priority seeking for artifacts that are only avaialble on virtual machines, detection sandboxes, or analyst machines. Once these items are detected the individual gscript can run KillSelf() which looks up its own PID and terminates the process. 

## KillSelf()
Calling this function will get the PID of the gscript master process and terminate this process. This itself is a fairly detectable event as opposed to invoking a crash.



