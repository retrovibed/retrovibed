// This is a generated file - do not edit.
//
// Generated from community/community.social.proto.

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

@$core.Deprecated('Use pluginPublisherDescriptor instead')
const PluginPublisher$json = {
  '1': 'PluginPublisher',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'path', '3': 2, '4': 1, '5': 9, '10': 'path'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
    {'1': 'mimetype', '3': 4, '4': 1, '5': 9, '10': 'mimetype'},
    {'1': 'created_at', '3': 5, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 6, '4': 1, '5': 9, '10': 'updated_at'},
  ],
};

/// Descriptor for `PluginPublisher`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pluginPublisherDescriptor = $convert.base64Decode(
    'Cg9QbHVnaW5QdWJsaXNoZXISDgoCaWQYASABKAlSAmlkEhIKBHBhdGgYAiABKAlSBHBhdGgSIA'
    'oLZGVzY3JpcHRpb24YAyABKAlSC2Rlc2NyaXB0aW9uEhoKCG1pbWV0eXBlGAQgASgJUghtaW1l'
    'dHlwZRIeCgpjcmVhdGVkX2F0GAUgASgJUgpjcmVhdGVkX2F0Eh4KCnVwZGF0ZWRfYXQYBiABKA'
    'lSCnVwZGF0ZWRfYXQ=');

@$core.Deprecated('Use communityPublisherDescriptor instead')
const CommunityPublisher$json = {
  '1': 'CommunityPublisher',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'community_id', '3': 2, '4': 1, '5': 9, '10': 'community_id'},
    {'1': 'publisher_id', '3': 3, '4': 1, '5': 9, '10': 'publisher_id'},
    {'1': 'created_at', '3': 4, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 5, '4': 1, '5': 9, '10': 'updated_at'},
  ],
};

/// Descriptor for `CommunityPublisher`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityPublisherDescriptor = $convert.base64Decode(
    'ChJDb21tdW5pdHlQdWJsaXNoZXISDgoCaWQYASABKAlSAmlkEiIKDGNvbW11bml0eV9pZBgCIA'
    'EoCVIMY29tbXVuaXR5X2lkEiIKDHB1Ymxpc2hlcl9pZBgDIAEoCVIMcHVibGlzaGVyX2lkEh4K'
    'CmNyZWF0ZWRfYXQYBCABKAlSCmNyZWF0ZWRfYXQSHgoKdXBkYXRlZF9hdBgFIAEoCVIKdXBkYX'
    'RlZF9hdA==');

@$core.Deprecated('Use communitySocialDescriptor instead')
const CommunitySocial$json = {
  '1': 'CommunitySocial',
  '2': [
    {
      '1': 'community',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.Community',
      '10': 'community'
    },
    {
      '1': 'enabled',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.retrovibed.community.CommunityPublisher',
      '10': 'enabled'
    },
  ],
};

/// Descriptor for `CommunitySocial`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communitySocialDescriptor = $convert.base64Decode(
    'Cg9Db21tdW5pdHlTb2NpYWwSPQoJY29tbXVuaXR5GAEgASgLMh8ucmV0cm92aWJlZC5jb21tdW'
    '5pdHkuQ29tbXVuaXR5Ugljb21tdW5pdHkSQgoHZW5hYmxlZBgCIAMoCzIoLnJldHJvdmliZWQu'
    'Y29tbXVuaXR5LkNvbW11bml0eVB1Ymxpc2hlclIHZW5hYmxlZA==');

