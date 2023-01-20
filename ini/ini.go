package ini

import "os"

// Путь к файлу в Obsidian, содержащего информацио о работе
var OBS_PATH string = os.ExpandEnv("$USERPROFILE\\Dropbox\\ObsidianDATA\\main\\Other\\Работа.md")

//Путь к сайту BELMETA
const BELMETA_URL string = "https://belmeta.com/vacansii?l=Минск&sf=900&sort=date"

//Путь к сайту RABOTA BY
const RABOTABY_URL string = "https://rabota.by/search/vacancy?text=&salary=990&area=1002&ored_clusters=true&order_by=publication_time&enable_snippets=true&only_with_salary=true&search_period=1"

// Определяем структуру вакансии
type Vacansy struct {
	Title      string
	Url        string
	Org_title  string
	Org_url    string
	Salary     string
	Expirience string
}
