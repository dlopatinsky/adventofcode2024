package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
)

type Direction int
const (
    RIGHT Direction = iota
    RIGHT_DOWN
    DOWN
    LEFT_DOWN
    LEFT
    LEFT_TOP
    TOP
    RIGHT_TOP
)

func ReadWordsAsRunes(filePath string) ([][]rune, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("Failed to read the word file: %v", err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    var runes [][]rune
    for scanner.Scan() {
        line := scanner.Text()
        runes = append(runes, []rune(line))
    }
    return runes, nil
}

func WordCount(runes [][]rune, word []rune) int {
    if len(word) == 0 {
        return 0
    }
    count := 0
    for y, line := range runes {
        firstLetters := FindRuneIndices(line, word[0])
        for _, x := range firstLetters {
            count += FullWordCount(runes, word, x, y)
        }
    }
    return count
}

func XWordCount(runes [][]rune) int {
    count := 0
    for y, line := range runes {
        lettersA := FindRuneIndices(line, 'A')
        for _, x := range lettersA {
            if IsAnX(runes, x, y) {
                count++
            }
        }
    }
    return count
}

func FindRuneIndices(slice []rune, target rune) []int {
    var indices []int
    for i, r := range slice {
        if r == target {
            indices = append(indices, i)
        }
    }
    return indices
}

func FullWordCount(runes [][]rune, word []rune, x, y int) int {
    count := 0
    for d := 0; d < 8; d++ {
        if CheckNextRune(runes, word, 1, x, y, Direction(d)) {
            count++
        }
    }
    return count
}

func CheckNextRune(runes [][]rune, word []rune, runeIndex, x, y int, d Direction) bool {
    if len(word) == runeIndex {
        return true
    }
    switch d {
    case RIGHT:
        if x == len(runes[y]) - 1 || runes[y][x + 1] != word[runeIndex] {
            return false
        }
        return CheckNextRune(runes, word, runeIndex + 1, x + 1, y, d)
    case RIGHT_DOWN:
        if x == len(runes[y]) - 1 || y == len(runes) - 1 || runes[y + 1][x + 1] != word[runeIndex] {
            return false
        }
        return CheckNextRune(runes, word, runeIndex + 1, x + 1, y + 1, d)
    case DOWN:
        if y == len(runes) - 1 || runes[y + 1][x] != word[runeIndex] {
            return false
        }
        return CheckNextRune(runes, word, runeIndex + 1, x, y + 1, d)
    case LEFT_DOWN:
        if x == 0 || y == len(runes) - 1 || runes[y + 1][x - 1] != word[runeIndex] {
            return false
        }
        return CheckNextRune(runes, word, runeIndex + 1, x - 1, y + 1, d)
    case LEFT:
        if x == 0 || runes[y][x - 1] != word[runeIndex] {
            return false
        }
        return CheckNextRune(runes, word, runeIndex + 1, x - 1, y, d)
    case LEFT_TOP:
        if x == 0 || y == 0 || runes[y - 1][x - 1] != word[runeIndex] {
            return false
        }
        return CheckNextRune(runes, word, runeIndex + 1, x - 1, y - 1, d)
    case TOP:
        if y == 0 || runes[y - 1][x] != word[runeIndex] {
            return false
        }
        return CheckNextRune(runes, word, runeIndex + 1, x, y - 1, d)
    case RIGHT_TOP:
        if x == len(runes[y]) - 1 || y == 0 || runes[y - 1][x + 1] != word[runeIndex] {
            return false
        }
        return CheckNextRune(runes, word, runeIndex + 1, x + 1, y - 1, d)
    default:
        return false
    }
}

func IsAnX(runes [][]rune, x, y int) bool {
    return x > 0 && y > 0 &&
        x < len(runes[y]) - 1 &&
        y < len(runes) - 1 &&
        ((runes[y - 1][x - 1] == 'M' &&
        runes[y + 1][x + 1] == 'S') || 
        (runes[y - 1][x - 1] == 'S' &&
        runes[y + 1][x + 1] == 'M')) &&
        ((runes[y - 1][x + 1] == 'M' &&
        runes[y + 1][x - 1] == 'S') || 
        (runes[y - 1][x + 1] == 'S' &&
        runes[y + 1][x - 1] == 'M'))
}

func main() {
    filePath := flag.String("f", "words", "Word file path")
    taskNum := flag.Int("t", 0, "Task number")
    flag.Parse()

    runes, err := ReadWordsAsRunes(*filePath)
    if err != nil {
        log.Fatal(err)
    }

    switch *taskNum {
    case 0:
        count := WordCount(runes, []rune{'X', 'M', 'A', 'S'})
        fmt.Println(count)
    case 1:
        count := XWordCount(runes)
        fmt.Println(count)
    default:
        log.Fatalf("Unknown task number: %d", *taskNum)
    }
}

