#!/bin/sh
set -e

apk add --no-cache python3 py3-pip sqlite
pip3 install --quiet --break-system-packages --root-user-action=ignore requests

if [ "${VERIFY_SSL:-false}" = "true" ]; then
  VERIFY_SSL_PY="True"
else
  VERIFY_SSL_PY="False"
fi

cat > /cad-calls/config.py << EOF
"""
Configuration for CAD Calls API.

This file contains configuration parameters for accessing the Police-to-Citizen portal.
Modify these values as needed to connect to different agencies or endpoints.
"""

# Police to Citizen API Configuration
BASE_URL = "${CAD_BASE_URL}"
AGENCY_ID = ${CAD_AGENCY_ID}

# API Endpoints - derived from base URL
API_ENDPOINTS = {
    "cadcalls": f"{BASE_URL}/api/CADCalls/{AGENCY_ID}"
}

# Request settings
API_SETTINGS = {
    "verify_ssl": ${VERIFY_SSL_PY},
    "timeout": 30,
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    "request_method": "POST"
}

# Default request parameters
DEFAULT_PARAMS = {
    "include_open": True,
    "include_closed": False,
    "take": 30,
    "skip": 0,
    "search_text": ""
}
EOF

chmod +x /cadbot

exec /cadbot \
  -token "${TELEGRAM_BOT_TOKEN}" \
  -python-script /cad-calls/direct_api_post.py \
  -db /data/cadbot.db \
  -interval "${CHECK_INTERVAL:-5}"
