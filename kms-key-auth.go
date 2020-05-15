package main

import (
    "github.com/Kong/go-pdk"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/kms"
    "github.com/aws/aws-sdk-go/aws"
    "fmt"

    b64 "encoding/base64"
)

type Config struct {
    EncApiKey string
    KmsRegion string
}

func New() interface{} {
    return &Config{}
}

func decryptKey(key string, region string) string {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(region)},
    )
    svc := kms.New(sess)

    data, err := b64.StdEncoding.DecodeString(key)

    result, err := svc.Decrypt(&kms.DecryptInput{CiphertextBlob: data})

    if err != nil {
        fmt.Println("Got error decrypting data: ", err)
    }

    unencrytped := string(result.Plaintext)

    return unencrytped
}

func (conf Config) Access(kong *pdk.PDK) {
    key, headerErr := kong.Request.GetHeader("apikey")

    if headerErr != nil {
        kong.Log.Err("Error getting header")
    }

    encApiKey := conf.EncApiKey

    region := conf.KmsRegion

    // region defaults to us-east-1 if not set
    if region == "" {
        region = "us-east-1"
    }

    unencrytpedKey := decryptKey(encApiKey, region)

    headers := make(map[string][]string)
    headers["Content-Type"] = append(headers["Content-Type"], "application/json")

    if headerErr != nil || unencrytpedKey == "" || key == "" || unencrytpedKey != key {
        kong.Response.Exit(403, "Invalid authentication credentials", headers) 
    }
}
