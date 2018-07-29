# Multiple Scripts

One of the massive powers of GSCIRPT is that you can compile any number of scripts into a single binary. This means individual gscripts are often written as atomic as possible, ussually accomplishing a single task. These multiple scripts are then selectivly compiled into a single binary based on the techniques you desire.

## N Number of Scripts

The final argument gscript takes is a path to a script, or any number of paths to a script. This must be the final command line argument. This can include wildcards or any type of path expansion in bash. 

## Multiple Script Examples

```sh
gscript c /path/to/gscript.gs /path/to/another_gscript.gs /path/to/a_last_gscript.gs
```

```sh
gscript c /path/to/*t.gs
```
