// This is a generated file - do not edit.
//
// Generated from content.addressable.storage.proto.

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

@$core.Deprecated('Use mediaDescriptor instead')
const Media$json = {
  '1': 'Media',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'account_id', '3': 2, '4': 1, '5': 9, '10': 'account_id'},
    {'1': 'uploaded_by', '3': 3, '4': 1, '5': 9, '10': 'uploaded_by'},
    {'1': 'created_at', '3': 4, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 5, '4': 1, '5': 9, '10': 'updated_at'},
    {'1': 'mimetype', '3': 6, '4': 1, '5': 9, '10': 'mimetype'},
    {'1': 'bytes', '3': 7, '4': 1, '5': 4, '10': 'bytes'},
    {'1': 'usage', '3': 8, '4': 1, '5': 4, '10': 'usage'},
    {'1': 'tombstoned_at', '3': 9, '4': 1, '5': 9, '10': 'tombstoned_at'},
  ],
  '9': [
    {'1': 20, '2': 1000},
  ],
};

/// Descriptor for `Media`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaDescriptor = $convert.base64Decode(
    'CgVNZWRpYRIOCgJpZBgBIAEoCVICaWQSHgoKYWNjb3VudF9pZBgCIAEoCVIKYWNjb3VudF9pZB'
    'IgCgt1cGxvYWRlZF9ieRgDIAEoCVILdXBsb2FkZWRfYnkSHgoKY3JlYXRlZF9hdBgEIAEoCVIK'
    'Y3JlYXRlZF9hdBIeCgp1cGRhdGVkX2F0GAUgASgJUgp1cGRhdGVkX2F0EhoKCG1pbWV0eXBlGA'
    'YgASgJUghtaW1ldHlwZRIUCgVieXRlcxgHIAEoBFIFYnl0ZXMSFAoFdXNhZ2UYCCABKARSBXVz'
    'YWdlEiQKDXRvbWJzdG9uZWRfYXQYCSABKAlSDXRvbWJzdG9uZWRfYXRKBQgUEOgH');

@$core.Deprecated('Use mediaSearchRequestDescriptor instead')
const MediaSearchRequest$json = {
  '1': 'MediaSearchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'offset', '3': 2, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 3, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 4, '2': 1000},
  ],
};

/// Descriptor for `MediaSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaSearchRequestDescriptor = $convert.base64Decode(
    'ChJNZWRpYVNlYXJjaFJlcXVlc3QSFAoFcXVlcnkYASABKAlSBXF1ZXJ5EhYKBm9mZnNldBgCIA'
    'EoBFIGb2Zmc2V0EhQKBWxpbWl0GAMgASgEUgVsaW1pdEoFCAQQ6Ac=');

@$core.Deprecated('Use mediaSearchResponseDescriptor instead')
const MediaSearchResponse$json = {
  '1': 'MediaSearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.cas.MediaSearchRequest',
      '10': 'next'
    },
    {
      '1': 'items',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.retrovibed.cas.Media',
      '10': 'items'
    },
  ],
};

/// Descriptor for `MediaSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaSearchResponseDescriptor = $convert.base64Decode(
    'ChNNZWRpYVNlYXJjaFJlc3BvbnNlEjYKBG5leHQYASABKAsyIi5yZXRyb3ZpYmVkLmNhcy5NZW'
    'RpYVNlYXJjaFJlcXVlc3RSBG5leHQSKwoFaXRlbXMYAiADKAsyFS5yZXRyb3ZpYmVkLmNhcy5N'
    'ZWRpYVIFaXRlbXM=');

@$core.Deprecated('Use mediaCreateRequestDescriptor instead')
const MediaCreateRequest$json = {
  '1': 'MediaCreateRequest',
  '2': [
    {
      '1': 'media',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.cas.Media',
      '10': 'media'
    },
  ],
};

/// Descriptor for `MediaCreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaCreateRequestDescriptor = $convert.base64Decode(
    'ChJNZWRpYUNyZWF0ZVJlcXVlc3QSKwoFbWVkaWEYASABKAsyFS5yZXRyb3ZpYmVkLmNhcy5NZW'
    'RpYVIFbWVkaWE=');

@$core.Deprecated('Use mediaCreateResponseDescriptor instead')
const MediaCreateResponse$json = {
  '1': 'MediaCreateResponse',
  '2': [
    {
      '1': 'media',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.cas.Media',
      '10': 'media'
    },
  ],
};

/// Descriptor for `MediaCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaCreateResponseDescriptor = $convert.base64Decode(
    'ChNNZWRpYUNyZWF0ZVJlc3BvbnNlEisKBW1lZGlhGAEgASgLMhUucmV0cm92aWJlZC5jYXMuTW'
    'VkaWFSBW1lZGlh');

