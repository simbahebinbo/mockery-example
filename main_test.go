package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/mock"
	"mockery-example/mocks"
)

func TestS3Mock(t *testing.T) {
	mockS3 := &mocks.S3API{}

	mockResultFn := func(input *s3.ListObjectsInput) *s3.ListObjectsOutput {
		output := &s3.ListObjectsOutput{}
		output.SetCommonPrefixes([]*s3.CommonPrefix{
			&s3.CommonPrefix{
				Prefix: aws.String("2017-01-01"),
			},
		})
		return output
	}

	// NB: .Return(...) must return the same signature as the method being mocked.
	//     In this case it's (*s3.ListObjectsOutput, error).
	mockS3.On("ListObjects", mock.MatchedBy(func(input *s3.ListObjectsInput) bool {
		return input.Delimiter != nil && *input.Delimiter == "/" && input.Prefix == nil
	})).Return(mockResultFn, nil)

	listingInput := &s3.ListObjectsInput{
		Bucket:    aws.String("foo"),
		Delimiter: aws.String("/"),
	}
	listingOutput, err := mockS3.ListObjects(listingInput)
	if err != nil {
		t.Fatalf("Error listing keys: %s", err)
	}

	for _, x := range listingOutput.CommonPrefixes {
		t.Logf("common prefix: %+v", *x)
	}
}
