package scripts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

/*
CheckTLS initiates a TLS assessment for the given domain using SSL Labs API
Args:
	domain string: The domain to assess
	startnew bool: Whether to start a new assessment or use cached results
Returns:
	map[string]interface{}: The assessment result
	error: Any error encountered during the process
*/
func CheckTLS(domain string, startnew bool) (map[string]interface{}, error) {
	// First SSL Labs API endpoints for TLS checking
	var SSL_Lab_Api_Entrypoint = fmt.Sprintf("https://api.ssllabs.com/api/v2/analyze?host=%s&publish=off&all=done&ignoreMismatch=on", domain)
	if startnew {
		SSL_Lab_Api_Entrypoint += "&startNew=on" // Indicates to start a new assessment
	} else {
		SSL_Lab_Api_Entrypoint += "&fromCache=on" // Indicates to use cached results if available
	}

	req, _ := http.NewRequest("GET", SSL_Lab_Api_Entrypoint, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	return result, nil

}

/*
pollUntilReady continuously polls the SSL Labs API until the TLS assessment for the domain is ready
Args:
	domain string: The domain to assess
Returns:
	map[string]interface{}: The final assessment result
	error: Any error encountered during the process
*/
func PollUntilReady(domain string) (map[string]interface{}, error) {
	for { // Bucle infinito hasta que llegue a un return o break
		result, err := CheckTLS(domain, false) // Llama a la funcion CheckTLS con startnew en false porque ya se inicio la evaluacion antes
		if err != nil {
			return nil, err
		}
		status := result["status"].(string)
		if status == "READY" {
			return result, nil
		} else if status == "ERROR" {
			return nil, fmt.Errorf("Error during TLS assessment for domain %s", domain)
		}
		time.Sleep(10 * time.Second) // Wait before polling again
	}
}

// func searchForDomainResults(domain string); error {
// 	_, err := CheckTLS(domain, true) // Start a new assessment
// 	if err != nil {
// 		return err
// 	}
// 	finalResult, err := pollUntilReady(domain) // Poll until the assessment is ready
// 	if err != nil {
// 		return err
// 	}

// 	return err

// }
