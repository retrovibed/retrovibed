// This is a generated file - do not edit.
//
// Generated from community.proto.

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
  ],
  '9': [
    {'1': 16, '2': 20},
    {'1': 20, '2': 1000},
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
    'bHRKBAgQEBRKBQgUEOgH');

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

@$core.Deprecated('Use publishedContentDescriptor instead')
const PublishedContent$json = {
  '1': 'PublishedContent',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'community_id', '3': 2, '4': 1, '5': 9, '10': 'community_id'},
    {'1': 'known_media_id', '3': 3, '4': 1, '5': 9, '10': 'known_media_id'},
    {'1': 'magnet_uri', '3': 4, '4': 1, '5': 9, '10': 'magnet_uri'},
    {'1': 'published_at', '3': 5, '4': 1, '5': 9, '10': 'published_at'},
    {'1': 'created_at', '3': 6, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 7, '4': 1, '5': 9, '10': 'updated_at'},
    {'1': 'archived_id', '3': 8, '4': 1, '5': 9, '10': 'archived_id'},
    {'1': 'title', '3': 9, '4': 1, '5': 9, '10': 'title'},
    {'1': 'description', '3': 10, '4': 1, '5': 9, '10': 'description'},
    {'1': 'mimetype', '3': 11, '4': 1, '5': 9, '10': 'mimetype'},
    {'1': 'library_id', '3': 12, '4': 1, '5': 9, '10': 'library_id'},
    {'1': 'oauth_google_id', '3': 13, '4': 1, '5': 9, '10': 'oauth_google_id'},
    {'1': 'encryption_seed', '3': 14, '4': 1, '5': 9, '10': 'encryption_seed'},
  ],
};

/// Descriptor for `PublishedContent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishedContentDescriptor = $convert.base64Decode(
    'ChBQdWJsaXNoZWRDb250ZW50Eg4KAmlkGAEgASgJUgJpZBIiCgxjb21tdW5pdHlfaWQYAiABKA'
    'lSDGNvbW11bml0eV9pZBImCg5rbm93bl9tZWRpYV9pZBgDIAEoCVIOa25vd25fbWVkaWFfaWQS'
    'HgoKbWFnbmV0X3VyaRgEIAEoCVIKbWFnbmV0X3VyaRIiCgxwdWJsaXNoZWRfYXQYBSABKAlSDH'
    'B1Ymxpc2hlZF9hdBIeCgpjcmVhdGVkX2F0GAYgASgJUgpjcmVhdGVkX2F0Eh4KCnVwZGF0ZWRf'
    'YXQYByABKAlSCnVwZGF0ZWRfYXQSIAoLYXJjaGl2ZWRfaWQYCCABKAlSC2FyY2hpdmVkX2lkEh'
    'QKBXRpdGxlGAkgASgJUgV0aXRsZRIgCgtkZXNjcmlwdGlvbhgKIAEoCVILZGVzY3JpcHRpb24S'
    'GgoIbWltZXR5cGUYCyABKAlSCG1pbWV0eXBlEh4KCmxpYnJhcnlfaWQYDCABKAlSCmxpYnJhcn'
    'lfaWQSKAoPb2F1dGhfZ29vZ2xlX2lkGA0gASgJUg9vYXV0aF9nb29nbGVfaWQSKAoPZW5jcnlw'
    'dGlvbl9zZWVkGA4gASgJUg9lbmNyeXB0aW9uX3NlZWQ=');

@$core.Deprecated('Use publishContentRequestDescriptor instead')
const PublishContentRequest$json = {
  '1': 'PublishContentRequest',
  '2': [
    {
      '1': 'published_content',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.PublishedContent',
      '10': 'published_content'
    },
    {
      '1': 'publish_mode',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.retrovibed.community.PublishMode',
      '10': 'publish_mode'
    },
  ],
  '9': [
    {'1': 3, '2': 4},
  ],
};

