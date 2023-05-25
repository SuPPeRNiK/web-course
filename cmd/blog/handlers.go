package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	Theme       string `db:"theme"`
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
	Theme       string `db:"theme"`
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

type adminPage struct {
	Header   []headeradmindata
	MainTop  []maintopdata
	MainInfo []maininfodata
	Content  []contentdata
}

type headeradmindata struct {
	Logo    string
	Avatar  string
	ExitURL string
	Exit    string
}

type maintopdata struct {
	Title    string
	Subtitle string
	Button   string
}

type maininfodata struct {
	Title   string
	Fields  []fieldsdata
	Preview []previewdata
}

type fieldsdata struct {
	Title          string
	Description    string
	AuthorName     string
	AuthorPhoto    string
	AuthorPhotoURL string
	Upload         string
	Date           string
	TitleImage     string
	BigImageURL    string
	SmallImageURL  string
	BigNote        string
	SmallNote      string
}

type previewdata struct {
	Article  []articledata
	PostCard []postcarddata
}

type articledata struct {
	Label    string
	FrameURL string
	Title    string
	Subtitle string
	Imageurl string
}

type postcarddata struct {
	Label          string
	FrameURL       string
	Imageurl       string
	Title          string
	Subtitle       string
	AuthorPhotoURL string
	AuthorName     string
	Data           string
}

type contentdata struct {
	Title   string
	Comment string
}

type loginpage struct {
	Background string
	Header     []headerlogindata
	Main       []mainlogindata
}

type headerlogindata struct {
	Escape string
	Title  string
}

type mainlogindata struct {
	Title  string
	Email  string
	Pass   string
	Button string
}

