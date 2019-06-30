# megantory ![BuildPic]

Search **all** your configured AWS profiles - FAST

* Groks your AWS credentials file for all the profiles defined therein
* Performs a free-text search against all profiles / regions / supported services
* Does all the searches mentioned above and processing with a high degree of concurrency
* Returns a list of found resources with breadcrumbs where to find them

## Prerequisites

Skip this if you are just downloading the binary (TODO)

* [GoLang Installation]
* [GoLang Environment Variables]


## Install

##### The GoLang way
```bash
go get -u github.com/Kaurin/megantory
```
or...

##### Binary install
TODO


## Usage
```bash
# Getting help
megantory help

# Searching
megantory search "my search string"
```

## Configuration

TODO


[GoLang installation]: https://golang.org/doc/install
[GoLang Environment Variables]: https://github.com/golang/go/wiki/SettingGOPATH
[BuildPic]: https://travis-ci.org/Kaurin/megantory.svg?branch=master "Master build"
