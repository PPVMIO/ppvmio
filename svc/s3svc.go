package svc

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var svc *s3.S3
var downloader *s3manager.Downloader

func init() {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		log.Println(err)
	}

	downloader = s3manager.NewDownloader(sess)
	svc = s3.New(sess)

}
func RetrieveObjects(prefix string) []*s3.Object {
	bucket := "photos-ppvmio-public"
	resp, _ := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	contents := resp.Contents

	n := 0
	for _, ele := range contents {
		if strings.HasPrefix(*ele.Key, prefix) && strings.Compare(*ele.Key, prefix) != 0 {
			contents[n] = ele
			n++
		}
	}
	contents = contents[:n]
	return contents

}
