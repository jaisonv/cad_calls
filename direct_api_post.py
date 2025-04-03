import json
import requests
import logging
import os
import argparse
import shutil
from datetime import datetime
from config import BASE_URL, AGENCY_ID, DEFAULT_PARAMS, API_SETTINGS, API_ENDPOINTS

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

def make_cadcalls_request(agency_id=AGENCY_ID, include_open=DEFAULT_PARAMS["include_open"], include_closed=DEFAULT_PARAMS["include_closed"], take=DEFAULT_PARAMS["take"], skip=DEFAULT_PARAMS["skip"], search_text=DEFAULT_PARAMS["search_text"]):
    """
    Make a direct POST request to the CADCalls API with specified parameters.
    
    Args:
        agency_id (int): Agency ID to fetch data for
        include_open (bool): Whether to include open calls
        include_closed (bool): Whether to include closed calls
        take (int): Number of records to retrieve
        skip (int): Number of records to skip (for pagination)
        search_text (str): Optional search text to filter calls
        
    Returns:
        tuple: (response object, output filename or None)
    """
    # Create output directory
    os.makedirs("cadcalls_results", exist_ok=True)
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    
    # Generate a more descriptive filename with site info
    site_name = BASE_URL.replace("https://", "").replace("http://", "").split(".")[0]
    
    # Target URL for API call - get the full URL from the config
    api_endpoint = API_ENDPOINTS.get("cadcalls", "")
    if api_endpoint.startswith("/"):
        api_url = f"{BASE_URL}{api_endpoint}"
    else:
        api_url = api_endpoint
        
    logger.info(f"Using API endpoint: {api_url}")
    
    # Create session with browser-like headers
    session = requests.Session()
    session.verify = API_SETTINGS.get("verify_ssl", False)
    
    # Step 1: First visit the main site to get cookies
    logger.info(f"Visiting {site_name} site for cookies")
    session.headers.update({
        'User-Agent': API_SETTINGS.get("user_agent", "Mozilla/5.0"),
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8',
        'Accept-Language': 'en-US,en;q=0.9',
    })
    
    # Try several URL patterns for the main site
    try:
        main_response = session.get(BASE_URL, timeout=API_SETTINGS.get("timeout", 30))
        main_response.raise_for_status()
        logger.info(f"Successfully connected to {BASE_URL}")
        logger.debug(f"Main site cookies: {session.cookies.get_dict()}")
    except requests.exceptions.RequestException as e:
        logger.error(f"Error accessing main site {BASE_URL}: {e}")
        return None, None
    
    # Step 2: Try to visit the CADCalls page using various common patterns
    cad_url = f"{BASE_URL}/CADCalls"
    logger.info(f"Attempting to visit CADCalls page at {cad_url}")
    
    try:
        cadcalls_response = session.get(cad_url, timeout=API_SETTINGS.get("timeout", 30))
        cadcalls_response.raise_for_status()
        logger.info(f"Successfully accessed CADCalls page at {cad_url}")
    except requests.exceptions.RequestException as e:
        logger.warning(f"Error accessing CADCalls page at {cad_url}: {e}")
        # Try alternative URL format that some sites might use
        alternative_cad_url = f"{BASE_URL}/Home/CADCalls"
        try:
            logger.info(f"Trying alternative CADCalls URL: {alternative_cad_url}")
            cadcalls_response = session.get(alternative_cad_url, timeout=API_SETTINGS.get("timeout", 30))
            cadcalls_response.raise_for_status()
            logger.info(f"Successfully accessed alternative CADCalls page at {alternative_cad_url}")
        except requests.exceptions.RequestException as e:
            logger.error(f"Failed to access both CADCalls URLs: {e}")
            logger.warning("Continuing with API request anyway...")
    
    # Update headers for API request
    session.headers.update({
        'Accept': 'application/json, text/javascript, */*; q=0.01',
        'Content-Type': 'application/json',
        'Origin': BASE_URL,
        'Referer': cad_url,
        'X-Requested-With': 'XMLHttpRequest'
    })
    
    # Check for CSRF token
    if 'XSRF-TOKEN' in session.cookies:
        csrf_token = session.cookies['XSRF-TOKEN']
        logger.info(f"Found CSRF token: {csrf_token[:15]}...")
        session.headers.update({
            'X-XSRF-TOKEN': csrf_token
        })
    else:
        logger.warning("No CSRF token found in cookies. Some sites may require this for API access.")
    
    # Request payload exactly as specified - fix the parameters to avoid 500 error
    payload = {
        "IncludeOpenCalls": include_open,
        "IncludeClosedCalls": include_closed,
        "IncludeCount": True,
        "PagingOptions": {
            "SortOptions": [
                {
                    "Name": "StartTime",
                    "SortDirection": "Descending",
                    "Sequence": 1
                }
            ],
            "Take": take,
            "Skip": skip
        },
        "FilterOptionsParameters": {
            "IntersectionSearch": True,
            "SearchText": search_text,
            "Parameters": []  # Leave this as an empty array regardless of search_text
        }
    }
    
    # Make the API request based on configured method
    logger.info(f"Making API request to {api_url}")
    logger.debug(f"Headers: {session.headers}")
    
    request_method = API_SETTINGS.get("request_method", "POST").upper()
    response = None
    
    try:
        if request_method == "POST":
            logger.info(f"Using POST method with payload: {json.dumps(payload, indent=2)}")
            response = session.post(api_url, json=payload, timeout=API_SETTINGS.get("timeout", 30))
        else:  # GET
            # Create query parameters for GET (simplified version of the POST payload)
            params = {
                'includeOpen': str(include_open).lower(),
                'includeClosed': str(include_closed).lower(),
                'take': take,
                'skip': skip,
                'searchText': search_text
            }
            logger.info(f"Using GET method with params: {params}")
            response = session.get(api_url, params=params, timeout=API_SETTINGS.get("timeout", 30))
        
        status_code = response.status_code
        content_type = response.headers.get('Content-Type', '')
        
        logger.info(f"Response status: {status_code}")
        logger.info(f"Response content type: {content_type}")
        
        # Save response regardless of status code
        output_filename = None
        
        if response.text:
            if 'application/json' in content_type:
                try:
                    # Save as JSON with site name in filename
                    output_filename = f"cadcalls_results/{site_name}_cadcalls_{agency_id}_{timestamp}.json"
                    with open(output_filename, "w") as f:
                        json.dump(response.json(), f, indent=2)
                    logger.info(f"Saved JSON response to {output_filename}")
                except ValueError:
                    # Not valid JSON, save as text
                    output_filename = f"cadcalls_results/{site_name}_cadcalls_{agency_id}_{timestamp}.txt"
                    with open(output_filename, "w") as f:
                        f.write(response.text)
                    logger.info(f"Saved text response to {output_filename}")
            else:
                # Save as HTML
                output_filename = f"cadcalls_results/{site_name}_cadcalls_{agency_id}_{timestamp}.html"
                with open(output_filename, "w") as f:
                    f.write(response.text)
                logger.info(f"Saved HTML response to {output_filename}")
        
        # If response was not successful but we got content, add debug info
        if status_code != 200:
            logger.warning(f"API request returned non-200 status: {status_code}")
            debug_filename = f"cadcalls_results/{site_name}_debug_info_{timestamp}.txt"
            with open(debug_filename, "w") as f:
                f.write(f"URL: {api_url}\n")
                f.write(f"Status Code: {status_code}\n")
                f.write(f"Headers Sent:\n{json.dumps(dict(session.headers), indent=2)}\n\n")
                f.write(f"Cookies:\n{json.dumps(dict(session.cookies.get_dict()), indent=2)}\n\n")
                f.write(f"Payload:\n{json.dumps(payload, indent=2)}\n\n")
                f.write(f"Response Headers:\n{json.dumps(dict(response.headers), indent=2)}\n\n")
                f.write(f"Response Content:\n{response.text[:2000]}\n\n")
            logger.info(f"Saved debug information to {debug_filename}")
        
        return response, output_filename
    except Exception as e:
        logger.error(f"Error making API request: {e}")
        return None, None

