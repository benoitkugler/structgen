package test

import "encoding/json"

// Code generated by structgen/interfaces. DO NOT EDIT

// union2Wrapper may be used as replacements for union2
// when working with JSON
type union2Wrapper struct {
	Data union2
}

func (out *union2Wrapper) UnmarshalJSON(src []byte) error {
	var wr struct {
		Data json.RawMessage
		Kind int
	}
	err := json.Unmarshal(src, &wr)
	if err != nil {
		return err
	}
	switch wr.Kind {
	case 0:
		var data member1
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data

	default:
		panic("exhaustive switch")
	}
	return err
}

func (item union2Wrapper) MarshalJSON() ([]byte, error) {
	type wrapper struct {
		Data interface{}
		Kind int
	}
	var wr wrapper
	switch data := item.Data.(type) {
	case member1:
		wr = wrapper{Kind: 0, Data: data}

	default:
		panic("exhaustive switch")
	}
	return json.Marshal(wr)
}

func (ct ITFSlice) MarshalJSON() ([]byte, error) {
	tmp := make([]union2Wrapper, len(ct))
	for i, v := range ct {
		tmp[i].Data = v
	}
	return json.Marshal(tmp)
}

func (ct *ITFSlice) UnmarshalJSON(data []byte) error {
	var tmp []union2Wrapper
	err := json.Unmarshal(data, &tmp)
	*ct = make(ITFSlice, len(tmp))
	for i, v := range tmp {
		(*ct)[i] = v.Data
	}
	return err
}

// union1Wrapper may be used as replacements for union1
// when working with JSON
type union1Wrapper struct {
	Data union1
}

func (out *union1Wrapper) UnmarshalJSON(src []byte) error {
	var wr struct {
		Data json.RawMessage
		Kind int
	}
	err := json.Unmarshal(src, &wr)
	if err != nil {
		return err
	}
	switch wr.Kind {
	case 0:
		var data member1
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case 1:
		var data member2
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case 2:
		var data member3
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data

	default:
		panic("exhaustive switch")
	}
	return err
}

func (item union1Wrapper) MarshalJSON() ([]byte, error) {
	type wrapper struct {
		Data interface{}
		Kind int
	}
	var wr wrapper
	switch data := item.Data.(type) {
	case member1:
		wr = wrapper{Kind: 0, Data: data}
	case member2:
		wr = wrapper{Kind: 1, Data: data}
	case member3:
		wr = wrapper{Kind: 2, Data: data}

	default:
		panic("exhaustive switch")
	}
	return json.Marshal(wr)
}
