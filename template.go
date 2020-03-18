package gofast

import (
	"bytes"
	"fmt"
	"io"
)

type message struct {
	content string
}

//Template is template to decode message
type Template struct {
	Units []TemplateUnit
	pmap  *PMap
}

//TemplateUnit is the unit for FAST message Decoding
type TemplateUnit interface {
	Decode(in io.Reader, pmap *PMap, msg *bytes.Buffer) error
	Reset()
}

//Sequence is one kind of template unit. And it's comprised of this:  1. group length ； 2. field group repeats group length of times.
type Sequence struct {
	value []byte
	group []TemplateUnit //Field group in this seqence.
	glen  uint           //Length of field group in this sequence. This should be always in the
	gpmap PMap           //Pmap of the group
}

//Field is one kind of template unit
type Field struct {
	tag          string
	preValue     interface{}
	isMandatory  bool
	initValue    interface{}
	decodeByType DecodeFunc
}

//Decode decodes the field with no operator. Field with NoOperator does not require to check PMap as the field should be always in FAST message.
func (f *Field) Decode(in io.Reader, pmap *PMap, msg *bytes.Buffer) error {
	buf := ReadBinaryByStopbit(in)
	fieldContentInString, err := f.decodeByType(buf, !f.isMandatory, f)
	if err != nil {
		return err
	}
	if fieldContentInString != nil {
		appendField(msg, f.tag, *fieldContentInString)
	} else {
		if f.isMandatory {
			return fmt.Errorf("null value decoded of mandatory field: %s ", f.tag)
		}
	}
	return nil
}

//Reset resets the field with no operator. Since it does not use previous status, Reset does nothing.
func (f *Field) Reset() {}

//FieldWithConstantOperator is field with Constant operator
type FieldWithConstantOperator Field

//Decode decodes the field with Constant operator. Field with constant operator never encoded in FAST message. When PMap bit of the field is set,
//the decoder should set this field with a constant in output message.
func (f *FieldWithConstantOperator) Decode(in io.Reader, pmap *PMap, msg *bytes.Buffer) error {
	if f.isMandatory || pmap.HasNextPresenceBit() {
		appendField(msg, f.tag, fmt.Sprintf("%v", f.initValue))
		return nil
	}
	return nil
}

//Reset resets field with ConstantOperator.
func (f *FieldWithConstantOperator) Reset() {}

//FieldWithDefaultOperator is
type FieldWithDefaultOperator Field

//Decode is decode function of FieldWithDefaultOperator
func (f *FieldWithDefaultOperator) Decode(in io.Reader, pmap *PMap, msg *bytes.Buffer) error {
	var fieldContentInString *string
	if pmap.HasNextPresenceBit() {
		var err error
		fast := ReadBinaryByStopbit(in)
		fieldContentInString, err = f.decodeByType(fast, !f.isMandatory, (*Field)(f))
		if err != nil {
			return err
		}
	} else {
		if f.initValue != nil {
			*fieldContentInString = fmt.Sprintf("%v", f.initValue)
		}
	}

	if fieldContentInString == nil {
		if f.isMandatory {
			return fmt.Errorf("null value decoded of mandatory field: %s ", f.tag)
		}
	} else {
		appendField(msg, f.tag, *fieldContentInString)
	}
	return nil
}

//Reset method of FieldWithDefaultOperator
func (f *FieldWithDefaultOperator) Reset() {
	//Nothing need to do as the default value(initValue) does not change after creation.
}

//FieldWithCopyOperator is field with copy operator
type FieldWithCopyOperator Field

//Decode method of FieldWithCopyOperator
func (f *FieldWithCopyOperator) Decode(in io.Reader, pmap *PMap, msg *bytes.Buffer) error {
	var fieldContentInString *string
	if pmap.HasNextPresenceBit() {
		var err error
		fast := ReadBinaryByStopbit(in)
		fieldContentInString, err = f.decodeByType(fast, f.isMandatory, (*Field)(f))
		if err != nil {
			return err
		}
	} else {
		if f.preValue != nil {
			*fieldContentInString = fmt.Sprintf("%v", f.preValue)
		}
	}
	if fieldContentInString == nil {
		if f.isMandatory {
			return fmt.Errorf("null value decoded of mandatory field: %s ", f.tag)
		}
	} else {
		appendField(msg, f.tag, *fieldContentInString)
	}
	return nil
}

//Reset method of FieldWithCopyOperator
func (f *FieldWithCopyOperator) Reset() {
	f.preValue = f.initValue
}

type FieldWithIncreaseOperator Field

func (f *FieldWithIncreaseOperator) Decode(in io.Reader, pmap *PMap, msg *bytes.Buffer) error {

}
func (f *FieldWithIncreaseOperator) Reset() {
}

type FieldWithDeltaOperator Field

func (f *FieldWithDeltaOperator) Decode(in io.Reader, pmap *PMap, msg *bytes.Buffer) error {

}
func (f *FieldWithDeltaOperator) Reset() {
}

type FieldWithTailOperator Field

func (f *FieldWithTailOperator) Decode(in io.Reader, pmap *PMap, msg *bytes.Buffer) error {

}
func (f *FieldWithTailOperator) Reset() {
}

//TypeDecoder is the decoder of Field. Different type of field should has its own type of decoder.
type DecodeFunc func(fastMsg []byte, isNullable bool, f *Field) (*string, error)

func appendField(msg *bytes.Buffer, tag, value string) {
	msg.WriteString(tag)
	msg.WriteRune('=')
	msg.WriteString(value)
}
