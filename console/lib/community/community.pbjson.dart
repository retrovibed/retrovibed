// This is a generated file - do not edit.
//
// Generated from community/community.proto.

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

@$core.Deprecated('Use publishModeDescriptor instead')
const PublishMode$json = {
  '1': 'PublishMode',
  '2': [
    {'1': 'UNLISTED', '2': 0},
    {'1': 'LISTED', '2': 1},
    {'1': 'SYNDICATED', '2': 2},
  ],
};

/// Descriptor for `PublishMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List publishModeDescriptor = $convert.base64Decode(
    'CgtQdWJsaXNoTW9kZRIMCghVTkxJU1RFRBAAEgoKBkxJU1RFRBABEg4KClNZTkRJQ0FURUQQAg'
    '==');

@$core.Deprecated('Use communityDescriptor instead')
const Community$json = {
  '1': 'Community',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'account_id', '3': 2, '4': 1, '5': 9, '10': 'account_id'},
    {'1': 'created_at', '3': 4, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 5, '4': 1, '5': 9, '10': 'updated_at'},
    {'1': 'mimetype', '3': 6, '4': 1, '5': 9, '10': 'mimetype'},
    {'1': 'domain', '3': 7, '4': 1, '5': 9, '10': 'domain'},
    {'1': 'description', '3': 8, '4': 1, '5': 9, '10': 'description'},
    {'1': 'entropy', '3': 9, '4': 1, '5': 9, '10': 'entropy'},
    {'1': 'bytes', '3': 10, '4': 1, '5': 4, '10': 'bytes'},
    {'1': 'subscribed_at', '3': 11, '4': 1, '5': 9, '10': 'subscribed_at'},
    {
      '1': 'default_publish_mode',
      '3': 12,
      '4': 1,
      '5': 14,
      '6': '.retrovibed.community.PublishMode',
      '10': 'default_publish_mode'
    },
    {'1': 'hidden', '3': 13, '4': 1, '5': 8, '10': 'hidden'},
    {'1': 'url', '3': 14, '4': 1, '5': 9, '10': 'url'},
    {'1': 'adult', '3': 15, '4': 1, '5': 8, '10': 'adult'},
    {'1': 'default_ttl', '3': 16, '4': 1, '5': 4, '10': 'default_ttl'},
    {
      '1': 'default_language',
      '3': 17,
      '4': 1,
      '5': 9,
      '10': 'default_language'
    },
    {'1': 'last_sync_at', '3': 1000, '4': 1, '5': 9, '10': 'last_sync_at'},
  ],
  '9': [
    {'1': 18, '2': 1000},
  ],
};

/// Descriptor for `Community`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityDescriptor = $convert.base64Decode(
    'CglDb21tdW5pdHkSDgoCaWQYASABKAlSAmlkEh4KCmFjY291bnRfaWQYAiABKAlSCmFjY291bn'
    'RfaWQSHgoKY3JlYXRlZF9hdBgEIAEoCVIKY3JlYXRlZF9hdBIeCgp1cGRhdGVkX2F0GAUgASgJ'
    'Ugp1cGRhdGVkX2F0EhoKCG1pbWV0eXBlGAYgASgJUghtaW1ldHlwZRIWCgZkb21haW4YByABKA'
    'lSBmRvbWFpbhIgCgtkZXNjcmlwdGlvbhgIIAEoCVILZGVzY3JpcHRpb24SGAoHZW50cm9weRgJ'
    'IAEoCVIHZW50cm9weRIUCgVieXRlcxgKIAEoBFIFYnl0ZXMSJAoNc3Vic2NyaWJlZF9hdBgLIA'
    'EoCVINc3Vic2NyaWJlZF9hdBJVChRkZWZhdWx0X3B1Ymxpc2hfbW9kZRgMIAEoDjIhLnJldHJv'
    'dmliZWQuY29tbXVuaXR5LlB1Ymxpc2hNb2RlUhRkZWZhdWx0X3B1Ymxpc2hfbW9kZRIWCgZoaW'
    'RkZW4YDSABKAhSBmhpZGRlbhIQCgN1cmwYDiABKAlSA3VybBIUCgVhZHVsdBgPIAEoCFIFYWR1'
    'bHQSIAoLZGVmYXVsdF90dGwYECABKARSC2RlZmF1bHRfdHRsEioKEGRlZmF1bHRfbGFuZ3VhZ2'
    'UYESABKAlSEGRlZmF1bHRfbGFuZ3VhZ2USIwoMbGFzdF9zeW5jX2F0GOgHIAEoCVIMbGFzdF9z'
    'eW5jX2F0SgUIEhDoBw==');

