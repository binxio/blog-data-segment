package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/sts"
	"log"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func GetSession() *session.Session {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	return sess
}

func GetSTS(sess *session.Session) *sts.STS {
	return sts.New(sess)
}

func GetKinesis(sess *session.Session) *kinesis.Kinesis {
	return kinesis.New(sess, aws.NewConfig().WithRegion("eu-west-1"))
}

func ShowCallerId(svc *sts.STS) {
	request := sts.GetCallerIdentityInput{}
	id, err := svc.GetCallerIdentity(&request)
	if err != nil {
		panic(err)
	}
	fmt.Printf("AccountId: %s, Arn: %s, UserId: %s", *id.Account, *id.Arn, *id.UserId)
}

func CreatePerson(name string, age int) Person {
	return Person{Name: name, Age: age}
}

func SerializePerson(person Person) []byte {
	data, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}
	return []byte(string(data) + "\n")
}

func CreateRecord(streamName string, partitionKey string, data []byte) *kinesis.PutRecordInput {
	return &kinesis.PutRecordInput{
		PartitionKey: aws.String(partitionKey),
		StreamName:   aws.String(streamName),
		Data:         data,
	}
}

func PutRecords(streamName string, partitionKey string, svc *kinesis.Kinesis) {
	for i := 0; i < 100000; i++ {
		person := CreatePerson(fmt.Sprintf("dennis %d", i), i)
		data := SerializePerson(person)
		record := CreateRecord(streamName, partitionKey, data)
		res, err := svc.PutRecord(record)
		log.Println("Publishing:", string(data), hex.EncodeToString(data), res)
		if err != nil {
			panic(err)
		}
	}
}

func ProcessRecords(records *[]*kinesis.Record) {
	for _, record := range *records {
		log.Println(string(record.Data))
	}
}

func GetNextRecords(it *string, svc *kinesis.Kinesis) {
	records, err := svc.GetRecords(&kinesis.GetRecordsInput{
		ShardIterator: it,
	})
	if err != nil {
		panic(err)
	}
	ProcessRecords(&records.Records)
	GetNextRecords(records.NextShardIterator, svc)
}


func GetRecords(streamName string, shardId string, svc *kinesis.Kinesis) {
	it, err := svc.GetShardIterator(&kinesis.GetShardIteratorInput{
		StreamName:        aws.String(streamName),
		ShardId:           aws.String(shardId),
		ShardIteratorType: aws.String("TRIM_HORIZON"),
	})
	if err != nil {
		panic(err)
	}
	GetNextRecords(it.ShardIterator, svc)
}
