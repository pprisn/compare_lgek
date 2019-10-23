package main

import "github.com/kcasctiv/dbf3"

func main(){

// Create new
file := dbf3.New(dbf3.WithLang(dbf3.Lang38))

// Change language driver
//file.SetLang(newDriver)

/*
// Get values
fields := file.Fields()
for idx := 0; idx < file.Rows(); idx++ {
    for _, field := range fields {
        value, _ := file.Get(idx, field.Name())
        fmt.Println(value)
    }
}

// Set value
err := file.Set(idx, "field_name", "value")

// Add row
idx, err := file.NewRow()

// Delete row
err := file.DelRow(idx)

// Check row is deleted
deleted, err := file.Deleted(idx)
*/


// Add field
_ = file.AddField("field_char", dbf3.Character, 10, 0)
_ = file.AddField("field_date", dbf3.Character, 10, 0)
_ = file.AddField("field_num", dbf3.Numeric, 10, 2)

// Set value
// Add row
idx, _ := file.NewRow()
_= file.Set(idx, "field_char", ` Слово`)
_= file.Set(idx, "field_date", ` 22-10-19`)
_= file.Set(idx, "field_num", `1234.27`)


// Save (into file)
_ = file.SaveFile("filename.dbf")

// Save (into writer)
//err := file.Save(writer)

}