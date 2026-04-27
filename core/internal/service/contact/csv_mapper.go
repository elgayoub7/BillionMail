package contact

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// headerAliases maps standard field names to common CSV header variations
var headerAliases = map[string][]string{
	"email":      {"email", "e-mail", "email address", "mail", "courriel"},
	"first_name": {"first_name", "firstname", "first name", "fname", "given name", "prenom", "prénom", "givenname"},
	"last_name":  {"last_name", "lastname", "last name", "lname", "surname", "family name", "nom", "nom de famille"},
	"company":    {"company", "organization", "organisation", "org", "company name", "societe", "société", "entreprise"},
	"phone":      {"phone", "telephone", "tel", "phone number", "mobile", "cell", "téléphone"},
	"title":      {"title", "job title", "position", "role", "poste", "fonction"},
	"website":    {"website", "url", "site", "web", "site web", "domain"},
	"city":       {"city", "ville", "town"},
	"country":    {"country", "pays", "nation"},
	"industry":   {"industry", "sector", "secteur", "industrie"},
	"linkedin":   {"linkedin", "linkedin url", "linkedin_url"},
}

// CSVParseResult holds the result of parsing a CSV file
type CSVParseResult struct {
	Headers     []string            `json:"headers"`
	Rows        []map[string]string `json:"rows"`
	TotalRows   int                 `json:"total_rows"`
	Mapping     map[string]string   `json:"mapping"`
	Unmapped    []string            `json:"unmapped"`
	InvalidRows int                 `json:"invalid_rows"`
	Duplicates  int                 `json:"duplicates"`
}

// MapHeaders auto-maps CSV headers to standard field names
func MapHeaders(headers []string) map[string]string {
	mapping := make(map[string]string)
	mapped := make(map[string]bool)

	for _, header := range headers {
		normalized := strings.TrimSpace(strings.ToLower(header))
		if normalized == "" {
			continue
		}

		for field, aliases := range headerAliases {
			if mapped[field] {
				continue
			}
			for _, alias := range aliases {
				if normalized == alias {
					mapping[header] = field
					mapped[field] = true
					break
				}
			}
			if mapped[field] {
				break
			}
		}

		if _, ok := mapping[header]; !ok {
			mapping[header] = "" // unmapped
		}
	}

	return mapping
}

// GetUnmappedHeaders returns headers that weren't auto-mapped
func GetUnmappedHeaders(mapping map[string]string) []string {
	var unmapped []string
	for header, field := range mapping {
		if field == "" {
			unmapped = append(unmapped, header)
		}
	}
	return unmapped
}

// ParseCSV parses CSV content and returns structured result
func ParseCSV(content string, maxPreview int) *CSVParseResult {
	result := &CSVParseResult{}
	reader := csv.NewReader(strings.NewReader(content))

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return result
	}
	result.Headers = headers
	result.Mapping = MapHeaders(headers)
	result.Unmapped = GetUnmappedHeaders(result.Mapping)

	// Read rows
	seen := make(map[string]bool)
	rowNum := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.InvalidRows++
			continue
		}
		rowNum++

		row := make(map[string]string)
		for i, val := range record {
			if i < len(headers) {
				row[headers[i]] = strings.TrimSpace(val)
			}
		}

		// Check email validity for duplicate detection
		email := ""
		for h, field := range result.Mapping {
			if field == "email" {
				email = strings.ToLower(row[h])
				break
			}
		}

		if email != "" && strings.Contains(email, "@") {
			if seen[email] {
				result.Duplicates++
				continue
			}
			seen[email] = true
		} else {
			result.InvalidRows++
			continue
		}

		result.TotalRows++
		if maxPreview <= 0 || len(result.Rows) < maxPreview {
			result.Rows = append(result.Rows, row)
		}
	}

	return result
}

// MapCSVToContacts converts parsed CSV rows to contact format using column mapping
func MapCSVToContacts(rows []map[string]string, columnMapping map[string]string) ([]g.Map, error) {
	var contacts []g.Map

	for _, row := range rows {
		email := ""
		attribs := make(map[string]string)

		for csvCol, targetField := range columnMapping {
			val, ok := row[csvCol]
			if !ok || val == "" {
				continue
			}

			switch targetField {
			case "email":
				email = val
			case "skip":
				// Skip this column
			default:
				// All other fields become attributes
				attribs[targetField] = val
			}
		}

		if email == "" {
			continue
		}

		contact := g.Map{
			"email":   email,
			"attribs": attribs,
		}
		contacts = append(contacts, contact)
	}

	if len(contacts) == 0 {
		return nil, fmt.Errorf("no valid contacts found")
	}

	return contacts, nil
}
