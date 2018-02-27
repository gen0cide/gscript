//timeout:3
//import:/Users/flint/Downloads/tater.jpg

function BeforeDeploy() {
  var tater = Asset("tater.jpg");
  var ts = Timestamp();
  var fn = "/tmp/" + ts + "_tater.jpg";
  WriteFile(fn, tater);
  LogInfo("Wrote file to " + fn);
  return true;
}

function Deploy() {
  Sleep(10);
  return true;
}

function AfterDeploy() {
  LogInfo("Completed sleeping!");
  return true;
}

function OnError() {
  LogError("This gscript encountered an error.");
  return true;
}