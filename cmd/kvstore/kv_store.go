package main

import (
	"errors"
	"fmt"
	"github.com/exsql-io/kv-store/pkg/lib/kvstore"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	ErrInvalidNumberOfArguments = errors.New("invalid number of arguments")
)

func main() {
	app := &cli.App{
		Name:    "kvs",
		Usage:   "in-memory key value store",
		Version: kvstore.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "set",
				Usage: "set the value of a key",
				Args:  true,
				Action: func(cCtx *cli.Context) error {
					kvs, err := kvstore.New(cCtx.String("path"))
					if err != nil {
						return err
					}

					defer func() {
						if err := kvs.Close(); err != nil {
							log.Printf("kvstore.Close: %v", err)
						}
					}()

					if cCtx.Args().Len() != 2 {
						return ErrInvalidNumberOfArguments
					}

					return kvs.Set(cCtx.Args().Get(0), cCtx.Args().Get(1))
				},
			},
			{
				Name:  "get",
				Usage: "get the value of a key",
				Action: func(cCtx *cli.Context) error {
					kvs, err := kvstore.New(cCtx.String("path"))
					if err != nil {
						return err
					}

					defer func() {
						if err := kvs.Close(); err != nil {
							log.Printf("kvstore.Close: %v", err)
						}
					}()

					if cCtx.Args().Len() != 1 {
						return ErrInvalidNumberOfArguments
					}

					value, exists, err := kvs.Get(cCtx.Args().First())
					if err != nil {
						return err
					}

					if exists {
						fmt.Println(fmt.Sprintf("value: %s", value))
					} else {
						fmt.Println("nil")
					}

					return nil
				},
			},
			{
				Name:  "rm",
				Usage: "remove the value of a key",
				Action: func(cCtx *cli.Context) error {
					kvs, err := kvstore.New(cCtx.String("path"))
					if err != nil {
						return err
					}

					defer func() {
						if err := kvs.Close(); err != nil {
							log.Printf("kvstore.Close: %v", err)
						}
					}()

					if cCtx.Args().Len() != 1 {
						return ErrInvalidNumberOfArguments
					}

					return kvs.Remove(cCtx.Args().First())
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
