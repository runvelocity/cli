/*
Copyright © 2023 Utibeabasi Umanah utibeabasiumanah6@gmail.com
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Fetch all functions",
	Run: func(cmd *cobra.Command, args []string) {
		managerUrl := viper.Get("managerurl").(string)
		log.Println(text.FgGreen.Sprintf("Fetching functions"))
		getRequest, err := http.NewRequest("GET", managerUrl+"/functions", nil)
		check(err)
		getRequest.Header.Set("Content-Type", "application/json")
		getResponse, err := http.DefaultClient.Do(getRequest)
		check(err)
		resBody := readBody(getResponse.Body)
		// Check the response
		if getResponse.StatusCode != http.StatusOK {
			log.Fatalln(text.BgRed.Sprintf("Error occured while fetching function: %s", resBody["message"]))
		}

		var rows []table.Row
		for _, v := range resBody["functions"].([]interface{}) {
			v = v.(map[string]interface{})
			rows = append(rows, table.Row{v.(map[string]interface{})["name"], v.(map[string]interface{})["uuid"], v.(map[string]interface{})["status"]})
		}

		tableRowHeader := table.Row{"Name", "UUID", "STatus"}
		tw := table.NewWriter()
		tw.AppendHeader(tableRowHeader)
		tw.AppendRows(rows)
		fmt.Println(tw.Render())

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
