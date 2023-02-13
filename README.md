
# Web Crawler

## Usage

Build the project
```
$ go build cmd/crawlerapp/main.go
```

Run the project
```
$ ./main https://google.com
```

**Notes**: 
- You must only supply one URL for the program to run
- You can supply a URL in any of the following forms
    - `https://example.com`
    - `http://example.com`
    - `www.example.com`
    - `example.com`

## Tests

Run the tests
```
$ go test ./... -v
```

Test coverage
```
ok  	github.com/beccafields/crawler/crawler	0.003s	coverage: 91.9% of statements
```