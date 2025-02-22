package main

import (
	"errors"
	"gen/factors"
	"os"
)

func main() {
	//fmt.Println("choose gen function")

	dir, err := os.OpenRoot("../arit/bin/")

	fileName := "./factors_file"
	tableName := "./factors_table"

	file, err := dir.Create(fileName)
	table, err2 := dir.Create(tableName)
	if err != nil || err2 != nil {
		panic(errors.Join(err, err2))
	}

	factors.Make(file, table)
}
