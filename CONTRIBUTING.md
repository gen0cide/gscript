# Contributing

## Contributing To Genesis Engine or Compiler

## Contributing To The Standard Library

- Must have cause to implement the function in Golang (see GOLANG.md)
- Must have an associated issue
- Must have an associated PR against the LATEST version of the dev branch
- Must compile without errors for all supported target architectures
  - If package is OS specific, must be implemented in OS package
- Must pass all linter checks
- Must have an associated unit test(s) written in Golang
- Must have an associated genesis script that implements the function standalone in a gscript for in VM testing
- Must follow guidelines for implementing standard library functions (STANDARD_LIBRARY_GUIDELINES.md)
- Must include a commit updating the standard library documentation (using the function template) and include:
  - Name
  - Description
  - Author
  - Method Signature
  - Arguments and their types
  - Returns and their types
  - Example usage
  - How to handle failure
- Must have a revision to the CHANGELOG.md

## Contributing Example Scripts

## Contributing To Documentation
