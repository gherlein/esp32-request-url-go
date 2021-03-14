package main

import (
    "context"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/s3"
    "log"
    "time"
    "github.com/google/uuid"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/iotdataplane"
)

type MyEvent struct {
        Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {

    sess, err := session.NewSession(&aws.Config{
     	      Region: aws.String("us-west-1")},
	      
    )	      
    svc := s3.New(sess)

  
    id := uuid.New()
    file := id.String()    

    req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
        Bucket: aws.String("esp32-rekognition-497939524935"),
        Key:    aws.String(file),
    })
    urlStr, err := req.Presign(15 * time.Minute)

    if err != nil {
        log.Println("Failed to sign request", err)
    }

    log.Println("file: ", file)
    log.Println("URL:", urlStr)

    log.Println("starting mqtt")

    sessiot, err := session.NewSession(&aws.Config{
     	      Region: aws.String("us-west-1"),
	      Endpoint: aws.String("a11cxor2mooi1v-ats.iot.us-west-1.amazonaws.com")},   
    )	      

//    mqtt := iotdataplane.New(session.New(), aws.NewConfig().WithRegion("us-west-1"))
    //mqtt := iotdataplane.New(session.New())

    mqtt := iotdataplane.New(sessiot)
    params := &iotdataplane.PublishInput{
		Topic:   aws.String("esp32/sub/url"), // Required
		Payload: []byte(urlStr),
		Qos:     aws.Int64(1),
     }
     resp, err := mqtt.Publish(params)
     log.Println("resp: ",resp)

     if err != nil {
	log.Println("error: ",err.Error())
     }

    return urlStr, nil
}

func main() {
        lambda.Start(HandleRequest)
}
