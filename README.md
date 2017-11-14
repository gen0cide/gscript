# gscript

Malware Dropper Scripting Engine

**WARNING: THIS IS SHIT CODE. USE AT YOUR OWN RISK.**

## Description

GENESIS Scripting (gscript for short) is a technology I've developed to allow dynamic runtime execution of malware installation based on parameters determined at runtime.

Inspiration for this comes from my old AutoRune™ days and from the need for malware to basically become self aware without a bunch of duplicate overhead code.

GScript uses a JS V8 Virtual Machine to interpret your genesis script and allow it to hook into the malware initialization.

## What is GENESIS btw?

GENESIS was created by @vyrus, @gen0cide, and @ahhh.db for dynamically bundling multiple payloads into one dropper for faster deployment of implants for the CCDC Red Team.

GENESIS Script is a virtual machine to allow intelligent deployment of those payloads.

## Notes
```

// Globals (Can be user defined)
var source_url    = null;
var source_folder = null;
var packed_file   = null;
var source_bytes  = null;
var file_dest     = null;
var exec_args     = null;
var timeout       = 180000; // 3min

// Read-Only (Defined at runtime)
var user_info = null;
var hostname  = null;
var ip_addrs  = null;
var os        = null;
var platform  = null;
var error     = null;

// Genesis Hooks (Can be user defined)
function BeforeDeploy() {};
function Deploy()       {};
function AfterDeploy()  {};
function OnError()      {};

// Core Functions (Should *not* be overridden!)
// use these method signatures as references in
// other functions or hooks above.
function Halt()                               { return true;  };
function DeleteFile(path)                     { return false; };
function WriteFile(path, bytes, perms)        { return false; };
function ExecuteFile(path, args)              { return false; };
function AppendFile(path, bytes)              { return false; };
function ReplaceInFile(path, target, replace) { return false; };
function Signal(pid, signal)                  { return false; };
function Implode()                            { return true;  };

// Investigatory Functions (Should *not* be overridden!)
// These functions allow you to get information about
// a given system in order to make decisions based off
// the context the runtime is executing in.
function LocalUserExists(username)      { return true;  };
function ProcExistsWithName(name)       { return true;  };
function CanReadFile(path)              { return true;  };
function CanWriteFile(path)             { return true;  };
function CanExecFile(path)              { return true;  };
function FileExists(path)               { return true;  };
function DirExists(path)                { return true;  };
function FileContains(path, match)      { return true;  };
function IsVM()                         { return true;  };
function IsAWS()                        { return true;  };
function HasPublicIP()                  { return true;  };
function CanMakeTCPConn(dst, port)      { return true;  };
function ExpectedDNS(query, type, resp) { return true;  };
function CanMakeHTTPConn(url)           { return true;  };
function DetectSSLMITM(url, cert_fp)    { return true;  };
function CmdSuccessful(cmd)             { return true;  };
function CanPing(dst)                   { return true;  };
function TCPPortInUse(port)             { return true;  };
function UDPPortInUse(port)             { return true;  };
function ExistsInPath(progname)         { return true;  };
function CanSudo()                      { return true;  };
function Matches(string, match)         { return true;  };
function CanSSHLogin(ip, port, u, p)    { return true;  };

// Utility functions
function RetrieveFileFromURL(url)       { return bytes;  };
function DNSQuery(question, type)       { return resp;   };
function HTTPRequest(m, u, b, h)        { return resp;   };
function Cmd(cmd)                       { return resp;   };
function MD5(bytes)                     { return hash;   };
function SHA1(bytes)                    { return hash;   };
function B64Decode(string)              { return bytes;  };
function B64Encode(bytes)               { return string; };
function Timestamp()                    { return epoch;  };
function CPUStats()                     { return stats;  };
function MemStats()                     { return stats;  };
function SSHCmd(ip, port, u, p, cmd)    { return output; };
function Sleep(seconds)                 { return true;   };
function GetTweet(tid)                  { return true;   };
function GetDirsInPath()                { return dirs;   };
function EnvVars()                      { return vars;   };
function GetEnv(varname)                { return value;  };
function FileCreateTime(path)           { return epoch;  };
function FileModifyTime(path)           { return epoch;  };
function LoggedInUsers()                { return users;  };
function UsersRunningProcs()       { return users;  };
function ServeFileOverHTTP(file, port)  { return true;   };

```

## TODO

 * Implement All Functions
 * Implement Global Variables
 * Implement Runtime Variable Loading
 * Implement Hook Callers
 * Implement CLI Framework (cmd/gscript)
 * Implement Compiler / Crypter

## Credits

Shoutouts to the homies:

 * vyrus
 * ahhh.db
 * cmccsec
 * carnal0wnage
 * indi303