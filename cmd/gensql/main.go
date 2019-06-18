package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const sqlFormat = "INSERT INTO `%s`.`%s` (%s) VALUES (%s);"

func main() {
	var (
		flagInput   = flag.String("i", "", "csv format file")
		flagOuput   = flag.String("o", "", "output to file or stdout, default, os.Stdout")
		flagDBName  = flag.String("db", "", "db name to write sql")
		flagTblName = flag.String("tbl", "", "table name to write sql")
	)

	flag.Parse()

	c := &config{
		input:   *flagInput,
		output:  *flagOuput,
		dbName:  *flagDBName,
		tblName: *flagTblName,
	}

	if err := c.parse(); err != nil {
		log.Printf("[Error]: c.parse() got err: %v", err)
		return
	}

	if err := c.process(); err != nil {
		log.Printf("[Error]: c.process() got err: %v", err)
		return
	}
}

type config struct {
	input   string
	output  string
	dbName  string
	tblName string

	w            io.Writer
	r            *csv.Reader
	fieldsFormat string
}

func (c *config) parse() error {
	if c.input == "" {
		return errors.New("input cannot be empty")
	}
	fd, err := os.Open(c.input)
	if err != nil {
		return err
	}
	c.r = csv.NewReader(fd)

	if c.output == "" {
		c.w = os.Stdout
		log.Println("[Info] output is empty, write to Stdout")
	} else {
		fd, err := os.OpenFile(c.output, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.New("could not open output file")
		}
		c.w = fd
	}

	if c.dbName == "" {
		return errors.New("dbname cannot be empty")
	}

	if c.tblName == "" {
		return errors.New("tablename cannot be empty")
	}

	return nil
}

func (c *config) process() error {
	all, err := c.r.ReadAll()
	if err != nil {
		return err
	}

	sqls := make([]string, len(all)-1)
	for lineIdx, line := range all {
		if lineIdx == 0 {
			// true: 字段行
			quoteSlice(line, "`")
			c.fieldsFormat = strings.Join(line, ",")
			continue
		}
		// false: 数据行
		sqls[lineIdx-1] = c.assemble(line)
	}

	return c.writeTo(sqls)
}

func (c *config) writeTo(sqls []string) (err error) {
	for _, sql := range sqls {
		if _, err = io.WriteString(c.w, sql+"\n"); err != nil {
			return err
		}
	}
	return nil
}

func (c *config) assemble(data []string) string {
	quoteSlice(data, "'")
	d := strings.Join(data, ",")
	return fmt.Sprintf(sqlFormat, c.dbName, c.tblName, c.fieldsFormat, d)
}

func quoteSlice(sl []string, wrap string) {
	for idx, v := range sl {
		sl[idx] = quote(v, wrap)
	}
}

func quote(s string, wrap string) string {
	if s == "NULL" {
		return s
	}
	return wrap + s + wrap
}
