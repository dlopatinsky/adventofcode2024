package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
)

type (
    Position struct {
        X int
        Y int
    }
    Garden map[rune][]Position
    Region map[Position]struct{}
)

var directions = []Position{
    {-1, 0},
    {1, 0},
    {0, -1},
    {0, 1},
}

var sideDirections = []Position{
    {-1, -1},
    {-1, 1},
    {1, -1},
    {1, 1},
}

func ReadGardenMap(filePath string) (Garden, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    scanner := bufio.NewScanner(file)

    garden := make(Garden)
    y := 0
    for scanner.Scan() {
        line := scanner.Text()
        for x, f := range line {
            garden[f] = append(garden[f], Position{x, y})
        }
        y++
    }
    return garden, nil
}

func (garden *Garden) Regions() []Region {
    var regions []Region
    for k, v := range *garden {
        visited := make(map[Position]struct{})
        for _, p := range v {
            if _, ok := visited[p]; !ok {
                region := garden.dfs(k, p, visited)
                regions = append(regions, region)
            }
        }
    }
    return regions
}

func (garden *Garden) dfs(fruit rune, pos Position, visited map[Position]struct{}) Region {
    region := Region{pos: struct{}{}}
    positions := (*garden)[fruit]
    for _, dir := range directions {
        nextPos := Position{pos.X + dir.X, pos.Y + dir.Y}
        for _, p := range positions {
            if _, ok := visited[p]; !ok && nextPos == p {
                visited[p] = struct{}{}
                for k, v := range garden.dfs(fruit, nextPos, visited) {
                    region[k] = v
                }
            }
        }
    }
    return region
}

func (region *Region) Price() int {
    area := region.area()
    perimeter := region.perimeter()
    return area * perimeter 
}

func (region *Region) BulkPrice() int {
    area := region.area()
    sides := region.sides()
    return area * sides
}

func (region *Region) area() int {
    return len(*region)
}

func (region *Region) perimeter() int {
    perim := 0
    for pos := range *region {
        for _, dir := range directions {
            if _, ok := (*region)[Position{pos.X + dir.X, pos.Y + dir.Y}]; !ok {
                perim++
            }
        }
    }
    return perim
}

func(region *Region) sides() int {
    sides := 0
    for p := range *region {
        for _, dir := range sideDirections {
            top := Position{p.X, p.Y + dir.Y}
            left := Position{p.X + dir.X, p.Y}
            _, hasTop := (*region)[top]
            _, hasLeft := (*region)[left]
            // convex corner
            if !hasTop && !hasLeft {
                sides++
                continue
            }
            topLeft := Position{p.X + dir.X, p.Y + dir.Y}
            _, hasTopLeft := (*region)[topLeft]
            // concave corner
            if hasTop && hasLeft && !hasTopLeft {
                sides++
                continue
            }
        }
    }
    return sides
}

func main() {
    filePath := flag.String("f", "garden", "Garden map file path")
    taskNum := flag.Int("t", 0, "Task number")
    flag.Parse()

    garden, err := ReadGardenMap(*filePath)
    if err != nil {
        log.Fatalf("Error reading the garden map file: %v", err)
    }
    regions := garden.Regions()

    switch *taskNum {
    case 0:
        total := 0
        for _, r := range regions {
            total += r.Price()
        }
        fmt.Println(total)
    case 1:
        total := 0
        for _, r := range regions {
            total += r.BulkPrice()
        }
        fmt.Println(total)
    default:
        log.Fatalf("Invalid task number: %v", *taskNum)
    }
}

