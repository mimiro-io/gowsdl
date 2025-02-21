// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var typesTmpl = `
{{define "SimpleType"}}
	{{$typeName := replaceReservedWords .Name | makePublic}}
	{{if .Doc}} {{.Doc | comment}} {{end}}
	{{if ne .List.ItemType ""}}
		type {{$typeName}} []{{toGoType .List.ItemType false | removePointerFromType}}
	{{else if ne .Union.MemberTypes ""}}
		type {{$typeName}} string
	{{else if .Union.SimpleType}}
		type {{$typeName}} string
	{{else if .Restriction.Base}}
		type {{$typeName}} {{toGoType .Restriction.Base false | removePointerFromType}}
    {{else}}
		type {{$typeName}} interface{}
	{{end}}

	{{if .Restriction.Enumeration}}
	const (
		{{with .Restriction}}
			{{range .Enumeration}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{$typeName}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$typeName}} = "{{goString .Value}}" {{end}}
		{{end}}
	)
	{{end}}
{{end}}

{{define "ComplexContent"}}
	{{$baseType := toGoType .Extension.Base false}}
	{{ if $baseType }}
		{{$baseType}}
	{{end}}

	{{template "Elements" .Extension.Sequence}}
	{{template "Elements" .Extension.Choice}}
	{{template "Elements" .Extension.SequenceChoice}}
	{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "Attributes"}}
    {{ $targetNamespace := getNS }}
	{{range .}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
		{{ if ne .Type "" }}
			{{ normalize .Name | makeFieldPublic}} {{toGoType .Type false}} ` + "`" + `xml:"{{with $targetNamespace}}{{.}} {{end}}{{.Name}},attr,omitempty" json:"{{.Name}},omitempty"` + "`" + `
		{{ else }}
			{{ normalize .Name | makeFieldPublic}} string ` + "`" + `xml:"{{with $targetNamespace}}{{.}} {{end}}{{.Name}},attr,omitempty" json:"{{.Name}},omitempty"` + "`" + `
		{{ end }}
	{{end}}
{{end}}

{{define "SimpleContent"}}
	Value {{toGoType .Extension.Base false}} ` + "`xml:\",chardata\" json:\"-,\"`" + `
	{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "ComplexTypeInline"}}
	{{replaceReservedWords .Name | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}struct {
	{{with .ComplexType}}
		{{if ne .ComplexContent.Extension.Base ""}}
			{{template "ComplexContent" .ComplexContent}}
		{{else if ne .SimpleContent.Extension.Base ""}}
			{{template "SimpleContent" .SimpleContent}}
		{{else}}
			{{template "Elements" .Sequence}}
			{{template "Elements" .Choice}}
			{{template "Elements" .SequenceChoice}}
			{{template "Elements" .All}}
			{{template "Attributes" .Attributes}}
		{{end}}
	{{end}}
	} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
{{end}}

{{define "Elements"}}
	{{range .}}
		{{if ne .Ref ""}}
			{{removeNS .Ref | replaceReservedWords  | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{toGoType .Ref .Nillable }} ` + "`" + `xml:"{{.Ref | removeNS}},omitempty" json:"{{.Ref | removeNS}},omitempty"` + "`" + `
		{{else}}
		{{if not .Type}}
			{{if .SimpleType}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .SimpleType.List.ItemType ""}}
					{{ normalize .Name | makeFieldPublic}} []{{toGoType .SimpleType.List.ItemType false}} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
				{{else}}
					{{ normalize .Name | makeFieldPublic}} {{toGoType .SimpleType.Restriction.Base false}} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
				{{end}}
			{{else}}
				{{template "ComplexTypeInline" .}}
			{{end}}
		{{else}}
			{{if .Doc}}{{.Doc | comment}} {{end}}
			{{replaceAttrReservedWords .Name | makeFieldPublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{toGoType .Type .Nillable }} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + ` {{end}}
		{{end}}
	{{end}}
{{end}}

{{define "Any"}}
	{{range .}}
		Items     []string ` + "`" + `xml:",any" json:"items,omitempty"` + "`" + `
	{{end}}
{{end}}

{{range .Schemas}}
	{{ $targetNamespace := setNS .TargetNamespace }}

	{{range .SimpleType}}
		{{template "SimpleType" .}}
	{{end}}

	{{range .Elements}}
		{{$name := .Name}}
		{{$typeName := replaceReservedWords $name | makePublic}}
		{{if not .Type}}
			{{/* ComplexTypeLocal */}}
			{{with .ComplexType}}
				type {{$typeName}} struct {
					XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$name}}\"`" + `
					{{if ne .ComplexContent.Extension.Base ""}}
						{{template "ComplexContent" .ComplexContent}}
					{{else if ne .SimpleContent.Extension.Base ""}}
						{{template "SimpleContent" .SimpleContent}}
					{{else}}
						{{template "Elements" .Sequence}}
						{{template "Any" .Any}}
						{{template "Elements" .Choice}}
						{{template "Elements" .SequenceChoice}}
						{{template "Elements" .All}}
						{{template "Attributes" .Attributes}}
					{{end}}
				}
			{{end}}
			{{/* SimpleTypeLocal */}}
			{{with .SimpleType}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .List.ItemType ""}}
					type {{$typeName}} []{{toGoType .List.ItemType false | removePointerFromType}}
				{{else if ne .Union.MemberTypes ""}}
					type {{$typeName}} string
				{{else if .Union.SimpleType}}
					type {{$typeName}} string
				{{else if .Restriction.Base}}
					type {{$typeName}} {{toGoType .Restriction.Base false | removePointerFromType}}
				{{else}}
					type {{$typeName}} interface{}
				{{end}}

				{{if .Restriction.Enumeration}}
				const (
					{{with .Restriction}}
						{{range .Enumeration}}
							{{if .Doc}} {{.Doc | comment}} {{end}}
							{{$typeName}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$typeName}} = "{{goString .Value}}" {{end}}
					{{end}}
				)
				{{end}}
			{{end}}
		{{else}}
			{{$type := toGoType .Type .Nillable | removePointerFromType}}
			{{if ne ($typeName) ($type)}}
				type {{$typeName}} {{$type}}
				{{if eq ($type) ("soap.XSDDateTime")}}
					func (xdt {{$typeName}}) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
						return soap.XSDDateTime(xdt).MarshalXML(e, start)
					}

					func (xdt *{{$typeName}}) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
						return (*soap.XSDDateTime)(xdt).UnmarshalXML(d, start)
					}
				{{else if eq ($type) ("soap.XSDDate")}}
					func (xd {{$typeName}}) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
						return soap.XSDDate(xd).MarshalXML(e, start)
					}

					func (xd *{{$typeName}}) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
						return (*soap.XSDDate)(xd).UnmarshalXML(d, start)
					}
				{{else if eq ($type) ("soap.XSDTime")}}
					func (xt {{$typeName}}) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
						return soap.XSDTime(xt).MarshalXML(e, start)
					}

					func (xt *{{$typeName}}) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
						return (*soap.XSDTime)(xt).UnmarshalXML(d, start)
					}
				{{end}}
			{{end}}
		{{end}}
	{{end}}
	{{ $validNamepace := false }}
	{{ range $key, $value := .Xmlns }}
		{{ if eq $value $targetNamespace }}
			{{ $validNamepace = true }}
		{{ end }}
	{{ end }}

	{{range .ComplexTypes}}
		{{/* ComplexTypeGlobal */}}
		{{$typeName := replaceReservedWords .Name | makePublic}}
		{{if and (eq (len .SimpleContent.Extension.Attributes) 0) (eq (toGoType .SimpleContent.Extension.Base false) "string") }}
			type {{$typeName}} string
		{{else}}
			type {{$typeName}} struct {
				{{$type := findNameByType .Name}}
				{{if and (notHasSuffix $type "Response") ($validNamepace)}}
					{{if ne .Name $type}}
						XMLName xml.Name ` + "`xml:\"{{$type}}\"`" + `
					{{else}}	
						XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$type}}\"`" + `
					{{end}}	
				{{end}}

				{{if ne .ComplexContent.Extension.Base ""}}
					{{template "ComplexContent" .ComplexContent}}
				{{else if ne .SimpleContent.Extension.Base ""}}
					{{template "SimpleContent" .SimpleContent}}
				{{else}}
					{{template "Elements" .Sequence}}
					{{template "Any" .Any}}
					{{template "Elements" .Choice}}
					{{template "Elements" .SequenceChoice}}
					{{template "Elements" .All}}
					{{template "Attributes" .Attributes}}
				{{end}}
			}

