# Auditing Other gscripts

If you are the genesis builder it can help to have a methodolog to quickly audit other peoples gscripts to make sure they run properly in testing and the final binaries.

## Javascript Linter
Setting your text editor / IDE up w/ a javascript linter to parse .gs files will help for quickly spotting syntax errors.

## Build Individual Scripts
Building individual scripts w/ their associated assets against their target platforms is a good way to see if the native libararies and script will even build, or if there is a golang type mismatch, for example.

## Sandbox / Automated Testing
It can help to have a private sandbox or automated testing deployedment, to test final binary builds in a fast and painless way. Testing builds of individual gscripts automatically is a good way to target quick errors w/ the script at runtime. 

