// genesis script

//import:/Users/flint/Downloads/tater.jpg

function BeforeDeploy() {
 return true; 
}

function Deploy() {
  var tater = Asset("tater.jpg");
  var ts = RandomString(12);
  var fn = "/tmp/" + ts.value + "_tater.jpg";
  var fileRet = WriteFile(fn, tater.fileData, 0777);
  DebugConsole();
  return true;
}

function AfterDeploy() {
  return true;
}
