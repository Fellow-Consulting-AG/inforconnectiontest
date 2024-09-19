# Infor-Test CLI Application

`Infor-Test` is a simple Go-based CLI application designed to fetch an OAuth 2.0 access token from a specified authorization server using credentials stored in a `.ionapi` configuration file. This tool is useful for applications interacting with the Infor CloudSuite API using service account credentials.

The project is set up with GitHub Actions to automate the build and release process for both macOS and Windows binaries. Each time you push to the `main` branch, new binaries are created and released on GitHub.

## Features

- Cross-platform builds for macOS and Windows.
- Automatically increments version numbers starting from `v1.0.0`.
- Automates the release process via GitHub Actions.
- Reads credentials from a `.ionapi` file for fetching an access token.

## Prerequisites

- **Go (Golang)**: Install Go from [the official Go website](https://golang.org/doc/install).
- **GitHub**: To utilize the GitHub Actions workflow, ensure your project is hosted on GitHub.

## How to Use the Application

### Step 1: Prepare the `.ionapi` File

The `.ionapi` file is a JSON file that contains the service account credentials and endpoints needed to authenticate with the Infor authorization server. The structure of the `.ionapi` file is as follows:

```json
{
  "ci": "Your-Client-ID",
  "cs": "Your-Client-Secret",
  "iu": "https://mingle-ionapi.inforcloudsuite.com",
  "pu": "https://mingle-sso.inforcloudsuite.com:443/Your-Tenant/as/",
  "oa": "authorization.oauth2",
  "ot": "token.oauth2",
  "or": "revoke_token.oauth2",
  "saak": "Your-Service-Account-Access-Key",
  "sask": "Your-Service-Account-Secret-Key"
}
```


•	ci: Client ID
•	cs: Client Secret
•	iu: Base URL for calling the ION API Gateway
•	pu: Base URL for calling the authorization server
•	ot: Path for accessing tokens (combined with pu)
•	saak: Service Account Access Key (used as username)
•	sask: Service Account Secret Key (used as password)

### Step 2: Build or Download the Binary

Option 1: Download Pre-Built Binaries (Recommended)

The GitHub Actions workflow automatically builds and releases macOS and Windows binaries after each push to the main branch.

	1.	Visit the Releases page on GitHub.
	2.	Download the appropriate binary for your operating system:
	•	macOS: Infor-test-mac
	•	Windows: Infor-test-windows.exe

Step 2: Run the Application

	1.	Place your .ionapi file in the same directory as the binary or provide its path when running the binary.
	2.	Run the binary with the .ionapi file as an argument.
	•	On macOS

```bash
    ./Infor-test-mac ./path/to/your/INFOR-DOC2.ionapi
```

	On Windows:

```bash
Infor-test-windows.exe ./path/to/your/INFOR-DOC2.ionapi
```

Example Output

```
Loaded .ionapi file with the following values:
Client ID: ******
Client Secret: ******
Username (SAAK): ******
Password (SASK): ******
Access Token URL: https://mingle-sso.inforcloudsuite.com:443/Your-Tenant/as/token.oauth2
Sending POST request to: https://mingle-sso.inforcloudsuite.com:443/Your-Tenant/as/token.oauth2
Form Data: grant_type=password&scope=&username=******&password=******
Access Token: ******
✅ Connection successful! Access token obtained successfully.
```