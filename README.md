# ofxgo
A library for querying OFX servers and parsing the responses and an example command-line client.

## Goals
The main purpose of this project is to provide a library to make it easier to query your financial information with OFX from the comfort of Golang, without having to marshal/unmarshal to SGML or XML. The library does *not* intend to abstract away all of the details of the OFX specification, which would be very difficult to do well. Instead, it exposes the OFX SGML/XML hierarchy as structs which mostly resemble it.

Because the OFX specification is rather... 'comprehensive,' it can be difficult for those unfamiliar with it to figure out where to start. To that end, I have created a sample command-line client which uses the library to do simple tasks (currently it only lists your accounts and queries for balances and transactions). My hope is that by studying its code, new users will be able to figure out how to use the library much faster than staring at the OFX specification (or my comprehensive documentation). The command-line client also serves as an easy way for me to test/debug the library with actual financial institutions, which frequently have 'quirks' in their implementations. The command-line client can be found in the [cmd/ofx directory](https://github.com/aclindsa/ofxgo/tree/master/cmd/ofx) of this repository.

## Documentation
(Limited) API documentation can be found at https://godoc.org/github.com/aclindsa/ofxgo

### Installing the command-line client
To install the command-line client and test it out, you may do the following:

$ go get -v github.com/aclindsa/ofxgo/cmd/ofx && go install -v github.com/aclindsa/ofxgo/cmd/ofx