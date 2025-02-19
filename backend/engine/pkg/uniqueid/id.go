package uniqueid

import (
	"fmt"
	"math/big"
	"os"
	"regexp"
	"sync/atomic"
	"time"
)

var counter uint64 = 10000000000000000000 // Start from a 20-digit number

// GenerateBaseId creates a unique, **fixed 20-digit** identifier using timestamp & counter
func GenerateBaseId() *big.Int {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "1234"
	}

	// Extract numeric digits from hostname
	re := regexp.MustCompile("\\D")
	hostDigits := re.ReplaceAllString(hostname, "")
	if len(hostDigits) > 4 {
		hostDigits = hostDigits[:4]
	}
	if len(hostDigits) == 0 {
		hostDigits = "1234"
	}
	hostDigits = fmt.Sprintf("%04s", hostDigits)

	// Generate a **20-digit ID** using a counter
	timestamp := time.Now().UnixMilli() % 10000000000 // Keep last 10 digits
	counterValue := atomic.AddUint64(&counter, 1) % 10000000000

	uniqueIdStr := fmt.Sprintf("%d%010d", timestamp, counterValue)
	uniqueId, _ := new(big.Int).SetString(uniqueIdStr, 10)
	return uniqueId
}