@$core.Deprecated('Use socialsSearchRequestDescriptor instead')
const SocialsSearchRequest$json = {
  '1': 'SocialsSearchRequest',
  '2': [
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 1, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `SocialsSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List socialsSearchRequestDescriptor = $convert.base64Decode(
    'ChRTb2NpYWxzU2VhcmNoUmVxdWVzdBIXCgZvZmZzZXQYhAcgASgEUgZvZmZzZXQSFQoFbGltaX'
    'QYhQcgASgEUgVsaW1pdEoFCAEQhAdKBgiGBxDoBw==');

@$core.Deprecated('Use socialsSearchResponseDescriptor instead')
const SocialsSearchResponse$json = {
  '1': 'SocialsSearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.SocialsSearchRequest',
      '10': 'next'
    },
    {
      '1': 'items',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.retrovibed.community.CommunitySocial',
      '10': 'items'
    },
    {
      '1': 'catalog',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.retrovibed.community.PluginPublisher',
      '10': 'catalog'
    },
  ],
};

/// Descriptor for `SocialsSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List socialsSearchResponseDescriptor = $convert.base64Decode(
    'ChVTb2NpYWxzU2VhcmNoUmVzcG9uc2USPgoEbmV4dBgBIAEoCzIqLnJldHJvdmliZWQuY29tbX'
    'VuaXR5LlNvY2lhbHNTZWFyY2hSZXF1ZXN0UgRuZXh0EjsKBWl0ZW1zGAIgAygLMiUucmV0cm92'
    'aWJlZC5jb21tdW5pdHkuQ29tbXVuaXR5U29jaWFsUgVpdGVtcxI/CgdjYXRhbG9nGAMgAygLMi'
    'UucmV0cm92aWJlZC5jb21tdW5pdHkuUGx1Z2luUHVibGlzaGVyUgdjYXRhbG9n');

@$core.Deprecated('Use pluginPublisherCreateResponseDescriptor instead')
const PluginPublisherCreateResponse$json = {
  '1': 'PluginPublisherCreateResponse',
  '2': [
    {
      '1': 'publisher',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.PluginPublisher',
      '10': 'publisher'
    },
  ],
};

/// Descriptor for `PluginPublisherCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pluginPublisherCreateResponseDescriptor =
    $convert.base64Decode(
        'Ch1QbHVnaW5QdWJsaXNoZXJDcmVhdGVSZXNwb25zZRJDCglwdWJsaXNoZXIYASABKAsyJS5yZX'
        'Ryb3ZpYmVkLmNvbW11bml0eS5QbHVnaW5QdWJsaXNoZXJSCXB1Ymxpc2hlcg==');

@$core.Deprecated('Use pluginPublisherDeleteResponseDescriptor instead')
const PluginPublisherDeleteResponse$json = {
  '1': 'PluginPublisherDeleteResponse',
  '2': [
    {
      '1': 'publisher',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.PluginPublisher',
      '10': 'publisher'
    },
  ],
};

/// Descriptor for `PluginPublisherDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pluginPublisherDeleteResponseDescriptor =
    $convert.base64Decode(
        'Ch1QbHVnaW5QdWJsaXNoZXJEZWxldGVSZXNwb25zZRJDCglwdWJsaXNoZXIYASABKAsyJS5yZX'
        'Ryb3ZpYmVkLmNvbW11bml0eS5QbHVnaW5QdWJsaXNoZXJSCXB1Ymxpc2hlcg==');

@$core.Deprecated('Use communityPublisherEnableResponseDescriptor instead')
const CommunityPublisherEnableResponse$json = {
  '1': 'CommunityPublisherEnableResponse',
  '2': [
    {
      '1': 'enabled',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.CommunityPublisher',
      '10': 'enabled'
    },
  ],
};

/// Descriptor for `CommunityPublisherEnableResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityPublisherEnableResponseDescriptor =
    $convert.base64Decode(
        'CiBDb21tdW5pdHlQdWJsaXNoZXJFbmFibGVSZXNwb25zZRJCCgdlbmFibGVkGAEgASgLMigucm'
        'V0cm92aWJlZC5jb21tdW5pdHkuQ29tbXVuaXR5UHVibGlzaGVyUgdlbmFibGVk');

@$core.Deprecated('Use communityPublisherDisableResponseDescriptor instead')
const CommunityPublisherDisableResponse$json = {
  '1': 'CommunityPublisherDisableResponse',
  '2': [
    {
      '1': 'disabled',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.CommunityPublisher',
      '10': 'disabled'
    },
  ],
};

/// Descriptor for `CommunityPublisherDisableResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityPublisherDisableResponseDescriptor =
    $convert.base64Decode(
        'CiFDb21tdW5pdHlQdWJsaXNoZXJEaXNhYmxlUmVzcG9uc2USRAoIZGlzYWJsZWQYASABKAsyKC'
        '5yZXRyb3ZpYmVkLmNvbW11bml0eS5Db21tdW5pdHlQdWJsaXNoZXJSCGRpc2FibGVk');
