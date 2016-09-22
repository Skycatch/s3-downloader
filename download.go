package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Bucket is the s3 bucket
var (
	Bucket    = aws.String("hardcoded")
	Prefix    = aws.String("hardcoded")
	Directory = "testData"
)

// S3Handler hello world
type S3Handler struct {
	session           *session.Session
	manager           *s3manager.Downloader
	bytes             int64
	bufferChannel     chan *s3.Object
	bufferChannelSize int
	channel           chan *s3.Object
	channelSize       int
	waitGroup         sync.WaitGroup
	connvar           int
	successCount      int
	failureCount      int
}

func (handler *S3Handler) listObjectsPages() error {
	defer handler.waitGroup.Done()

	client := s3.New(session.New())
	params := &s3.ListObjectsV2Input{Bucket: Bucket, Prefix: Prefix}
	err := client.ListObjectsV2Pages(params, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		fmt.Println("Number of items:", len(page.Contents), "Last Page:", lastPage)
		handler.waitGroup.Add(len(page.Contents))
		for _, element := range page.Contents {
			handler.bufferChannel <- element
			handler.channelSize++
		}

		return true
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (handler *S3Handler) getObjectsAsync() {
	concurrency := 50

	go func() {
		for item := range handler.bufferChannel {
			handler.bufferChannelSize++
			handler.channel <- item
		}
	}()

	for i := 0; i < concurrency; i++ {
		go func() {
			for item := range handler.channel {
				handler.getObject(item)
			}
		}()
	}
	handler.waitGroup.Wait()
	close(handler.channel)
}

func (handler *S3Handler) getObject(item *s3.Object) {
	defer handler.waitGroup.Done()
	defer func() {
		handler.bufferChannelSize--
	}()
	_, name := filepath.Split(*item.Key)

	if name == "" {
		fmt.Printf("Key (%s) has no file\n", *item.Key)
	} else {

		file := filepath.Join(Directory, *item.Key)
		if err := os.MkdirAll(filepath.Dir(file), 0775); err != nil {
			panic(err)
		}

		fd, err := os.Create(file)
		if err != nil {
			log.Fatal("Failed to create file", err)
			panic(err)
		}
		defer fd.Close()

		params := &s3.GetObjectInput{Bucket: Bucket, Key: item.Key}
		numBytes, err := handler.manager.Download(fd, params)

		handler.bytes += numBytes
		handler.channelSize--

		if err != nil {
			handler.failureCount++
			fmt.Println("Failed to download file", err)
			return
		}
		handler.successCount++
		fmt.Printf("Downloaded s3://%s/%s to %s size: %d bytes - channels: %d buff: %d\n", *Bucket, *item.Key, file, numBytes, handler.channelSize, handler.bufferChannelSize)
	}
}

func (handler *S3Handler) getObjectOld(item *s3.Object) {
	defer handler.waitGroup.Done()

	path := strings.Split(*item.Key, "/")
	saveName := path[len(path)-1]
	fmt.Println("Downloading", saveName, "Path", path)

	name := "skycatch-processing-jobs"

	file, err := os.Create(saveName)
	if err != nil {
		log.Fatal("Failed to create file", err)
	}
	defer file.Close()

	numBytes, err := handler.manager.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(name),
			Key:    aws.String(*item.Key),
		})
	if err != nil {
		fmt.Println("Failed to download file", err)
		return
	}

	fmt.Println("Downloaded file", file.Name(), numBytes, "bytes", handler.bytes, "total", handler.channelSize, "channels")
	return
}

func (handler *S3Handler) download() bool {
	return true
}

func (handler *S3Handler) initialize() {
	handler.connvar = runtime.NumCPU()
	awsSession, err := session.NewSession()
	handler.session = awsSession

	transport := &http.Transport{
		MaxIdleConnsPerHost: 3000,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	s3Client := s3.New(session.New(), &aws.Config{
		HTTPClient: httpClient,
	})

	handler.bufferChannel = make(chan *s3.Object, 1000)
	handler.channel = make(chan *s3.Object, handler.connvar)
	handler.waitGroup.Add(1)

	handler.manager = s3manager.NewDownloaderWithClient(s3Client, func(d *s3manager.Downloader) {
		d.PartSize = (100 * 5 * 1024 * 1024)
		d.Concurrency = 50
	})
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
