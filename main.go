package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/fatih/color"
    "github.com/shirou/gopsutil/v3/cpu"
    "github.com/shirou/gopsutil/v3/disk"
    "github.com/shirou/gopsutil/v3/host"
    "github.com/shirou/gopsutil/v3/mem"
    "github.com/vbauerster/mpb/v7"
    "github.com/vbauerster/mpb/v7/decor"
)

// Сервисы и связанные с ними хосты
var services = map[string][]string{
    "Service A": {"host1.example.com", "host2.example.com"},
    "Service B": {"host3.example.com", "host4.example.com"},
}

// Информация о сертификате
type CertificateInfo struct {
    Host        string    `json:"host"`
    Issuer      string    `json:"issuer"`
    ValidFrom   time.Time `json:"valid_from"`
    ValidUntil  time.Time `json:"valid_until"`
    SerialNumber string   `json:"serial_number"`
    SystemTime  string    `json:"system_time"`
    Uptime      uint64    `json:"uptime"`
    CPUUsage    float64   `json:"cpu_usage"`
    RAMUsage    float64   `json:"ram_usage"`
}

// Моковые данные для примера
func getCertificateInfo(host string) CertificateInfo {
    info := hostInfo()
    return CertificateInfo{
        Host:       host,
        Issuer:     "Let's Encrypt",
        ValidFrom:  time.Now().Add(-time.Hour * 24 * 30),
        ValidUntil: time.Now().Add(time.Hour * 24 * 60),
        SerialNumber: "1234567890",
        SystemTime: info.SystemTime,
        Uptime:     info.Uptime,
        CPUUsage:   info.CPUUsage,
        RAMUsage:   info.RAMUsage,
    }
}

type HostInfo struct {
    SystemTime string  `json:"system_time"`
    Uptime     uint64  `json:"uptime"`
    CPUUsage   float64 `json:"cpu_usage"`
    RAMUsage   float64 `json:"ram_usage"`
}

func hostInfo() HostInfo {
    sysTime, _ := host.BootTime()
    bootTime := time.Unix(int64(sysTime), 0)
    upTime := time.Since(bootTime).Seconds()
    cpuPercent, _ := cpu.Percent(time.Second, false)
    memInfo, _ := mem.VirtualMemory()
    return HostInfo{
        SystemTime: time.Now().Format(time.RFC1123),
        Uptime:     uint64(upTime),
        CPUUsage:   cpuPercent[0],
        RAMUsage:   memInfo.UsedPercent,
    }
}

// Обработчик запроса для получения информации о сервисах
func servicesHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(services)
}

// Обработчик запроса для получения информации о сертификате по хосту
func certificateHandler(w http.ResponseWriter, r *http.Request) {
    host := r.URL.Query().Get("host")
    if host == "" {
        http.Error(w, "Host is required", http.StatusBadRequest)
        return
    }
    certInfo := getCertificateInfo(host)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(certInfo)
}

// Обработчик запроса для импорта данных из JSON файла
func importServicesHandler(w http.ResponseWriter, r *http.Request) {
    var importedServices map[string][]string
    err := json.NewDecoder(r.Body).Decode(&importedServices)
    if err != nil {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }
    services = importedServices
    w.WriteHeader(http.StatusOK)
}

// Обработчик запроса для экспорта данных в JSON файл
func exportServicesHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(services)
}

// Функция для демонстрации прогресса
func simulateProgress() {
    p := mpb.New(mpb.WithWaitGroup(nil))
    total := 100
    bar := p.AddBar(int64(total),
        mpb.PrependDecorators(
            decor.Name("Progress"),
            decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
        ),
        mpb.AppendDecorators(
            decor.Percentage(),
        ),
    )
    for i := 0; i < total; i++ {
        time.Sleep(100 * time.Millisecond)
        bar.Increment()
    }
    p.Wait()
}

func main() {
    color.Cyan("Starting SSL Certificate Checker...")
    go simulateProgress() // Запуск анимации прогресса в фоновом потоке

    http.HandleFunc("/services", servicesHandler)
    http.HandleFunc("/certificate", certificateHandler)
    http.HandleFunc("/import", importServicesHandler)
    http.HandleFunc("/export", exportServicesHandler)
    fmt.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}
