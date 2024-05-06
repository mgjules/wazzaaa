# Wazzaaa

Sends a whatsapp message to multiple recipients

<!--toc:start-->

- [Wazzaaa](#wazzaaa)
  - [Installation](#installation)
  - [Usage](#usage)

<!--toc:end-->

## Installation

- You can install using the [latest released binary](https://github.com/mgjules/wazzaaa/releases/latest).

- **OR** using Go:

  ```sh
  go install github.com/mgjules/wazzaaa@latest
  ```

## Usage

```sh
❯ ./wazzaaa --help
Usage of ./wazzaaa:
  -message string
      message to send
  -recipients string
      list of comma separated phone numbers
```

Example:

```sh
wazzaaa --recipients "xxxxxxxxxx,yyyyyyyyy" --message "Hello there..."
```

When executed for the first time, it will generate a QR code in the terminal.
Link your WhatsApp account to wazzaaa using the QR code.
More info about how to link a device to your WhatsApp account [here](https://faq.whatsapp.com/1317564962315842/?cms_platform=web).
