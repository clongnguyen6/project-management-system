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