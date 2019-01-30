package gitpage

import (
	"fmt"
	"io/ioutil"

	"github.com/microcosm-cc/bluemonday"

	"github.com/russross/blackfriday"
)

// Tree interface describes the interfaces that must be implemented
type Tree interface {
	Root() string
}

// Get retrieves the page
func Get(git Tree, page string) string {
	path := git.Root()
	markdownPage := fmt.Sprintf("%s/%s.md", path, page)

	file, _ := ioutil.ReadFile(markdownPage)
	unsafe := blackfriday.Run(file)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	return string(html)
}
