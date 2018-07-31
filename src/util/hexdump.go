package util

import (
	"bufio"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/util/collections"
)

// Annotation : Denotes a chunk of data to be dumped
type Annotation struct {
	Name string
	Size int
}

// IAnnotationsGenerator : Interface to implement by types to use HexDump
type IAnnotationsGenerator interface {
	GenerateAnnotations() []Annotation
}

// IAnnotationsIterator : Interface to implement by types to use HexDumpFromIterator
type IAnnotationsIterator interface {
	Next() (Annotation, bool)
}

func writeHexdumpMember(offset int, size int, writer io.Writer, buffer []byte, name string) {
	var hexBuff = make([]string, size)
	var j = 0
	if offset+size > len(buffer) {
		panic(encoder.ErrBufferUnderflow)
	}
	for i := offset; i < offset+size; i++ {
		hexBuff[j] = strconv.FormatInt(int64(buffer[i]), 16)
		j++
	}
	for i := 0; i < len(hexBuff); i++ {
		if len(hexBuff[i]) == 1 {
			hexBuff[i] = "0" + hexBuff[i]
		}
	}

	var sliceContents = getSliceContentsString(hexBuff, offset)
	var serialized = encoder.Serialize(sliceContents + " " + name + "\n")

	f := bufio.NewWriter(writer)
	defer f.Flush()
	f.Write(serialized[4:])

}

func getSliceContentsString(sl []string, offset int) string {
	var res = ""
	var counter = 0
	var currentOff = offset
	if offset != -1 {
		var hex = strconv.FormatInt(int64(offset), 16)
		var l = len(hex)
		for i := 0; i < 4-l; i++ {
			hex = "0" + hex
		}
		hex = "0x" + hex
		res += hex + " | "
	}
	for i := 0; i < len(sl); i++ {
		counter++
		res += sl[i] + " "
		if counter == 16 {
			if i != len(sl)-1 {
				res = strings.TrimRight(res, " ")
				res += "\n"
				currentOff += 16
				if offset != -1 {
					//res += "         " //9 spaces
					var hex = strconv.FormatInt(int64(currentOff), 16)
					var l = len(hex)
					for i := 0; i < 4-l; i++ {
						hex = "0" + hex
					}
					hex = "0x" + hex
					res += hex + " | "
				}
				counter = 0
			} else {
				res += "..."
				return res
			}
		}
	}
	for i := 0; i < (16 - counter); i++ {
		res += "..."
	}
	res += "..."
	return res
}

func printFinalHex(i int, writer io.Writer) {
	var finalHex = strconv.FormatInt(int64(i), 16)
	var l = len(finalHex)
	for i := 0; i < 4-l; i++ {
		finalHex = "0" + finalHex
	}
	finalHex = "0x" + finalHex
	finalHex = finalHex + " | "

	var serialized = encoder.Serialize(finalHex)

	f := bufio.NewWriter(writer)
	defer f.Flush()
	f.Write(serialized[4:])
}

// HexDump : Returns hexdump of buffer according to annotations, via writer
func HexDump(buffer []byte, annotations []Annotation, writer io.Writer) {
	var currentOffset = 0

	for _, element := range annotations {
		writeHexdumpMember(currentOffset, element.Size, writer, buffer, element.Name)
		currentOffset += element.Size
	}

	printFinalHex(currentOffset, writer)
}

// HexDumpFromIterator : Returns hexdump of buffer according to annotationsIterator, via writer
func HexDumpFromIterator(buffer []byte, annotationsIterator IAnnotationsIterator, writer io.Writer) {
	var currentOffset = 0

	var current, valid = annotationsIterator.Next()

	for {
		if !valid {
			break
		}
		writeHexdumpMember(currentOffset, current.Size, writer, buffer, current.Name)
		currentOffset += current.Size
		current, valid = annotationsIterator.Next()
	}

	printFinalHex(currentOffset, writer)
}

