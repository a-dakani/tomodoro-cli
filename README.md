# Tomodoro CLI

The Tomodoro CLI is a command-line interface for the [tomodoro.de](https://app.tomodoro.de/) website, offering synchronized pomodoro timers right in your command line.

With this interactive CLI, you can easily fetch teams, start or pause timers, and keep them in sync with other team members.

## Prerequisites

Before using the Tomodoro CLI, ensure that you have the following prerequisites installed on your system:

- Go (version 1.2 or later)

If you don't have Go installed, you can download it from the official Go website: [Download Go](https://golang.org/dl/)

## Installation

Follow the steps below to install the Tomodoro CLI:

1. Clone the repository by running the following command:

   ```shell
   git clone https://github.com/a-dakani/tomodoro-cli.git
   ```

2. Navigate to the project directory:

   ```shell
   cd tomodoro-cli
   ```

3. Install the CLI using the following command (requires sudo access):

   ```shell
   make install
   ```

4. To run the CLI, execute the following command:

   ```shell
   tomodoro
   ```

## Uninstallation

If you decide to remove the Tomodoro CLI, you can do so by running the following command (requires sudo access):

```shell
make uninstall
```

This will uninstall the CLI and remove the config files from your system.

## Demo

![](https://github.com/a-dakani/tomodoro-cli/blob/main/demo.gif?raw=true)