def display_cad_calls(data):
    """Display CAD calls data in a readable format"""
    if 'CADCalls' in data and isinstance(data['CADCalls'], list):
        calls = data['CADCalls']
        total = data.get('Total', len(calls))
        
        print(f"\n=== CAD Calls ({total} total) ===\n")
        
        if len(calls) == 0:
            print("No calls found matching your criteria.")
            return
            
        for i, call in enumerate(calls, 1):
            start_time = call.get('StartTime', 'Unknown')
            if start_time and 'T' in start_time:
                # Format the time for better readability
                date_part, time_part = start_time.split('T')
                time_part = time_part.split('-')[0]  # Remove timezone part
                start_time = f"{date_part} {time_part}"
            
            print(f"Call #{i}:")
            print(f"  Status: {call.get('CallType', 'Unknown')}")
            print(f"  Time: {start_time}")
            print(f"  Nature: {call.get('Nature', 'Unknown')}")
            print(f"  Address: {call.get('Address', 'Unknown')}")
            print(f"  Agency: {call.get('Agency', 'Unknown')}")
            print(f"  Incident ID: {call.get('IncidentId', 'Unknown')}")
            if call.get('HasLocation'):
                print(f"  Location: ({call.get('Latitude', 'Unknown')}, {call.get('Longitude', 'Unknown')})")
            print()
    else:
        print("No CAD calls data found in response")

