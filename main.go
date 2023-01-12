// go build -ldflags -H=windowsgui -o JobsToObsidian.exe main.go
package main

import (
	"JobSearching/belmeta"
	"JobSearching/praca"
	"JobSearching/rabota"
	"log"
	"os"
	"time"
)

func main() {
	// Путь к файлу в Obsidian, содержащего информацио о работе
	OBS_PATH := os.ExpandEnv("$USERPROFILE\\Dropbox\\ObsidianDATA\\main\\Other\\РаботаTEST.md")

	//Получаем данные из различных сайтов
	rabota_jobs := rabota.GetVac()
	praca_jobs := praca.GetVac()
	belmeta_jobs := belmeta.GetVac()

	//Создаём файл хранения всех данных в Obsidian
	file, err := os.Create(OBS_PATH)
	if err != nil {
		log.Fatal(err)
	}
	os.OpenFile(OBS_PATH, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	now := time.Now().Format("02 Jan 06 15:04 MST")
	file.Write([]byte("#Работа\n" + now + "\n"))

	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(rabota_jobs); i++ {
		file.Write([]byte("[" + rabota_jobs[i].Title + "]" + "(" + rabota_jobs[i].Url + ")" + "\n" + "[" + rabota_jobs[i].Org_title + "]" + "(" + rabota_jobs[i].Org_url + ")" + "\n" + rabota_jobs[i].Salary + "\n" + rabota_jobs[i].Expirience + "\n\n"))
	}

	for i := 0; i < len(praca_jobs); i++ {
		file.Write([]byte("[" + praca_jobs[i].Title + "]" + "(" + praca_jobs[i].Url + ")" + "\n" + "[" + praca_jobs[i].Org_title + "]" + "(" + praca_jobs[i].Org_url + ")" + "\n" + praca_jobs[i].Salary + "\n" + praca_jobs[i].Expirience + "\n\n"))
	}

	for i := 0; i < len(belmeta_jobs); i++ {
		file.Write([]byte("[" + belmeta_jobs[i].Title + "]" + "(" + belmeta_jobs[i].Url + ")" + "\n" + "[" + belmeta_jobs[i].Org_title + "]" + "(" + belmeta_jobs[i].Org_url + ")" + "\n" + belmeta_jobs[i].Salary + "\n" + belmeta_jobs[i].Expirience + "\n\n"))
	}
	file.Close()
}
