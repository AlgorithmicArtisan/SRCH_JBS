package praca

import (
	"JobSearching/ini"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
)

var attempts int

// Определяем срез для структур типа vacansy
var overall_jobs []ini.Vacansy

// Добавляет вакансию в срез вакансий
func newVacansy(t, u, ot, ou, s, e string) {
	overall_jobs = append(overall_jobs, ini.Vacansy{
		Title:      t,
		Url:        u,
		Org_title:  ot,
		Org_url:    ou,
		Salary:     s,
		Expirience: e,
	})
}

// Возвращает индекс страницы на которой было найдено начало бесплатных обьявлений.
func getPageIndex() (int, error) {
	var page_index int
	lookfor := []byte(`standard vac-small`)
	for {
		// Получаем байт-код итерируемой стриницы
		response, err := http.Get(fmt.Sprintf("https://praca.by/search/vacancies/?page=%d&search[cities][Минск]=1&search[query]=&search[cities-radius][Минск]=1&search[query-text-params][headline]=0&form-submit-btn=Искать", page_index))
		if err != nil {
			log.Fatal("Ошибка при обращении к странице сайта: " + response.Status)
		}
		// Преобразуем слайс байт
		data, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

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

func GetVac() []ini.Vacansy {
	//Получаем индекс страницы на Praca.by, содержащей в себе бесплатные обьявления.
	page_index, err := getPageIndex()
	if err != nil {
		log.Fatal(err)
	}

	c := colly.NewCollector()

	//Находим родительский класс иеррархия которого нас интересует
	c.OnHTML("li.standard.vac-small", func(element *colly.HTMLElement) {
		aElement := element.DOM
		// Название
		title := aElement.Find("a.vac-small__title-link").Text()
		//Ссылка
		url, _ := aElement.Find("a.vac-small__title-link").Attr("href")
		//Название организации
		org_title := aElement.Find("a.vac-small__organization").Text()
		//Ссылка организации
		org_url, _ := aElement.Find("a.vac-small__organization").Attr("href")
		//Зарплата
		salary := aElement.Find("span.salary-dotted").Text()
		//Требуемый опыт
		expirience := aElement.Find("div.vac-small__experience").Text()

		title = strings.TrimSpace(title)
		url = strings.TrimSpace(url)
		org_title = strings.TrimSpace(org_title)
		org_url = strings.TrimSpace(org_url)
		salary = strings.TrimSpace(salary)
		expirience = strings.TrimSpace(expirience)

		if len(overall_jobs) != 10 {
			newVacansy(title, url, org_title, org_url, salary, expirience)
		} else {
			return
		}
	})

	// Запуск парсинга
	c.Visit(fmt.Sprintf("https://praca.by/search/vacancies/?page=%d&search[cities][Минск]=1&search[query]=&search[cities-radius][Минск]=1&search[query-text-params][headline]=0&form-submit-btn=Искать", page_index))

	if len(overall_jobs) < 1 {
		GetVac()
		attempts++
	} else if attempts == 2 {
		log.Fatal("Ошибка запроса на Praca.by")
	}

	return overall_jobs[:]
}
