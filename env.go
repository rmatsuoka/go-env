package env

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"golang.org/x/exp/constraints"
)

type envSetter struct {
	v     value
	doc   string
	typ   string
	where where
}

type where struct {
	file string
	line int
}

func (w where) String() string {
	return fmt.Sprintf("%s:%d", w.file, w.line)
}

func whereami() where {
	var pcs [1]uintptr
	// skip [runtime.Callers, this function]
	runtime.Callers(3, pcs[:])
	f, _ := runtime.CallersFrames(pcs[:]).Next()
	return where{
		file: f.File,
		line: f.Line,
	}
}

var defaultEnvSetter = make(map[string]*envSetter)

func Signed[T constraints.Signed](key string, value T, doc string) *T {
	i := value
	where := whereami()
	defaultEnvSetter[key] = &envSetter{
		typ:   fmt.Sprintf("%T", i),
		where: where,
		doc:   doc,
		v:     signedValue[T]{&i},
	}
	return (*T)(&i)
}

func String(key, value, doc string) *string {
	str := value
	where := whereami()
	defaultEnvSetter[key] = &envSetter{
		typ:   "string",
		where: where,
		doc:   doc,
		v:     stringValue{&str},
	}
	return &str
}

func Parse() {
	for key, setter := range defaultEnvSetter {
		v, ok := os.LookupEnv(key)
		if !ok {
			continue
		}
		if err := setter.v.set(v); err != nil {
			log.Println(err)
		}
	}
	if os.Getenv("HELP") != "" {
		usage()
	}
}

func usage() {
	for key, setter := range defaultEnvSetter {
		fmt.Fprintf(
			os.Stderr,
			"# %s\n# at %s\n%s='%s'\n\n",
			setter.doc, setter.where, key, setter.v,
		)
	}
	os.Exit(1)
}

type value interface {
	set(string) error
	String() string
}

type signedValue[T constraints.Signed] struct {
	ptr *T
}

func (v signedValue[T]) set(s string) error {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*v.ptr = T(i)
	return err
}

func (v signedValue[T]) String() string {
	return strconv.FormatInt(int64(*v.ptr), 10)
}

type stringValue struct {
	ptr *string
}

func (v stringValue) set(s string) error {
	*v.ptr = s
	return nil
}

func (v stringValue) String() string { return *v.ptr }
