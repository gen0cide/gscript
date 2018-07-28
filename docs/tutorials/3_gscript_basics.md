# GSCRIPT Basics
Writing gscripts is very simple, the syntax is a very loose JavaScript, giving you all native JavaScript functions. 

## Things to know
- JavaScript Engine 5
That gives you these native functions: 

- You will want to be familiar w/ the native Lib
These are helper GoLang functions we wrote just for GSCRIPT
    - You will want to read the native lib docs 
    These will help you understand how to handle native lib return objects

- You will want to read the GoDocs of any GoLang Native Library you use
These are very helpful for understand the objects that things will return

- You will want to Debug and Log when testing your scripts
Read more about that in 4 and 5

- You will want to use Obfuscation when deploying your scripts to your targets
Read more about that in 9_5

## Writing Your First Script
To start your first script, try a simple function like the one that follows:
```
function Deploy() {
    console.log("Hello World");
    var out = G.rand.GetAlphaNumericString(4);
    console.log("out uppecase: "+ out.toUpperCase());
}
```

## Compiling Your First Script
You will want to compile the above script with logging, such that you can see the output of your randomly generated string
./gscript compile --enable-logging --obfuscation-level 3 ./hello_world.gs

## More simple examples
https://github.com/ahhh/gscripts/tree/master/attack/multi
