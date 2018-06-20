
# !!! WARNING - BREAKING CHANGES NOTICE !!!

<b>A new version of gscript will be released on August 10th, 2018 during a presentation at the Defcon Security Conference.<br/>Expect MAJOR breaking API changes and other instabilities before then. Official v1 release will follow semantic versioning guidelines from that point on.</b>

# gscript

Genesis Scripting Engine

<div align="center"><img src="http://www.city-church.org.uk/sites/default/files/styles/sidebar_series_image/public/series/image/genesis-in-the-beginning.jpg?itok=EJFz0LWt" alt="Genesis Logo"/></div>

## Description

Genesis Scripting (gscript for short) is a technology I've developed to enable more intelligent malware stagers. Typically, stagers are pretty dumb. Most stagers are unique to the malware they deploy and do not allow for "bundling" of multiple payloads. Sophisticated attackers *do* in fact bundle their payloads, which makes runtime uncertainty even more assured.

GScript changes that. GScript allows for dynamic execution logic *per payload*. This is achieved by a modified Javascript runtime that is statically embedded in the final stager binary. This runtime/virtual machine runs "hook" functions that you've defined in your script, checking to ensure the script wishes to proceed after each hook.

This has significant benefits over traditional tactics:

 * Scripts are far more "sandboxed" from each other. If you're bundling 10 payloads and 1 of them has a syntax error in its script, with gscript, only that scripts VM dies, not the entire program.
 * GScript's VM, while sandboxed, has native hooks injected into it. This allows the VM to interact with things outside of the VM (filesystem, network, registry, etc.).
 * These functions are by and large, completely cross platform. This allows someone to only learn GScript to write scripts without having to learn a different programming language.
 * Execution is also parallelized using the very effective go routine paradigm, resulting in much faster execution of stagers with multiple payloads.
 
This development process is incredibly efficient with our gscript CLI utility.

The VM's custom runtime (Engine) is referred commonly as "GSE" - Genesis Scripting Engine.

## Architecture

### Compiler

The compiler is what translates your gscripts and their assets into a finalized binary. Some features of the compiler:

 * Support native binary compilation for all major operating systems: Windows, Linux, OS X
 * Can support large numbers of scripts and assets into a single executable.
 * Built-in lossless compression and obfuscation of both scripts and embedded assets.
 * **VERY FAST**. Compilation times generally less than 5 seconds.
 * Post compilation obfuscation to remove references to the library.
 * Defaults to a null logger for the final binary (no output ever!), but can be overridden to inject a development logger into the final binary for testing.

### VM Engine

The final binary contains the gscript engine with all scripts and their imported assets. It will initialize VMs, one for each script, and execute them generally in parallel (theres priority overrides, but more on that below!). 

The VMs cannot interact with one another and errors in one VM will be gracefully handled by the engine. This prohibits one VM from causing instability or fatal errors preventing other scripts from executing.

