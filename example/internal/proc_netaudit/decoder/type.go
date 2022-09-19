package decoder

type ApplicationPortocol uint32

var decoderMap = map[ApplicationPortocol]Decoder{}

func RegisterDecoder(num uint32, decoder Decoder) ApplicationPortocol {
	if _, ok := decoderMap[ApplicationPortocol(num)]; ok {
		panic("protocol already registered")
	}
	decoderMap[ApplicationPortocol(num)] = decoder
	return ApplicationPortocol(num)
}
