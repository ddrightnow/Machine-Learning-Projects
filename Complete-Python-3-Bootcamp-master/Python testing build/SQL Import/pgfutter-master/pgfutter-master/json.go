package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Try to JSON decode the bytes
func tryUnmarshal(b []byte) error {
	var v interface{}
	err := json.Unmarshal(b, &v)
	return err
}

//Copy JSON Rows and return list of errors
func copyJSONRows(i *Import, reader *bufio.Reader, ignoreErrors bool) (error, int, int) {
	success := 0
	failed := 0

	for {
		// ReadBytes instead of a Scanner because it can deal with very long lines
		// which happens often with big JSON objects
		line, err := reader.ReadBytes('\n')

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			err = fmt.Errorf("%s: %s", err, line)
			return err, success, failed
		}

		err = tryUnmarshal(line)
		if err != nil {
			failed++
			if ignoreErrors {
				os.Stderr.WriteString(string(line))
				continue
			} else {
				err = fmt.Errorf("%s: %s", err, line)
				return err, success, failed
			}
		}

		err = i.AddRow(string(line))
		if err != nil {
			failed++
			if ignoreErrors {
				os.Stderr.WriteString(string(line))
				continue
			} else {
				err = fmt.Errorf("%s: %s", err, line)
				return err, success, failed
			}
		}

		success++
	}

	return nil, success, failed
}

func importJSON(filename string, connStr string, schema string, tableName string, ignoreErrors bool, dataType string) error {

	db, err := connect(connStr, schema)
	if err != nil {
		return err
	}
	defer db.Close()

	i, err := NewJSONImport(db, schema, tableName, "data", dataType)
	if err != nil {
		return err
	}

	var success, failed int
	if filename == "" {
		reader := bufio.NewReader(os.Stdin)
		err, success, failed = copyJSONRows(i, reader, ignoreErrors)
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		bar := NewProgressBar(file)
		reader := bufio.NewReader(io.TeeReader(file, bar))
		bar.Start()
		err, success, failed = copyJSONRows(i, reader, ignoreErrors)
		bar.Finish()
	}

	if err != nil {
		lineNumber := success + failed
		return fmt.Errorf("line %d: %s", lineNumber, err)
	}

	fmt.Println(fmt.Sprintf("%d rows imported into %s.%s", success, schema, tableName))

	if ignoreErrors && failed > 0 {
		fmt.Println(fmt.Sprintf("%d rows could not be imported into %s.%s and have been written to stderr.", failed, schema, tableName))
	}

	return i.Commit()
}
