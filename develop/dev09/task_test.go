package main_test

import (
	"bytes"
	"dev09/internal/app"
	"dev09/internal/creator"
	"dev09/internal/downloader"
	ex "dev09/internal/extractor"
	"dev09/utils"
	"golang.org/x/net/html"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func newTestServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			w.Write([]byte(`<html>
            <body>
                <a href="/link1.css">Link 1</a>
                <a href="/gojo/jojo">Link 2</a>
            </body>
            </html>`))
		case "/gojo/jojo":
			w.Write([]byte(`<html>
            <body>
                <a href="/link1.css">Link 1</a>
                <a href="/link2">Link 2</a>
            </body>
            </html>`))
		}

	}))

	return ts
}

func TestCreateFileAndDirWithHTMLExt(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	file, err := creator.CreateFile("https://www.marinohall.com/peterhof/banketnyj-zal-na-finskom-zalive")
	if err != nil {
		t.Errorf("Failed to create file %v: %v", file, err)
	}
	defer file.Close()

	if _, err = os.Stat("www.marinohall.com/peterhof/"); os.IsNotExist(err) {
		t.Errorf("Directory was not created at path: www.marinohall.com/peterhof/")
	}

	if _, err = os.Stat("www.marinohall.com/peterhof/banketnyj-zal-na-finskom-zalive.html"); os.IsNotExist(err) {
		t.Errorf("File was not created at path: ")
	}
}

func TestCreateFileAndDirForRootHTML(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	file, err := creator.CreateFile("https://www.marinohall.com/")
	if err != nil {
		t.Errorf("Failed to create file %v: %v", file, err)
	}
	defer file.Close()

	if _, err = os.Stat("www.marinohall.com/"); os.IsNotExist(err) {
		t.Errorf("Directory was not created at path: www.marinohall.com/")
	}

	if _, err = os.Stat("www.marinohall.com/index.html"); os.IsNotExist(err) {
		t.Errorf("File was not created at path: www.marinohall.com/index.html")
	}
}

func TestCreateFileAndDirWithAnotherExt(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	file, err := creator.CreateFile("https://www.marinohall.com/aboba/jojo.css")
	if err != nil {
		t.Errorf("Failed to create file %v: %v", file, err)
	}
	defer file.Close()

	if _, err = os.Stat("www.marinohall.com/aboba"); os.IsNotExist(err) {
		t.Errorf("Directory was not created at path: www.marinohall.com/")
	}

	if _, err = os.Stat("www.marinohall.com/aboba/jojo.css"); os.IsNotExist(err) {
		t.Errorf("File was not created at path: www.marinohall.com/aboba/jojo.css")
	}
}

// func TestDownloadPageAndExtractLinks(t *testing.T) {
func TestCrawlFileCreate(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	ts := newTestServer()
	defer ts.Close()

	var visited = make(map[string]bool)
	app.Crawl(ts.URL+"/", ts.URL, &visited)

	if _, err := os.Stat(strings.TrimPrefix(ts.URL, "http://") + "/index.html"); os.IsNotExist(err) {
		t.Errorf("File was not created at path: %v/index.html", strings.TrimPrefix(ts.URL, "http://"))
	}

	if _, err := os.Stat(strings.TrimPrefix(ts.URL, "http://") + "/link2.html"); os.IsNotExist(err) {
		t.Errorf("File was not created at path: %v/link2.html", strings.TrimPrefix(ts.URL, "http://"))
	}

	if _, err := os.Stat(strings.TrimPrefix(ts.URL, "http://") + "/link1.css"); os.IsNotExist(err) {
		t.Errorf("File was not created at path: %v/link1.css", strings.TrimPrefix(ts.URL, "http://"))
	}

	if _, err := os.Stat(strings.TrimPrefix(ts.URL, "http://") + "/gojo/jojo.html"); os.IsNotExist(err) {
		t.Errorf("File was not created at path: %v/gojo/jojo.html", strings.TrimPrefix(ts.URL, "http://"))
	}
}

