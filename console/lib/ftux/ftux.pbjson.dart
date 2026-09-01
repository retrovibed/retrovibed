// This is a generated file - do not edit.
//
// Generated from ftux/ftux.proto.

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

@$core.Deprecated('Use communitySuggestionsDescriptor instead')
const CommunitySuggestions$json = {
  '1': 'CommunitySuggestions',
  '2': [
    {
      '1': 'community',
      '3': 1000,
      '4': 3,
      '5': 11,
      '6': '.retrovibed.community.Community',
      '10': 'community'
    },
  ],
  '9': [
    {'1': 1, '2': 1000},
  ],
};

/// Descriptor for `CommunitySuggestions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communitySuggestionsDescriptor = $convert.base64Decode(
    'ChRDb21tdW5pdHlTdWdnZXN0aW9ucxI+Cgljb21tdW5pdHkY6AcgAygLMh8ucmV0cm92aWJlZC'
    '5jb21tdW5pdHkuQ29tbXVuaXR5Ugljb21tdW5pdHlKBQgBEOgH');

@$core.Deprecated('Use subscribeCommunitiesRequestDescriptor instead')
const SubscribeCommunitiesRequest$json = {
  '1': 'SubscribeCommunitiesRequest',
  '2': [
    {'1': 'community_id', '3': 1, '4': 3, '5': 9, '10': 'community_id'},
  ],
  '9': [
    {'1': 2, '2': 1000},
  ],
};

/// Descriptor for `SubscribeCommunitiesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscribeCommunitiesRequestDescriptor =
    $convert.base64Decode(
        'ChtTdWJzY3JpYmVDb21tdW5pdGllc1JlcXVlc3QSIgoMY29tbXVuaXR5X2lkGAEgAygJUgxjb2'
        '1tdW5pdHlfaWRKBQgCEOgH');

@$core.Deprecated('Use subscribeCommunitiesResponseDescriptor instead')
const SubscribeCommunitiesResponse$json = {
  '1': 'SubscribeCommunitiesResponse',
};

/// Descriptor for `SubscribeCommunitiesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscribeCommunitiesResponseDescriptor =
    $convert.base64Decode('ChxTdWJzY3JpYmVDb21tdW5pdGllc1Jlc3BvbnNl');