/// Descriptor for `PublishContentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishContentRequestDescriptor = $convert.base64Decode(
    'ChVQdWJsaXNoQ29udGVudFJlcXVlc3QSVAoRcHVibGlzaGVkX2NvbnRlbnQYASABKAsyJi5yZX'
    'Ryb3ZpYmVkLmNvbW11bml0eS5QdWJsaXNoZWRDb250ZW50UhFwdWJsaXNoZWRfY29udGVudBJF'
    'CgxwdWJsaXNoX21vZGUYAiABKA4yIS5yZXRyb3ZpYmVkLmNvbW11bml0eS5QdWJsaXNoTW9kZV'
    'IMcHVibGlzaF9tb2RlSgQIAxAE');

@$core.Deprecated('Use publishContentResponseDescriptor instead')
const PublishContentResponse$json = {
  '1': 'PublishContentResponse',
  '2': [
    {
      '1': 'published_content',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.PublishedContent',
      '10': 'published_content'
    },
  ],
};

/// Descriptor for `PublishContentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishContentResponseDescriptor = $convert.base64Decode(
    'ChZQdWJsaXNoQ29udGVudFJlc3BvbnNlElQKEXB1Ymxpc2hlZF9jb250ZW50GAEgASgLMiYucm'
    'V0cm92aWJlZC5jb21tdW5pdHkuUHVibGlzaGVkQ29udGVudFIRcHVibGlzaGVkX2NvbnRlbnQ=');

@$core.Deprecated('Use publishedContentListRequestDescriptor instead')
const PublishedContentListRequest$json = {
  '1': 'PublishedContentListRequest',
  '2': [
    {'1': 'community_id', '3': 1, '4': 1, '5': 9, '10': 'community_id'},
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 2, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `PublishedContentListRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishedContentListRequestDescriptor =
    $convert.base64Decode(
        'ChtQdWJsaXNoZWRDb250ZW50TGlzdFJlcXVlc3QSIgoMY29tbXVuaXR5X2lkGAEgASgJUgxjb2'
        '1tdW5pdHlfaWQSFwoGb2Zmc2V0GIQHIAEoBFIGb2Zmc2V0EhUKBWxpbWl0GIUHIAEoBFIFbGlt'
        'aXRKBQgCEIQHSgYIhgcQ6Ac=');

@$core.Deprecated('Use publishedContentListResponseDescriptor instead')
const PublishedContentListResponse$json = {
  '1': 'PublishedContentListResponse',
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
      '1': 'next',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.PublishedContentListRequest',
      '10': 'next'
    },
    {
      '1': 'items',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.retrovibed.community.PublishedContent',
      '10': 'items'
    },
  ],
};

/// Descriptor for `PublishedContentListResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishedContentListResponseDescriptor = $convert.base64Decode(
    'ChxQdWJsaXNoZWRDb250ZW50TGlzdFJlc3BvbnNlEj0KCWNvbW11bml0eRgBIAEoCzIfLnJldH'
    'JvdmliZWQuY29tbXVuaXR5LkNvbW11bml0eVIJY29tbXVuaXR5EkUKBG5leHQYAiABKAsyMS5y'
    'ZXRyb3ZpYmVkLmNvbW11bml0eS5QdWJsaXNoZWRDb250ZW50TGlzdFJlcXVlc3RSBG5leHQSPA'
    'oFaXRlbXMYAyADKAsyJi5yZXRyb3ZpYmVkLmNvbW11bml0eS5QdWJsaXNoZWRDb250ZW50UgVp'
    'dGVtcw==');

@$core.Deprecated('Use communityMetricDescriptor instead')
const CommunityMetric$json = {
  '1': 'CommunityMetric',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'community_id', '3': 2, '4': 1, '5': 9, '10': 'community_id'},
    {'1': 'period_start', '3': 3, '4': 1, '5': 9, '10': 'period_start'},
    {'1': 'period_end', '3': 4, '4': 1, '5': 9, '10': 'period_end'},
    {'1': 'subscribers', '3': 5, '4': 1, '5': 13, '10': 'subscribers'},
  ],
};

