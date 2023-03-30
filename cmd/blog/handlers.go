package main

import (
	"html/template"
	"log"
	"net/http"
)

type indexPage struct {
	Top        []topdata
	MainNavbar []mainnavdata
	Content    []contentdata
	Down       []downdata
}

// Top

type topdata struct {
	Background string
	Header     []blockdata
	Headline   []headlinedata
}

type blockdata struct {
	Logo   string
	Navbar []navdata
}

type navdata struct {
	First  string
	Second string
	Third  string
	Fourth string
}

type headlinedata struct {
	Title    string
	Subtitle string
	Button   string
}

// !Top
// MainNavbar

type mainnavdata struct {
	First  string
	Second string
	Third  string
	Fourth string
	Fiveth string
	Sixth  string
}

// !MainNavbar
// Content

type contentdata struct {
	FeaturedPosts []featuredcontent
	MostRecent    []mostrecentcontent
}

type featuredcontent struct {
	Title string
	Posts []featuredposts
}

type featuredposts struct {
	Post1 []featuredpost
	Post2 []featuredpost
}

type featuredpost struct {
	Background string
	Title      string
	Subtitle   string
	PostInfo   []postinfo
}

type mostrecentcontent struct {
	Title string
	Posts []mostrecentposts
}

type mostrecentposts struct {
	Post1 []mostrecentpost
	Post2 []mostrecentpost
	Post3 []mostrecentpost
	Post4 []mostrecentpost
	Post5 []mostrecentpost
	Post6 []mostrecentpost
}

type mostrecentpost struct {
	Background string
	Title      string
	Subtitle   string
	PostInfo   []postinfo
}

type postinfo struct {
	Image string
	Name  string
	Date  string
}

// !Content
// Down

type downdata struct {
	Background string
	Header     []downheaderdata
	Footer     []blockdata
}

type downheaderdata struct {
	Title  string
	Button string
}

