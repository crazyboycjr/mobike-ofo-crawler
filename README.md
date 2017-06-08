# mobike-ofo-crawler

A crawler script, fetching real time bike positions from Mobike API and ofo API, saving the record in the Posgresql database.

## Build
Change to the project directory, and
```
go build -o $GOPATH/bin/mocrawler github.com/crazyboycjr/mobike-ofo-crawler
```
The executable binary file should locate in `$GOPATH/bin/` directory if build successfully.

## Usage

Save the correct tokens in token.txt
```
mocrawler --token=tokens.txt --output=data.txt
```
or
```
mocrawler --token=tokens.txt --host=127.0.0.1 --port=23000
```
