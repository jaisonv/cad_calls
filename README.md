# CAD Calls API Client

This script allows you to fetch Computer Aided Dispatch (CAD) calls from police departments using the Police-to-Citizen portal.

> **⚠️ Important Note:** This script does not work with all Police-to-Citizen portal departments.

> **⚠️ Legal Disclaimer:** This tool is provided for informational and educational purposes only. The data accessed is publicly available information. Users are responsible for using this tool in compliance with applicable laws and the terms of service of the websites being accessed. Do not use this tool for harassment, stalking, or any illegal activities. The authors are not responsible for any misuse or consequences thereof.

## Configuration

The application uses a configuration file (`config.py`) to store settings for accessing police department portals.

### Basic Configuration

Set the `BASE_URL` and `AGENCY_ID` in `config.py` to the desired police department.

```python
# Base settings in config.py
BASE_URL = "https://southmiamipdfl.policetocitizen.com"  # Base URL for the police department
AGENCY_ID = 386  # Agency ID for the specific police department
```

To use a different police department, simply change these values in `config.py`. No need to modify the script.

### Advanced Configuration

You can also configure additional settings in `config.py`:

```python
# Request settings for handling API connections
API_SETTINGS = {
    "verify_ssl": False,        # Set to True in production if HTTPS certificates are valid
    "timeout": 30,              # Request timeout in seconds
    "user_agent": "Mozilla/5.0" # Browser user agent string
}

# Default request parameters
DEFAULT_PARAMS = {
    "include_open": True,     # Whether to include open/active calls
    "include_closed": False,  # Whether to include closed/resolved calls
    "take": 30,               # Number of records to retrieve per request
    "skip": 0,                # Number of records to skip (for pagination)
    "search_text": ""         # Text to search for in call records
}
```

## Usage

First, create, activate the virtual environment and install the requirements:

```bash
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

Basic usage:

```bash
python direct_api_post.py
```

Command line options:

```
--agency-id INT    Agency ID (default: from config.py)
--open             Include open calls (default: True)
--closed           Include closed calls (default: False)
--take INT         Number of records to retrieve (default: 30)
--skip INT         Number of records to skip (default: 0)
--search TEXT      Search text to filter calls
--output FILE      Output file path (default: auto-generated)
--quiet            Suppress console output
```

Examples:

```bash
# Get the last 50 calls for the configured agency
python direct_api_post.py --take 50

# Get closed calls only (with no open calls)
python direct_api_post.py --closed

# Get both open and closed calls
python direct_api_post.py --closed --open

# Search for specific call types
python direct_api_post.py --search "Traffic"

# Output to a specific file
python direct_api_post.py --output my_calls.json
```

## Output

Results are saved in the `cadcalls_results` directory with filenames that include the site name and timestamp. You can also specify a custom output file path with the `--output` option.

The script creates detailed debug information files when API calls fail, which can help diagnose issues with specific departments.

## Troubleshooting

If you encounter issues connecting to a specific police department:

1. Check that the `BASE_URL` and `AGENCY_ID` are correct
2. Look at the debug information in the `cadcalls_results` directory
3. Try increasing the `timeout` value in the `API_SETTINGS` configuration
4. Some departments may require SSL verification - try setting `verify_ssl` to `True`
5. If you get a 500 error, try changing the `request_method` to "GET" in API_SETTINGS

### Known Compatibility Issues

Some police departments have implemented enhanced security measures that block automated requests. In our testing:

- **Working departments**: Tyler PD, South Miami PD, Wood County Sheriff
- **Non-working departments**: Lynchburg PD (blocked by WAF/firewall)

If you encounter "Request Rejected" errors or consistent 500 status codes, the department may be using security measures like:
- Web Application Firewalls (WAF)
- Bot detection
- Request fingerprinting
- CAPTCHA or other human verification

For these departments, you may need to use a browser automation tool like Selenium instead of direct API requests.

## Requirements

- Python 3.6+
- Requests library 

## Ethical Use Guidelines

This tool is designed to access publicly available data in an automated way. Please follow these guidelines:

1. **Respect Rate Limits**: Avoid making excessive requests that could burden the servers. 

2. **Public Data Only**: Only access data that is intentionally made public through these portals.

3. **Privacy Considerations**: While this data is public, be mindful that it may contain sensitive information about incidents and locations.

4. **Legal Compliance**: Ensure your use complies with all local, state, and federal laws regarding data access and use.

5. **Terms of Service**: Using this tool should not violate the terms of service of the websites being accessed.

6. **Academic and Journalistic Use**: This tool is well-suited for research, journalism, or public safety awareness purposes.

Remember that CAD data represents real emergency situations and should be treated with appropriate sensitivity. 