// maiState internal state for struct in message annotation iterator
type eaiState struct {
	Object       reflect.Value
	CurrentField int
	MaxField     int
	CurrentIndex int
	FieldPath    string
}

// EncoderAnnotationsIterator iterate over buffer annotations delimiting encoded object fields
type EncoderAnnotationsIterator struct {
	iterStates *collections.Stack
	state      *eaiState
}

// NewEncoderAnnotationsIterator : Initializes struct EncoderAnnotationsIterator
func NewEncoderAnnotationsIterator(object interface{}) EncoderAnnotationsIterator {
	var eaiState = eaiState{}
	eaiState.Object = reflect.ValueOf(object)
	eaiState.CurrentField = 0
	eaiState.MaxField = reflect.Indirect(reflect.ValueOf(eaiState.Object)).NumField()
	eaiState.CurrentIndex = -1
	eaiState.FieldName = ""

	var eai = EncoderAnnotationsIterator{
		&collections.Stack{},
		eaiState,
	}
	eai.state = &eaiState

	return eai
}

// Next : Yields next element of EncoderAnnotationsIterator
func (mai *EncoderAnnotationsIterator) Next() (util.Annotation, bool) {
	if !mai.LengthCalled {
		mai.LengthCalled = true
		return util.Annotation{Size: 4, Name: "Length"}, true
	}
	if !mai.PrefixCalled {
		mai.PrefixCalled = true
		return util.Annotation{Size: 4, Name: "Prefix"}, true

	}
	for mai.state.CurrentField >= mai.state.MaxField {
		if tip, exists := mai.fieldIterators.Pop(); exists {
			mai.state, _ = tip.(*maiState)
		} else {
			return util.Annotation{}, false
		}
	}

	var i = mai.state.CurrentField
	var j = mai.state.CurrentIndex

	var v = reflect.Indirect(reflect.ValueOf(mai.state.Object))
	t := v.Type()
	vF := v.Field(i)
	f := t.Field(i)
	for f.PkgPath != "" && i < mai.state.maxField {
		i++
		mai.state.CurrentField++
		mai.state.CurrentIndex = -1
		j = -1
		if i < mai.state.MaxField {
			f = t.Field(i)
			if f.Type.Kind() == reflect.Slice {
				if _, omitempty := encoder.ParseTag(f.Tag.Get("enc")); omitempty {
					if i == mai.state.MaxField-1 && mai.fieldsIterator.Len() == 0 {
						vF = v.Field(i)
						if vF.Len() == 0 {
							// Last field is empty slice. Nothing further tokens
							return util.Annotation{}, false
						}
					} else {
						panic(encoder.ErrInvalidOmitEmpty)
					}
				}
			}
		} else {
			return util.Annotation{}, false
		}
	}
	if f.Tag.Get("enc") != "-" {
		if vF.CanSet() || f.Name != "_" {
			if v.Field(i).Kind() == reflect.Slice {
				if mai.state.CurrentIndex == -1 {
					mai.state.CurrentIndex = 0
					return util.Annotation{
						Size: 4,
						Name: mai.state.FieldPath + f.Name + "#length",
					}, true
				}
				sliceLen := v.Field(i).Len()
				mai.CurrentIndex++
				if mai.CurrentIndex < sliceLen {
					// Emit annotation for slice item
					return util.Annotation{Size: len(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())[4:]), Name: f.Name + "[" + strconv.Itoa(j) + "]"}, true
				}
				// No more annotation tokens for current slice field
				mai.CurrentIndex = -1
				mai.CurrentField++
				if sliceLen > 0 {
					// Emit annotation for last item
					return util.Annotation{Size: len(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())[4:]), Name: f.Name + "[" + strconv.Itoa(j) + "]"}, true
				}
				// Zero length slice. Start over
				return mai.Next()
			}

			mai.CurrentField++
			return util.Annotation{Size: len(encoder.Serialize(v.Field(i).Interface())), Name: f.Name}, true

		}
	}

	return util.Annotation{}, false
}
