// Package dashfuncs holds functions for dashport
package dashfuncs

import (
	"bytes"
	"dashport/internal/jsonstructs"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func setPushHeaders(r *http.Request, apitoken string) {
	r.Header.Set("Authorization", apitoken)
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
}

func addPullHeaders(r *http.Request, apitoken string) {
	r.Header.Add("Authorization", apitoken)
	r.Header.Add("accept", "application/json; charset=utf-8")
}

// ReadConfig reads config from config path
func ReadConfig() (jsonstructs.Configuration, error) {
	file, _ := os.Open(os.Getenv("HOME") + "/.config/dashport/" + "dashportcfg.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := jsonstructs.Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		return configuration, err
	}
	return configuration, nil
}

// DashHandler handles all dashboard requests
func DashHandler(r jsonstructs.ReqParts, b io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(r.Method, r.URL, b)
	if err != nil {
		var resp = http.Response{}
		return &resp, err
	}

	switch r.Method {
	case http.MethodGet, http.MethodDelete:
		addPullHeaders(req, r.Token)
	case http.MethodPost, http.MethodPut:
		setPushHeaders(req, r.Token)
	default:
		return nil, fmt.Errorf("\033[31mcould not handle %s request for: %s\033[m", r.Method, r.Action)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, err
	}

	return resp, nil

}

// DashCopy gets dashboard from id and returns modified body to caller
func DashCopy(r *jsonstructs.ReqParts) (*bytes.Reader, error) {
	dashOut := jsonstructs.Dashboard{}
	var durl string

	r.Method = http.MethodGet
	r.URL += "/dashboards/" + r.ID

	dashResp, err := DashHandler(*r, nil)
	if err != nil {
		return nil, err
	}

	jsonErr := json.NewDecoder(dashResp.Body).Decode(&dashOut)
	if jsonErr != nil {
		return nil, fmt.Errorf("\033[31%s\033[m", jsonErr)
	}

	durl, r.Token, err = BuildURL(r.Denv, r.Conf)
        if err != nil {
                return nil, err
        }

	switch r.Action {
	case "clone":
		dashOut.ID = ""
		dashOut.DashboardMetadata.Owner = ""
		r.Method = http.MethodPost
		r.URL = durl + "/dashboards/"
	case "update":
		dashOut.ID = r.Did
		r.Method = http.MethodPut
		r.URL = durl + "/dashboards/" + r.Did
	}

	dashOutb, err := json.Marshal(&dashOut)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(dashOutb), nil

}

// BuildURL builds url from environment
func BuildURL(env string, configuration jsonstructs.Configuration) (url, apitoken string, err error) {

	for _, c := range configuration.APIConfig.Tenants {
		switch c.Env {
		case env:
			url := c.URL
			apitoken := "Api-Token " + c.Token
			return url, apitoken, nil
		}
	}
	return "", "", fmt.Errorf("\033[31munknown tenant env: %s\033[m", env)
}

// DashResp handles status and prints result from http.Response
func DashResp(resp *http.Response, doWhat string) error {

	if resp.StatusCode > 204 {
		var dashError jsonstructs.DashError
		jsonErr := json.NewDecoder(resp.Body).Decode(&dashError)
		if jsonErr != nil {
			return fmt.Errorf("\033[31%s\033[m", jsonErr)
		}

		jsonRes, err := json.MarshalIndent(&dashError, "", "\t")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Doom red
		fmt.Fprintf(os.Stdout, "\033[31m%s\033[m\n", "Fail")
		fmt.Fprintf(os.Stdout, "%s\n", jsonRes)
		return nil
	}

	switch doWhat {
	case "update":
		fmt.Fprintf(os.Stdout, "\033[32m%s\033[m\n", "Dashboard updated succesfully.")
		return nil
	case "delete":
		fmt.Fprintf(os.Stdout, "\033[32m%s\033[m\n", "Dashboard deleted succesfully.")
		return nil
	case "clone":
		var dashSuccess jsonstructs.DashSuccess
		jsonErr := json.NewDecoder(resp.Body).Decode(&dashSuccess)
		if jsonErr != nil {
			return fmt.Errorf("\033[31%s\033[m", jsonErr)
		}

		jsonRes, err := json.MarshalIndent(&dashSuccess, "", "\t")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		fmt.Fprintf(os.Stdout, "%s\n", jsonRes)
		return nil
	case "print":
		var dashboard jsonstructs.Dashboard
		jsonErr := json.NewDecoder(resp.Body).Decode(&dashboard)
		if jsonErr != nil {
			return fmt.Errorf("\033[31%s\033[m", jsonErr)
		}

		jsonRes, err := json.MarshalIndent(&dashboard, "", "\t")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		fmt.Fprintf(os.Stdout, "%s\n", jsonRes)
		return nil
	case "printall":
		var dashboards jsonstructs.DashboardsAll
		jsonErr := json.NewDecoder(resp.Body).Decode(&dashboards)
		if jsonErr != nil {
			return fmt.Errorf("\033[31%s\033[m", jsonErr)
		}

		jsonRes, err := json.MarshalIndent(&dashboards, "", "\t")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		fmt.Fprintf(os.Stdout, "%s\n", jsonRes)
		return nil
	default:
		return fmt.Errorf("unknown response handling error exiting")

	}

}
