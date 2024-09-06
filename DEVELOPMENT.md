# Development

## Prerequisites

To get started with development, ensure you have the following installed:

- **Go** version 1.23

## Setup

1. **Install Go 1.23:**

   Follow the instructions on the [official Go website](https://golang.org/doc/install) to install Go 1.23.

2. **Install `clitest`:**

   Follow the instructions on the [official clitest homepage](https://github.com/aureliojargas/clitest) to install clitest.

## Running tests

We have go unit tests but also a small acceptance test suite.

1. **Running unit tests**
  
  They are placed in `./tests` folder, so please run:

  ```bash
  go test -v ./tests/...
  ```

2. **Running acceptance tests**

  ```bash
  clitest examples/**
  ```

## Building kudasai

To build a working kudasai binary, you have several alternatives.

1. **Build the code locally**

  This will generate a binary in the root folder of the project.

  ```bash
  go build -o kudosai main.go
  ```

2. **Install the package**

  This will generate a binary in your `$GOBIN` folder, or in `$HOME/go/bin` by default.

  ```bash
  go install
  ```
