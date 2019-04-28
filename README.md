# opwire-testa

> Testing toolkit for building opwire-agent modules

## Usage

### Executing tests

#### Command line syntax

```shell
./opwire-testa run \
  --test-dirs=... \
  --test-dirs=... \
  --incl-files=file-name-pattern-or-regexp \
  --excl-files=file-name-pattern-or-regexp \
  --tags="+tag1,+tag2,-skipped-case1,-skipped-case2"
```

Use `help` sub-command to see more details for arguments:

```shell
./opwire-testa help run
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

Use `help` sub-command to see more details for arguments:

```shell
./opwire-testa req help curl
```

### Extracting curl command from a testcase

#### Command line syntax

```shell
./opwire-testa gen curl \
    --test-dirs=... \
    --test-dirs=... \
    --test-file=file-name-pattern \
    --test-name=test-case-name-pattern
```

Use `help` sub-command to see more details for arguments:

```shell
./opwire-testa gen help curl
```

## License

MIT

See [LICENSE](LICENSE) to see the full text.
