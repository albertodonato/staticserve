package main

import (
	"flag"
	"log"
	"os"
	"text/template"

	"github.com/albertodonato/h2static/server"
	"github.com/albertodonato/h2static/version"
)

const helpHeaderTemplate = `
{{.Name}} {{.Version}} - Tiny static web server with TLS and HTTP/2 support.

Usage of {{.Name}}:

`

// NewStaticServerFromCmdline returns a a StaticServer parsing cmdline args.
func NewStaticServerFromCmdline(fs *flag.FlagSet, args []string) (*server.StaticServer, error) {
	conf := &server.StaticServerConfig{}
	fs.StringVar(&conf.Addr, "addr", ":8080", "address and port to listen on")
	fs.StringVar(
		&conf.PasswordFile, "basic-auth", "",
		`password file for Basic Auth (each line should be in the form "user:SHA512-hash")`)
	fs.StringVar(&conf.Dir, "dir", ".", "directory to serve")
	fs.BoolVar(
		&conf.DisableLookupWithSuffix, "disable-lookup-with-suffix", false,
		"disable matching files with .htm(l) suffix for paths without suffix")
	fs.BoolVar(&conf.DisableH2, "disable-h2", false, "disable HTTP/2 support")
	fs.BoolVar(&conf.ShowDotFiles, "show-dotfiles", false, "show files whose name starts with a dot")
	fs.BoolVar(&conf.Log, "log", false, "log requests")
	fs.StringVar(&conf.TLSCert, "tls-cert", "", "certificate file for TLS connections")
	fs.StringVar(&conf.TLSKey, "tls-key", "", "key file for TLS connections")
	fs.Usage = func() {
		printHeader(fs)
		fs.PrintDefaults()
		fs.Output().Write([]byte{'\n'})
	}
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return server.NewStaticServer(*conf)
}

func printHeader(fs *flag.FlagSet) {
	tpl := template.Must(template.New("helpHeader").Parse(helpHeaderTemplate))
	if err := tpl.Execute(fs.Output(), version.App); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetPrefix(version.App.Name + " ")
	server, err := NewStaticServerFromCmdline(flag.CommandLine, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
