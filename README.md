# TODO
- [x] Read live telemetry (from a live session or replay)
- [x] Read data from a stored `.ibt` file
- [x] Allow to export the data to an `.ibt` file
- [x] Allow to export the session info data to a `.yaml` file
- [ ] Add the message broadcasting system
- [ ] Explore a more convenient API for fetching the data for the SDK user
- [ ] Change the pattern in which the data is fetched from the telemetry and
how it is exported into `.ibt` files


# About
This project is a simple Go SDK for the popular iRacing racing simulator.
It has the capabilites to:
- Read live data (live session or replay)
- Read data from a `.ibt` telemetry file

It should run on Linux, MacOS and Windows. With the caveat that live sessions
only happen on Windows (that I know about), therefore Linux and MacOS can only
read data from telemetry files.


## Usage
The SDK instance is created by calling `goirsdk.Init(Reader, exportTelem, exportYAML)`
- `Reader` is a variable that implements the interface:
```go
type Reader interface {
	io.Reader
	io.ReaderAt
	io.ReadCloser
}
```
To read data from a `.ibt` file, the user should pass the `*os.File` of it, and
to read live telemetry the user should pass nil

- `exportTelem` should be an empty string if the user doesn't want to export
the data, otherwise pass a string with the path for the destination telemetry
file

- `exportYAML` is just like the exportTelem but for the session info `yaml` data

### Example
```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ESilva15/goirsdk"
)

func msToKph(v float32) int {
	return int((3600 * v) / 1000)
}

func main() {
	// Open the data source file
	file, err := os.Open("/path/to/ibtFile")
	if err != nil {
		log.Fatalf("Failed to open IBT file: %v", err)
	}

	// Instantiate our iRacing SDK instance
	irsdk, err := goirsdk.Init(file, "", "")
	if err != nil {
		log.Fatalf("Failed to create iRacing interface: %v", err)
	}
	defer irsdk.Close()

	// Set up a loop to iterate our data
	mainLoopTicker := time.NewTicker(time.Second / 60)
	defer mainLoopTicker.Stop()

	for {
		// Update the data that the SDK is holding with the next tick
		_, err := irsdk.Update(100 * time.Millisecond)
		if err != nil {
			log.Printf("could not update data: %v", err)
			continue
		}

		// Vehicle Movement data gathered from the names we can find on the
		// telemetry_docs.pdf file
		// - I wish to make this less verbose if possible
		if _, ok := irsdk.Vars.Vars["Gear"]; !ok {
			log.Fatal("Field `Gear` doesn't exist")
		}

		if _, ok := irsdk.Vars.Vars["RPM"]; !ok {
			log.Fatal("Field `RPM` doesn't exist")
		}

		if _, ok := irsdk.Vars.Vars["Speed"]; !ok {
			log.Fatal("Field `Speed` doesn't exist")
		}

		gear := int32(irsdk.Vars.Vars["Gear"].Value.(int))
		rpm := int32(irsdk.Vars.Vars["RPM"].Value.(float32))
		speed := int32(msToKph(irsdk.Vars.Vars["Speed"].Value.(float32)))

		fmt.Printf("\033[?25l\033[2J\033[H")
		fmt.Printf("Gear: %d, RPM: %d, Speed: %d", gear, rpm, speed)

		<-mainLoopTicker.C
	}
}
```


## SharedMem
I vendored in the code from [hidez8891/shm](https://github.com/hidez8891/shm) 
since the repo has been archived. I took the opportunity to update some of its
code.
