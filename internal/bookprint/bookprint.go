/*
 * Copyright (C) 2023 Stefan Kühnel
 *
 * SPDX-License-Identifier: EUPL-1.2
 */

package bookprint

import (
	"html/template"
	"os"
	"path/filepath"

	"stefanco.de/bookprint/internal/book"
)

type Config struct {
	File        []byte
	OutputDir   string
	TemplateDir string
}

// ToDo: Show debug info when a template is not existent.
// ToDo: Check if it would be reliable to skip not existent templates.

func New(config *Config) error {
	b, err := book.New(config.File)
	if err != nil {
		return err
	}

	err = createIndex(b, config)
	if err != nil {
		return err
	}

	err = createMap(b, config)
	if err != nil {
		return err
	}

	err = createPages(b, config)
	if err != nil {
		return err
	}

	return nil
}

func createIndex(b *book.Book, config *Config) error {
	templateFile := filepath.Join(config.TemplateDir, "index.html")
	templateFileName := filepath.Base(templateFile)

	outputFile, err := os.Create(filepath.Join(config.OutputDir, templateFileName))
	if err != nil {
		return err
	}

	t, err := template.New(templateFileName).ParseFiles(templateFile)
	if err != nil {
		return err
	}

	err = t.Execute(outputFile, b)
	if err != nil {
		return err
	}

	err = outputFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func createMap(b *book.Book, config *Config) error {
	templateFile := filepath.Join(config.TemplateDir, "map.html")
	templateFileName := filepath.Base(templateFile)

	outputFile, err := os.Create(filepath.Join(config.OutputDir, templateFileName))
	if err != nil {
		return err
	}

	t, err := template.New(templateFileName).ParseFiles(templateFile)
	if err != nil {
		return err
	}

	err = t.Execute(outputFile, b)
	if err != nil {
		return err
	}

	err = outputFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func createPages(b *book.Book, config *Config) error {
	templateFile := filepath.Join(config.TemplateDir, "page.html")
	templateFileName := filepath.Base(templateFile)

	t, err := template.New(templateFileName).ParseFiles(templateFile)
	if err != nil {
		return err
	}

	type Page struct {
		MetaData *book.MetaData
		Page     *book.Page
	}

	for _, page := range b.Pages {
		outputFile, err := os.Create(filepath.Join(config.OutputDir, page.Path))
		if err != nil {
			return err
		}

		err = t.Execute(outputFile, &Page{
			MetaData: b.MetaData,
			Page:     page,
		})
		if err != nil {
			return err
		}

		err = outputFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
