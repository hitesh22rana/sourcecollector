# sourcecollector

A simple tool to consolidate multiple files into a single `.txt` file. Perfect for feeding your files to AI tools without any fuss.

## Getting Started

### Installation

You can install the `sourcecollector` CLI tool in two ways:

#### 1. Using `go install`

```bash
go install -u github.com/hitesh22rana/sourcecollector@latest
```

This will download and install the `sourcecollector` binary in your `$GOPATH/bin` directory.

#### 2. Running Locally

To run `sourcecollector` locally, follow these steps:

1. Clone this repository to your local machine
    ```bash
    git clone https://github.com/hitesh22rana/sourcecollector.git
    ```

2. Build and run the application using `make`
    ```bash
    make run input=/path/to/input/directory output=/path/to/output/file.txt
    ```

    Replace `/path/to/input/directory` and `/path/to/output/file.txt` with the actual paths you want to use for input and output, respectively.

    Alternatively, you can build and run the application manually:
    ```bash
    $ go build -o bin/sourcecollector cmd/cli/main.go
    $ ./bin/sourcecollector --input /path/to/input/directory --output /path/to/output/file.txt
    ```

### Usage

After installing or building the `sourcecollector` CLI tool, you can run it with the following command:

```bash
sourcecollector --input /path/to/input/directory --output /path/to/output/file.txt
```

Replace `/path/to/input/directory` and `/path/to/output/file.txt` with the actual paths you want to use for input and output, respectively.

The `--input` flag is required and specifies the input directory path. The `--output` flag is optional and specifies the output file path (default is `output.txt`).

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/hitesh22rana/sourcecollector/blob/main/LICENSE) file for details.