# Tesla GPS

Tesla GPS is a simple command-line tool to stream your Tesla vehicle's GPS coordinates in real-time and save them in GPX format.

## Vehicle Sleep

API-calls can prevent the vehicle from going to sleep, which would cause the battery of the vehicle to drain while idle at a higher rate than normal (vampire drain). The software will try to prevent this, using a [similar method to TeslaFi](https://support.teslafi.com/knowledge-bases/2/articles/640-enabling-sleep-settings-to-limit-vampire-loss).

If you miss gps-data (or you want your vehicle to stay awake), run the software with the `--wakeup` flag. This will keep the vehicle awake.
If you notice that your vehicle's battery is draining faster than normal, stop the software and file an issue.


## Prerequisites

### Installing Go

If you don't have Go installed, you can download it from the [official Go website](https://golang.org/dl/) and follow the [installation instructions](https://golang.org/doc/install) for your operating system.

Make sure to set your `GOPATH` environment variable and add `$GOPATH/bin` to your `PATH` variable, as described in the [Go documentation](https://golang.org/doc/gopath_code.html#GOPATH).

### Obtaining a Tesla API Token

To use Tesla GPS, you need a Tesla API token. You can obtain one by following the instructions in the [Tesla JSON API (Unofficial) Documentation](https://tesla-api.timdorr.com/api-basics/authentication).

Alternatively, you can use the Android 'Tesla Tokens' app to obtain the token: [Google Play Store](https://play.google.com/store/apps/details?id=net.leveugle.teslatokens&hl=en&gl=US)

Create a JSON file with the following structure, replacing `my-refresh-token` with the refresh token you obtained:

```json
{
  "refresh_token": "my-refresh-token"
}
```

## Build

1. Clone this repository:
```
git clone https://github.com/tijszwinkels/tesla-gps.git
```
2. Change the directory to the repository:
```
cd tesla-gps
```
3. Build the executable:
```
go build
```

## Usage
To stream GPS data and save it to a GPX file, run:
```
./tesla-gps --token /path/to/your/token > output.gpx
```

Replace /path/to/your/token with the path to your Tesla API token file.

To view the live updates of the GPX file while it's being written, open another terminal and run:
```
tail -f output.gpx
```

## License

This project is released under the MIT License.