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
	partitionKey := "1"
	sess := common.GetSession()
	common.ShowCallerId(common.GetSTS(sess))
	common.PutRecords(streamName, partitionKey, common.GetKinesis(sess))
	log.Println("Done")
}
