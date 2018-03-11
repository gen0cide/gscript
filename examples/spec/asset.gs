//import:/Users/flint/Downloads/tater.jpg

var data = null;

function BeforeDeploy() {
 data = Asset("tater.jpg").fileData;
 return true; 
}

function Deploy() {
  var ts = Timestamp();
  var filename = "/tmp/" + ts + "_tater.jpg";
  DebugConsole();
  WriteFile(filename, data, 0644);
  LogInfo("Wrote file to " + filename);
  return true;
}

function AfterDeploy() {
  return true;
}
