# Implementing Native Code

## Converting a tool to a lib
Sometimes you will want to call an existing GoLang library from gscript, but the tool needs a bit of tweaking to be used in gscript
We can see a great example of this w/ the conversion of GoRedLoot to gloot
### Some Reasons you may want to do this
#### Adding a non-main package
Many golang tools are written inside of main, or in the main package. 
It will help to change these tools to a library so they can be called from GSCRIPT

#### Creating Helper Type Functions
Many golang tools implmenet custom types, unfortuantly GSCRIPT can't use custom types from native libraries directly, unless a function returns that type.

## Writing your own native lib
Sometimes you will need to write your tools in GoLang so that they can be used in GSCRIPT