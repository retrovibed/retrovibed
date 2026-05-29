// This is a generated file - do not edit.
//
// Generated from meta.wireguard.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use wireguardDescriptor instead')
const Wireguard$json = {
  '1': 'Wireguard',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'created_at', '3': 2, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 3, '4': 1, '5': 9, '10': 'updated_at'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {'1': 'default', '3': 5, '4': 1, '5': 8, '10': 'default'},
    {'1': 'port', '3': 6, '4': 1, '5': 13, '10': 'port'},
    {'1': 'dns_rate_limit', '3': 7, '4': 1, '5': 13, '10': 'dns_rate_limit'},
    {
      '1': 'maximum_connections',
      '3': 8,
      '4': 1,
      '5': 4,
      '10': 'maximum_connections'
    },
  ],
};

/// Descriptor for `Wireguard`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardDescriptor = $convert.base64Decode(
    'CglXaXJlZ3VhcmQSDgoCaWQYASABKAlSAmlkEh4KCmNyZWF0ZWRfYXQYAiABKAlSCmNyZWF0ZW'
    'RfYXQSHgoKdXBkYXRlZF9hdBgDIAEoCVIKdXBkYXRlZF9hdBIgCgtkZXNjcmlwdGlvbhgEIAEo'
    'CVILZGVzY3JpcHRpb24SGAoHZGVmYXVsdBgFIAEoCFIHZGVmYXVsdBISCgRwb3J0GAYgASgNUg'
    'Rwb3J0EiYKDmRuc19yYXRlX2xpbWl0GAcgASgNUg5kbnNfcmF0ZV9saW1pdBIwChNtYXhpbXVt'
    'X2Nvbm5lY3Rpb25zGAggASgEUhNtYXhpbXVtX2Nvbm5lY3Rpb25z');

@$core.Deprecated('Use wireguardSearchRequestDescriptor instead')
const WireguardSearchRequest$json = {
  '1': 'WireguardSearchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'offset', '3': 2, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 3, '4': 1, '5': 4, '10': 'limit'},
  ],
};

/// Descriptor for `WireguardSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardSearchRequestDescriptor =
    $convert.base64Decode(
        'ChZXaXJlZ3VhcmRTZWFyY2hSZXF1ZXN0EhQKBXF1ZXJ5GAEgASgJUgVxdWVyeRIWCgZvZmZzZX'
        'QYAiABKARSBm9mZnNldBIUCgVsaW1pdBgDIAEoBFIFbGltaXQ=');

@$core.Deprecated('Use wireguardSearchResponseDescriptor instead')
const WireguardSearchResponse$json = {
  '1': 'WireguardSearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.WireguardSearchRequest',
      '10': 'next'
    },
    {
      '1': 'items',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.meta.Wireguard',
      '10': 'items'
    },
  ],
};

/// Descriptor for `WireguardSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardSearchResponseDescriptor = $convert.base64Decode(
    'ChdXaXJlZ3VhcmRTZWFyY2hSZXNwb25zZRIwCgRuZXh0GAEgASgLMhwubWV0YS5XaXJlZ3Vhcm'
    'RTZWFyY2hSZXF1ZXN0UgRuZXh0EiUKBWl0ZW1zGAIgAygLMg8ubWV0YS5XaXJlZ3VhcmRSBWl0'
    'ZW1z');

@$core.Deprecated('Use wireguardUpdateRequestDescriptor instead')
const WireguardUpdateRequest$json = {
  '1': 'WireguardUpdateRequest',
  '2': [
    {
      '1': 'Wireguard',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.Wireguard',
      '10': 'wireguard'
    },
  ],
};

/// Descriptor for `WireguardUpdateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardUpdateRequestDescriptor =
    $convert.base64Decode(
        'ChZXaXJlZ3VhcmRVcGRhdGVSZXF1ZXN0Ei0KCVdpcmVndWFyZBgBIAEoCzIPLm1ldGEuV2lyZW'
        'd1YXJkUgl3aXJlZ3VhcmQ=');

@$core.Deprecated('Use wireguardUpdateResponseDescriptor instead')
const WireguardUpdateResponse$json = {
  '1': 'WireguardUpdateResponse',
  '2': [
    {
      '1': 'Wireguard',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.Wireguard',
      '10': 'wireguard'
    },
  ],
};

/// Descriptor for `WireguardUpdateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardUpdateResponseDescriptor =
    $convert.base64Decode(
        'ChdXaXJlZ3VhcmRVcGRhdGVSZXNwb25zZRItCglXaXJlZ3VhcmQYASABKAsyDy5tZXRhLldpcm'
        'VndWFyZFIJd2lyZWd1YXJk');

