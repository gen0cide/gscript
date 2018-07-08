//import:/Users/flint/Downloads/tater.jpg

function BeforeDeploy() {
  LogInfo("BeforeDeploy()");
  return true;
}

function Deploy() {
  var firstHash = MD5("helloworld");
  var secondHash = MD5(Asset("tater.jpg").fileData);
  DebugConsole();
  return true;
}

function AfterDeploy() {
  LogInfo("AfterDeploy()");
  return true;
}
