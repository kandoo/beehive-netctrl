package nom

import (
	"bytes"
	"encoding/gob"
	"reflect"

	"github.com/golang/glog"
	"github.com/soheilhy/beehive/bh"
)

// Object is the interface of all structs in the network object model.
type Object interface {
	// GobDecode decodes the object from a byte array using the Gob encoding.
	GobDecode(b []byte) error
	// GobEncode encodes the object into a byte array using the Gob encoding.
	GobEncode() ([]byte, error)
	// JSONDecode decodes the object from a byte array using the JSON encoding.
	JSONDecode(b []byte) error
	// JSONEncode encodes the object into a byte array using the JSON encoding.
	JSONEncode() ([]byte, error)
	// UID returns a unique ID of this object. This ID is unique in the network
	// among all other objects.
	UID() UID
}

// GobDecode decodes the object from b using Gob.
func ObjGobDecode(obj interface{}, b []byte) error {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	return dec.Decode(obj)
}

// GobEncode encodes the object into a byte array using Gob.
func ObjGobEncode(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DictGet(d bh.Dictionary, k bh.Key, obj Object) error {
	v, err := d.Get(k)
	if err != nil {
		return err
	}

	if err = obj.GobDecode(v); err != nil {
		glog.Errorf("Error in decoding %s from dictionary %s: %v",
			reflect.TypeOf(obj).String(), d.Name(), err)
		return err
	}

	return nil
}

func DictPut(d bh.Dictionary, k bh.Key, obj Object) error {
	v, err := obj.GobEncode()
	if err != nil {
		return err
	}

	d.Put(k, v)
	return nil
}
