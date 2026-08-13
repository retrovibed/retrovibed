// This is a generated file - do not edit.
//
// Generated from community/community.publish.proto.

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
    {'1': 'encryption_seed', '3': 12, '4': 1, '5': 9, '10': 'encryption_seed'},
    {'1': 'bytes', '3': 13, '4': 1, '5': 4, '10': 'bytes'},
    {'1': 'library_id', '3': 1000, '4': 1, '5': 9, '10': 'library_id'},
    {
      '1': 'oauth_google_id',
      '3': 1001,
      '4': 1,
      '5': 9,
      '10': 'oauth_google_id'
    },
  ],
  '9': [
    {'1': 14, '2': 1000},
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
    'GgoIbWltZXR5cGUYCyABKAlSCG1pbWV0eXBlEigKD2VuY3J5cHRpb25fc2VlZBgMIAEoCVIPZW'
    '5jcnlwdGlvbl9zZWVkEhQKBWJ5dGVzGA0gASgEUgVieXRlcxIfCgpsaWJyYXJ5X2lkGOgHIAEo'
    'CVIKbGlicmFyeV9pZBIpCg9vYXV0aF9nb29nbGVfaWQY6QcgASgJUg9vYXV0aF9nb29nbGVfaW'
    'RKBQgOEOgH');

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

@$core.Deprecated('Use publishContentDeleteRequestDescriptor instead')
const PublishContentDeleteRequest$json = {
  '1': 'PublishContentDeleteRequest',
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

/// Descriptor for `PublishContentDeleteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishContentDeleteRequestDescriptor =
    $convert.base64Decode(
        'ChtQdWJsaXNoQ29udGVudERlbGV0ZVJlcXVlc3QSVAoRcHVibGlzaGVkX2NvbnRlbnQYASABKA'
        'syJi5yZXRyb3ZpYmVkLmNvbW11bml0eS5QdWJsaXNoZWRDb250ZW50UhFwdWJsaXNoZWRfY29u'
        'dGVudA==');

@$core.Deprecated('Use publishContentDeleteResponseDescriptor instead')
const PublishContentDeleteResponse$json = {
  '1': 'PublishContentDeleteResponse',
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

/// Descriptor for `PublishContentDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishContentDeleteResponseDescriptor =
    $convert.base64Decode(
        'ChxQdWJsaXNoQ29udGVudERlbGV0ZVJlc3BvbnNlElQKEXB1Ymxpc2hlZF9jb250ZW50GAEgAS'
        'gLMiYucmV0cm92aWJlZC5jb21tdW5pdHkuUHVibGlzaGVkQ29udGVudFIRcHVibGlzaGVkX2Nv'
        'bnRlbnQ=');

@$core.Deprecated('Use publishedContentSearchRequestDescriptor instead')
const PublishedContentSearchRequest$json = {
  '1': 'PublishedContentSearchRequest',
  '2': [
    {'1': 'community_id', '3': 1, '4': 1, '5': 9, '10': 'community_id'},
    {'1': 'query', '3': 2, '4': 1, '5': 9, '10': 'query'},
    {'1': 'sync', '3': 3, '4': 1, '5': 9, '10': 'sync'},
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 4, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `PublishedContentSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishedContentSearchRequestDescriptor = $convert.base64Decode(
    'Ch1QdWJsaXNoZWRDb250ZW50U2VhcmNoUmVxdWVzdBIiCgxjb21tdW5pdHlfaWQYASABKAlSDG'
    'NvbW11bml0eV9pZBIUCgVxdWVyeRgCIAEoCVIFcXVlcnkSEgoEc3luYxgDIAEoCVIEc3luYxIX'
    'CgZvZmZzZXQYhAcgASgEUgZvZmZzZXQSFQoFbGltaXQYhQcgASgEUgVsaW1pdEoFCAQQhAdKBg'
    'iGBxDoBw==');

@$core.Deprecated('Use publishedContentSearchResponseDescriptor instead')
const PublishedContentSearchResponse$json = {
  '1': 'PublishedContentSearchResponse',
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
      '6': '.retrovibed.community.PublishedContentSearchRequest',
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

/// Descriptor for `PublishedContentSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishedContentSearchResponseDescriptor = $convert.base64Decode(
    'Ch5QdWJsaXNoZWRDb250ZW50U2VhcmNoUmVzcG9uc2USPQoJY29tbXVuaXR5GAEgASgLMh8ucm'
    'V0cm92aWJlZC5jb21tdW5pdHkuQ29tbXVuaXR5Ugljb21tdW5pdHkSRwoEbmV4dBgCIAEoCzIz'
    'LnJldHJvdmliZWQuY29tbXVuaXR5LlB1Ymxpc2hlZENvbnRlbnRTZWFyY2hSZXF1ZXN0UgRuZX'
    'h0EjwKBWl0ZW1zGAMgAygLMiYucmV0cm92aWJlZC5jb21tdW5pdHkuUHVibGlzaGVkQ29udGVu'
    'dFIFaXRlbXM=');
