package main

import (
    "encoding/json"
    "fmt"
    "github.com/tealeg/xlsx"
    "io/ioutil"
)

func main() {
    excelFileName := "test.xlsx"
    xlFile, err := xlsx.OpenFile(excelFileName)
    var stops []ShuttleStop
    for _, sheet := range xlFile.Sheets {
        for _, row := range sheet.Rows {
            for _, cell := range row.Cells {
                stops[cell] := cell.String()
            }
        }
    }

    // export as json
    stopsJson, err := json.Marshal(stops)
    if err != nil {
        println("error:", err)
    }
    err = ioutil.WriteFile("schedule.json", stopsJson, 0644)
}
