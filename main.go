// package main

// import (
// 	"database/sql"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"time"

// 	_ "go-echo/docs"

// 	"github.com/go-playground/validator/v10"
// 	"github.com/golang-jwt/jwt/v5" // JWT
// 	"github.com/joho/godotenv"

// 	// "github.com/joho/godotenv"
// 	"go-echo/pkg/middleware"

// 	echojwt "github.com/labstack/echo-jwt/v4"
// 	"github.com/labstack/echo/v4"

// 	// "github.com/labstack/echo/v4"
// 	_ "github.com/lib/pq"
// 	echoSwagger "github.com/swaggo/echo-swagger"
// )

// // =====================
// // JWT Config
// // =====================
// var jwtSecret = []byte("my-super-secret") // ganti di .env biar aman

// // Struct untuk claim JWT
// type JwtCustomClaims struct {
// 	Email string `json:"email"`
// 	jwt.RegisteredClaims
// }

// // =====================
// // Custom Validator
// // =====================
// type CustomValidator struct {
// 	validator *validator.Validate
// }

// func (cv *CustomValidator) Validate(i interface{}) error {
// 	return cv.validator.Struct(i)
// }

// // =====================
// // Struct User
// // =====================
// type User struct {
// 	ID           int    `json:"id"`
// 	Name         string `json:"name" validate:"required"`
// 	Email        string `json:"email" validate:"required,email"`
// 	DepartmentID int    `json:"department_id" validate:"required"`
// }

// // Login request struct
// type LoginRequest struct {
// 	Email string `json:"email" validate:"required,email"`
// }

// // =====================
// // Global DB Connection
// // =====================
// var db *sql.DB

// // =====================
// // Main Function
// // =====================
// func main() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// 	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// 	db = connectDB(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

// 	e := echo.New()
// 	e.Validator = &CustomValidator{validator: validator.New()}

// 	// Swagger
// 	e.GET("/swagger/*", echoSwagger.WrapHandler)

// 	// Login route (no auth)
// 	e.POST("/login", Login)

// 	// Middleware JWT
// 	e.Use(middleware.Logger())
// 	e.Use(middleware.Recover())

// 	r := e.Group("/")

// 	// Protect all routes except login & swagger
// 	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
// 		SigningKey: jwtSecret,
// 	}))

// 	e.Use(echojwt.WithConfig(echojwt.Config{
// 		SigningKey: jwtSecret,
// 	}))

// 	// Routes
// 	e.GET("/", func(c echo.Context) error {
// 		return c.String(http.StatusOK, "Welcome to the User API with JWT")
// 	})
// 	e.GET("users", GetUsers)
// 	e.GET("users/:id", GetUserByID)
// 	e.POST("users", CreateUser)
// 	e.PUT("users/:id", UpdateUser)
// 	e.DELETE("users/:id", DeleteUser)

// 	e.Logger.Fatal(e.Start(":7070"))
// }

// // =====================
// // DB Connection Function
// // =====================
// func connectDB(username string, password string, dbName string) *sql.DB {
// 	connStr := "user=" + username + " password=" + password + " dbname=" + dbName + " sslmode=disable"

// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil
// 	}

// 	if err = db.Ping(); err != nil {
// 		log.Fatal(err)
// 		return nil
// 	}

// 	log.Default().Println("Connected to the database successfully")
// 	return db
// }

// // =====================
// // Handler: Create User
// // =====================
// func CreateUser(c echo.Context) error {
// 	var newUser User
// 	if err := c.Bind(&newUser); err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
// 	}
// 	if err := c.Validate(&newUser); err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
// 	}

// 	query := `INSERT INTO "e-meeting"."User" (user_id, name, email, department_id)
// 	          VALUES ($1, $2, $3, $4)`

// 	_, err := db.Exec(query, newUser.ID, newUser.Name, newUser.Email, newUser.DepartmentID)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
// 	}

// 	return c.JSON(http.StatusCreated, newUser)
// }

// // =====================
// // Handler: Get User By ID
// // =====================
// func GetUserByID(c echo.Context) error {
// 	id := c.Param("id")
// 	idInt, err := strconv.Atoi(id)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
// 	}

// 	query := `SELECT user_id, name, email, department_id FROM "e-meeting"."User" WHERE user_id = $1`
// 	row := db.QueryRow(query, idInt)

// 	var user User
// 	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.DepartmentID); err != nil {
// 		if err == sql.ErrNoRows {
// 			return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
// 		}
// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
// 	}

// 	return c.JSON(http.StatusOK, user)
// }

// // =====================
// // Handler: Get All Users
// // =====================
// func GetUsers(c echo.Context) error {
// 	query := `SELECT user_id, name, email, department_id FROM "e-meeting"."User"`
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
// 	}
// 	defer rows.Close()

// 	var users []User
// 	for rows.Next() {
// 		var user User
// 		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.DepartmentID); err != nil {
// 			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
// 		}
// 		users = append(users, user)
// 	}

// 	return c.JSON(http.StatusOK, users)
// }

// // =====================
// // Handler: Update User
// // =====================
// func UpdateUser(c echo.Context) error {
// 	id := c.Param("id")
// 	idInt, err := strconv.Atoi(id)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
// 	}

// 	var updatedUser User
// 	if err := c.Bind(&updatedUser); err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
// 	}
// 	if err := c.Validate(&updatedUser); err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
// 	}

// 	query := `UPDATE "e-meeting"."User"
// 			  SET name = $1, email = $2, department_id = $3
// 			  WHERE user_id = $4`

// 	result, err := db.Exec(query, updatedUser.Name, updatedUser.Email, updatedUser.DepartmentID, idInt)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
// 	}

// 	updatedUser.ID = idInt
// 	return c.JSON(http.StatusOK, updatedUser)
// }

// // =====================
// // Handler: Delete User
// // =====================
// func DeleteUser(c echo.Context) error {
// 	id := c.Param("id")
// 	idInt, err := strconv.Atoi(id)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
// 	}

// 	query := `DELETE FROM "e-meeting"."User" WHERE user_id = $1`
// 	result, err := db.Exec(query, idInt)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
// 	}

// 	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfully"})
// }

// // =====================
// // Handler: Login (JWT)
// // =====================
// func Login(c echo.Context) error {
// 	var req LoginRequest
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
// 	}
// 	if err := c.Validate(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
// 	}

// 	// (dummy check, bisa ganti query DB user table)
// 	// sekarang siapa aja yang masukin email valid langsung dapet token
// 	claims := &JwtCustomClaims{
// 		Email: req.Email,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)), // token expired 2 menit
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	t, err := token.SignedString(jwtSecret)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate token"})
// 	}

// 	return c.JSON(http.StatusOK, echo.Map{
// 		"token": t,
// 	})
// }
