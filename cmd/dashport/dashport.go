// Package main
package main

import (
	"dashport/internal/dashfuncs"
	"dashport/internal/jsonstructs"
	"net/http"
	"flag"
	"fmt"
	"os"
)

func errHandle(e error) int {
	fmt.Fprintf(os.Stdout, "%s\n", e)
	return 1
}

func main() {

	configuration, err := dashfuncs.ReadConfig()
	if err != nil {
		os.Exit(errHandle(err))
	}

	oenv := flag.String("oenv", "", "Set origin tenant environment to copy or print from")
	denv := flag.String("denv", "", "Set destination tenant environment to copy to")
	id := flag.String("id", "", "Set the id of the dashboard you want to see/copy.")
	did := flag.String("did", "", "Set the id of the destination dashboard you want to update.")
	act := flag.String("act", "", "print, printall, clone, delete and update")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *oenv == "" {
		fmt.Fprintf(os.Stdout, "%s\n", "Oenv needs to be set")
		flag.PrintDefaults()
		os.Exit(1)
	}

	url, apitoken, err := dashfuncs.BuildURL(*oenv, configuration)
	if err != nil {
		os.Exit(errHandle(err))
	}

	reqparts := jsonstructs.ReqParts{
		Action: *act,
		Env:    *oenv,
		Denv:   *denv,
		ID:     *id,
		Did:    *did,
		URL:    url,
		Token:  apitoken,
		Conf:   configuration,
	}

	switch *act {
	case "print", "printall":
		if *id == "" && *act == "print" {
			fmt.Printf("%s.\n", "print action needs -id and -oenv option")
			flag.Usage()
			os.Exit(1)
		}

		reqparts.Method = http.MethodGet
		reqparts.URL += "/dashboards/" + *id
		resp, err := dashfuncs.DashHandler(reqparts, nil)
		if err != nil {
			os.Exit(errHandle(err))
		}

		err = dashfuncs.DashResp(resp, *act)
		if err != nil {
			os.Exit(errHandle(err))
		}

	case "clone", "update":
		if *act == "clone" && *id == "" || *denv == "" {
			fmt.Printf("%s.\n", "clone action needs -oenv, -denv and -id options.")
			flag.Usage()
			os.Exit(1)
		}
		if *act == "update" && (*id == "" || *did == "" || *denv == "") {
			fmt.Printf("%s.\n", "update action needs -oenv, -denv, -id and -did options.")
			flag.Usage()
			os.Exit(1)
		}
		dbody, err := dashfuncs.DashCopy(&reqparts)
		if err != nil {
			os.Exit(errHandle(err))
		}

		resp, err := dashfuncs.DashHandler(reqparts, dbody)
		if err != nil {
			os.Exit(errHandle(err))
		}

		err = dashfuncs.DashResp(resp, *act)
		if err != nil {
			os.Exit(errHandle(err))
		}

	case "delete":
		if *id == "" || *oenv == "" {
			fmt.Printf("%s.\n", "delete action needs -id and -oenv option")
			flag.Usage()
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Deleting dashboard with id %s\n", *id)

		reqparts.Method = http.MethodDelete
		reqparts.URL += "/dashboards/" + *id

		resp, err := dashfuncs.DashHandler(reqparts, nil)
		if err != nil {
			os.Exit(errHandle(err))
		}

		err = dashfuncs.DashResp(resp, *act)
		if err != nil {
			os.Exit(errHandle(err))
		}

	default:
		fmt.Printf("%s.\n", "Unknown action")
	}

}
