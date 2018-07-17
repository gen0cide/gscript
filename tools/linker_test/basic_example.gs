//go_import:time as gotimelib
//go_import:os as oslib

function Deploy() {
  filetowrite = "/tmp/nikki.txt"

  timeobj = gotimelib.Now()
  timestring = timeobj.String()

  G.file.WriteFileFromString(filetowrite, timestring)

  console.log(timestring)

  DebugConsole()
}