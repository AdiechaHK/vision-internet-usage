package main

import (
	"encoding/json"
	"fmt"
	"hello-heroku/data"
	"hello-heroku/schedule"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello word !")
}

func main() {

	var rootCmd = &cobra.Command{
		Use:   "hugo",
		Short: "Hugo is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
					  love by spf13 and friends in Go.
					  Complete documentation is available at http://hugo.spf13.com`,
		Run: func(cmd *cobra.Command, args []string) {
			// scheduleFlag := cmd.PersistentFlags().Bool("scheduler", false, "Use --scheduler to run scheduled command instead of serving web app.")
			isScheduler, err := cmd.Flags().GetBool("scheduler")
			if err != nil {
				panic(err)
			}
			fmt.Println("hello")
			if isScheduler {
				fmt.Println("Yes schedule")
				schedule.TestFun()
			} else {
				port := os.Getenv("PORT")
				http.HandleFunc("/data", func(w http.ResponseWriter, _ *http.Request) {
					lst := data.GetData()
					w.Header().Set("Content-Type", "application/json")
					if lst != nil {
						json.NewEncoder(w).Encode(lst)
					} else {
						json.NewEncoder(w).Encode(data.DataCollection{})
					}
				})
				http.Handle("/", http.FileServer(http.Dir("./static")))
				fmt.Println("Listening on port: " + port)
				http.ListenAndServe(":"+port, nil)
			}
		},
	}

	rootCmd.PersistentFlags().Bool("scheduler", false, "Use --scheduler to run scheduled command instead of serving web app.")

	rootCmd.Execute()
}
