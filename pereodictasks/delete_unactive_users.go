package pereodictasks

import (
	"log"
	"time"

	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/robfig/cron/v3"
)

func StartPereodicTasks() error {
    c := cron.New()
    _, err := c.AddFunc("@every 10m", func() {
        if err := deleteInactiveUsers(); err != nil {
            log.Printf("Error deleting inactive users: %v", err)
        } else {
            log.Println("Inactive users deleted successfully")
        }
    })
    if err != nil {
        log.Fatalf("could not schedule task: %v", err)
    }
    log.Println("Scheduled task to delete inactive users every 10 minute")
    c.Start()
    return nil
}

func deleteInactiveUsers() error {
    log.Println("Running deleteInactiveUsers task")
    criticalTime := time.Now().Add(-24 * time.Hour)
    var unActiveUsers []models.User
    if err := config.DB.Model(&models.User{}).Where("is_active = ? AND created_at < ?", false, criticalTime).Find(&unActiveUsers).Error; err != nil {
        return err
    }
    tx := config.DB.Begin()
    if tx.Error != nil {
        return tx.Error
    }

    for _, user := range unActiveUsers {
        if err := tx.Delete(&user).Error; err != nil {
            tx.Rollback()
            return err
        }
    }

    if err := tx.Commit().Error; err != nil {
        return err
    }
    log.Println("Deleted inactive users:", len(unActiveUsers))
    return nil
}