// This is a generated file - do not edit.
//
// Generated from media.known.proto.

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

@$core.Deprecated('Use knownDescriptor instead')
const Known$json = {
  '1': 'Known',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'rating', '3': 2, '4': 1, '5': 2, '10': 'rating'},
    {'1': 'adult', '3': 3, '4': 1, '5': 8, '10': 'adult'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {'1': 'summary', '3': 5, '4': 1, '5': 9, '10': 'summary'},
    {'1': 'image', '3': 6, '4': 1, '5': 9, '10': 'image'},
    {'1': 'released', '3': 7, '4': 1, '5': 9, '10': 'released'},
    {'1': 'mimetype', '3': 8, '4': 1, '5': 9, '10': 'mimetype'},
    {'1': 'source', '3': 9, '4': 1, '5': 9, '10': 'source'},
    {'1': 'uid', '3': 10, '4': 1, '5': 9, '10': 'uid'},
  ],
};

/// Descriptor for `Known`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List knownDescriptor = $convert.base64Decode(
    'CgVLbm93bhIOCgJpZBgBIAEoCVICaWQSFgoGcmF0aW5nGAIgASgCUgZyYXRpbmcSFAoFYWR1bH'
    'QYAyABKAhSBWFkdWx0EiAKC2Rlc2NyaXB0aW9uGAQgASgJUgtkZXNjcmlwdGlvbhIYCgdzdW1t'
    'YXJ5GAUgASgJUgdzdW1tYXJ5EhQKBWltYWdlGAYgASgJUgVpbWFnZRIaCghyZWxlYXNlZBgHIA'
    'EoCVIIcmVsZWFzZWQSGgoIbWltZXR5cGUYCCABKAlSCG1pbWV0eXBlEhYKBnNvdXJjZRgJIAEo'
    'CVIGc291cmNlEhAKA3VpZBgKIAEoCVIDdWlk');

@$core.Deprecated('Use knownSearchRequestDescriptor instead')
const KnownSearchRequest$json = {
  '1': 'KnownSearchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'adult', '3': 2, '4': 1, '5': 8, '10': 'adult'},
    {'1': 'language', '3': 3, '4': 1, '5': 9, '10': 'language'},
    {'1': 'mimetype', '3': 4, '4': 1, '5': 9, '10': 'mimetype'},
    {
      '1': 'released',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.meta.DateRange',
      '10': 'released'
    },
    {'1': 'source', '3': 6, '4': 3, '5': 9, '10': 'source'},
    {'1': 'id', '3': 7, '4': 3, '5': 9, '10': 'id'},
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 8, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `KnownSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List knownSearchRequestDescriptor = $convert.base64Decode(
    'ChJLbm93blNlYXJjaFJlcXVlc3QSFAoFcXVlcnkYASABKAlSBXF1ZXJ5EhQKBWFkdWx0GAIgAS'
    'gIUgVhZHVsdBIaCghsYW5ndWFnZRgDIAEoCVIIbGFuZ3VhZ2USGgoIbWltZXR5cGUYBCABKAlS'
    'CG1pbWV0eXBlEisKCHJlbGVhc2VkGAUgASgLMg8ubWV0YS5EYXRlUmFuZ2VSCHJlbGVhc2VkEh'
    'YKBnNvdXJjZRgGIAMoCVIGc291cmNlEg4KAmlkGAcgAygJUgJpZBIXCgZvZmZzZXQYhAcgASgE'
    'UgZvZmZzZXQSFQoFbGltaXQYhQcgASgEUgVsaW1pdEoFCAgQhAdKBgiGBxDoBw==');

@$core.Deprecated('Use knownSearchResponseDescriptor instead')
const KnownSearchResponse$json = {
  '1': 'KnownSearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.KnownSearchRequest',
      '10': 'next'
    },
    {'1': 'items', '3': 2, '4': 3, '5': 11, '6': '.media.Known', '10': 'items'},
  ],
};

/// Descriptor for `KnownSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List knownSearchResponseDescriptor = $convert.base64Decode(
    'ChNLbm93blNlYXJjaFJlc3BvbnNlEi0KBG5leHQYASABKAsyGS5tZWRpYS5Lbm93blNlYXJjaF'
    'JlcXVlc3RSBG5leHQSIgoFaXRlbXMYAiADKAsyDC5tZWRpYS5Lbm93blIFaXRlbXM=');

@$core.Deprecated('Use knownMatchRequestDescriptor instead')
const KnownMatchRequest$json = {
  '1': 'KnownMatchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
  ],
};