/// Descriptor for `CommunityMetric`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityMetricDescriptor = $convert.base64Decode(
    'Cg9Db21tdW5pdHlNZXRyaWMSDgoCaWQYASABKAlSAmlkEiIKDGNvbW11bml0eV9pZBgCIAEoCV'
    'IMY29tbXVuaXR5X2lkEiIKDHBlcmlvZF9zdGFydBgDIAEoCVIMcGVyaW9kX3N0YXJ0Eh4KCnBl'
    'cmlvZF9lbmQYBCABKAlSCnBlcmlvZF9lbmQSIAoLc3Vic2NyaWJlcnMYBSABKA1SC3N1YnNjcm'
    'liZXJz');

@$core.Deprecated('Use publishedContentMetricDescriptor instead')
const PublishedContentMetric$json = {
  '1': 'PublishedContentMetric',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'published_content_id',
      '3': 2,
      '4': 1,
      '5': 9,
      '10': 'published_content_id'
    },
    {'1': 'period_start', '3': 3, '4': 1, '5': 9, '10': 'period_start'},
    {'1': 'period_end', '3': 4, '4': 1, '5': 9, '10': 'period_end'},
    {'1': 'archivers', '3': 5, '4': 1, '5': 13, '10': 'archivers'},
    {'1': 'bytes', '3': 6, '4': 1, '5': 3, '10': 'bytes'},
    {'1': 'revenue', '3': 7, '4': 1, '5': 3, '10': 'revenue'},
  ],
};

/// Descriptor for `PublishedContentMetric`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishedContentMetricDescriptor = $convert.base64Decode(
    'ChZQdWJsaXNoZWRDb250ZW50TWV0cmljEg4KAmlkGAEgASgJUgJpZBIyChRwdWJsaXNoZWRfY2'
    '9udGVudF9pZBgCIAEoCVIUcHVibGlzaGVkX2NvbnRlbnRfaWQSIgoMcGVyaW9kX3N0YXJ0GAMg'
    'ASgJUgxwZXJpb2Rfc3RhcnQSHgoKcGVyaW9kX2VuZBgEIAEoCVIKcGVyaW9kX2VuZBIcCglhcm'
    'NoaXZlcnMYBSABKA1SCWFyY2hpdmVycxIUCgVieXRlcxgGIAEoA1IFYnl0ZXMSGAoHcmV2ZW51'
    'ZRgHIAEoA1IHcmV2ZW51ZQ==');

@$core.Deprecated('Use communityMetricsRequestDescriptor instead')
const CommunityMetricsRequest$json = {
  '1': 'CommunityMetricsRequest',
  '2': [
    {'1': 'community_id', '3': 1, '4': 1, '5': 9, '10': 'community_id'},
    {'1': 'period', '3': 2, '4': 1, '5': 9, '10': 'period'},
    {'1': 'start_date', '3': 3, '4': 1, '5': 9, '10': 'start_date'},
    {'1': 'end_date', '3': 4, '4': 1, '5': 9, '10': 'end_date'},
  ],
};

/// Descriptor for `CommunityMetricsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityMetricsRequestDescriptor = $convert.base64Decode(
    'ChdDb21tdW5pdHlNZXRyaWNzUmVxdWVzdBIiCgxjb21tdW5pdHlfaWQYASABKAlSDGNvbW11bm'
    'l0eV9pZBIWCgZwZXJpb2QYAiABKAlSBnBlcmlvZBIeCgpzdGFydF9kYXRlGAMgASgJUgpzdGFy'
    'dF9kYXRlEhoKCGVuZF9kYXRlGAQgASgJUghlbmRfZGF0ZQ==');

@$core.Deprecated('Use communityMetricsResponseDescriptor instead')
const CommunityMetricsResponse$json = {
  '1': 'CommunityMetricsResponse',
  '2': [
    {
      '1': 'summary',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.retrovibed.community.CommunityMetric',
      '10': 'summary'
    },
    {'1': 'total_archivers', '3': 2, '4': 1, '5': 5, '10': 'total_archivers'},
    {
      '1': 'items',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.retrovibed.community.PublishedContentMetric',
      '10': 'items'
    },
  ],
};