{{if notHasSuffix $type "Response"}}
func (r {{$typeName}}) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := reflect.ValueOf(r)
	t := reflect.TypeOf(r)
	if start.Name.Space == "" {
		parentXmlNameField, _ := t.FieldByName("XMLName")
		parentTypeTag := parentXmlNameField.Tag.Get("xml")
		if parentTypeTag != "" {
			tokens := strings.Split(parentTypeTag, " ")
			//start.Name.Space = tokens[0]
			if len(tokens) > 1 {
				start.Name.Local = "ns2:" + tokens[1]
			}
			start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns:ns2"}, Value: tokens[0]})
		}
	} else {
		start.Name.Space = ""
	}
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}
	fields := reflect.TypeOf(r).NumField()
	for i := 0; i < fields; i++ {
		f := reflect.TypeOf(r).Field(i)
		if f.Name != "XMLName" && !f.Anonymous {
			fieldRef := v.FieldByName(f.Name)
			ft := fieldRef.Type()
			n := xml.Name{}
			if ft.Kind() == reflect.Ptr {
				elm :=  ft.Elem()
				if elm.Kind() == reflect.Ptr {
					xmlNameType, _ := elm.FieldByName("XMLName")
					typeTag := xmlNameType.Tag.Get("xml")
					if typeTag != "" {
						tokens := strings.Split(typeTag, " ")
						n.Space = tokens[0]
						if len(tokens) > 1 {
							n.Local = tokens[1]
						}
					}
				}
			}
			//
			xmlTag := f.Tag.Get("xml")
			omitEmpty := false
			if xmlTag != "" {
				tokens := strings.Split(xmlTag, ",")
				xmlTagName := tokens[0]
				n.Local = xmlTagName
				if len(start.Attr) > 0 {
					n.Space = start.Attr[0].Value
				}
				if len(tokens) > 1 {
					omitEmpty = tokens[1] == "omitempty"
				}
			}
			if !omitEmpty || !reflect.ValueOf(r).Field(i).IsZero() {
				err := e.EncodeElement(reflect.ValueOf(r).Field(i).Interface(), xml.StartElement{Name: n})
				if err != nil {
					return err
				}
			}
		}
	}
	err = e.EncodeToken(start.End())
	if err != nil {
		return err
	}
	return nil
}
{{end}}


		{{end}}
	{{end}}
{{end}}
`