func TestCrawlCheckPathFileInFilesAndLinks(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	ts := newTestServer()
	defer ts.Close()

	var visited = make(map[string]bool)
	app.Crawl(ts.URL+"/", ts.URL, &visited)

	if !visited[ts.URL+"/"] && !visited[ts.URL+"/link1.css"] && !visited[ts.URL+"/gojo/jojo"] && !visited[ts.URL+"/link2"] {
		t.Errorf("Not all links was visited")
	}

	content, err := os.ReadFile(strings.TrimPrefix(ts.URL, "http://") + "/gojo/jojo.html")
	if err != nil {
		t.Errorf("Failed to read file %v: %v", strings.TrimPrefix(ts.URL, "http://")+"/gojo/jojo.html", err)
	}

	if !strings.Contains(string(content), "\"../link1.css\"") && !strings.Contains(string(content), "\"../link2.html\"") {
		t.Errorf("Not corrected path in file")
	}

	content, err = os.ReadFile(strings.TrimPrefix(ts.URL, "http://") + "/index.html")
	if err != nil {
		t.Errorf("Failed to read file %v: %v", strings.TrimPrefix(ts.URL, "http://")+"/index.html", err)
	}

	if !strings.Contains(string(content), "\"link1.css\"") && !strings.Contains(string(content), "\"gojo/jojo.html\"") {
		t.Errorf("Not corrected path in file")
	}
}

func TestExtractLinksFromCSS(t *testing.T) {
	links := make(map[string]bool)
	expectedLinks := map[string]bool{
		"'aboba/jojo.css'":   true,
		"\"foo.css":          true,
		"souja/boy/gege.css": true,
	}
	bufStr := bytes.NewBufferString("url('aboba/jojo.css') url(\"foo.css#bar\") lol kek url(\"data:\"randomdata)" +
		"url(souja/boy/gege.css)")

	ex.ExtractLinksFromCSS(&links, bufStr)

	if len(expectedLinks) == len(links) {
		for link := range links {
			if !expectedLinks[link] {
				t.Errorf("Link was not correctly extracted from CSS: %v", link)
			}
		}
	} else {
		t.Errorf("Not the expected number of links: %v, got: %v", len(expectedLinks), len(links))
	}
}

func TestExtractHTMLLinks(t *testing.T) {
	links := make(map[string]bool)
	expectedLinks := map[string]bool{
		"/styles/main.css": true,
		"/images/logo.png": true,
		"/page":            true,
		"/images/bg.jpg":   true,
	}
	htmlDoc := `
		<html>
			<head>
				<link rel="stylesheet" href="/styles/main.css">
			</head>
			<body>
				<img src="/images/logo.png">
				<a href="/page">Go to page</a>
				<div style="background-image: url('/images/bg.jpg')"></div>
			</body>
		</html>
	`

	doc, err := html.Parse(strings.NewReader(htmlDoc))
	if err != nil {
		t.Errorf("Failed to parse HTML: %v", err)
	}

	ex.CorrectAndExtractHTMLLinks("www.example.com/gojo/jojo", doc, &links)

	if len(expectedLinks) == len(links) {
		for link := range links {
			if !expectedLinks[link] {
				t.Errorf("Link was not correctly extracted from CSS: %v", link)
			}
		}
	} else {
		t.Errorf("Not the expected number of links: %v, got: %v", len(expectedLinks), len(links))
	}
}

func TestCorrectHTMLLinks(t *testing.T) {
	links := make(map[string]bool)
	htmlDoc := `
		<html>
			<head>
				<link rel="stylesheet" href="/styles/main.css">
			</head>
			<body>
				<img src="/images/logo.png">
				<a href="/page">Go to page</a>
				<div style="background-image: url('/images/bg.jpg')"></div>
			</body>
		</html>
	`

	doc, err := html.Parse(strings.NewReader(htmlDoc))
	if err != nil {
		t.Errorf("Failed to parse HTML: %v", err)
	}

	var buf bytes.Buffer

	ex.CorrectAndExtractHTMLLinks("www.example.com/gojo/jojo", doc, &links)
	err = html.Render(&buf, doc)
	if err != nil {
		t.Errorf("Failed to render HTML: %v", err)
	}

	if !strings.Contains(buf.String(), "\"../styles/main.css\"") ||
		!strings.Contains(buf.String(), "\"../images/logo.png\"") ||
		!strings.Contains(buf.String(), "\"../page.html\"") ||
		!strings.Contains(buf.String(), "(../images/bg.jpg)") {
		t.Errorf("Not corrected path in file")
	}

	buf.Reset()

	doc, err = html.Parse(strings.NewReader(htmlDoc))
	if err != nil {
		t.Errorf("Failed to parse HTML: %v", err)
	}

	ex.CorrectAndExtractHTMLLinks("www.example.com/", doc, &links)
	err = html.Render(&buf, doc)
	if err != nil {
		t.Errorf("Failed to render HTML: %v", err)
	}

	if !strings.Contains(buf.String(), "\"styles/main.css\"") ||
		!strings.Contains(buf.String(), "\"images/logo.png\"") ||
		!strings.Contains(buf.String(), "\"page.html\"") ||
		!strings.Contains(buf.String(), "(images/bg.jpg)") {
		t.Errorf("Not corrected path in file")
	}
}

