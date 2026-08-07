package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	baseURL := "http://localhost:8080/api/v1"

	endpoints := []struct {
		Method string
		URL    string
		Body   string
		Desc   string
	}{
		// Auth Endpoints
		{"POST", "/auth/register", `{"name":"Test User","email":"test12345@test.com","password":"password123","role":"patient"}`, "Auth: Register"},
		{"POST", "/auth/login", `{"email":"test12345@test.com","password":"password123"}`, "Auth: Login"},
		{"GET", "/auth/profile", "", "Auth: Profile (Without Token)"}, // Should fail but structurally return success: false
		{"POST", "/auth/logout", "", "Auth: Logout"},
		{"POST", "/auth/forgot-password", `{"email":"test12345@test.com"}`, "Auth: Forgot Password"},

		// Patient Endpoints
		{"GET", "/patients", "", "Patient: Index"},
		{"GET", "/patients/9999", "", "Patient: Show"},
		{"POST", "/patients", `{}`, "Patient: Store (Validation Error expected)"},
		{"PUT", "/patients/9999", `{}`, "Patient: Update"},
		{"DELETE", "/patients/9999", "", "Patient: Destroy"},

		// Medical Record Endpoints
		{"GET", "/medical-records", "", "MedicalRecord: Index"},
		{"GET", "/patients/9999/medical-records", "", "MedicalRecord: History"},
		{"POST", "/medical-records", `{}`, "MedicalRecord: Store"},
		{"PUT", "/medical-records/9999", `{}`, "MedicalRecord: Update"},
		{"DELETE", "/medical-records/9999", "", "MedicalRecord: Destroy"},

		// Appointment Endpoints
		{"GET", "/appointments", "", "Appointment: Index"},
		{"POST", "/appointments", `{}`, "Appointment: Store"},
		{"PUT", "/appointments/9999", `{}`, "Appointment: Update"},
		{"POST", "/appointments/9999/cancel", "", "Appointment: Cancel"},
		{"POST", "/appointments/9999/reschedule", `{}`, "Appointment: Reschedule"},
	}

	fmt.Println("=== API Migration Structural Test ===")
	fmt.Printf("Targeting Base URL: %s\n\n", baseURL)

	passed := 0
	failed := 0

	for _, ep := range endpoints {
		fmt.Printf("Testing: %s [%s %s]\n", ep.Desc, ep.Method, ep.URL)
		
		var req *http.Request
		var err error
		if ep.Body != "" {
			req, err = http.NewRequest(ep.Method, baseURL+ep.URL, bytes.NewBuffer([]byte(ep.Body)))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, err = http.NewRequest(ep.Method, baseURL+ep.URL, nil)
		}

		if err != nil {
			fmt.Printf("  [FAIL] Failed to create request: %v\n", err)
			failed++
			continue
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("  [FAIL] Request failed. Is the server running? %v\n", err)
			failed++
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Printf("  [FAIL] Response is not valid JSON: %v\n", err)
			failed++
			continue
		}

		// Check for Laravel ApiResponse Trait parity
		_, hasSuccess := result["success"]
		_, hasMessage := result["message"]
		
		if !hasSuccess || !hasMessage {
			fmt.Printf("  [FAIL] Missing 'success' or 'message' wrapper. Got: %s\n", string(body))
			failed++
			continue
		}

		// Check inner structures
		isSuccess := result["success"] == true
		if isSuccess {
			if _, hasData := result["data"]; !hasData {
				fmt.Printf("  [FAIL] Successful response missing 'data' field. Got: %s\n", string(body))
				failed++
				continue
			}
		} else {
			// Error response might have 'errors' if it's 422
			if resp.StatusCode == 422 {
				if _, hasErrors := result["errors"]; !hasErrors {
					fmt.Printf("  [FAIL] 422 Error response missing 'errors' object. Got: %s\n", string(body))
					failed++
					continue
				}
			}
		}

		fmt.Printf("  [PASS] Status: %d, Response structurally matches Laravel.\n", resp.StatusCode)
		passed++
	}

	fmt.Printf("\n--- Test Summary ---\n")
	fmt.Printf("Total: %d | Passed: %d | Failed: %d\n", len(endpoints), passed, failed)
	
	if failed == 0 {
		fmt.Println("🎉 ALL ENDPOINTS MATCH LARAVEL STRUCTURAL PARITY!")
	} else {
		fmt.Println("❌ SOME ENDPOINTS FAILED.")
	}
}
