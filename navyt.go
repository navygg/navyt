package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

// Config config.toml
type Config struct {
	SdkID     string
	HasAppUID bool
	Host      string
	Port      int
	LogDir    string
	LogFile   string
}

// THandler implements http.Handler
type THandler struct {
	logger *log.Logger
}

// ServeHTTP response req
func (t *THandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	p := r.Form["name"][0]
	t.logger.Printf("%s, %s", host, p)

	rsp, err := http.Get("https://api.mch.weixin.qq.com/pay/orderquery")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("http.Get() err: %v", err)))
		return
	}
	defer rsp.Body.Close()

	_, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.logger.Printf("ioutil.ReadAll() err: %v", err)
	}

	w.Write([]byte(fmt.Sprintf("got: %s, status: %d", p, rsp.StatusCode)))
}

func start(fileName string) error {
	fmt.Println(fileName)
	config := &Config{}
	if _, err := toml.DecodeFile(fileName, config); err != nil {
		fmt.Printf("config file parse error: %v\n", err)
		os.Exit(-1)
	}

	var err error
	config.LogDir, err = filepath.Abs(path.Clean(config.LogDir))
	if err != nil {
		fmt.Printf("abs error: %v\n", err)
		os.Exit(-1)
	}

	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		fmt.Printf("mkdirall error: %v\n", err)
		os.Exit(-1)
	}

	if config.LogFile == "" {
		fmt.Println("LogFile is nil")
		os.Exit(-1)
	}

	logFileName := path.Join(config.LogDir, path.Clean(config.LogFile))
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("open %s error: %v\n", logFileName, err)
		os.Exit(-1)
	}
	defer logFile.Close()

	logger := log.New(logFile, config.SdkID+" ", log.LstdFlags|log.Lshortfile)
	h := THandler{logger: logger}

	http.HandleFunc("/hi", h.ServeHTTP)
	server := http.Server{
		Addr: fmt.Sprintf("%v:%v", config.Host, config.Port),
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("start server error: %v\n", err)
		os.Exit(-1)
	}

	return nil
}

func main() {
	cmdRun := &cobra.Command{
		Use:   "run config.toml",
		Short: "run server",
		Long:  "run server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("need config file")
			}

			if err := start(args[0]); err != nil {
				fmt.Printf("start error: %+v\n", err)
			}
			return nil
		},
	}

	rootCmd := &cobra.Command{
		Use:  "navyt server",
		Long: "navyt server",
	}

	rootCmd.AddCommand(cmdRun)
	rootCmd.Execute()
}
