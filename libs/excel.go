package libs

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/extrame/xls"
	"github.com/tealeg/xlsx"
)

func SimpleReadExcel(r io.Reader) ([][][]string, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	xlFile, err := xls.OpenReader(bytes.NewReader(bs), "utf-8")
	if err != nil {
		file, err := xlsx.OpenBinary(bs)
		if err != nil {
			return nil, err
		}
		return file.ToSlice()
	} else {
		count := xlFile.NumSheets()
		table := make([][][]string, count)
		for i := 0; i < count; i++ {
			sheet := xlFile.GetSheet(i)
			nrow := len(sheet.Rows)
			table[i] = make([][]string, nrow)
			for irow := 0; irow < nrow; irow++ {
				row := sheet.Rows[uint16(irow)]
				ncol := len(row.Cols)
				table[i][irow] = make([]string, ncol)
				for icol := 0; icol < ncol; icol++ {
					col := row.Cols[uint16(icol)]
					table[i][irow][icol] = strings.Join(col.String(xlFile), "|")
				}
			}
		}
		return table, nil
	}
}