type createPostRequest struct {
	Title           string `json:"title"`
	SubTitle        string `json:"subtitle"`
	AuthorName      string `json:"authorname"`
	AuthorPhoto     string `json:"authorphoto"`
	AuthorPhotoName string `json:"authorphotoname"`
	Data            string `json:"data"`
	BigImage        string `json:"bigimage"`
	BigImageName    string `json:"bigimagename"`
	SmallImage      string `json:"smallimage"`
	SmallImageName  string `json:"smallimagename"`
	Content         string `json:"content"`
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

func admin(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/admin.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	data := adminPage{
		Header:   AdminHeader(),
		MainTop:  AdminMainTop(),
		MainInfo: AdminMainInfo(),
		Content:  AdminContent(),
	}

	err = ts.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/login.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	data := loginpage{
		Background: "../static/Sources/Backgrounds/login_background.png",
		Header:     LoginHeader(),
		Main:       LoginMain(),
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

func featuredPost(db *sqlx.DB) ([]*featuredpostdata, error) {
	const query = `
		SELECT
		  post_id,
		  title,
		  subtitle,
		  author,
		  author_url,
		  publish_date,
		  image_url,
		  theme
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
		  image_url,
		  theme
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

func AdminHeader() []headeradmindata {
	return []headeradmindata{
		{
			Logo:    "../static/Sources/svg_files/escape_author_white.svg",
			Avatar:  "../static/Sources/authors/avatar.jpg",
			ExitURL: "/login",
			Exit:    "../static/Sources/svg_files/log_out.svg",
		},
	}
}

func AdminMainTop() []maintopdata {
	return []maintopdata{
		{
			Title:    "New Post",
			Subtitle: "Fill out the form bellow and publish your article",
			Button:   "Publish",
		},
	}
}

func AdminMainInfo() []maininfodata {
	return []maininfodata{
		{
			Title:   "Main Information",
			Fields:  fields(),
			Preview: preview(),
		},
	}
}

func fields() []fieldsdata {
	return []fieldsdata{
		{
			Title:          "Title",
			Description:    "Short description",
			AuthorName:     "Author Name",
			AuthorPhoto:    "Author Photo",
			AuthorPhotoURL: "../static/Sources/svg_files/photo_icon.svg",
			Upload:         "Upload",
			Date:           "Publish Date",
			TitleImage:     "Hero image",
			BigImageURL:    "../static/Sources/Backgrounds/hero_image_big.png",
			SmallImageURL:  "../static/Sources/Backgrounds/hero_image_small.png",
			BigNote:        "Size up to 10mb. Format: png, jpeg, gif.",
			SmallNote:      "Size up to 5mb. Format: png, jpeg, gif.",
		},
	}
}

func preview() []previewdata {
	return []previewdata{
		{
			Article:  article(),
			PostCard: postcard(),
		},
	}
}

func article() []articledata {
	return []articledata{
		{
			Label:    "Article preview",
			FrameURL: "../static/Sources/Backgrounds/aritcle_frame.png",
			Title:    "New Post",
			Subtitle: "Please, enter any description",
			Imageurl: "../static/Sources/Backgrounds/image_not_selected.png",
		},
	}
}

func postcard() []postcarddata {
	return []postcarddata{
		{
			Label:          "Post card preview",
			FrameURL:       "../static/Sources/Backgrounds/post_card_frame.png",
			Imageurl:       "../static/Sources/Backgrounds/image_not_selected.png",
			Title:          "New Post",
			Subtitle:       "Please, enter any description",
			AuthorPhotoURL: "../static/Sources/svg_files/photo_icon.svg",
			AuthorName:     "Enter author name",
			Data:           "4/19/2023",
		},
	}
}

func AdminContent() []contentdata {
	return []contentdata{
		{
			Title:   "Content",
			Comment: "Post content (plain text)",
		},
	}
}

func LoginHeader() []headerlogindata {
	return []headerlogindata{
		{
			Escape: "../static/Sources/svg_files/escape_author_white.svg",
			Title:  "Log in to start creating",
		},
	}
}

func LoginMain() []mainlogindata {
	return []mainlogindata{
		{
			Title:  "Log In",
			Email:  "Email",
			Pass:   "Password",
			Button: "Log In",
		},
	}
}

func createPost(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "1Error", 500)
			log.Println(err.Error())
			return
		}

		var req createPostRequest

		err = json.Unmarshal(reqData, &req)
		if err != nil {
			http.Error(w, "2Error", 500)
			log.Println(err.Error())
			return
		}

		fileAuthorImg := req.AuthorPhoto[strings.IndexByte(req.AuthorPhoto, ',')+1:]
		authorImg, err := base64.StdEncoding.DecodeString(fileAuthorImg)
		if err != nil {
			http.Error(w, "2Error", 500)
			log.Println(err.Error())
			return
		}
		fileAuthor, err := os.Create("static/Sources/Authors/" + req.AuthorPhotoName)
		if err != nil {
			http.Error(w, "2Error", 500)
			log.Println(err.Error())
			return
		}
		_, err = fileAuthor.Write(authorImg)
		if err != nil {
			http.Error(w, "2Error", 500)
			log.Println(err.Error())
			return
		}

		fileBigImage := req.BigImage[strings.IndexByte(req.BigImage, ',')+1:]
		bigImg, err := base64.StdEncoding.DecodeString(fileBigImage)
		if err != nil {
			http.Error(w, "2Error", 500)
			log.Println(err.Error())
			return
		}
		fileBig, err := os.Create("static/Sources/Backgrounds/" + req.BigImageName)
		if err != nil {
			http.Error(w, "2Error", 500)
			log.Println(err.Error())
			return
		}
		_, err = fileBig.Write(bigImg)
		if err != nil {
			http.Error(w, "2Error", 500)
			log.Println(err.Error())
			return
		}

		fileSmallImage := req.SmallImage[strings.IndexByte(req.SmallImage, ',')+1:]
		smallImg, err := base64.StdEncoding.DecodeString(fileSmallImage)
		if err != nil {
			http.Error(w, "2Error", 500)
			log.Println(err.Error())
			return
		}
		fileSmall, err := os.Create("static/Sources/Authors/" + req.AuthorPhotoName)
		if err != nil {
			http.Error(w, "2Error", 500)
			log.Println(err.Error())
			return
		}
		_, err = fileSmall.Write(smallImg)
		if err != nil {
			http.Error(w, "2Error", 500)
			log.Println(err.Error())
			return
		}

		err = saveData(db, req)

		if err != nil {
			http.Error(w, "bd", 500)
			log.Println(err.Error())
			return
		}

		return
	}
}

func saveData(db *sqlx.DB, req createPostRequest) error {
	const query = `
		INSERT INTO
			post
		(
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_url,
			content
		)
		VALUES
		(
			?,
			?,
			?,
			CONCAT('../../static/Sources/Authors/', ?),
			?,
			CONCAT('../../static/Sources/Backgrounds/', ?),
			?
		)
	`

	_, err := db.Exec(query, req.Title, req.SubTitle, req.AuthorName, req.AuthorPhotoName, req.Data, req.BigImageName, req.Content)
	return err
}
