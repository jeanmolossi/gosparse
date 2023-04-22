package gosparse

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jeanmolossi/gosparse/internal/filter"
	"github.com/jeanmolossi/gosparse/internal/include"
	"github.com/jeanmolossi/gosparse/internal/pagination"
	"github.com/jeanmolossi/gosparse/internal/sort"
	"github.com/jeanmolossi/gosparse/internal/sparsefieldsets"
)

// getValueAndValidate recebe a interface e trata para que
// o reflect.Value seja correspondente à uma estrutura.
//
// getValueAndValidate pode receber uma referência para uma estrutura
// ou o valor de uma estrutura.
//
// getValueAndValidate também valida se de fato recebeu uma estrutura.
// caso seja algo diferente de uma estrutura será retornado um erro.
func getValueAndValidate(s any) (reflect.Value, error) {
	value := reflect.ValueOf(s)

	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return value, fmt.Errorf("can extract from structs only")
	}

	return value, nil
}

// extractTag recebe um StructField e extrai as configurações de
// querystring aceitas a partir da Tagname (gosparse).
func extractTag(typ reflect.StructField) *config {
	tag := typ.Tag.Get(Tagname)
	if tag == "" || tag == "-" {
		return nil
	}

	result := extractor(tag)
	return &result
}

// handleTags recebe um reflect.Value de uma estrutura e trata para
// ter um objeto de configuração válido para montar um GoSparse
func handleTags(v reflect.Value) (map[string]config, error) {
	fields := map[string]config{}

	for i := 0; i < v.NumField(); i++ {
		typ := v.Type().Field(i)

		conf := *extractTag(typ)
		fields[conf.Name] = conf

		kind := typ.Type.Kind()
		if conf.Relation && (kind == reflect.Ptr || kind == reflect.Struct) {
			field, err := getValueAndValidate(v.Field(i).Interface())
			if err != nil {
				return nil, err
			}

			res, err := handleTags(field)
			if err != nil {
				return nil, err
			}

			for _, rel := range res {
				k := []string{conf.Name, rel.Name}
				fields[strings.Join(k, ".")] = rel
			}
		}
	}

	return fields, nil
}

// Tag extractor -----------------------------------

var Tagname = "gosparse"

type config struct {
	// Name corresponde ao nome do campo que será aceito na querystring
	Name string
	// Select indica se o campo pode ser utilizado no parâmetro "fields"
	//
	// 	fields[name]
	Select bool
	// Sort indica se é um campo aceito para o parâmetro "sort"
	//
	//	sort=name
	Sort bool
	// Filter indica se é um campo válido para o parâmetro "filter"
	//
	//	filter[name]
	Filter bool
	// Relation indica se é um campo válido para o parâmetro "include"
	//
	//	include=name
	Relation bool
}

// extractor recebe a tag do campo e trata para que seja retornado
// um objeto de configuração válido.
func extractor(tag string) config {
	configs := strings.Split(tag, ";")

	c := config{}

	for _, conf := range configs {
		name, done := strings.CutPrefix(conf, "name:")
		if done {
			c.Name = name
			continue
		}

		if strings.HasPrefix(conf, "select") {
			c.Select = true
			continue
		}

		if strings.HasPrefix(conf, "sort") {
			c.Sort = true
			continue
		}

		if strings.HasPrefix(conf, "filter") {
			c.Filter = true
			continue
		}

		if strings.HasPrefix(conf, "relation") {
			c.Relation = true
			continue
		}
	}

	return c
}

// Real extract

// Extract recebe interface e trata para que seja montado um Gosparse
// baseado na tag "gosparse" da estrutura
func Extract(s any) (Gosparse, error) {
	value, err := getValueAndValidate(s)
	if err != nil {
		return Gosparse{}, nil
	}

	extracted, err := handleTags(value)
	if err != nil {
		return Gosparse{}, nil
	}

	gs := Gosparse{
		Include:    *include.New(),
		Fieldset:   *sparsefieldsets.New(),
		Filter:     *filter.New(),
		Pagination: *pagination.New(),
		Sort:       *sort.New(),
	}

	relations := make([]string, 0, len(extracted))
	fields := make([]string, 0, len(extracted))
	filters := make([]string, 0, len(extracted))
	sorter := make([]string, 0, len(extracted))

	for field, conf := range extracted {
		if conf.Relation {
			relations = append(relations, field)
		}

		if conf.Select {
			fields = append(fields, field)
		}

		if conf.Filter {
			filters = append(filters, field)
		}

		if conf.Sort {
			sorter = append(sorter, field)
		}
	}

	AcceptRelations(relations...)(&gs)
	AcceptFields(fields...)(&gs)
	AcceptFilters(filters...)(&gs)
	AcceptSortBy(sorter...)(&gs)

	return gs, nil
}
