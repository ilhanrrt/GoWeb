package logger

import (
	"fmt"
	"os"
	"time"
)

const (
	// ANSI escape kodları (sadece terminalde çalışır)
	reset = "\033[0m"
	green = "\033[32m" // Yeşil renk, olumlu işlemler için
	red   = "\033[31m" // Kırmızı renk, olumsuz işlemler için
)

// Loglama işlemini yapan fonksiyon
func LogAction(adminID int, adminType string, action string) {
	logFileName := "admin_logs.txt"
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Dosya açma hatası:", err)
		return
	}
	defer file.Close()

	logMessage := fmt.Sprintf("Tarih: %s, Kullanıcı ID: %d, Admin Statü: %s, İşlem: %s\n", time.Now().Format("2006-01-02 15:04:05"), adminID, adminType, action)

	if _, err := file.WriteString(logMessage); err != nil {
		fmt.Println("Dosyaya yazma hatası:", err)
	} else {
		fmt.Println("Log kaydedildi:", logMessage)
	}
}

// Loglama işlemini yapan fonksiyon
func LogSaveMessage(action string) {
	logFileName := "admin_logs.txt"
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Dosya açma hatası:", err)
		return
	}
	defer file.Close()

	logMessage := fmt.Sprintf("Tarih: %s, İşlem: %s\n", time.Now().Format("2006-01-02 15:04:05"), action)

	if _, err := file.WriteString(logMessage); err != nil {
		fmt.Println("Dosyaya yazma hatası:", err)
	} else {
		fmt.Println("Log kaydedildi:", logMessage)
	}
}

// LogActionPanel: Panel işlemlerini loglar
func LogActionPanel(adminID int, adminType string, action string, isSuccess bool) {
	logFileName := "panel_logs.txt"
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Dosya açma hatası:", err)
		return
	}
	defer file.Close()

	// Başarılı ve başarısız işlem durumları için başlık
	status := "INFO"
	color := green

	if !isSuccess {
		status = "ERROR"
		color = red
	}

	// Mesajın başına status ekleniyor
	logMessage := fmt.Sprintf("[%s] Tarih: %s, Kullanıcı ID: %d, Admin Statü: %s, İşlem: %s\n",
		status, time.Now().Format("2006-01-02 15:04:05"), adminID, adminType, action)

	// Terminale renkli log yazdır
	fmt.Printf("%s%s%s", color, logMessage, reset)

	// Dosyaya status ile yazma (renksiz, ama status bilgisi içeriyor)
	if _, err := file.WriteString(logMessage); err != nil {
		fmt.Println("Dosyaya yazma hatası:", err)
	}
}

// LogSaveMessagePanel: Paneldeki genel işlemleri loglar
func LogSaveMessagePanel(action string, isSuccess bool) {
	logFileName := "panel_logs.txt"
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Dosya açma hatası:", err)
		return
	}
	defer file.Close()

	// Başarılı ve başarısız işlem durumları için başlık
	status := "INFO"
	color := green

	if !isSuccess {
		status = "ERROR"
		color = red
	}

	// Mesajın başına status ekleniyor
	logMessage := fmt.Sprintf("[%s] Tarih: %s, İşlem: %s\n", status, time.Now().Format("2006-01-02 15:04:05"), action)

	// Terminale renkli log yazdır
	fmt.Printf("%s%s%s", color, logMessage, reset)

	// Dosyaya status ile yazma (renksiz, ama status bilgisi içeriyor)
	if _, err := file.WriteString(logMessage); err != nil {
		fmt.Println("Dosyaya yazma hatası:", err)
	}
}