func TestDownloadAndExtractLinks(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	ts := newTestServer()
	links, err := downloader.DownloadAndExtractLinks(ts.URL + "/")
	if err != nil {
		t.Errorf("Failed to download page and extract links: %v", err)
	}

	if len(links) == 2 {
		if !links["/gojo/jojo"] && !links["/link1.css"] {
			t.Errorf("Link was not correctly extracted from page: %v", links)
		}
	} else {
		t.Errorf("Not the expected number of links: 2, got: %v", len(links))
	}
	if _, err = os.Stat(strings.TrimPrefix(ts.URL, "http://") + "/index.html"); os.IsNotExist(err) {
		t.Errorf("File was not created at path: %v/index.html", strings.TrimPrefix(ts.URL, "http://"))
	}
}

func TestGetFullPath(t *testing.T) {
	fullPath, err := utils.GetFullPath("https://www.marinohall.com/ambar")
	if err != nil {
		t.Errorf("Failed to get full path: %v", err)
	}

	if fullPath != "www.marinohall.com/ambar" {
		t.Errorf("Full path was not correctly expected: www.marinohall.com/ambar, got: %v", fullPath)
	}
}

func TestGetPathToRoot(t *testing.T) {
	rootPath, err := utils.GetPathToRoot("https://www.marinohall.com/ambar")
	if err != nil {
		t.Errorf("Failed to get path to root: %v", err)
	}

	if rootPath != "" {
		t.Errorf("Root path was not correctly expected: \"\", got: %v", rootPath)
	}

	rootPath, err = utils.GetPathToRoot("https://www.marinohall.com/ambar/jojo")
	if err != nil {
		t.Errorf("Failed to get path to root: %v", err)
	}

	if rootPath != "../" {
		t.Errorf("Root path was not correctly expected: ../, got: %v", rootPath)
	}
}

func TestSetExtHTML(t *testing.T) {
	url := "https://www.marinohall.com/ambar"

	if filepath := utils.SetExtHTML(url); filepath != url+".html" {
		t.Errorf("SetExtHTML returned wrong filename: %v", filepath)
	}

	if filepath := utils.SetExtHTML(url + ".css"); filepath != url+".css" {
		t.Errorf("SetExtHTML returned wrong filename: %v", filepath)
	}

	if filepath := utils.SetExtHTML(url + ".png"); filepath != url+".png" {
		t.Errorf("SetExtHTML returned wrong filename: %v", filepath)
	}
}

func TestGetFileName(t *testing.T) {
	if fileName := utils.GetFileName("root/ambar"); fileName != "root/ambar.html" {
		t.Errorf("GetFileName returned wrong filename: %v, expected: root/ambar.html", fileName)
	}

	if fileName := utils.GetFileName("root/"); fileName != "root/index.html" {
		t.Errorf("GetFileName returned wrong filename: %v, expected: root/index.html", fileName)
	}

	if fileName := utils.GetFileName("root/ambar.css"); fileName != "root/ambar.css" {
		t.Errorf("GetFileName returned wrong filename: %v, expected: root/ambar.css", fileName)
	}

	if fileName := utils.GetFileName("root/ambar.png"); fileName != "root/ambar.png" {
		t.Errorf("GetFileName returned wrong filename: %v, expected: root/ambar.png", fileName)
	}
}

func TestNormalizeURL(t *testing.T) {
	if url := utils.NormalizeURL("/"); url != "index.html" {
		t.Errorf("NormalizeURL returned wrong filename: %v, expected: index.html", url)
	}

	if url := utils.NormalizeURL("/jojo.css"); url != "jojo.css" {
		t.Errorf("NormalizeURL returned wrong filename: %v, expected: jojo.css", url)
	}

	if url := utils.NormalizeURL("/foobar"); url != "foobar.html" {
		t.Errorf("NormalizeURL returned wrong filename: %v, expected: foobar.html", url)
	}
}

func TestConvertRelativeURLToAbsolute(t *testing.T) {
	expected := "https://example.com/jojo"

	absPath, err := utils.ConvertRelativeURLToAbsolute("https://example.com/", "/jojo")
	if err != nil {
		t.Errorf("Failed to get absolute path: %v", err)
	}

	if absPath != expected {
		t.Errorf("Absolute path was not correctly expected: %v, got: %v", expected, absPath)
	}

	absPath, err = utils.ConvertRelativeURLToAbsolute("https://example.com/foo/bar", "/jojo")
	if err != nil {
		t.Errorf("Failed to get absolute path: %v", err)
	}

	if absPath != expected {
		t.Errorf("Absolute path was not correctly expected: %v, got: %v", expected, absPath)
	}
}
