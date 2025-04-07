# Golang structs and parser for PTV Visum network file

W.I.P.

## How to use:
* Get the package:
    ```shell
    go get github.com/lddl/go-ptv-visum
    ```

* Example (full example is [here](./example/sample/main.go)):
    ```go
    package main

    import (
        "fmt"
        "os"
        "sort"
        "strings"

        ptvvisum "github.com/lddl/go-ptv-visum"
    )

    func main() {
        file, err := os.Open("./example/sample/example.net")
        if err != nil {
            fmt.Println(err)
            return
        }
        defer file.Close()
        ptvData, err := ptvvisum.ReadPTVFromFile(file)
        if err != nil {
            fmt.Println(err)
            return
        }

        fmt.Println("Version:")
        fmt.Printf("\tVersion: %s\n", ptvData.Version.Version)
        fmt.Printf("\tFileType: %s\n", ptvData.Version.FileType)
        fmt.Printf("\tLanguage: %s\n", ptvData.Version.Language)
        fmt.Printf("\tUnit: %s\n", ptvData.Version.Unit)

    }
    ```

* Those sections ARE NOT supported currently:
    * Table: Stops
    * Table: Stop areas
    * Table: Stop points
    * Table: Lines
    * Table: Line routes
    * Table: Line route items
    * Table: Time profiles
    * Table: Time profile items
    * Table: Vehicle journeys
    * Table: Vehicle journey sections
    * Table: Transfer walk times between stop areas
    * Table: Block versions
    * Table: Points of interest: State (32)
    * Table: Points of interest: County (33)
    * Table: Points of interest: Municipality (34)
    * Table: Legs
    * Table: Lanes
    * Table: Lane turns
    * Table: Crosswalks
