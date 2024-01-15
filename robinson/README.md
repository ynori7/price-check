# Robinson Price Check
Checks the price of a given Robinson vacation

## Usage
```
go run cmd/main.go --config config.yaml
```

Be sure to first copy config.yaml.dist to config.yaml and fill in the missing blanks

### Set up cronjob:

First, build the binary:
```
go build -o bin/robinsonpricecheck robinson/cmd/main.go
```

Then set up the cronjob:

```
0 14 * * 1 /path/to/goprojects/src/github.com/ynori7/price-check/bin/robinsonpricecheck --config /path/to/goprojects/src/github.com/ynori7/price-check/robinson/config.yaml >> /path/to/goprojects/src/github/ynori7/price-check/log/robinsonpricecheck.log 2>&1
```
