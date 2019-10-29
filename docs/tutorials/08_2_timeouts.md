# Timeouts

Timeouts will terminate a gscript at a specified time, in case it's hanging in some way. Unless specified w/ the `//timeouts: ` macro this defaults to `30` which is seconds.

## Macro

One of the most import macros, the `//timeouts: ` macro 
It lets you set the explicitly how long a gscirpt can run for. 
This is often useful for an operator to set to limit user submitted gscripts to a specified runtime.
