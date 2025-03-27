package excel

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func WithExcel(path string, headerNames []interface{}) chan<- []interface{} {
	excel := excelize.NewFile()
	defer excel.Close()

	writer, err := excel.NewStreamWriter("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	styleID, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "777777",
			Size:  14,
			Bold:  true,
		},
	})

	err = writer.SetRow("A1", headerNames, excelize.RowOpts{
		Height:       16,
		StyleID:      styleID,
		OutlineLevel: 0,
	})
	if err != nil {
		log.Fatal(err)
	}

	var (
		total = 1
		out   = make(chan []interface{}, 1)
	)

	go func() {
		for raw := range out {
			err = writer.SetRow(fmt.Sprintf("A%d", total+1), raw)
			if err != nil {
				return
			}
		}

		err = writer.Flush()
		if err != nil {
			log.Fatal(err)
		}

		err = excel.SaveAs(path)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return out
}
