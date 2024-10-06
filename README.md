# Nginx Log Analyzer

Nginx Log Analyzer is a powerful and flexible Go-based tool designed to analyze Nginx access logs. It helps identify patterns, potential security issues, and high-traffic IPs based on customizable criteria.

## Features

- Analyze compressed (.gz) or uncompressed Nginx access log files
- Filter log entries based on custom string matches and HTTP status codes
- Group log entries by IP address within customizable time periods
- Identify top IPs with the highest number of matching requests
- Flexible command-line options for easy customization of analysis parameters

## Installation

1. Ensure you have Go installed on your system. If not, download and install it from [golang.org](https://golang.org/).

2. Clone this repository:
   ```
   git clone https://github.com/lsemenenko/nginx-log-analyzer.git
   cd nginx-log-analyzer
   ```

3. Build the project:
   ```
   go build -o nginx-log-analyzer
   ```

## Usage

Run the analyzer with the following command:

```
./nginx-log-analyzer [options]
```

### Options

- `-log`: Log file pattern (default: "/var/log/nginx/access.log")
- `-match`: String to match in log lines (default: "wp-admin")
- `-status`: HTTP status code to count (default: "200")
- `-limit`: Number of top results to display (default: 10)
- `-period`: Time period for grouping (e.g., 10m, 1h) (default: 10m)

### Example

To analyze logs for WordPress admin login attempts (status 200), grouped in 1-hour periods, and display the top 10 results:

```
./nginx-log-analyzer -log="/path/to/logs/*.log" -match="wp-admin" -status="200" -limit=10 -period=1h
```

## Output

The tool will display a table with the following information:

- Rank: Position in the list of top IPs
- IP Address: The IP address with high activity
- Max Count: The maximum number of matching requests within the specified time period
- Period: The time range during which the max count occurred

## Versioning
This project follows Semantic Versioning (SemVer).
Version format: MAJOR.MINOR.PATCH (e.g., 1.2.3)

- MAJOR version increments denote incompatible API changes.
- MINOR version increments denote added functionality in a backwards-compatible manner.
- PATCH version increments denote backwards-compatible bug fixes.

## Contributing

Contributions to the Nginx Log Analyzer are welcome! Please feel free to submit pull requests, create issues or spread the word.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.