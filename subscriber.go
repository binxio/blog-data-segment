package main

import (
	"github.com/binxio/datasegment/common"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	streamName := os.Getenv("KINESIS_STREAM_NAME")
	if streamName == "" {
		panic("KINESIS_STREAM_NAME not set")
	}
	shardId := "0" // shard ids start from 0
	sess := common.GetSession()
	common.GetRecords(streamName, shardId, common.GetKinesis(sess))
	log.Println("Done")
}
