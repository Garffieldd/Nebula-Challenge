package scripts

import (
	"fmt"
	"github.com/tidwall/gjson"
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

	[]byte: The assessment result in byte format
	error: Any error encountered during the process
*/
func CheckTLS(domain string, startnew bool) ([]byte, error) {
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
	//var result []byte
	//json.Unmarshal(body, &result)
	return body, nil

}

/*
pollUntilReady continuously polls the SSL Labs API until the TLS assessment for the domain is ready
Args:

	domain string: The domain to assess

Returns:

	[]byte: The final assessment result in byte format
	error: Any error encountered during the process
*/
func PollUntilReady(domain string) ([]byte, error) {

	for { // Bucle infinito hasta que llegue a un return o break
		result, err := CheckTLS(domain, false) // Llama a la funcion CheckTLS con startnew en false porque ya se inicio la evaluacion antes
		if err != nil {
			return nil, err
		}
		status := gjson.GetBytes(result, "status").String()
		switch status {
		case "READY":
			return result, nil
		case "ERROR":
			return nil, fmt.Errorf("Error during TLS assessment for domain %s", domain)
		}
		time.Sleep(10 * time.Second) // Wait before polling again
	}
}
