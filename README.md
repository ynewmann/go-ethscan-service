# Fun EtherScan Api Implementation

## Run

To run simple runnable example:
```go
go run main.go 
```

### Flags
- `apiUrl` - set Api URL 
- `apiKey` - se Api key
- `useCache` - use cache (default `true`)
- `port` - set server address (default `:3333`) 
- `cacheSize` - set size of cache (default `1000`) 
- `cache-bckp` - set path to a backup of cache
- `load-cfg` - load config from a file. If set, all other flags will be ignored. 
