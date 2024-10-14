package giorpc

import (
	context "context"
	"encoding/base32"
	"encoding/base64"
	"fmt"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative example.proto

type Service struct {
	UnimplementedGiorpcServer
}

func (s *Service) Encode(ctx context.Context, request *EncodeRequest) (*EncodeResponse, error) {
	response := new(EncodeResponse)
	switch request.Encoding {
	case Encoding_ENCODING_BASE32_STANDARD:
		response.Encoded = base32.StdEncoding.EncodeToString([]byte(request.Decoded))
	case Encoding_ENCODING_BASE32_HEXIDECIMAL:
		response.Encoded = base32.HexEncoding.EncodeToString([]byte(request.Decoded))
	case Encoding_ENCODING_BASE64_STANDARD:
		response.Encoded = base64.StdEncoding.EncodeToString([]byte(request.Decoded))
	case Encoding_ENCODING_BASE64_URL_SAFE:
		response.Encoded = base64.URLEncoding.EncodeToString([]byte(request.Decoded))
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", request.Encoding)
	}
	return response, nil
}
func (s *Service) Decode(ctx context.Context, request *DecodeRequest) (*DecodeResponse, error) {
	response := new(DecodeResponse)
	switch request.Encoding {
	case Encoding_ENCODING_BASE32_STANDARD:
		decoded, err := base32.StdEncoding.DecodeString(request.Encoded)
		if err != nil {
			return nil, fmt.Errorf("failed to decode %q using %s", request.Encoded, request.Encoding)
		}
		response.Decoded = string(decoded)
	case Encoding_ENCODING_BASE32_HEXIDECIMAL:
		decoded, err := base32.HexEncoding.DecodeString(request.Encoded)
		if err != nil {
			return nil, fmt.Errorf("failed to decode %q using %s", request.Encoded, request.Encoding)
		}
		response.Decoded = string(decoded)
	case Encoding_ENCODING_BASE64_STANDARD:
		decoded, err := base64.StdEncoding.DecodeString(request.Encoded)
		if err != nil {
			return nil, fmt.Errorf("failed to decode %q using %s", request.Encoded, request.Encoding)
		}
		response.Decoded = string(decoded)
	case Encoding_ENCODING_BASE64_URL_SAFE:
		decoded, err := base64.URLEncoding.DecodeString(request.Encoded)
		if err != nil {
			return nil, fmt.Errorf("failed to decode %q using %s", request.Encoded, request.Encoding)
		}
		response.Decoded = string(decoded)
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", request.Encoding)
	}
	return response, nil
}
