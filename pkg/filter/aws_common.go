package filter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Create session in aws using specific aws-profile and aws-region specified as arguments
func createAwsSession(awsProfile string, awsRegion string) (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(awsRegion)},
		Profile: awsProfile,
	})
	CheckErr(err, "Could not create session!")
	return sess, err
}