package model

import (
	"time"

	"github.com/gsq/music_bakcend_micorservice/pkg/config"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB

type Song struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Artist    string    `gorm:"size:100;not null" json:"artist"`
	Album     string    `gorm:"size:100" json:"album"`
	Duration  uint      `gorm:"default:0" json:"duration"` // 秒
	URL       string    `gorm:"size:500;not null" json:"url"`
	CoverURL  string    `gorm:"size:500" json:"cover_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`                    // 主键，自增
	Username  string    `gorm:"size:50;not null;unique" json:"username"` // 用户名，唯一
	Email     string    `gorm:"size:100" json:"email"`                   // 邮箱，可选
	Password  string    `gorm:"size:255;not null" json:"password"`       // 密码（建议加密）
	CreatedAt time.Time `json:"created_at"`                              // 注册时间
	UpdatedAt time.Time `json:"updated_at"`                              // 更新时间
}

func init() {
	config.Connect()
	db = config.GetDB()
}

func GetSongs() []Song {
	var Songs []Song
	db.Find(&Songs)
	return Songs
}

func GetSongByName(name string) (*Song, error) {
	var song Song
	result := db.Where("title = ?", name).First(&song)
	if result.Error != nil {
		return nil, result.Error
	}
	return &song, nil
}

func SearchSongByKeyword(keyword string) ([]Song, error) {
	var songs []Song
	if err := db.Where("title Like ?", "%"+keyword+"%").Find(&songs).Error; err != nil {
		return nil, err
	}
	//fmt.Println(songs)
	return songs, nil
}

func CreateUser(username, email, password string) (*User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &User{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
