# Implementing Native Code

Native Code in GSCRIPT is when you use the //go_import: macro to reference a native GoLang library in GSCRIPT, allowing you to call those exported GoLang packages and functions directly from GSCRIPT. This allows Go code to be used when Javascript limitations are encountered.

## Converting a Tool Into a GSCRIPT Library

Sometimes you will want to call an existing GoLang library from GSCRIPT, but the tool may need a bit of tweaking to be used by GSCRIPT.

An example of this can be seen in the conversion of GoRedLoot to gloot:

#### GoRedLoot:
- __https://github.com/ahhh/GoRedLoot/blob/master/main.go__

#### Gloot:
- __https://github.com/ahhh/gloot/blob/master/gloot.go__

#### Gscript:
- __https://github.com/ahhh/GSCRIPTs/blob/master/attack/os_x/looter_example.gs__

### Reasons For Converting a Tool Into a GSCRIPT Libary: 

- Adding a non-main package: Many GoLang tools are written inside of main, or in the main package. These tools should be converted into a library so they can be called from GSCRIPT.

- Creating helper functions: GSCRIPT cannot use custom types from native libraries directly, unless a function returns that type. Helper functions could be used to make the type GSCRIPT compatable.

## Writing Your Own Native Library

Sometimes you will need to write your tools in GoLang so that they can be used in GSCRIPT.
