# Amazon Price Check
Amazon-Price-Check allows you to provide a list of Amazon product IDs (ASINs) to periodically check to see if they fall below a configured minimum price or become available on Kindle Unlimited. The primary intent now is for checking books. 

## Usage
```
go run main.go --config config.yaml
```

Be sure to first copy config.yaml.dist to config.yaml and fill in the missing blanks

### Set up cronjob:

First, build the binary:
```
go build -o bin/amazonpricecheck amazon/cmd/main.go
```

Then set up the cronjob:

```
0 14 * * 1 /path/to/goprojects/src/github.com/ynori7/price-check/bin/amazonpricecheck --config /path/to/goprojects/src/github.com/ynori7/price-check/amazon/config.yaml 
```