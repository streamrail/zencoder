package main

import (
	"encoding/json"
	"flag"
	zencoder "github.com/streamrail/zencoder"
	"log"
)

var (
	apiKey = flag.String("key", "", "your zencoder api key")
	input  = flag.String("input", "https://s3.amazonaws.com/mybucket/myfolder/default.mp4", "a video you want to zencode")
	//// encode a simpler video:
	// outputs = flag.String("outputs", "[{ \"url\": \"s3://mybucket/myfolder/output.mp4\", \"credentials\": \"s3_production\" }]", "outputs for the zencoded video")
	outputs = flag.String("outputs", "[{\"streaming_delivery_format\": \"dash\",\"video_bitrate\": 700,\"type\": \"segmented\",\"url\": \"s3://mybucket/myfolder/dash_test.mpd\", \"credentials\": \"s3_production\" }]", "outputs for the zencoded video")
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
