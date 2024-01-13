/*
Copyright © 2023 Utibeabasi Umanah utibeabasiumanah6@gmail.com
*/
package cmd

import (
	"log"
	"net/http"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var functionNameToDelete string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a function",
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

		for _, v := range resBody["functions"].([]interface{}) {
			v = v.(map[string]interface{})
			if v.(map[string]interface{})["name"] == functionNameToDelete {
				log.Println(text.FgGreen.Sprintf("Deleting function %s", functionNameToDelete))
				deleteRequest, err := http.NewRequest("DELETE", managerUrl+"/functions/"+v.(map[string]interface{})["uuid"].(string), nil)
				check(err)
				deleteResponse, err := http.DefaultClient.Do(deleteRequest)
				check(err)
				resBody := readBody(deleteResponse.Body)
				// Check the response
				if deleteResponse.StatusCode != http.StatusOK {
					log.Fatalln(text.BgRed.Sprintf("Error occured while deleting function: %s", resBody["message"]))
				}
				log.Println(text.FgGreen.Sprintf("Successfully deleted function %s", functionNameToDelete))
				return
			}
		}

		log.Fatalln(text.BgRed.Sprintf("Function not found"))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&functionNameToDelete, "name", "", "Function name to delete")
	err := deleteCmd.MarkFlagRequired("name")
	check(err)
}
