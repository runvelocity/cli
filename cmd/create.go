/*
Copyright © 2023 Utibeabasi Umanah utibeabasiumanah6@gmail.com
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rs/xid"
	"github.com/runvelocity/cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	functionName string
	filePath     string
	handler      string
)

func check(err error) {
	if err != nil {
		log.Fatalln(text.BgRed.Sprintf(err.Error()))
	}
}

func readBody(body io.ReadCloser) map[string]interface{} {
	var resp map[string]interface{}
	defer body.Close()
	b, err := io.ReadAll(body)
	check(err)
	err = json.Unmarshal(b, &resp)
	check(err)
	return resp
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new velocity function",
	Run: func(cmd *cobra.Command, args []string) {
		managerUrl := viper.Get("managerurl").(string)

		functionId := xid.New().String()

		file, err := os.Open(filePath)
		check(err)
		values := map[string]io.Reader{
			"code": file,
			"key":  strings.NewReader(functionId),
		}
		var uploadBuffer bytes.Buffer
		w := multipart.NewWriter(&uploadBuffer)
		err = utils.WriteBytes(values, w)
		check(err)
		uploadRequest, err := http.NewRequest("POST", managerUrl+"/upload", &uploadBuffer)
		check(err)
		uploadRequest.Header.Set("Content-Type", w.FormDataContentType())
		res, err := http.DefaultClient.Do(uploadRequest)
		check(err)
		// Check the response
		uploadResponse := readBody(res.Body)

		if res.StatusCode != http.StatusOK {
			log.Fatalln(text.BgRed.Sprintf("Error occured while uploading code: %s", uploadResponse["message"]))
		}

		// Create the function
		createFunctionArgs := map[string]io.Reader{
			"name":         strings.NewReader(functionName),
			"codeLocation": strings.NewReader(uploadResponse["Location"].(string)),
			"handler":      strings.NewReader(handler),
			"uuid":         strings.NewReader(functionId),
		}

		var createBuffer bytes.Buffer
		w = multipart.NewWriter(&createBuffer)
		err = utils.WriteBytes(createFunctionArgs, w)
		check(err)

		createRequest, err := http.NewRequest("POST", managerUrl+"/functions", &createBuffer)
		check(err)
		createRequest.Header.Set("Content-Type", w.FormDataContentType())
		createResponse, err := http.DefaultClient.Do(createRequest)
		check(err)
		resBody := readBody(createResponse.Body)
		// Check the response
		if createResponse.StatusCode != http.StatusCreated {
			log.Fatalln(text.BgRed.Sprintf("Error occured while creating function: %s", resBody["message"]))
		}

		log.Println(text.FgGreen.Sprintf("Created function %s", resBody["name"]))
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&functionName, "name", "", "Function name")
	createCmd.Flags().StringVar(&filePath, "file-path", "", "Path to the zip file containing the function code")
	createCmd.Flags().StringVar(&handler, "handler", "", "Function handler")
	err := createCmd.MarkFlagRequired("name")
	check(err)
	err = createCmd.MarkFlagRequired("file-path")
	check(err)
	err = createCmd.MarkFlagRequired("handler")
	check(err)
}
