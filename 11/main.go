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
    Stone int
    Stones map[Stone]int
)

func ReadStones(filePath string) (Stones, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    scanner := bufio.NewScanner(file)

    stones := make(Stones)
    for scanner.Scan() {
        line := scanner.Text()
        slice := strings.Split(line, " ")
        for _, s := range slice {
            num, err := strconv.Atoi(s)
            if err != nil {
                return nil, fmt.Errorf("Invalid number on a stone: %v, %v", s, err)
            }
            stones[Stone(num)]++
        }
    }
    return stones, nil
}

func (stones *Stones) Blink() {
    newStones := make(Stones)
    for stone, count := range *stones {
        if stone == 0 {
            newStones[1] += count
        } else {
            digits := stone.Digits()
            length := len(digits)
            if length % 2 == 0 {
                newStones[ToStone(digits[:length / 2])] += count
                newStones[ToStone(digits[length / 2:])] += count
            } else {
                newStones[stone * 2024] += count
            }
        }
        *stones = newStones
    }
}

func (stone *Stone) Digits() []int {
    if *stone == 0 {
        return []int{0}
    }
    number := int(*stone)
    var digits []int
    for number > 0 {
        digits = append(digits, number % 10)
        number /= 10
    }
    slices.Reverse(digits)
    return digits
}

func ToStone(digits []int) Stone {
    stone := 0
    a := 1
    for i := len(digits) - 1; i > -1; i-- {
        stone += digits[i] * a
        a *= 10
    }
    return Stone(stone)
}

func (stones *Stones) Count() int {
    count := 0
    for _, c := range *stones {
        count += c
    }
    return count
}

func main() {
    filePath := flag.String("f", "stones", "Stone file path")
    taskNum := flag.Int("t", 0, "Task number")
    flag.Parse()

    stones, err := ReadStones(*filePath)
    if err != nil {
        log.Fatalf("Error reading the stone file: %v", err)
    }

    switch *taskNum {
    case 0:
        for i := 0; i < 25; i++ {
            stones.Blink()
        }
        fmt.Println(stones.Count())
    case 1:
        for i := 0; i < 75; i++ {
            stones.Blink()
        }
        fmt.Println(stones.Count())
    default:
        log.Fatalf("Invalid task number: %v", *taskNum)
    }
}

