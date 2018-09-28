package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/amclees/hit/store"
)

var defaultStoreFactory store.StoreFactory = store.FileStoreFactory

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.New("No arguments received")
	}

	storeFactory := defaultStoreFactory
	if customFactory, err := getCustomStoreFactory(); err == nil  {
		storeFactory = customFactory
	} else if err.Error() != "default store used" {
		return err
	}

	st, err := storeFactory()
	if err != nil {
		return err
	}

	switch (os.Args[1]) {
		case "is":
			err := is(st)
			return err
		case "store":
			err := setStore()
			return err
		default:
			err := lookup(st)
			return err
	}
}

func getCustomStoreFactory() (store.StoreFactory, error) {
	return nil, errors.New("default store used")
}

func is(st *store.Store) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	el := store.Element{Value: wd, Command: "cd"}
	var key string
	l := len(os.Args)
	if l >= 3 {
		key = os.Args[2]
	}
	if l >= 4 {
		el.Command = os.Args[3]
	}
	if l >= 5 {
		el.Value = os.Args[4]
	}

	return (*st).Insert(key, el)
}

func setStore() error {
	return nil
}

func lookup(st *store.Store) error {
	key := os.Args[1]
	el, err := (*st).Lookup(key)
	if err != nil {
		return err
	}

	fmt.Print(el.Command, " ", el.Value)

	return nil
}
