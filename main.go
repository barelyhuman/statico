package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"net/http"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v2"
)

type Site struct {
	Name        string `yaml:"name"`
	Link        string `yaml:"link"`
	Description string `yaml:"description"`
}

type Config struct {
	Site          Site `yaml:"site"`
	TemplateNames struct {
		PostTemplateName      string `yaml:"post_template"`
		PostIndexTemplateName string `yaml:"post_index_template"`
		PageTemplateName      string `yaml:"page_template"`
		RSSTemplateName       string `yaml:"rss_template"`
	} `yaml:"template_names"`

	ContentPath    string   `yaml:"content_path"`
	TemplatesPath  string   `yaml:"templates_path"`
	OutPath        string   `yaml:"out_path"`
	PublicFolder   string   `yaml:"public_folder"`
	PostIndexPath  string   `yaml:"post_index_path"`
	IndexedFolders []string `yaml:"indexed_folders"`
	GenerateRss    bool     `yaml:"generate_rss"`
	RssOutPath     string   `yaml:"rss_out_path"`
	Port           string   `yaml:"port"`
}

type ImageURLGen string

type Metadata struct {
	Published  bool   `yaml:"published"`
	Title      string `yaml:"title"`
	Date       Date   `yaml:"date"`
	ImageURL   string `yaml:"image_url"`
	AGImageURL ImageURLGen
	Slug       string
	Content    string
	OutPath    string
}

// Post - container for both metadata and the content of the post
// used to differentiate between a post and just metadata as part of
// other items that might need metadata
type Post struct {
	Site    Site
	Meta    Metadata
	Content string
}

// BlogIndex , used for storing metadata indexes for creatings a blog index
// responsible for the /posts url to show a list of all available posts
type IndexedFiles struct {
	Site  Site
	Files []Metadata
}

// ATOMFeed - Information stuct containing all needed items for for creating a Atom RSS Feed
// TODO: has duplicate structs and can reuse existing structs
type ATOMFeed struct {
	Site struct {
		Name        string
		Link        string
		Description string
	}
	Posts []struct {
		Slug    string
		Title   string
		Link    string
		Content string
		Date    time.Time
	}
}

var (
	markdownProcessor   goldmark.Markdown
	parsedTemplates     *template.Template
	allFilesForIndexing []Metadata
	ConfigRef           *Config
)

var termColors = &TermColors{}

func main() {
	var configFile string
	termColors.Init()
	enableWatch := flag.Bool("watch", false, "Start statico in watch mode")
	enableWatchAlias := flag.Bool("w", false, "alias `-watch`")
	enableServe := flag.Bool("serve", false, "Enable file server")
	enableServeAlias := flag.Bool("s", false, "alias `-serve`")
	configFileFlag := flag.String("config", "", "Config file to use")
	configFileFlagAlias := flag.String("c", "", "alias `-config`")

	flag.Parse()

	configFile = "./config.yml"

	if len(*configFileFlag) > 0 {
		configFile = *configFileFlag
	}

	if len(*configFileFlagAlias) > 0 {
		configFile = *configFileFlagAlias
	}

	ConfigRef = &Config{}
	err := ConfigRef.readConfig(configFile)
	if err != nil {
		log.Fatal("Error reading config: ", err)
	}

	Statico()

	if *enableWatch || *enableWatchAlias || *enableServe || *enableServeAlias {
		waitForKill := make(chan int)

		if *enableWatch || *enableWatchAlias {
			go WatchFiles()
		}

		if *enableServe || *enableServeAlias {
			go ServeFiles()
		}

		<-waitForKill
	}

}

func Success(text string) string {
	return termColors.Bold(termColors.Green(text))
}

func Bullet(text string) string {
	return termColors.Reset(termColors.Bold(text))
}

