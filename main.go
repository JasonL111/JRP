// This program was modified by JasonL111 in 2024
// Base on https://github.com/awsdocs/aws-doc-sdk-examples


/*
 * Copyright 2024 JasonL111
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

 package main

 import (
   "context"
   "fmt"
   "log"
   "os"
   "time"
 
   "github.com/aws/aws-sdk-go-v2/aws"
   "github.com/aws/aws-sdk-go-v2/config"
   "github.com/aws/aws-sdk-go-v2/credentials"
   "github.com/aws/aws-sdk-go-v2/service/s3"
   "github.com/joho/godotenv"
 )
 
 // Presigner encapsulates the Amazon S3 presign actions.
 type Presigner struct {
   PresignClient *s3.PresignClient
 }
 
 // GetObject generates a presigned URL for downloading an object from S3-compatible storage.
 func (presigner Presigner) GetObject(
   ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (string, error) {
   request, err := presigner.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
     Bucket: aws.String(bucketName),
     Key:    aws.String(objectKey),
   }, func(opts *s3.PresignOptions) {
     opts.Expires = time.Duration(lifetimeSecs) * time.Second
   })
   if err != nil {
     log.Printf("Couldn't get a presigned request for %v:%v. Here's why: %v\n",
       bucketName, objectKey, err)
     return "", err
   }
   return request.URL, nil
 }
 
 func main() {
   err := godotenv.Load()
   if err != nil {
     log.Fatalf("Error loading .env file: %v", err)
   }
   // Backblaze B2 credentials
   b2KeyID := os.Getenv("B2_KEY_ID")
   b2ApplicationKey := os.Getenv("B2_APPLICATION_KEY")
   b2Endpoint := os.Getenv("B2_ENDPOINT")
   b2Region := os.Getenv("B2_REGION")
   bucketName := os.Getenv("B2_BUCKET_NAME")
   prefix := os.Getenv("PREFIX")
   lifetimeSecs := int64(3600) // Presigned URL valid for 1 hour
 
   // Check for required environment variables
   if b2KeyID == "" || b2ApplicationKey == "" || b2Endpoint == "" || b2Region == "" || bucketName == "" {
     log.Fatal("Missing required environment variables.")
   }
 
   // AWS Config with B2 credentials and custom endpoint
   cfg, err := config.LoadDefaultConfig(context.TODO(),
     config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
       b2KeyID,
       b2ApplicationKey,
       "", // Session token is not used for Backblaze B2
     )),
     config.WithRegion(b2Region),
     config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
       func(service, region string, options ...interface{}) (aws.Endpoint, error) {
         if service == s3.ServiceID {
           return aws.Endpoint{
             URL:           b2Endpoint,
             SigningRegion: b2Region,
           }, nil
         }
         return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
       })),
   )
   if err != nil {
     log.Fatalf("Failed to load configuration: %v", err)
   }
 
   // Initialize S3 and Presign clients with UsePathStyle option
   s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
     o.UsePathStyle = true
   })
   presignClient := s3.NewPresignClient(s3Client)
   presigner := Presigner{PresignClient: presignClient}
 
   // List objects in the bucket with optional prefix
   listObjectsInput := &s3.ListObjectsV2Input{
     Bucket: aws.String(bucketName),
   }
   if prefix != "" {
     listObjectsInput.Prefix = aws.String(prefix)
   }
 
   // Open file to save URLs
   outputFile := "presigned_urls.txt"
   file, err := os.Create(outputFile)
   if err != nil {
     log.Fatalf("Failed to create output file: %v", err)
   }
   defer file.Close()
 
   // Paginated listing and URL generation
   for {
     output, err := s3Client.ListObjectsV2(context.TODO(), listObjectsInput)
     if err != nil {
       log.Fatalf("Failed to list objects in bucket %s: %v", bucketName, err)
     }
 
     for _, item := range output.Contents {
       // Generate a Presigned URL for each object
       url, err := presigner.GetObject(context.TODO(), bucketName, *item.Key, lifetimeSecs)
       if err != nil {
         log.Printf("Failed to generate presigned URL for object %s: %v", *item.Key, err)
         continue
       }
 
       // Write object key and URL to file
       _, err = file.WriteString("Object Key: " + *item.Key + "\n")
       if err != nil {
         log.Fatalf("Failed to write to file: %v", err)
       }
       _, err = file.WriteString("Presigned URL: " + url + "\n\n")
       if err != nil {
         log.Fatalf("Failed to write to file: %v", err)
       }
     }
 
     // If no more pages, break the loop
     if output.IsTruncated == nil || !*output.IsTruncated {
       break
     }
 
     // Set continuation token for next page
     listObjectsInput.ContinuationToken = output.NextContinuationToken
   }
 
   log.Printf("Presigned URLs have been written to %s", outputFile)
 }
 