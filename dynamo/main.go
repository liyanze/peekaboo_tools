package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type (
	DynamoDbTest struct {
		DynamoDbClient *dynamodb.Client
		TableName      string
	}
)

func NewDynamoDbTest() DynamoDbTest {
	return DynamoDbTest{
		DynamoDbClient: initRecentConversationDynamoDB(),
		TableName:      "Movie",
	}
}

func (d DynamoDbTest) Do() {
	d.CreateMovieTable()
}

func initRecentConversationDynamoDB() *dynamodb.Client {
	// 配置 AWS 访问凭证和其他选项
	nopClient := smithyhttp.ClientDoFunc(func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Header:     http.Header{},
			Body:       ioutil.NopCloser(strings.NewReader("")),
		}, nil
	})

	client := dynamodb.NewFromConfig(aws.Config{
		Region:       "ap-northeast-1",
		BaseEndpoint: aws.String("http://localhost:18000"),
		DefaultsMode: aws.DefaultsModeStandard,
		Retryer: func() func() aws.Retryer {
			return nil
		}(),
		HTTPClient:       nopClient,
		RetryMaxAttempts: 10,
		RetryMode:        aws.RetryModeStandard,
		Credentials:      CredentialsProviderImpl{},
	})
	return client
}

type CredentialsProviderImpl struct {
}

func (CredentialsProviderImpl) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     "dummy",
		SecretAccessKey: "dummy",
		SessionToken:    "dummy",
	}, nil
}

func (d DynamoDbTest) CreateMovieTable() (*types.TableDescription, error) {
	var tableDesc *types.TableDescription
	table, err := d.DynamoDbClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("year"),
			AttributeType: types.ScalarAttributeTypeN,
		}, {
			AttributeName: aws.String("title"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("year"),
			KeyType:       types.KeyTypeHash,
		}, {
			AttributeName: aws.String("title"),
			KeyType:       types.KeyTypeRange,
		}},
		TableName: aws.String(d.TableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", d.TableName, err)
	} else {

		waiter := dynamodb.NewTableExistsWaiter(d.DynamoDbClient)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(d.TableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc = table.TableDescription
	}
	return tableDesc, err
}
