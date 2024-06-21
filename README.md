</div>

<div align="center"><h1>godrivefileuploader</h1></div>

<div align="center">
`godrivefileuploader` is a command-line tool for uploading files to Google Drive, implemented in Go.

It periodically uploads a specified directory to your Google Drive within a defined interval. 

![](https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square)
</div>


## Usage

1. Visit the [Google Drive API Quickstart Guide](https://developers.google.com/drive/api/quickstart/go#set_up_your_environment) to enable the Google Drive API for your project.
2. Complete the steps: _Enable the API_, _Configure the OAuth consent screen_ and _Authorize credentials for a desktop application_
3. During the _Create credentials for a desktop application_ step, save the _credentials.json_ file to your working directory or another directory. Provide the path to the credentials file as an argument to the application using the -c flag.

### Help message

```
gofileuploader [OPTIONS] QUERY...

Application Options:
  -P, --period=      Period of uploading files, such as '1m', '2h'. Default is 1 hour (default: 1h)
  -p, --path=        Absolute path of the directory
  -d, --duration=    Total duration of the uploader (default: 1h)
  -c, --credentials= Path to the credentials file (default: credentials.json)

Help Options:
  -h, --help         Show this help message
```

## Installation

### Prerequisites
* [Go](https://go.dev/doc/install) (version 1.18+)

### Build from source

```bash
git clone https://github.com/berbatt/godrivefileuploader.git
cd godrivefileuploader 
go install
```

Verify the installation by running:
```bash
godrivefileuploader --help
```

## Examples
### Example 1: Basic Usage

Upload a directory every 30 minutes for 2 hours:
```bash
godrivefileuploader -p /path/to/directory -P 30m -d 2h -c /path/to/credentials.json
```

### Example 2: Authorizing the Application
When you run the application for the first time, it will prompt you to authorize access to your Google Drive:
```bash
godrivefileuploader -p /path/to/directory -P 30m -d 2h -c /path/to/credentials.json
```
The output will include a URL that you need to visit to authorize the application:
```bash
Go to the following URL to authorize the application: 
https://accounts.google.com/o/oauth2/...
Refresh Token saved successfully.
```
1. Copy and paste the URL into your web browser or click the link.
2. Follow the prompts to authorize the application.
3. After authorizing, the application will save a refresh token and proceed with the upload process.

## Contributing

Pull requests are welcome.

## License

[MIT](https://choosealicense.com/licenses/mit/)

## Troubleshooting
### Common Issues
- **Invalid Credentials:** Ensure your credentials.json file is correctly placed and the path provided to the -c option is accurate.
- **Insufficient Permissions:** Verify that your Google Drive API settings include the necessary permissions.
- **Network Issues:** Check your internet connection and firewall settings.