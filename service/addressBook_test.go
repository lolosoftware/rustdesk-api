package service

import (
	"testing"

	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestEnsurePeerInCollection(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&model.AddressBookCollection{}, &model.AddressBook{}); err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	DB = db

	collection := &model.AddressBookCollection{UserId: 42, Name: "Support"}
	if err = db.Create(collection).Error; err != nil {
		t.Fatal(err)
	}

	peer := &model.Peer{
		Id:       "STA007525",
		UserId:   7,
		Username: "operator",
		Hostname: "workstation",
		Os:       "Windows 11",
	}
	addressBooks := &AddressBookService{}
	if err = addressBooks.EnsurePeerInCollection(peer, collection.Id); err != nil {
		t.Fatal(err)
	}
	if err = addressBooks.EnsurePeerInCollection(peer, collection.Id); err != nil {
		t.Fatal(err)
	}

	var entries []model.AddressBook
	if err = db.Find(&entries).Error; err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("got %d address book entries, want 1", len(entries))
	}
	entry := entries[0]
	if entry.UserId != collection.UserId || entry.CollectionId != collection.Id {
		t.Fatalf("entry owner/collection = %d/%d, want %d/%d", entry.UserId, entry.CollectionId, collection.UserId, collection.Id)
	}
	if entry.Id != peer.Id || entry.Username != peer.Username || entry.Hostname != peer.Hostname || entry.Platform != "Windows" {
		t.Fatalf("unexpected address book entry: %#v", entry)
	}
}

func TestEnsurePeerInCollectionRejectsUnknownCollection(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&model.AddressBookCollection{}, &model.AddressBook{}); err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	DB = db

	addressBooks := &AddressBookService{}
	if err = addressBooks.EnsurePeerInCollection(&model.Peer{Id: "123"}, 999); err == nil {
		t.Fatal("expected an error for an unknown address book collection")
	}
}

func TestFilterDuplicateAddressBookHostnames(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:duplicate-address-books?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&model.AddressBook{}); err != nil {
		t.Fatal(err)
	}

	previousDB := DB
	DB = db
	t.Cleanup(func() { DB = previousDB })

	entries := []*model.AddressBook{
		{Id: "1", UserId: 1, CollectionId: 10, Hostname: "PC-COMPTA", Platform: "Windows"},
		{Id: "2", UserId: 1, CollectionId: 10, Hostname: " pc-compta ", Platform: "Windows"},
		{Id: "3", UserId: 1, CollectionId: 11, Hostname: "PC-COMPTA", Platform: "Windows"},
		{Id: "4", UserId: 2, CollectionId: 10, Hostname: "PC-COMPTA", Platform: "Windows"},
		{Id: "5", UserId: 1, CollectionId: 10, Hostname: "SRV-FICHIERS", Platform: "Linux"},
		{Id: "6", UserId: 1, CollectionId: 10, Hostname: "SRV-FICHIERS", Platform: "Linux"},
		{Id: "7", UserId: 1, CollectionId: 10, Hostname: "POSTE-UNIQUE", Platform: "Windows"},
		{Id: "8", UserId: 1, CollectionId: 10, Hostname: ""},
		{Id: "9", UserId: 1, CollectionId: 10, Hostname: "   "},
	}
	if err = db.Create(&entries).Error; err != nil {
		t.Fatal(err)
	}

	addressBooks := &AddressBookService{}
	duplicates := addressBooks.List(1, 100, func(tx *gorm.DB) {
		addressBooks.FilterDuplicateHostnames(tx)
	})
	if duplicates.Total != 4 {
		t.Fatalf("got %d duplicate entries, want 4", duplicates.Total)
	}

	windowsDuplicates := addressBooks.List(1, 100, func(tx *gorm.DB) {
		tx.Where("platform like ?", "%Windows%")
		addressBooks.FilterDuplicateHostnames(tx)
	})
	if windowsDuplicates.Total != 2 {
		t.Fatalf("got %d Windows duplicate entries, want 2", windowsDuplicates.Total)
	}
}
