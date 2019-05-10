# opwire-testa

> Testing toolkit for building opwire-agent modules

## Usage

### Download `opwire-testa`

To download the latest `opwire-testa` on Linux/macOS/BSD systems, run:

```shell
curl https://opwire.org/opwire-testa/install.sh | bash
```

For other systems:

* Download the relevant [`opwire-testa`](https://github.com/opwire/opwire-testa/releases/latest) release,
* Extract the `opwire-testa` or `opwire-testa.exe` binary from the archive to project home folder (current directory).

### Executing tests

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
* `--tags` (`-g`): Conditional tags for selecting test cases. In above example, `label1`, `label2` are test case inclusion tags, `pending-case1`, `pending-case2` are test case exclusion tags.

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

* `--request` (`-X`): Specifies a custom request method to use when communicating with the HTTP server.
* `--url`: Specifies a URL to fetch.
* `--header` (`-H`): Extra header to include in the request when sending HTTP to a server.
* `--data` (`-d`): Sends the specified data in a POST/PUT/PATCH request to the HTTP server.
* `--export`: Renders this `request` in specific format instead of executing. Currently support only one format: `testcase`.
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

See [LICENSE](LICENSE) to see the full text.
