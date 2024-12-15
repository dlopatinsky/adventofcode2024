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

type DiskMap []int
type DiskLayout []int

func ReadDiskMap(filePath string) (DiskMap, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("Error reading the disk map file: %v", err)
    }
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanRunes)

    var diskMap DiskMap
    for scanner.Scan() {
        r := scanner.Text()
        n, err := strconv.Atoi(r)
        if err != nil {
            continue
        }
        diskMap = append(diskMap, n)
    }
    return diskMap, nil
}

func (diskMap *DiskMap) DiskLayout() DiskLayout {
    var layout DiskLayout
    id := 0
    for i := 0; i < len(*diskMap); i++ {
        n := (*diskMap)[i]
        for j := 0; j < n; j++ {
            layout = append(layout, id)
        }
        id++
        i++
        if i < len(*diskMap) {
            n = (*diskMap)[i]
            for j := 0; j < n; j++ {
                layout = append(layout, -1)
            }
        }
    }
    return layout
}

func (layout *DiskLayout) Compact() {
    rev := len(*layout) - 1
    for i, n := range *layout {
        if rev <= i {
            break
        }
        if n == -1 {
            for ; rev >= 0 && rev > i; rev-- {
                if (*layout)[rev] != -1 {
                    (*layout)[i], (*layout)[rev] = (*layout)[rev], (*layout)[i]
                    break
                }
            }
        }
    }
}

func (layout *DiskLayout) CompactBlocks() {
    for id := layout.maxId(); id > -1; id-- {
        fileAddress := slices.Index(*layout, id)
        fileSize := layout.blockSize(fileAddress)
        for i := 0; i < len(*layout); {
            freeAddress := slices.Index((*layout)[i:], -1) + i
            if freeAddress == i - 1 {
                break
            }
            freeSize := layout.blockSize(freeAddress)
            if fileAddress > freeAddress &&
                fileSize <= freeSize {
                layout.moveBlock(fileAddress, freeAddress)
                break
            }
            i = freeAddress + freeSize
        }
    }
}

func (layout *DiskLayout) Checksum() int {
    sum := 0
    for i, n := range *layout {
        if n != -1 {
            sum += i * n
        }
    }
    return sum
}

func (layout *DiskLayout) maxId() int {
    max := -1
    for _, v := range *layout {
        if v > max {
            max = v
        }
    }
    return max
}

func (layout *DiskLayout) blockSize(address int) int {
    size := 1
    id := (*layout)[address]
    for i := address + 1; i < len(*layout); i++ {
        if (*layout)[i] == id {
            size++
        } else {
            break
        }
    }
    return size
}

func (layout *DiskLayout) moveBlock(from, to int) {
    size := layout.blockSize(from)
    for i := 0; i < size; i++ {
        (*layout)[from+i], (*layout)[to+i] =
            (*layout)[to+i], (*layout)[from+i]
    }
}

func main() {
    filePath := flag.String("f", "disk_map", "Disk map file path")
    taskNum := flag.Int("t", 0, "Task number")
    flag.Parse()

    diskMap, err := ReadDiskMap(*filePath)
    if err != nil {
        log.Fatal(err)
    }
    layout := diskMap.DiskLayout()

    switch *taskNum {
    case 0:
        layout.Compact()
        fmt.Println(layout.Checksum())
    case 1:
        layout.CompactBlocks()
        fmt.Println(layout.Checksum())
    default:
        log.Fatalf("Unknown task number: %d", *taskNum)
    }
}

