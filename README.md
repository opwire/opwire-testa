# opwire-testa

`opwire-testa` is a simple API testing tool written in golang. It is originally developed to test APIs built by `opwire-agent`. The `opwire-testa` is the most convenient one to work with `opwire-agent`, but is also able to be used with other API services easily.

> `Developed by a programmer for other programmers`

## Usage

### Download `opwire-testa`

To download the latest `opwire-testa` on Linux/macOS/BSD systems, run:

```shell
curl https://opwire.org/opwire-testa/install.sh | bash
```

For other systems:

* Download the relevant [`opwire-testa`](https://github.com/opwire/opwire-testa/releases/latest) release,
* Extract the `opwire-testa` or `opwire-testa.exe` binary from the archive to the home folder of your project.

### Execute tests

#### Command line syntax

```shell
./opwire-testa run \
  --test-dirs=... \
  --test-dirs=... \
  --incl-files=tests/feature-1/*/*.yml \
  --incl-files=tests/feature-2/.* \
  --excl-files=tests/demo/* \
  --excl-files=tests/examples/* \
  --tags="+label1,+label2,-pending-case1,-pending-case2"
```

Command line options:

* `--test-dirs` (`-d`): Directories contain test suite files.
* `--incl-files` (`-i`): File inclusion patterns.
* `--excl-files` (`-e`): File exclusion patterns.
* `--test-name` (`-n`): Test title/name matching pattern.
* `--tags` (`-g`): Conditional tags for selecting test cases. In the above example, `label1`, `label2` are the two tags which include test cases, while `pending-case1`, `pending-case2` exclude test cases. To include test cases, the mandantory is not having any `pending-case1` or `pending-case2` selected.

Use `--help` flag to see more details for arguments:

```shell
./opwire-testa run --help
```

### Generating a testcase from a curl command

#### Illustration

![Generating a testcase flow](https://raw.github.com/opwire/opwire-testa/master/docs/assets/images/generating-a-testcase.png)

#### Step 1. Make an HTTP request with Insomnia

#### Step 2. Run `curl` command with `opwire-testa`

#### Step 3. Generate testcase from the request

#### Step 4. Append the testcase to a testsuite

#### Step 5. Verify the updated testsuite

#### Command line syntax

```shell
./opwire-testa req curl \
--request POST \
--url "http://localhost:17779/-"
--header "name1: value1" \
--header "name2: value2" \
--data '{
  "name": "opwire",
  "url": "https://opwire.org"
}'
--export "testcase"
```

Command line options:

* `--request` (`-X`): Specifies a customized request method to use when communicating with the HTTP server.
* `--url`: Specifies a URL to fetch.
* `--header` (`-H`): Specifies an extra header to include in the request when sending HTTP to the server.
* `--data` (`-d`): Specifies HTTP body in a POST/PUT/PATCH request to the HTTP server.
* `--export`: Renders this `request` in a specific format instead of executing. The only one format supported, currently is `testcase`.
* `--snapshot`: Alias of `--export=testcase`.

Use `--help` flag to see more details for arguments:

```shell
./opwire-testa req curl --help
```

### Extracting curl command from a testcase

#### Command line syntax

```shell
./opwire-testa gen curl \
  --test-dirs=... \
  --test-dirs=... \
  --incl-files=file-inclusion-pattern \
  --excl-files=file-exclusion-pattern \
  --test-name=test-case-name-pattern
```

Use `--help` flag to see more details for arguments:

```shell
./opwire-testa gen curl --help
```

## License

MIT

Refer to [LICENSE](LICENSE) to see the full text.
