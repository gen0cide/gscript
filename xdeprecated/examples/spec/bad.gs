// genesis script

var helloWorld = "helloworld";
var foo = MD5(helloWorld);
VMLogInfo(foo);
console.log("haha");
console.log("wat");
var dem_bytes = StringToByteArray(helloWorld);
console.log(dem_bytes);
var newString = ByteArrayToString(dem_bytes);
console.log(newString);

var fileTest = WriteFile("/tmp/foobar1234", -38882.5488202);
console.log(fileTest);

asdfpaoihsdf


