package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "sync"
)

type Order struct {
    ID       int     `json:"id"`
    UserID   int     `json:"user_id"`
    Product  string  `json:"product"`
    Amount   float64 `json:"amount"`
    UserName string  `json:"user_name,omitempty"`
}

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var (
    orders   = make(map[int]Order)
    ordersMu sync.RWMutex
    nextID   = 1
)

func main() {
    // Инициализация с тестовыми данными
    initOrders()

    http.HandleFunc("/orders", handleOrders)
    http.HandleFunc("/orders/", handleOrderByID)
    http.HandleFunc("/health", handleHealth)

    port := getEnv("PORT", "8081")
    log.Printf("Order Service запущен на порту %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initOrders() {
    orders[1] = Order{ID: 1, UserID: 1, Product: "Ноутбук", Amount: 75000}
    orders[2] = Order{ID: 2, UserID: 2, Product: "Смартфон", Amount: 35000}
    nextID = 3
}

func handleOrders(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    switch r.Method {
    case "GET":
        getAllOrders(w, r)
    case "POST":
        createOrder(w, r)
    default:
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

func getAllOrders(w http.ResponseWriter, r *http.Request) {
    ordersMu.RLock()
    orderList := make([]Order, 0, len(orders))
    for _, order := range orders {
        orderList = append(orderList, order)
    }
    ordersMu.RUnlock()

    // Обогащаем данные информацией о пользователях
    for i := range orderList {
        user, err := getUserFromService(orderList[i].UserID)
        if err == nil {
            orderList[i].UserName = user.Name
        }
    }

    json.NewEncoder(w).Encode(orderList)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
    var order Order
    if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ordersMu.Lock()
    order.ID = nextID
    nextID++
    orders[order.ID] = order
    ordersMu.Unlock()

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(order)
}

func handleOrderByID(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var id int
    fmt.Sscanf(r.URL.Path, "/orders/%d", &id)

    if r.Method == "GET" {
        getOrderByID(w, id)
    } else {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

func getOrderByID(w http.ResponseWriter, id int) {
    ordersMu.RLock()
    order, exists := orders[id]
    ordersMu.RUnlock()

    if !exists {
        http.Error(w, "Заказ не найден", http.StatusNotFound)
        return
    }

    // Получаем информацию о пользователе
    user, err := getUserFromService(order.UserID)
    if err == nil {
        order.UserName = user.Name
    }

    json.NewEncoder(w).Encode(order)
}

func getUserFromService(userID int) (*User, error) {
    userServiceURL := getEnv("USER_SERVICE_URL", "http://localhost:8080")
    url := fmt.Sprintf("%s/users/%d", userServiceURL, userID)

    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Ошибка запроса к User Service: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("пользователь не найден")
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var user User
    if err := json.Unmarshal(body, &user); err != nil {
        return nil, err
    }

    return &user, nil
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