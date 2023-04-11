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

var wakeup = flag.Bool("wakeup", false, "wake up the vehicle and keep it awake")
var verbose = flag.Bool("verbose", false, "verbose logging to stderr")
var tokenPath = flag.String("token", "", "path to token file")

const dontSleepAfterDrivingDuration = 30 * time.Minute
const tryToSleepDuration = 15 * time.Minute

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
	// Create client
	client, err := tesla.NewClient(ctx, tesla.WithTokenFile(tokenPath))

	// Get the first vehicle ID
	var id int64 = 0
	if err != nil {
		return err
	}
	vehicles, err := client.Vehicles()
	if err != nil {
		return err
	}
	if id == 0 {
		for _, vehicle := range vehicles {
			id = vehicle.ID
		}
	}
	vehicle, err := client.Vehicle(id)
	if err != nil {
		return err
	}

	if *wakeup {
		vehicle.Wakeup()
	}

	// Get the drive-state, and write it in gpx format
	var prevDriveState *tesla.DriveState

	// Expiry times
	var stayAwakeAfterDrivingExpiry time.Time
	var goToSleepExpiry time.Time

	for {
		time.Sleep(time.Millisecond * 800)

		// Get the vehicle state (Doesn't keep awake)
		vehicle, err := client.Vehicle(id)
		if *verbose {
			fmt.Fprintf(os.Stderr, "Vehicle state %v\n", vehicle.State)
		}

		// Handle desired sleeping behavior of the vehicle
		if !*wakeup {
			now := time.Now()

			if !stayAwakeAfterDrivingExpiry.IsZero() && now.After(stayAwakeAfterDrivingExpiry) {
				if *verbose {
					fmt.Fprintf(os.Stderr, "30m after driving timer expired. Setting 'go to sleep' timer.\n")
				}
				goToSleepExpiry = now.Add(tryToSleepDuration)
				stayAwakeAfterDrivingExpiry = time.Time{}
			}

			if !goToSleepExpiry.IsZero() && now.After(goToSleepExpiry) {
				if *verbose {
					fmt.Fprintf(os.Stderr, "15m 'go to sleep' timer expired. Restarting 'stay awake' timer again.\n")
				}
				goToSleepExpiry = time.Time{}
				stayAwakeAfterDrivingExpiry = now.Add(dontSleepAfterDrivingDuration)
			}

			if vehicle.State == "asleep" {
				goToSleepExpiry = time.Time{}
				stayAwakeAfterDrivingExpiry = time.Time{}
				if *verbose {
					fmt.Fprintf(os.Stderr, "Ssshh. Vehicle is sleeping. Not doing anything until it wakes up.\n")
					time.Sleep(time.Second * 30)
				}
				continue
			} else if !goToSleepExpiry.IsZero() {
				if *verbose {
					fmt.Fprintf(os.Stderr, "Singing lullabies, waiting for the car to go to sleep for 15m.\n")
					fmt.Fprintf(os.Stderr, "Please note; If the car starts driving within these 15m, we might miss it\n")
				}
				continue
			}
		}

		// Get the drive state (Does keep awake)
		driveState, err := vehicle.DriveState()
		if err != nil {
			if *verbose {
				fmt.Fprintf(os.Stderr, "Couldn't retrieve drivestate: %v", err)
			}
			continue
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "Shift state %v\n", driveState.ShiftState)
		}
		if prevDriveState != nil && driveState.Latitude == prevDriveState.Latitude && driveState.Longitude == prevDriveState.Longitude && driveState.GpsAsOf == prevDriveState.GpsAsOf {
			// Skip writing this point if it's identical to the previous one
			continue
		}

		// If the car becomes inactive and the timer isn't running yet, start the sleep timer
		if driveState.ShiftState != "D" && driveState.ShiftState != "R" {
			if stayAwakeAfterDrivingExpiry.IsZero() {
				stayAwakeAfterDrivingExpiry = time.Now().Add(dontSleepAfterDrivingDuration)
				if *verbose {
					fmt.Fprintf(os.Stderr, "Car became inactive. Setting 'stay awake' timer.\n")
				}
			}
			time.Sleep(time.Second * 4)
		} else if (driveState.ShiftState == "D" || driveState.ShiftState == "R") && !stayAwakeAfterDrivingExpiry.IsZero() {
			stayAwakeAfterDrivingExpiry = time.Time{}
			if *verbose {
				fmt.Fprintf(os.Stderr, "Car is active. Stopping 'stay awake' timer.\n")
			}
		}

		_ = writeTrkpt(driveState)
		prevDriveState = driveState
	}

	return nil
}
