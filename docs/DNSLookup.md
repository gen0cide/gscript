### `DNSQuery(question, type)`

Perform a DNS lookup.

##### Argument List

 * `question` (String) - DNS query question (eg: "twitter.com")
 * `type` (String) - DNS question type (A, CNAME, MX, etc.)

##### Return Type

`Object` - Reference `VMDNSQueryResponse` in `response_objects.go` for object details.

#Example 

This will perform a DNS lookup. 