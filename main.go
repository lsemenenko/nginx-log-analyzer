/*
  Nginx Log Analyzer

	Copyright (c) 2024 Leonid Semenenko
	https://github.com/lsemenenko/nginx-log-analyzer
	 
	Permission is hereby granted, free of charge, to any person
	obtaining a copy of this software and associated documentation
	files (the “Software”), to deal in the Software without
	restriction, including without limitation the rights to use,
	copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the
	Software is furnished to do so, subject to the following
	conditions:

	The above copyright notice and this permission notice shall be
	included in all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND,
	EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
	OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
	NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
	HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
	WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
	FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
	OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type IPCount struct {
	IP        string
	Count     int
	StartTime time.Time
	EndTime   time.Time
}

func main() {
	// Step 1: Define command-line flags
	logPattern := flag.String("log", "/var/log/nginx/access.log", "Log file pattern")
	matchString := flag.String("match", "wp-admin", "String to match in log lines")
	statusCode := flag.String("status", "200", "HTTP status code to count")
	resultLimit := flag.Int("limit", 10, "Number of top results to display")
	timePeriod := flag.Duration("period", 10*time.Minute, "Time period for grouping (e.g., 10m, 1h)")
	flag.Parse()

	// Step 2: Process logs and store relevant entries
	tempLogs, err := processLogs(*logPattern, *matchString, *statusCode)
	if err != nil {
		fmt.Println("Error processing logs:", err)
		return
	}

	// Step 3: Count status codes within specified time period for each IP
	ipPeriods, ipMaxCount := countStatusCodes(tempLogs, *timePeriod)

	// Step 4: Sort and get top IPs with highest counts
	topIPs := getTopIPs(ipPeriods, ipMaxCount, *resultLimit)

	// Step 5: Print results
	printResults(topIPs, *matchString, *statusCode, *timePeriod)
}

func processLogs(pattern, matchString, statusCode string) ([]string, error) {
	var tempLogs []string

	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		var reader io.Reader
		if strings.HasSuffix(file, ".gz") {
			gzReader, err := gzip.NewReader(f)
			if err != nil {
				return nil, err
			}
			defer gzReader.Close()
			reader = gzReader
		} else {
			reader = f
		}

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, matchString) && strings.Contains(line, " "+statusCode+" ") {
				fields := strings.Fields(line)
				if len(fields) > 4 {
					tempLogs = append(tempLogs, fmt.Sprintf("%s %s", fields[0], strings.Trim(fields[3], "[]")))
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return tempLogs, nil
}

func countStatusCodes(logs []string, period time.Duration) (map[string]int, map[string]IPCount) {
	ipPeriods := make(map[string]int)
	ipMaxCount := make(map[string]IPCount)

	for _, log := range logs {
		fields := strings.Fields(log)
		if len(fields) != 2 {
			continue
		}

		ip := fields[0]
		timestamp, err := time.Parse("02/Jan/2006:15:04:05", fields[1])
		if err != nil {
			continue
		}

		windowStart := timestamp.Truncate(period)
		windowEnd := windowStart.Add(period)
		periodKey := fmt.Sprintf("%s,%d", ip, windowStart.Unix())

		ipPeriods[periodKey]++

		count := ipPeriods[periodKey]
		if maxCount, ok := ipMaxCount[ip]; !ok || count > maxCount.Count {
			ipMaxCount[ip] = IPCount{
				IP:        ip,
				Count:     count,
				StartTime: windowStart,
				EndTime:   windowEnd,
			}
		}
	}

	return ipPeriods, ipMaxCount
}

func getTopIPs(ipPeriods map[string]int, ipMaxCount map[string]IPCount, limit int) []IPCount {
	var topIPs []IPCount
	for _, count := range ipMaxCount {
		topIPs = append(topIPs, count)
	}

	sort.Slice(topIPs, func(i, j int) bool {
		return topIPs[i].Count > topIPs[j].Count
	})

	if len(topIPs) > limit {
		topIPs = topIPs[:limit]
	}

	return topIPs
}

func printResults(topIPs []IPCount, matchString, statusCode string, period time.Duration) {
	fmt.Printf("Top %d IPs with the highest number of %s status codes for %s in a %s period:\n", len(topIPs), statusCode, matchString, period)
	fmt.Println("Rank | IP Address | Max Count | Period")
	fmt.Println("-----|------------|-----------|------------------------")

	for i, ip := range topIPs {
		fmt.Printf("%4d | %-10s | %9d | %s to %s\n",
			i+1, ip.IP, ip.Count,
			ip.StartTime.Format("02/Jan/2006:15:04:05"),
			ip.EndTime.Format("02/Jan/2006:15:04:05"))
	}
}