// Statico - static file and index generator
func Statico() {
	// Clean existing out directory and prefilled values
	err := os.RemoveAll(ConfigRef.OutPath)
	allFilesForIndexing = nil
	markdownProcessor = nil
	parsedTemplates = nil

	if err != nil {
		log.Fatalln(err)
	}

	// Copy the public directory to the out folder for the needed public assets
	err = CopyDir(ConfigRef.PublicFolder, ConfigRef.OutPath)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate Markdown Processor
	markdownProcessor = goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	// Parse all templates from templates directory to hold into
	// the global variable to use later with other dynamic templates generated later
	parsedTemplates, err = template.ParseGlob(ConfigRef.TemplatesPath + "/*")

	if err != nil {
		log.Fatal("Error parsing templates:\n", err)
	}

	// Convert the given content source from config into markdown/compiled HTML files
	err = convertDirectoryToMarkdown(ConfigRef.ContentPath)
	if err != nil {
		log.Fatalf("Failed to convert directory/file at path %v, Error: %v \n", ConfigRef.ContentPath, err)
	}

	// Generate Index files
	for _, indexPath := range ConfigRef.IndexedFolders {
		indexFile, err := os.Create(ConfigRef.OutPath + "/" + indexPath + "/index.html")
		var filesToIndex []Metadata

		if err != nil {
			log.Fatalf("failed to compile index file in folder %v,Error: %v \n ", indexPath, err)
		}

		defer indexFile.Close()

		for _, file := range allFilesForIndexing {
			if file.OutPath == indexPath {
				filesToIndex = append(filesToIndex, file)
			}
		}

		sort.Slice(filesToIndex, func(i, j int) bool {
			return filesToIndex[j].Date.Time.Before(filesToIndex[i].Date.Time)
		})

		err = parsedTemplates.ExecuteTemplate(indexFile, ConfigRef.TemplateNames.PostIndexTemplateName, IndexedFiles{Site: ConfigRef.Site, Files: filesToIndex})

		if err != nil {
			log.Fatalf("failed to compile index file template for folder %v,Error: %v \n ", indexPath, err)
		}

		indexFile.Sync()
	}

	// Generate the rss feed if the config variable is set to true
	if ConfigRef.GenerateRss {
		feed := ATOMFeed{}

		feed.Site = struct {
			Name        string
			Link        string
			Description string
		}{
			Name:        ConfigRef.Site.Name,
			Link:        ConfigRef.Site.Link,
			Description: ConfigRef.Site.Description,
		}

		for _, fileIndex := range allFilesForIndexing {
			feed.Posts = append(feed.Posts,
				struct {
					Slug    string
					Title   string
					Link    string
					Content string
					Date    time.Time
				}{
					Slug:    fileIndex.Slug,
					Title:   fileIndex.Title,
					Link:    feed.Site.Link + "/" + fileIndex.Slug,
					Date:    fileIndex.Date.Time,
					Content: fileIndex.Content,
				},
			)
		}

		rssWriter, err := os.Create(ConfigRef.RssOutPath)
		if err != nil {
			log.Fatal(err)
		}
		defer rssWriter.Close()

		err = parsedTemplates.ExecuteTemplate(rssWriter, ConfigRef.TemplateNames.RSSTemplateName, feed)

		if err != nil {
			log.Fatal(err)
		}

		rssWriter.Sync()
	}

	fmt.Println(
		Success("Static Files Generated to : ") + Bullet(ConfigRef.OutPath),
	)
}

func (cfg *Config) readConfig(configFilePath string) error {
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(file), &cfg)
	if err != nil {
		return err
	}

	return nil
}

// changeFileExtension - simple replace statement wrapper to change the extension of files
func changeFileExtension(fileName string, checkFor string, replaceWith string) string {
	return strings.Replace(fileName, checkFor, replaceWith, 1)
}

// writeToBlog - Write the given html to the blog output
// and compile the html using the given the metadata
func writeToBlog(fileNameHTML string, metadata Metadata) {
	fileToWrite, err := os.Create(ConfigRef.OutPath + "/" + fileNameHTML)
	if err != nil {
		panic(err)
	}
	defer fileToWrite.Close()

	post := Post{
		Site:    ConfigRef.Site,
		Meta:    metadata,
		Content: metadata.Content,
	}

	err = parsedTemplates.ExecuteTemplate(fileToWrite, ConfigRef.TemplateNames.PostTemplateName, post)
	if err != nil {
		log.Fatal(err)
	}

	fileToWrite.Sync()
}

