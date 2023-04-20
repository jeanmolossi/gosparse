package filter

// Predicate é um tipo para definir um enum de predicados
// aceitos nos campos do parâmetro "fields"
type Predicate int

const (
	NONE Predicate = iota
	EQ
	NEQ
	IN
	NIN
	GT
	GTE
	LT
	LTE
	BLANK
	NULL
	NOT_NULL
	START
	END
)

func getPredicate(pre string) Predicate {
	if pre == "" {
		return NONE
	}

	preMap := map[string]Predicate{
		"eq":      EQ,
		"neq":     NEQ,
		"in":      IN,
		"nin":     NIN,
		"gt":      GT,
		"gte":     GTE,
		"lt":      LT,
		"lte":     LTE,
		"blank":   BLANK,
		"null":    NULL,
		"notnull": NOT_NULL,
		"start":   START,
		"end":     END,
	}

	if predicate, found := preMap[pre]; found {
		return predicate
	}

	return NONE
}