/// Descriptor for `KnownMatchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List knownMatchRequestDescriptor = $convert
    .base64Decode('ChFLbm93bk1hdGNoUmVxdWVzdBIUCgVxdWVyeRgBIAEoCVIFcXVlcnk=');

@$core.Deprecated('Use knownLookupRequestDescriptor instead')
const KnownLookupRequest$json = {
  '1': 'KnownLookupRequest',
};

/// Descriptor for `KnownLookupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List knownLookupRequestDescriptor =
    $convert.base64Decode('ChJLbm93bkxvb2t1cFJlcXVlc3Q=');

@$core.Deprecated('Use knownLookupResponseDescriptor instead')
const KnownLookupResponse$json = {
  '1': 'KnownLookupResponse',
  '2': [
    {'1': 'known', '3': 1, '4': 1, '5': 11, '6': '.media.Known', '10': 'known'},
  ],
};

/// Descriptor for `KnownLookupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List knownLookupResponseDescriptor = $convert.base64Decode(
    'ChNLbm93bkxvb2t1cFJlc3BvbnNlEiIKBWtub3duGAEgASgLMgwubWVkaWEuS25vd25SBWtub3'
    'du');

@$core.Deprecated('Use knownDownloadRequestDescriptor instead')
const KnownDownloadRequest$json = {
  '1': 'KnownDownloadRequest',
};

/// Descriptor for `KnownDownloadRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List knownDownloadRequestDescriptor =
    $convert.base64Decode('ChRLbm93bkRvd25sb2FkUmVxdWVzdA==');

@$core.Deprecated('Use knownDownloadResponseDescriptor instead')
const KnownDownloadResponse$json = {
  '1': 'KnownDownloadResponse',
  '2': [
    {'1': 'known', '3': 1, '4': 1, '5': 11, '6': '.media.Known', '10': 'known'},
  ],
};

/// Descriptor for `KnownDownloadResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List knownDownloadResponseDescriptor = $convert.base64Decode(
    'ChVLbm93bkRvd25sb2FkUmVzcG9uc2USIgoFa25vd24YASABKAsyDC5tZWRpYS5Lbm93blIFa2'
    '5vd24=');

@$core.Deprecated('Use knownCreateRequestDescriptor instead')
const KnownCreateRequest$json = {
  '1': 'KnownCreateRequest',
  '2': [
    {'1': 'known', '3': 1, '4': 1, '5': 11, '6': '.media.Known', '10': 'known'},
  ],
};

/// Descriptor for `KnownCreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List knownCreateRequestDescriptor = $convert.base64Decode(
    'ChJLbm93bkNyZWF0ZVJlcXVlc3QSIgoFa25vd24YASABKAsyDC5tZWRpYS5Lbm93blIFa25vd2'
    '4=');

@$core.Deprecated('Use knownCreateResponseDescriptor instead')
const KnownCreateResponse$json = {
  '1': 'KnownCreateResponse',
  '2': [
    {'1': 'known', '3': 1, '4': 1, '5': 11, '6': '.media.Known', '10': 'known'},
  ],
};

/// Descriptor for `KnownCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List knownCreateResponseDescriptor = $convert.base64Decode(
    'ChNLbm93bkNyZWF0ZVJlc3BvbnNlEiIKBWtub3duGAEgASgLMgwubWVkaWEuS25vd25SBWtub3'
    'du');

