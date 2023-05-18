package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type indexPage struct {
	Top             []topdata
	MainNavbar      []mainnavdata
	FeaturedTitle   string
	FeaturedPost    []*featuredpostdata
	MostRecentTitle string
	MostRecentPost  []*mostrecentpostdata
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
	PostID      string `db:"post_id"`
	Title       string `db:"title"`
	Subtitle    string `db:"subtitle"`
	Author      string `db:"author"`
	Authorurl   string `db:"author_url"`
	Publishdate string `db:"publish_date"`
	Imageurl    string `db:"image_url"`
	PostURL     string
}

type mostrecentpostdata struct {
	PostID      string `db:"post_id"`
	Title       string `db:"title"`
	Subtitle    string `db:"subtitle"`
	Author      string `db:"author"`
	Authorurl   string `db:"author_url"`
	Publishdate string `db:"publish_date"`
	Imageurl    string `db:"image_url"`
	PostURL     string
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
	Title    string `db:"title"`
	Subtitle string `db:"subtitle"`
	Image    string `db:"image_url"`
	Content  string `db:"content"`
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

func post(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := mux.Vars(r)["postID"]

		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid order id", 404)
			log.Println(err)
			return
		}

		post, err := postByID(db, postID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Order not found", 404)
				log.Println(err)
				return
			}

			http.Error(w, "Server Error", 500)
			log.Println(err)
			return
		}

		ts, err := template.ParseFiles("pages/post.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}

		data := postIndexPage{
			Top:         posttop(),
			PostContent: post,
			Down:        down(),
		}

		err = ts.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
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

func featuredPost(db *sqlx.DB) ([]*featuredpostdata, error) {
	const query = `
		SELECT
		  post_id,
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

	var featuredposts []*featuredpostdata

	err := db.Select(&featuredposts, query)
	if err != nil {
		return nil, err
	}

	for _, post := range featuredposts {
		post.PostURL = "/post/" + post.PostID
	}

	fmt.Println(featuredposts)

	return featuredposts, nil
}

func mostrecentPost(db *sqlx.DB) ([]*mostrecentpostdata, error) {
	const query = `
		SELECT
		  post_id,
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

	var mostrecent []*mostrecentpostdata

	err := db.Select(&mostrecent, query)
	if err != nil {
		return nil, err
	}

	for _, post := range mostrecent {
		post.PostURL = "/post/" + post.PostID
	}

	fmt.Println(mostrecent)

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
			Header: postblock(),
		},
	}
}

func postblock() []blockdata {
	return []blockdata{
		{
			Logo:   "../static/Sources/Backgrounds/Escape2.svg",
			Navbar: navbardata(),
		},
	}
}

func postByID(db *sqlx.DB, postID int) ([]postcontentdata, error) {
	const query = `
		SELECT
		  title,
		  subtitle,
		  image_url,
		  content
		FROM
		  post
	    WHERE
		  post_id = ?
	`

	var post []postcontentdata

	err := db.Select(&post, query, postID)
	if err != nil {
		return nil, err
	}

	return post, nil
}
