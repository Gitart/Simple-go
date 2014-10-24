/ RethinkDB utilities functions for gorethink driver
 
package rdbu
 
import (
	"fmt"
	rdb "github.com/dancannon/gorethink"
	"reflect"
)
 
// Map is a convinient type alias for DB operations
type Map map[string]interface{}
 
// Error indices this package error.
type Error struct {
	msg string
	e   error
}
 
func (err Error) Error() string {
	if err.e == nil {
		return err.msg
	}
	return err.msg + err.e.Error()
}
 
// Exists returns true if `expr` query returns results.
// expr: ReQL *singleRowSelection* or scalar. Use `ExistsMulti` otherwise.
func Exists(expr rdb.RqlTerm, s *rdb.Session) (bool, error) {
	r, err := expr.RunRow(s)
	return err == nil && !r.IsNil(), err
}
 
// ExistsMulti returns true if `expr` query is not empty.
// expr: ReQL *selection*.
func ExistsMulti(expr rdb.RqlTerm, s *rdb.Session) (bool, error) {
	r, err := expr.IsEmpty().RunRow(s)
	if err != nil {
		return false, err
	}
	var isEmpty bool
	err = r.Scan(&isEmpty)
	return err == nil && !r.IsNil() && !isEmpty, err
}
 
// Scan fetches a result from a *singleRowSelection* into a reference to `dest`.
// dest: pointer to single element.
func Scan(expr rdb.RqlTerm, dest interface{}, s *rdb.Session) (bool, error) {
	row, err := expr.RunRow(s)
	if err != nil {
		return false, err
	}
	if row.IsNil() {
		return false, nil
	}
	return true, row.Scan(dest)
}
 
// ScanAll returns the results of *selection* query in user supplied slice
// slice: pointer to slice
func ScanAll(expr rdb.RqlTerm, slice interface{}, s *rdb.Session) error {
	rows, err := expr.Run(s)
	if err != nil {
		return err
	}
	slicePtrValue := reflect.ValueOf(slice)
	if slicePtrValue.Kind() != reflect.Ptr {
		return Error{"GetAllRows: You need to pass a pointer to slice.", nil}
	}
	sliceValue := slicePtrValue.Elem()
	if sliceValue.Kind() != reflect.Slice {
		return Error{"GetAllRows: You need to pass a slice.", nil}
	}
 
	tempSlice := reflect.MakeSlice(sliceValue.Type(), 0, 0)
	for rows.Next() {
		elem := reflect.New(sliceValue.Type().Elem())
		if err := rows.Scan(elem.Interface()); err != nil {
			return Error{fmt.Sprint("Scan failed. Destination element type was: ",
				reflect.TypeOf(elem.Interface())), err}
		}
		tempSlice = reflect.Append(tempSlice, elem.Elem())
	}
	sliceValue.Set(tempSlice)
	return nil
}
 
// ScanFirst - It's like `ScanAll` but extracts only the first element.
// Returns true if a query result contains at least one element.
// dest: pointer to single element.
func ScanFirst(expr rdb.RqlTerm, dest interface{}, s *rdb.Session) (bool, error) {
	rows, err := expr.Run(s)
	if err != nil {
		return false, err
	}
	if !rows.Next() || rows.IsNil() {
		return false, nil
	}
	return true, rows.Scan(dest)
 
}