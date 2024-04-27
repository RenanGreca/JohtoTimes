package list

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"johtotimes.com/src/post"
	T "johtotimes.com/src/templates"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	category := req.PathValue("category")
	fmt.Println(req)
	fmt.Println("Listing: " + category)

	ListPage(category).Render(req.Context(), w)
}

func ListPage(slug string) templ.Component {
	log.Println("List page: ", slug)
	// entries, err := os.ReadDir("web/posts")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	var posts []post.Post
	// for _, e := range entries {
	// 	// post := internal.ParseMarkdown("web/posts/" + e.Name())
	// 	// if post.Category == slug {
	// 	// 	posts = append(posts, post)
	// 	// }
	// }
	list := List("Category: "+slug, posts)

	return T.Base("Johto Times", list)
}
