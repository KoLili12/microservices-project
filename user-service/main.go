package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "sync"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var (
    users   = make(map[int]User)
    usersMu sync.RWMutex
    nextID  = 1
)

func main() {
    // Инициализация с тестовыми данными
    initUsers()

    http.HandleFunc("/users", handleUsers)
    http.HandleFunc("/users/", handleUserByID)
    http.HandleFunc("/health", handleHealth)

    port := getEnv("PORT", "8080")
    log.Printf("User Service запущен на порту %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initUsers() {
    users[1] = User{ID: 1, Name: "Иван Иванов", Email: "ivan@example.com"}
    users[2] = User{ID: 2, Name: "Мария Петрова", Email: "maria@example.com"}
    nextID = 3
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    switch r.Method {
    case "GET":
        getAllUsers(w, r)
    case "POST":
        createUser(w, r)
    default:
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
    usersMu.RLock()
    defer usersMu.RUnlock()

    userList := make([]User, 0, len(users))
    for _, user := range users {
        userList = append(userList, user)
    }

    json.NewEncoder(w).Encode(userList)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    usersMu.Lock()
    user.ID = nextID
    nextID++
    users[user.ID] = user
    usersMu.Unlock()

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func handleUserByID(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var id int
    fmt.Sscanf(r.URL.Path, "/users/%d", &id)

    if r.Method == "GET" {
        getUserByID(w, id)
    } else {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

func getUserByID(w http.ResponseWriter, id int) {
    usersMu.RLock()
    user, exists := users[id]
    usersMu.RUnlock()

    if !exists {
        http.Error(w, "Пользователь не найден", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}