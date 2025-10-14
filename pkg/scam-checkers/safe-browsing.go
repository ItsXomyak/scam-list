package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Структуры для формирования JSON-запроса к API
type ThreatEntry struct {
	URL string `json:"url"`
}

type ThreatInfo struct {
	ThreatTypes      []string      `json:"threatTypes"`
	PlatformTypes    []string      `json:"platformTypes"`
	ThreatEntryTypes []string      `json:"threatEntryTypes"`
	ThreatEntries    []ThreatEntry `json:"threatEntries"`
}

type ClientInfo struct {
	ClientID      string `json:"clientId"`
	ClientVersion string `json:"clientVersion"`
}

type SafeBrowsingRequest struct {
	Client     ClientInfo `json:"client"`
	ThreatInfo ThreatInfo `json:"threatInfo"`
}

// Структуры для разбора JSON-ответа от API
type Match struct {
	ThreatType string `json:"threatType"`
}

type SafeBrowsingResponse struct {
	Matches []Match `json:"matches"`
}

// checkGoogleSafeBrowsing проверяет URL с помощью Google Safe Browsing API.
// Возвращает строку с результатом и ошибку, если что-то пошло не так.
func checkGoogleSafeBrowsing(urlToCheck string, apiKey string) (string, error) {
	// 1. Формируем URL для запроса
	apiEndpoint := fmt.Sprintf("https://safebrowsing.googleapis.com/v4/threatMatches:find?key=%s", apiKey)

	// 2. Собираем тело запроса, используя наши структуры
	requestBody := SafeBrowsingRequest{
		Client: ClientInfo{
			ClientID:      "scam-list-project-go", // Имя вашего проекта
			ClientVersion: "1.0.0",
		},
		ThreatInfo: ThreatInfo{
			ThreatTypes:      []string{"MALWARE", "SOCIAL_ENGINEERING", "UNWANTED_SOFTWARE", "POTENTIALLY_HARMFUL_APPLICATION"},
			PlatformTypes:    []string{"ANY_PLATFORM"},
			ThreatEntryTypes: []string{"URL"},
			ThreatEntries:    []ThreatEntry{{URL: urlToCheck}},
		},
	}

	// 3. Кодируем тело запроса в JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("ошибка кодирования JSON: %w", err)
	}

	// 4. Отправляем POST-запрос
	resp, err := http.Post(apiEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем, что сервер ответил успешно
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API вернул ошибку %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 5. Читаем и декодируем ответ
	var apiResponse SafeBrowsingResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return "", fmt.Errorf("ошибка декодирования JSON-ответа: %w", err)
	}

	// 6. Анализируем ответ
	if len(apiResponse.Matches) > 0 {
		// Если в поле 'matches' есть данные, значит, угроза найдена
		threatType := apiResponse.Matches[0].ThreatType
		return fmt.Sprintf("🔴 Опасно! Google определил угрозу: %s", threatType), nil
	}

	// Если поле 'matches' пустое, угроз нет
	return "✅ Безопасно (по данным Google)", nil
}

func main() {
	// ❗ Замените на свой ключ API
	const apiKey = "AIzaSyCWAgdUXotMti82Xnxs_aeS-G6RtN2NLS0"

	// --- Примеры использования ---

	// 1. Проверка вредоносного сайта
	maliciousURL := "http://testsafebrowsing.appspot.com/s/phishing.html"
	result1, err := checkGoogleSafeBrowsing(maliciousURL, apiKey)
	if err != nil {
		fmt.Printf("Ошибка при проверке '%s': %v\n", maliciousURL, err)
	} else {
		fmt.Printf("Проверка '%s': %s\n", maliciousURL, result1)
	}

	// 2. Проверка безопасного сайта
	safeURL := "https://google.com"
	result2, err := checkGoogleSafeBrowsing(safeURL, apiKey)
	if err != nil {
		fmt.Printf("Ошибка при проверке '%s': %v\n", safeURL, err)
	} else {
		fmt.Printf("Проверка '%s': %s\n", safeURL, result2)
	}
}