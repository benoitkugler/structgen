package orm

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/benoitkugler/structgen/loader"
)

// Constraint is an arbitrary constraint
// Enums value, defined as #<TypeName>.<VarName> are replaced by their content
type Constraint struct {
	goTableName string
	constraint  string
}

func NewConstraint(c loader.Comment, lookupTable map[string]string) (Constraint, error) {
	constraint, err := parseComment(c.Content, lookupTable)
	if err != nil {
		return Constraint{}, err
	}
	return Constraint{goTableName: c.TypeName, constraint: constraint}, nil
}

var re = regexp.MustCompile(`#[\w_]+\.[\w_]+`)

func parseComment(comment string, lookupTable map[string]string) (string, error) {
	var notFound []string
	replacer := func(enumTag string) string {
		enumTag = enumTag[1:] // we remove #
		value, has := lookupTable[enumTag]
		if !has { // accumulate error
			notFound = append(notFound, enumTag)
		}
		return value
	}
	comment = re.ReplaceAllStringFunc(comment, replacer)
	if len(notFound) > 0 {
		return "", fmt.Errorf("unknown enum names : %s", strings.Join(notFound, " ; "))
	}
	return comment, nil
}

func (u Constraint) Render() string {
	if u.goTableName == "" {
		return u.constraint
	}
	return fmt.Sprintf(`ALTER TABLE %s %s;`, tableName(u.goTableName), u.constraint)
}

type ForeignKeyConstraint struct {
	sqlField       string
	sqlSourceTable string
	sqlTargetTable string
	deleteAction   string
}

func (s ForeignKeyConstraint) Render() string {
	onDelete := ""
	if s.deleteAction != "" {
		onDelete = "ON DELETE " + s.deleteAction
	}
	return fmt.Sprintf("ALTER TABLE %s ADD FOREIGN KEY(%s) REFERENCES %s %s;",
		s.sqlSourceTable, s.sqlField, s.sqlTargetTable, onDelete)
}

func (s SQLField) ForeignConstraint(tableGoName string) (ForeignKeyConstraint, bool) {
	targetTable := s.ForeignKey()
	ct := ForeignKeyConstraint{
		sqlField:       s.SQLName,
		sqlSourceTable: tableName(tableGoName),
		sqlTargetTable: targetTable,
		deleteAction:   s.onDeleteConstraint(),
	}
	return ct, targetTable != ""
}

// check if the constraint is of the form UNIQUE(col1, col2, ...)
var reUnique = regexp.MustCompile(`(?i)ADD UNIQUE\s?\((.*)\)`)

// return the column made unique by a UNIQUE(col) constraint
// or the empty string(
func IsUniqueConstraint(ct loader.Comment) string {
	matchs := reUnique.FindStringSubmatch(ct.Content)
	if len(matchs) > 0 {
		cols := strings.Split(matchs[1], ",")
		if len(cols) == 1 { // unique column
			return strings.TrimSpace(cols[0])
		}
	}
	return ""
}