@$core.Deprecated('Use wireguardTouchRequestDescriptor instead')
const WireguardTouchRequest$json = {
  '1': 'WireguardTouchRequest',
};

/// Descriptor for `WireguardTouchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardTouchRequestDescriptor =
    $convert.base64Decode('ChVXaXJlZ3VhcmRUb3VjaFJlcXVlc3Q=');

@$core.Deprecated('Use wireguardTouchResponseDescriptor instead')
const WireguardTouchResponse$json = {
  '1': 'WireguardTouchResponse',
  '2': [
    {
      '1': 'Wireguard',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.Wireguard',
      '10': 'wireguard'
    },
  ],
};

/// Descriptor for `WireguardTouchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardTouchResponseDescriptor =
    $convert.base64Decode(
        'ChZXaXJlZ3VhcmRUb3VjaFJlc3BvbnNlEi0KCVdpcmVndWFyZBgBIAEoCzIPLm1ldGEuV2lyZW'
        'd1YXJkUgl3aXJlZ3VhcmQ=');

@$core.Deprecated('Use wireguardUploadRequestDescriptor instead')
const WireguardUploadRequest$json = {
  '1': 'WireguardUploadRequest',
};

/// Descriptor for `WireguardUploadRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardUploadRequestDescriptor =
    $convert.base64Decode('ChZXaXJlZ3VhcmRVcGxvYWRSZXF1ZXN0');

@$core.Deprecated('Use wireguardUploadResponseDescriptor instead')
const WireguardUploadResponse$json = {
  '1': 'WireguardUploadResponse',
  '2': [
    {
      '1': 'Wireguard',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.Wireguard',
      '10': 'wireguard'
    },
  ],
};

/// Descriptor for `WireguardUploadResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardUploadResponseDescriptor =
    $convert.base64Decode(
        'ChdXaXJlZ3VhcmRVcGxvYWRSZXNwb25zZRItCglXaXJlZ3VhcmQYASABKAsyDy5tZXRhLldpcm'
        'VndWFyZFIJd2lyZWd1YXJk');

@$core.Deprecated('Use wireguardCurrentRequestDescriptor instead')
const WireguardCurrentRequest$json = {
  '1': 'WireguardCurrentRequest',
};

/// Descriptor for `WireguardCurrentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardCurrentRequestDescriptor =
    $convert.base64Decode('ChdXaXJlZ3VhcmRDdXJyZW50UmVxdWVzdA==');

@$core.Deprecated('Use wireguardCurrentResponseDescriptor instead')
const WireguardCurrentResponse$json = {
  '1': 'WireguardCurrentResponse',
  '2': [
    {
      '1': 'Wireguard',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.Wireguard',
      '10': 'wireguard'
    },
    {'1': 'online', '3': 2, '4': 1, '5': 8, '10': 'online'},
    {'1': 'ip', '3': 3, '4': 1, '5': 9, '10': 'ip'},
    {'1': 'ip4', '3': 4, '4': 1, '5': 9, '10': 'ip4'},
  ],
};

/// Descriptor for `WireguardCurrentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardCurrentResponseDescriptor = $convert.base64Decode(
    'ChhXaXJlZ3VhcmRDdXJyZW50UmVzcG9uc2USLQoJV2lyZWd1YXJkGAEgASgLMg8ubWV0YS5XaX'
    'JlZ3VhcmRSCXdpcmVndWFyZBIWCgZvbmxpbmUYAiABKAhSBm9ubGluZRIOCgJpcBgDIAEoCVIC'
    'aXASEAoDaXA0GAQgASgJUgNpcDQ=');

@$core.Deprecated('Use wireguardDeleteRequestDescriptor instead')
const WireguardDeleteRequest$json = {
  '1': 'WireguardDeleteRequest',
};

/// Descriptor for `WireguardDeleteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardDeleteRequestDescriptor =
    $convert.base64Decode('ChZXaXJlZ3VhcmREZWxldGVSZXF1ZXN0');

@$core.Deprecated('Use wireguardDeleteResponseDescriptor instead')
const WireguardDeleteResponse$json = {
  '1': 'WireguardDeleteResponse',
  '2': [
    {
      '1': 'Wireguard',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.Wireguard',
      '10': 'wireguard'
    },
  ],
};

/// Descriptor for `WireguardDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardDeleteResponseDescriptor =
    $convert.base64Decode(
        'ChdXaXJlZ3VhcmREZWxldGVSZXNwb25zZRItCglXaXJlZ3VhcmQYASABKAsyDy5tZXRhLldpcm'
        'VndWFyZFIJd2lyZWd1YXJk');
