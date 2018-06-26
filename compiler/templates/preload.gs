function StringToByteArray(s) {
  var data = [];
  for (var i = 0; i < s.length; i++) {
    data.push(s.charCodeAt(i));
  }
  return data;
}

function ByteArrayToString(a) {
  return String.fromCharCode.apply(String, a);
}

function Dump(obj) {
  return "\n" + JSON.stringify(obj, null, 2);
}

function BeforeDeploy() {
  return true;
}

function Deploy() {
  return false;
}

function AfterDeploy() {
  return true;
}

function OnError() {
  return false;
}

function Sleep(seconds) {
  var start = new Date().getTime();
  for (var i = 0; i < 1e7; i++) {
    if ((new Date().getTime() - start) > (seconds * 1000)) {
      break;
    }
  }
}

function DebugConsole() {
  return true;
}