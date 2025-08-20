package model

import (
	"time"

	"github.com/gsq/music_bakcend_micorservice/database"
	"github.com/jinzhu/gorm"
)

// PlayHistory 播放历史模型
type PlayHistory struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	UserID   uint      `gorm:"not null" json:"user_id"`
	SongID   uint      `gorm:"not null" json:"song_id"`
	PlayedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"played_at"`
	Song     Song      `gorm:"foreignKey:SongID" json:"song"`
}

// 自动迁移表结构
func init() {
	db = database.GetDB()
	db.AutoMigrate(&PlayHistory{})
}

// CreatePlayHistory 创建或更新播放历史记录
func CreatePlayHistory(userID, songID uint) error {
	var history PlayHistory
	// 查找用户是否已播放过该歌曲
	err := db.Where("user_id = ? AND song_id = ?", userID, songID).First(&history).Error

	if err != nil {
		// 如果不存在，创建新记录
		if gorm.IsRecordNotFoundError(err) {
			history = PlayHistory{
				UserID:   userID,
				SongID:   songID,
				PlayedAt: time.Now(),
			}
			return db.Create(&history).Error
		}
		// 其他错误
		return err
	}

	// 如果已存在，更新播放时间
	history.PlayedAt = time.Now()
	return db.Save(&history).Error
}

// GetUserPlayHistory 获取用户播放历史（每个歌曲只保留最新一条记录）
func GetUserPlayHistory(userID uint, limit int) ([]PlayHistory, error) {
	var histories []PlayHistory

	// 使用子查询获取每个用户每个歌曲的最新播放记录
	err := db.Preload("Song").
		Where("id IN (SELECT MAX(id) FROM play_histories WHERE user_id = ? GROUP BY song_id)", userID).
		Order("played_at DESC").
		Limit(limit).
		Find(&histories).Error
	if err != nil {
		return nil, err
	}

	return histories, nil
}

// DeletePlayHistory 删除用户的播放历史记录
func DeletePlayHistory(userID, songID uint) error {
	// 删除用户的特定歌曲播放历史
	return db.Where("user_id = ? AND song_id = ?", userID, songID).Delete(&PlayHistory{}).Error
}
