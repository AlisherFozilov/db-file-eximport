package db_file_eximport

import (
	"database/sql"
	"io/ioutil"
)

//--------NO WARRANTIES!

type MapperRowTo func(rows *sql.Rows) (interface{}, error)
type MapperInterfaceSliceTo func([]interface{}) interface{}
type Marshaller func(interface{}) ([]byte, error)

func ExportToFile(
	db *sql.DB,
	getDataFromDbSQL string,
	filename string,
	mapRow MapperRowTo,
	marshal Marshaller,
	mapDataSlice MapperInterfaceSliceTo) error {

	rows, err := db.Query(getDataFromDbSQL)
	if err != nil {
		return err
	}
	defer func() {
		err = rows.Close()
	}()
	var dataSlice []interface{}
	for rows.Next() {
		dataElement, err := mapRow(rows)
		if err != nil {
			return err
		}
		dataSlice = append(dataSlice, dataElement)
	}
	exportData := mapDataSlice(dataSlice)
	data, err := marshal(exportData)
	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		return err
	}
	return nil
}


type MapperBytesTo func([]byte) ([]interface{}, error)

func ImportFromFile(
	db *sql.DB,
	filename string,
	mapBytes MapperBytesTo,
	insertToDB func(interface{}, *sql.DB) error,
) error {
	itemsData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	sliceData, err := mapBytes(itemsData)

	for _, datum := range sliceData {
		err = insertToDB(datum, db)
		if err != nil {
			return err
		}
	}

	return nil
}