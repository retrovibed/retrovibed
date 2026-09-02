// This is a generated file - do not edit.
//
// Generated from media/media.filesystem.proto.

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

@$core.Deprecated('Use filesystemSearchRequestDescriptor instead')
const FilesystemSearchRequest$json = {
  '1': 'FilesystemSearchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'mimetypes', '3': 2, '4': 3, '5': 9, '10': 'mimetypes'},
    {'1': 'hidden', '3': 3, '4': 1, '5': 8, '10': 'hidden'},
    {'1': 'directory_id', '3': 4, '4': 1, '5': 9, '10': 'directory_id'},
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 5, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `FilesystemSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filesystemSearchRequestDescriptor = $convert.base64Decode(
    'ChdGaWxlc3lzdGVtU2VhcmNoUmVxdWVzdBIUCgVxdWVyeRgBIAEoCVIFcXVlcnkSHAoJbWltZX'
    'R5cGVzGAIgAygJUgltaW1ldHlwZXMSFgoGaGlkZGVuGAMgASgIUgZoaWRkZW4SIgoMZGlyZWN0'
    'b3J5X2lkGAQgASgJUgxkaXJlY3RvcnlfaWQSFwoGb2Zmc2V0GIQHIAEoBFIGb2Zmc2V0EhUKBW'
    'xpbWl0GIUHIAEoBFIFbGltaXRKBQgFEIQHSgYIhgcQ6Ac=');

@$core.Deprecated('Use filesystemSearchResponseDescriptor instead')
const FilesystemSearchResponse$json = {
  '1': 'FilesystemSearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.FilesystemSearchRequest',
      '10': 'next'
    },
    {'1': 'items', '3': 2, '4': 3, '5': 11, '6': '.media.Media', '10': 'items'},
    {
      '1': 'breadcrumb',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.media.Media',
      '10': 'breadcrumb'
    },
  ],
};

/// Descriptor for `FilesystemSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filesystemSearchResponseDescriptor = $convert.base64Decode(
    'ChhGaWxlc3lzdGVtU2VhcmNoUmVzcG9uc2USMgoEbmV4dBgBIAEoCzIeLm1lZGlhLkZpbGVzeX'
    'N0ZW1TZWFyY2hSZXF1ZXN0UgRuZXh0EiIKBWl0ZW1zGAIgAygLMgwubWVkaWEuTWVkaWFSBWl0'
    'ZW1zEiwKCmJyZWFkY3J1bWIYAyADKAsyDC5tZWRpYS5NZWRpYVIKYnJlYWRjcnVtYg==');

@$core.Deprecated('Use filesystemCreateRequestDescriptor instead')
const FilesystemCreateRequest$json = {
  '1': 'FilesystemCreateRequest',
  '2': [
    {'1': 'directory_id', '3': 1, '4': 1, '5': 9, '10': 'directory_id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `FilesystemCreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filesystemCreateRequestDescriptor =
    $convert.base64Decode(
        'ChdGaWxlc3lzdGVtQ3JlYXRlUmVxdWVzdBIiCgxkaXJlY3RvcnlfaWQYASABKAlSDGRpcmVjdG'
        '9yeV9pZBISCgRuYW1lGAIgASgJUgRuYW1l');

@$core.Deprecated('Use filesystemCreateResponseDescriptor instead')
const FilesystemCreateResponse$json = {
  '1': 'FilesystemCreateResponse',
  '2': [
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
  ],
};

/// Descriptor for `FilesystemCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filesystemCreateResponseDescriptor =
    $convert.base64Decode(
        'ChhGaWxlc3lzdGVtQ3JlYXRlUmVzcG9uc2USIgoFbWVkaWEYASABKAsyDC5tZWRpYS5NZWRpYV'
        'IFbWVkaWE=');

@$core.Deprecated('Use filesystemMoveRequestDescriptor instead')
const FilesystemMoveRequest$json = {
  '1': 'FilesystemMoveRequest',
  '2': [
    {'1': 'directory_id', '3': 1, '4': 1, '5': 9, '10': 'directory_id'},
  ],
};

/// Descriptor for `FilesystemMoveRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filesystemMoveRequestDescriptor = $convert.base64Decode(
    'ChVGaWxlc3lzdGVtTW92ZVJlcXVlc3QSIgoMZGlyZWN0b3J5X2lkGAEgASgJUgxkaXJlY3Rvcn'
    'lfaWQ=');

@$core.Deprecated('Use filesystemMoveResponseDescriptor instead')
const FilesystemMoveResponse$json = {
  '1': 'FilesystemMoveResponse',
  '2': [
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
  ],
};

/// Descriptor for `FilesystemMoveResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filesystemMoveResponseDescriptor =
    $convert.base64Decode(
        'ChZGaWxlc3lzdGVtTW92ZVJlc3BvbnNlEiIKBW1lZGlhGAEgASgLMgwubWVkaWEuTWVkaWFSBW'
        '1lZGlh');

@$core.Deprecated('Use filesystemDeleteResponseDescriptor instead')
const FilesystemDeleteResponse$json = {
  '1': 'FilesystemDeleteResponse',
  '2': [
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
    {
      '1': 'removed',
      '3': 2,
      '4': 1,
      '5': 4,
      '8': {'6': 1},
      '10': 'removed',
    },
  ],
};

/// Descriptor for `FilesystemDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filesystemDeleteResponseDescriptor =
    $convert.base64Decode(
        'ChhGaWxlc3lzdGVtRGVsZXRlUmVzcG9uc2USIgoFbWVkaWEYASABKAsyDC5tZWRpYS5NZWRpYV'
        'IFbWVkaWESHAoHcmVtb3ZlZBgCIAEoBEICMAFSB3JlbW92ZWQ=');
