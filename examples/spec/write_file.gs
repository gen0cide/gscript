//import:/Users/flint/Downloads/service.exe

function BeforeDeploy() { 
  return true; 
}

function MakeID() {
  var text = "";
  var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

  for (var i = 0; i < 5; i++)
    text += possible.charAt(Math.floor(Math.random() * possible.length));

  return text;
}

function Deploy() {
  return WriteFile("C:/Users/galaxy/Desktop/" + MakeID() + ".exe", Asset("service.exe"), "0644");
}

function AfterDeploy() { 
  var file_hash = MD5(Asset("service.exe"));
  LogInfo("Asset MD5: " + file_hash);
  return true;
}

// //import:/Users/flint/Downloads/service.exe

// function BeforeDeploy() {
//   LogInfo("BeforeDeploy()");
//   return true;
// }

// function Deploy() {
//   return WriteFile("Z:/Public/service.exe", Asset("service.exe"), "0755");
// }

// function AfterDeploy() {
//   LogInfo("AfterDeploy()");
//   return true;
// }
