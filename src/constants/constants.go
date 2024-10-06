package constants

type Link struct {
	Name string
	URL  string
	Icon string
}

type CSS struct {
	Path string
}

type Colors struct {
	Background string
	White      string
	Text       string
	NeonRed    string
	NeonBlue   string
}

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
	{Name: "Issues", URL: "/archive/issues"},
	{Name: "Interviews", URL: "/category/interview"},
	{Name: "Features", URL: "/category/feature"},
	{Name: "Opinions", URL: "/category/opinion"},
	{Name: "Collections", URL: "/category/collection"},
	{Name: "Memories", URL: "/category/memory"},
	{Name: "Mailbag", URL: "/mailbag"},
	{Name: "News", URL: "/news"},
}

var Tabs = []Link{
	{Name: "Issues", URL: "/issues", Icon: "book@2x.png"},
	{Name: "Archive", URL: "/archive", Icon: "archivebox@2x.png"},
	{Name: "Search", URL: "/search", Icon: "magnifyingglass@2x.png"},
	{Name: "Community", URL: "/community", Icon: "person.3@2x.png"},
	{Name: "About", URL: "/about", Icon: "info.square@2x.png"},
}

const Title = "Johto Times"

// const URL = "https://johtotimes.renangreca.com"
// const URL = "http://localhost:3000"
const AssetPath = "web"
const StylesPath = AssetPath + "/styles"
const ScriptsPath = AssetPath + "/scripts"
const ImgPath = AssetPath + "/img"
const AudioPath = AssetPath + "/audio"
const IconsPath = ImgPath + "/icons"
const PostsPath = AssetPath + "/posts"
const MailbagPath = AssetPath + "/mailbag"
const NewsPath = AssetPath + "/news"
const IssuesPath = AssetPath + "/issues"
const CategoriesPath = AssetPath + "/categories"

var PostTypePath = map[byte]string{
	'P': PostsPath,
	'M': MailbagPath,
	'N': NewsPath,
	'I': IssuesPath,
}

// Add to this as we want more css files
var Stylesheets = [...]string{
	"/fonts.css",
	"/theme.css",
	"/style.css",
	"/fonts.css",
	"/header.css",
	"/list.css",
	"/tabbar.css",
	"/archive.css",
	"/single.css",
}
