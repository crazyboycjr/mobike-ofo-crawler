# mobike-ofo-crawler

A crawler script, fetching real time bike positions from Mobike API and ofo API, saving the record in the Posgresql database.

## Build
Change to the project directory, and
```
GOPATH=`pwd` go install main
```
The executable binary file should locate in `bin/` directory if build successfully.

## Usage

Save the correct tokens in token.txt
```
./main --token=tokens.txt --output=data.txt
```
or
```
./main --token=tokens.txt --host=127.0.0.1 --port=23000
```
