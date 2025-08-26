Go ADB MemLogger
A lightweight Go utility to continuously log memory usage from a connected Android device. This tool is designed to run alongside automated tests (like UIAutomator) to help developers correlate specific user actions with their impact on application and system memory.

The output is a simple, timestamped CSV file, perfect for graphing and analysis.

Key Features
Continuous Polling: Gathers memory snapshots at a configurable interval.

Detailed Metrics: Captures key memory indicators, including an app's Total PSS, Java Heap, and the system's available RAM.

CSV Output: Generates a universally compatible memory_log.csv file, ready for import into Excel, Google Sheets, or any data analysis tool.

Minimal Performance Impact: Uses the native Android dumpsys tool, ensuring the measurement process has a negligible effect on the device's performance.

Ideal for Automated UI Testing: Designed to be run in the background to provide a clear performance log for a full UI test suite.

Prerequisites
Before you begin, ensure you have the following installed and configured on your system:

Go: Version 1.18 or later.

Android SDK Platform-Tools: The adb command must be available in your system's PATH.

A Connected Android Device: With USB Debugging enabled.

Installation & Setup
Place the Code: Save the Go program into a file named main.go.

Initialize Go Module: From your terminal, in the same directory as your main.go file, run the following command to initialize the Go module.

Bash

go mod init adb-memlogger
Usage
Basic Execution
Configure the Logger: Open main.go and modify the constants at the top of the file to set your target application's package name and the desired polling interval.

Go

const (
    // The package name of the app you want to monitor
    packageName = "com.your.package.name"
    // How many seconds to wait between each poll
    pollInterval = 5 * time.Second
)
Run the Program: With your Android device connected, run the logger from your terminal:

Bash

go run .
The program will start logging memory data to memory_log_go.csv in the same directory. Press Ctrl+C to stop.

Integrating with UIAutomator Tests
This is the primary use case. The following workflow allows you to automatically capture memory data for the exact duration of your test run.

Bash

#!/bin/bash

echo "✅ Starting memory logger in the background..."
# Start the Go logger and send it to the background
go run . &
# Save the Process ID (PID) of the logger
LOGGER_PID=$!

echo "▶️ Running UIAutomator test suite..."
# Execute your UI test command (replace with your actual command)
adb shell am instrument -w com.example.myapp.test/androidx.test.runner.AndroidJUnitRunner

echo "⏹️ UI tests finished. Stopping memory logger..."
# Stop the background logger using its PID
kill $LOGGER_PID

echo "📊 Memory log saved to memory_log_go.csv. Analysis complete."
Output File (memory_log_go.csv)
The script generates a CSV file with the following columns:

Timestamp    PSS_Total_kB    Java_Heap_kB    System_Free_RAM_kB
2025-08-25 22:10:53    154321    65536    2345678
2025-08-25 22:10:58    158910    71234    2310987
...    ...    ...    ...

Export to Sheets
This data can be easily graphed to visualize memory usage over the course of your test.

Future Enhancements (TODO)
[ ] Implement command-line flags to specify the package name and poll interval at runtime, removing the need to edit the source code.

[ ] Create a master shell script (run_test.sh) to automate the entire process.

[ ] Add CPU usage monitoring by parsing dumpsys cpuinfo.

[ ] Add support for targeting a specific device by its serial number if multiple devices are connected.

License
This project is licensed under the MIT License.
