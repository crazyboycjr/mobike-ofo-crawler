# mobike-ofo-crawler

A crawler script, fetching real time bike positions from Mobike API and ofo API~~, saving the records in the Posgresql database~~.

## Build
Change to the project directory, and
```
go build -o $GOPATH/bin/mocrawler github.com/crazyboycjr/mobike-ofo-crawler
```
The executable binary file should locate in `$GOPATH/bin/` directory if build successfully.

## Usage

Save the correct tokens in token.txt, then the following command to read tokens from `tokens.txt` and save the result to `data.txt`
```
mocrawler --token=tokens.txt --output=data.txt
```
the default is read from `./tokens.txt` and output to `./data.txt`

~~To directly save the results to databases, use~~, haven't implemented this feature.
```
mocrawler --token=tokens.txt --host=127.0.0.1 --port=23000
```

Moreover, you can specify which module to run, use
```
mocrawler --module mobike --module ofo
```
and specify the concurrency number of each module
```
mocrawler --module mobike -Cmobike 10 --module ofo -Cofo 1
```
