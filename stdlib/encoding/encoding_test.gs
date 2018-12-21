// Example gscript template
// Title: Multi Platform Encoding Examples
// Author: ahhh
// Purpose: Testing a bunch of data encoding on different platforms
// Gscript version: 1.0.0

//priority:150
//timeout:150

function Deploy() {
    console.log("Starting Base64e Command");
    var b64 = G.encoding.EncodeBase64("hello world");
    console.log("b64e: "+b64);

    console.log("Starting Base64d Command");
    var decoded = G.encoding.DecodeBase64("aGVsbG8gd29ybGQ=");
    console.log("b64d: "+decoded[0]);
    console.log("b64d errors: "+ decoded[1]);

    console.log("Starting EncodeStringAsBytes");
    var bytes = G.encoding.EncodeStringAsBytes("test");
    console.log("bytes: "+ bytes);

    console.log("EncodingBytesAsString");
    var bstring = G.encoding.EncodeBytesAsString(bytes);
    console.log("bytes to string: "+ bstring);
    
    return true;
}
