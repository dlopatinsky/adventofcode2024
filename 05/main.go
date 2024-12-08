package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "slices"
    "strconv"
    "strings"
)

type (
    Rules map[int][]int
    Update []int
) 

func ReadRulesAndUpdates(filePath string) (Rules, []Update, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, nil, fmt.Errorf("Error reading the pages file: %v", err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)

    rules := make(Rules)
    var updates []Update
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "|") {
            err := AddRule(rules, line)
            if err != nil {
                return nil, nil, fmt.Errorf("Error reading a rule: %v", err)
            }
        } else if strings.Contains(line, ",") {
            err := AddUpdate(&updates, line)
            if err != nil {
                return nil, nil, fmt.Errorf("Error reading an update: %v", err)
            }
        }
    }
    return rules, updates, nil
}

func AddRule(rules Rules, line string) error {
    slice := strings.Split(line, "|")
    a, err := strconv.Atoi(slice[0])
    if err != nil {
        return err
    }
    b, err := strconv.Atoi(slice[1])
    if err != nil {
        return err
    }
    rules[b] = append(rules[b], a)
    return nil
}

func AddUpdate(updates *[]Update, line string) error {
    slice := strings.Split(line, ",")
    var newUpdate Update
    for _, s := range slice {
        page, err := strconv.Atoi(s)
        if err != nil {
            return err
        }
        newUpdate = append(newUpdate, page)
    }
    *updates = append(*updates, newUpdate)
    return nil
}

func FilterCorrectUpdates(updates []Update, rules Rules, keepCorrect bool) []Update {
    var filtered []Update
    for _, u := range updates {
        if u.IsCorrect(rules) == keepCorrect {
            filtered = append(filtered, u)
        }
    }
    return filtered
}

func (update Update) IsCorrect(rules Rules) bool {
    for i := 0; i < len(update) - 1 ; i++ {
        a, b := update[i], update[i + 1]
        if r, ok := rules[b]; !ok || slices.Index(r, a) == -1 {
            return false
        }
    }
    return true
}

func SumMiddlePages(updates []Update) int {
    sum := 0
    for _, u := range updates {
        sum += u[len(u) / 2]
    }
    return sum
}

func (update *Update) Fix(rules Rules) {
    slices.SortFunc(*update, func(a, b int) int {
        if _, ok := rules[b]; a == b || !ok {
            return 0
        }
        if slices.Contains(rules[b], a) {
            return -1
        }
        return 1
    })
}

func main() {
    filePath := flag.String("f", "pages", "Pages file path")
    taskNum := flag.Int("t", 0, "Task number")
    flag.Parse()

    rules, updates, err := ReadRulesAndUpdates(*filePath)
    if err != nil {
        log.Fatal(err)
    }

    switch *taskNum {
    case 0:
        correctUpdates := FilterCorrectUpdates(updates, rules, true)
        fmt.Println(SumMiddlePages(correctUpdates))
    case 1:
        incorrectUpdates := FilterCorrectUpdates(updates, rules, false)
        for _, u := range incorrectUpdates {
            u.Fix(rules)
        }
        fmt.Println(SumMiddlePages(incorrectUpdates))
    default:
        log.Fatalf("Unknown task number: %d", *taskNum)
    }
}

