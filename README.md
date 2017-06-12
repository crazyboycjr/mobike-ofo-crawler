# mobike-ofo-crawler

A crawler script, fetching real time bike positions from Mobike API and ofo API<s>, saving the records in the Posgresql database</s>.

## Build
```
go get github.com/crazyboycjr/mobike-ofo-crawler
go build -o $GOPATH/bin/mocrawler github.com/crazyboycjr/mobike-ofo-crawler
```
The executable binary file should locate in `$GOPATH/bin/` directory if build successfully.

## Usage

In order to correctly fetch mobike data, you must obtain the `accesstoken` and `mobileno` from a real user. The `accesstoken` and `mobileno` can be captured with the help of some packet sniffers like wireshark when you are using the wechat mini program.
Afterwards, save the tokens as the form of `token.txt.example` in `token.txt`, then the following command will read tokens from `token.txt` and save the result to `data.txt`
```
mocrawler --token=token.txt --output=data.txt
```
the default is read from `./tokens.txt` and output to `./data.txt`

~~To directly save the results to databases, use~~, haven't implemented this feature.
```
mocrawler --token=token.txt --host=127.0.0.1 --port=23000
```

Moreover, you can specify which module to run, use
```
mocrawler --module mobike --module ofo
```
and specify the concurrency number of each module
```
mocrawler --module mobike -Cmobike 10 --module ofo -Cofo 1
```

## Details

Please look at [details.txt](https://github.com/crazyboycjr/mobike-ofo-crawler/blob/master/details.txt)
