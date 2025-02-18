# src/utils/id.py

"""
Unique ID Generator

Generates unique numeric IDs based on:
1. Current timestamp (milliseconds).
2. A counter to avoid duplicate timestamps.
3. The machine's hostname (digits only).

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 19, 2025
"""

import os
import time

def generate_unique_id():
    """
    Generate a globally unique numeric ID.

    - Uses the current timestamp in milliseconds.
    - Includes a rolling counter for uniqueness.
    - Appends a machine-specific identifier.
    
    Returns:
        int: A unique numeric ID.
    
    Last Updated: February 19, 2025
    """
    hostname = ''.join(filter(str.isdigit, os.uname()[1]))[:4] or '1234'
    timestamp = int(time.time() * 1000)
    counter_value = generate_unique_id.counter % 1000
    generate_unique_id.counter += 1
    return int(f"{timestamp}{str(counter_value).zfill(3)}{hostname.zfill(4)}")

# Initialize counter
generate_unique_id.counter = 0
