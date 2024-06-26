package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/bogdanbratsky/golang-itego-api/config"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

// email представляет структуру данных, получаемых из конфига
type email struct {
	from     string
	to       string
	host     string
	port     int
	password string
}

// FormData представляет структуру данных, получаемых из формы
type FormData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Consent bool   `json:"consent"`
}

// Map для хранения времени последнего запроса с каждого IP
var requestTimes = make(map[string]time.Time)
var mu sync.Mutex

const requestInterval = 1 * time.Minute // минимальный интервал между запросами

// SendEmailHandler обрабатывает POST-запросы для отправки email
func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request received:", r.Method, r.URL)

	ip := r.RemoteAddr
	mu.Lock()
	lastRequestTime, exists := requestTimes[ip]
	if exists && time.Since(lastRequestTime) < requestInterval {
		mu.Unlock()
		// Если запрос поступает слишком часто, возвращаем ошибку 429
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}
	// Обновляем время последнего запроса для текущего IP
	requestTimes[ip] = time.Now()
	mu.Unlock()

	// Разбор JSON-запроса
	var data FormData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Проверка на обязательные поля
	if data.Name == "" || (!data.Consent) || (data.Email == "" && data.Phone == "") {
		http.Error(w, "Missing required fields or consent not given", http.StatusBadRequest)
		return
	}

	if err := config.InitConfig(); err != nil {
		log.Println("Не удалось прочитать config.yaml:", err)
	}

	emailConfig := email{
		from:     viper.GetString("email.from"),
		to:       viper.GetString("email.to"),
		host:     viper.GetString("email.host"),
		port:     viper.GetInt("email.port"),
		password: viper.GetString("email.password"),
	}

	// Настройка письма
	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig.from)
	m.SetHeader("To", emailConfig.to)
	m.SetHeader("Subject", "New Form Submission")

	// Установка тела письма в зависимости от того, что было предоставлено
	body := "Name: " + data.Name + "\n"
	if data.Email != "" {
		body += "Email: " + data.Email + "\n"
	}
	if data.Phone != "" {
		body += "Phone: " + data.Phone + "\n"
	}
	m.SetBody("text/plain", body)

	// Настройка SMTP
	d := gomail.NewDialer(
		emailConfig.host,
		emailConfig.port,
		emailConfig.from,
		emailConfig.password,
	)

	// Отправка письма
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Could not send email: %v", err)
		http.Error(w, "Could not send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}
