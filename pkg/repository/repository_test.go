package repository

import (
	"fmt"
	"github.com/pxecore/pxecore/pkg/entity"
	"go.uber.org/atomic"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestConcurrency(t *testing.T) {
	repositories := [...]Repository{newMemoryRepositoryTest(t)}

	for _, repository := range repositories {
		t.Run(fmt.Sprint("Test_", reflect.TypeOf(repository).Elem().Name()), func(t *testing.T) {
			runOpenConcurrencyTest(t, repository)
			runReadWriteConcurrencyTest(t, repository)
		})
	}
}

func TestCRUD(t *testing.T) {
	repositories := [...]Repository{newMemoryRepositoryTest(t)}

	for _, repository := range repositories {
		t.Run(fmt.Sprint("Test_", reflect.TypeOf(repository).Elem().Name()), func(t *testing.T) {
			runIndividualHostCRUD(t, repository)
			runIndividualGroupCRUD(t, repository)
			runIndividualTemplateCRUD(t, repository)
		})
	}
}

func newMemoryRepositoryTest(t *testing.T) Repository {
	m, err := newMemoryRepository(make(map[string]interface{}))
	if err != nil {
		t.Fatal("Error newMemoryRepositoryTest - ", err)
	}
	return m
}

func runIndividualHostCRUD(t *testing.T, m Repository) {
	if err := m.Write(func(s Session) error {
		return s.Host().Create(entity.Host{
			ID:            "10",
			HardwareAddr:  []string{"86-53-25-6A-E0-D4"},
			TrapMode:      true,
			TrapTriggered: true,
			Vars:          map[string]string{"foo": "bar"},
			GroupID:       "",
			TemplateID:    "",
		})
	}); err != nil {
		t.Fatal("runIndividualHostCRUD - error creating ", err)
	}
	if err := m.Read(func(s Session) error {
		h, err := s.Host().Get("10")
		if err != nil {
			return err
		}
		if h.ID != "10" || h.HardwareAddr[0] != "86-53-25-6A-E0-D4" || h.Vars["foo"] != "bar" {
			t.Fatal("Invalid stored data - ", h)
		}
		h, err = s.Host().FindByHardwareAddr("86-53-25-6A-E0-D4")
		if err != nil {
			return err
		}
		if h.ID != "10" || h.HardwareAddr[0] != "86-53-25-6A-E0-D4" || h.Vars["foo"] != "bar" {
			t.Fatal("Invalid stored data - ", h)
		}
		return nil
	}); err != nil {
		t.Fatal("runIndividualHostCRUD - error creating ", err)
	}

	if err := m.Write(func(s Session) error {
		return s.Host().Update(entity.Host{
			ID:            "10",
			HardwareAddr:  []string{"86-53-25-6A-E0-D5"},
			TrapMode:      false,
			TrapTriggered: false,
			Vars:          map[string]string{"bar": "foo"},
			GroupID:       "",
			TemplateID:    "",
		})
	}); err != nil {
		t.Fatal("runIndividualHostCRUD - error creating ", err)
	}
	if err := m.Read(func(s Session) error {
		_, err := s.Host().Get("10")
		if err != nil {
			return err
		}
		_, err = s.Host().FindByHardwareAddr("86-53-25-6A-E0-D5")
		if err != nil {
			return err
		}
		_, err = s.Host().FindByHardwareAddr("86-53-25-6A-E0-D4")
		if err == nil {
			t.Fatal("runIndividualHostCRUD - deleted HardwareAddr returned ", err)
		}
		return nil
	}); err != nil {
		t.Fatal("runIndividualHostCRUD - error creating ", err)
	}
}

func runIndividualGroupCRUD(t *testing.T, m Repository) {
	if err := m.Write(func(s Session) error {
		return s.Group().Create(entity.Group{
			ID:                "11",
			Vars:              map[string]string{"foo": "bar"},
			HostsIDs:          []string{"host"},
			GroupIDs:          []string{"GroupID"},
			ParentID:          "",
			TemplateID: "defaultTemplate",
		})
	}); err != nil {
		t.Fatal("runIndividualHostCRUD - error creating ", err)
	}
	if err := m.Read(func(s Session) error {
		h, err := s.Group().Get("11")
		if err != nil {
			return err
		}
		if h.ID != "11" || h.HostsIDs[0] != "host" || h.Vars["foo"] != "bar" {
			t.Fatal("Invalid stored data - ", h)
		}
		if err != nil {
			return err
		}
		if h.ID != "11" || h.HostsIDs[0] != "host" || h.Vars["foo"] != "bar" {
			t.Fatal("Invalid stored data - ", h)
		}
		return nil
	}); err != nil {
		t.Fatal("runIndividualHostCRUD - error creating ", err)
	}

	if err := m.Write(func(s Session) error {
		return s.Group().Update(entity.Group{
			ID:                "11",
			Vars:              map[string]string{"foo2": "bar2"},
			HostsIDs:          []string{"host2"},
			GroupIDs:          []string{"GroupID2"},
			ParentID:          "",
			TemplateID: "defaultTemplate",
		})
	}); err != nil {
		t.Error("runIndividualHostCRUD - error creating ", err)
	}
	if err := m.Read(func(s Session) error {
		h, err := s.Group().Get("11")
		if err != nil {
			return err
		}
		if h.GroupIDs[0] != "GroupID2" {
			t.Error("Invalid stored data - ", h)
		}
		return nil
	}); err != nil {
		t.Error("runIndividualHostCRUD - error creating ", err)
	}
}
func runIndividualTemplateCRUD(t *testing.T, m Repository) {
	if err := m.Write(func(s Session) error {
		return s.Template().Create(entity.Template{
			ID:       "id",
			Template: "template",
		})
	}); err != nil {
		t.Fatal("runIndividualHostCRUD - error creating ", err)
	}
	if err := m.Read(func(s Session) error {
		h, err := s.Template().Get("id")
		if err != nil {
			return err
		}
		if h.Template != "template" {
			t.Fatal("Invalid stored data - ", h)
		}
		return nil
	}); err != nil {
		t.Fatal("runIndividualHostCRUD - error creating ", err)
	}

	if err := m.Write(func(s Session) error {
		return s.Template().Update(entity.Template{
			ID:       "id",
			Template: "template2",
		})
	}); err != nil {
		t.Fatal("runIndividualHostCRUD - error creating ", err)
	}
	if err := m.Read(func(s Session) error {
		h, err := s.Template().Get("id")
		if err != nil {
			return err
		}
		if h.Template != "template2" {
			t.Fatal("Invalid stored data - ", h)
		}
		return nil
	}); err != nil {
		t.Fatal("runIndividualHostCRUD - error creating ", err)
	}
}

func runOpenConcurrencyTest(t *testing.T, m Repository) {
	s, _ := m.Open(true)
	if err := s.Close(); err != nil {
		t.Fatal("OK case failed - error: ", err)
	}

	if err := s.Close(); err == nil {
		t.Fatal("KO case failed - session should be close - error: ")
	}

	s, _ = m.Open(false)
	if err := s.Close(); err != nil {
		t.Fatal("OK case failed - error: ", err)
	}

	if err := s.Close(); err == nil {
		t.Fatal("KO case failed - session should be close - error: ")
	}
}

func runReadWriteConcurrencyTest(t *testing.T, m Repository) {
	var wg sync.WaitGroup
	var a atomic.Bool
	for i := 1; i <= 5; i++ {
		id := i
		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			defer wg.Done()
			switch id {
			case 3, 4:
				_ = m.Write(func(s Session) error {
					if !a.CAS(false, true) {
						t.Fatal("write operation did't lock. Expected: ", id, " - Returned:", a)
					}
					time.Sleep(100)
					if !a.CAS(true, false) {
						t.Fatal("write operation did't lock. Expected: ", id, " - Returned:", a)
					}
					return nil
				})
				break
			default:
				_ = m.Read(func(s Session) error {
					a.CAS(false, true)
					time.Sleep(100)
					a.CAS(true, false)
					return nil
				})
			}
		}(id, &wg)
	}
	wg.Wait()
}
