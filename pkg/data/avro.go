package data

import "github.com/hamba/avro"

var (
	avroInts    = avro.MustParse(`{"type": "array", "items": "long"}`)
	avroFloats  = avro.MustParse(`{"type": "array", "items": "double"}`)
	avroStrings = avro.MustParse(`{"type": "array", "items": "string"}`)
	avroObjects = avro.MustParse(`{"type": "array", "items": {
			"type": "record",
			"name": "Object",
			"fields" : [
				{"name": "I", "type": "long"},
				{"name": "F", "type": "double"},
				{"name": "T", "type": "boolean"},
				{"name": "S", "type": "string"},
				{"name": "B", "type": "bytes"}
			]
		}
	}`)
)
