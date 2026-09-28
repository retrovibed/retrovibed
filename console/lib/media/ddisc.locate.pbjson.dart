// This is a generated file - do not edit.
//
// Generated from media/ddisc.locate.proto.

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

@$core.Deprecated('Use locateDescriptor instead')
const Locate$json = {
  '1': 'Locate',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'created_at', '3': 2, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 3, '4': 1, '5': 9, '10': 'updated_at'},
    {'1': 'known_media_id', '3': 4, '4': 1, '5': 9, '10': 'known_media_id'},
    {
      '1': 'located_torrent_id',
      '3': 5,
      '4': 1,
      '5': 9,
      '10': 'located_torrent_id'
    },
    {'1': 'query', '3': 6, '4': 1, '5': 9, '10': 'query'},
    {'1': 'mimetype', '3': 7, '4': 1, '5': 9, '10': 'mimetype'},
    {'1': 'tombstoned_at', '3': 8, '4': 1, '5': 9, '10': 'tombstoned_at'},
    {'1': 'autodownload', '3': 9, '4': 1, '5': 8, '10': 'autodownload'},
    {'1': 'adult', '3': 10, '4': 1, '5': 8, '10': 'adult'},
    {'1': 'attempts', '3': 11, '4': 1, '5': 13, '10': 'attempts'},
    {'1': 'next_check_at', '3': 12, '4': 1, '5': 9, '10': 'next_check_at'},
  ],
};

/// Descriptor for `Locate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locateDescriptor = $convert.base64Decode(
    'CgZMb2NhdGUSDgoCaWQYASABKAlSAmlkEh4KCmNyZWF0ZWRfYXQYAiABKAlSCmNyZWF0ZWRfYX'
    'QSHgoKdXBkYXRlZF9hdBgDIAEoCVIKdXBkYXRlZF9hdBImCg5rbm93bl9tZWRpYV9pZBgEIAEo'
    'CVIOa25vd25fbWVkaWFfaWQSLgoSbG9jYXRlZF90b3JyZW50X2lkGAUgASgJUhJsb2NhdGVkX3'
    'RvcnJlbnRfaWQSFAoFcXVlcnkYBiABKAlSBXF1ZXJ5EhoKCG1pbWV0eXBlGAcgASgJUghtaW1l'
    'dHlwZRIkCg10b21ic3RvbmVkX2F0GAggASgJUg10b21ic3RvbmVkX2F0EiIKDGF1dG9kb3dubG'
    '9hZBgJIAEoCFIMYXV0b2Rvd25sb2FkEhQKBWFkdWx0GAogASgIUgVhZHVsdBIaCghhdHRlbXB0'
    'cxgLIAEoDVIIYXR0ZW1wdHMSJAoNbmV4dF9jaGVja19hdBgMIAEoCVINbmV4dF9jaGVja19hdA'
    '==');

@$core.Deprecated('Use locateSearchRequestDescriptor instead')
const LocateSearchRequest$json = {
  '1': 'LocateSearchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'id', '3': 2, '4': 3, '5': 9, '10': 'id'},
    {'1': 'attempts_min', '3': 3, '4': 1, '5': 4, '10': 'attempts_min'},
    {'1': 'attempts_max', '3': 4, '4': 1, '5': 4, '10': 'attempts_max'},
    {'1': 'pending', '3': 5, '4': 1, '5': 8, '10': 'pending'},
    {'1': 'completed', '3': 6, '4': 1, '5': 8, '10': 'completed'},
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 7, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `LocateSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locateSearchRequestDescriptor = $convert.base64Decode(
    'ChNMb2NhdGVTZWFyY2hSZXF1ZXN0EhQKBXF1ZXJ5GAEgASgJUgVxdWVyeRIOCgJpZBgCIAMoCV'
    'ICaWQSIgoMYXR0ZW1wdHNfbWluGAMgASgEUgxhdHRlbXB0c19taW4SIgoMYXR0ZW1wdHNfbWF4'
    'GAQgASgEUgxhdHRlbXB0c19tYXgSGAoHcGVuZGluZxgFIAEoCFIHcGVuZGluZxIcCgljb21wbG'
    'V0ZWQYBiABKAhSCWNvbXBsZXRlZBIXCgZvZmZzZXQYhAcgASgEUgZvZmZzZXQSFQoFbGltaXQY'
    'hQcgASgEUgVsaW1pdEoFCAcQhAdKBgiGBxDoBw==');

@$core.Deprecated('Use locateSearchResponseDescriptor instead')
const LocateSearchResponse$json = {
  '1': 'LocateSearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.ddisc.LocateSearchRequest',
      '10': 'next'
    },
    {
      '1': 'items',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.ddisc.Locate',
      '10': 'items'
    },
  ],
};

/// Descriptor for `LocateSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locateSearchResponseDescriptor = $convert.base64Decode(
    'ChRMb2NhdGVTZWFyY2hSZXNwb25zZRIuCgRuZXh0GAEgASgLMhouZGRpc2MuTG9jYXRlU2Vhcm'
    'NoUmVxdWVzdFIEbmV4dBIjCgVpdGVtcxgCIAMoCzINLmRkaXNjLkxvY2F0ZVIFaXRlbXM=');

@$core.Deprecated('Use locateLookupRequestDescriptor instead')
const LocateLookupRequest$json = {
  '1': 'LocateLookupRequest',
};

/// Descriptor for `LocateLookupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locateLookupRequestDescriptor =
    $convert.base64Decode('ChNMb2NhdGVMb29rdXBSZXF1ZXN0');

@$core.Deprecated('Use locateLookupResponseDescriptor instead')
const LocateLookupResponse$json = {
  '1': 'LocateLookupResponse',
  '2': [
    {
      '1': 'locate',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.ddisc.Locate',
      '10': 'locate'
    },
  ],
};

/// Descriptor for `LocateLookupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locateLookupResponseDescriptor = $convert.base64Decode(
    'ChRMb2NhdGVMb29rdXBSZXNwb25zZRIlCgZsb2NhdGUYASABKAsyDS5kZGlzYy5Mb2NhdGVSBm'
    'xvY2F0ZQ==');

