package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
	"log"
)

const userPrefix = "user_"
const tabPrefix = "tab_"


func prefix(prefix string, key string) []byte {
	return []byte(fmt.Sprintf("%s%s", prefix, key))
}

type Store struct {
	db *badger.DB
}

func NewStore(location string) (*Store, error) {
	db, err := badger.Open(badger.DefaultOptions(location))
	if err != nil {
		return nil, err
	}
	return &Store{
		db,
	}, nil
}

func (s Store) Close() {
	err := s.db.Close()
	if err != nil {
		log.Fatalf("%v", err)
	}
}

type User struct {
	Name     string
	Password []byte
	Admin    bool
	Tabs  []uuid.UUID
}

type Tab struct {
	Id    uuid.UUID
	Owner string
	Public bool // Visible on home page?
	Contents string // JSON encoded tab contents
}

func (s *Store) CreateUser(user User) error {
	return s.db.Update(func(txn *badger.Txn) error {
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(&user)
		if err != nil {
			return err
		}

		return txn.Set(prefix(userPrefix, user.Name), b.Bytes())
	})
}

func (s *Store) GetUser(name string) (User, error) {
	var res User
	return res, s.db.View(func(txn *badger.Txn) error {
		entry, err := txn.Get(prefix(userPrefix, name))
		if err != nil {
			return err
		}
		return entry.Value(func(val []byte) error {
			return json.NewDecoder(bytes.NewBuffer(val)).Decode(&res)
		})
	})
}

func (s *Store) CountUsers() (int, error) {
	res := 0
	return res, s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		for it.Seek([]byte(userPrefix)); it.ValidForPrefix([]byte(userPrefix)); it.Next() {
			res += 1
		}

		it.Close()

		return nil
	})
}

func (s *Store) CountAdminUsers() (int, error) {
	users, err := s.GetUsers()
	if err != nil {
		return 0, err
	}

	num := 0
	for _, user := range users {
		if user.Admin {
			num += 1
		}
	}

	return num, nil
}

func (s Store) UpdateUser(user *User) error {
	return s.db.Update(func(txn *badger.Txn) error {
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(&user)
		if err != nil {
			return err
		}

		return txn.Set(prefix(userPrefix, user.Name), b.Bytes())
	})
}

func (s Store) GetTab(id uuid.UUID) (*Tab, error) {
	var res *Tab
	return res, s.db.View(func(txn *badger.Txn) error {
		entry, err := txn.Get(prefix(tabPrefix, id.String()))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil {
			return err
		}
		return entry.Value(func(val []byte) error {
			return json.NewDecoder(bytes.NewBuffer(val)).Decode(&res)
		})
	})
}

func (s Store) CreateTab(tab Tab) error {
	err := s.AddTabToUser(tab.Owner, tab.Id)
	if err != nil {
		return err
	}

	return s.db.Update(func(txn *badger.Txn) error {
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(&tab)
		if err != nil {
			return err
		}

		return txn.Set(prefix(tabPrefix, tab.Id.String()), b.Bytes())
	})
}

func (s Store) AddTabToUser(owner string, id uuid.UUID) error {
	return s.db.Update(func(txn *badger.Txn) error {
		entry, err := txn.Get(prefix(userPrefix, owner))
		if err != nil {
			return err
		}

		var user User
		err = entry.Value(func(val []byte) error {
			return json.NewDecoder(bytes.NewBuffer(val)).Decode(&user)
		})
		if err != nil {
			return err
		}

		user.Tabs = append(user.Tabs, id)

		var b bytes.Buffer
		err = json.NewEncoder(&b).Encode(&user)
		if err != nil {
			return err
		}

		return txn.Set(prefix(userPrefix, user.Name), b.Bytes())
	})
}

func (s Store) RmTabFromUser(owner string, id uuid.UUID) error {
	return s.db.Update(func(txn *badger.Txn) error {
		entry, err := txn.Get(prefix(userPrefix, owner))
		if err != nil {
			return err
		}

		var user User
		err = entry.Value(func(val []byte) error {
			return json.NewDecoder(bytes.NewBuffer(val)).Decode(&user)
		})
		if err != nil {
			return err
		}

		var toRemove []int
		for i, a := range user.Tabs {
			if a == id {
				toRemove = append(toRemove, i)
			}
		}

		for _, r := range toRemove {
			user.Tabs = append(user.Tabs[:r], user.Tabs[r+1:]...)
		}

		var b bytes.Buffer
		err = json.NewEncoder(&b).Encode(&user)
		if err != nil {
			return err
		}

		return txn.Set(prefix(userPrefix, user.Name), b.Bytes())
	})
}

