package list

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"johtotimes.com/src/database"
	"johtotimes.com/src/post"
	"johtotimes.com/src/templates"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	category := req.PathValue("category")
	fmt.Println(req)
	fmt.Println("Listing: " + category)

	ListPage(category).Render(req.Context(), w)
}

func ListPage(slug string) templ.Component {
	log.Println("List page: ", slug)

	var posts []post.Post

	list := List("Category: "+slug, posts)

	return templates.Base("Johto Times", list)
}

func Mailbag(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling request to " + req.URL.Path)
	page := 0
	pagestr := req.PathValue("page")
	if len(pagestr) > 0 {
		var err error
		page, err = strconv.Atoi(pagestr)
		if err != nil {
			log.Fatalf("Error converting %q to int\n", req.PathValue("page"))
		}
	}
	PostType('M', page).Render(req.Context(), w)
}

func PostType(t byte, page int) templ.Component {
	log.Printf("Rendering home page")
	db := database.Connect()
	defer db.Close()
	retrievedPosts, err := db.Posts.GetPage('M', page, 10)
	if err != nil {
		log.Fatal(err)
	}
	posts := []post.Post{}
	for _, p := range retrievedPosts {
		posts = append(posts, p)
	}
	log.Printf("Found %d posts", len(posts))

	list := List("", posts)

	return templates.Base("Johto Times", list)
}

func HomePage() templ.Component {
	log.Printf("Rendering home page")
	db := database.Connect()
	defer db.Close()
	retrievedPosts, err := db.Posts.GetPage('P', 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	posts := []post.Post{}
	for _, p := range retrievedPosts {
		posts = append(posts, p)
	}
	log.Printf("Found %d posts", len(posts))

	list := List("", posts)

	return templates.Base("Johto Times", list)
}
