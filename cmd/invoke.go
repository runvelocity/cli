/*
Copyright © 2023 Utibeabasi Umanah utibeabasiumanah6@gmail.com
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var functionNameToInvoke string

// invokeCmd represents the invoke command
var invokeCmd = &cobra.Command{
	Use:   "invoke",
	Short: "Invoke a function",
	Run: func(cmd *cobra.Command, args []string) {
		payload := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			if err != nil {
				log.Fatalln("A fatal error occured: " + err.Error())
			}
		}
		invokeRequest, err := http.NewRequest("POST", managerUrl+"/invoke/"+functionNameToInvoke, bytes.NewBuffer(jsonPayload))
		invokeRequest.Header.Set("Content-type", "application/json")
		if err != nil {
			log.Fatalln("A fatal error occured: " + err.Error())
		}
		log.Println(text.FgGreen.Sprintf("Invoking function %s", functionNameToInvoke))
		invokeResponse, err := http.DefaultClient.Do(invokeRequest)
		if err != nil {
			log.Fatalln("A fatal error occured: " + err.Error())
		}
		resBody := readBody(invokeResponse.Body)
		// Check the response
		if invokeResponse.StatusCode != http.StatusOK {
			log.Fatalln(text.BgRed.Sprintf("Error occured while invoking function: %s", resBody["message"]))
		} else {
			log.Println(text.FgGreen.Sprintf("%v", resBody["result"]))
		}
	},
}

func init() {
	rootCmd.AddCommand(invokeCmd)
	invokeCmd.Flags().StringVar(&functionNameToInvoke, "name", "", "Function to invoke")
	err := invokeCmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatalln("A fatal error occured: " + err.Error())
	}
}
