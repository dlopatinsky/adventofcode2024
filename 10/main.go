package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "slices"
    "strconv"
)

type (
    HeightMap [][]int
    Position struct {
        X int
        Y int
    }
)

func ReadHeightMap(filePath string) (HeightMap, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    scanner := bufio.NewScanner(file)

    var heightMap HeightMap
    for scanner.Scan() {
        line := scanner.Text()
        row := make([]int, len(line))
        for x, h := range line {
            row[x], err = strconv.Atoi(string(h))
            if err != nil {
                return nil, fmt.Errorf("Invalid height value: %v", h)
            }
        }
        heightMap = append(heightMap, row)
    }
    return heightMap, nil
}

func (heightMap *HeightMap) Trailheads() []Position {
    var positions []Position
    for y, row := range *heightMap {
        for x, h := range row {
            if h == 0 {
                positions = append(positions, Position{x, y})
            }
        }
    }
    return positions
}

func (heightMap *HeightMap) ReachableNines(pos Position) []Position {
    current := (*heightMap)[pos.Y][pos.X]
    if current == 9 {
        return []Position{pos}
    }
    var reachable []Position 
    for _, dx := range []int{1, -1} {
        nextX := pos.X + dx
        if nextX < len((*heightMap)[pos.Y]) && nextX >= 0 &&
        (*heightMap)[pos.Y][nextX] == current + 1 {
            for _, p := range heightMap.ReachableNines(Position{nextX, pos.Y}) {
                if !slices.Contains(reachable, p) {
                    reachable = append(reachable, p)
                }
            }
        }
    }
    for _, dy := range []int{1, -1} {
        nextY := pos.Y + dy
        if nextY < len(*heightMap) && nextY >= 0 &&
            (*heightMap)[nextY][pos.X] == current + 1 {
            for _, p := range heightMap.ReachableNines(Position{pos.X, nextY}) {
                if !slices.Contains(reachable, p) {
                    reachable = append(reachable, p)
                }
            }
        }
    }
    return reachable
}

func (heightMap *HeightMap) Rating(pos Position) int {
    current := (*heightMap)[pos.Y][pos.X]
    if current == 9 {
        return 1
    }
    rating := 0
    for _, dx := range []int{1, -1} {
        nextX := pos.X + dx
        if nextX < len((*heightMap)[pos.Y]) && nextX >= 0 &&
        (*heightMap)[pos.Y][nextX] == current + 1 {
            rating += heightMap.Rating(Position{nextX, pos.Y})
        }
    }
    for _, dy := range []int{1, -1} {
        nextY := pos.Y + dy
        if nextY < len(*heightMap) && nextY >= 0 &&
            (*heightMap)[nextY][pos.X] == current + 1 {
            rating += heightMap.Rating(Position{pos.X, nextY})
        }
    }
    return rating
}

func main() {
    filePath := flag.String("f", "height_map", "Height map file path")
    taskNum := flag.Int("t", 0, "Task number")
    flag.Parse()

    heightMap, err := ReadHeightMap(*filePath)
    if err != nil {
        log.Fatalf("Error reading the height map file: %v", err)
    }

    switch *taskNum {
    case 0:
        trailheads := heightMap.Trailheads()
        score := 0
        for _, t := range trailheads {
            score += len(heightMap.ReachableNines(t))
        }
        fmt.Println(score)
    case 1:
        trailheads := heightMap.Trailheads()
        ratings := 0
        for _, t := range trailheads {
            ratings += heightMap.Rating(t)
        }
        fmt.Println(ratings)
    default:
        log.Fatalf("Invalid task number: %d", *taskNum)
    }
}

