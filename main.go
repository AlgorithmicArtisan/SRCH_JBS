// go build -ldflags -H=windowsgui -o JobsToObsidian.exe main.go
package main

import (
	sound "JobSearching/PlaySound"
	"JobSearching/belmeta"
	"JobSearching/ini"
	"JobSearching/praca"
	"JobSearching/rabota"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

// Функция проверки файла на предмет его существования
func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func main() {
	//Получаем данные из различных сайтов и совмещаем их в единый слайс
	jobs := append([]ini.Vacansy{}, belmeta.GetVac()...)
	jobs = append(jobs, praca.GetVac()...)
	jobs = append(jobs, rabota.GetVac()...)

	jobs_backup := make([]ini.Vacansy, len(jobs))
	copy(jobs_backup, jobs)

	var newvac_count int

	//Сравниваем данные из прошлой итерации с новыми, чтобы выявить наиболее свежие вакансии
	if checkFileExists("jobs.json") {
		// Загрузка среза из файла
		fmt.Println("Файл существует!")
		var loadedJobs []ini.Vacansy
		file, _ := os.Open("jobs.json")
		json.NewDecoder(file).Decode(&loadedJobs)
		file.Close()

		//Свапаем и подсвечиваем свежайшие предложения
		var countains int
		for y := 0; y < len(jobs); y++ {
			for x := 0; x < len(loadedJobs); x++ {
				if jobs[y].Title != loadedJobs[x].Title || jobs[y].Org_title != loadedJobs[x].Org_title {
					countains++
				}
			}
			if countains == len(loadedJobs) {
				fmt.Println(jobs[y].Title + " Свежая Вакансия! <--------------------------------------------------------------------------------")
				jobs[y].Title = "\U0001F6A9 " + jobs[y].Title
				jobs[y], jobs[newvac_count] = jobs[newvac_count], jobs[y]
				newvac_count++
			} else {
				fmt.Println(jobs[y].Title + " Уже есть в базе")
			}
			countains = 0
		}
	} else {
		fmt.Println("Файл несуществует!")
	}

	//Создаём файл хранения всех данных в Obsidian
	file, err := os.Create(ini.OBS_PATH)
	if err != nil {
		log.Fatal(err)
	}
	os.OpenFile(ini.OBS_PATH, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// Создаём временной штамп
	now := time.Now().Format("02 Jan 06 15:04 MST\n")
	file.Write([]byte("#Работа\n" + now + "\n"))

	if err != nil {
		log.Fatal(err)
	}
	//Заносим данные в файл
	for i := 0; i < len(jobs); i++ {
		file.Write([]byte("[" + jobs[i].Title + "]" + "(" + jobs[i].Url + ")" + "\n" + "[" + jobs[i].Org_title + "]" + "(" + jobs[i].Org_url + ")" + "\n" + jobs[i].Salary + "\n" + jobs[i].Expirience + "\n\n"))
	}

	file.Close()

	//Создаём бэкап данных в файл JSON для последующей обработки в следующем цикле
	jfile, _ := os.Create("jobs.json")
	json.NewEncoder(jfile).Encode(jobs_backup)
	jfile.Close()

	//Воспроизводим звуковое оповещение при наличии новых вакансий за последние 15 минут
	if newvac_count > 0 {
		sound.SndPlaySoundW("alert.wav", sound.SND_SYNC)
	}
}
