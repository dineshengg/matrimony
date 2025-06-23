package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"

	routing "github.com/qiangxue/fasthttp-routing"
	log "github.com/sirupsen/logrus"
	viper "github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

var pageCache = make(map[string][]byte)

func IsRedirectInitialized() bool {
	// Check cache first
	pageCache["home/home.html"] = nil
	pageCache["login/login.html"] = nil
	pageCache["login/newuser.hmtl"] = nil
	pageCache["profile/profile.html"] = nil
	pageCache["contact/contact.html"] = nil
	pageCache["profile/account.html"] = nil

	for key, _ := range pageCache {
		data, status := inCache(key)
		if status != true {
			log.Printf("Error loading page %s:\n", key)
			delete(pageCache, key) // Remove from cache if not found
		} else {
			pageCache[key] = data // Cache the loaded page
		}
		if len(data) == 0 {
			fmt.Printf("Page %s is empty, removing from cache\n", key)
			delete(pageCache, key) // Remove empty pages from cache
		}
	}
	return len(pageCache) > 0
}

// only redirect the url will not load the page dynamically
func Redirect(ctx *routing.Context, url string) error {
	if len(url) > 0 {
		ctx.Redirect(url, fasthttp.StatusSeeOther)
	} else {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.WriteString(fmt.Sprintf("Invalid URL: %s", url))
	}
	return nil
}

func LoadPageWithValues(ctx *routing.Context, url string, data interface{}) error {
	// Check if the page is in cache
	page, found := inCache(url)
	if !found {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.WriteString(fmt.Sprintf("Page not found: %s", url))
		return fmt.Errorf("page not found: %s", url)
	}

	//use html template to fill the data in the page
	if data != nil {
		log.Println(page)
		return nil
	}

	return nil
}

// RenderTemplatePage loads the HTML page, parses it as a template, and executes it with the provided data.
// Example: RenderTemplatePage(ctx, "login/newuser.html", map[string]string{"Name": "User"})
func RenderTemplatePage(ctx *routing.Context, url string, data interface{}) ([]byte, error) {
	// pageBytes, found := inCache(url)
	// if !found || len(pageBytes) == 0 {
	// 	ctx.SetStatusCode(fasthttp.StatusNotFound)
	// 	ctx.WriteString(fmt.Sprintf("Page not found: %s", url))
	// 	return fmt.Errorf("page not found: %s", url)
	// }
	pagebytes, err := LoadPage(url)
	if err != nil {
		log.Errorf("Error loading page %s: %v", url, err)
		return nil, fmt.Errorf("error loading page %s: %v", url, err)
	}

	tmpl, err := template.New(url).Parse(string(pagebytes))
	if err != nil {
		log.Errorf("Template parse error for %s: %v", url, err)
		return nil, fmt.Errorf("template parse error for %s: %v", url, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Errorf("Template execution error for %s: %v", url, err)
		return nil, fmt.Errorf("template execution error for %s: %v", url, err)
	}
	return buf.Bytes(), nil
}

func inCache(url string) ([]byte, bool) {
	//cache the page in memory so that when we need to redirect will take from cache instead of file system
	page, found := pageCache[url]
	if found {
		return page, true
	}
	// If not found in cache, try to load from file system
	data, err := LoadPage(url)
	if err != nil {
		return nil, false
	}
	pageCache[url] = data
	// Return the loaded page
	return data, true
}

func LoadPage(url string) ([]byte, error) {
	//load the page from file system or cache
	//if not found in cache then load from file system and cache it
	//if not found in file system then return error
	// Read file from file system
	relative_path := viper.GetString("resources.filepath")
	url = fmt.Sprintf("%s/%s", relative_path, url)
	log.Debug(url)
	data, err := os.ReadFile(url)
	if err != nil {
		return nil, err
	}
	return data, nil
}
