# `DNSQuestion(target, request)`

Issue a DNS query.

## Argument List

 * `target` (string) - Target of DNS query.
 * `request` (string) - DNS query type. Ex: A, NS, TXT, ...

## Return Type

 * `obj.answer` (string) - Result of DNS query.
 * `obj.runtimeError (error) - Error.

## Example

```js
var obj = DNSQuestion("10.10.10.1", "A")
// obj.answer (string)
// obj.runtimeError (error)
```

