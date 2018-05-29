package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"strings"
)

var parallelRequest bool

func queryCheck(args []string) string {
	interfaceAndPassword()
	queryInit := []string{}
	queryFinal := ""
	if len(args) == 0 {
		fmt.Println("\033[31mError: Must set at least one query")
		os.Exit(0)
	} else if len(args) == 1 {
		queryFinal = url.PathEscape(args[0])
	} else {
		for i := 0; i < len(args); i++ {
			queryInit = append(queryInit, args[i])
		}
		queryFinal = url.PathEscape(strings.Join(queryInit, ";"))
	}
	return queryFinal
}

func runQuery(cmd *cobra.Command, args []string) error {

	if parallelRequest == true {
		printParallelQueries(args)
	}

	queryFinal := queryCheck(args)
	queryUrl := "queryResults?query=" + queryFinal
	bodyBytes := connect("GET", queryUrl, nil)
	results := normalizeQuery(bodyBytes)

	if len(results) == 1 {
		jsonData, _ := json.MarshalIndent(results[0], "", "  ")
		fmt.Fprintf(os.Stdout, "%s", jsonData)
	} else {
		jsonData, _ := json.MarshalIndent(results, "", "  ")
		fmt.Fprintf(os.Stdout, "%s", jsonData)
	}

	return nil
}

var query = &cobra.Command{
	Use:   "query",
	Short: "Runs one or more ROQL queries",
	Long:  "\033[93mRuns one or more ROQL queries and returns parsed results\033[0m \033[0;32m\n\nSingle Query Example: \033[0m \n$ osvc-rest query \"DESCRIBE\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE \033[0;32m\n\nMultiple Queries Example:\033[0m \n$ osvc-rest query \"SELECT * FROM INCIDENTS LIMIT 100\" \"SELECT * FROM SERVICEPRODUCTS LIMIT 100\" \"SELECT * FROM SERVICECATEGORIES LIMIT 100\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE \033[0;32m\n\nParallel Queries Example:\033[0m \n$ osvc-rest query \"SELECT * FROM INCIDENTS LIMIT 20000\" \"SELECT * FROM INCIDENTS Limit 20000 OFFSET 20000\" \"SELECT * FROM INCIDENTS Limit 20000 OFFSET 40000\" \"SELECT * FROM INCIDENTS Limit 20000 OFFSET 60000\" \"SELECT * FROM INCIDENTS Limit 20000 OFFSET 80000\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE",
	RunE:  runQuery,
}

func init() {
	query.Flags().BoolVarP(&parallelRequest, "parallel", "l", false, "Runs queries in parallel")
	RootCmd.AddCommand(query)
}