![GScript Architecture Diagram](https://i.imgur.com/C1avFt7.png)

The Engine has been designed to be lean and free from bloated imports. It's come a long way, but there will be more improvements here in the future as well.

#### Execution Flow

The VM expects your scripts to implement three functions (called "hooks" in our documentation):

 * `BeforeDeploy()` - *intended to function as the reconnaissance function to determine if the script should attempt to deploy it's second stage.*
 * `Deploy()` - *contains the stage two deployment logic.*
 * `AfterDeploy()` - *allows the script to clean up after a successful deployment.*
 
If any of these functions returns `false`, the engine will prevent subsequent execution of the functions proceeding for that script.

When either hooks return false or all hooks run successfully, the VM will destroy itself and notify the scheduling engine that the VM is finished executing.

## CLI

The gscript command line utility is the primary interface for writing, testing, debugging, and compiling your Genesis binaries.

```
NAME:
   gscript - Command Line application for interacting with the GENESIS Scripting Engine.

USAGE:
   gscript [global options] command [command options] [arguments...]

COMMANDS:
     compile, c  Compile genesis scripts into a stand alone binary.
     new, n      Writes an example gscript to either the given path or STDOUT.
     run, r      Run a Genesis script locally (Careful not to infect yourself).
     shell, s    Run an interactive gscript REPL.
     test, t     Check a GSE script for syntax errors.
     update, u   Update the gscript CLI binary to the latest version.
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d    Enable verbose output of runtime debugging logs.
   --help, -h     show help
   --version, -v  print the version
```

You can use `-h` on any subcommand as well to see subcommand specific options.

## Installation

You'll need the following dependencies installed and configured on your system. Regardless of the target OS, we recommend (and generally only test) the CLI utility on Linux/OSX based systems. You can still compile binaries **for**  Windows though!

 * Golang v1.10 or higher
 
After you have both of those, run:

```
$ go get github.com/gen0cide/gscript/cmd/gscript
```


## Scripting

Writing a Genesis script is simple and easy to learn. A few things to know about writing genesis scripts:

 * A genesis script (.gs) is written in Javascript. This allows any editors or IDE that can work with Javascript to also apply to genesis scripts.
 * The VM targets ES5, so ES6 syntax will not be valid.
 * You can define your own functions and variables within your script.
 * GScript injects a standard library of functions into the VM to facilitate interaction outside of the Javascript sandbox. These functions, while called from Javascript, are actually implemented natively in code. For more information about function development, take a look at [DEVELOPERS.md](DEVELOPERS.md).
 * Logging functions (LogInfo, LogError, etc.) will log output when using the `run` subcommand, or when a final binary is compiled with the flag `--enable-logging`. You do not have to remove your logging statements to keep your final binary from logging output - the compiler does this by default.

### Basic Example

In the Execution Flow section, we covered the "hook" functions a gscript needs to implement. Here is the most basic example of a valid Genesis script:

```
// A standard javascript comment. Totally valid!

function BeforeDeploy() {
	return true;
}

function Deploy() {
	return true;
}

function AfterDeploy() {
	return true;
}
```

And that's it! The only requirement is that you implement those three functions and return `true` when you wish your script to continue to the next function, or `false` when you do not want the VM to continue executing this script.

Now let's use some of the built in functions:

```
// A standard javascript comment. Totally valid!

// We can declare our own globals
var val = null;

// As well as our own custom functions.
function Greet(name) {
	return "Hello, " + name;
}

function BeforeDeploy() {
	// LogInfo is a VM function. Also: LogDebug(), LogWarn(), LogError()
	LogInfo("*** EXAMPLE BEFORE DEPLOY ***");
	
	// HOSTNAME is a constant injected into the VM at runtime by the Engine
	// See the section below on "Embedded Variables" for more info.
	val = HOSTNAME;
	return true;
}

function Deploy() {
	// We can reference our val variable within these functions.
	LogInfo(Greet(val));
	return true;
}

function AfterDeploy() {
	LogInfo("*** EXAMPLE AFTER DEPLOY ***");
	return true;
}
```

Now let's talk about test and run our script! 

#### Testing your script

I saved the above script to `example.gs` and ran the command `gscript test example.gs`:

![Testing Script Example](https://i.imgur.com/y6AjwG5.png)

#### Running your script
Now what about if we want to execute the script? Let's use the `run` subcommand (`gscript run example.gs`)

![Running Script Example](https://i.imgur.com/H8TcOwV.png)


#### Debugging your script

Lets say there was an unexpected error and I wanted to investigate. You can use the `DebugConsole()` function (only enabled with `run`, not in a final binary) to drop into an interactive gscript REPL.


I've modified the `Deploy()` function of our script to look like this:

```
function Deploy() {
	DebugConsole();

	// We can reference our val variable within these functions.
	LogInfo(Greet(val));
	return true;
}
```

Now when I re-run our script with the same command above, a breakpoint is set and I can now use the interactive debugger:

![Debugging Script Example](https://i.imgur.com/2TUZpiT.png)

Notice the auto complete will hint at special VM functions that GScript has injected into the Javascript VM.

#### Compiling your script

So you've done everything except build a final binary. Let's walk you through what to do here.

We will use the following command:

```
$ gscript compile --outfile ./example.bin example.gs
```
If you do not specify the `--os` or `--arch` flag, GScript will default to whatever your operating system and architecture is. Here's the results:

![Compiling Script Example](https://i.imgur.com/PvXXK4c.png)

And now if I run our binary, what do you think we'll see?

![Normal Execution Script Example](https://i.imgur.com/wg6Fh9z.png)

Aha! Where's all our output? Remember, logging is **DISABLED** by default in a final binary. Let's add some debug logging into the compiled binary:

```
$ gscript compile --outfile ./example2.bin --enable-logging example.gs
```

![Compiling With Logging Example](https://i.imgur.com/uyK4VtN.png)

And now if I run the binary:

![Running Binary with Logging Example](https://i.imgur.com/hdY1UyK.png)

Notice the original `example.bin` did nothing, while the `example2.bin` allowed you to see execution output.

You can also combine multiple scripts into the same compiled binary. Here's an example of the syntax:

```
$ gscript compile --outfile /tmp/foo.bin a.gs b.gs c.gs
```

This will embed all three GScripts into the final binary, isolate them from one another, and parallelize their execution.

Next, we will explore how to embed files into your binary and look at a few new VM functions (Full VM function documentation is linked to in the **VM Functions** section).

### Compiler Macros

GScript has it's own compiler macro implementation. This lets you add syntactically valid Javascript statements while extending the functionality of the VM. The most common use for this is embedding files (also called assets).

I've copied a file "tater.jpg" into our example directory. In addition, I've noted our current working directory as well as taken an MD5 of the image.

![Asset Example](https://i.imgur.com/CRVUOHD.png)

With the asset in the directory, I'm going to write a new genesis script to `example-asset.gs`:

```
// A standard javascript comment. Totally valid!

// Adding tater.jpg to this file
//import:/tmp/gscript_demo/tater.jpg
// ^ Example usage of the "import" compiler macro

function BeforeDeploy() {
	return true;
}

function Deploy() {
	// We will use a new function Timestamp() to get a unique filename.
	var filename = "/tmp/gscript_demo/" + Timestamp().value + ".jpg";	
	
	// WriteFile() is part of the file package in the standard library.
	// Asset() is how you retrieve your asset's data (core package).
	WriteFile(filename, Asset("tater.jpg").fileData, 0644);

	// Still need to return true!
	return true;
}

function AfterDeploy() {
	return true;
}
```

Line 4 is the compiler macro:

```
//import:/tmp/gscript_demo/tater.jpg
```

When you either perform a `run` or a `compile`, the Engine will embed "tater.jpg" in the runtime (retrieving with Asset function). When compiled, this asset is embedded within the binary itself. Now lets compile and run it!

![Compiling an imported asset](https://i.imgur.com/kbS8DlL.png)

I've highlighted the line where the compiler found the asset. Now if I run it, I will see the new JPG file in my directory. An MD5 of the files show them to be identical.

![Embedded Asset Final Binary Example](https://i.imgur.com/xIsmSEM.png)

Some final notes about file assets:

 * You can embed multiple assets per script, just make sure they have unique filenames.
 * You will reference the file in your script by the filename: `Asset("tater.jpg")`. Whatever the file was called, thats how it will be referenced in the runtime.
 * Embedding works as well for the `gscript run` command. This allows you to test your assets easily.
 * There are more macros below. All are optional, but available if you need them.

#### List of Macros

| Macro Name | Default     | Example                                                 | Purpose                                                                                                      |
|---------------|----------|---------------------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| `import`   | `null` | `//import:/tmp/gscript_demo/tater.jpg` | Embeds a local file into your scripts runtime. |
| `import_url`   | `null` | `//import_url:https://example.com/tater.jpg` | Embeds a remote file into your scripts runtime. |
| `timeout`    | `30` | `//timeout:120`                                             | Overrides the default timeout for this specific VM instance.                                                                         |
| `priority`    | `100`  | `//priority:50`                           | *Overrides the runtime priority of a VM.                                                         

* *The priority macro allows you to tune the execution flow of your final binary. If you compile multiple scripts (`gscript compile --outfile /tmp/foo.bin a.gs b.gs c.gs`) your scripts by default will all execute in parallel since the default priority is 100. Any scripts that override this will be executed either before or after the default. The lower the priority, the earlier in the execution chain it will occur. Note that each unique priority blocks until all VMs for that priority have returned.*

*For example, if `a.gs` and `b.gs` have a priority of 15, and `c.gs` has a priority of 100, `c.gs` will not begin execution until BOTH `a.gs` and `b.gs` have returned or timed out.*


### Embedded Variables

These variables are pre-defined and injected into the Genesis VM at runtime:

| Variable Name | Type     | Example                                                 | Purpose                                                                                                      |
|---------------|----------|---------------------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| `USER_INFO`   | `object` | `{uid: 0, gid: 0, username: "root", home_dir: "/root"}` | Information about the User. Will be basically whatever is returned with https://golang.org/pkg/os/user/#User |
| `HOSTNAME`    | `string` | `example01`                                             | The hostname of the machine.                                                                                 |
| `IP_ADDRS`    | `array`  | `["127.0.0.1","192.168.1.5"]`                           | The IP addresses of the machine.                                                                             |
| `OS`          | `string` | `linux`                                                 | The operating system (basically `runtime.GOOS`)                                                              |
| `ARCH`        | `string` | `amd64`                                                 | The CPU architecture (basically `runtime.GOARCH`)                                                            |

### VM Functions

The VM functions are documented in the Godoc for the Engine. It includes examples of Javascript calling as well as examples of all the values functions return within the VM.

https://godoc.org/github.com/gen0cide/gscript/engine

### Interactive Debugger

In addition to the `DebugConsole()` command in your script, you can explore the GScript VM directly from an interactive REPL using the `gscript shell` command:

![Interactive REPL Example](https://i.imgur.com/EPuzHLJ.png)

## What is GENESIS btw?

GENESIS was created by @vyrus, @gen0cide, @emperorcow, and @ahhh for dynamically bundling multiple payloads into one dropper for faster deployment of implants for the CCDC Red Team.

For more information on this work we do every year, see my blog post outlining our toolbox:

 * <https://alexlevinson.wordpress.com/2017/05/09/know-your-opponent-my-ccdc-toolbox/>

Inspiration for this comes from my old AutoRune™ days and from the need for malware to basically become self aware without a bunch of duplicate overhead code.

## Credits

Shoutouts to the homies:

 * vyrus
 * ahhh
 * cmccsec
 * carnal0wnage
 * indi303
 * emperorcow
 * rossja
 * jackson5
