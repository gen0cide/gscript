# Checking Compatibility

It is important to check the os and arch of each portion included in your gscripts as you will be building to a single final target in the end.

## Checking Native Libraries

Some native libraries are not cross platform for platform specific functions. Make sure the specified native libraries can build to the target platform and arch before building the final binaries.

## Checking Assets

Assets are usually less of an issue but each asset should be tested and working independently before putting them in gscript.

## Using OS & Arch
