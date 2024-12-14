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

type (
    Equation struct {
        Value    int
        Operands []int
    }
)

func EquationSum(scanner bufio.Scanner, concat bool) (int, error) {
    sum := 0
    for scanner.Scan() {
        line := scanner.Text()
        slice := strings.SplitAfter(line, ":")

        v := slice[0][:len(slice[0])-1]
        value, err := strconv.Atoi(v)
        if err != nil {
            return sum, fmt.Errorf("Invalid test value: %v", v)
        }
        slice[1] = strings.Trim(slice[1], " ")
        var operands []int
        for _, n := range strings.Split(slice[1], " ") {
            o, err := strconv.Atoi(n)
            if err != nil {
                return sum, fmt.Errorf("Invalid operand: %v", n)
            }
            operands = append(operands, o)
        }

        equation := Equation{value, operands}
        if equation.IsValid(concat) {
            sum += equation.Value
        }
    }
    if err := scanner.Err(); err != nil {
        return sum, err
    }
    return sum, nil
}

func (equation *Equation) IsValid(concat bool) bool {
    return equation.evaluate(1, equation.Operands[0], concat)
}

func (equation *Equation) evaluate(current, value int, concat bool) bool {
    if current == len(equation.Operands) {
        return equation.Value == value
    }

    op := equation.Operands[current]
    current++

    result := equation.evaluate(current, value+op, concat) ||
        equation.evaluate(current, value*op, concat)
    if concat {
        result = result || equation.evaluate(current, Concat(value, op), concat)
    }
    return result
}

func Concat(a, b int) int {
    c := 10
    for b%c != b {
        c *= 10
    }
    return a*c + b
}

func main() {
    filePath := flag.String("f", "equations", "Equation file path")
    taskNum := flag.Int("t", 0, "Task number (0-1)")
    flag.Parse()

    file, err := os.Open(*filePath)
    if err != nil {
        log.Fatalf("Error reading equation file: %v", err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)

    switch *taskNum {
    case 0:
        sum, err := EquationSum(*scanner, false)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(sum)
    case 1:
        sum, err := EquationSum(*scanner, true)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(sum)
    default:
        log.Fatalf("Unknown task number: %d", *taskNum)
    }
}

