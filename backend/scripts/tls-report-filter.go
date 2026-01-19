package scripts

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

/*
Struct created to hold the filtered TLS report information
*/
type FilteredTLSReport struct {
	Host        string             `json:"host"`
	WebProtocol string             `json:"webProtocol"`
	Endpoints   []FilteredEndpoint `json:"endpoints"` // List of filtered endpoints
	Summary     string             `json:"summary"`
	Timestamp   time.Time          `json:"timestamp"`
}

/*
Struct created to hold the filtered endpoint information (that is in the FilteredTLSReport struct)
*/
type FilteredEndpoint struct {
	IPAddress                string               `json:"ipAddress"`
	Grade                    string               `json:"grade"`
	HasWarnings              bool                 `json:"hasWarnings"`
	IsExceptional            bool                 `json:"isExceptional"`
	Certificate              *FilteredCertificate `json:"certificate"`
	Protocols                []string             `json:"protocols"`
	NegotiatedCipherStrength float64              `json:"negotiatedCipherStrength"`
	MaxCipherStrength        float64              `json:"maxCipherStrength"`
	HasWeakCiphers           bool                 `json:"hasWeakCiphers"`
	HSTS                     string               `json:"hsts"`
	Server                   string               `json:"server"`
	ChainIssues              int64                `json:"issues"`
}

/*
Struct created to hold the filtered certificate information (that is in the FilteredTLSReport->Endpoint struct)
*/
type FilteredCertificate struct {
	Subject       string  `json:"subject"`
	Issuer        string  `json:"issuer"`
	ValidityYears float64 `json:"validityYears"`
	ExpiresInDays float64 `json:"expiresInDays"`
}

/*
FilterSSLReport is a function that assembles the new parsed report object and returns it
Args:

	rawReport []byte: the report info in byte format,so the gjson library can use it

Returns:

	*FilteredTLSReport: Pointer of the FilterTLSReport Struct
	error: Any error encountered during the process
*/
func FilterSSLReport(rawReport []byte) (*FilteredTLSReport, error) {

	if len(rawReport) == 0 {
		return nil, fmt.Errorf("rawReport is empty (no data received from SSL Labs)")
	}

	if !gjson.ValidBytes(rawReport) {
		return nil, fmt.Errorf("rawReport is not valid JSON")
	}

	report := &FilteredTLSReport{
		Host:        gjson.GetBytes(rawReport, "host").String(),
		WebProtocol: gjson.GetBytes(rawReport, "protocol").String(),
		Timestamp:   time.Now(),
	}

	endpointsData := gjson.GetBytes(rawReport, "endpoints")
	if !endpointsData.Exists() || !endpointsData.IsArray() {
		return report, nil
	}

	var filteredEndpoints []FilteredEndpoint

	for _, endpoint := range endpointsData.Array() {
		fe := FilteredEndpoint{
			IPAddress:                endpoint.Get("ipAddress").String(),
			Grade:                    endpoint.Get("grade").String(),
			HasWarnings:              endpoint.Get("hasWarnings").Bool(),
			IsExceptional:            endpoint.Get("isExceptional").Bool(),
			Protocols:                extractProtocols(endpoint),
			NegotiatedCipherStrength: endpoint.Get("details.suites.list[0].cipherStrength").Float(),
			MaxCipherStrength:        captureMaxCipherStrength(endpoint),
			HasWeakCiphers:           existsWeakCipher(endpoint),
			HSTS:                     endpoint.Get("details.hstsPolicy.status").String(),
			Server:                   endpoint.Get("details.serverSignature").String(),
			ChainIssues:              endpoint.Get("details.chain.issues").Int(),
			Certificate:              extractCertificateData(endpoint),
		}

		filteredEndpoints = append(filteredEndpoints, fe)
	}

	report.Endpoints = filteredEndpoints
	report.Summary = generateSummary(report)

	return report, nil
}

/*
extractCertificateData assembles the certificate object.
Args:

	endpoint gjson.Result: The endpoint result given by the library gjson

Returns:

	*FilteredCertificate: Pointer of the FilteredCertificate Struct
	error: Any error encountered during the process
*/
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

