package main

import (
	"database/sql"
	"fmt"
	_ "github.com/nakagami/firebirdsql"
)

func main() {
	var a string
	//user:password@servername[:port_number]/database_name_or_file[?params1=value1[&param2=value2]...]
	conn, err := sql.Open("firebirdsql", "SNPopurey:~BSbBZFa@10.193.124.48/D:/FBDATA/PPSMAIN/PPSWEB.FDB")
	if err != nil {
		fmt.Println(err)
        }

	defer conn.Close()

	/*
	   	sele := `SELECT PAYMENTS.pay_barcode as a from PAYMENTS
	   left outer join PARTNER as P on PAYMENTS.PARTNER_ID = P.ID
	   left outer join MAINPOSTOFFICE as M on PAYMENTS.MAINPOSTOFFICE_ID = M.ID
	   left outer join transfers as tr on Payments.transfer_id = tr.ID
	   left outer join registry_po as rpo on tr.REGISTRY_PO_ID = rpo.ID
	   where PAYMENTS.DATE_PAYM >= "01.10.2019" and  rpo.NUM= 119165  and
	     ((PAYMENTS.PARTNER_ID = 623) or (PAYMENTS.PARTNER_ID =621) or (PAYMENTS.PARTNER_ID=622))
	        order by  PAYMENTS.FCNUMBER`
	*/
	sele := "SELECT PAYMENTS.pay_barcode as a from PAYMENTS where PAYMENTS.DATE_PAYM >= '10.10.2019'"

	rows, err := conn.Query(sele)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v \n", rows)
		for rows.Next() {
			rows.Scan(&a)
			fmt.Println(a)

		}
	}

}
