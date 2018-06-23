### `Signal(pid, signal)`

Send a signal to another process.

##### Argument List

 * `pid` (String) - Process ID you wish to signal
 * `signal` (Integer) - Type of signal you wish to send (9, 15, etc.)

##### Return Type

`boolean` (true = success, false = error)

#Example

This sends a signal to another process. 