// pbvalidate validates pbjson files against a protobuf message
//
// Usage:
//
//    pbvalidate -f $workspace/auth/api/config.proto -I /,$workspace,$workspace/vendor/github.com/googleapis/googleapis -m auth.NamespaceConfig  /tmp/namespace_config.json
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/mkmik/stringlist"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"k8s.io/klog"
)

var (
	protoFileName = flag.String("f", "", "Proto schema files")
	protoMessage  = flag.String("m", "", "Proto message")
	importPaths   = stringlist.Flag("I", "Proto include path")
	strict        = flag.Bool("strict", false, "Strict validation")
)

func findMessage(fds []*desc.FileDescriptor, name string) *desc.MessageDescriptor {
	for _, fd := range fds {
		if md := fd.FindMessage(name); md != nil {
			return md
		}
	}
	return nil
}

// Two steps validation
// 1 - Make sure that the json file can be unmarshalled with the provided .proto file
// 2 (if strict flag provided) - Make sure that the result of marshalling is the same than the provided json file
// It will detect properties not being camelCase, integers quoted as strings and so on which are proto3 json serializer valid.
// This is specially useful when the json files are consumed by non protobuffers APIs
func run(fileName string, protoMessage string, importPaths []string, src string, strict bool) error {
	p := &protoparse.Parser{
		ImportPaths: importPaths,
	}
	fds, err := p.ParseFiles(fileName)
	if err != nil {
		return fmt.Errorf("parsing %q: %w", fileName, err)
	}
	md := findMessage(fds, protoMessage)

	if md == nil {
		return fmt.Errorf("cannot find message %q", protoMessage)
	}
	m := dynamic.NewMessage(md)

	b, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading file %q: %w", src, err)
	}
	if err := m.UnmarshalJSON(b); err != nil {
		return fmt.Errorf("parsing %q: %w", src, err)
	}

	if !strict {
		return nil
	}

	// ReMarshall the file and make sure that the result is the same than the provided json file (src)
	reMarshall, err := m.MarshalJSON()
	if err != nil {
		return fmt.Errorf("re-marshalling error %w", err)
	}

	differ := diff.New()
	d, err := differ.Compare(b, reMarshall)
	if err != nil {
		return fmt.Errorf("diff error %w", err)
	}

	if d.Modified() {
		var leftJSON map[string]interface{}
		json.Unmarshal(b, &leftJSON)

		formatter := formatter.NewAsciiFormatter(leftJSON, formatter.AsciiFormatterConfig{})
		diffString, err := formatter.Format(d)
		if err != nil {
			return fmt.Errorf("diff format error %w", err)
		}
		return fmt.Errorf("The provided file does not match the schema\n %s", diffString)
	}

	return nil
}

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()

	flag.Parse()
	src := flag.Arg(0)

	if err := run(*protoFileName, *protoMessage, *importPaths, src, *strict); err != nil {
		klog.Exitf("%v", err)
	}
}
