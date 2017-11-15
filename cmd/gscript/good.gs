// genesis script

var helloWorld = "helloworld";
var foo = MD5(helloWorld);
VMLogInfo(foo);
console.log("haha");
DebugConsole();
console.log("wat");
var dem_bytes = StringToByteArray(helloWorld);
console.log(dem_bytes);
var newString = ByteArrayToString(dem_bytes);
console.log(newString);