// writeToPage - write to a simple page instead of blog and use the simple
// page template to create the output file, again uses metadata as needed
func writeToPage(fileNameHTML string, content []byte, metadata Metadata) {
	fileToWrite, err := os.Create(ConfigRef.OutPath + "/" + fileNameHTML)
	if err != nil {
		panic(err)
	}
	defer fileToWrite.Close()

	var toHTML bytes.Buffer

	if err := markdownProcessor.Convert(content, &toHTML); err != nil {
		panic(err)
	}

	post := Post{
		Site:    ConfigRef.Site,
		Meta:    metadata,
		Content: toHTML.String(),
	}

	err = parsedTemplates.ExecuteTemplate(fileToWrite, ConfigRef.TemplateNames.PageTemplateName, post)
	if err != nil {
		log.Fatal(err)
	}

	fileToWrite.Sync()
}

// Convert the given directory into markdown and write the processed data to
// the out folder
func convertDirectoryToMarkdown(srcFolder string) error {
	// define prefixes for reading and writing
	pathPrefix := srcFolder
	outPathPrefix := "/"

	// if the source folder and the content folder don't match,
	// aka, you are not on the root directory for processing so add the nested values to the path
	if srcFolder != ConfigRef.ContentPath {
		pathPrefix = ConfigRef.ContentPath + "/" + srcFolder
		outPathPrefix = strings.Replace(srcFolder, "./", "", 1)
	}

	info, err := os.Stat(pathPrefix)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("given source is not a directory")
	}

	files, err := ioutil.ReadDir(pathPrefix)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			err := convertDirectoryToMarkdown(file.Name())
			if err != nil {
				return err
			}
			continue
		}

		err := handleUnprocessedTemplate(pathPrefix, outPathPrefix, file)
		if err != nil {
			return err
		}
	}

	return nil
}

// Check if the extension matches that of a markdown file
func isMarkdownFile(file os.FileInfo) bool {
	extension := strings.SplitN(file.Name(), ".", -1)
	return extension[len(extension)-1] == "md"
}

// Check if the extension matches that of a html file
func isHTMLFile(file os.FileInfo) bool {
	extension := strings.SplitN(file.Name(), ".", -1)
	return extension[len(extension)-1] == "html"
}

// normalizeOtherFileTypeName - as the name specifies, make sure the name doesn't have any spaces when
// used to create a html file
func normalizeOtherFileTypeName(file os.FileInfo) string {
	fileName := changeFileExtension(file.Name(), ".md", "")
	fileName = strings.ReplaceAll(fileName, "-", "")
	return fileName
}

//  handleUnprocessedTemplate - handle the given path, output , file to process
//  by converting an Markdown file with Frontmatter into a blog file
//  a normal markdown file into a simple html file and
//  a html file into a dynamically compiled go html template
func handleUnprocessedTemplate(pathPrefix string, outPathPrefix string, file os.FileInfo) error {
	var err error
	var isHTMLFileBool = isHTMLFile(file)

	os.MkdirAll(ConfigRef.OutPath+"/"+outPathPrefix, os.ModePerm)

	fileData, err := ioutil.ReadFile(pathPrefix + "/" + file.Name())
	if err != nil {
		log.Fatal("Unable to read file:"+pathPrefix+"/"+file.Name()+"\n Error:", err)
	}

	if !isMarkdownFile(file) && !isHTMLFileBool {
		fmt.Println(termColors.Dim("Skipping file, not a markdown file:"), termColors.Bold(file.Name()))
		return nil
	}

	if isHTMLFileBool {
		err = handleHTMLFile(file, fileData, outPathPrefix)
		if err != nil {
			return err
		}
	} else if isMarkdownWithFrontMatter(fileData) {
		return handleMarkdownFile(file, fileData, outPathPrefix)
	} else {
		return handleOtherFile(file, fileData, outPathPrefix)
	}
	return nil
}

// isMarkdownWithFrontMatter - check if the markdown file has any frontmatter
func isMarkdownWithFrontMatter(fileData []byte) bool {
	parts := bytes.SplitN(fileData, []byte("---"), 3)
	return len(parts) == 3
}