if __name__ == "__main__":
    # Parse command-line arguments
    parser = argparse.ArgumentParser(description="Fetch CAD Calls data from Tyler Police-to-Citizen portal")
    parser.add_argument("--agency-id", type=int, default=AGENCY_ID, help=f"Agency ID (default: {AGENCY_ID})")
    parser.add_argument("--open", action="store_true", default=DEFAULT_PARAMS["include_open"], help=f"Include open calls (default: {DEFAULT_PARAMS['include_open']})")
    parser.add_argument("--closed", action="store_true", default=DEFAULT_PARAMS["include_closed"], help=f"Include closed calls (default: {DEFAULT_PARAMS['include_closed']})")
    parser.add_argument("--take", type=int, default=DEFAULT_PARAMS["take"], help=f"Number of records to retrieve (default: {DEFAULT_PARAMS['take']})")
    parser.add_argument("--skip", type=int, default=DEFAULT_PARAMS["skip"], help=f"Number of records to skip (default: {DEFAULT_PARAMS['skip']})")
    parser.add_argument("--search", type=str, default=DEFAULT_PARAMS["search_text"], help="Search text to filter calls")
    parser.add_argument("--output", type=str, help="Output file path (default: auto-generated)")
    parser.add_argument("--quiet", action="store_true", help="Suppress console output")
    
    args = parser.parse_args()
    
    # Set logging level
    if args.quiet:
        logging.getLogger().setLevel(logging.WARNING)
    
    logger.info(f"Starting CAD Calls retrieval for agency {args.agency_id}")
    
    response, output_file = make_cadcalls_request(
        agency_id=args.agency_id,
        include_open=args.open,
        include_closed=args.closed,
        take=args.take,
        skip=args.skip,
        search_text=args.search
    )
    
    if response and response.status_code == 200:
        try:
            data = response.json()
            
            # Display data
            if not args.quiet:
                display_cad_calls(data)
            
            # Save to custom output file if specified
            if args.output and output_file:
                shutil.copy(output_file, args.output)
                logger.info(f"Copied data to {args.output}")
            
            logger.info("CAD calls data successfully retrieved")
            
            # Return success code
            exit(0)
        except ValueError:
            logger.error("Response is not valid JSON")
            exit(1)
    else:
        status = response.status_code if response else "No response"
        logger.error(f"Request failed with status: {status}")
        
        # If we got a response but it failed, try to capture error details
        if response and response.text:
            logger.error(f"Error details: {response.text[:200]}")
            
        exit(1) 