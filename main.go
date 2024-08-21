package main

import (
	"log"
	"time"
)

type Acc struct {
	Login string
	Pass  string
}

func main() {

	saturday, sunday := GetNextWeekendDates()
	sunday = sunday.Add(time.Hour * 24)
	saturday = saturday.Add(-time.Hour * 24)
	text := "Здравствуйте. Если Вы хотите сообщить о нежелательной реакции на лекарственный препарат, свяжитесь с уполномоченным лицом по фармаконадзору по телефону +375445867515 или по электронной почте head_pv@ft.by"
	fromDate := saturday.Format("20060102") + "133000Z"
	toDate := sunday.Format("20060102") + "045900Z"

	for _, acc := range getAccs() {
		log.Print(acc.Login)
		sendToZimbra(acc, text, fromDate, toDate)
	}

}
