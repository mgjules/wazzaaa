# Wazzaaa

Sends a whatsapp message to multiple recipients

## Installation

- You can install using the [latest released binary](https://github.com/mgjules/wazzaaa/releases/latest).

- **OR** using Go:

    ```shell
    go install github.com/mgjules/wazzaaa@latest
    ```

## Usage

```shell
❯ ./wazzaaa --help
Usage of ./wazzaaa:
  -message string
    	message to send
  -recipients string
    	list of comma separated phone numbers
```

Example:

```shell
wazzaaa --recipients "xxxxxxxxxx,yyyyyyyyy" --message "Hello there..."
```

When executed for the first time, it will generate a QR code in the terminal.
Link your whatsapp account to wazzaaa using the QR code.
More info about how to link a device to your whatsapp account [here](https://faq.whatsapp.com/1317564962315842/?cms_platform=web).
