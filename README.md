# qlcplus-http-api

`qlcplus-http-api` is a small Go application to provide an HTTP API for QLC+ using its websocket API.

Each time a request is made to one of `qlcplus-http-api`'s HTTP endpoints it makes a websocket connection to QLC+,
writes a request based on the endpoint URL, and reads the response to return.

## Usage

Download the latest binary for your platform from the
[GitHub releases](https://github.com/wwwil/qlcplus-http-api/releases), or [build](#Building) it yourself, then run:

```bash
./qlcplus-http-api --qlcplus localhost:9999 --http localhost:8888
```

The homepage will show a table of IDs and names for discovered QLC+ objects, currently only virtual console widgets are
supported.

The status of a widget can then be retrieved wither by ID or name with an HTTP request, for example using `curl`:

```bash
$ curl localhost:8888/widgets/id/22
255
$ curl localhost:8888/widgets/name/Full%20Brightness
255
```

Note the use of `%20` for space ` ` characters in widget names.

The widget can then be interacted with using `curl` to make a HTTP POST request:

```bash
$ curl --data "255" localhost:8888/widgets/name/Full%20Brightness
BUTTON|255
```

The response will show the widget type and the value sent.

Currently only buttons and sliders are supported. The interaction behaviour depends on the widget type; buttons only
support sending a value of `255` to indicate a button press, sliders can be sent values between `0` and `255` to set a
slider position.

## Installation

The `qlcplus-http-api` can be run as a systemd unit. Copy the binary to the correct location, add the provided systemd
unit file, then reload and start the service:

```bash
cp ./qlcplus-http-api /usr/local/bin/qlcplus-http-api
cp ./deploy/qlcplus-http-api.service /etc/systemd/system/qlcplus-http-api.service
systemctl daemon-reload
systemctl enable qlcplus-http-api
systemctl start qlcplus-http-api
```

Modify the systemd unit file as required to pass the correct argument values to `qlcplus-http-api`.

View and follow the logs with:

```bash
journalctl -f -u qlcplus-http-api
```

## Building

Use the [`make.sh`](./make.sh) script to build for your platform:

```bash
./make build
``` 

## Development

Refer to the [QLC+ Web API test page](https://www.qlcplus.org/Test_Web_API.html) for more details on how to interact
with QLC+ over a websocket. You may need to save the page locally for it to work correctly.

Known issues:

- When setting a widget to a value it already has QLC+ will not respond. This mostly affects slider widgets as buttons
  are usually only pressed briefly. This causes the read request to hang until it hits the timeout, it does not affect
  other requests.
