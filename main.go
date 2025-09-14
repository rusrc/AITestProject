package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Member представляет участника команды
type Member struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Role         string        `json:"role"`
	Avatar       string        `json:"avatar"`
	SortOrder    int           `json:"sort_order"`
	Achievements []Achievement `json:"achievements"`
	CreatedAt    time.Time     `json:"created_at"`
}

// Achievement представляет ачивку
type Achievement struct {
	ID        int       `json:"id"`
	MemberID  int       `json:"member_id"`
	Image     string    `json:"image"`
	Category  string    `json:"category"` // "positive" или "negative"
	CreatedAt time.Time `json:"created_at"`
}

// CSV хранилище данных
type DataStore struct {
	members           map[int]*Member
	achievements      map[int]*Achievement
	nextMemberID      int
	nextAchievementID int
	mutex             sync.RWMutex
	membersFile       string
	achievementsFile  string
}

// Создание нового хранилища
func NewDataStore() *DataStore {
	ds := &DataStore{
		members:           make(map[int]*Member),
		achievements:      make(map[int]*Achievement),
		nextMemberID:      1,
		nextAchievementID: 1,
		membersFile:       "data/members.csv",
		achievementsFile:  "data/achievements.csv",
	}

	// Загружаем данные из CSV файлов
	ds.loadFromCSV()
	return ds
}

// Загрузка данных из CSV файлов
func (ds *DataStore) loadFromCSV() {
	// Создаем папку data если её нет
	os.MkdirAll("data", 0755)

	// Загружаем участников
	ds.loadMembersFromCSV()

	// Загружаем ачивки
	ds.loadAchievementsFromCSV()
}

// Загрузка участников из CSV
func (ds *DataStore) loadMembersFromCSV() {
	file, err := os.Open(ds.membersFile)
	if err != nil {
		// Файл не существует, создаем пустой
		ds.saveMembersToCSV()
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Ошибка чтения CSV участников: %v", err)
		return
	}

	// Пропускаем заголовок
	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 6 {
			continue
		}

		id, _ := strconv.Atoi(record[0])
		sortOrder, _ := strconv.Atoi(record[4])
		createdAt, _ := time.Parse(time.RFC3339, record[5])

		member := &Member{
			ID:           id,
			Name:         record[1],
			Role:         record[2],
			Avatar:       record[3],
			SortOrder:    sortOrder,
			Achievements: []Achievement{},
			CreatedAt:    createdAt,
		}

		ds.members[id] = member
		if id >= ds.nextMemberID {
			ds.nextMemberID = id + 1
		}
	}
}

// Загрузка ачивок из CSV
func (ds *DataStore) loadAchievementsFromCSV() {
	file, err := os.Open(ds.achievementsFile)
	if err != nil {
		// Файл не существует, создаем пустой
		ds.saveAchievementsToCSV()
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Ошибка чтения CSV ачивок: %v", err)
		return
	}

	// Пропускаем заголовок
	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 5 {
			continue
		}

		id, _ := strconv.Atoi(record[0])
		memberID, _ := strconv.Atoi(record[1])
		createdAt, _ := time.Parse(time.RFC3339, record[4])

		achievement := &Achievement{
			ID:        id,
			MemberID:  memberID,
			Image:     record[2],
			Category:  record[3],
			CreatedAt: createdAt,
		}

		ds.achievements[id] = achievement
		if id >= ds.nextAchievementID {
			ds.nextAchievementID = id + 1
		}

		// Добавляем ачивку к участнику
		if member, exists := ds.members[memberID]; exists {
			member.Achievements = append(member.Achievements, *achievement)
		}
	}
}

