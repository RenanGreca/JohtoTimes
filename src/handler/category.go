package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"johtotimes.com/src/category"
	"johtotimes.com/src/database"
	"johtotimes.com/src/internal"
	"johtotimes.com/src/templates"
)

func CategoryHandler(w http.ResponseWriter, req *http.Request) {
	category := req.PathValue("category")
	fmt.Println(req)
	fmt.Println("Listing: " + category)

	page := getPageNumber(req)
	CategoryPage(category, page).Render(req.Context(), w)
}

func CategoryPage(slug string, page int) templ.Component {
	log.Println("List category page: ", slug)

	db := database.Connect()
	defer db.Close()
	posts, err := db.Posts.GetByCategorySlug(slug, page, 10)
	if err != nil {
		log.Fatal(err)
	}

	cat := category.GetFromFile(internal.CategoriesPath, slug)
	description := unsafe(cat.Contents)
	plural := cat.Metadata.Plural

	list := templates.ListTemplate(plural, slug, description, posts)

	return templates.Base("Johto Times", list)
}
