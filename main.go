package main

import (
    "fmt"
    "io"
    "log"
    "net"
    "net/http"
    "sync"

    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
)

type ServerConfig struct {
    Path string `mapstructure:"path"`
    IP   string `mapstructure:"ip"`
    Port string `mapstructure:"port"`
}

type Config struct {
    LogLevel   string         `mapstructure:"log_level"`
    HTTPPort   int            `mapstructure:"http_port"`
    UDPAddress string         `mapstructure:"udp_address"`
    ProxyKey   string         `mapstructure:"proxy_key"`
    Servers    []ServerConfig `mapstructure:"servers"`
}

var (
    config     Config
    configLock sync.RWMutex
)

func main() {
    loadConfig()

    for _, server := range config.Servers {
        http.HandleFunc(server.Path, makeHandler(server))
    }

    addr := fmt.Sprintf(":%d", config.HTTPPort)
    logInfo("Listening on %s...", addr)
    log.Fatal(http.ListenAndServe(addr, nil))
}

func loadConfig() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")

    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("Error reading config: %v", err)
    }

    if err := viper.Unmarshal(&config); err != nil {
        log.Fatalf("Error unmarshaling config: %v", err)
    }

    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        logInfo("Config file changed, reloading...")
        var newConfig Config
        if err := viper.Unmarshal(&newConfig); err != nil {
            log.Printf("Error unmarshaling new config: %v", err)
            return
        }
        configLock.Lock()
        config = newConfig
        configLock.Unlock()
    })
}

func makeHandler(server ServerConfig) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        logDebug("[%s] Received %s request from %s", server.Path, r.Method, r.RemoteAddr)

        defer r.Body.Close()
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Failed to read body", http.StatusInternalServerError)
            return
        }

        configLock.RLock()
        udpAddr := config.UDPAddress
        proxyKey := config.ProxyKey
        configLock.RUnlock()

        bodyStr := string(body)
        if len(bodyStr) < 2 || bodyStr[:2] != "L " {
            bodyStr = "L " + bodyStr
        }

        logDebug("[%s] Log line: %s", server.Path, bodyStr)

        logLine := fmt.Sprintf(" PROXY Key=%s %s:%sPROXY \n%s", proxyKey,server.IP, server.Port, bodyStr)

        conn, err := net.Dial("udp", udpAddr)
        if err != nil {
            log.Printf("[%s] UDP dial error: %v", server.Path, err)
            return
        }
        defer conn.Close()

        _, err = conn.Write([]byte(logLine))
        if err != nil {
            log.Printf("[%s] UDP write error: %v", server.Path, err)
        } else {
            logDebug("[%s] Sent to %s", server.Path, udpAddr)
        }
    }
}

func logDebug(format string, v ...interface{}) {
    configLock.RLock()
    defer configLock.RUnlock()
    if config.LogLevel == "debug" {
        log.Printf("[DEBUG] "+format, v...)
    }
}

func logInfo(format string, v ...interface{}) {
    configLock.RLock()
    defer configLock.RUnlock()
    if config.LogLevel == "debug" || config.LogLevel == "info" {
        log.Printf("[INFO] "+format, v...)
    }
}
