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

func ReadReports (filePath string) ([]Report, error) {
    file, err := os.Open(filePath)
    if err != nil {
	return nil, err
    }
    defer file.Close()

    var reports []Report
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
	line := scanner.Text()
	slice := strings.Fields(line)
	report := make(Report, len(slice))
	for i, v := range slice {
	    n, err := strconv.Atoi(v)
	    if err != nil {
		return nil, fmt.Errorf("Invalid number: %s", v)
	    }
	    report[i] = n
	}
	reports = append(reports, report)
    }
    return reports, nil
}

func CountSafeReports(reports []Report) int {
    count := 0
    for _, r := range reports {
	if r.UnsafeLevelCount() == 0 {
	    count++
	}
    }
    return count
}

func CountSafeReportsWithProblemDampener(reports []Report) int {
    count := 0
    for _, r := range reports {
	if r.UnsafeLevelCount() < 2 {
	    count++
	}
    }
    return count
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
    
    reports, err := ReadReports(*filePath)
    if err != nil {
        log.Fatalf("Error reading reports: %v", err)
    }

    switch *taskNum {
    case 0:
	count := CountSafeReports(reports)
	fmt.Println(count)
    case 1:
	count := CountSafeReportsWithProblemDampener(reports)
	fmt.Println(count)
    default:
        log.Fatalf("Unknown task number: %d", *taskNum)
    }
}

