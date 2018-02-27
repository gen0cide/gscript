function BeforeDeploy() { 
  LogTester();
  return true;
}

function Deploy() {
  LogDebug("testing debug logging");
  LogInfo("testing info logging");
  LogWarn("testing warn logging");
  LogError("testing error logging");
  LogFatal("testing fatal logging");
  return true;
}

function AfterDeploy() {
  LogInfo("Test complete.");
  return true;
}
