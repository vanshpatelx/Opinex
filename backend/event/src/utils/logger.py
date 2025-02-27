# src/utils/logger.py

"""
Logger Utility

Configures application-wide logging format and log levels.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 19, 2025
"""

import logging

# Configure logging
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")

# Create logger instance
logger = logging.getLogger("event-service")
