# Compiler Macros

## Included
### //go_import: 
This command will take a path to a local go library that is discoverable via the GOPATH. This is equivalent to a go import style declaration at the top of a GoLang program. 
### //import:
This macro will take a path to a binary file to be included in the final payload and addressable by the gscript as an asset, using the (AssetFuncs)

## Interpreted
### //priority:
this is a value 0-1000 (defaults 100) that specifies the order that the gscript VM should be executed in.
### //timeout:
This is a value 0-1000 (defaults 30), representing the maximum number of seconds a VM should run for.

## Sugested 
### // title: 
the name of the gscript file 
### // author: 
the person who wrote the gscript
### // purpose: 
a short description of the gscript or how it is intended to be used
### // gscript_version: 
The minimum version of the gscript compiler needed for this gscript
### // ATT&CK:
A link to the MITRE ATT&CK technique the gscript emulates
### // Tactic:
Additional operational advice about the gscript, such as how to trigger an implant or further leverage a technique 
### // Using: 
A link or description to any imported assets that the gscript uses