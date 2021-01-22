package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/sandro/lorca"
)

func main() {
	ui, _ := lorca.New("data:text/html,"+url.PathEscape(`
  <html>
  <head><title>Hello</title></head>
  <body>
    <input type="file" multiple>
    <script>
      window.files = []
      document.querySelector("input").addEventListener("change", async (e) => {
        let el = e.target
        for (file of el.files) {
				  window.files.push(file)
				}
        let paths = await filesAdded("files", window.files.length)
        console.log(paths)
      })
    </script>
  </body>
  </html>
`), "", 480, 320)
	defer ui.Close()

	ui.Bind("filesAdded", func(listVar string, length int) []string {
		paths := make([]string, length)
		for i := 0; i < length; i++ {
			tm, err := ui.EvalRaw(fmt.Sprintf("window.%s[%d]", listVar, i))
			if err != nil {
				log.Printf("error evaluating window.%s. err: %s\n", listVar, err)
			}
			path, err := lorca.GetFilePath(ui, tm.Result.Result.ObjectID)
			if err != nil {
				log.Println(err)
			} else {
				paths[i] = path
			}
		}
		return paths
	})

	// Wait for the browser window to be closed
	<-ui.Done()
}
