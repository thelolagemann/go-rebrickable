# go-rebrickable

[![go-rebrickable release (latest SemVer)](https://img.shields.io/github/v/release/thelolagemann/go-rebrickable?sort=semver)](https://github.com/thelolagemann/go-rebrickable/releases)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/thelolagemann/go-rebrickable)

go-rebrickable is a Go client library for accessing the [Rebrickable API V3](https://rebrickable.com/api/v3/docs/).

## Installation

go-rebrickable is compatible with modern Go releases in module mode, to install the latest version:
```shell
go get github.com/thelolagemann/go-rebrickable
```

Or to install a specific version from the list of [versions](https://github.com/thelolagemann/go-rebrickable/releases)
```shell
go get github.com/thelolagemann/go-rebrickable@x.y.z
```

## Usage

*Please note: error handling has been omitted from examples for the sake of brevity*

In order to use the rebrickable API, you must have a valid API key. Generate one [here](https://rebrickable.com/api/) 
if you don't already have one.

```go
package main

import (
	"fmt"
	rbrick "github.com/thelolagemann/go-rebrickable"
)

func main() {
	client := rbrick.NewLEGOClient("API_KEY")
	color, _ := client.Color(212)
	fmt.Println(color.Name)
	
	// Outputs: 
	// Bright Light Blue
}
```

You can pass in a number of options to the `NewLEGOClient` function in order to further configure the client. For example, 
to use a custom HTTP client.

```go
httpClient := &http.Client{Timeout: time.Second * 30}
client, _ := rbrick.NewLEGOClient(apiKey, rbrick.HTTPClient(httpClient))
```

Several endpoints accept additional query parameters in order to filter your search. For example, to use a page size of
5:

```go
colors, _ := client.Colors(rbrick.PageSize(5))
```

## TODOs

* [ ] implement user methods
* [ ] implement all query parameters
* [ ] improve test cases
* [ ] document differences between Client and LEGOClient
* [ ] add wider range of examples#
* [ ] download and query local

### Contributing

All contributions are welcome! For any bug reports/feature requests, feel free to open an issue or submit a pull request.

### License

This project is licensed under the MIT license - see the [LICENSE](https://github.com/thelolagemann/go-rebrickable/blob/master/LICENSE) file for details.