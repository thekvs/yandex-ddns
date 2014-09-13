# Yandex DNS Dynamic DNS Tool
### Use Yandex DNS as a dynamic DNS service

This project uses Yandex DNS API to provide a dynamic DNS service.

## Installing
This project is written in the [Go](http://golang.org/) programming language and to build it you need to install Go compiler and set some enviroment variables. [Here is an instructions on how to do it](http://golang.org/doc/install). After you've done it, run the following command in your shell:
```
$ go get github.com/thekvs/yandex-ddns
```
and this will build the binary in ```$GOPATH/bin```.

## Configuration file options
```yandex-ddns``` uses JSON format for configuration file. Below is a list of supported configuration options.

* ```"token"``` -- Yandex DNS API token. Thi is a mandatory option.
* ```"domain"``` -- Main domain registered at Yandex DNS. This is a mandatory option.
* ```"subdomain"``` -- Subdomain of the registered domain, leave empty if you want to use main domain.
* ```"ttl"``` -- TTL value for DNS record which has to be in the range [900, 1209600], if omitted current default value is used.
* ```"logfile"``` -- File where to write logs, if omitted messages will be printen on stderr or stdout.

## Usage
To update your IP every 15 minutes install in you crontab file something like this:
```
*/15 * * * *    /path/to/yandex-ddns -config /path/to/yandex-ddns.json
```

## Licensing
All source code included in this distribution is covered by the MIT License found in the LICENSE file.
