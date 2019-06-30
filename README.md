# megantory ![BuildPic]

Search **all** your configured AWS profiles - FAST

* Groks your AWS credentials file for all the profiles defined therein
* Performs a free-text search against all profiles / regions / supported services / supported resources
* Does all the searches mentioned above and processing with a high degree of concurrency
* Returns a list of found resources with breadcrumbs where to find them

## Supported Services and Resources

* EC2
  * Address (EIP)
  * Instance

This list should be expanded as the project grows

## Installing
##### Linux

```bash
sudo curl https://github.com/Kaurin/megantory/releases/latest/download/megantory.linux -Lo /usr/local/bin/megantory && sudo chmod +x /usr/local/bin/megantory
```

##### MacOS

```bash
sudo curl https://github.com/Kaurin/megantory/releases/latest/download/megantory.macos -Lo /usr/local/bin/megantory && sudo chmod +x /usr/local/bin/megantory
```

##### Windows

Download the binary as listed below

## Uninstalling
##### Linux and MacOS
```bash
sudo rm /usr/local/bin/megantory
```

##### Windows
Remove the binary from where you downloaded it from

## Download

If you just want to download the binaries, you can use these links, or click the [Releases] page

[Linux binary]
[Windows binary]
[Macos binary]

## Usage
##### MacOS and Linux
```bash
# Getting help
megantory help

# Searching
megantory search "my search string"

# Print version information
megantory version
```

##### Windows
Same as the above, just use `megantory.exe` from your download location

## Configuration

TODO


## Develop

##### Prerequisites

* [GoLang Installation]
* [GoLang Environment Variables]

##### Setting up

This project uses `go mod`, so it is recommended to clone it outside of `$GOPATH`. 

For example:
```bash
mkdir -p ~/mydir
cd ~/mydir
git clone git@github.com:Kaurin/megantory.git
cd ~/mydir/megantory
go run main.go # Downloads dependancies automatically through the power of go mod!
```



##### Contributing
If you are willing to help in any way, please contact me via Github or the Issues page.

Pull requests are welcome!


[GoLang installation]: https://golang.org/doc/install
[GoLang Environment Variables]: https://github.com/golang/go/wiki/SettingGOPATH
[BuildPic]: https://travis-ci.org/Kaurin/megantory.svg?branch=master "Master build"
[Linux binary]: https://github.com/Kaurin/megantory/releases/latest/download/megantory.linux
[MacOS binary]: https://github.com/Kaurin/megantory/releases/latest/download/megantory.macos
[Windows binary]: https://github.com/Kaurin/megantory/releases/latest/download/megantory.exe
[Releases]: https://github.com/Kaurin/megantory/releases