/*
calculateCertValidity grab the notBefore and notAfter info from the report and transforms them into validityYears and expiresInDays.
Args:

	endpoint gjson.Result: The endpoint result given by the library gjson

Returns:

	validityYears float64: A float number that indicates in how many years is the certificate valid for
	expiresInDays float64: A float number that indicates in how many days is expiring the certificate
*/
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

/*
extractProtocols assembles the slice of the protocols names.
Args:

	endpoint gjson.Result: The endpoint result given by the library gjson

Returns:

	protocolsArray []string: Slice with the protocols names
*/
func extractProtocols(endpoint gjson.Result) (protocolsArray []string) {

	endpoint.Get("details.protocols").ForEach(func(_, p gjson.Result) bool {
		protocolsArray = append(protocolsArray, p.Get("name").String()+" "+p.Get("version").String())
		return true
	})

	return protocolsArray

}

/*

captureMaxCipherStrength search for the maximum value of CipherStrength in the Report
Args:

			endpoint gjson.Result: The endpoint result given by the library gjson

	Returns:
			maximumMaxCipherStrength float64: the maximum value of CipherStrenghtr founded in the report

*/

func captureMaxCipherStrength(endpoint gjson.Result) (maximumMaxCipherStrength float64) {
	suiteList := endpoint.Get("details.suites.list").Array()
	maximumMaxCipherStrength = 0.0
	for _, value := range suiteList {
		if value.Get("cipherStrength").Float() > maximumMaxCipherStrength {
			maximumMaxCipherStrength = value.Get("cipherStrength").Float()
		}
	}

	return maximumMaxCipherStrength
}

/*existsWeakCipher search for a value of ciphererStrenght <= 112
Args:

			endpoint gjson.Result: The endpoint result given by the library gjson

	Returns:
			existWeakCipher bool: returns true if exist, false otherwise

*/

func existsWeakCipher(endpoint gjson.Result) (existWeakCipher bool) {
	suiteList := endpoint.Get("details.suites.list").Array()
	existWeakCipher = false
	weakThreshold := 112.0
	for _, value := range suiteList {
		if value.Get("cipherStrength").Float() < weakThreshold {
			existWeakCipher = true
		}
	}

	return existWeakCipher
}

/* generateSummary builds a string based on the domain information and additionally builds a verdict.
Args:

		reportInfo *FilteredTLSReport: is the report info so far and will be used to put together the summary

Returns:
		string: the complete summary + verdict

*/

