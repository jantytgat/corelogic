package controllers

import (
	"fmt"
	"github.com/jantytgat/corelogic/internal/models"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type FrameworkController struct {
	Frameworks map[string]models.Framework

	Release models.Release
	//SortedOlderVersions []string

	Expressions     map[string]string
	Fields          map[string]string
	SortedFieldKeys []string

	SectionData map[string][]string
}

func (c *FrameworkController) Parse() error {
	fmt.Printf("Parsing %s\n", c.Release.GetSemanticVersion())
	framework := c.Frameworks[c.Release.GetSemanticVersion()]
	installExpressions, err := framework.GetExpressions("install", nil)
	if err != nil {
		return err
	}
	uninstallExpressions, err := framework.GetExpressions("uninstall", nil)
	if err != nil {
		return err
	}
	fields, err := c.collectFieldsFromFramework(c.Release.GetSemanticVersion())
	if err != nil {
		return err
	}

	for k, v := range installExpressions {
		fmt.Println(k, "\t", v)
	}
	for k, v := range uninstallExpressions {
		fmt.Println(k, "\t", v)
	}
	for k, v := range fields {
		fmt.Println(k, "\t", v)
	}

	return nil
}

func (c *FrameworkController) GetSortedOlderVersions() []string {
	keys := make([]string, 0, len(c.Frameworks))
	for k := range c.Frameworks {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	return keys
}

func (c *FrameworkController) GetOutput(version string, kind string, tagFilter []string) ([]string, error) {
	//defer general.FinishTimer(general.StartTimer("Frameworks " + f.Release.GetVersionAsString() + " get " + kind + " output"))

	var output []string

	var err error
	log.Printf("Get output for %s\n", version)
	framework := c.Frameworks[version]
	log.Println(framework.Release)
	log.Println(len(framework.Packages))
	c.Expressions, err = framework.GetExpressions(kind, tagFilter)
	if err != nil {
		log.Fatal(err)
		return output, err
	}

	c.Fields, err = c.collectFieldsFromFramework(version)
	if err != nil {
		log.Fatal(err)
		return output, err
	}

	c.setSortedFieldKeys(c.Fields)
	c.unfoldExpressions(version)
	framework.SortPrefixes(c.Frameworks[version].Prefixes)
	c.SectionData = make(map[string][]string)
	c.collectExpressionsPerSection(version)

	for _, p := range c.Frameworks[version].Prefixes {
		output = append(output, "### "+p.Section)
		output = append(output, c.SectionData[p.Section]...)
		output = append(output, "##########################")
	}

	return output, err
}

func (c *FrameworkController) collectFieldsFromFramework(version string) (map[string]string, error) {
	//defer general.FinishTimer(general.StartTimer("Frameworks " + f.Release.GetVersionAsString() + " get fields"))

	framework := c.Frameworks[version]
	fields, err := framework.GetFields()
	if err != nil {
		log.Fatal(err)
	}

	return c.unfoldFields(version, fields), err
}

func (c *FrameworkController) unfoldFields(version string, fields map[string]string) map[string]string {
	//defer general.FinishTimer(general.StartTimer("Frameworks " + f.Release.GetVersionAsString() + " unfold fields"))

	//framework := c.Frameworks[version]

	re := regexp.MustCompile(`<<[a-zA-Z0-9_.]*/[a-zA-Z0-9_]*>>`)
	for key := range fields {
		loop := true
		for loop {
			foundKeys := re.FindAllString(fields[key], -1)
			for _, foundKey := range foundKeys {
				searchKey := strings.ReplaceAll(foundKey, "<<", "")
				searchKey = strings.ReplaceAll(searchKey, ">>", "")
				fields[key] = strings.ReplaceAll(fields[key], foundKey, fields[searchKey])
			}

			if !re.MatchString(fields[key]) {
				loop = false
			}
		}

		//for k := range framework.GetPrefixMap() {
		//	if !strings.Contains(fields[key], "<<") {
		//		break
		//	}
		//
		//	fields[key] = strings.ReplaceAll(fields[key], "<<"+k+">>", framework.GetPrefixWithVersion(k))
		//}
	}

	return fields
}

func (c *FrameworkController) setSortedFieldKeys(fields map[string]string) {
	//defer general.FinishTimer(general.StartTimer("Frameworks " + f.Release.GetVersionAsString() + " set sorted field keys"))

	fieldKeys := make([]string, 0, len(fields))
	for f := range fields {
		fieldKeys = append(fieldKeys, f)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(fieldKeys)))

	c.SortedFieldKeys = fieldKeys
}

func (c *FrameworkController) unfoldExpressions(version string) {
	//defer general.FinishTimer(general.StartTimer("Frameworks " + f.Release.GetVersionAsString() + " unfold expressions"))

	wg := &sync.WaitGroup{}
	ch := make(chan models.UnfoldedExpressionData)

	count := 0
	for k := range c.Expressions {
		wg.Add(1)
		count++
		go c.unfoldExpressionHandler(version, k, ch, wg)
	}
	wg.Add(1)
	go c.unfoldedExpressionCollector(count, ch, wg)

	wg.Wait()
	close(ch)
}

func (c *FrameworkController) unfoldExpressionHandler(version string, elementName string, ch chan<- models.UnfoldedExpressionData, wg *sync.WaitGroup) {
	defer wg.Done()

	output := models.UnfoldedExpressionData{
		Key:   elementName,
		Value: c.Expressions[elementName],
	}

	output.Value = c.replaceDataInExpression(version, output.Value)
	ch <- output
}

func (c *FrameworkController) replaceDataInExpression(version string, expression string) string {
	if expression != "" {
		expression = c.replaceFieldsInExpression(expression)
		expression = c.replacePrefixesInExpression(version, expression)
		expression = strings.TrimSuffix(expression, "\n")
	}

	return expression
}

func (c *FrameworkController) replaceFieldsInExpression(expression string) string {
	re := regexp.MustCompile(`<<[a-zA-Z0-9_.]*/[a-zA-Z0-9_]*>>`)

	loop := true
	for loop {
		foundKeys := re.FindAllString(expression, -1)
		for _, foundKey := range foundKeys {
			searchKey := strings.ReplaceAll(foundKey, "<<", "")
			searchKey = strings.ReplaceAll(searchKey, ">>", "")
			expression = strings.ReplaceAll(expression, foundKey, c.Fields[searchKey])
		}

		if !re.MatchString(expression) {
			loop = false
		}
	}

	return expression
}

func (c *FrameworkController) replacePrefixesInExpression(version string, expression string) string {
	framework := c.Frameworks[version]

	// Replace prefixes in expressions
	for p := range framework.GetPrefixMap() {
		if !strings.Contains(expression, "<<") {
			break
		}
		expression = strings.ReplaceAll(expression, "<<"+p+">>", framework.GetPrefixWithVersion(p))
	}

	return expression
}

func (c *FrameworkController) unfoldedExpressionCollector(count int, ch <-chan models.UnfoldedExpressionData, wg *sync.WaitGroup) {
	defer wg.Done()
	completed := false

	var expressions = make(map[string]string)
	for !completed {
		select {
		case data, ok := <-ch:
			if !ok {
				completed = true
			}
			expressions[data.Key] = data.Value
			count--
		default:
			if count == 0 {
				completed = true
			}
		}
	}
	c.Expressions = expressions
}

func (c *FrameworkController) collectExpressionsPerSection(version string) {
	//defer general.FinishTimer(general.StartTimer("Frameworks " + f.Release.GetVersionAsString() + " collect expressions per section"))

	wg := &sync.WaitGroup{}
	globalChannel := make(chan models.SectionData)

	count := 0
	for _, p := range c.Frameworks[version].Prefixes {
		sectionChannel := make(chan string)

		wg.Add(1)
		go c.collectExpressionsForSection(p.Section, sectionChannel, wg)

		wg.Add(1)
		go c.sectionExpressionCollector(p.Section, globalChannel, sectionChannel, wg)
		count++
	}

	wg.Add(1)
	go c.ExpressionCollector(count, globalChannel, wg)

	wg.Wait()
}

func (c *FrameworkController) collectExpressionsForSection(sectionName string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for k := range c.Expressions {
		if strings.Contains(k, sectionName) {
			if c.Expressions[k] != "" {
				ch <- c.Expressions[k]
			}
		}
	}
	close(ch)
}

func (c *FrameworkController) sectionExpressionCollector(sectionName string, globalChannel chan<- models.SectionData, ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	completed := false

	var expressions []string
	for !completed {
		select {
		case data, ok := <-ch:
			if !ok {
				completed = true
			} else {
				expressions = append(expressions, data)
			}
		}
	}

	globalChannel <- models.SectionData{
		Name:        sectionName,
		Expressions: expressions,
	}
}

func (c *FrameworkController) ExpressionCollector(count int, ch <-chan models.SectionData, wg *sync.WaitGroup) {
	defer wg.Done()
	completed := false

	for !completed {
		select {
		case sectionData, ok := <-ch:
			if !ok {
				completed = true
			}
			c.SectionData[sectionData.Name] = sectionData.Expressions
			count--
		default:
			if count == 0 {
				completed = true
			}
		}
	}
}
