package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// Input struct for JSON decoding
type Input1 struct {
	Values []string `json:"values"`
}

type Input struct {
	Values []float64 `json:"values"`
}

func readEpsFromUser(input string) []string {
	indexToElement := map[int]string{
		1:  "ПЛ-110 кВ",
		2:  "ПЛ-35 кВ",
		3:  "ПЛ-10 кВ",
		4:  "КЛ-10 кВ (траншея)",
		5:  "КЛ-10 кВ (кабельний канал)",
		6:  "Т-110 кВ",
		7:  "Т-35 кВ",
		8:  "Т-10 кВ (кабельна мережа 10 кВ)",
		9:  "Т-10 кВ (повітряна мережа 10 кВ)",
		10: "В-110 кВ (елегазовий)",
		11: "В-10 кВ (малооливний)",
		12: "В-10 кВ (вакуумний)",
		13: "АВ-0.38 кВ",
		14: "ЕД 6, 10 кВ",
		15: "ЕД 0,38 кВ",
	}

	var userSelectedKeys []string
	selectedOptions := strings.Fields(input)

	for _, option := range selectedOptions {
		index, err := strconv.Atoi(option)
		if err == nil {
			if key, exists := indexToElement[index]; exists {
				userSelectedKeys = append(userSelectedKeys, key)
			} else {
				fmt.Println("Invalid input detected. Try again.")
				return nil
			}
		} else {
			fmt.Println("Invalid input detected. Try again.")
			return nil
		}
	}

	return userSelectedKeys
}

func calculateTask1(input string, nString string) string {
	epsElements := map[string]map[string]float64{
		"ПЛ-110 кВ":          {"omega": 0.07, "tv": 10.0, "mu": 0.167, "tp": 35.0},
		"ПЛ-35 кВ":           {"omega": 0.02, "tv": 8.0, "mu": 0.167, "tp": 35.0},
		"ПЛ-10 кВ":           {"omega": 0.02, "tv": 10.0, "mu": 0.167, "tp": 35.0},
		"КЛ-10 кВ (траншея)": {"omega": 0.03, "tv": 44.0, "mu": 1.0, "tp": 9.0},
		"КЛ-10 кВ (кабельний канал)": {"omega": 0.005, "tv": 17.5, "mu": 1.0, "tp": 9.0},
		"Т-110 кВ": {"omega": 0.015, "tv": 100.0, "mu": 1.0, "tp": 43.0},
		"Т-35 кВ":  {"omega": 0.02, "tv": 80.0, "mu": 1.0, "tp": 28.0},
		"Т-10 кВ (кабельна мережа 10 кВ)":  {"omega": 0.005, "tv": 60.0, "mu": 0.5, "tp": 10.0},
		"Т-10 кВ (повітряна мережа 10 кВ)": {"omega": 0.05, "tv": 60.0, "mu": 0.5, "tp": 10.0},
		"В-110 кВ (елегазовий)":            {"omega": 0.01, "tv": 30.0, "mu": 0.1, "tp": 30.0},
		"В-10 кВ (малооливний)":            {"omega": 0.02, "tv": 15.0, "mu": 0.33, "tp": 15.0},
		"В-10 кВ (вакуумний)":              {"omega": 0.01, "tv": 15.0, "mu": 0.33, "tp": 15.0},
		"АВ-0.38 кВ":                       {"omega": 0.05, "tv": 4.0, "mu": 0.33, "tp": 10.0},
		"ЕД 6, 10 кВ":                      {"omega": 0.1, "tv": 160.0, "mu": 0.5, "tp": 0.0},
		"ЕД 0,38 кВ":                       {"omega": 0.1, "tv": 50.0, "mu": 0.5, "tp": 0.0},
	}

	userKeys := readEpsFromUser(input)

	if len(userKeys) == 0 {
		fmt.Println("Введені некоректні дані. Спробуйте ще раз.")
	}

	var omegaSum, tRecovery, maxTp float64

	for _, key := range userKeys {
		if element, exists := epsElements[key]; exists {
			omegaSum += element["omega"]
			tRecovery += element["omega"] * element["tv"]
			if element["tp"] > maxTp {
				maxTp = element["tp"]
			}
		}
	}

	n, err := strconv.ParseFloat(nString, 64)

	if err != nil {
		fmt.Println("Error:", err)
	}

	omegaSum += 0.03 * n
	tRecovery += 0.06 * n
	tRecovery /= omegaSum

	kAP := omegaSum * tRecovery / 8760
	kPP := 1.2 * maxTp / 8760
	omegaDK := 2 * 0.295 * (kAP + kPP)
	omegaDKS := omegaDK + 0.02

	output := fmt.Sprintf("Частота відмов одноколової системи: %.2f рік^-1\n", omegaSum) +
		fmt.Sprintf("Середня тривалість відновлення: %.2f год\n", tRecovery) +
		fmt.Sprintf("Коефіцієнт аварійного простою: %.2f\n", kAP) +
		fmt.Sprintf("Коефіцієнт планового простою: %.2f\n", kPP) +
		fmt.Sprintf("Частота відмов одночасно двох кіл двоколової системи: %.2f рік^-1\n", omegaDK) +
		fmt.Sprintf("Частота відмов двоколової системи з урахуванням секційного вимикача: %.2f рік^-1\n", omegaDKS)

	return output
}

func calculateTask2(zPerA, zPerP, omega, t, kp, Pm, Tm float64) string {
	MWA := omega * t * Pm * Tm
	MWP := kp * Pm * Tm
	M := zPerA*MWA + zPerP*MWP

	output := fmt.Sprintf(
		"Математичне сподівання аварійного недовідпущення: %.0f кВт * год\n"+
			"Математичне сподівання планового недовідпущення: %.0f кВт * год\n"+
			"Математичне сподівання збитків від перервання електропостачання: %.0f грн\n",
		math.Round(MWA), math.Round(MWP), math.Round(M),
	)
	return output
}

func calculator1Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input Input1
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if len(input.Values) != 2 {
		http.Error(w, "Invalid number of inputs", http.StatusBadRequest)
		return
	}
	result := calculateTask1(input.Values[0], input.Values[1])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

func calculator2Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if len(input.Values) != 7 {
		http.Error(w, "Invalid number of inputs", http.StatusBadRequest)
		return
	}
	result := calculateTask2(input.Values[0], input.Values[1], input.Values[2], input.Values[3],
		input.Values[4], input.Values[5], input.Values[6])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/calculator1", calculator1Handler)
	http.HandleFunc("/api/calculator2", calculator2Handler)

	fmt.Println("Server running at http://localhost:8085")
	http.ListenAndServe(":8085", nil)
}