// Сохранение участников в CSV
func (ds *DataStore) saveMembersToCSV() {
	file, err := os.Create(ds.membersFile)
	if err != nil {
		log.Printf("Ошибка создания CSV участников: %v", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Записываем заголовок
	writer.Write([]string{"id", "name", "role", "avatar", "sort_order", "created_at"})

	// Записываем данные
	for _, member := range ds.members {
		writer.Write([]string{
			strconv.Itoa(member.ID),
			member.Name,
			member.Role,
			member.Avatar,
			strconv.Itoa(member.SortOrder),
			member.CreatedAt.Format(time.RFC3339),
		})
	}
}

// Сохранение ачивок в CSV
func (ds *DataStore) saveAchievementsToCSV() {
	file, err := os.Create(ds.achievementsFile)
	if err != nil {
		log.Printf("Ошибка создания CSV ачивок: %v", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Записываем заголовок
	writer.Write([]string{"id", "member_id", "image", "category", "created_at"})

	// Записываем данные
	for _, achievement := range ds.achievements {
		writer.Write([]string{
			strconv.Itoa(achievement.ID),
			strconv.Itoa(achievement.MemberID),
			achievement.Image,
			achievement.Category,
			achievement.CreatedAt.Format(time.RFC3339),
		})
	}
}

// Добавление участника
func (ds *DataStore) AddMember(name, role, avatar string) *Member {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	member := &Member{
		ID:           ds.nextMemberID,
		Name:         name,
		Role:         role,
		Avatar:       avatar,
		SortOrder:    ds.nextMemberID, // Используем ID как порядок сортировки
		Achievements: []Achievement{},
		CreatedAt:    time.Now(),
	}

	ds.members[ds.nextMemberID] = member
	ds.nextMemberID++

	// Сохраняем в CSV
	ds.saveMembersToCSV()

	return member
}

// Получение всех участников
func (ds *DataStore) GetMembers() []*Member {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	members := make([]*Member, 0, len(ds.members))
	for _, member := range ds.members {
		// Копируем участника с его ачивками
		memberCopy := *member
		memberCopy.Achievements = make([]Achievement, len(member.Achievements))
		copy(memberCopy.Achievements, member.Achievements)
		members = append(members, &memberCopy)
	}

	// Сортируем по SortOrder
	for i := 0; i < len(members)-1; i++ {
		for j := i + 1; j < len(members); j++ {
			if members[i].SortOrder > members[j].SortOrder {
				members[i], members[j] = members[j], members[i]
			}
		}
	}

	return members
}

// Добавление ачивки
func (ds *DataStore) AddAchievement(memberID int, image, category string) (*Achievement, error) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	// Проверяем существование участника
	if _, exists := ds.members[memberID]; !exists {
		return nil, fmt.Errorf("участник с ID %d не найден", memberID)
	}

	achievement := &Achievement{
		ID:        ds.nextAchievementID,
		MemberID:  memberID,
		Image:     image,
		Category:  category,
		CreatedAt: time.Now(),
	}

	ds.achievements[ds.nextAchievementID] = achievement
	ds.nextAchievementID++

	// Добавляем ачивку к участнику
	ds.members[memberID].Achievements = append(ds.members[memberID].Achievements, *achievement)

	// Сохраняем в CSV
	ds.saveAchievementsToCSV()

	return achievement, nil
}

// Создание папки для загрузок
func ensureUploadDir() error {
	dirs := []string{"uploads/avatars", "uploads/achievements"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// Сохранение загруженного файла
func saveUploadedFile(fileBytes []byte, filename string) (string, error) {
	// Создаем уникальное имя файла
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".jpg" // По умолчанию jpg
	}
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Определяем папку по типу файла
	var uploadDir string
	if strings.Contains(filename, "avatar") || strings.Contains(filename, "member") {
		uploadDir = "uploads/avatars"
	} else {
		uploadDir = "uploads/achievements"
	}

	filePath := filepath.Join(uploadDir, newFilename)

	// Сохраняем файл
	if err := os.WriteFile(filePath, fileBytes, 0644); err != nil {
		return "", err
	}

	return filePath, nil
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Обработчик получения всех участников
func getMembersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	members := dataStore.GetMembers()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    members,
	})
}

// Обработчик создания участника
func createMemberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	role := r.FormValue("role")

	if name == "" || role == "" {
		http.Error(w, "Имя и роль обязательны", http.StatusBadRequest)
		return
	}

	// Обработка загруженного файла
	var avatarPath string
	if file, _, err := r.FormFile("image"); err == nil {
		defer file.Close()

		// Читаем файл
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
			return
		}

		// Сохраняем файл
		avatarPath, err = saveUploadedFile(fileBytes, "avatar.jpg")
		if err != nil {
			http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
			return
		}
	} else {
		// Используем дефолтный аватар
		avatarPath = "assets/images/avatar.svg"
	}

	// Создаем участника
	member := dataStore.AddMember(name, role, avatarPath)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    member,
	})
}

// Обработчик добавления ачивки
func addAchievementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	memberIDStr := r.FormValue("member_id")
	category := r.FormValue("category")

	if memberIDStr == "" || category == "" {
		http.Error(w, "ID участника и категория обязательны", http.StatusBadRequest)
		return
	}

	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		http.Error(w, "Неверный ID участника", http.StatusBadRequest)
		return
	}

	// Обработка загруженного файла
	var imagePath string
	if file, _, err := r.FormFile("image"); err == nil {
		defer file.Close()

		// Читаем файл
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
			return
		}

		// Сохраняем файл
		imagePath, err = saveUploadedFile(fileBytes, "achievement.jpg")
		if err != nil {
			http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Файл изображения обязателен", http.StatusBadRequest)
		return
	}

	// Добавляем ачивку
	achievement, err := dataStore.AddAchievement(memberID, imagePath, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    achievement,
	})
}

// Глобальное хранилище данных
var dataStore *DataStore

func main() {
	// Создаем хранилище данных
	dataStore = NewDataStore()

	// Создаем папки для загрузок и данных
	if err := ensureUploadDir(); err != nil {
		log.Fatal("Ошибка создания папок:", err)
	}

	// Создаем папку для данных
	os.MkdirAll("data", 0755)

	// Создаем роутер
	r := mux.NewRouter()

	// Добавляем CORS middleware
	r.Use(corsMiddleware)

	// Регистрируем роуты
	r.HandleFunc("/api/members", getMembersHandler).Methods("GET")
	r.HandleFunc("/api/members", createMemberHandler).Methods("POST")
	r.HandleFunc("/api/achievements", addAchievementHandler).Methods("POST")

	// Статические файлы
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))

	// Запускаем сервер
	port := ":5000"
	fmt.Printf("Сервер запущен на http://localhost%s\n", port)
	fmt.Println("Доступные эндпоинты:")
	fmt.Println("  GET  /api/members - получить всех участников")
	fmt.Println("  POST /api/members - создать участника")
	fmt.Println("  POST /api/achievements - добавить ачивку")

	log.Fatal(http.ListenAndServe(port, r))
}
