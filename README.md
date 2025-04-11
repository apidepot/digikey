# digikey

Go client for accessing the DigiKey API.

[![GoDoc][godoc badge]][godoc link]
[![Go Report Card][report badge]][report card]
[![License Badge][license badge]][LICENSE.txt]

## Overview

The [digikey][] package provides a Go-based client for the [DigiKey
API][dk-api]. To access the [DigiKey API][dk-api] your client application must
be registered with DigiKey and to make production calls to the API, developers
must be a member of an organization. To learn more see the [DigiKey API
Resources][dk-resources].

## Installation

```bash
$ go get github.com/apidepot/digikey
```

## Examples

Examples are available at <https://github.com/apidepot/digikey-examples/>.

## Implementation Status

This library is currently in alpha status and is changing frequently. Not
everything is implemented, including a list of what is implemented.

## Contributing

Contributions are welcome! To contribute please:

1. Fork the repository
2. Create a feature branch
3. Code
4. Submit a [pull request][]

### Testing

Instead of using [GNU Make][make], this project uses [Just][] as its
task/command runner.

Prior to submitting a [pull request][], please run:

```bash
$ just check    # formats, vets, and unit tests the code
$ just lint     # lints code using staticcheck
```

To update and view the test coverage report:

```bash
$ just cover
```

#### Integration Testing

To perform the integration tests run:

```bash
$ make int
```

Prior to doing so, you'll need to create a `config_test.toml` file with your IEX
Cloud API Token and the base URL. It is recommended to use your sandbox token
and the sandbox URL, so as to not be charged credits when running the
integration tests. Sandbox tokens start with `Tpk_` instead of `pk_` for
non-sandbox tokens. Using the sandbox does make integration a little more
difficult, since results are scrambled in sandbox mode.

Example `config_test.toml` file:

```toml
Token = "Tpk_your_iexcloud_test_token"
BaseURL = "https://sandbox.iexapis.com/v1"
```

## License

[digikey][] is released under the MIT license. Please see the
[LICENSE.txt][] file for more information.

[digikey]: https://github.com/apidepot/digikey
[dk-api]: https://developer.digikey.com/
[dk-resources]: https://developer.digikey.com/resources
[godoc badge]: https://godoc.org/github.com/apidepot/digikey?status.svg
[godoc link]: https://godoc.org/github.com/apidepot/digikey
[just]: https://just.systems/
[LICENSE.txt]: https://github.com/apidepot/digikey/blob/master/LICENSE.txt
[license badge]: https://img.shields.io/badge/license-MIT-blue.svg
[make]: https://www.gnu.org/software/make/
[pull request]: https://help.github.com/articles/using-pull-requests
[report badge]: https://goreportcard.com/badge/github.com/apidepot/digikey
[report card]: https://goreportcard.com/report/github.com/apidepot/digikey
