package handler

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"johtotimes.com/src/assert"
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
	body := categoryPage(category, page)

	render(body, isHTMX(req), category, w)
	// categoryPage(category, page).Render(req.Context(), w)
}

func categoryPage(slug string, page int) templ.Component {
	assert.LogDebug("List category page: ", slug)

	db := database.Connect()
	defer db.Close()
	posts := db.Posts.GetByCategorySlug(slug, page, 10)

	cat := model.GetCategoryFromFile(constants.CategoriesPath, slug)
	description := renderHTML(cat.Description)

	return templates.ListTemplate(cat.Plural, slug, description, posts, page)
	// return templates.Base("Johto Times", list)
}