@$core.Deprecated('Use communitySearchRequestDescriptor instead')
const CommunitySearchRequest$json = {
  '1': 'CommunitySearchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'offset', '3': 2, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 3, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 4, '2': 1000},
  ],
};

/// Descriptor for `CommunitySearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communitySearchRequestDescriptor =
    $convert.base64Decode(
        'ChZDb21tdW5pdHlTZWFyY2hSZXF1ZXN0EhQKBXF1ZXJ5GAEgASgJUgVxdWVyeRIWCgZvZmZzZX'
        'QYAiABKARSBm9mZnNldBIUCgVsaW1pdBgDIAEoBFIFbGltaXRKBQgEEOgH');

@$core.Deprecated('Use communitySearchResponseDescriptor instead')
const CommunitySearchResponse$json = {
  '1': 'CommunitySearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.CommunitySearchRequest',
      '10': 'next'
    },
    {
      '1': 'items',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.retrovibed.community.Community',
      '10': 'items'
    },
  ],
};

/// Descriptor for `CommunitySearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communitySearchResponseDescriptor = $convert.base64Decode(
    'ChdDb21tdW5pdHlTZWFyY2hSZXNwb25zZRJACgRuZXh0GAEgASgLMiwucmV0cm92aWJlZC5jb2'
    '1tdW5pdHkuQ29tbXVuaXR5U2VhcmNoUmVxdWVzdFIEbmV4dBI1CgVpdGVtcxgCIAMoCzIfLnJl'
    'dHJvdmliZWQuY29tbXVuaXR5LkNvbW11bml0eVIFaXRlbXM=');

@$core.Deprecated('Use communityCreateRequestDescriptor instead')
const CommunityCreateRequest$json = {
  '1': 'CommunityCreateRequest',
  '2': [
    {
      '1': 'community',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.Community',
      '10': 'community'
    },
  ],
};

/// Descriptor for `CommunityCreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityCreateRequestDescriptor =
    $convert.base64Decode(
        'ChZDb21tdW5pdHlDcmVhdGVSZXF1ZXN0Ej0KCWNvbW11bml0eRgBIAEoCzIfLnJldHJvdmliZW'
        'QuY29tbXVuaXR5LkNvbW11bml0eVIJY29tbXVuaXR5');

@$core.Deprecated('Use communityCreateResponseDescriptor instead')
const CommunityCreateResponse$json = {
  '1': 'CommunityCreateResponse',
  '2': [
    {
      '1': 'community',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.Community',
      '10': 'community'
    },
  ],
};

/// Descriptor for `CommunityCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityCreateResponseDescriptor =
    $convert.base64Decode(
        'ChdDb21tdW5pdHlDcmVhdGVSZXNwb25zZRI9Cgljb21tdW5pdHkYASABKAsyHy5yZXRyb3ZpYm'
        'VkLmNvbW11bml0eS5Db21tdW5pdHlSCWNvbW11bml0eQ==');

@$core.Deprecated('Use communityFindRequestDescriptor instead')
const CommunityFindRequest$json = {
  '1': 'CommunityFindRequest',
};

/// Descriptor for `CommunityFindRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityFindRequestDescriptor =
    $convert.base64Decode('ChRDb21tdW5pdHlGaW5kUmVxdWVzdA==');

@$core.Deprecated('Use communityFindResponseDescriptor instead')
const CommunityFindResponse$json = {
  '1': 'CommunityFindResponse',
  '2': [
    {
      '1': 'community',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.Community',
      '10': 'community'
    },
  ],
};

/// Descriptor for `CommunityFindResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityFindResponseDescriptor = $convert.base64Decode(
    'ChVDb21tdW5pdHlGaW5kUmVzcG9uc2USPQoJY29tbXVuaXR5GAEgASgLMh8ucmV0cm92aWJlZC'
    '5jb21tdW5pdHkuQ29tbXVuaXR5Ugljb21tdW5pdHk=');

@$core.Deprecated('Use communityUploadRequestDescriptor instead')
const CommunityUploadRequest$json = {
  '1': 'CommunityUploadRequest',
  '2': [
    {
      '1': 'community',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.Community',
      '10': 'community'
    },
  ],
};

/// Descriptor for `CommunityUploadRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityUploadRequestDescriptor =
    $convert.base64Decode(
        'ChZDb21tdW5pdHlVcGxvYWRSZXF1ZXN0Ej0KCWNvbW11bml0eRgBIAEoCzIfLnJldHJvdmliZW'
        'QuY29tbXVuaXR5LkNvbW11bml0eVIJY29tbXVuaXR5');

