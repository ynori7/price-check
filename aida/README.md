# Aldiana Price Check
Checks the price of a given Aida cruise

## Usage
```
go run cmd/main.go --config config.yaml
```

Be sure to first copy config.yaml.dist to config.yaml and fill in the missing blanks

### Set up cronjob:

First, build the binary:
```
go build -o bin/aidapricecheck aida/cmd/main.go
```

Then set up the cronjob:

```
0 14 * * * /path/to/goprojects/src/github.com/ynori7/price-check/bin/aidapricecheck --config /path/to/goprojects/src/github.com/ynori7/price-check/aida/config.yaml >> /path/to/goprojects/src/github/ynori7/price-check/log/aidapricecheck.log 2>&1
```