@$core.Deprecated('Use mediaUploadRequestDescriptor instead')
const MediaUploadRequest$json = {
  '1': 'MediaUploadRequest',
  '2': [
    {
      '1': 'media',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.cas.Media',
      '10': 'media'
    },
  ],
};

/// Descriptor for `MediaUploadRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaUploadRequestDescriptor = $convert.base64Decode(
    'ChJNZWRpYVVwbG9hZFJlcXVlc3QSKwoFbWVkaWEYASABKAsyFS5yZXRyb3ZpYmVkLmNhcy5NZW'
    'RpYVIFbWVkaWE=');

@$core.Deprecated('Use mediaUploadResponseDescriptor instead')
const MediaUploadResponse$json = {
  '1': 'MediaUploadResponse',
  '2': [
    {
      '1': 'media',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.cas.Media',
      '10': 'media'
    },
  ],
};

/// Descriptor for `MediaUploadResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaUploadResponseDescriptor = $convert.base64Decode(
    'ChNNZWRpYVVwbG9hZFJlc3BvbnNlEisKBW1lZGlhGAEgASgLMhUucmV0cm92aWJlZC5jYXMuTW'
    'VkaWFSBW1lZGlh');

@$core.Deprecated('Use mediaDownloadRequestDescriptor instead')
const MediaDownloadRequest$json = {
  '1': 'MediaDownloadRequest',
};

/// Descriptor for `MediaDownloadRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaDownloadRequestDescriptor =
    $convert.base64Decode('ChRNZWRpYURvd25sb2FkUmVxdWVzdA==');

@$core.Deprecated('Use mediaCompletedRequestDescriptor instead')
const MediaCompletedRequest$json = {
  '1': 'MediaCompletedRequest',
};

/// Descriptor for `MediaCompletedRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaCompletedRequestDescriptor =
    $convert.base64Decode('ChVNZWRpYUNvbXBsZXRlZFJlcXVlc3Q=');

@$core.Deprecated('Use mediaCompletedResponseDescriptor instead')
const MediaCompletedResponse$json = {
  '1': 'MediaCompletedResponse',
  '2': [
    {
      '1': 'media',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.cas.Media',
      '10': 'media'
    },
  ],
};

/// Descriptor for `MediaCompletedResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaCompletedResponseDescriptor =
    $convert.base64Decode(
        'ChZNZWRpYUNvbXBsZXRlZFJlc3BvbnNlEisKBW1lZGlhGAEgASgLMhUucmV0cm92aWJlZC5jYX'
        'MuTWVkaWFSBW1lZGlh');

@$core.Deprecated('Use mediaFindRequestDescriptor instead')
const MediaFindRequest$json = {
  '1': 'MediaFindRequest',
};

/// Descriptor for `MediaFindRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaFindRequestDescriptor =
    $convert.base64Decode('ChBNZWRpYUZpbmRSZXF1ZXN0');

@$core.Deprecated('Use mediaFindResponseDescriptor instead')
const MediaFindResponse$json = {
  '1': 'MediaFindResponse',
  '2': [
    {
      '1': 'media',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.cas.Media',
      '10': 'media'
    },
  ],
};

/// Descriptor for `MediaFindResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaFindResponseDescriptor = $convert.base64Decode(
    'ChFNZWRpYUZpbmRSZXNwb25zZRIrCgVtZWRpYRgBIAEoCzIVLnJldHJvdmliZWQuY2FzLk1lZG'
    'lhUgVtZWRpYQ==');

@$core.Deprecated('Use mediaDeleteRequestDescriptor instead')
const MediaDeleteRequest$json = {
  '1': 'MediaDeleteRequest',
};

/// Descriptor for `MediaDeleteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaDeleteRequestDescriptor =
    $convert.base64Decode('ChJNZWRpYURlbGV0ZVJlcXVlc3Q=');

@$core.Deprecated('Use mediaDeleteResponseDescriptor instead')
const MediaDeleteResponse$json = {
  '1': 'MediaDeleteResponse',
  '2': [
    {
      '1': 'media',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.cas.Media',
      '10': 'media'
    },
  ],
};

/// Descriptor for `MediaDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaDeleteResponseDescriptor = $convert.base64Decode(
    'ChNNZWRpYURlbGV0ZVJlc3BvbnNlEisKBW1lZGlhGAEgASgLMhUucmV0cm92aWJlZC5jYXMuTW'
    'VkaWFSBW1lZGlh');
