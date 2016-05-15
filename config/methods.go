package config

import (
    "fmt"
)

// Getf сокращение для fmt.Sprintf
func Getf(format string, a ...interface {}) string {
    return fmt.Sprintf(format, a...)
}

// get выполняет роль тернарного оператора:
// возвращает первый параметр, если он не пустой, в противном случае - второй
func get(check string, dflt string) string {
    if check == "" {
        return dflt
    }
    return check

}

// bget действует подобно get, за исключением того,
// что первый параметр должен быть булевым выражением
func bget(check bool, result string, dflt string) string {
    if check == true {
        return result
    }
    return dflt

}