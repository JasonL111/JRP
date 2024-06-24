package main

import (
  "fmt"
  "os"
  "strconv"
  "time"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/joho/godotenv"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    fmt.Fprintln(os.Stderr, "failed to load .env file:", err)
  }

  newSession, err := session.NewSession(&aws.Config{
    Credentials: credentials.NewStaticCredentials(
      os.Getenv("S3_KEY_ID"),
      os.Getenv("S3_SECRET"),
      "",
    ),
    Region:           aws.String(os.Getenv("S3_REGOIN")),
    Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
    S3ForcePathStyle: aws.Bool(true),
  })
  if err != nil {
    fmt.Fprintln(os.Stderr, "failed to create session:", err)
    return
  }

  bucket := os.Getenv("S3_BUCKET")
  prefix := os.Getenv("S3_PREFIX")
  durationHours, err := strconv.Atoi(os.Getenv("S3_DURATION_HOURS"))
  if err != nil {
    fmt.Fprintln(os.Stderr, "S3_DURATION_HOURS is not an integer:", err)
    return
  }

  client := s3.New(newSession)

  var continuationToken *string = nil

  for {
    listObjectsV2Output, err := client.ListObjectsV2(&s3.ListObjectsV2Input{
      ContinuationToken: continuationToken,
      Prefix:            aws.String(prefix),
      Bucket:            aws.String(bucket),
    })
    if err != nil {
      fmt.Fprintln(os.Stderr, "failed to list objects:", err)
      return
    }

    for _, object := range listObjectsV2Output.Contents {
      getObjectReq, _ := client.GetObjectRequest(&s3.GetObjectInput{
        Key:    object.Key,
        Bucket: aws.String(bucket),
      })

      presign, err := getObjectReq.Presign(time.Hour * time.Duration(durationHours))
      if err != nil {
        panic(err)
      }
      fmt.Println(*object.Key)
      fmt.Println()
      fmt.Println(presign)
      fmt.Println()
    }

    if !(*listObjectsV2Output.IsTruncated) {
      // all objects are listed
      return
    }

    continuationToken = listObjectsV2Output.ContinuationToken
  }

}
