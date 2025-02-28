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

counter = 0

def generate_unique_id():
    """
    Generates a unique 16-digit identifier using:
    - 10-digit Unix timestamp in seconds.
    - 4-digit auto-incrementing counter.
    - 2-digit hostname-derived segment.

    Returns:
        int: A unique 16-digit numeric ID
    """
    global counter
    
    timestamp = int(time.time())  # 10-digit Unix timestamp (seconds)
    counter = (counter + 1) % 10000  # 4-digit counter (rolls over at 9999)
    counter_value = str(counter).zfill(4)  # Zero-padded 4-digit counter
    
    hostname_digits = ''.join(filter(str.isdigit, os.uname().nodename))[:2]
    hostname_part = hostname_digits if hostname_digits else '12'  # Default to '12' if empty
    
    unique_id_str = f"{timestamp}{counter_value}{hostname_part}"  # Construct the ID string
    
    return int(unique_id_str)  # Convert to integer for consistency