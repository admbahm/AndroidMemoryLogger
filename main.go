package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "time"
)

// --- Configuration ---
const (
    // The package name of the app you want to monitor
    packageName = "com.google.android.youtube"
    // The name of the output CSV file
    outputFile = "memory_log_go.csv"
    // How many seconds to wait between each poll
    pollInterval = 5 * time.Second
)

// main is the entry point of our program
func main() {
    // 1. Create the CSV file and write the header
    file, err := os.Create(outputFile)
    if err != nil {
        log.Fatalf("❌ Failed to create output file: %s", err)
    }
    defer file.Close() // Ensure the file is closed when main exits

    writer := csv.NewWriter(file)
    header := []string{"Timestamp", "PSS_Total_kB", "Java_Heap_kB", "System_Free_RAM_kB"}
    if err := writer.Write(header); err != nil {
        log.Fatalf("❌ Failed to write header to CSV: %s", err)
    }
    writer.Flush() // Immediately write the header to the file

    fmt.Printf("✅ Starting memory log for '%s'.\n", packageName)
    fmt.Printf("   Logging to '%s' every %s. Press Ctrl+C to stop.\n", outputFile, pollInterval)

    // 2. Start the infinite loop to poll and log data
    for {
        // Get the full memory dump from the device
        meminfoOutput, err := exec.Command("adb", "shell", "dumpsys", "meminfo", packageName).Output()
        if err != nil {
            log.Printf("⚠️ Could not get app meminfo (is app running?): %s", err)
            time.Sleep(pollInterval) // Wait before retrying
            continue
        }
        
        // Also get the system-wide memory dump for Free RAM
        sysMeminfoOutput, err := exec.Command("adb", "shell", "dumpsys", "meminfo").Output()
        if err != nil {
            log.Printf("⚠️ Could not get system meminfo: %s", err)
            time.Sleep(pollInterval) // Wait before retrying
            continue
        }

        // 3. Extract the specific metrics using regular expressions
        pssTotal := extractMetric(string(meminfoOutput), `TOTAL\s+(\d+)`)
        javaHeap := extractMetric(string(meminfoOutput), `Java Heap:\s+(\d+)`)
        freeRam := extractMetric(string(sysMeminfoOutput), `Free RAM:\s+(\d+)`)
        
        // Get the current timestamp
        timestamp := time.Now().Format("2006-01-02 15:04:05")

        // 4. Write the new data row to the CSV
        record := []string{timestamp, pssTotal, javaHeap, freeRam}
        if err := writer.Write(record); err != nil {
            log.Printf("⚠️ Failed to write record to CSV: %s", err)
        }
        writer.Flush() // Ensure the data is written to the file immediately

        fmt.Println(strings.Join(record, ","))

        // Wait for the next interval
        time.Sleep(pollInterval)
    }
}

// extractMetric searches through a text block for a line matching a regex pattern
// and returns the first captured numeric group.
func extractMetric(output string, pattern string) string {
    re := regexp.MustCompile(pattern)
    matches := re.FindStringSubmatch(output)
    if len(matches) > 1 {
        return matches[1] // The first capture group is our number
    }
    return "0" // Return 0 if not found
}
