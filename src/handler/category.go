package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"johtotimes.com/src/constants"
	"johtotimes.com/src/database"
	"johtotimes.com/src/model"
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
	posts := db.Posts.GetByCategorySlug(slug, page, 10)

	cat := model.GetCategoryFromFile(constants.CategoriesPath, slug)
	description := renderHTML(cat.Description)

	list := templates.ListTemplate(cat.Plural, slug, description, posts)

	return templates.Base("Johto Times", list)
}
