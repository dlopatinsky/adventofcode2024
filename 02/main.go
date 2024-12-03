package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "strconv"   
    "strings"
)

type Report []int

func CountSafeReports(scanner bufio.Scanner, maxUnsafeLevels int) (int, error) {
    count := 0
    for scanner.Scan() {
	line := scanner.Text()
	slice := strings.Fields(line)
	report := make(Report, len(slice))
	for i, v := range slice {
	    n, err := strconv.Atoi(v)
	    if err != nil {
		return 0, fmt.Errorf("Invalid number: %s", v)
	    }
	    report[i] = n
	}
	if report.UnsafeLevelCount() <= maxUnsafeLevels {
	    count++
	}
    }
    return count, nil
}

func (r Report) UnsafeLevelCount() int {
    length := len(r)
    if length < 2 {
	return 0
    }
    diffs := make([]int, length - 1)
    for i := 0; i < length - 1; i++ {
	diffs[i] = r[i + 1] - r[i]
    }

    return GetUnsafeDiffCount(diffs)
}

func GetUnsafeDiffCount(diffs []int) int {
    var unsafe []int
    for i, d := range diffs {
	if Abs(d) < 1 || Abs(d) > 3 {
	    unsafe = append(unsafe, i)
	}
    }
    pos, neg := GetPosAndNegDiffIndices(diffs)
    if len(neg) > len(pos) {
	unsafe = append(unsafe, pos...)
    } else {
	unsafe = append(unsafe, neg...)
    }
    return len(unsafe)
}

func Abs(i int) int {
    if i < 0 {
	return -i
    } else {
	return i
    }
}

func GetPosAndNegDiffIndices(diffs []int) ([]int, []int) {
    var pos []int
    var neg []int
    for i, d := range diffs {
	if d < 0 {
	    neg = append(neg, i)
	} else {
	    pos = append(pos, i)
	}
    }
    return pos, neg
}

func main() {
    filePath := flag.String("f", "reports", "Report file path")
    taskNum := flag.Int("t", 0, "Task number (0-1)")
    flag.Parse()
    
    file, err := os.Open(*filePath)
    if err != nil {
	log.Fatalf("Error reading report file: %v", err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)

    switch *taskNum {
    case 0:
	count, err := CountSafeReports(*scanner, 0)
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Println(count)
    case 1:
	count, err := CountSafeReports(*scanner, 1)
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Println(count)
    default:
        log.Fatalf("Unknown task number: %d", *taskNum)
    }
}