@$core.Deprecated('Use knownLatestRequestDescriptor instead')
const KnownLatestRequest$json = {
  '1': 'KnownLatestRequest',
  '2': [
    {
      '1': 'released',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.DateRange',
      '10': 'released'
    },
    {'1': 'adult', '3': 2, '4': 1, '5': 8, '10': 'adult'},
    {'1': 'language', '3': 3, '4': 1, '5': 9, '10': 'language'},
    {'1': 'mimetype', '3': 4, '4': 1, '5': 9, '10': 'mimetype'},
    {'1': 'source', '3': 5, '4': 3, '5': 9, '10': 'source'},
    {'1': 'id', '3': 6, '4': 3, '5': 9, '10': 'id'},
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 7, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `KnownLatestRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List knownLatestRequestDescriptor = $convert.base64Decode(
    'ChJLbm93bkxhdGVzdFJlcXVlc3QSKwoIcmVsZWFzZWQYASABKAsyDy5tZXRhLkRhdGVSYW5nZV'
    'IIcmVsZWFzZWQSFAoFYWR1bHQYAiABKAhSBWFkdWx0EhoKCGxhbmd1YWdlGAMgASgJUghsYW5n'
    'dWFnZRIaCghtaW1ldHlwZRgEIAEoCVIIbWltZXR5cGUSFgoGc291cmNlGAUgAygJUgZzb3VyY2'
    'USDgoCaWQYBiADKAlSAmlkEhcKBm9mZnNldBiEByABKARSBm9mZnNldBIVCgVsaW1pdBiFByAB'
    'KARSBWxpbWl0SgUIBxCEB0oGCIYHEOgH');

@$core.Deprecated('Use knownLatestResponseDescriptor instead')
const KnownLatestResponse$json = {
  '1': 'KnownLatestResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.KnownLatestRequest',
      '10': 'next'
    },
    {'1': 'items', '3': 2, '4': 3, '5': 11, '6': '.media.Known', '10': 'items'},
  ],
};

/// Descriptor for `KnownLatestResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List knownLatestResponseDescriptor = $convert.base64Decode(
    'ChNLbm93bkxhdGVzdFJlc3BvbnNlEi0KBG5leHQYASABKAsyGS5tZWRpYS5Lbm93bkxhdGVzdF'
    'JlcXVlc3RSBG5leHQSIgoFaXRlbXMYAiADKAsyDC5tZWRpYS5Lbm93blIFaXRlbXM=');

@$core.Deprecated('Use recommendationSearchRequestDescriptor instead')
const RecommendationSearchRequest$json = {
  '1': 'RecommendationSearchRequest',
  '2': [
    {'1': 'mimetype', '3': 1, '4': 1, '5': 9, '10': 'mimetype'},
    {'1': 'adult', '3': 2, '4': 1, '5': 8, '10': 'adult'},
    {'1': 'language', '3': 3, '4': 1, '5': 9, '10': 'language'},
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 4, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `RecommendationSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recommendationSearchRequestDescriptor = $convert.base64Decode(
    'ChtSZWNvbW1lbmRhdGlvblNlYXJjaFJlcXVlc3QSGgoIbWltZXR5cGUYASABKAlSCG1pbWV0eX'
    'BlEhQKBWFkdWx0GAIgASgIUgVhZHVsdBIaCghsYW5ndWFnZRgDIAEoCVIIbGFuZ3VhZ2USFwoG'
    'b2Zmc2V0GIQHIAEoBFIGb2Zmc2V0EhUKBWxpbWl0GIUHIAEoBFIFbGltaXRKBQgEEIQHSgYIhg'
    'cQ6Ac=');

@$core.Deprecated('Use recommendationSearchResponseDescriptor instead')
const RecommendationSearchResponse$json = {
  '1': 'RecommendationSearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.RecommendationSearchRequest',
      '10': 'next'
    },
    {'1': 'items', '3': 2, '4': 3, '5': 11, '6': '.media.Known', '10': 'items'},
  ],
};

/// Descriptor for `RecommendationSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recommendationSearchResponseDescriptor =
    $convert.base64Decode(
        'ChxSZWNvbW1lbmRhdGlvblNlYXJjaFJlc3BvbnNlEjYKBG5leHQYASABKAsyIi5tZWRpYS5SZW'
        'NvbW1lbmRhdGlvblNlYXJjaFJlcXVlc3RSBG5leHQSIgoFaXRlbXMYAiADKAsyDC5tZWRpYS5L'
        'bm93blIFaXRlbXM=');

@$core.Deprecated('Use recommendationFindRequestDescriptor instead')
const RecommendationFindRequest$json = {
  '1': 'RecommendationFindRequest',
};

/// Descriptor for `RecommendationFindRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recommendationFindRequestDescriptor =
    $convert.base64Decode('ChlSZWNvbW1lbmRhdGlvbkZpbmRSZXF1ZXN0');

@$core.Deprecated('Use recommendationFindResponseDescriptor instead')
const RecommendationFindResponse$json = {
  '1': 'RecommendationFindResponse',
  '2': [
    {
      '1': 'recommendation',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.Known',
      '10': 'recomendation'
    },
  ],
};

/// Descriptor for `RecommendationFindResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recommendationFindResponseDescriptor =
    $convert.base64Decode(
        'ChpSZWNvbW1lbmRhdGlvbkZpbmRSZXNwb25zZRIzCg5yZWNvbW1lbmRhdGlvbhgBIAEoCzIMLm'
        '1lZGlhLktub3duUg1yZWNvbWVuZGF0aW9u');

@$core.Deprecated('Use recommendationDeleteRequestDescriptor instead')
const RecommendationDeleteRequest$json = {
  '1': 'RecommendationDeleteRequest',
};

/// Descriptor for `RecommendationDeleteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recommendationDeleteRequestDescriptor =
    $convert.base64Decode('ChtSZWNvbW1lbmRhdGlvbkRlbGV0ZVJlcXVlc3Q=');

@$core.Deprecated('Use recommendationDeleteResponseDescriptor instead')
const RecommendationDeleteResponse$json = {
  '1': 'RecommendationDeleteResponse',
  '2': [
    {
      '1': 'recommendation',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.Known',
      '10': 'recomendation'
    },
  ],
};

/// Descriptor for `RecommendationDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recommendationDeleteResponseDescriptor =
    $convert.base64Decode(
        'ChxSZWNvbW1lbmRhdGlvbkRlbGV0ZVJlc3BvbnNlEjMKDnJlY29tbWVuZGF0aW9uGAEgASgLMg'
        'wubWVkaWEuS25vd25SDXJlY29tZW5kYXRpb24=');

@$core.Deprecated('Use recommendationRefreshRequestDescriptor instead')
const RecommendationRefreshRequest$json = {
  '1': 'RecommendationRefreshRequest',
};

/// Descriptor for `RecommendationRefreshRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recommendationRefreshRequestDescriptor =
    $convert.base64Decode('ChxSZWNvbW1lbmRhdGlvblJlZnJlc2hSZXF1ZXN0');

@$core.Deprecated('Use recommendationRefreshResponseDescriptor instead')
const RecommendationRefreshResponse$json = {
  '1': 'RecommendationRefreshResponse',
  '2': [
    {
      '1': 'recommendation',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.Known',
      '10': 'recomendation'
    },
  ],
};

/// Descriptor for `RecommendationRefreshResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recommendationRefreshResponseDescriptor =
    $convert.base64Decode(
        'Ch1SZWNvbW1lbmRhdGlvblJlZnJlc2hSZXNwb25zZRIzCg5yZWNvbW1lbmRhdGlvbhgBIAEoCz'
        'IMLm1lZGlhLktub3duUg1yZWNvbWVuZGF0aW9u');
