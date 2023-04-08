# Tesla GPS

Tesla GPS is a simple command-line tool to stream your Tesla vehicle's GPS coordinates in real-time and save them in GPX format.

## ⚠️ Warning

This tool keeps your Tesla vehicle awake while running, which will prevent it from going into sleep mode and will increase phantom drain. This issue will be addressed in a later release.

## Prerequisites

To use Tesla GPS, you need a Tesla API token. You can obtain one by following the instructions in the [Tesla JSON API (Unofficial) Documentation](https://tesla-api.timdorr.com/api-basics/authentication).

Alternatively, you can use the Android 'Tesla Tokens' app to obtain the token: [Tesla Tokens app in the Google Play Store](https://play.google.com/store/apps/details?id=net.leveugle.teslatokens&hl=en&gl=US).

Create a JSON file (f.e. `tesla.token`), with the following structure, replacing `my-refresh-token` with the refresh token you obtained:

```json
{
        "refresh_token": "my-refresh-token"
}
```

## Installation

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