package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
)

const AWS_REGION = "ap-northeast-1"
const DYNAMO_ENDPOINT = "http://localhost:18000"

func initDynamoDB() *dynamodb.DynamoDB {
	// 配置 AWS 访问凭证和其他选项
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(AWS_REGION),
		Endpoint: aws.String(DYNAMO_ENDPOINT),
		Credentials: credentials.NewStaticCredentials(
			"dummy", "dummy", "dummy"),
	}))
	return dynamodb.New(sess)
}

type (
	DynamoDbV1Test struct {
		DynamoDbClient *dynamodb.DynamoDB
		TableName      string
	}
)

func NewDynamoDbV1Test() DynamoDbV1Test {
	return DynamoDbV1Test{
		DynamoDbClient: initDynamoDB(),
		TableName:      "Movie",
	}
}

func (d DynamoDbV1Test) Do() {
	d.createTable()
}

type User struct {
	UserID string `dynamo:"UserID,hash"`
	Name   string `dynamo:"Name,range"`
	Age    int    `dynamo:"Age"`
	Text   string `dynamo:"Text"`
}

func (d DynamoDbV1Test) createTable() {
	createTableInput := &dynamodb.CreateTableInput{
		TableName: aws.String(d.TableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	_, err := d.DynamoDbClient.CreateTable(createTableInput)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	fmt.Println("Successfully created table 'ExampleTable'")
	err = d.DynamoDbClient.WaitUntilTableExists(&dynamodb.DescribeTableInput{
		TableName: aws.String("ExampleTable"),
	})
	if err != nil {
		log.Fatalf("Wait for table exists failed. Here's why: %v", err)
	}

	fmt.Println("Table 'ExampleTable' is now active")
}
