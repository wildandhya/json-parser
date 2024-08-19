package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
)

func main() {
	app := kingpin.New("jp", "Simple JSON Parser")

	fileCommand := app.Command("file", "Filename")
	fileCommandArg := fileCommand.Arg("filename", "file path").String()
	fileCommand.Action(func(pc *kingpin.ParseContext) error {
		file := *fileCommandArg
		fileByte, err := os.ReadFile(file)
		if err != nil{
			return err
		}

		lexer := NewLexer(string(fileByte))
		tokens, err := lexer.GetTokens()
		if err != nil{
			return err
		}

		parser := NewParser(tokens)
		result, err := parser.Parse()
		if err != nil{
			log.Printf("%v", err)
			os.Exit(1)
		}
		log.Println(parser.toJSONString(result))
		return nil
	})

	_, err := app.Parse(os.Args[1:])
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
