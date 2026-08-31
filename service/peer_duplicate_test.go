package service

import (
	"testing"

	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestFilterDuplicateHostnames(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:duplicate-hostnames?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.AutoMigrate(&model.Peer{}); err != nil {
		t.Fatalf("migrate peers: %v", err)
	}

	previousDB := DB
	DB = db
	t.Cleanup(func() { DB = previousDB })

	peers := []*model.Peer{
		{Id: "1", Hostname: "PC-COMPTA"},
		{Id: "2", Hostname: " pc-compta "},
		{Id: "3", Hostname: "SRV-FICHIERS"},
		{Id: "4", Hostname: "SRV-FICHIERS"},
		{Id: "5", Hostname: "POSTE-UNIQUE"},
		{Id: "6", Hostname: ""},
		{Id: "7", Hostname: "   "},
	}
	if err := db.Create(&peers).Error; err != nil {
		t.Fatalf("create peers: %v", err)
	}

	result := (&PeerService{}).List(1, 100, func(tx *gorm.DB) {
		(&PeerService{}).FilterDuplicateHostnames(tx)
	})

	if result.Total != 4 {
		t.Fatalf("expected 4 duplicate peers, got %d", result.Total)
	}
	wantIDs := []string{"1", "2", "3", "4"}
	for index, wantID := range wantIDs {
		if result.Peers[index].Id != wantID {
			t.Errorf("peer %d: expected ID %s, got %s", index, wantID, result.Peers[index].Id)
		}
	}
}
