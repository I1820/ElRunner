/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 21-08-2018
 * |
 * | File Name:     codec_test.go
 * +===============================================
 */

package actions

var (
	herID = "i1820-el"

	herTextName = "Kiana"
	herCborName = []byte{0x67, 0x22, 0x4B, 0x69, 0x61, 0x6E, 0x61, 0x22}

	locationCbor   = []byte{0xA3, 0x63, 0x6C, 0x61, 0x74, 0x0A, 0x63, 0x6C, 0x6E, 0x67, 0x0A, 0x6B, 0x74, 0x65, 0x6D, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x1E}
	locationObject = map[string]interface{}{
		"temperature": 30.0,
		"_location": map[string]interface{}{
			"type":        "Point",
			"coordinates": []interface{}{10.0, 10.0},
		},
		"lat": 10.0,
		"lng": 10.0,
		"loc": map[string]interface{}{
			"type":        "Point",
			"coordinates": []interface{}{10.0, 10.0},
		},
	}

	codecOne = `
from codec import Codec
import cbor
class I1820(Codec):
    def decode(self, data):
        return cbor.loads(data)
    def encode(self, data):
        return cbor.dumps(data)
`
	codecTwo = `
from codec import Codec
import cbor
class I1820(Codec):
    thing_location = 'loc'
    def decode(self, data):
        print("Hello")
        d = cbor.loads(data)
        print(d)
        d['loc'] = self.create_location(d['lat'], d['lng'])
        return d
    def encode(self, data):
        return cbor.dumps(data)
`
)

func (as *ActionSuite) Test_CodecsResource_Encode_Decode_1() {
	// Create
	var id string
	resc := as.JSON("/api/codecs").Post(codeReq{ID: herID, Code: codecOne})
	as.Equalf(200, resc.Code, "Error: %s", resc.Body.String())
	resc.Bind(&id)
	as.Equal(id, herID)

	// Decode
	var dResult interface{}
	resd := as.JSON("/api/codecs/%s/decode", herID).Post(herCborName)
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
	resd.Bind(&dResult)
	as.Equal("\""+herTextName+"\"", dResult)

	// Encode
	var eResult []byte
	rese := as.JSON("/api/codecs/%s/encode", herID).Post(herTextName)
	as.Equalf(200, rese.Code, "Error: %s", rese.Body.String())
	rese.Bind(&eResult)
	as.Equal(herCborName, eResult)
}

func (as *ActionSuite) Test_CodecsResource_Encode_Decode_2() {
	// Create
	var id string
	resc := as.JSON("/api/codecs").Post(codeReq{ID: herID, Code: codecTwo})
	as.Equalf(200, resc.Code, "Error: %s", resc.Body.String())
	resc.Bind(&id)
	as.Equal(id, herID)

	// Decode
	var dResult interface{}
	resd := as.JSON("/api/codecs/%s/decode", herID).Post(locationCbor)
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
	resd.Bind(&dResult)
	as.Equal(locationObject, dResult)
}