func generateSummary(reportInfo *FilteredTLSReport) string {
	if reportInfo == nil || len(reportInfo.Endpoints) == 0 {
		return "No valid information could be obtained from the TLS analysis"
	}

	// initializing the variables
	bestPriority := getGradePriority("F")
	bestGrade := "F"
	hasWarningsAny := false
	isExceptionalAny := false
	hasTLS13 := false
	hasHSTS := false
	hasWeakCiphersAny := false
	minExpiresDays := 9999.0
	chainIssuesAny := int64(0)

	for _, endpoint := range reportInfo.Endpoints {
		currentPriority := getGradePriority(endpoint.Grade)
		if currentPriority > bestPriority {
			bestPriority = currentPriority
			bestGrade = endpoint.Grade
		}
		if endpoint.HasWarnings {
			hasWarningsAny = true
		}
		if endpoint.IsExceptional {
			isExceptionalAny = true
		}
		if contains(endpoint.Protocols, "TLS 1.3") {
			hasTLS13 = true
		}
		if strings.Contains(strings.ToLower(endpoint.HSTS), "present") {
			hasHSTS = true
		}
		if endpoint.HasWeakCiphers {
			hasWeakCiphersAny = true
		}
		if endpoint.ChainIssues > chainIssuesAny {
			chainIssuesAny = endpoint.ChainIssues
		}
		if endpoint.Certificate != nil && endpoint.Certificate.ExpiresInDays < minExpiresDays {
			minExpiresDays = endpoint.Certificate.ExpiresInDays
		}
	}

	var sb strings.Builder

	// Introduction
	sb.WriteString(fmt.Sprintf(" Análisis TLS para %s - Calificación general: %s", reportInfo.Host, bestGrade))

	// Protocols
	if hasTLS13 {
		sb.WriteString(" - Soporta TLS 1.3 (excelente nivel de seguridad actual).")
	} else if containsAny(reportInfo.Endpoints[0].Protocols, "TLS 1.2") {
		sb.WriteString(" - Soporta TLS 1.2, pero sin TLS 1.3 (aceptable, pero no óptimo en 2026).")
	} else {
		sb.WriteString(" - Protocolos obsoletos o inseguros detectados.")
	}

	// Cipher strength
	if len(reportInfo.Endpoints) > 0 {
		ep := reportInfo.Endpoints[0] //first as reference
		if ep.MaxCipherStrength >= 256 {
			sb.WriteString(" - Cifrado fuerte (hasta 256 bits).")
		} else if ep.MaxCipherStrength >= 128 {
			sb.WriteString(" - Cifrado aceptable (128 bits).")
		} else {
			sb.WriteString(" - Cifrado débil detectado.")
		}

		if hasWeakCiphersAny {
			sb.WriteString(" - Atención: hay suites cifradas débiles habilitadas.")
		}
	}

	// HSTS
	if hasHSTS {
		sb.WriteString(" - HSTS está activo (buena protección contra downgrade).")
	} else {
		sb.WriteString(" - Sin HSTS → vulnerable a ataques de downgrade (HTTP plano posible).")
	}

	// Certificate
	if minExpiresDays > 30 {
		sb.WriteString(" - Certificado válido por más de 30 días.")
	} else if minExpiresDays > 0 {
		sb.WriteString(fmt.Sprintf(" - Certificado expira en %.1f días → renovar pronto.", minExpiresDays))
	} else {
		sb.WriteString(" - Certificado expirado o inválido → sitio inseguro.")
	}

	// Warnings and exceptions
	if hasWarningsAny {
		sb.WriteString(" - Existen advertencias menores en la configuración.")
	}
	if isExceptionalAny {
		sb.WriteString(" - Al menos un endpoint tiene configuración excepcional.")
	}
	if chainIssuesAny > 0 {
		sb.WriteString(" - Problemas detectados en la cadena de certificados.\n")
	}

	// Final verdict
	var verdict string
	switch {
	case bestGrade == "A+" || (bestGrade == "A" && isExceptionalAny && !hasWarningsAny && hasTLS13 && hasHSTS):
		verdict = " Excelente"
	case bestGrade == "A" || (bestGrade == "A-" && hasTLS13 && hasHSTS):
		verdict = " Buena"
	case bestGrade == "B" || bestGrade == "C":
		verdict = " Aceptable (se recomienda mejorar)"
	case bestGrade == "D" || bestGrade == "E":
		verdict = " Deficiente (riesgo alto)"
	default:
		verdict = " Muy mala (sitio inseguro)"
	}

	sb.WriteString(" VEREDICTO FINAL: " + verdict)

	return sb.String()
}

/*
getGradePriority is used to manage a numerical priority scheme for letter grades.

Args:

	grade string: is the letter that represents a grade

Return:

	int: Returns the numeric priority (the higher the better)
*/
func getGradePriority(grade string) int {
	var gradePriority = map[string]int{
		"A+": 11,
		"A":  10,
		"A-": 9,
		"B+": 8,
		"B":  7,
		"B-": 6,
		"C+": 5,
		"C":  4,
		"C-": 3,
		"D+": 2,
		"D":  1,
		"D-": 0,
		"E+": -1,
		"E":  -2,
		"E-": -3,
		"F":  -10,
	}

	if p, ok := gradePriority[grade]; ok {
		return p
	}
	return -100 // unknown, priority low
}

/*
contains checks if a given item exists in a slice of strings.

Args:

	slice []string: The slice of strings to search in.
	item string: The string value to look for.

Returns:

	bool: true if the item is found in the slice, false otherwise.
*/
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

/*
containsAny checks if any of the provided items exist in a slice of strings.

Args:

	slice []string: The slice of strings to search in.
	items ...string: One or more strings to look for.

Returns:

	bool: true if at least one of the items is found in the slice,
	      false if none are found.
*/
func containsAny(slice []string, items ...string) bool {
	for _, item := range items {
		if contains(slice, item) {
			return true
		}
	}
	return false
}
