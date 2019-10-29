# Managing Team Persistence

One of the features of GSCRIPT is being able to have co-operative, single payloads. This means multiple red team members can sumbit arbitrary payloads and their associated gscripts which map to their desired persistence technique. In this way a single person can manage all of the team's payloads and compile them into a final bianry. This person can also audit the payloads and gscripts to make sure there are no conflicting binaries or persistence locations, such that two agents will stomp each other in some race condition. 

## Testing

It is important to test the implants on a wide variety of systems before wrapping them in gscript. Further, it is important to test individual gscripts before testing multiple scripts in a single binary. When testing use no obfuscation and enable logging. During final testing use the desired obfuscation level and no logging before running the final payload.

## Priority

One of the biggest advantages of gscript is leveraging the priority to allow certain payloads to run before others. In this way the genesis builder can select payloads that disable security controls to run before other implants, protecting other team members' dropped binaries.
