package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "regexp"
    "strconv"
)

func GetMulSum(scanner bufio.Scanner, ignoreDo bool) int {
    sum := 0
    var r *regexp.Regexp
    if ignoreDo {
        r = regexp.MustCompile(`(mul)\((\d{1,3}),(\d{1,3})\)`)
    } else {
        r = regexp.MustCompile(`(?:((mul)\((\d{1,3}),(\d{1,3})\))|((do)\(\))|((don't)\(\)))`)
    }
    doing := true
    for scanner.Scan() {
        line := scanner.Text()
        matches := r.FindAllStringSubmatch(line, -1)
        for _, match := range matches {
            if ignoreDo {
                r, err := strconv.Atoi(match[2])
                if err != nil {
                    continue
                }
                l, err := strconv.Atoi(match[3])
                if err != nil {
                    continue
                }
                sum += r * l
            } else {
                if match[6] == "do" {
                    doing = true
                } else if match[8] == "don't" {
                    doing = false
                } else if doing && match[2] == "mul" {
                    r, err := strconv.Atoi(match[3])
                    if err != nil {
                        continue
                    }
                    l, err := strconv.Atoi(match[4])
                    if err != nil {
                        continue
                    }
                    sum += r * l
                }
            }
        }
    }
    return sum
}

func main() {
    filePath := flag.String("f", "memory", "Memory file path")
    taskNum := flag.Int("t", 0, "Task number (0-1)")
    flag.Parse()

    file, err := os.Open(*filePath)
    if err != nil {
        log.Fatalf("Failed to read the memory file: %v", err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)

    switch *taskNum {
    case 0:
        sum := GetMulSum(*scanner, true)
        fmt.Println(sum)
    case 1:
        sum := GetMulSum(*scanner, false)
        fmt.Println(sum)
    default:
        log.Fatalf("Unknown task number: %d", *taskNum)
    }
}

