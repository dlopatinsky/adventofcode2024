package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "sort"
    "strconv"
    "strings"
)

func ReadLocationIds(filePath string) ([]int, []int, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, nil, err
    }
    defer file.Close()

    var rightList []int
    var leftList []int

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        slice := strings.Fields(line)
        if len(slice) < 2 {
            continue
        }
        right, err := strconv.Atoi(slice[0])
        if err != nil {
            return nil, nil, fmt.Errorf("Invalid number in the right list: %s", slice[0])
        }
        rightList = append(rightList, right)
        
        left, err := strconv.Atoi(slice[1])
        if err != nil {
            return nil, nil, fmt.Errorf("Invalid number in the left list: %s", slice[1])
        }
        leftList = append(leftList, left)
    }
    if err := scanner.Err(); err != nil {
        return nil, nil, err
    }

    return rightList, leftList, nil
}

func CalculateDistance(rightList, leftList []int) int {
    SortSlice(rightList)
    SortSlice(leftList)

    distance := 0
    for i, r := range rightList {
        l := leftList[i]
        distance += Abs(r - l)
    }

    return distance
}

func CalculateSimilarity(rightList, leftList []int) int {
    m := make(map[int]int)

    similarity := 0
    for _, r := range rightList {
        similarity += GetSimilarity(r, leftList, m) * r
    }

    return similarity
}

func GetSimilarity(rightNum int, leftList []int, m map[int]int) int {
    if val, ok := m[rightNum]; ok {
        return val
    } else {
        m[rightNum] = CountOccurrences(rightNum, leftList)
        return m[rightNum]
    }
}

func SortSlice(s []int) {
    sort.Slice(s, func(i, j int) bool {
        return s[i] < s[j]
    })
}

func Abs(i int) int {
    if i < 0 {
        return -i
    }
    return i
}

func CountOccurrences(v int, s []int) int {
    count := 0
    for _, e := range s {
        if e == v {
            count++
        }
    }
    return count
}

func main() {
    filePath := flag.String("f", "location_ids", "Location ID file path")
    taskNum := flag.Int("t", 0, "Task number (0-1)")
    flag.Parse()

    rightList, leftList, err := ReadLocationIds(*filePath)
    if err != nil {
        log.Fatalf("Error reading location IDs: %v", err)
    }

    switch *taskNum {
    case 0:
        distance := CalculateDistance(rightList, leftList)
        fmt.Println(distance)
    case 1:
        similarity := CalculateSimilarity(rightList, leftList)
        fmt.Println(similarity)
    default:
        log.Fatalf("unknown task number: %d", *taskNum)
    }
}

