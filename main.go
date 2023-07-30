package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Collection struct {
	Items map[string]Item `json:"items"`
}

type Item struct {
	Identifier string                 `json:"identifier"`
	Properties map[string]interface{} `json:"properties"`
	CreatedAt  int64                  `json:"created_at"`
	UpdatedAt  int64                  `json:"updated_at"`
}

type Store struct {
	Collections map[string]Collection
}

type Log struct {
	Collection string `json:"collection"`
	Identifier string `json:"identifier"`
	Old        Item   `json:"old"`
	New        Item   `json:"new"`
	Timestamp  int64  `json:"timestamp"`
}

func CreateOrAppendLog(l Log) {
	f, err := os.OpenFile("./logs.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	jsonData, err := json.Marshal(l)
	if err != nil {
		log.Fatal(err)
	}

	jsonData = append(jsonData, []byte("\n")...)

	f.Write(jsonData)
}

func (s *Store) CreateCollection(collection string) {
	s.Collections[collection] = Collection{
		Items: make(map[string]Item),
	}
}

func (s *Store) CreateOrUpdateItem(collection string, identifier string, properties map[string]interface{}) {

	if _, ok := s.Collections[collection]; !ok {
		s.CreateCollection(collection)
	}

	l := Log{
		Collection: collection,
		Identifier: identifier,
		Timestamp:  time.Now().Unix(),
	}

	if _, ok := s.Collections[collection].Items[identifier]; !ok {
		s.Collections[collection].Items[identifier] = Item{
			Identifier: identifier,
			Properties: properties,
			CreatedAt:  time.Now().Unix(),
			UpdatedAt:  time.Now().Unix(),
		}
	} else {
		l.Old = s.Collections[collection].Items[identifier]
		s.Collections[collection].Items[identifier] = Item{
			Identifier: identifier,
			Properties: properties,
			CreatedAt:  s.Collections[collection].Items[identifier].CreatedAt,
			UpdatedAt:  time.Now().Unix(),
		}
	}

	l.New = s.Collections[collection].Items[identifier]

	CreateOrAppendLog(l)

	storeToDir(s, "./collection")
}

func fileToCollection(dir string, file os.FileInfo) Collection {
	// read json

	f, err := os.Open(dir + "/" + file.Name())
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	var collection Collection
	json.Unmarshal(data, &collection)

	return collection
}

func dirToStore(dir string) *Store {
	d, err := os.Open(dir)
	if err != nil {
		// if dir not exist
		if os.IsNotExist(err) {
			os.Mkdir(dir, 0755)
			return &Store{
				Collections: make(map[string]Collection),
			}
		}

		log.Fatal(err)
	}

	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		return &Store{
			Collections: make(map[string]Collection),
		}
	}

	store := Store{
		Collections: make(map[string]Collection),
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			collection := fileToCollection(dir, file)
			store.Collections[strings.Replace(file.Name(), ".json", "", 1)] = collection
		}
	}

	return &store
}

func storeToDir(store *Store, dir string) {
	for collection, data := range store.Collections {
		f, err := os.Create(dir + "/" + collection + ".json")
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		jsonData, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			log.Fatal(err)
		}

		f.Write(jsonData)
	}
}

func main() {

	// fiber instance
	app := fiber.New(
		fiber.Config{
			Concurrency: 1,
		},
	)

	// routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	// put collection
	app.Put("/collection/:collection/:identifier", func(c *fiber.Ctx) error {

		// store instance
		store := dirToStore("./collection")

		collection := c.Params("collection")
		identifier := c.Params("identifier")

		payload := make(map[string]interface{})
		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		store.CreateOrUpdateItem(collection, identifier, payload)
		return c.JSON(payload)
	})

	app.Get("*.html", func(c *fiber.Ctx) error {
		path := c.Params("*")

		paths := []string{`/collection\/[a-z]\/i`, "/logs/i"}

		// regex match
		for _, p := range paths {
			if match, _ := regexp.MatchString(p, path); match {
				c.Send([]byte("path not exist"))
				return c.SendStatus(404)
			}
		}

		// return file
		f := "./" + path + ".json"

		if path == "logs" {
			f = "./logs.jsonl"
		}

		if _, err := os.Stat(f); os.IsNotExist(err) {
			c.Send([]byte("f not existos"))
			return c.SendStatus(404)
		}

		file, err := os.Open(f)
		if err != nil {
			c.Send([]byte("f cannot open"))

			return c.SendStatus(500)

		}

		defer file.Close()

		// read file
		data, err := ioutil.ReadAll(file)
		if err != nil {
			c.Send([]byte("f cannot read"))
			return c.SendStatus(500)
		}

		c.Type("html")

		if path == "logs" {
			return c.Send([]byte(UI(string(jsonLToJSON(data)))))
		}

		return c.Send([]byte(UI(string(data))))
	})

	// get
	app.Get("/*", func(c *fiber.Ctx) error {
		path := c.Params("*")

		paths := []string{`/collection\/[a-z]\.json/i`, `/logs.jsonl/i`}

		// regex match
		for _, p := range paths {
			if match, _ := regexp.MatchString(p, path); match {
				c.Send([]byte("f not existr"))
				return c.SendStatus(404)
			}
		}

		// return file
		f := "./" + path

		if _, err := os.Stat(f); os.IsNotExist(err) {
			c.Send([]byte("f not exist"))
			return c.SendStatus(404)
		}

		file, err := os.Open(f)
		if err != nil {
			c.Send([]byte("f cannot open"))
			return c.SendStatus(500)
		}

		defer file.Close()

		// read file
		b, err := ioutil.ReadAll(file)
		if err != nil {
			c.Send([]byte("f cannot read"))
			return c.SendStatus(500)
		}

		return c.Send(b)

	})

	// listen on port 3000
	app.Listen(":3000")
}

func UI(json string) string {
	return `<iframe id="jsoncrackEmbed" src="https://jsoncrack.com/widget"></iframe>

	<script>
	const jsonCrackEmbed = document.querySelector("iframe");
	
	const json = JSON.stringify(` + json + `);
	
	window?.addEventListener("message", (event) => {
	jsonCrackEmbed.contentWindow.postMessage({
		json
	}, "*");
	});
	</script>
	
	<style>
	body {
	margin: 0;
	padding: 0;
	}
	
	section {
	width: 100%;
	height: 100vh;
	display: flex;
	flex-direction: column;
	}
	
	textarea {
	width: 100%;
	height: 100%;
	}
	
	div {
	display: flex;
	width: 100%;
	height: 150px;
	}
	
	#jsoncrackEmbed {
	flex: 1;
	order: 2;
	border: none;
	width: 100%;
	height: 100vh;
	}
	</style>
	`
}

func jsonLToJSON(jsonl []byte) []byte {
	sl := strings.Split(string(jsonl), "\n")
	return []byte("[" + strings.Join(sl, ",") + "]")
}
