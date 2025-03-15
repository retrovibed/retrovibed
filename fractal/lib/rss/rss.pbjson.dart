//
//  Generated code. Do not modify.
//  source: rss.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use feedDescriptor instead')
const Feed$json = {
  '1': 'Feed',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'created_at', '3': 2, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 3, '4': 1, '5': 9, '10': 'updated_at'},
    {'1': 'next_check', '3': 4, '4': 1, '5': 9, '10': 'next_check'},
    {'1': 'description', '3': 5, '4': 1, '5': 9, '10': 'description'},
    {'1': 'url', '3': 6, '4': 1, '5': 9, '10': 'url'},
    {'1': 'autodownload', '3': 7, '4': 1, '5': 8, '10': 'autodownload'},
    {'1': 'autoarchive', '3': 8, '4': 1, '5': 8, '10': 'autoarchive'},
    {'1': 'contributing', '3': 9, '4': 1, '5': 8, '10': 'contributing'},
  ],
};

/// Descriptor for `Feed`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List feedDescriptor = $convert.base64Decode(
    'CgRGZWVkEg4KAmlkGAEgASgJUgJpZBIeCgpjcmVhdGVkX2F0GAIgASgJUgpjcmVhdGVkX2F0Eh'
    '4KCnVwZGF0ZWRfYXQYAyABKAlSCnVwZGF0ZWRfYXQSHgoKbmV4dF9jaGVjaxgEIAEoCVIKbmV4'
    'dF9jaGVjaxIgCgtkZXNjcmlwdGlvbhgFIAEoCVILZGVzY3JpcHRpb24SEAoDdXJsGAYgASgJUg'
    'N1cmwSIgoMYXV0b2Rvd25sb2FkGAcgASgIUgxhdXRvZG93bmxvYWQSIAoLYXV0b2FyY2hpdmUY'
    'CCABKAhSC2F1dG9hcmNoaXZlEiIKDGNvbnRyaWJ1dGluZxgJIAEoCFIMY29udHJpYnV0aW5n');

@$core.Deprecated('Use feedSearchRequestDescriptor instead')
const FeedSearchRequest$json = {
  '1': 'FeedSearchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 2, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `FeedSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List feedSearchRequestDescriptor = $convert.base64Decode(
    'ChFGZWVkU2VhcmNoUmVxdWVzdBIUCgVxdWVyeRgBIAEoCVIFcXVlcnkSFwoGb2Zmc2V0GIQHIA'
    'EoBFIGb2Zmc2V0EhUKBWxpbWl0GIUHIAEoBFIFbGltaXRKBQgCEIQHSgYIhgcQ6Ac=');

@$core.Deprecated('Use feedSearchResponseDescriptor instead')
const FeedSearchResponse$json = {
  '1': 'FeedSearchResponse',
  '2': [
    {'1': 'next', '3': 1, '4': 1, '5': 11, '6': '.rss.FeedSearchRequest', '10': 'next'},
    {'1': 'items', '3': 2, '4': 3, '5': 11, '6': '.rss.Feed', '10': 'items'},
  ],
};

/// Descriptor for `FeedSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List feedSearchResponseDescriptor = $convert.base64Decode(
    'ChJGZWVkU2VhcmNoUmVzcG9uc2USKgoEbmV4dBgBIAEoCzIWLnJzcy5GZWVkU2VhcmNoUmVxdW'
    'VzdFIEbmV4dBIfCgVpdGVtcxgCIAMoCzIJLnJzcy5GZWVkUgVpdGVtcw==');

@$core.Deprecated('Use feedCreateRequestDescriptor instead')
const FeedCreateRequest$json = {
  '1': 'FeedCreateRequest',
  '2': [
    {'1': 'feed', '3': 1, '4': 1, '5': 11, '6': '.rss.Feed', '10': 'feed'},
  ],
};

/// Descriptor for `FeedCreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List feedCreateRequestDescriptor = $convert.base64Decode(
    'ChFGZWVkQ3JlYXRlUmVxdWVzdBIdCgRmZWVkGAEgASgLMgkucnNzLkZlZWRSBGZlZWQ=');

@$core.Deprecated('Use feedCreateResponseDescriptor instead')
const FeedCreateResponse$json = {
  '1': 'FeedCreateResponse',
  '2': [
    {'1': 'feed', '3': 1, '4': 1, '5': 11, '6': '.rss.Feed', '10': 'feed'},
  ],
};

/// Descriptor for `FeedCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List feedCreateResponseDescriptor = $convert.base64Decode(
    'ChJGZWVkQ3JlYXRlUmVzcG9uc2USHQoEZmVlZBgBIAEoCzIJLnJzcy5GZWVkUgRmZWVk');

@$core.Deprecated('Use feedUpdateRequestDescriptor instead')
const FeedUpdateRequest$json = {
  '1': 'FeedUpdateRequest',
  '2': [
    {'1': 'feed', '3': 1, '4': 1, '5': 11, '6': '.rss.Feed', '10': 'feed'},
  ],
};

/// Descriptor for `FeedUpdateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List feedUpdateRequestDescriptor = $convert.base64Decode(
    'ChFGZWVkVXBkYXRlUmVxdWVzdBIdCgRmZWVkGAEgASgLMgkucnNzLkZlZWRSBGZlZWQ=');

@$core.Deprecated('Use feedUpdateResponseDescriptor instead')
const FeedUpdateResponse$json = {
  '1': 'FeedUpdateResponse',
  '2': [
    {'1': 'feed', '3': 1, '4': 1, '5': 11, '6': '.rss.Feed', '10': 'feed'},
  ],
};

/// Descriptor for `FeedUpdateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List feedUpdateResponseDescriptor = $convert.base64Decode(
    'ChJGZWVkVXBkYXRlUmVzcG9uc2USHQoEZmVlZBgBIAEoCzIJLnJzcy5GZWVkUgRmZWVk');

@$core.Deprecated('Use feedDeleteRequestDescriptor instead')
const FeedDeleteRequest$json = {
  '1': 'FeedDeleteRequest',
};

/// Descriptor for `FeedDeleteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List feedDeleteRequestDescriptor = $convert.base64Decode(
    'ChFGZWVkRGVsZXRlUmVxdWVzdA==');

@$core.Deprecated('Use feedDeleteResponseDescriptor instead')
const FeedDeleteResponse$json = {
  '1': 'FeedDeleteResponse',
  '2': [
    {'1': 'feed', '3': 1, '4': 1, '5': 11, '6': '.rss.Feed', '10': 'feed'},
  ],
};

/// Descriptor for `FeedDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List feedDeleteResponseDescriptor = $convert.base64Decode(
    'ChJGZWVkRGVsZXRlUmVzcG9uc2USHQoEZmVlZBgBIAEoCzIJLnJzcy5GZWVkUgRmZWVk');

