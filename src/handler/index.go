package handler

import (
	"log"

	"github.com/a-h/templ"
	"johtotimes.com/src/database"
	"johtotimes.com/src/list"
	"johtotimes.com/src/post"
	T "johtotimes.com/src/templates"
)

func IndexPage() templ.Component {
	log.Printf("Rendering index page")
	db := database.Connect()
	defer db.Close()
	retrievedPosts, err := db.Posts.GetPage(10, 10)
	if err != nil {
		log.Fatal(err)
	}
	posts := []post.Post{}
	for _, p := range retrievedPosts {
		posts = append(posts, p)
		log.Printf("Found post %s\n", p.Title)
	}
	log.Printf("%d", len(posts))

	list := list.List("", posts)

	return T.Base("Johto Times", list)
}
