// Example gscript template
// Title: Multi Platform Randomness Example
// Author: ahhh
// Purpose: Testing randomness on different platforms
// Gscript version: 1.0.0

//priority:150
//timeout:150

function Deploy() {

    console.log("Starting RandomInt");
    out1 = G.rand.RandomInt(12, 20);
    console.log("out: "+out1);

    console.log("Starting GetAlphaNumericString");
    var out2 = G.rand.GetAlphaNumericString(12);
    console.log("out: " +out2);
    console.log("out upper: "+ out2.toUpperCase());
    console.log("out lower: "+ out2.toLowerCase());
    
    console.log("Starting GetAlphaString");
    var out3 = G.rand.GetAlphaString(10);
    console.log("out: "+ out3);
    console.log("out upper: "+ out3.toUpperCase());
    console.log("out lower: "+ out3.toLowerCase());

    console.log("Starting GetAlphaNumericSpecialString");
    var out4 = G.rand.GetAlphaNumericSpecialString(7);
    console.log("out: "+ out4);
    console.log("out upper: "+ out4.toUpperCase());
    console.log("out lower: "+ out4.toLowerCase());

    console.log("Starting GetBools");
    var out5 = G.rand.GetBool();
    console.log("out: "+ out5);

    return true;
}
