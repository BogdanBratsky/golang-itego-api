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

type email struct {
	from     string
	to       string
	host     string
	port     int
	password string
}

type FormData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Wish    string `json:"wish"`
	Problem string `json:"problem"`
	Consent bool   `json:"consent"`
}

var requestTimes = make(map[string]time.Time)
var mu sync.Mutex

const requestInterval = 1 * time.Minute

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request received:", r.Method, r.URL)

	ip := r.RemoteAddr
	mu.Lock()
	lastRequestTime, exists := requestTimes[ip]
	if exists && time.Since(lastRequestTime) < requestInterval {
		mu.Unlock()
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}
	requestTimes[ip] = time.Now()
	mu.Unlock()

	var data FormData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !data.Consent || (data.Name == "" && data.Email == "" && data.Phone == "" && data.Wish == "" && data.Problem == "") {
		http.Error(w, "Consent not given or no information provided", http.StatusBadRequest)
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

	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig.from)
	m.SetHeader("To", emailConfig.to)
	m.SetHeader("Subject", "New Form Submission")

	body := ""
	if data.Name != "" {
		body += "Name: " + data.Name + "\n"
	}
	if data.Email != "" {
		body += "Email: " + data.Email + "\n"
	}
	if data.Phone != "" {
		body += "Phone: " + data.Phone + "\n"
	}
	if data.Wish != "" {
		body += "Wish: " + data.Wish + "\n"
	}
	if data.Problem != "" {
		body += "Problem: " + data.Problem + "\n"
	}
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(
		emailConfig.host,
		emailConfig.port,
		emailConfig.from,
		emailConfig.password,
	)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Could not send email: %v", err)
		http.Error(w, "Could not send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}
