package main

// import (
// 	"fmt"
// 	"log"
// 	"time"
// )

// func test() {
// 	log.Println("Starting GoSearch")
// 	log.Println("Creating indexmap...")
// 	app := initApp("test app")
// 	log.Println("")
// 	input, _ := parseArbJSON(`{"example": "hey oi", "ho": "yo hey hey hey hey"}`)
// 	input2, _ := parseArbJSON(`{"example": "hey no", "ho": "yo hey hey hey hey"}`)
// 	app.addIndex(input)
// 	app.addIndex(input2)
// 	log.Println(app.indexes)
// 	log.Println("### SEARCHING...")
// 	start := time.Now()
// 	search, _ := app.search("oi", make([]string, 0))
// 	end := time.Now()
// 	log.Println("### SEARCH RESULT:", search)
// 	log.Println("### SEARCH TIME:", end.Sub(start))

// 	// log.Println("### TESTING WITH LARGE DATA SET...")
// 	// app = initApp("Movie Data")
// 	// jsonFile, err := os.Open("./moviedata.json")
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// }
// 	// log.Println("### Successfully Opened users.json...")
// 	// defer jsonFile.Close()
// 	// byteValue, _ := ioutil.ReadAll(jsonFile)

// 	// var result []map[string]interface{}
// 	// json.Unmarshal([]byte(byteValue), &result)
// 	// start = time.Now()
// 	// for num, item := range result {
// 	// 	id := app.addIndex(item)
// 	// 	if num == 1701 {
// 	// 		log.Println("doc", id)
// 	// 	}
// 	// }
// 	// end = time.Now()
// 	// log.Println("### Added", len(result), "records in", end.Sub(start))

// 	// start = time.Now()
// 	// search = app.search("Deux", make([]string, 0))
// 	// end = time.Now()
// 	// log.Println(search)
// 	// log.Println("### All fields search done in", end.Sub(start))
// 	// fields := make([]string, 1)
// 	// fields[0] = "year" // notice how it auto casts the int into a string
// 	// start = time.Now()
// 	// search = app.search("19938", fields)
// 	// end = time.Now()
// 	// log.Println(search)
// 	// log.Println("### Specific fields search done in", end.Sub(start))
// 	// burnfields := make([]string, 3)
// 	// burnfields[0] = "year"
// 	// burnfields[1] = "info"
// 	// burnfields[2] = "title"
// 	// start = time.Now()
// 	// search = app.search("Deux", burnfields)
// 	// end = time.Now()
// 	// log.Println(search)
// 	// log.Println("### Many fields search done in", end.Sub(start))

// 	log.Println("### BUILDING ENORMOUS DATA SET")
// 	app = initApp("Big Set")
// 	start = time.Now()
// 	for i := 0; i < 1000000; i++ {
// 		input, _ = parseArbJSON(fmt.Sprintf(`{"example": "hey%d", "no": "ho%d", "te": "to%d", "fe": "fo%d", "re": "re%d", "to": "bo%d", "aa": "aa%d", "bb": "bb%d", "cc": "cc%d", "ee": "ee%d"}`, i, i, i, i, i, i, i, i, i, i))
// 		app.addIndex(input)
// 	}
// 	end = time.Now()
// 	log.Println("### BIG CONSTRUCTION:", end.Sub(start))

// 	log.Println("### SEARCHING...")
// 	start = time.Now()
// 	search, _ = app.search("ho5000", make([]string, 0))
// 	end = time.Now()
// 	log.Println(search)
// 	log.Println("### SEARCH DONE IN", end.Sub(start))

// 	log.Println("### SEARCHING...")
// 	start = time.Now()
// 	search, _ = app.search("cc500000", make([]string, 0))
// 	end = time.Now()
// 	log.Println(search)
// 	log.Println("### SEARCH DONE IN", end.Sub(start))

// 	log.Println("### SEARCHING...")
// 	start = time.Now()
// 	search, _ = app.search("ee999900", make([]string, 0))
// 	end = time.Now()
// 	log.Println(search)
// 	log.Println("### SEARCH DONE IN", end.Sub(start))

// 	log.Println("### SPECIFIC SEARCH")
// 	fields := make([]string, 1)
// 	fields[0] = "example"
// 	start = time.Now()
// 	search, _ = app.search("hey9000", fields)
// 	end = time.Now()
// 	log.Println(search)
// 	log.Println("### SPECIFIC SEARCH DONE IN", end.Sub(start))

// 	log.Println("### SPECIFIC SEARCH")
// 	fields = make([]string, 1)
// 	fields[0] = "cc"
// 	start = time.Now()
// 	search, _ = app.search("cc500000", fields)
// 	end = time.Now()
// 	log.Println(search)
// 	log.Println("### SPECIFIC SEARCH DONE IN", end.Sub(start))

// 	log.Println("### SPECIFIC SEARCH")
// 	fields = make([]string, 1)
// 	fields[0] = "ee"
// 	start = time.Now()
// 	search, _ = app.search("ee999900", fields)
// 	end = time.Now()
// 	log.Println(search)
// 	log.Println("### SPECIFIC SEARCH DONE IN", end.Sub(start))
// }
