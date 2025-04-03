"""
Configuration for CAD Calls API.

This file contains configuration parameters for accessing the Police-to-Citizen portal.
Modify these values as needed to connect to different agencies or endpoints.
"""

# Police to Citizen API Configuration
# IMPORTANT: Replace these placeholder values with your own police department's information
BASE_URL = "REPLACE_WITH_YOUR_BASE_URL"  # Example: "https://southmiamipdfl.policetocitizen.com"
AGENCY_ID = 0  # Replace with your agency ID number (Example: 386 for South Miami PD)

# API Endpoints - derived from base URL
API_ENDPOINTS = {
    "cadcalls": f"{BASE_URL}/api/CADCalls/{AGENCY_ID}"
}

# Request settings
# Modify these if needed for specific department APIs
API_SETTINGS = {
    "verify_ssl": False,        # Set to True in production if HTTPS certificates are valid
    "timeout": 30,              # Request timeout in seconds
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    "request_method": "POST"    # (may need to be changed for other departments)
}

# Default request parameters
# These values are used when no command line arguments are provided
DEFAULT_PARAMS = {
    "include_open": True,     # Whether to include open/active calls
    "include_closed": False,  # Whether to include closed/resolved calls
    "take": 30,               # Number of records to retrieve per request
    "skip": 0,                # Number of records to skip (for pagination)
    "search_text": ""         # Text to search for in call records
} 