@$core.Deprecated('Use communityUploadResponseDescriptor instead')
const CommunityUploadResponse$json = {
  '1': 'CommunityUploadResponse',
  '2': [
    {
      '1': 'community',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.Community',
      '10': 'community'
    },
  ],
};

/// Descriptor for `CommunityUploadResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityUploadResponseDescriptor =
    $convert.base64Decode(
        'ChdDb21tdW5pdHlVcGxvYWRSZXNwb25zZRI9Cgljb21tdW5pdHkYASABKAsyHy5yZXRyb3ZpYm'
        'VkLmNvbW11bml0eS5Db21tdW5pdHlSCWNvbW11bml0eQ==');

@$core.Deprecated('Use communityDeleteRequestDescriptor instead')
const CommunityDeleteRequest$json = {
  '1': 'CommunityDeleteRequest',
};

/// Descriptor for `CommunityDeleteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityDeleteRequestDescriptor =
    $convert.base64Decode('ChZDb21tdW5pdHlEZWxldGVSZXF1ZXN0');

@$core.Deprecated('Use communityDeleteResponseDescriptor instead')
const CommunityDeleteResponse$json = {
  '1': 'CommunityDeleteResponse',
  '2': [
    {
      '1': 'community',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.Community',
      '10': 'community'
    },
  ],
};

/// Descriptor for `CommunityDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityDeleteResponseDescriptor =
    $convert.base64Decode(
        'ChdDb21tdW5pdHlEZWxldGVSZXNwb25zZRI9Cgljb21tdW5pdHkYASABKAsyHy5yZXRyb3ZpYm'
        'VkLmNvbW11bml0eS5Db21tdW5pdHlSCWNvbW11bml0eQ==');

@$core.Deprecated('Use communityUpdateRequestDescriptor instead')
const CommunityUpdateRequest$json = {
  '1': 'CommunityUpdateRequest',
  '2': [
    {
      '1': 'community',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.Community',
      '10': 'community'
    },
  ],
};

/// Descriptor for `CommunityUpdateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityUpdateRequestDescriptor =
    $convert.base64Decode(
        'ChZDb21tdW5pdHlVcGRhdGVSZXF1ZXN0Ej0KCWNvbW11bml0eRgBIAEoCzIfLnJldHJvdmliZW'
        'QuY29tbXVuaXR5LkNvbW11bml0eVIJY29tbXVuaXR5');

@$core.Deprecated('Use communityUpdateResponseDescriptor instead')
const CommunityUpdateResponse$json = {
  '1': 'CommunityUpdateResponse',
  '2': [
    {
      '1': 'community',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.Community',
      '10': 'community'
    },
  ],
};

/// Descriptor for `CommunityUpdateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityUpdateResponseDescriptor =
    $convert.base64Decode(
        'ChdDb21tdW5pdHlVcGRhdGVSZXNwb25zZRI9Cgljb21tdW5pdHkYASABKAsyHy5yZXRyb3ZpYm'
        'VkLmNvbW11bml0eS5Db21tdW5pdHlSCWNvbW11bml0eQ==');

@$core.Deprecated('Use communitySubscribeRequestDescriptor instead')
const CommunitySubscribeRequest$json = {
  '1': 'CommunitySubscribeRequest',
};

/// Descriptor for `CommunitySubscribeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communitySubscribeRequestDescriptor =
    $convert.base64Decode('ChlDb21tdW5pdHlTdWJzY3JpYmVSZXF1ZXN0');

@$core.Deprecated('Use communitySubscribeResponseDescriptor instead')
const CommunitySubscribeResponse$json = {
  '1': 'CommunitySubscribeResponse',
};

/// Descriptor for `CommunitySubscribeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communitySubscribeResponseDescriptor =
    $convert.base64Decode('ChpDb21tdW5pdHlTdWJzY3JpYmVSZXNwb25zZQ==');

@$core.Deprecated('Use youTubeStatusDescriptor instead')
const YouTubeStatus$json = {
  '1': 'YouTubeStatus',
  '2': [
    {'1': 'linked', '3': 1, '4': 1, '5': 8, '10': 'linked'},
    {'1': 'id', '3': 2, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `YouTubeStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List youTubeStatusDescriptor = $convert.base64Decode(
    'Cg1Zb3VUdWJlU3RhdHVzEhYKBmxpbmtlZBgBIAEoCFIGbGlua2VkEg4KAmlkGAIgASgJUgJpZA'
    '==');
