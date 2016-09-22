package main

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/codegangsta/cli"
)

func main() {
	r := syscall.Rlimit{}

	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &r)

	fmt.Println("Open file current:", r.Cur)
	fmt.Println("Open file maximum:", r.Max)
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)
	fmt.Println(numcpu, "logical processors")

	startingTime := time.Now().UTC()

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:   "get",
			Usage:  "get objects",
			Action: Download,
		},
	}
	app.Run(os.Args)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)

	fmt.Printf("Native [%v]\nMilliseconds [%d]\nSeconds [%.3f]\n",
		duration,
		duration.Nanoseconds()/1e6,
		duration.Seconds())
}

// Download Comment
func Download(c *cli.Context) {
	handler := &S3Handler{}
	handler.initialize()
	go handler.listObjectsPages()
	handler.getObjectsAsync()
	fmt.Printf("Bytes %d\nFailed %d\nCompleted %d\n", handler.bytes, handler.failureCount, handler.successCount)
}
