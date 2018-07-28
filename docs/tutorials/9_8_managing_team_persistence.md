# Managing Team Persistence

One of the features of GSCRIPT is being able to have co-operative, single payloads. This means multiple peaople, red team members, can sumbit arbitrary payloads and their associated gscripts which map to their desired persistence technique. In this way, a single person, the genesis builder we can call them, can manage all of the teams payloads and compile them into a final bianry. This person can also audit the payloads and gscripts to make sure there are no conflicting binaries or persistence locations, such that two agents will stomp each other in some race condition. 

## Testing

Its important to test the implants on a wide varrity of systems before wrapping them in gscript. Further, it's important to test individual gscripts before testing multiple scripts in a single binary. When testing use no obfuscation and enbale logging, however make sure to do a final test w/ the desired obfuscation level and no logging before running the final payload like that. It can help to have automated test infrastructure such as a vagrant farm to check for errors, however one can also make do with local vms and snapshots. Don't test or enable debugging on victim machiens.

## Priority

One of the biggest advantages of gscript is leveraging the priority to allow certain payloads to run before others. In this way the genesis builder can select payloads that disable security controls to run before other implants, protecting other peoples dropped binaries.

## Persistence SpreadSheet

You can maintain a spread sheet like the one seen below to track each contributions implants: 
- gscripts,
- general ATT&CK techniques (persistence techniques)
- file locations 
- target platforms / arch 
- test status 
- etc...
