# src/utils/id.py

import os
import time

def generate_unique_id():
    hostname = ''.join(filter(str.isdigit, os.uname()[1]))[:4] or '1234'
    timestamp = int(time.time() * 1000)
    counter_value = generate_unique_id.counter % 1000
    generate_unique_id.counter += 1
    return int(f"{timestamp}{str(counter_value).zfill(3)}{hostname.zfill(4)}")

generate_unique_id.counter = 0
