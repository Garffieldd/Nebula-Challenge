package scripts

import (
	"math"
	"time"

	"github.com/tidwall/gjson"
)

/*
Struct created to hold the filtered TLS report information
*/
type FilteredTLSReport struct {
	Host         string              `json:"host"`
	OverallGrade string              `json:"overallGrade"`
	Endpoints    []FilteredEndpoints `json:"endpoints"` /// List of filtered endpoints
	Summary      string              `json:"summary"`
	Timestamp    time.Time           `json:"timestamp"`
}

/*
Struct created to hold the filtered endpoint information (that is in the FilteredTLSReport struct)
*/
type FilteredEndpoints struct {
	IPAddress                string   `json:"ipAddress"`
	Grade                    string   `json:"grade"`
	Certificate              []FilteredCertificate   `json:"certificate"`
	Protocols                []string `json:"protocols"`
	NegotiatedChiperStrength int      `json:"negotiatedCipherStrength"`
	MaxCipherStrength        int      `json:"maxCipherStrength"`
	HasWeakCiphers           bool     `json:"hasWeakCiphers"`
	HSTS                     string   `json:"hsts"`
	Server                   string   `json:"server"`
	Issues                   []string `json:"issues"`
}


/*
Struct created to hold the filtered certificate information (that is in the FilteredTLSReport->Endpoint struct)
*/
type FilteredCertificate struct {
	Subject       string  `json:subject`
	Issuer        string  `json:issuer`
	ValidityYears float64 `json:validityYears`
	ExpiresInDays float64 `json:expiresInDays`
}

func FilteredSSLReport(rawReport []byte) (*FilteredTLSReport, error) {
	report := &FilteredTLSReport{
		Host:      gjson.GetBytes(rawReport, "host").String(),
		Timestamp: time.Now(),
	}

	// endpoints := gjson.GetBytes(rawReport, "endpoints").Array()
	// for _, endpoint := range endpoints {
	// 	ipAdress := endpoint.Get("ipAdress").String()
	// 	grade := endpoint.Get("grade").String()
	// 	certificate := extractCertificateData(endpoint)
	// 	protocols := extractProtocols(endpoint)

	// }

	return report, nil
}

func extractCertificateData(endpoint gjson.Result) *FilteredCertificate {
	validityYears, expireInDays := calculateCertValidity(endpoint)
	certificate := &FilteredCertificate{
		Subject:       endpoint.Get("details.cert.subject").String(),
		Issuer:        endpoint.Get("details.cert.issuerSubject").String(),
		ValidityYears: validityYears,
		ExpiresInDays: expireInDays,
	}

	return certificate

}

func calculateCertValidity(endpoint gjson.Result) (validityYears float64, expiresInDays float64) {
	notBeforeMs := endpoint.Get("details.cert.notBefore").Float()
	notAfterMs := endpoint.Get("details.cert.notAfter").Float()

	if notBeforeMs == 0 || notAfterMs == 0 {
		return 0, 0
	}

	notBefore := time.UnixMilli(int64(notBeforeMs))
	notAfter := time.UnixMilli(int64(notAfterMs))

	validityDuration := notAfter.Sub(notBefore)
	validityDays := validityDuration.Hours() / 24
	validityYears = validityDays / 365.25

	currentTime := time.Now().UTC()
	expiresDuration := notAfter.Sub(currentTime)
	expiresInDays = expiresDuration.Hours() / 24

	validityYears = math.Round(validityYears*10) / 10
	expiresInDays = math.Round(expiresInDays*10) / 10

	return validityYears, expiresInDays
}

func extractProtocols(endpoint gjson.Result) (protocolsArray []string) {
	protocolsResult := []string{}
	protocols := endpoint.Get("details.protocols").Array()
	for _, protocol := range protocols {
		protocolsResult = append(protocolsResult, protocol.Get("name").String())
	}

	return protocolsResult

}
