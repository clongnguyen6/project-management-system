project-management/
├── cmd/
│   ├── main.go               # Điểm khởi động ứng dụng
│   ├── migrate.go            # Tệp riêng để chạy migrations
├── config/
│   └── config.go             # Quản lý cấu hình ứng dụng
├── internal/
│   ├── models/
│   │   ├── project.go         # Định nghĩa model Project
│   │   ├── user.go            # Định nghĩa model User
│   │   ├── task.go            # Định nghĩa model Task
│   │   ├── role.go            # Định nghĩa model Role
│   │   ├── project_user.go    # Định nghĩa model ProjectUser
│   ├── handlers/
│   │   ├── project.go         # Handlers cho Project
│   │   ├── user.go            # Handlers cho User
│   │   ├── task.go            # Handlers cho Task
│   │   ├── role.go            # Handlers cho Role
│   │   ├── project_user.go    # Handlers cho ProjectUser
│   ├── router/
│   │   └── router.go          # Định nghĩa router
│   ├── database/
│   │   ├── connection.go      # Kết nối database
│   │   ├── migrations.go      # Quản lý migrations
│   ├── middlewares/
│   │   └── middlewares.go     # Middleware xử lý logs, lỗi, CORS
│   ├── utils/
│       └── helpers.go         # Các hàm phụ trợ
└── go.mod                    # Tệp module Go


// Swagger:
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/http-swagger/v2
go get -u github.com/swaggo/swag
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/http-swagger

Run swagger:
export PATH=$(go env GOPATH)/bin:$PATH
swag init -g cmd/project-management-system-api/main.go --parseDependency --parseInternal



Add auth for swagger:
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.


Example:
// Trong file model
type Project struct {
    gorm.Model
    Name        string    
    Description string    
    Users       []User    `gorm:"many2many:user_projects"`
}

type User struct {
    gorm.Model
    Projects    []Project `gorm:"many2many:user_projects"`
}


// Test
go install rsc.io/uncover@latest
go test -coverprofile=c.out // with coverage
uncover c.out // Show missing coverage