func (s Store) GetUserTabs(user *User) ([]Tab, error) {
	if user.Tabs == nil {
		return []Tab{}, nil
	}

	res := make([]Tab, len(user.Tabs))

	return res, s.db.View(func(txn *badger.Txn) error {
		for i, id := range user.Tabs {
			entry, err := txn.Get(prefix(tabPrefix, id.String()))
			if err == badger.ErrKeyNotFound {
				return nil
			}
			if err != nil {
				return err
			}
			err = entry.Value(func(val []byte) error {
				return json.NewDecoder(bytes.NewBuffer(val)).Decode(&res[i])
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s Store) RmTab(tab *Tab) error {
	if err := s.RmTabFromUser(tab.Owner, tab.Id); err != nil {
		return err
	}

	return s.db.Update(func(txn *badger.Txn) error {

		return txn.Delete(prefix(tabPrefix, tab.Id.String()))
	})
}

func (s Store) GetUsers() ([]User, error) {
	var res []User
	return res, s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		for it.Seek([]byte(userPrefix)); it.ValidForPrefix([]byte(userPrefix)); it.Next() {

			var user User
			err := it.Item().Value(func(val []byte) error {
				return json.NewDecoder(bytes.NewBuffer(val)).Decode(&user)
			})

			if err != nil {
				return err
			}

			res = append(res, user)
		}

		it.Close()

		return nil
	})
}

func (s Store) GetTabs() ([]Tab, error) {
	var res []Tab
	return res, s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		for it.Seek([]byte(tabPrefix)); it.ValidForPrefix([]byte(tabPrefix)); it.Next() {

			var tab Tab
			err := it.Item().Value(func(val []byte) error {
				return json.NewDecoder(bytes.NewBuffer(val)).Decode(&tab)
			})

			if err != nil {
				return err
			}

			res = append(res, tab)
		}

		it.Close()

		return nil
	})
}

func (s Store) GetPublicTabs() ([]Tab, error) {
	var res []Tab
	return res, s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		for it.Seek([]byte(tabPrefix)); it.ValidForPrefix([]byte(tabPrefix)); it.Next() {

			var tab Tab
			err := it.Item().Value(func(val []byte) error {
				return json.NewDecoder(bytes.NewBuffer(val)).Decode(&tab)
			})

			if !tab.Public {
				continue
			}

			if err != nil {
				return err
			}

			res = append(res, tab)
		}

		it.Close()

		return nil
	})
}

func (s Store) RmUser(name string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		var user User
		entry, err := txn.Get(prefix(userPrefix, name))
		if err != nil {
			return err
		}

		err = entry.Value(func(val []byte) error {
			return json.NewDecoder(bytes.NewBuffer(val)).Decode(&user)
		})
		if err != nil {
			return err
		}

		for _, id := range user.Tabs {
			if err != nil {
				return err
			}
			err = txn.Delete(prefix(tabPrefix, id.String()))
			if err != nil {
				return err
			}
		}

		return txn.Delete(prefix(userPrefix, name))
	})
}

func (s Store) SetAdmin(name string, value bool) error {
	return s.db.Update(func(txn *badger.Txn) error {
		var user User
		entry, err := txn.Get(prefix(userPrefix, name))
		if err != nil {
			return err
		}
		err = entry.Value(func(val []byte) error {
			return json.NewDecoder(bytes.NewBuffer(val)).Decode(&user)
		})
		if err != nil {
			return err
		}

		user.Admin = value

		var b bytes.Buffer
		err = json.NewEncoder(&b).Encode(&user)
		if err != nil {
			return err
		}

		return txn.Set(prefix(userPrefix, user.Name), b.Bytes())
	})
}

func (s Store) SetTab(id uuid.UUID, tab *Tab) error {
	return s.db.Update(func(txn *badger.Txn) error {
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(&tab)
		if err != nil {
			return err
		}

		return txn.Set(prefix(tabPrefix, id.String()), b.Bytes())
	})
}


