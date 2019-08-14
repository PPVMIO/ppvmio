package svc
import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
)

var svc *s3.S3

func init() {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		log.Println(err)
	}

	svc = s3.New(sess)

}
func RetrievePhotos(n int) []*s3.Object {
	bucket := "photos-ppvmio-public"
	resp, _ := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	return resp.Contents
}