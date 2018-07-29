//go_import:github.com/gen0cide/gscript/x/svc as svc
//import:/tmp/svctest-win.exe

var service_bin_path = "C:\\WINDOWS\\system32\\megadeath.exe"


var serviceSettings = {
  name: "gscript_eicar_service",
  display_name: "gscript eicar service",
  description: "genesis engine eicar example service",
  arguments: [],
  executable_path: service_bin_path,
  working_directory: "C:\\WINDOWS\\system32",
  options: {}
}

function Deploy() {
  console.log("Writing binary to disk...")

  filedata = GetAssetAsString("svctest-win.exe")

  var quitAfterDebug = false;
  DebugConsole()

  if (quitAfterDebug == true) {
    return
  }

  errchk = G.file.WriteFileFromString(service_bin_path, filedata[0])
  if (errchk !== undefined) {
    console.error("Error writing file: " + errchk.Error())
    DebugConsole()
    return
  }

  console.log("Creating new service object...")

  var svcObj = svc.NewFromJSON(serviceSettings)
  if (svcObj[1] !== undefined) {
    console.error("Error creating service: " + svcObj[1].Error())
    DebugConsole()
    return
  }

  console.log("Checking service config sanity...")

  confchk = svcObj[0].CheckConfig(true)
  if (confchk[1] !== undefined || confchk[0] === false) {
    console.error("Error checking config: " + confchk[1].Error())
    DebugConsole()
    return
  }

  console.log("Installing service...")

  installchk = svcObj[0].Install(true)
  if (installchk !== undefined) {
    console.error("Error installing service: " + installchk.Error())
    DebugConsole()
    return
  }

  console.log("Starting service...")

  startchk = svcObj[0].Start()
  if (startchk !== undefined) {
    console.error("Error starting service: " + startchk.Error())
    DebugConsole()
    return
  }

  DebugConsole()
}