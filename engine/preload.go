package engine

var VMPreload = `
function StringToByteArray(s) {
  var data = [];
  for (var i = 0; i < s.length; i++ ) {
    data.push(s.charCodeAt(i));
  }
  return data;
}

function ByteArrayToString(a) {
  return String.fromCharCode.apply(String, a);
}

function DumpObjectIndented(obj, indent) {
  var result = "";
  if (indent == null) indent = "";

  for (var property in obj) {
    var value = obj[property];
    if (typeof value == 'string') {
      value = "'" + value + "'";
    }
    else if (typeof value == 'object') {
      if (value instanceof Array) {
        value = "[ " + value + " ]";
      } else {
        var od = DumpObjectIndented(value, indent + "  ");        
        value = "\n" + indent + "{\n" + od + "\n" + indent + "}";
      }
    }
    result += indent + "'" + property + "' : " + value + ",\n";
  }
  return result.replace(/,\n$/, "");
}

function BeforeDeploy() {
  return false;
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
    if ((new Date().getTime() - start) > (seconds * 1000)){
      break;
    }
  }
}
`
