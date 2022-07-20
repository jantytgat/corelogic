package controllers

import (
	"github.com/jantytgat/corelogic/internal/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

type YamlController struct{}

func (c *YamlController) Load(version string) (FrameworkController, error) {
	//defer general.FinishTimer(general.StartTimer("Loading framework " + version))
	framework := models.Framework{}
	rootDir := "assets/framework/" + version
	var source []byte
	var err error

	source, err = ioutil.ReadFile(rootDir + "/framework.yaml")
	if err != nil {
		log.Fatal(err)
		return FrameworkController{Framework: framework}, err
	}

	err = yaml.Unmarshal(source, &framework)
	if err != nil {
		log.Fatal(err)
		return FrameworkController{Framework: framework}, err
	}

	framework.Packages = []models.Package{}

	subDirs, err := ioutil.ReadDir(rootDir + "/packages")
	if err != nil {
		log.Fatal(err)
		return FrameworkController{Framework: framework}, err
	}

	for _, d := range subDirs {
		if d.IsDir() {
			var p models.Package
			p, err = c.getPackagesFromDirectory(rootDir, d.Name())
			if err != nil {
				return FrameworkController{Framework: framework}, err
			}
			framework.Packages = append(framework.Packages, p)
		}
	}

	return FrameworkController{Framework: framework}, err
}

func (c *YamlController) getPackagesFromDirectory(rootDir string, directoryName string) (models.Package, error) {
	// defer general.FinishTimer(general.StartTimer("GetPackagesFromDirectory " + rootDir + "/packages/" + directoryName))

	myPackage := models.Package{
		Name:    directoryName,
		Modules: []models.Module{},
	}

	files, err := ioutil.ReadDir(rootDir + "/packages/" + myPackage.Name)
	if err != nil {
		log.Fatal(err)
		return myPackage, err
	}

	for _, f := range files {
		if !f.IsDir() {
			if filepath.Ext(f.Name()) == ".yaml" {
				// log.Println(f.Name())
				var module models.Module
				module, err = c.getModuleFromFile(rootDir + "/packages/" + myPackage.Name + "/" + f.Name())
				if err != nil {
					return myPackage, err
				}
				myPackage.Modules = append(myPackage.Modules, module)
			}
		} else {
			var modules []models.Module
			modules, err = c.getModulesFromDirectory(rootDir + "/packages/" + myPackage.Name + "/" + f.Name())
			if err != nil {
				return myPackage, err
			}
			myPackage.Modules = append(myPackage.Modules, modules...)
		}
	}
	return myPackage, err
}

func (c *YamlController) getModuleFromFile(filePath string) (models.Module, error) {
	// defer general.FinishTimer(general.StartTimer("GetModuleFromFile " + filePath))

	module := models.Module{}

	moduleSource, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
		return module, err
	}

	err = yaml.Unmarshal(moduleSource, &module)
	if err != nil {
		log.Fatal(err)
	}

	return module, err
}

func (c *YamlController) getModulesFromDirectory(filePath string) ([]models.Module, error) {
	// defer general.FinishTimer(general.StartTimer("GetModulesFromDirectory " + filePath))

	var modules []models.Module

	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
		return modules, err
	}

	for _, f := range files {
		if !f.IsDir() {
			if filepath.Ext(f.Name()) == ".yaml" {
				// log.Println(f.Name())
				module, err := c.getModuleFromFile(filePath + "/" + f.Name())
				if err != nil {
					log.Fatal(err)
					return modules, err
				}
				modules = append(modules, module)
			}
		}
	}

	return modules, err
}
