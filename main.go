package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	// _ "github.com/mattn/go-sqlite3" // need cgo, bud support windows
	_ "modernc.org/sqlite" // cgo free, but don't support windows

	// color
	"github.com/gookit/color"
	"github.com/liuzl/gocc"
)

type Entry struct {
	Char  string
	Forms []string
}

//go:embed cangjie.db
var dbFile embed.FS

var (
	green = color.FgGreen.Render
	blue  = color.FgBlue.Render
	red   = color.FgRed.Render

	bgLightYellow = color.BgLightYellow.Render
)

const HelpInfo = `Usage:
  cj [query]		查询
  cj -r | --reinstall	重建数据库`

func init() {

	// if not exist cangjie.db
	// then create it
	if _, err := os.Stat(dbPath("cangjie.db")); os.IsNotExist(err) {
		err = installDB()
		if err != nil {
			panic(err)
		}
	}
}

func main() {

	querys := parseArgs(os.Args[1:])

	for _, query := range querys {
		entry := new(Entry)
		entry, err := doQuery(query)
		if err != nil {
			fmt.Println(err)
		}

		if len(entry.Forms) > 1 {
			if entry.Forms[0] == entry.Forms[1] {
				entry.Forms = entry.Forms[:1]
			}
		}

		for _, form := range entry.Forms {
			fmt.Printf(bgLightYellow(red(" > %s < ")), entry.Char)
			fmt.Print(green(" "))
			for _, char := range form {
				fmt.Print(green(atohan(string(char))))
			}
			fmt.Print(" ")
			for _, char := range form {
				fmt.Printf("%s", blue(string(char)))
			}

			fmt.Println()
		}
	}
}

func parseArgs(args []string) (querys []string) {
	var err error
	if len(args) == 0 {
		fmt.Println(HelpInfo)
		os.Exit(1)
	}

	if args[0] == "-r" || args[0] == "--reinstall" {
		err = reinstallDB()
		if err != nil {
			panic(err)
		}
		color.Redln("cangjie.db 重建成功")
		os.Exit(0)
	}

	for _, arg := range args {
		for _, query := range arg {
			querys = append(querys, string(query))
			querys = append(querys, convertToSimpleChinese(string(query)))
		}
	}

	querys = removeDuplicates(strings.Join(querys, ""))
	return

}

func removeDuplicates(input string) []string {
	charMap := make(map[rune]bool)
	result := []string{}

	for _, char := range input {
		if !charMap[char] {
			charMap[char] = true
			result = append(result, string(char))
		}
	}

	return result
}

func convertToSimpleChinese(arg string) string {
	t2s, err := gocc.New("t2s")
	if err != nil {
		log.Fatal(err)
	}

	simpleChinese, err := t2s.Convert(arg)
	if err != nil {
		fmt.Println(err)
	}

	return simpleChinese

}

func doQuery(query string) (*Entry, error) {
	db, err := sql.Open("sqlite", dbPath("cangjie.db"))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	q := fmt.Sprintf("select form from entry where char='%s'", query)
	row, err := db.Query(q)
	if err != nil {
		// fmt.Print("faild to query")
		return nil, err
	}

	defer row.Close()

	var form string
	for row.Next() {
		err = row.Scan(&form)
		if err != nil {
			return nil, err
		}
	}
	// fmt.Println(form)

	var result []string
	err = json.Unmarshal([]byte(form), &result)
	if err != nil {
		return nil, err
	}

	entry := Entry{
		Char:  query,
		Forms: result,
	}

	return &entry, nil

}

// Ascii to Hans
func atohan(s string) string {

	var result string

	table := map[string]string{
		"a": "日",
		"b": "月",
		"c": "金",
		"d": "木",
		"e": "水",
		"f": "火",
		"g": "土",
		"h": "竹",
		"i": "戈",
		"j": "十",
		"k": "大",
		"l": "中",
		"m": "一",
		"n": "弓",
		"o": "人",
		"p": "心",
		"q": "手",
		"r": "口",
		"s": "尸",
		"t": "廿",
		"u": "山",
		"v": "女",
		"w": "田",
		"x": "難",
		"y": "卜",
		"z": "符",
	}

	for _, char := range s {
		result += table[string(char)]
	}
	return result
}

// reinstallDB is a function that reinstall the database.
// It removes the database if it already exists,
// and then installs it from the embedded file.
//
// It returns an error if there is any error during the process.
func reinstallDB() error {
	if _, err := os.Stat(dbPath("cangjie.db")); os.IsNotExist(err) {
		err = os.Remove(dbPath("cangjie.db"))
		if err != nil {
			return err
		}
	}

	if err := installDB(); err != nil {
		return err
	}

	return nil
}

// installDB installs the database.
//
// It reads the contents of the "cangjie.db" file using the dbFile.ReadFile function
// and creates a new "cangjie.db" file using the os.Create function. It then writes the
// data read from the file to the new file using the cangjie.Write function. Finally,
// it calls the cangjie.Sync function to ensure that the data is written to disk,
// and closes the file using the defer statement. It returns an error if any of the
// file operations fail.
//
// The installDB function has no parameters and returns an error.
func installDB() error {

	data, err := dbFile.ReadFile("cangjie.db")
	if err != nil {
		return err
	}

	cangjie, err := os.Create(dbPath("cangjie.db"))
	if err != nil {
		return err
	}

	_, err = cangjie.Write(data)
	if err != nil {
		return err
	}

	err = cangjie.Sync()
	if err != nil {
		return err
	}
	defer func(cangjie *os.File) {
		err := cangjie.Close()
		if err != nil {
			panic("Fail to close file....")
		}
	}(cangjie)

	return nil
}

// dbPath returns the path to the SQLite database file.
//
// It takes a dbName string parameter, which specifies the name of the database.
// It returns a string, which is the path to the SQLite database file.
// It will choose the User Cache Dir to store the database file.
func dbPath(dbName string) string {
	// exe, err := os.Executable()
	// if err != nil {
	// 	panic(err)
	// }
	// exePath := filepath.Dir(exe)

	// sqlitePath := filepath.Join(exePath, dbName)
	// return sqlitePath

	// use User Cache Dir to store cangjie.db

	// if not exist %UserCacheDir%\cangjie folder
	// then create it

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(filepath.Join(cacheDir, "cangjie"))
	if os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(cacheDir, "cangjie"), os.ModePerm)
		if err != nil {
			panic("Fail to create cache dir")
		}
	}

	return filepath.Join(cacheDir, "cangjie", dbName)
}
