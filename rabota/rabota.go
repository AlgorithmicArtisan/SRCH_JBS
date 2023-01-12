package rabota

import (
	"log"
	"strings"

	"github.com/gocolly/colly"
)

const RABOTABY_URL = "https://rabota.by/search/vacancy?text=&salary=990&area=1002&ored_clusters=true&order_by=publication_time&enable_snippets=true&only_with_salary=true&search_period=1"

var attempts int

// Определяем структуру вакансии
type vacansy struct {
	Title      string
	Url        string
	Org_title  string
	Org_url    string
	Salary     string
	Expirience string
}

// Определяем срез для структур типа vacansy
var overall_jobs []vacansy

// Добавляет вакансию в срез вакансий
func newVacansy(t, u, ot, ou, s, e string) {
	overall_jobs = append(overall_jobs, vacansy{
		Title:      t,
		Url:        u,
		Org_title:  ot,
		Org_url:    ou,
		Salary:     s,
		Expirience: e,
	})
}

func GetVac() []vacansy {
	c := colly.NewCollector()

	c.OnHTML("div.serp-item[data-qa]", func(element *colly.HTMLElement) {
		aElement := element.DOM
		title := aElement.Find("a.serp-item__title").Text()
		url, _ := aElement.Find("a.serp-item__title").Attr("href")
		salary := aElement.Find("span.bloko-header-section-3").Text()
		org_title := aElement.Find("a.bloko-link.bloko-link_kind-tertiary").Text()
		raw_org_url, _ := aElement.Find("a.bloko-link.bloko-link_kind-tertiary").Attr("href")
		org_url := "https://rabota.by/" + raw_org_url

		title = strings.TrimSpace(title)
		url = strings.TrimSpace(url)
		org_title = strings.TrimSpace(org_title)
		org_url = strings.TrimSpace(org_url)
		salary = strings.TrimSpace(salary)

		if len(overall_jobs) != 20 {
			newVacansy(title, url, org_title, org_url, salary, "Опыт не указан")
		} else {
			return
		}
	})

	c.Visit(RABOTABY_URL)

	if len(overall_jobs) < 1 {
		GetVac()
		attempts++
	} else if attempts == 2 {
		log.Fatal("Ошибка запроса на Rabota.by")
	}

	return overall_jobs[:]
}
