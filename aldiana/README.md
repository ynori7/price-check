# Aldiana Price Check
Checks the price of a given Aldiana vacation

## Usage
```
go run cmd/main.go --config config.yaml
```

Be sure to first copy config.yaml.dist to config.yaml and fill in the missing blanks

### Set up cronjob:

First, build the binary:
```
go build -o bin/aldianapricecheck aldiana/cmd/main.go
```

Then set up the cronjob:

```
0 14 * * 1 /path/to/goprojects/src/github.com/ynori7/price-check/bin/aldianapricecheck --config /path/to/goprojects/src/github.com/ynori7/price-check/aldiana/config.yaml >> /path/to/goprojects/src/github/ynori7/price-check/log/aldianapricecheck.log 2>&1
```