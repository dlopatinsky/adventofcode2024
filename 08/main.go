package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "slices"
)

type (
    Position struct {
        x int
        y int
    }
    AntennaMap struct {
        antennas map[rune][]Position
        width    int
        height   int
    }
)

func ReadAntennaMap(filePath string) (AntennaMap, error) {
    file, err := os.Open(filePath)
    var antennaMap AntennaMap
    if err != nil {
        return antennaMap, err
    }
    scanner := bufio.NewScanner(file)

    antennaMap.antennas = make(map[rune][]Position)
    y := 0
    for scanner.Scan() {
        line := scanner.Text()
        for x, f := range line {
            if f != '.' {
                antennaMap.antennas[f] = append(antennaMap.antennas[f], Position{x, y})
            }
        }
        antennaMap.width = len(line)
        y++
        antennaMap.height = y
    }

    return antennaMap, nil
}

func (antennaMap *AntennaMap) AntinodeCount(resonantHarmonics bool) int {
    var antinodes []Position
    for _, positions := range antennaMap.antennas {
        for _, a := range positions {
            for _, b := range positions {
                if a != b {
                    dx := a.x - b.x
                    dy := a.y - b.y
                    var newAntinodes []Position
                    p := a
                    for p.x >= 0 && p.x < antennaMap.width &&
                        p.y >= 0 && p.y < antennaMap.height {
                        p.x += dx
                        p.y += dy
                        newAntinodes = append(newAntinodes, p)
                        if !resonantHarmonics {
                            break
                        }
                    }
                    p = a
                    for p.x >= 0 && p.x < antennaMap.width &&
                        p.y >= 0 && p.y < antennaMap.height {
                        p.x -= dx
                        p.y -= dy
                        newAntinodes = append(newAntinodes, p)
                        if !resonantHarmonics {
                            break
                        }
                    }
                    for _, p := range newAntinodes {
                        if p.x >= 0 && p.x < antennaMap.width &&
                            p.y >= 0 && p.y < antennaMap.height &&
                            !slices.Contains(antinodes, p) {
                            if !resonantHarmonics && slices.Contains(positions, p) {
                                continue
                            }
                            antinodes = append(antinodes, p)
                        }
                    }
                }
            }
        }
    }
    return len(antinodes)
}

func main() {
    filePath := flag.String("f", "map", "Map file path")
    taskNum := flag.Int("t", 0, "Task number")
    flag.Parse()

    antennaMap, err := ReadAntennaMap(*filePath)
    if err != nil {
        log.Fatal(err)
    }

    switch *taskNum {
    case 0:
        count := antennaMap.AntinodeCount(false)
        fmt.Println(count)
    case 1:
        count := antennaMap.AntinodeCount(true)
        fmt.Println(count)
    default:
        log.Fatalf("Unknown task number: %d", *taskNum)
    }
}

