package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "slices"
)

type Direction int
const (
    UP  = iota
    LEFT
    DOWN
    RIGHT
)

type (
    Position struct {
        X int;
        Y int;
    }
    LabMap struct {
        StartPos Position;
        Obstacles []Position;
        Width int;
        Height int;
    }
) 

func ReadLabMap(filePath string) (LabMap, error) {
    var labMap LabMap
    file, err := os.Open(filePath)
    if err != nil {
        return labMap, fmt.Errorf("Error reading the map file: %v", err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)

    y := 0
    for scanner.Scan() {
        line := scanner.Text()
        for x, p := range line {
            switch p {
            case '.':
            case '#':
                labMap.Obstacles = append(labMap.Obstacles,
                    Position {x, y})
            case '^':
                labMap.StartPos = Position {x, y}
            default:
                return labMap, fmt.Errorf("Unknown map symbol: %v", p)
            }
        }
        labMap.Width = len(line);
        y++
    }
    labMap.Height = y;
    return labMap, nil
}

func (labMap *LabMap) VisitedPositions() []Position {
    var visited []Position
    pos := labMap.StartPos
    direction := UP
    guardLeft := false
    for !guardLeft {
        switch direction {
        case UP:
            if pos.Y == 0 {
                guardLeft = true
            }
            nextPos := Position {pos.X, pos.Y - 1}
            if slices.Contains(labMap.Obstacles, nextPos) {
                direction = LEFT
            } else {
                pos = nextPos
            }
        case LEFT:
            if pos.X == labMap.Width - 1 {
                guardLeft = true
            }
            nextPos := Position {pos.X + 1, pos.Y}
            if slices.Contains(labMap.Obstacles, nextPos) {
                direction = DOWN
            } else {
                pos = nextPos
            }
        case DOWN:
            if pos.Y == labMap.Height - 1{
                guardLeft = true
            }
            nextPos := Position {pos.X, pos.Y + 1}
            if slices.Contains(labMap.Obstacles, nextPos) {
                direction = RIGHT
            } else {
                pos = nextPos
            }
        case RIGHT:
            if pos.X == 0 {
                guardLeft = true
            }
            nextPos := Position {pos.X - 1, pos.Y}
            if slices.Contains(labMap.Obstacles, nextPos) {
                direction = UP
            } else {
                pos = nextPos
            }
        }
        if !slices.Contains(visited, pos) {
            visited = append(visited, pos)
        }
    }
    return visited
}

type Pair struct {
    pos Position;
    dir Direction;
}

func (labMap *LabMap) Loops(obstacle Position) bool {
    var visited []Pair
    pos := labMap.StartPos
    direction := UP
    guardLeft := false
    obstacles := append(labMap.Obstacles, obstacle)
    for !guardLeft {
        switch direction {
        case UP:
            if pos.Y == 0 {
                guardLeft = true
            }
            nextPos := Position {pos.X, pos.Y - 1}
            if slices.Contains(obstacles, nextPos) {
                direction = LEFT
            } else {
                pos = nextPos
            }
        case LEFT:
            if pos.X == labMap.Width - 1 {
                guardLeft = true
            }
            nextPos := Position {pos.X + 1, pos.Y}
            if slices.Contains(obstacles, nextPos) {
                direction = DOWN
            } else {
                pos = nextPos
            }
        case DOWN:
            if pos.Y == labMap.Height - 1{
                guardLeft = true
            }
            nextPos := Position {pos.X, pos.Y + 1}
            if slices.Contains(obstacles, nextPos) {
                direction = RIGHT
            } else {
                pos = nextPos
            }
        case RIGHT:
            if pos.X == 0 {
                guardLeft = true
            }
            nextPos := Position {pos.X - 1, pos.Y}
            if slices.Contains(obstacles, nextPos) {
                direction = UP
            } else {
                pos = nextPos
            }
        }
        pair := Pair {pos, Direction(direction)};
        if !slices.Contains(visited, pair) {
            visited = append(visited, pair)
        } else {
            return true
        }
    }
    return false
}

func main() {
    filePath := flag.String("f", "map", "Map file path")
    taskNum := flag.Int("t", 0, "Task number")
    flag.Parse()

    labMap, err := ReadLabMap(*filePath)
    if err != nil {
        log.Fatal(err)
    }

    switch *taskNum {
    case 0:
        fmt.Println(len(labMap.VisitedPositions()))
    case 1:
        visited := labMap.VisitedPositions()
        count := 0
        for _, pos := range visited {
            if labMap.Loops(pos) {
                count++
            }
        }
        fmt.Println(count)
    default:
        log.Fatalf("Unknown task number: %d", *taskNum)
    }
}

