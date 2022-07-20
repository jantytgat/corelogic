package controllers

import (
	"github.com/jantytgat/corelogic/internal/models"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type FrameworkController struct {
	Framework models.Framework

	Expressions     map[string]string
	Fields          map[string]string
	SortedFieldKeys []string

	SectionData map[string][]string
}

func (c *FrameworkController) GetOutput(kind string, tagFilter []string) ([]string, error) {
	//defer general.FinishTimer(general.StartTimer("Framework " + f.Release.GetVersionAsString() + " get " + kind + " output"))

	var output []string

	var err error
	c.Expressions, err = c.Framework.GetExpressions(kind, tagFilter)
	if err != nil {
		log.Fatal(err)
		return output, err
	}

	c.Fields, err = c.collectFieldsFromFramework()
	if err != nil {
		log.Fatal(err)
		return output, err
	}

	c.setSortedFieldKeys(c.Fields)
	c.unfoldExpressions()
	c.Framework.SortPrefixes(c.Framework.Prefixes)
	c.SectionData = make(map[string][]string)
	c.collectExpressionsPerSection()

	for _, p := range c.Framework.Prefixes {
		output = append(output, "### "+p.Section)
		output = append(output, c.SectionData[p.Section]...)
		output = append(output, "##########################")
	}

	return output, err
}

func (c *FrameworkController) collectFieldsFromFramework() (map[string]string, error) {
	//defer general.FinishTimer(general.StartTimer("Framework " + f.Release.GetVersionAsString() + " get fields"))

	fields, err := c.Framework.GetFields()
	if err != nil {
		log.Fatal(err)
	}

	return c.unfoldFields(fields), err
}

func (c *FrameworkController) unfoldFields(fields map[string]string) map[string]string {
	//defer general.FinishTimer(general.StartTimer("Framework " + f.Release.GetVersionAsString() + " unfold fields"))

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

		for k := range c.Framework.GetPrefixMap() {
			if !strings.Contains(fields[key], "<<") {
				break
			}

			fields[key] = strings.ReplaceAll(fields[key], "<<"+k+">>", c.Framework.GetPrefixWithVersion(k))
		}
	}

	return fields
}

func (c *FrameworkController) setSortedFieldKeys(fields map[string]string) {
	//defer general.FinishTimer(general.StartTimer("Framework " + f.Release.GetVersionAsString() + " set sorted field keys"))

	fieldKeys := make([]string, 0, len(fields))
	for f := range fields {
		fieldKeys = append(fieldKeys, f)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(fieldKeys)))

	c.SortedFieldKeys = fieldKeys
}

func (c *FrameworkController) unfoldExpressions() {
	//defer general.FinishTimer(general.StartTimer("Framework " + f.Release.GetVersionAsString() + " unfold expressions"))

	wg := &sync.WaitGroup{}
	ch := make(chan models.UnfoldedExpressionData)

	count := 0
	for k := range c.Expressions {
		wg.Add(1)
		count++
		go c.unfoldExpressionHandler(k, ch, wg)
	}
	wg.Add(1)
	go c.unfoldedExpressionCollector(count, ch, wg)

	wg.Wait()
	close(ch)
}

func (c *FrameworkController) unfoldExpressionHandler(elementName string, ch chan<- models.UnfoldedExpressionData, wg *sync.WaitGroup) {
	defer wg.Done()

	output := models.UnfoldedExpressionData{
		Key:   elementName,
		Value: c.Expressions[elementName],
	}

	output.Value = c.replaceDataInExpression(output.Value)
	ch <- output
}

func (c *FrameworkController) replaceDataInExpression(expression string) string {
	if expression != "" {
		expression = c.replaceFieldsInExpression(expression)
		expression = c.replacePrefixesInExpression(expression)
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

func (c *FrameworkController) replacePrefixesInExpression(expression string) string {
	// Replace prefixes in expressions
	for p := range c.Framework.GetPrefixMap() {
		if !strings.Contains(expression, "<<") {
			break
		}
		expression = strings.ReplaceAll(expression, "<<"+p+">>", c.Framework.GetPrefixWithVersion(p))
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

func (c *FrameworkController) collectExpressionsPerSection() {
	//defer general.FinishTimer(general.StartTimer("Framework " + f.Release.GetVersionAsString() + " collect expressions per section"))

	wg := &sync.WaitGroup{}
	globalChannel := make(chan models.SectionData)

	count := 0
	for _, p := range c.Framework.Prefixes {
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