// !Down
func index(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/index.html") // Главная страница блога
	if err != nil {
		http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
		log.Println(err.Error())                    // Используем стандартный логгер для вывода ошбики в консоль
		return                                      // Не забываем завершить выполнение ф-ии
	}

	data := indexPage{
		Top:        top(),
		MainNavbar: mainnavbardata(),
		Content:    content(),
		Down:       down(),
	}

	err = ts.Execute(w, data) // Заставляем шаблонизатор вывести шаблон в тело ответа
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/the-road-ahead.html") // Главная страница блога
	if err != nil {
		http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
		log.Println(err.Error())                    // Используем стандартный логгер для вывода ошбики в консоль
		return                                      // Не забываем завершить выполнение ф-ии
	}

	data := indexPage{}

	err = ts.Execute(w, data) // Заставляем шаблонизатор вывести шаблон в тело ответа
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func top() []topdata {
	return []topdata{
		{
			Background: "../static/Sources/Backgrounds/header-1section-background.jpg",
			Header:     block(),
			Headline:   headline(),
		},
	}
}

func block() []blockdata {
	return []blockdata{
		{
			Logo:   "../static/Sources/Backgrounds/Escape.svg",
			Navbar: navbardata(),
		},
	}
}

func navbardata() []navdata {
	return []navdata{
		{
			First:  "HOME",
			Second: "CATEGORIES",
			Third:  "ABOUT",
			Fourth: "CONTACT",
		},
	}
}

func mainnavbardata() []mainnavdata {
	return []mainnavdata{
		{
			First:  "Nature",
			Second: "Photography",
			Third:  "Relaxation",
			Fourth: "Vacation",
			Fiveth: "Travel",
			Sixth:  "Adventure",
		},
	}
}

func headline() []headlinedata {
	return []headlinedata{
		{
			Title:    "Let's do it together",
			Subtitle: "We travel the world in search of stories. Come along for the ride.",
			Button:   "View Latest Posts",
		},
	}
}

func content() []contentdata {
	return []contentdata{
		{
			FeaturedPosts: featuredcontents(),
			MostRecent:    mostrecentcontents(),
		},
	}
}

func featuredcontents() []featuredcontent {
	return []featuredcontent{
		{
			Title: "Futured Posts",
			Posts: featuredcontentposts(),
		},
	}
}

func featuredcontentposts() []featuredposts {
	return []featuredposts{
		{
			Post1: featuredpost_id1(),
			Post2: featuredpost_id2(),
		},
	}
}

func featuredpost_id1() []featuredpost {
	return []featuredpost{
		{
			Background: "../static/Sources/Backgrounds/the-road-ahead-background.jpg",
			Title:      "The Road Ahead",
			Subtitle:   "The road ahead might be paved - it might not be.",
			PostInfo:   featured_postinfo_id1(),
		},
	}
}

func featuredpost_id2() []featuredpost {
	return []featuredpost{
		{
			Background: "../static/Sources/Backgrounds/from-top-down-background.jpg",
			Title:      "From Top Down",
			Subtitle:   "Once a year, go someplace you’ve never been before.",
			PostInfo:   featured_postinfo_id2(),
		},
	}
}

func featured_postinfo_id1() []postinfo {
	return []postinfo{
		{
			Image: "../static/Sources/Authors/Mat-Vogles.svg",
			Name:  "Mat Vogles",
			Date:  "September 25, 2015",
		},
	}
}

func featured_postinfo_id2() []postinfo {
	return []postinfo{
		{
			Image: "../static/Sources/Authors/Mat-Vogles.svg",
			Name:  "Mat Vogles",
			Date:  "September 25, 2015",
		},
	}
}

func mostrecentcontents() []mostrecentcontent {
	return []mostrecentcontent{
		{
			Title: "Most Recent",
			Posts: mostrecentcontentposts(),
		},
	}
}

func mostrecentcontentposts() []mostrecentposts {
	return []mostrecentposts{
		{
			Post1: mostrecentpost_id1(),
			Post2: mostrecentpost_id2(),
			Post3: mostrecentpost_id3(),
			Post4: mostrecentpost_id4(),
			Post5: mostrecentpost_id5(),
			Post6: mostrecentpost_id6(),
		},
	}
}

func mostrecentpost_id1() []mostrecentpost {
	return []mostrecentpost{
		{
			Background: "../static/Sources/Backgrounds/Still-Stading-background.jpg",
			Title:      "Still Standing Tall",
			Subtitle:   "Life begins at the end of your comfort zone.",
			PostInfo:   mostrecent_postinfo_id1(),
		},
	}
}

func mostrecentpost_id2() []mostrecentpost {
	return []mostrecentpost{
		{
			Background: "../static/Sources/Backgrounds/Sunny-Side-Up-background.jpg",
			Title:      "Sunny Side Up",
			Subtitle:   "No place is ever as bad as they tell you it’s going to be.",
			PostInfo:   mostrecent_postinfo_id2(),
		},
	}
}

func mostrecentpost_id3() []mostrecentpost {
	return []mostrecentpost{
		{
			Background: "../static/Sources/Backgrounds/Water-Falls-background.jpg",
			Title:      "Water Falls",
			Subtitle:   "We travel not to escape life, but for life not to escape us.",
			PostInfo:   mostrecent_postinfo_id3(),
		},
	}
}

func mostrecentpost_id4() []mostrecentpost {
	return []mostrecentpost{
		{
			Background: "../static/Sources/Backgrounds/Through-the-Mist-background.jpg",
			Title:      "Through the Mist",
			Subtitle:   "Travel makes you see what a tiny place you occupy in the world.",
			PostInfo:   mostrecent_postinfo_id4(),
		},
	}
}

func mostrecentpost_id5() []mostrecentpost {
	return []mostrecentpost{
		{
			Background: "../static/Sources/Backgrounds/Awaken-Early-background.jpg",
			Title:      "Awaken Early",
			Subtitle:   "Not all those who wander are lost.",
			PostInfo:   mostrecent_postinfo_id5(),
		},
	}
}

func mostrecentpost_id6() []mostrecentpost {
	return []mostrecentpost{
		{
			Background: "../static/Sources/Backgrounds/Try-it-Always-background.jpg",
			Title:      "Try it Always",
			Subtitle:   "The world is a book, and those who do not travel read only one page.",
			PostInfo:   mostrecent_postinfo_id6(),
		},
	}
}

func mostrecent_postinfo_id1() []postinfo {
	return []postinfo{
		{
			Image: "../static/Sources/Authors/William-wong.svg",
			Name:  "William Wang",
			Date:  "9/25/2015",
		},
	}
}

func mostrecent_postinfo_id2() []postinfo {
	return []postinfo{
		{
			Image: "../static/Sources/Authors/Mat-Vogles.svg",
			Name:  "Mat Vogles",
			Date:  "9/25/2015",
		},
	}
}

func mostrecent_postinfo_id3() []postinfo {
	return []postinfo{
		{
			Image: "../static/Sources/Authors/Mat-Vogles.svg",
			Name:  "Mat Vogles",
			Date:  "9/25/2015",
		},
	}
}

func mostrecent_postinfo_id4() []postinfo {
	return []postinfo{
		{
			Image: "../static/Sources/Authors/William-wong.svg",
			Name:  "William Wang",
			Date:  "9/25/2015",
		},
	}
}

func mostrecent_postinfo_id5() []postinfo {
	return []postinfo{
		{
			Image: "../static/Sources/Authors/Mat-Vogles.svg",
			Name:  "Mat Vogles",
			Date:  "9/25/2015",
		},
	}
}

func mostrecent_postinfo_id6() []postinfo {
	return []postinfo{
		{
			Image: "../static/Sources/Authors/Mat-Vogles.svg",
			Name:  "Mat Vogles",
			Date:  "9/25/2015",
		},
	}
}

func down() []downdata {
	return []downdata{
		{
			Background: "../static/Sources/Backgrounds/footer-background1.jpg",
			Header:     downheader(),
			Footer:     block(),
		},
	}
}

func downheader() []downheaderdata {
	return []downheaderdata{
		{
			Title:  "Stay in Touch",
			Button: "Submit",
		},
	}
}
