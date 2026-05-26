// This is a generated file - do not edit.
//
// Generated from media.recent.proto.

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

@$core.Deprecated('Use recentSearchRequestDescriptor instead')
const RecentSearchRequest$json = {
  '1': 'RecentSearchRequest',
  '2': [
    {
      '1': 'created',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.DateRange',
      '10': 'created'
    },
    {'1': 'mimetype', '3': 2, '4': 1, '5': 9, '10': 'mimetype'},
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 3, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `RecentSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recentSearchRequestDescriptor = $convert.base64Decode(
    'ChNSZWNlbnRTZWFyY2hSZXF1ZXN0EikKB2NyZWF0ZWQYASABKAsyDy5tZXRhLkRhdGVSYW5nZV'
    'IHY3JlYXRlZBIaCghtaW1ldHlwZRgCIAEoCVIIbWltZXR5cGUSFwoGb2Zmc2V0GIQHIAEoBFIG'
    'b2Zmc2V0EhUKBWxpbWl0GIUHIAEoBFIFbGltaXRKBQgDEIQHSgYIhgcQ6Ac=');

@$core.Deprecated('Use recentSearchResponseDescriptor instead')
const RecentSearchResponse$json = {
  '1': 'RecentSearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.RecentSearchRequest',
      '10': 'next'
    },
    {
      '1': 'items',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.media.RecentRecordRequest',
      '10': 'items'
    },
  ],
};

/// Descriptor for `RecentSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recentSearchResponseDescriptor = $convert.base64Decode(
    'ChRSZWNlbnRTZWFyY2hSZXNwb25zZRIuCgRuZXh0GAEgASgLMhoubWVkaWEuUmVjZW50U2Vhcm'
    'NoUmVxdWVzdFIEbmV4dBIwCgVpdGVtcxgCIAMoCzIaLm1lZGlhLlJlY2VudFJlY29yZFJlcXVl'
    'c3RSBWl0ZW1z');

@$core.Deprecated('Use recentRecordRequestDescriptor instead')
const RecentRecordRequest$json = {
  '1': 'RecentRecordRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'media', '3': 2, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
    {'1': 'duration', '3': 3, '4': 1, '5': 4, '10': 'duration'},
    {'1': 'position', '3': 4, '4': 1, '5': 4, '10': 'position'},
    {
      '1': 'query',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.media.MediaSearchRequest',
      '10': 'query'
    },
  ],
};

/// Descriptor for `RecentRecordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recentRecordRequestDescriptor = $convert.base64Decode(
    'ChNSZWNlbnRSZWNvcmRSZXF1ZXN0Eg4KAmlkGAEgASgJUgJpZBIiCgVtZWRpYRgCIAEoCzIMLm'
    '1lZGlhLk1lZGlhUgVtZWRpYRIaCghkdXJhdGlvbhgDIAEoBFIIZHVyYXRpb24SGgoIcG9zaXRp'
    'b24YBCABKARSCHBvc2l0aW9uEi8KBXF1ZXJ5GAUgASgLMhkubWVkaWEuTWVkaWFTZWFyY2hSZX'
    'F1ZXN0UgVxdWVyeQ==');

@$core.Deprecated('Use recentRecordResponseDescriptor instead')
const RecentRecordResponse$json = {
  '1': 'RecentRecordResponse',
};

/// Descriptor for `RecentRecordResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recentRecordResponseDescriptor =
    $convert.base64Decode('ChRSZWNlbnRSZWNvcmRSZXNwb25zZQ==');

@$core.Deprecated('Use recentDeleteRequestDescriptor instead')
const RecentDeleteRequest$json = {
  '1': 'RecentDeleteRequest',
};

/// Descriptor for `RecentDeleteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recentDeleteRequestDescriptor =
    $convert.base64Decode('ChNSZWNlbnREZWxldGVSZXF1ZXN0');

@$core.Deprecated('Use recentDeleteResponseDescriptor instead')
const RecentDeleteResponse$json = {
  '1': 'RecentDeleteResponse',
};

/// Descriptor for `RecentDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recentDeleteResponseDescriptor =
    $convert.base64Decode('ChRSZWNlbnREZWxldGVSZXNwb25zZQ==');
