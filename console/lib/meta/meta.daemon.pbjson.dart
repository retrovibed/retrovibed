//
//  Generated code. Do not modify.
//  source: meta.daemon.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use daemonDescriptor instead')
const Daemon$json = {
  '1': 'Daemon',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'created_at', '3': 2, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 3, '4': 1, '5': 9, '10': 'updated_at'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {'1': 'hostname', '3': 5, '4': 1, '5': 9, '10': 'hostname'},
  ],
};

/// Descriptor for `Daemon`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List daemonDescriptor = $convert.base64Decode(
    'CgZEYWVtb24SDgoCaWQYASABKAlSAmlkEh4KCmNyZWF0ZWRfYXQYAiABKAlSCmNyZWF0ZWRfYX'
    'QSHgoKdXBkYXRlZF9hdBgDIAEoCVIKdXBkYXRlZF9hdBIgCgtkZXNjcmlwdGlvbhgEIAEoCVIL'
    'ZGVzY3JpcHRpb24SGgoIaG9zdG5hbWUYBSABKAlSCGhvc3RuYW1l');

@$core.Deprecated('Use daemonSearchRequestDescriptor instead')
const DaemonSearchRequest$json = {
  '1': 'DaemonSearchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'offset', '3': 2, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 3, '4': 1, '5': 4, '10': 'limit'},
  ],
};

/// Descriptor for `DaemonSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List daemonSearchRequestDescriptor = $convert.base64Decode(
    'ChNEYWVtb25TZWFyY2hSZXF1ZXN0EhQKBXF1ZXJ5GAEgASgJUgVxdWVyeRIWCgZvZmZzZXQYAi'
    'ABKARSBm9mZnNldBIUCgVsaW1pdBgDIAEoBFIFbGltaXQ=');

@$core.Deprecated('Use daemonSearchResponseDescriptor instead')
const DaemonSearchResponse$json = {
  '1': 'DaemonSearchResponse',
  '2': [
    {'1': 'next', '3': 1, '4': 1, '5': 11, '6': '.meta.DaemonSearchRequest', '10': 'next'},
    {'1': 'items', '3': 2, '4': 3, '5': 11, '6': '.meta.Daemon', '10': 'items'},
  ],
};

/// Descriptor for `DaemonSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List daemonSearchResponseDescriptor = $convert.base64Decode(
    'ChREYWVtb25TZWFyY2hSZXNwb25zZRItCgRuZXh0GAEgASgLMhkubWV0YS5EYWVtb25TZWFyY2'
    'hSZXF1ZXN0UgRuZXh0EiIKBWl0ZW1zGAIgAygLMgwubWV0YS5EYWVtb25SBWl0ZW1z');

@$core.Deprecated('Use daemonCreateRequestDescriptor instead')
const DaemonCreateRequest$json = {
  '1': 'DaemonCreateRequest',
  '2': [
    {'1': 'daemon', '3': 1, '4': 1, '5': 11, '6': '.meta.Daemon', '10': 'daemon'},
  ],
};

/// Descriptor for `DaemonCreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List daemonCreateRequestDescriptor = $convert.base64Decode(
    'ChNEYWVtb25DcmVhdGVSZXF1ZXN0EiQKBmRhZW1vbhgBIAEoCzIMLm1ldGEuRGFlbW9uUgZkYW'
    'Vtb24=');

@$core.Deprecated('Use daemonCreateResponseDescriptor instead')
const DaemonCreateResponse$json = {
  '1': 'DaemonCreateResponse',
  '2': [
    {'1': 'daemon', '3': 1, '4': 1, '5': 11, '6': '.meta.Daemon', '10': 'daemon'},
  ],
};

/// Descriptor for `DaemonCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List daemonCreateResponseDescriptor = $convert.base64Decode(
    'ChREYWVtb25DcmVhdGVSZXNwb25zZRIkCgZkYWVtb24YASABKAsyDC5tZXRhLkRhZW1vblIGZG'
    'FlbW9u');

@$core.Deprecated('Use daemonLookupRequestDescriptor instead')
const DaemonLookupRequest$json = {
  '1': 'DaemonLookupRequest',
};

/// Descriptor for `DaemonLookupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List daemonLookupRequestDescriptor = $convert.base64Decode(
    'ChNEYWVtb25Mb29rdXBSZXF1ZXN0');

@$core.Deprecated('Use daemonLookupResponseDescriptor instead')
const DaemonLookupResponse$json = {
  '1': 'DaemonLookupResponse',
  '2': [
    {'1': 'daemon', '3': 1, '4': 1, '5': 11, '6': '.meta.Daemon', '10': 'daemon'},
  ],
};

/// Descriptor for `DaemonLookupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List daemonLookupResponseDescriptor = $convert.base64Decode(
    'ChREYWVtb25Mb29rdXBSZXNwb25zZRIkCgZkYWVtb24YASABKAsyDC5tZXRhLkRhZW1vblIGZG'
    'FlbW9u');

@$core.Deprecated('Use daemonUpdateRequestDescriptor instead')
const DaemonUpdateRequest$json = {
  '1': 'DaemonUpdateRequest',
  '2': [
    {'1': 'daemon', '3': 1, '4': 1, '5': 11, '6': '.meta.Daemon', '10': 'daemon'},
  ],
};

/// Descriptor for `DaemonUpdateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List daemonUpdateRequestDescriptor = $convert.base64Decode(
    'ChNEYWVtb25VcGRhdGVSZXF1ZXN0EiQKBmRhZW1vbhgBIAEoCzIMLm1ldGEuRGFlbW9uUgZkYW'
    'Vtb24=');

@$core.Deprecated('Use daemonUpdateResponseDescriptor instead')
const DaemonUpdateResponse$json = {
  '1': 'DaemonUpdateResponse',
  '2': [
    {'1': 'daemon', '3': 1, '4': 1, '5': 11, '6': '.meta.Daemon', '10': 'daemon'},
  ],
};

/// Descriptor for `DaemonUpdateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List daemonUpdateResponseDescriptor = $convert.base64Decode(
    'ChREYWVtb25VcGRhdGVSZXNwb25zZRIkCgZkYWVtb24YASABKAsyDC5tZXRhLkRhZW1vblIGZG'
    'FlbW9u');

@$core.Deprecated('Use daemonDisableRequestDescriptor instead')
const DaemonDisableRequest$json = {
  '1': 'DaemonDisableRequest',
  '2': [
    {'1': 'daemon', '3': 1, '4': 1, '5': 11, '6': '.meta.Daemon', '10': 'daemon'},
  ],
};

/// Descriptor for `DaemonDisableRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List daemonDisableRequestDescriptor = $convert.base64Decode(
    'ChREYWVtb25EaXNhYmxlUmVxdWVzdBIkCgZkYWVtb24YASABKAsyDC5tZXRhLkRhZW1vblIGZG'
    'FlbW9u');

@$core.Deprecated('Use daemonDisableResponseDescriptor instead')
const DaemonDisableResponse$json = {
  '1': 'DaemonDisableResponse',
  '2': [
    {'1': 'daemon', '3': 1, '4': 1, '5': 11, '6': '.meta.Daemon', '10': 'daemon'},
  ],
};

/// Descriptor for `DaemonDisableResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List daemonDisableResponseDescriptor = $convert.base64Decode(
    'ChVEYWVtb25EaXNhYmxlUmVzcG9uc2USJAoGZGFlbW9uGAEgASgLMgwubWV0YS5EYWVtb25SBm'
    'RhZW1vbg==');

