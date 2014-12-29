# Zencoder

Zencoder client for Go

## usage

- import the package with an alias for easy usage
```go
import (
	zencoder "github.com/streamrail/zencoder"
)
```
- init zencoder.Client with a zencoder.Options pointer reference. default options:
```go
type Options struct {
	ApiKey       string // mandatory
	ApiEndpoint  string // https://app.zencoder.com/api/v2/jobs
	ResponseType string // application/json
	Timeout      int // 30 seconds
}
```

- for example, if you have a string flag called apiKey with your Zencoder key:
```go
if client, err := zencoder.NewClient(&zencoder.Options{
		ApiKey: *apiKey,
	}); err != nil {
		log.Fatal(err.Error())
} 
```
- client.Zencode function returns a JSON object (map[string]interface{}):
```go
if res, err := client.Zencode(*input, o); err != nil {
	log.Fatal(err.Error())
} else {
	str, _ := json.Marshal(res)
	log.Printf("response: %s", str)
}
```

## examples

check out example/example.go. 

- example: zencode a video into [MPEG-DASH](https://app.zencoder.com/docs/guides/encoding-settings/dash) and saving the result on an S3 Bucket (read [this](https://app.zencoder.com/docs/guides/getting-started/working-with-s3) to understand how to setup Zencoder with your S3 credentials). 

this would be a POST to the Zencoder API with a JSON that looks like this:
```json
{
  "input": "http://s3.amazonaws.com/zencodertesting/test.mov",
  "outputs": [
    {
      "streaming_delivery_format": "dash",
      "video_bitrate": 700,
      "type": "segmented",
      "url": "s3://mybucket/dash-examples/sbr/rendition.mpd"
    }
  ]
}

```

```go
package main

import (
	"encoding/json"
	"flag"
	zencoder "github.com/streamrail/zencoder"
	"log"
)

var (
	apiKey  = flag.String("key", "", "your zencoder api key")
	input   = flag.String("input", "http://s3.amazonaws.com/zencodertesting/test.mov", "a video you want to zencode")
	outputs = flag.String("outputs", "[{\"streaming_delivery_format\": \"dash\",\"video_bitrate\": 700,\"type\": \"segmented\",\"url\": \"s3://mybucket/dash-examples/sbr/rendition.mpd\"}]", "outputs for the zencoded video")
)

func main() {
	flag.Parse()

	if *apiKey == "" {
		log.Fatalf("Required flag: -key")
	}

	if *input == "" {
		log.Fatalf("Required flag: -input")
	}

	if *outputs == "" {
		log.Fatalf("Required flag: -outputs")
	}

	var o []map[string]interface{}
	err := json.Unmarshal([]byte(*outputs), &o)
	if err != nil {
		log.Fatal(err.Error())
	}

	if client, err := zencoder.NewClient(&zencoder.Options{
		ApiKey: *apiKey,
	}); err != nil {
		log.Fatal(err.Error())
	} else {
		if res, err := client.Zencode(*input, o); err != nil {
			log.Fatal(err.Error())
		} else {
			str, _ := json.Marshal(res)
			log.Printf("response: %s", str)
		}
	}
}

```

## google app engine

google app engine does not support using the net/http package directly. the [zencoder-gae](http://github.com/streamrail/zencoder-gae) package uses appengine/urlfetch instead 