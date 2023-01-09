// go build -ldflags -H=windowsgui -o JobsToObsidian.exe main.go
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	//"reflect"

	"github.com/gocolly/colly"
)

func getError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getPageIndex() (int, error) {
	var page_index int
	lookfor := []byte(`standard vac-small`)
	for {
		response, err := http.Get(fmt.Sprintf("https://praca.by/search/vacancies/?page=%d&search[cities][Минск]=1&search[query]=&search[cities-radius][Минск]=1&search[query-text-params][headline]=0&form-submit-btn=Искать", page_index))
		getError(err)
		data, err := io.ReadAll(response.Body)
		getError(err)
		defer response.Body.Close()
		if bytes.Index(data, lookfor) != -1 {
			break
		} else {
			//fmt.Println("На странице " + fmt.Sprint(page_index) + " бесплатных обьявлений не найдено!")
			page_index++
			if page_index > 150 {
				return -1, errors.New("Превышен лимит на обработку страниц!")
			}
		}
	}
	return page_index, nil
}

func main() {
	const OBS_PATH string = "C:\\Users\\Администратор\\Dropbox\\ObsidianDATA\\main\\Other\\Работа.md"
	//Получаем индекс страницы, содержащей в себе бесплатные обьявления.
	page_index, err := getPageIndex()
	if err != nil {
		log.Fatal(err)
		return
	}

	//Создаём файл хранения данных
	file, err := os.Create(OBS_PATH)
	if err != nil {
		log.Fatal(err)
	}
	os.OpenFile(OBS_PATH, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	file.Write([]byte("#Работа\n\n"))
	defer file.Close()

	c := colly.NewCollector(
		// Посещение только домена praca.by
		colly.AllowedDomains("praca.by"),
	)

	c.OnHTML("li.standard.vac-small", func(element *colly.HTMLElement) {
		// Найти элемент <a class="vac-small__title-link"> внутри текущего элемента
		aElement := element.DOM
		// Получить текст этого элемента и ссылку
		text := aElement.Find("a.vac-small__title-link").Text()
		url, _ := aElement.Find("a.vac-small__title-link").Attr("href")
		Org_url, _ := aElement.Find("a.vac-small__organization").Attr("href")
		Org_text := aElement.Find("a.vac-small__organization").Text()
		raw_salary := aElement.Find("span.salary-dotted").Text()
		expirience := aElement.Find("div.vac-small__experience").Text()

		//Модифицируем строку с зарплатой
		re := regexp.MustCompile(`[\d-]+`)
		raw_2_salary := re.FindAllString(raw_salary, -1)
		salary := strings.Join(raw_2_salary, " ")

		//Заносим информацию в файл md
		os.OpenFile(OBS_PATH, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		file.Write([]byte("[" + text + "]" + "(" + url + ")" + " " + "[" + Org_text + "]" + "(" + Org_url + ")" + "\n" + salary + " " + expirience + "\n\n\n"))
		//fmt.Println("Done!")
	})
	defer file.Close()

	// Показать адрес посещаемой страницы
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	// Запуск парсинга
	c.Visit(fmt.Sprintf("https://praca.by/search/vacancies/?page=%d&search[cities][Минск]=1&search[query]=&search[cities-radius][Минск]=1&search[query-text-params][headline]=0&form-submit-btn=Искать", page_index))

}
