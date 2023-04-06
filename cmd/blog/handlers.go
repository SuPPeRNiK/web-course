package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type indexPage struct {
	Top             []topdata
	MainNavbar      []mainnavdata
	FeaturedTitle   string
	FeaturedPost    []featuredpostdata
	MostRecentTitle string
	MostRecentPost  []mostrecentpostdata
	Down            []downdata
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

type featuredpostdata struct {
	Title       string `db:"title"`
	Subtitle    string `db:"subtitle"`
	Author      string `db:"author"`
	Authorurl   string `db:"author_url"`
	Publishdate string `db:"publish_date"`
	Imageurl    string `db:"image_url"`
}

type mostrecentpostdata struct {
	Title       string `db:"title"`
	Subtitle    string `db:"subtitle"`
	Author      string `db:"author"`
	Authorurl   string `db:"author_url"`
	Publishdate string `db:"publish_date"`
	Imageurl    string `db:"image_url"`
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

type postIndexPage struct {
	Top         []posttopdata
	PostContent []postcontentdata
	Down        []downdata
}

type posttopdata struct {
	Header   []blockdata
	Title    string
	Subtitle string
}

type postcontentdata struct {
	Background string
	P1         string
	P2         string
	P3         string
	P4         string
}

func index(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		featuredposts, err := featuredPost(db)
		if err != nil {
			http.Error(w, "Error", 500)
			log.Println(err)
			return
		}

		mostrecent, err := mostrecentPost(db)
		if err != nil {
			http.Error(w, "Error", 500)
			log.Println(err)
			return
		}

		ts, err := template.ParseFiles("pages/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}

		data := indexPage{
			Top:             top(),
			MainNavbar:      mainnavbardata(),
			FeaturedTitle:   "Featured Posts",
			FeaturedPost:    featuredposts,
			MostRecentTitle: "Most Recent",
			MostRecentPost:  mostrecent,
			Down:            down(),
		}

		err = ts.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/the-road-ahead.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	data := postIndexPage{
		Top:         posttop(),
		PostContent: postpage(),
		Down:        down(),
	}

	err = ts.Execute(w, data)
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

func featuredPost(db *sqlx.DB) ([]featuredpostdata, error) {
	const query = `
		SELECT
		  title,
		  subtitle,
		  author,
		  author_url,
		  publish_date,
		  image_url
		FROM
		  post
		WHERE featured = 1
	`

	var featuredposts []featuredpostdata

	err := db.Select(&featuredposts, query)
	if err != nil {
		return nil, err
	}

	return featuredposts, nil
}

func mostrecentPost(db *sqlx.DB) ([]mostrecentpostdata, error) {
	const query = `
		SELECT
		  title,
		  subtitle,
		  author,
		  author_url,
		  publish_date,
		  image_url
		FROM
		  post
		WHERE featured = 0
	`

	var mostrecent []mostrecentpostdata

	err := db.Select(&mostrecent, query)
	if err != nil {
		return nil, err
	}

	return mostrecent, nil
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

func posttop() []posttopdata {
	return []posttopdata{
		{
			Header:   block(),
			Title:    "The Road Ahead",
			Subtitle: "The road ahead might be paved - it might not be.",
		},
	}
}

func postpage() []postcontentdata {
	return []postcontentdata{
		{
			Background: "../static/Sources/Backgrounds/the-road-ahead-background.jpg",
			P1:         "Dark spruce forest frowned on either side the frozen waterway. The trees had been stripped by a recent wind of their white covering of frost, and they seemed to lean towards each other, black and ominous, in the fading light. A vast silence reigned over the land. The land itself was a desolation, lifeless, without movement, so lone and cold that the spirit of it was not even that of sadness. There was a hint in it of laughter, but of a laughter more terrible than any sadness — a laughter that was mirthless as the smile of the sphinx, a laughter cold as the frost and partaking of the grimness of infallibility. It was the masterful and incommunicable wisdom of eternity laughing at the futility of life and the effort of life. It was the Wild, the savage, frozen-hearted Northland Wild.",

			P2: "But there was life, abroad in the land and defiant. Down the frozen waterway toiled a string of wolfish dogs. Their bristly fur was rimed with frost. Their breath froze in the air as it left their mouths, spouting forth in spumes of vapour that settled upon the hair of their bodies and formed into crystals of frost. Leather harness was on the dogs, and leather traces attached them to a sled which dragged along behind. The sled was without runners. It was made of stout birch-bark, and its full surface rested on the snow. The front end of the sled was turned up, like a scroll, in order to force down and under the bore of soft snow that surged like a wave before it. On the sled, securely lashed, was a long and narrow oblong box. There were other things on the sled—blankets, an axe, and a coffee-pot and frying-pan; but prominent, occupying most of the space, was the long and narrow oblong box.",

			P3: "In advance of the dogs, on wide snowshoes, toiled a man. At the rear of the sled toiled a second man. On the sled, in the box, lay a third man whose toil was over,—a man whom the Wild had conquered and beaten down until he would never move nor struggle again. It is not the way of the Wild to like movement. Life is an offence to it, for life is movement; and the Wild aims always to destroy movement. It freezes the water to prevent it running to the sea; it drives the sap out of the trees till they are frozen to their mighty hearts; and most ferociously and terribly of all does the Wild harry and crush into submission man—man who is the most restless of life, ever in revolt against the dictum that all movement must in the end come to the cessation of movement.",

			P4: "But at front and rear, unawed and indomitable, toiled the two men who were not yet dead. Their bodies were covered with fur and soft-tanned leather. Eyelashes and cheeks and lips were so coated with the crystals from their frozen breath that their faces were not discernible. This gave them the seeming of ghostly masques, undertakers in a spectral world at the funeral of some ghost. But under it all they were men, penetrating the land of desolation and mockery and silence, puny adventurers bent on colossal adventure, pitting themselves against the might of a world as remote and alien and pulseless as the abysses of space.",
		},
	}
}
