// -----------------------------------------------------
// GENESIS Script Engine
// https://github.com/gen0cide/gscript
//
// Example Script
// 
// *Note: .gs files are just javascript! Enjoy!
// -----------------------------------------------------

// -----------------------------------------------------
// COMPILER MACROS (Remove comment and modify to use)
// //import:/path/to/local/file.exe
// //import_url:https://website.com/with/local/file.exe
// In both of those examples, you can retrieve the asset
// by using Asset("file.exe") (returns byte array).

// -----------------------------------------------------
// GLOBALS
var test_var = "hello, world";
var test_num = 5;

// -----------------------------------------------------
// HOOKS
// Your final script must implement these methods.
// If any method returns false, the VM will cease
// execution and not continue to the subsequent
// functions.

// -----------------------------------------------------
// BeforeDeploy() is meant to allow you the opportunity
// to investigate the target system to determine if you
// even want to proceed to the Deploy() step. You might:
// - Check to see if you have the right architecture.
// - Check to see if apache is installed.
// - Check to see if your payload has already been dropped.
function BeforeDeploy() {
  LogInfo("This is an example of a debug command.");
  LogWarn("This will only show when you use the \"run\" subcommand.");
  LogError("It will *not* print when run from a compiled binary.");
  return true;
}

// -----------------------------------------------------
// Deploy() is where you actually deploy your payload.
// Remember to return true if it deploys successfully.
function Deploy() {
  if(test_num > 2) {
    LogInfo("test_var = " + test_var);
  }
  return true;
}

// -----------------------------------------------------
// AfterDeploy() allows you to clean up or validate deployment.
function AfterDeploy() {
  return true;
}


