package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//DBConfig holds database connection object
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func init() {
	SetupDatabase()
	SetupUserRole()
}

//Model objects for gorm
type Model struct {
	ID        uint64     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

//Subsidiaries structures the list of subsidiaries
type Subsidiaries struct {
	Model
	Subsidiary string `gorm:"varchar(100)" json:"subsidiary"`
}

//RoleValues holds the values for user roles
type RoleValues struct {
	VMSAdminOfficer  uint64
	HROfficer        uint64
	PASTeamLead      uint64
	FrontDeskOfficer uint64
	TaskAdmin        uint64
}

//Departments holds the list of deparments and subsidiaries belonging
type Departments struct {
	Model
	Department   string `gorm:"type:varchar(100)" json:"department"`
	Subsidiary   string `gorm:"type:varchar(100)" json:"subsidiary"`
	SubsidiaryID uint64 `gorm:"type:int(10)" json:"subsidiary_id"`
	Head         string `gorm:"type:varchar(10)" json:"head"`
	HeadID       uint64 `gorm:"type:int(10)" json:"head_id"`
}

//User struct shows models for users
type User struct {
	Model
	FullName     string `gorm:"type:varchar(100)" json:"full_name" form:"full_name"`
	Password     string `gorm:"type:varchar(100)" json:"password" form:"password"`
	Department   string `gorm:"type:varchar(100)" json:"department"`
	DepartmentID uint64 `gorm:"type:int(10)" json:"department_id"`
	JobTitle     string `gorm:"type:varchar(100)" json:"job_title"`
	Gender       string `gorm:"type:varchar(100)" json:"gender"`
	Location     string `gorm:"type:varchar(100)" json:"location"`
	Subsidiary   string `gorm:"type:varchar(100)" json:"subsidiary"`
	SubsidiaryID uint64 `gorm:"type:int(10)" json:"subsidiary_id"`
	LinkedIn     string `gorm:"type:varchar(100)" json:"linked_in"`
	Number       string `gorm:"type:varchar(100)"  json:"number"`
	Email        string `gorm:"type:varchar(100); unique_index" json:"email" form:"email"`
	Role         uint64 `gorm:"type:int(10)" json:"role"`
	Image        string `gorm:"type:varchar(250)" json:"image"`
}

//AttemptedLogin logs any user trying to log in
type AttemptedLogin struct {
	gorm.Model
	Username string `gorm:"type:varchar(100)" json:"username"`
	Status   string `gorm:"type:varchar(100)" json:"status"`
	Message  string `gorm:"type:varchar(100)" json:"message"`
}

//Roles shows a model of user roles and permission
type Roles struct {
	Model
	Code   uint64 `gorm:"type:int(10)" json:"role_id"`
	Role   string `gorm:"type:varchar(150)" json:"role"`
	User   string `gorm:"type:varchar(100)" json:"user"`
	UserID uint64 `gorm:"type:int(10)" json:"user_id"`
}

//PasswordRecoveryData holds needed data to reset password
type PasswordRecoveryData struct {
	Model
	Email string `gorm:"type:varchar(100)" json:"email"`
	Code  string `gorm:"type:varchar(100)" json:"code"`
}

//ResetPasswordData holds the needed data to reset a password
type ResetPasswordData struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

//ValidateRole holds data needed to validate a user role
type ValidateRole struct {
	UserID   uint64 `json:"user_id"`
	RoleCode uint64 `json:"role_code"`
}
