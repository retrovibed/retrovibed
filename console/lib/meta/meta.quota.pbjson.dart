// This is a generated file - do not edit.
//
// Generated from meta.quota.proto.

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

@$core.Deprecated('Use quotaDescriptor instead')
const Quota$json = {
  '1': 'Quota',
  '2': [
    {'1': 'sku', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'account_id', '3': 2, '4': 1, '5': 9, '10': 'account_id'},
    {'1': 'created_at', '3': 3, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 4, '4': 1, '5': 9, '10': 'updated_at'},
    {'1': 'renewed_at', '3': 5, '4': 1, '5': 9, '10': 'renewed_at'},
    {'1': 'description', '3': 6, '4': 1, '5': 9, '10': 'description'},
    {'1': 'adjustable', '3': 7, '4': 1, '5': 8, '10': 'adjustable'},
    {'1': 'maximum', '3': 8, '4': 1, '5': 3, '10': 'maximum'},
    {'1': 'credits', '3': 9, '4': 1, '5': 3, '10': 'credits'},
    {'1': 'reserved', '3': 10, '4': 1, '5': 3, '10': 'reserved'},
    {'1': 'consumed', '3': 11, '4': 1, '5': 3, '10': 'consumed'},
  ],
};

/// Descriptor for `Quota`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List quotaDescriptor = $convert.base64Decode(
    'CgVRdW90YRIPCgNza3UYASABKAlSAmlkEh4KCmFjY291bnRfaWQYAiABKAlSCmFjY291bnRfaW'
    'QSHgoKY3JlYXRlZF9hdBgDIAEoCVIKY3JlYXRlZF9hdBIeCgp1cGRhdGVkX2F0GAQgASgJUgp1'
    'cGRhdGVkX2F0Eh4KCnJlbmV3ZWRfYXQYBSABKAlSCnJlbmV3ZWRfYXQSIAoLZGVzY3JpcHRpb2'
    '4YBiABKAlSC2Rlc2NyaXB0aW9uEh4KCmFkanVzdGFibGUYByABKAhSCmFkanVzdGFibGUSGAoH'
    'bWF4aW11bRgIIAEoA1IHbWF4aW11bRIYCgdjcmVkaXRzGAkgASgDUgdjcmVkaXRzEhoKCHJlc2'
    'VydmVkGAogASgDUghyZXNlcnZlZBIaCghjb25zdW1lZBgLIAEoA1IIY29uc3VtZWQ=');

@$core.Deprecated('Use quotaSearchRequestDescriptor instead')
const QuotaSearchRequest$json = {
  '1': 'QuotaSearchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'offset', '3': 2, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 3, '4': 1, '5': 4, '10': 'limit'},
  ],
};

/// Descriptor for `QuotaSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List quotaSearchRequestDescriptor = $convert.base64Decode(
    'ChJRdW90YVNlYXJjaFJlcXVlc3QSFAoFcXVlcnkYASABKAlSBXF1ZXJ5EhYKBm9mZnNldBgCIA'
    'EoBFIGb2Zmc2V0EhQKBWxpbWl0GAMgASgEUgVsaW1pdA==');

@$core.Deprecated('Use quotaSearchResponseDescriptor instead')
const QuotaSearchResponse$json = {
  '1': 'QuotaSearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.QuotaSearchRequest',
      '10': 'next'
    },
    {'1': 'items', '3': 2, '4': 3, '5': 11, '6': '.meta.Quota', '10': 'items'},
  ],
};

/// Descriptor for `QuotaSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List quotaSearchResponseDescriptor = $convert.base64Decode(
    'ChNRdW90YVNlYXJjaFJlc3BvbnNlEiwKBG5leHQYASABKAsyGC5tZXRhLlF1b3RhU2VhcmNoUm'
    'VxdWVzdFIEbmV4dBIhCgVpdGVtcxgCIAMoCzILLm1ldGEuUXVvdGFSBWl0ZW1z');

@$core.Deprecated('Use quotaUpdateRequestDescriptor instead')
const QuotaUpdateRequest$json = {
  '1': 'QuotaUpdateRequest',
  '2': [
    {'1': 'quota', '3': 1, '4': 1, '5': 11, '6': '.meta.Quota', '10': 'quota'},
  ],
};

/// Descriptor for `QuotaUpdateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List quotaUpdateRequestDescriptor = $convert.base64Decode(
    'ChJRdW90YVVwZGF0ZVJlcXVlc3QSIQoFcXVvdGEYASABKAsyCy5tZXRhLlF1b3RhUgVxdW90YQ'
    '==');

@$core.Deprecated('Use quotaUpdateResponseDescriptor instead')
const QuotaUpdateResponse$json = {
  '1': 'QuotaUpdateResponse',
  '2': [
    {'1': 'quota', '3': 1, '4': 1, '5': 11, '6': '.meta.Quota', '10': 'quota'},
  ],
};

/// Descriptor for `QuotaUpdateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List quotaUpdateResponseDescriptor = $convert.base64Decode(
    'ChNRdW90YVVwZGF0ZVJlc3BvbnNlEiEKBXF1b3RhGAEgASgLMgsubWV0YS5RdW90YVIFcXVvdG'
    'E=');

@$core.Deprecated('Use quotaFindRequestDescriptor instead')
const QuotaFindRequest$json = {
  '1': 'QuotaFindRequest',
};

/// Descriptor for `QuotaFindRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List quotaFindRequestDescriptor =
    $convert.base64Decode('ChBRdW90YUZpbmRSZXF1ZXN0');

@$core.Deprecated('Use quotaFindResponseDescriptor instead')
const QuotaFindResponse$json = {
  '1': 'QuotaFindResponse',
  '2': [
    {'1': 'quota', '3': 1, '4': 1, '5': 11, '6': '.meta.Quota', '10': 'quota'},
  ],
};

/// Descriptor for `QuotaFindResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List quotaFindResponseDescriptor = $convert.base64Decode(
    'ChFRdW90YUZpbmRSZXNwb25zZRIhCgVxdW90YRgBIAEoCzILLm1ldGEuUXVvdGFSBXF1b3Rh');
