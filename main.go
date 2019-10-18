//list_qr.csv - содержит список QR ЛГЭК
//input.csv - список оплат вошедших на перечисление в п.п
//Требуется заполнить поле LSH значением PersAcc из QR кода в списке list_qr.csv
//Результатом является файл result_.csv с модификацией поля LSH данными из PersAcc=
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"flag"
	"os"
	"strings"
	"time"
)
var dN [][]string
var d1 [][]string

var file_qr     = flag.String("file_qr",   "list_qr.csv" , `Файл со списком QR кодов оплат, полученный из почтовой платежной системы`)
var file_input  = flag.String("file_input","input.csv"   , `Файл со списком оплат, сформированный в платежной почтовой системе`)
var file_result = flag.String("file_result","result_.csv", `Имя файла куда сохранить результат работы`)

func main() {

	flag.Parse() //парсим флаги запуска программы
	var floger *os.File
	var err error
	if floger, err = os.OpenFile("compare_lgek.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) ; err != nil {
		panic(err)
	}
	defer floger.Close()


	var iline int
	fout, err := os.Create(*file_result)
	if err != nil {
		log.Fatalf("Cannot open result_csvs %s\n", path, err.Error())
	}
	defer fout.Close()

        //file_QR, file_input := path() //заменил на работу с flag

	//загрузим файлы в массивы
   
        log.SetOutput(floger)
        t0 := time.Now()
        log.Printf("СТАРТ %v \n", t0)

	separ := rune('|')
	d1 = loadCSV(*file_qr, separ)
	separ = rune(';')
	dN = loadCSV(*file_input, separ)

	for _, line := range dN {
		//vline = strings.Join(line, ";") // объединим поля в строку с применением разделителя ;
		//if strings.Contains(vline, key){ // если найдено в строке ключевое слово
		iline = iline + 1
                if (iline == 2) {
                log.Printf("Обработка реестра с данными по платежному поручению %s от %s \n", line[0],line[1])
		}
		//PLP;DT_BANK;ID_DB;LSH;PRINT_SUM;DATA;TIME_OPL;SUMMA;FIO;ADRES;HV;GV;LIV;OT;OT_KUB;OT_VREM;PENY;REGNO;CHEQUE;FEEPRC;MANUAL
		//116836;16.10.2019;398043;1111735;0919;14.10.2019;08:09:34;48,57;;;177,0000;0,0000;0;;;;0,00;;60336;1,3;N
		fmt.Fprintf(fout, "%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;\n",
			line[0], line[1], line[2], isfound(line[3]), line[4], line[5], line[6], line[7], line[8], line[9], line[10],
			line[11], line[12], line[13], line[14], line[15], line[16], line[17], line[18], line[19], line[20])
		//	}

	}
       fmt.Println(iline)
       t1 := time.Now()
       log.Printf("Обработано %d записей.",iline)
       log.Printf("Завершение работы, общее время выполнения %v сек.\n", t1.Sub(t0))

}

func isfound(key string) string {
	var vline string
	if key == "" {
		return ""
	}
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
			return line[9][8:]
		}
	}
	return key

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

func path() (string, string) {
	if len(os.Args) < 3 {
		return "list_qr.csv", "input.csv"
	}
	return os.Args[1], os.Args[2]
}
