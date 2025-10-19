package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const baseURL = "http://localhost:9091"

type Line struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Station struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	LineID string `json:"line_id"`
	Status string `json:"status"`
}

type Incident struct {
	Line            string `json:"line"`
	Station         string `json:"station"`
	Timestamp       string `json:"timestamp"`
	DurationMinutes int    `json:"duration_minutes"`
	IncidentType    string `json:"incident_type"`
}

func main() {
	fmt.Println("🚀 Starting to populate Transport Analytics database with sample data...")
	fmt.Println()

	// Wait for API to be ready
	if err := waitForAPI(); err != nil {
		fmt.Printf("❌ Error: API is not available: %v\n", err)
		fmt.Println("Please ensure the application is running with: docker-compose up or make run")
		os.Exit(1)
	}

	// Create Lines
	fmt.Println("📍 Creating lines...")
	lines := []string{
		"North-South Line",
		"East-West Line",
		"Circle Line",
		"Downtown Line",
		"Thomson-East Coast Line",
	}

	lineIDs := make(map[string]string)
	for _, lineName := range lines {
		line, err := createLine(lineName)
		if err != nil {
			fmt.Printf("❌ Failed to create line %s: %v\n", lineName, err)
			continue
		}
		lineIDs[lineName] = line.ID
		fmt.Printf("  ✓ Created line: %s (ID: %s)\n", lineName, line.ID)
	}
	fmt.Println()

	// Create Stations
	fmt.Println("🚉 Creating stations...")
	stations := map[string][]string{
		"North-South Line": {
			"Jurong East", "Bukit Batok", "Bukit Gombak", "Choa Chu Kang",
			"Yew Tee", "Kranji", "Marsiling", "Woodlands", "Admiralty",
			"Sembawang", "Canberra", "Yishun", "Khatib", "Yio Chu Kang",
			"Ang Mo Kio", "Bishan", "Braddell", "Toa Payoh", "Novena",
			"Newton", "Orchard", "Somerset", "Dhoby Ghaut", "City Hall",
			"Raffles Place", "Marina Bay",
		},
		"East-West Line": {
			"Pasir Ris", "Tampines", "Simei", "Tanah Merah", "Bedok",
			"Kembangan", "Eunos", "Paya Lebar", "Aljunied", "Kallang",
			"Lavender", "Bugis", "City Hall", "Raffles Place", "Tanjong Pagar",
			"Outram Park", "Tiong Bahru", "Redhill", "Queenstown",
		},
		"Circle Line": {
			"Dhoby Ghaut", "Bras Basah", "Esplanade", "Promenade",
			"Nicoll Highway", "Stadium", "Mountbatten", "Dakota",
			"Paya Lebar", "MacPherson", "Tai Seng", "Bartley",
			"Serangoon", "Lorong Chuan", "Bishan", "Marymount",
		},
		"Downtown Line": {
			"Bukit Panjang", "Cashew", "Hillview", "Beauty World",
			"King Albert Park", "Sixth Avenue", "Tan Kah Kee", "Botanic Gardens",
			"Stevens", "Newton", "Little India", "Rochor", "Bugis",
			"Promenade", "Bayfront", "Downtown", "Telok Ayer", "Chinatown",
		},
		"Thomson-East Coast Line": {
			"Woodlands North", "Woodlands", "Woodlands South", "Springleaf",
			"Lentor", "Mayflower", "Bright Hill", "Upper Thomson",
			"Caldecott", "Mount Pleasant", "Stevens", "Napier", "Orchard Boulevard",
		},
	}

	stationIDs := make(map[string]string)
	for lineName, stationNames := range stations {
		lineID, ok := lineIDs[lineName]
		if !ok {
			continue
		}
		for _, stationName := range stationNames {
			station, err := createStation(stationName, lineID, "active")
			if err != nil {
				fmt.Printf("❌ Failed to create station %s: %v\n", stationName, err)
				continue
			}
			key := fmt.Sprintf("%s|%s", lineName, stationName)
			stationIDs[key] = station.ID
			fmt.Printf("  ✓ Created station: %s on %s (ID: %s)\n", stationName, lineName, station.ID)
		}
	}
	fmt.Println()

	// Create Incidents
	fmt.Println("⚠️  Creating sample incidents (400+ over last 90 days)...")

	// Create random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Incident types with their probabilities
	incidentTypes := []string{"mechanical", "signal", "power", "mechanical", "signal", "mechanical", "weather", "other"}

	// Generate 420 incidents spread over 90 days
	incidentCount := 0
	successCount := 0
	failCount := 0

	// Create a list of all line-station pairs
	lineStationPairs := []struct {
		Line    string
		Station string
	}{}

	for lineName, stationNames := range stations {
		for _, stationName := range stationNames {
			lineStationPairs = append(lineStationPairs, struct {
				Line    string
				Station string
			}{lineName, stationName})
		}
	}

	// Generate 420 incidents
	for i := 0; i < 420; i++ {
		// Random day in the last 90 days
		daysAgo := rng.Intn(90)

		// Random hour (weighted toward peak hours 7-9am and 5-8pm)
		hour := rng.Intn(24)
		if rng.Float32() < 0.4 { // 40% during peak hours
			if rng.Float32() < 0.5 {
				hour = 7 + rng.Intn(3) // 7-9am
			} else {
				hour = 17 + rng.Intn(4) // 5-8pm
			}
		}

		// Random minute and second
		minute := rng.Intn(60)
		second := rng.Intn(60)

		// Calculate timestamp
		timestamp := time.Now().AddDate(0, 0, -daysAgo).
			Truncate(24 * time.Hour).
			Add(time.Duration(hour) * time.Hour).
			Add(time.Duration(minute) * time.Minute).
			Add(time.Duration(second) * time.Second)

		// Random duration (5-120 minutes, weighted toward shorter durations)
		var duration int
		r := rng.Float32()
		if r < 0.5 {
			duration = 5 + rng.Intn(16) // 5-20 mins (50%)
		} else if r < 0.8 {
			duration = 20 + rng.Intn(31) // 20-50 mins (30%)
		} else {
			duration = 50 + rng.Intn(71) // 50-120 mins (20%)
		}

		// Random incident type
		incidentType := incidentTypes[rng.Intn(len(incidentTypes))]

		// Random line-station pair
		pair := lineStationPairs[rng.Intn(len(lineStationPairs))]

		incidentCount++
		err := createIncident(pair.Line, pair.Station, timestamp, duration, incidentType)
		if err != nil {
			failCount++
			// Only print first few errors to avoid spam
			if failCount <= 5 {
				fmt.Printf("  ⚠ Failed to create incident #%d: %v\n", incidentCount, err)
			}
			continue
		}
		successCount++

		// Print progress every 50 incidents
		if successCount%50 == 0 {
			fmt.Printf("  ✓ Created %d incidents...\n", successCount)
		}
	}

	fmt.Printf("  ✓ Successfully created %d incidents (%d failed)\n", successCount, failCount)
	fmt.Println()

	fmt.Println("✅ Database population complete!")
	fmt.Println()
	fmt.Println("📊 Summary:")
	fmt.Printf("  - Lines created: %d\n", len(lineIDs))
	fmt.Printf("  - Stations created: %d\n", len(stationIDs))
	fmt.Printf("  - Incidents created: %d\n", successCount)
	fmt.Println()
	fmt.Println("You can now view the data at:")
	fmt.Println("  - Swagger UI: http://localhost:9091/swagger/")
	fmt.Println("  - API: http://localhost:9091")
}

func waitForAPI() error {
	fmt.Println("⏳ Waiting for API to be ready...")
	for i := 0; i < 30; i++ {
		resp, err := http.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			fmt.Println("✓ API is ready!")
			fmt.Println()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
		fmt.Print(".")
	}
	return fmt.Errorf("timeout waiting for API")
}

func createLine(name string) (*Line, error) {
	payload := map[string]string{"name": name}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(baseURL+"/lines", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}

	var line Line
	if err := json.NewDecoder(resp.Body).Decode(&line); err != nil {
		return nil, err
	}
	return &line, nil
}

func createStation(name, lineID, status string) (*Station, error) {
	payload := map[string]string{
		"name":    name,
		"line_id": lineID,
		"status":  status,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(baseURL+"/stations", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}

	var station Station
	if err := json.NewDecoder(resp.Body).Decode(&station); err != nil {
		return nil, err
	}
	return &station, nil
}

func createIncident(line, station string, timestamp time.Time, duration int, incidentType string) error {
	payload := map[string]interface{}{
		"line":             line,
		"station":          station,
		"timestamp":        timestamp.Format(time.RFC3339),
		"duration_minutes": duration,
		"incident_type":    incidentType,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(baseURL+"/incidents", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