//  handleOtherFile - handle files that are not html templates or blog (with frontmatter)
//  based markdown and just a simple markdown file
func handleOtherFile(file os.FileInfo, fileData []byte, outPathPrefix string) error {
	metadata := Metadata{}
	fileNameHTML := changeFileExtension(file.Name(), ".md", ".html")
	metadata.Slug = outPathPrefix + "/" + fileNameHTML
	name := normalizeOtherFileTypeName(file)
	metadata.Title = toTitleCase(name)
	writeToPage(outPathPrefix+"/"+fileNameHTML, fileData, metadata)
	return nil
}

// handleHTMLFile - handle creation and writing of a html file by dynamically creating
// a template , compiling it and then writing it with the same name as the source
// used to create the compiled index.html page of this blog
func handleHTMLFile(file os.FileInfo, fileData []byte, outPathPrefix string) error {
	fileNameHTML := changeFileExtension(file.Name(), ".md", ".html")
	dynTemplateName := file.Name() + "HTML"
	newTemplate := parsedTemplates.New(file.Name() + "HTML")
	parsedHtmlFile, err := newTemplate.Parse(string(fileData))
	if err != nil {
		return err
	}
	writeParsedHTML(ConfigRef.OutPath+"/"+outPathPrefix+"/"+fileNameHTML, parsedHtmlFile, dynTemplateName)
	return nil
}

// handleMarkdownFile - specifically handle markdown files with frontmatter available
// and pass them to compile with the blog-post styled template
func handleMarkdownFile(file os.FileInfo, fileData []byte, outPathPrefix string) error {
	metadata := Metadata{}
	fileNameHTML := changeFileExtension(file.Name(), ".md", ".html")
	metadata.Slug = outPathPrefix + "/" + fileNameHTML
	parts := bytes.SplitN(fileData, []byte("---"), 3)
	if len(parts) != 3 {
		return nil
	}

	if isInSlice(ConfigRef.IndexedFolders, outPathPrefix) {
		metadata.OutPath = outPathPrefix
	}

	err := yaml.Unmarshal(parts[1], &metadata)

	if err != nil {
		return fmt.Errorf("failed to covert frontmatter of file:%v with error: %v \n ", file.Name(), err)
	}

	if metadata.Published {
		var toHTML bytes.Buffer

		err = markdownProcessor.Convert(parts[2], &toHTML)
		if err != nil {
			return err
		}

		metadata.Content = toHTML.String()

		if len(metadata.ImageURL) <= 0 {
			metadata.AGImageURL = ImageURLGen("https://og.reaper.im/api?title=" + metadata.Title + "&fontSize=8&subtitle=at " + ConfigRef.Site.Link)
		}

		allFilesForIndexing = append(allFilesForIndexing, metadata)
		writeToBlog(outPathPrefix+"/"+fileNameHTML, metadata)
	}
	return nil
}

// writeParsedHTML - write the parsed html file to the needed out file
func writeParsedHTML(filePath string, templates *template.Template, templateName string) {
	fileToWrite, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer fileToWrite.Close()
	templates.ExecuteTemplate(fileToWrite, templateName, nil)
	fileToWrite.Sync()
}

// toTitleCase - simple string conversion for strings to title case style
// eg: a fox => A Fox
func toTitleCase(toConv string) string {
	parts := strings.SplitN(toConv, " ", -1)
	result := []string{}
	for _, word := range parts {
		result = append(result,
			strings.ToUpper(string(word[0]))+word[1:],
		)
	}
	return strings.Join(result, " ")
}

func isInSlice(sliceToCheck []string, toSearch string) bool {
	for _, item := range sliceToCheck {
		if item == toSearch {
			return true
		}
	}
	return false
}

func ServeFiles() {
	fs := http.FileServer(http.Dir(ConfigRef.OutPath))
	http.Handle("/", fs)

	port := "3000"

	if len(ConfigRef.Port) > 0 {
		port = ConfigRef.Port
	}

	log.Println("Listening on :" + port + "...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
