// Code generated by structgen. DO NOT EDIT
	
	typedef JSON = Map<String, dynamic>; // alias to shorten JSON convertors

	
int intFromJson(dynamic json) {
			return json as int;
		}

		dynamic intToJson(int item) {
			return item;
		}
		
List<int> List_intFromJson(dynamic json) {
		return (json as List<dynamic>).map(intFromJson).toList();
	}

	dynamic List_intToJson(List<int> item) {
		return item.map(intToJson).toList();
	}
	

		// github.com/benoitkugler/structgen/dart-types/test.concret1
		class concret1 implements itfName {
		final List<int> List2;
final int V;

		concret1(this.List2, this.V);
		}
		
		
	concret1 concret1FromJson(JSON json) {
		return concret1(
			List_intFromJson(json['List2']),
intFromJson(json['V'])
		);
	}
	
	JSON concret1ToJson(concret1 item) {
		return {
			"List2" : List_intToJson(item.List2),
"V" : intToJson(item.V)
		};
	}
	
	
	
double doubleFromJson(dynamic json) {
			return json as double;
		}

		dynamic doubleToJson(double item) {
			return item;
		}
		

		// github.com/benoitkugler/structgen/dart-types/test.concret2
		class concret2 implements itfName {
		final double D;

		concret2(this.D);
		}
		
		
	concret2 concret2FromJson(JSON json) {
		return concret2(
			doubleFromJson(json['D'])
		);
	}
	
	JSON concret2ToJson(concret2 item) {
		return {
			"D" : doubleToJson(item.D)
		};
	}
	
	
	
// Corresponding Go code
	/*
	func itfNameUnmarshallJSON(src []byte) (itfName, error) {
		type wrapper struct {
			Data json.RawMessage
			Kind int
		}
		var wr wrapper
		err := json.Unmarshal(src, &wr)
		if err != nil {
			return nil, err
		}
		switch wr.Kind {
			case 0:
			var out concret1
			err = json.Unmarshal(wr.Data, &out)
			return out, err
	case 1:
			var out concret2
			err = json.Unmarshal(wr.Data, &out)
			return out, err
	
		default:
			panic("exhaustive switch")
		}
	}
	
func itfNameMarshallJSON(item itfName) ([]byte, error) {
		type wrapper struct {
				Data interface{}  
				Kind int         
		}
		var out wrapper
		switch item.(type) {
		case concret1:
			out = wrapper{Kind: 0, Data: item}
		case concret2:
			out = wrapper{Kind: 1, Data: item}
		
		default:
			panic("exhaustive switch")
		}
		return json.Marshal(out)
	}
	
	*/ 
	abstract class itfName {}
	itfName itfNameFromJson(dynamic json_) {
		final json = json_ as JSON;
		final kind = json['Kind'] as int;
		final data = json['Data'];
		switch (kind) {
			case  0:
			return concret1FromJson(data);
case  1:
			return concret2FromJson(data);
		default:
			throw ("unexpected type");
		}
	}
	
JSON itfNameToJson(itfName item) {
		if (item is concret1) {
			return {'Kind': 0, 'Data': concret1ToJson(item)};
		}else if (item is concret2) {
			return {'Kind': 1, 'Data': concret2ToJson(item)};
		} else {
			throw ("unexpected type");
		}	
	}
	
List<itfName> List_itfNameFromJson(dynamic json) {
		return (json as List<dynamic>).map(itfNameFromJson).toList();
	}

	dynamic List_itfNameToJson(List<itfName> item) {
		return item.map(itfNameToJson).toList();
	}
	
// github.com/benoitkugler/structgen/dart-types/test.ListV
typedef ListV = List<itfName>;


		// github.com/benoitkugler/structgen/dart-types/test.model
		class model  {
		final itfName Value;
final int A;
final ListV L;

		model(this.Value, this.A, this.L);
		}
		
		
	model modelFromJson(JSON json) {
		return model(
			itfNameFromJson(json['Value']),
intFromJson(json['A']),
List_itfNameFromJson(json['L'])
		);
	}
	
	JSON modelToJson(model item) {
		return {
			"Value" : itfNameToJson(item.Value),
"A" : intToJson(item.A),
"L" : List_itfNameToJson(item.L)
		};
	}
	
	
	
