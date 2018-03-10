# Function Development

Below is a tutorial on how to develop a VM function in gscript. 

## Location

Functions are defined in `engine/lib_*.go` files. Each "package" has it's own `lib_` file. Please place your methods in the correct folder.

Your function configuration will be placed in `functions.yml` under the corresponding package section.


## Rules / Known Limitations

Functions are defined in `engine/lib_*.go` files. Each "package" has it's own `lib_` file. Please place your methods in the correct folder.

Example: 

 ```
 func (e *Engine) WriteFile(path string, fileData []byte, perms int64) (int, error)
          ^       ^          ^                                           ^
          1       2          3                                           4
 ```
 
 1. Method must be defined as a **pointer receiver** to the `Engine` type.
 2. Method must begin with a capital, and be camel cased appropriately. Follow Golang naming conventions. Your Golang function name is what your function inside the VM will be named.
 3. Properly name your variables - descriptive and unique! Also, do not use any built in Golang reserved keywords or types. Arguments should also be lowercase, but continuing the camelCase convention of Golang.
 4. The function must return something. Whatever that something is should be able to be marshaled into JSON. If your function doesn't return something (even a boolean) you're going to have a bad time.

## Tutorial
For this example, I will implement the `WriteFile()` command.


### 1) Write Golang Function
First, implement your function following the guidelines above.


```
func (e *Engine) WriteFile(path string, fileData []byte, perms int64) (int, error) {
	err := ioutil.WriteFile(path, fileData, os.FileMode(uint32(perms)))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Error writing the file: %s", err.Error())
		return 0, err
	}
	return len(fileData), nil
}
```

Since this has to do with files, we will place it in `engine/lib_file.go`.

**IMPORTANT** - HANDLE ALL YOUR ERRORS. If you don't, you'll break gscript. Make them fail gracefully.

### 2) Write Function Config

Now we will open `functions.yml` and copy the following block into the `file package` section.

```
- name: WriteFile
  description: Writes data from a byte array to a file with the given permissions.
  author: Alex
  package: file
  args:
    - label: path
      gotype: string
    - label: fileData
      gotype: '[]byte'
    - label: perms
      gotype: int64
  returns:
    - label: bytesWritten
      gotype: int
      return: true
    - label: fileError
      gotype: error
      return: true
```

 * `name:` - This is the name of your function as you write it in the Golang lib.
 * `descripton:` - This is what will be generated in the function documentation. Be accurate.
 * `author:` - This is you!
 * `package:` - This is the engine package we determined above. (more around this in development)
 * `args:` - This is where you will document your function arguments.
 	* `label:` - variable name as displayed in the function definition.
 	* `gotype:` - The type in Golang of that argument. Note the VM can be wacky, it tries to do type conversion but sometime it gets it wrong.
 * `returns:` - This is where you'll document the return values of your method.
 	* `label:` - This will be the field assigned to the Javascript return object.
 	* `gotype:` - This is the Golang type of this argument.
 	* `return:` - TRUE or FALSE, should this be included in the return object.
 	
To show you what the return will look like in the VM, here's a brief pseudocode example:

```
var ret = WriteFile("/tmp/tater.jpg", Asset("tater.jpg"), 0644);
// ret.bytesWritten = 488204
// ret.fileError = nil
```

### 3) Build gscript

After saving `engine/lib_file.go` and `functions.yml`, now it's time to generate and build!

In the root of the repo is a script called `build_cli.sh`. Run this to build gscript. If you've got type mismatches or problems with anything in steps 1 or 2, you should get errors here.

Note that generation is *part* of the `build_cli.sh` script. You should see it's output!

### 4) Done!

Now you can look at `FUNCTIONS.md` and see your new documentation rendered. 
