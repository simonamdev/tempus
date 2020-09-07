package db

import (
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // needed for sqlite
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Projects []Project
}

func (u *User) BeforeCreate() (err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	u.Password = string(bytes)
	return
}

type Project struct {
	gorm.Model
	Title      string `gorm:"not null"`
	Entries    []ProjectEntry
	EntryTypes []ProjectEntryType
	UserID     uint
}

type ProjectEntryType struct {
	gorm.Model
	Title     string `gorm:"unique;not null"`
	ProjectID uint
	Colour    string
}

type ProjectEntry struct {
	gorm.Model
	EntryTypeID uint
	ProjectID   uint
	OpenTime    time.Time
	CloseTime   *time.Time
}

func Open() {
	var err error
	DB, err = gorm.Open("sqlite3", "data/tempus.db")
	if err != nil {
		log.Println("Unable to open DB!")
		panic(err)
	}
	DB.LogMode(true)
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Project{})
	DB.AutoMigrate(&ProjectEntryType{})
	DB.AutoMigrate(&ProjectEntry{})
}

func Clear() {
	info, err := os.Stat("./tempus.db")
	if err != nil {
		log.Println(err.Error())
		return
	}
	if !info.IsDir() {
		err := os.Remove("./tempus.db")
		if err != nil {
			panic(err)
		}
	}
}

const isoParseFormat = "2006-01-02 15:04:05-07:00"

func Seed() {
	// TODO: set default admin password in env
	baseUser := User{Name: "Admin", Email: "simon@agius-muscat.net", Password: "admin123"}

	notFound := DB.Where("email = ?", baseUser.Email).Find(&User{}).RecordNotFound()
	if notFound {
		DB.Create(&baseUser)
	}
	count := 0
	err := DB.Model(&Project{}).Count(&count).Error
	if err != nil {
		panic(err)
	}
	if count == 0 {
		project := Project{Title: "Test Project", UserID: 1}
		err := DB.Create(&project).Error
		if err != nil {
			panic(err)
		}
	}

}
