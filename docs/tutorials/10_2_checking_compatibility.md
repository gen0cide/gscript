# Checking Compatibility
Its important to check the os and arch of the things your including in your gscripts, as you will be building to a single final target in the end.


## Checking Native Libs
Some native libs are not cross platform, this both offers a lot of advantage (for platform specific functions), as well as potential compile and runtime errors
Make sure your native libs can build to your target platform and arch before building your final binaries

## Checking Assets
Assets are ussually less of an issue in terms of gscript errors, however you will want to test all of your assets and make sure they work independantly before putting them in gscript

## Using OS & Arch