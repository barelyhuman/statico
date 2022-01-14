package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v2"
)

// Enums, defined with the prefix `E` to save the common variable names
// from getting polluted
const (
	EMarkdownFile = iota
	EHTMLFile
	EOtherFile
)

type FileManifestItem struct {
	Input, Output, OutputDir string
	InputType                int
	NeedsIndexing            bool
	IndexingFolder           string
	FileMeta                 map[string]interface{}
	Slug                     string
}

type Statico struct {
	config        *Config
	templates     *template.Template
	filesManifest []FileManifestItem
}

func (app *Statico) ServeFiles() {
	fmt.Println(Dim("Starting Dev Server..."))
	fs := http.FileServer(http.Dir(app.config.OutPath))
	http.Handle("/", fs)

	port := "3000"

	if len(app.config.Port) > 0 {
		port = app.config.Port
	}

	fmt.Println(Bullet(">> Listening on " + Info(":"+port)))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (app *Statico) WatchFiles() error {
	w := watcher.New()
	// experimenting with 2
	// TODO: add it to the config so people can change if needed
	w.SetMaxEvents(2)

	go handleFileChanges(w, app)

	if err := w.AddRecursive(app.config.ContentPath); err != nil {
		return err
	}

	if err := w.AddRecursive(app.config.PublicFolder); err != nil {
		return err
	}

	if err := w.AddRecursive(app.config.TemplatesPath); err != nil {
		return err
	}

	if err := w.AddRecursive(app.config.PublicFolder); err != nil {
		return err
	}

	if err := w.Start(time.Millisecond * 500); err != nil {
		return err
	}

	return nil
}

func handleFileChanges(w *watcher.Watcher, app *Statico) error {
	for {
		select {
		case <-w.Event:
			fmt.Println(
				Dim(logTime()) +
					Info("Recompiling ..."),
			)
			app.Build()
			fmt.Println(
				Dim(logTime()) +
					Success("Compiled! Static files generated to: "+Bullet(app.config.OutPath)),
			)
		case err := <-w.Error:
			return err
		case <-w.Closed:
			return nil
		}
	}
}

// simple time log prefix for logging statements
func logTime() string {
	return "[" + time.Now().Format(
		"15:04:05",
	) + "] "
}

func (app *Statico) Build() {
	app.CleanOutPath()
	app.CopyPublic()
	app.ReadTemplates()
	app.ReadContentFiles()
	app.ProcessFilesManifest()
	app.CreateIndexFiles()
}

func (app *Statico) CleanOutPath() {
	bail(os.RemoveAll(app.config.OutPath))
}

func (app *Statico) ReadTemplates() {
	parsedTemplates, err := template.ParseGlob(app.config.TemplatesPath + "/*")
	app.templates = parsedTemplates
	bail(err)
}

func (app *Statico) ReadContentFiles() {
	app.filesManifest = readDirectoryRecursive(app.config.ContentPath, app.config)
}

// Copy the public folder into the output folder
func (app *Statico) CopyPublic() {
	bail(
		CopyDir(app.config.PublicFolder, app.config.OutPath),
	)
}

func readDirectoryRecursive(source string, config *Config) []FileManifestItem {
	// define prefixes for reading and writing
	pathPrefix := source
	outPathPrefix := config.OutPath
	var paths []FileManifestItem

	// if the source folder and the content folder don't match,
	// aka, you are not on the root directory for processing so add the nested values to the path
	if source != config.ContentPath {
		pathPrefix = config.ContentPath + "/" + source
		outPathPrefix = outPathPrefix + "/" + strings.Replace(source, "./", "", 1)
	}

	// throw an error if the given source directory is not actually
	// a directory
	info, err := os.Stat(pathPrefix)
	bail(err)
	if !info.IsDir() {
		bail(
			fmt.Errorf("given source is not a directory"),
		)
	}

	// get all files from the given directory
	files, err := ioutil.ReadDir(pathPrefix)
	bail(err)

	for _, file := range files {
		if file.IsDir() {
			// recursively get files inside the source that also
			// can be converted to html
			_paths := readDirectoryRecursive(file.Name(), config)
			paths = append(paths, _paths...)
			bail(err)
			continue
		}

		fileName := file.Name()
		inputType := EOtherFile
		toIndex := false
		var indexingFolder string

		if isMarkdownFile(file) {
			fileName = changeFileExtension(file.Name(), ".md", ".html")
			inputType = EMarkdownFile
		}

		if isHTMLFile(file) {
			inputType = EHTMLFile
		}

		for _, i := range config.IndexedFolders {
			if source == i {
				indexingFolder = source
				toIndex = true
			}
		}

		// add the needed files to the manifest of files
		paths = append(paths, FileManifestItem{
			Input:          pathPrefix + "/" + file.Name(),
			Output:         outPathPrefix + "/" + fileName,
			OutputDir:      outPathPrefix + "/",
			InputType:      inputType,
			NeedsIndexing:  toIndex,
			IndexingFolder: indexingFolder,
			Slug:           contructUrlFromPaths(config.OutPath, outPathPrefix+"/"+fileName),
		})
	}
	return paths
}

func (app *Statico) ProcessFilesManifest() {

	markdownProcessor := goldmark.New(
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

	for manifestIndex, file := range app.filesManifest {
		var toHTML bytes.Buffer
		var metadata = make(map[string]interface{})
		file.FileMeta = make(map[string]interface{})

		fileData, err := os.ReadFile(file.Input)
		bail(err)

		if file.InputType == EHTMLFile {
			CopyFile(file.Input, file.Output)
			continue
		}

		if file.InputType != EMarkdownFile {
			continue
		}

		markdownBytes := fileData

		parts := bytes.SplitN(fileData, []byte("---"), 3)
		if len(parts) == 3 {
			markdownBytes = parts[2]
			yaml.Unmarshal(parts[1], &metadata)
		}

		if err := markdownProcessor.Convert(markdownBytes, &toHTML); err != nil {
			bail(err)
		}

		err = os.MkdirAll(file.OutputDir, os.ModePerm)
		bail(err)

		fileToWrite, err := os.Create(file.Output)
		bail(err)
		defer fileToWrite.Close()

		metadata["slug"] = file.Slug
		app.filesManifest[manifestIndex].FileMeta = metadata

		err = app.templates.ExecuteTemplate(fileToWrite, app.config.TemplateNames.PageTemplateName, struct {
			Site    Site
			Meta    map[string]interface{}
			Content template.HTML
		}{
			Site:    app.config.Site,
			Meta:    metadata,
			Content: template.HTML(toHTML.String()),
		})
		bail(err)

		fileToWrite.Sync()
	}

	// read input file from the manifest
	// write to output file from the manifest
}

func (app *Statico) CreateIndexFiles() {
	fmt.Println(Info("[Statico] Compiling Index Templates"))
	for _, toIndex := range app.config.IndexedFolders {
		indexFile, err := os.Create(path.Join(app.config.OutPath, toIndex, "index.html"))
		bail(err)
		defer indexFile.Close()

		var filesMeta []map[string]interface{}

		for _, file := range app.filesManifest {
			if file.NeedsIndexing && file.IndexingFolder == toIndex {
				filesMeta = append(filesMeta, file.FileMeta)
			}
		}

		sort.Slice(filesMeta, func(i, j int) bool {
			dateOne, err := parseStringToDate(
				InterfaceToString(filesMeta[j]["date"]),
			)
			bail(err)
			dateTwo, err := parseStringToDate(
				InterfaceToString(filesMeta[i]["date"]),
			)
			bail(err)
			return dateOne.Time.Before(dateTwo.Time)
		})

		err = app.templates.ExecuteTemplate(indexFile, app.config.TemplateNames.PostIndexTemplateName, struct {
			Site  Site
			Files []map[string]interface{}
			Meta  []map[string]interface{}
		}{
			Site:  app.config.Site,
			Files: filesMeta,
		},
		)

		if err != nil {
			fmt.Println(Warn("[Warn] " + err.Error()))
		}

	}
}

func contructUrlFromPaths(outPath string, filePath string) string {
	url := strings.Replace(filePath, outPath+"/", "", 1)
	return url
}

func parseStringToDate(buf string) (Date, error) {
	var tt time.Time
	var err error

	formatsToCheck := []string{
		"02-01-2006",
		"2006-01-02",
		"02/01/2006",
		"2/01/2006",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05.999999999Z07:00",
	}

	for _, format := range formatsToCheck {
		tt, err = time.Parse(format, strings.TrimSpace(buf))
		if err != nil {
			continue
		} else {
			break
		}
	}

	if err != nil {
		return Date{}, err
	}

	return Date{
		Time: tt,
	}, nil
}
