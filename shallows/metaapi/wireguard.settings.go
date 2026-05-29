package metaapi

import "google.golang.org/protobuf/encoding/protojson"

func (t *Wireguard) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *Wireguard) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *WireguardSearchRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *WireguardSearchRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *WireguardSearchResponse) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *WireguardSearchResponse) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *WireguardUpdateRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *WireguardUpdateRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *WireguardUpdateResponse) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *WireguardUpdateResponse) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *WireguardTouchRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *WireguardTouchRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *WireguardTouchResponse) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *WireguardTouchResponse) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *WireguardUploadRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *WireguardUploadRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *WireguardUploadResponse) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *WireguardUploadResponse) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *WireguardCurrentRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *WireguardCurrentRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *WireguardCurrentResponse) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *WireguardCurrentResponse) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *WireguardDeleteRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *WireguardDeleteRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *WireguardDeleteResponse) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *WireguardDeleteResponse) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}
