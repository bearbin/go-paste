go-paste
========

A tool to upload text to pastebin. Fairly simple, and easy to use.
Support for other paste services is planned for the future.

## Installation

    go get github.com/bearbin/go-paste

## Examples

Put a file:

    go-paste put example-file

Put from stdin:

    go-paste put -

Put your data on fpaste, not pastebin:

    go-paste -s fpaste put example-file

Get a paste:

	go-paste get http://pastebin.com/ZTBUm4B2

## Supported Paste Services

 - pastebin.com
 - fpaste.org

## Development

go-paste is Public Domain software and Pull Requests are welcome.
