package internal

// This should rarely change once we've decided which links to have
var SocialLinks = [...]Link{
	{Name: "instagram", URL: "https://instagram.com/"},
	{Name: "twitter", URL: "https://twitter.com"},
	{Name: "youtube", URL: "https://youtube.com"},
	{Name: "twitch", URL: "https://twitch.tv"},
	{Name: "discord", URL: "https://discordapp.com/"},
}

// Probably this should become dynamic according to some DB data
var Categories = [...]Link{
	{Name: "Interviews", URL: "/interviews"},
	{Name: "Features", URL: "/features"},
	{Name: "Opinions", URL: "/opinions"},
	{Name: "Collections", URL: "/collections"},
	{Name: "Memories", URL: "/memories"},
	{Name: "Mailbag", URL: "/mailbag"},
}

const Title = "Johto Times"

const URL = "https://johtotimes.renangreca.com"
const AssetPath = "web"
const StylesPath = "/" + AssetPath + "/styles"
const ScriptsPath = "/" + AssetPath + "/scripts"
const ImgPath = AssetPath + "/img"
const IconsPath = ImgPath + "/icons"

// Add to this as we want more css files
var Stylesheets = [...]string{
	"/fonts.css",
	"/theme.css",
	"/style.css",
	"/fonts.css",
	"/header.css",
	"/list.css",
}
