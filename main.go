package main

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"resty.dev/v3"
)

// CEPResponse represents the response from the Brasil API
type CEPResponse struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {
	s := server.NewMCPServer(
		"Brazilian ZIP code (CEP) Lookup 🇧🇷",
		"1.0.0",
	)

	cepTool := mcp.NewTool("lookup_address",
		mcp.WithDescription("Get address information from a Brazilian ZIP code (CEP). Accepts formats like 01310-100 or 01310100."),
		mcp.WithString("cep",
			mcp.Required(),
			mcp.Description("Brazilian ZIP code (CEP) with or without hyphen (e.g., 01310-100 or 01310100)"),
		),
	)

	s.AddTool(cepTool, handleCEPLookup)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func handleCEPLookup(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	cepValue, ok := args["cep"]
	if !ok {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "CEP parameter is required",
				},
			},
			IsError: true,
		}, nil
	}

	cep, ok := cepValue.(string)
	if !ok {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "CEP must be a string",
				},
			},
			IsError: true,
		}, nil
	}

	cleanCEP := cleanCEP(cep)
	if !isValidCEP(cleanCEP) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Invalid CEP format. Expected 8 digits (e.g., 01310100 or 01310-100)",
				},
			},
			IsError: true,
		}, nil
	}

	address, err := fetchAddressFromAPI(ctx, cleanCEP)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Failed to fetch address: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	response := formatAddress(address)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: response,
			},
		},
		IsError: false,
	}, nil
}

// cleanCEP removes non-numeric characters from CEP
func cleanCEP(cep string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(cep, "")
}

// isValidCEP checks if the CEP has exactly 8 digits
func isValidCEP(cep string) bool {
	return len(cep) == 8 && regexp.MustCompile(`^\d{8}$`).MatchString(cep)
}

// fetchAddressFromAPI calls the Brasil API and returns the address information
func fetchAddressFromAPI(ctx context.Context, cep string) (*CEPResponse, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	client := resty.New().
		SetTimeout(5 * time.Second)
	defer func() {
		if err := client.Close(); err != nil {
			fmt.Printf("failed to close client: %v\n", err)
		}
	}()

	var cepResponse CEPResponse
	var errorResponse map[string]any

	resp, err := client.R().
		WithContext(ctx).
		SetResult(&cepResponse).
		SetError(&errorResponse).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.IsError() {
		switch resp.StatusCode() {
		case http.StatusNotFound:
			return nil, fmt.Errorf("CEP not found")
		default:
			return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode(), resp.String())
		}
	}

	return &cepResponse, nil
}

// formatAddress creates a human-readable string from the CEP response
func formatAddress(address *CEPResponse) string {
	var parts []string

	if address.Street != "" {
		parts = append(parts, fmt.Sprintf("Street: %s", address.Street))
	}
	if address.Neighborhood != "" {
		parts = append(parts, fmt.Sprintf("Neighborhood: %s", address.Neighborhood))
	}
	if address.City != "" {
		parts = append(parts, fmt.Sprintf("City: %s", address.City))
	}
	if address.State != "" {
		parts = append(parts, fmt.Sprintf("State: %s", address.State))
	}
	if address.CEP != "" {
		parts = append(parts, fmt.Sprintf("CEP: %s", address.CEP))
	}

	return strings.Join(parts, "\n")
}
