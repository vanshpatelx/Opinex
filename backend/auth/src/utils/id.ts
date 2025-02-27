// src/utils/id.ts
/**
    Unique ID Generator

    Generates a unique 16-character identifier using:
    - 2-digit system hostname-derived segment.
    - 10-digit Unix timestamp in seconds.
    - 4-digit auto-incrementing counter.

    Features:
    - Ensures uniqueness across calls by combining hostname, timestamp, and a counter.
    - Uses BigInt to handle large numeric identifiers.

    Dependencies:
    - os (System hostname access)

    Author: Vansh Patel (remotevansh@gmail.com)  
    Date: February 27, 2025  
 */

import * as os from 'os';

let counter = 0;

function generateUniqueId(): bigint {
  const timestamp = Math.floor(Date.now() / 1000); // 10-digit timestamp (seconds)
  const counterValue = (++counter % 10000).toString().padStart(4, '0'); // 4-digit counter
  const hostnamePart = os.hostname().replace(/\D/g, '').slice(0, 2) || '12'; // 2-digit hostname

  const uniqueIdString = `${timestamp}${counterValue}${hostnamePart}`; // 16-digit ID
  return BigInt(uniqueIdString);
}

export { generateUniqueId };
