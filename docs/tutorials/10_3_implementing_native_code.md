# Implementing Native Code
Native Code in gscript is when you use the //go_import: macro to reference a native GoLang library in gscript, allowing you to call those exported golang packages and functions directly from gscript, and do powerfull thigns you otherwise couldn't do in JavaScript.

## Converting a tool to a lib
Sometimes you will want to call an existing GoLang library from gscript, but the tool needs a bit of tweaking to be used in gscript
We can see a great example of this w/ the conversion of GoRedLoot to gloot

GoRedLoot:
https://github.com/ahhh/GoRedLoot/blob/master/main.go
Gloot:
https://github.com/ahhh/gloot/blob/master/gloot.go
Gscript:
https://github.com/ahhh/gscripts/blob/master/attack/os_x/looter_example.gs

### Some Reasons you may want to do this
#### Adding a non-main package
Many golang tools are written inside of main, or in the main package. 
It will help to change these tools to a library so they can be called from GSCRIPT

#### Creating Helper Type Functions
Many golang tools implmenet custom types, unfortuantly GSCRIPT can't use custom types from native libraries directly, unless a function returns that type.

## Writing your own native lib
Sometimes you will need to write your tools in GoLang so that they can be used in GSCRIPT