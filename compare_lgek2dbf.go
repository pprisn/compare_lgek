//
//list_qr.csv - содержит список QR ЛГЭК
//input.dbf - список оплат вошедших на перечисление в п.п
//Требуется заполнить поле LSH значением PersAcc из QR кода в списке list_qr.csv
//Результатом является файл file_result.dbf с модификацией поля LSH данными из PersAcc=
//Для успешного выполнения программе требуется передать список list_qr.csv соответствующий оплатам в файле file_input.dbf

package main

import (
	"encoding/csv"
	"fmt"
	"github.com/LindsayBradford/go-dbf/godbf"
	"io"
	"os"
	"unicode/utf8"
	"log"
	"flag"
	"strings"
	"time"
)

var file_qr     = flag.String("file_qr",    "list_qr.csv"  , `Файл со списком QR кодов оплат, полученный из почтовой платежной системы`)
var file_input  = flag.String("file_input", "IM.dbf"      , `Файл со списком оплат, сформированный в платежной почтовой системе`)
var file_result = flag.String("file_result","IM_.dbf"     ,`Имя файла куда сохранить результат работы`)

//Слайс для хранения строк из файла  file_qr.csv
var d1 [][]string

func main() {

	flag.Parse() //парсим флаги запуска программы
	var floger *os.File
	var err error
	if floger, err = os.OpenFile("compare_lgek2dbf.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) ; err != nil {
		panic(err)
	}
	defer floger.Close()

	//Наполним слайс d1 данными
	separ := rune('|')
	d1 = loadCSV(*file_qr, separ)



	dbfTable, err := godbf.NewFromFile(*file_input, "866")
	if err != nil {
		panic(err)
	}
        fields := dbfTable.Fields()

        log.SetOutput(floger)
        t0 := time.Now()
        log.Printf("СТАРТ %v \n", t0)

        //Определение параметров для csv.NewWriter в поток os.Stdout
	delimiter := ";"
	comma, _ := utf8.DecodeRuneInString(delimiter)
	out := csv.NewWriter(os.Stdout)
	out.Comma = comma

        //Полечение заголовка из dbf файла в виде []byte
	r, err := os.Open(*file_input)
	if err != nil {
		panic(err)
	}
	defer r.Close()

        //Растеч размера заголовка dbf файла
	sizeHeader := len(fields)*32 + 32 + 1
        //Инициализация объема памяти для []byte заголовка 
	header := make([]byte, sizeHeader)
        //Чтение sizeHeader числа байт из потока r в массив header[:]
	_, err = io.ReadFull(r, header[:])

        // Обнуление числа записей в структуре заголовка 
        header[4],header[5],header[6],header[7] = 0,0,0,0
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", header)
        //---------------------------------------------
        // Создание структуры для нового файла dbf из массива []byte 
	newdbfTable, _ := godbf.NewFromByteArray(header, "866")

	//fmt.Printf("%+v\n", fields)
	/*
	   [
	    {name:PLP fieldType:67 length:8 decimalPlaces:0 fieldStore:[80 76 80 0 0 0 0 0 0 0 0 67 0 0 0 0 8 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:DT_BANK fieldType:68 length:8 decimalPlaces:0 fieldStore:[68 84 95 66 65 78 75 0 0 0 0 68 0 0 0 0 8 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:ID_DB fieldType:67 length:10 decimalPlaces:0 fieldStore:[73 68 95 68 66 0 0 0 0 0 0 67 0 0 0 0 10 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:LSH fieldType:67 length:24 decimalPlaces:0 fieldStore:[76 83 72 0 0 0 0 0 0 0 0 67 0 0 0 0 24 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:PRINT_SUM fieldType:67 length:13 decimalPlaces:0 fieldStore:[80 82 73 78 84 95 83 85 77 0 0 67 0 0 0 0 13 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:DATA fieldType:68 length:8 decimalPlaces:0 fieldStore:[68 65 84 65 0 0 0 0 0 0 0 68 0 0 0 0 8 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:TIME_OPL fieldType:67 length:10 decimalPlaces:0 fieldStore:[84 73 77 69 95 79 80 76 0 0 0 67 0 0 0 0 10 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:SUMMA fieldType:78 length:20 decimalPlaces:2 fieldStore:[83 85 77 77 65 0 0 0 0 0 0 78 0 0 0 0 20 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:FIO fieldType:67 length:50 decimalPlaces:0 fieldStore:[70 73 79 0 0 0 0 0 0 0 0 67 0 0 0 0 50 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:ADRES fieldType:67 length:60 decimalPlaces:0 fieldStore:[65 68 82 69 83 0 0 0 0 0 0 67 0 0 0 0 60 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:HV fieldType:78 length:12 decimalPlaces:4 fieldStore:[72 86 0 0 0 0 0 0 0 0 0 78 0 0 0 0 12 4 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:GV fieldType:78 length:12 decimalPlaces:4 fieldStore:[71 86 0 0 0 0 0 0 0 0 0 78 0 0 0 0 12 4 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:LIV fieldType:78 length:10 decimalPlaces:0 fieldStore:[76 73 86 0 0 0 0 0 0 0 0 78 0 0 0 0 10 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:OT fieldType:78 length:12 decimalPlaces:4 fieldStore:[79 84 0 0 0 0 0 0 0 0 0 78 0 0 0 0 12 4 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:OT_KUB fieldType:78 length:12 decimalPlaces:3 fieldStore:[79 84 95 75 85 66 0 0 0 0 0 78 0 0 0 0 12 3 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:OT_VREM fieldType:78 length:12 decimalPlaces:0 fieldStore:[79 84 95 86 82 69 77 0 0 0 0 78 0 0 0 0 12 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:PENY fieldType:78 length:12 decimalPlaces:2 fieldStore:[80 69 78 89 0 0 0 0 0 0 0 78 0 0 0 0 12 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:REGNO fieldType:67 length:10 decimalPlaces:0 fieldStore:[82 69 71 78 79 0 0 0 0 0 0 67 0 0 0 0 10 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:CHEQUE fieldType:67 length:10 decimalPlaces:0 fieldStore:[67 72 69 81 85 69 0 0 0 0 0 67 0 0 0 0 10 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:FEEPRC fieldType:67 length:10 decimalPlaces:0 fieldStore:[70 69 69 80 82 67 0 0 0 0 0 67 0 0 0 0 10 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	    {name:MANUAL fieldType:67 length:1 decimalPlaces:0 fieldStore:[77 65 78 85 65 76 0 0 0 0 0 67 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]}
	   ]
	*/

        //Выделение объема памяти для слайса fieldRow 
	fieldRow := make([]string, len(fields))
        //Наполнение слайса fieldRow значениями имен полей fields
	for i := 0; i < len(fields); i++ {
		fieldRow[i] = fields[i].Name()
	}

	//Вывод заголовка dbf в формае csv 
	out.Write(fieldRow)
	out.Flush()

	// Вывод строк Output rows
        var idx int
        var N int
        N = len(fields)

	for i := 0; i < dbfTable.NumberOfRecords(); i++ {
		row := dbfTable.GetRowAsSlice(i)
		        // Добавление записей в новую структуру newdbfTable
			idx = newdbfTable.AddNewRecord()
                        // Записываем значения полей
			for j:=0; j < N; j++ {
                                if(j == 3){
                                      row[j] = isfound( row[j] )
                                      fmt.Printf("%+v\n", row[j])
				}
				_ = newdbfTable.SetFieldValueByName(idx, fields[j].Name(), row[j] )
				
			}
                //Вывод записей из dbf файла в формае csv на экран
		out.Write(row)
		out.Flush()


	}

	newdbfTable.SaveFile(*file_result)

       t1 := time.Now()
       log.Printf("Обработано %d записей.",dbfTable.NumberOfRecords())
       log.Printf("Завершение работы, общее время выполнения %v сек.\n", t1.Sub(t0))

	return

}


func loadCSV(path string, separ rune) [][]string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", path, err.Error())
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = separ          //';' разделитель полей в файле scv
	r.Comment = '#'          // символ комментария
	r.LazyQuotes = true      // разрешить ковычки в полях
	rows, err := r.ReadAll() //прочитать весь файл в массив [][]string
	if err != nil {
		log.Fatalln("Cannot read CSV data:", err.Error())
	}
	return rows
}

func isfound(key string) string {
	var vline string
	if key == "" {
		return ""
	}
        //Ищем значение key по строкам слайса d1
	for _, line := range d1 {
		vline = strings.Join(line, ";")   // объединим поля в строку с применением разделителя ;
		if strings.Contains(vline, key) { // если найдено в строке ключевое слово 
			//ST00012|
			//Name=РСО АО ЛГЭК|
			//PersonalAcc=40702810435000104067|
			//BankName=Липецкое отд № 8593 ПАО Сбербанк|
			//BIC=044206604|
			//CorrespAcc=30101810800000000604|
			//PayeeINN=4825066916|
			//KPP=482501001|
			//Sum=26411|
			//PersAcc=100050001238|
			//PaymPeriod=0919|TechCode=01|
			return line[9][8:] //PersAcc=100050001238
		}
	}
	return key

}
