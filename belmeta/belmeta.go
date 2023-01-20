package belmeta

import (
	"JobSearching/ini"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

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

var attempts int

func GetVac() []ini.Vacansy {
	c := colly.NewCollector()

	c.OnHTML("article.job", func(element *colly.HTMLElement) {
		aElement := element.DOM
		title := aElement.Find("h2.title").Text()
		raw_url, _ := aElement.Find("h2.title").Find("a").Attr("href")
		org_title := aElement.Find("div.job-data.company").Text()
		salary := aElement.Find("div.job-data.salary").Text()
		//Избавляемся в строке c url от всего, кроме ID
		re := regexp.MustCompile(`[0-9]+`)
		raw_url_purification := re.FindAllString(raw_url, -1)
		ID := strings.Join(raw_url_purification, " ")
		//Присваиваем ссылке ID
		url := ("https://belmeta.com/viewjob?id=" + ID)

		if len(overall_jobs) != 10 {
			title = strings.TrimSpace(title)
			url = strings.TrimSpace(url)
			org_title = strings.TrimSpace(org_title)
			salary = strings.TrimSpace(salary)
			newVacansy(title, url, org_title, "", salary, "Опыт не указан")
			c.Visit(url)
		} else {
			return
		}
	})

	c.Visit(ini.BELMETA_URL)

	if len(overall_jobs) < 1 {
		GetVac()
		attempts++
	} else if attempts == 2 {
		log.Fatal("Ошибка запроса на Belmeta.by")
		return nil
	}

	return overall_jobs[:]

}
