package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bogosj/tesla"
)

var tokenPath = flag.String("token", "", "path to token file")

func main() {
	flag.Parse()

	// Get the authentication token
	if *tokenPath == "" {
		fmt.Println("--token must be specified")
		os.Exit(1)
	}

	// Write GPX header to stdout
	writeGPXHeader()

	// Handle termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Run the main program
	done := make(chan struct{})
	go func() {
		if err := run(context.Background(), *tokenPath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		close(done)
	}()

	select {
	case <-sigChan:
		// Write GPX footer to stdout upon receiving a termination signal
		writeGPXFooter()
	case <-done:
	}

}

func writeGPXHeader() {
	fmt.Println("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	fmt.Println("<gpx version=\"1.1\" creator=\"Created by Tesla-gps (https://github.com/tijszwinkels/tesla-gps)\" xmlns=\"http://www.topografix.com/GPX/1/1\">")
	fmt.Println("<trk>")
	fmt.Println("<trkseg>")
}

func writeGPXFooter() {
	fmt.Println("</trkseg>")
	fmt.Println("</trk>")
	fmt.Println("</gpx>")
}

func writeTrkpt(driveState *tesla.DriveState) error {
	fmt.Printf("<trkpt lat=\"%v\" lon=\"%v\">\n", driveState.Latitude, driveState.Longitude)
	fmt.Printf("<time>%v</time>\n", time.Unix(driveState.GpsAsOf, 0).UTC().Format(time.RFC3339))
	fmt.Printf("</trkpt>\n")
	return nil
}

func run(ctx context.Context, tokenPath string) error {
	c, err := tesla.NewClient(ctx, tesla.WithTokenFile(tokenPath))

	var id int64 = 0
	if err != nil {
		return err
	}

	v, err := c.Vehicles()
	if err != nil {
		return err
	}

	if id == 0 {
		for i, v := range v {
			if i > 0 {
				fmt.Println("----")
			}
			id = v.ID
		}
	}
	vh, err := c.Vehicle(id)
	if err != nil {
		return err
	}
	vh.Wakeup()
	for {
		driveState, err := vh.DriveState()
		if err != nil {
			fmt.Errorf("Couldn't retrieve drivestate: %v", err)
			continue
		}
		_ = writeTrkpt(driveState)
		time.Sleep(time.Millisecond * 800)
	}

	return nil
}