@$core.Deprecated('Use locateCreateRequestDescriptor instead')
const LocateCreateRequest$json = {
  '1': 'LocateCreateRequest',
  '2': [
    {
      '1': 'locate',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.ddisc.Locate',
      '10': 'locate'
    },
  ],
};

/// Descriptor for `LocateCreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locateCreateRequestDescriptor = $convert.base64Decode(
    'ChNMb2NhdGVDcmVhdGVSZXF1ZXN0EiUKBmxvY2F0ZRgBIAEoCzINLmRkaXNjLkxvY2F0ZVIGbG'
    '9jYXRl');

@$core.Deprecated('Use locateCreateResponseDescriptor instead')
const LocateCreateResponse$json = {
  '1': 'LocateCreateResponse',
  '2': [
    {
      '1': 'locate',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.ddisc.Locate',
      '10': 'locate'
    },
  ],
};

/// Descriptor for `LocateCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locateCreateResponseDescriptor = $convert.base64Decode(
    'ChRMb2NhdGVDcmVhdGVSZXNwb25zZRIlCgZsb2NhdGUYASABKAsyDS5kZGlzYy5Mb2NhdGVSBm'
    'xvY2F0ZQ==');

@$core.Deprecated('Use locateDeleteResponseDescriptor instead')
const LocateDeleteResponse$json = {
  '1': 'LocateDeleteResponse',
  '2': [
    {
      '1': 'locate',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.ddisc.Locate',
      '10': 'locate'
    },
  ],
};

/// Descriptor for `LocateDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locateDeleteResponseDescriptor = $convert.base64Decode(
    'ChRMb2NhdGVEZWxldGVSZXNwb25zZRIlCgZsb2NhdGUYASABKAsyDS5kZGlzYy5Mb2NhdGVSBm'
    'xvY2F0ZQ==');

@$core.Deprecated('Use locateRetryRequestDescriptor instead')
const LocateRetryRequest$json = {
  '1': 'LocateRetryRequest',
  '2': [
    {'1': 'reset_attempts', '3': 1, '4': 1, '5': 8, '10': 'reset_attempts'},
  ],
};

/// Descriptor for `LocateRetryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locateRetryRequestDescriptor = $convert.base64Decode(
    'ChJMb2NhdGVSZXRyeVJlcXVlc3QSJgoOcmVzZXRfYXR0ZW1wdHMYASABKAhSDnJlc2V0X2F0dG'
    'VtcHRz');

@$core.Deprecated('Use locateRetryResponseDescriptor instead')
const LocateRetryResponse$json = {
  '1': 'LocateRetryResponse',
  '2': [
    {
      '1': 'locate',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.ddisc.Locate',
      '10': 'locate'
    },
  ],
};

/// Descriptor for `LocateRetryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locateRetryResponseDescriptor = $convert.base64Decode(
    'ChNMb2NhdGVSZXRyeVJlc3BvbnNlEiUKBmxvY2F0ZRgBIAEoCzINLmRkaXNjLkxvY2F0ZVIGbG'
    '9jYXRl');