/// Descriptor for `CommunityMetricsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List communityMetricsResponseDescriptor = $convert.base64Decode(
    'ChhDb21tdW5pdHlNZXRyaWNzUmVzcG9uc2USPwoHc3VtbWFyeRgBIAEoCzIlLnJldHJvdmliZW'
    'QuY29tbXVuaXR5LkNvbW11bml0eU1ldHJpY1IHc3VtbWFyeRIoCg90b3RhbF9hcmNoaXZlcnMY'
    'AiABKAVSD3RvdGFsX2FyY2hpdmVycxJCCgVpdGVtcxgDIAMoCzIsLnJldHJvdmliZWQuY29tbX'
    'VuaXR5LlB1Ymxpc2hlZENvbnRlbnRNZXRyaWNSBWl0ZW1z');

@$core.Deprecated('Use metricsSyncRequestDescriptor instead')
const MetricsSyncRequest$json = {
  '1': 'MetricsSyncRequest',
  '2': [
    {'1': 'community_id', '3': 1, '4': 1, '5': 9, '10': 'community_id'},
    {'1': 'since', '3': 2, '4': 1, '5': 9, '10': 'since'},
  ],
};

/// Descriptor for `MetricsSyncRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricsSyncRequestDescriptor = $convert.base64Decode(
    'ChJNZXRyaWNzU3luY1JlcXVlc3QSIgoMY29tbXVuaXR5X2lkGAEgASgJUgxjb21tdW5pdHlfaW'
    'QSFAoFc2luY2UYAiABKAlSBXNpbmNl');

@$core.Deprecated('Use metricsSyncResponseDescriptor instead')
const MetricsSyncResponse$json = {
  '1': 'MetricsSyncResponse',
  '2': [
    {
      '1': 'community_metrics',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.retrovibed.community.CommunityMetric',
      '10': 'community_metrics'
    },
    {
      '1': 'content_metrics',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.retrovibed.community.PublishedContentMetric',
      '10': 'content_metrics'
    },
    {'1': 'synced_at', '3': 3, '4': 1, '5': 9, '10': 'synced_at'},
    {'1': 'complete', '3': 4, '4': 1, '5': 8, '10': 'complete'},
  ],
};

/// Descriptor for `MetricsSyncResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricsSyncResponseDescriptor = $convert.base64Decode(
    'ChNNZXRyaWNzU3luY1Jlc3BvbnNlElMKEWNvbW11bml0eV9tZXRyaWNzGAEgAygLMiUucmV0cm'
    '92aWJlZC5jb21tdW5pdHkuQ29tbXVuaXR5TWV0cmljUhFjb21tdW5pdHlfbWV0cmljcxJWCg9j'
    'b250ZW50X21ldHJpY3MYAiADKAsyLC5yZXRyb3ZpYmVkLmNvbW11bml0eS5QdWJsaXNoZWRDb2'
    '50ZW50TWV0cmljUg9jb250ZW50X21ldHJpY3MSHAoJc3luY2VkX2F0GAMgASgJUglzeW5jZWRf'
    'YXQSGgoIY29tcGxldGUYBCABKAhSCGNvbXBsZXRl');

@$core.Deprecated('Use metricsSyncProgressDescriptor instead')
const MetricsSyncProgress$json = {
  '1': 'MetricsSyncProgress',
  '2': [
    {'1': 'status', '3': 1, '4': 1, '5': 9, '10': 'status'},
    {
      '1': 'community_metrics_count',
      '3': 2,
      '4': 1,
      '5': 5,
      '10': 'community_metrics_count'
    },
    {
      '1': 'content_metrics_count',
      '3': 3,
      '4': 1,
      '5': 5,
      '10': 'content_metrics_count'
    },
    {'1': 'error', '3': 4, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `MetricsSyncProgress`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricsSyncProgressDescriptor = $convert.base64Decode(
    'ChNNZXRyaWNzU3luY1Byb2dyZXNzEhYKBnN0YXR1cxgBIAEoCVIGc3RhdHVzEjgKF2NvbW11bm'
    'l0eV9tZXRyaWNzX2NvdW50GAIgASgFUhdjb21tdW5pdHlfbWV0cmljc19jb3VudBI0ChVjb250'
    'ZW50X21ldHJpY3NfY291bnQYAyABKAVSFWNvbnRlbnRfbWV0cmljc19jb3VudBIUCgVlcnJvch'
    'gEIAEoCVIFZXJyb3I=');

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
