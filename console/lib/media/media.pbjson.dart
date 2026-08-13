// This is a generated file - do not edit.
//
// Generated from media/media.proto.

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
    {'1': 'description', '3': 2, '4': 1, '5': 9, '10': 'description'},
    {'1': 'mimetype', '3': 3, '4': 1, '5': 9, '10': 'mimetype'},
    {'1': 'image', '3': 4, '4': 1, '5': 9, '10': 'image'},
    {'1': 'archive_id', '3': 5, '4': 1, '5': 9, '10': 'archive_id'},
    {'1': 'torrent_id', '3': 6, '4': 1, '5': 9, '10': 'torrent_id'},
    {'1': 'created_at', '3': 7, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 8, '4': 1, '5': 9, '10': 'updated_at'},
    {'1': 'known_media_id', '3': 9, '4': 1, '5': 9, '10': 'known_media_id'},
    {'1': 'encryption_seed', '3': 10, '4': 1, '5': 9, '10': 'encryption_seed'},
    {'1': 'uri', '3': 11, '4': 1, '5': 9, '10': 'uri'},
  ],
};

/// Descriptor for `Media`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaDescriptor = $convert.base64Decode(
    'CgVNZWRpYRIOCgJpZBgBIAEoCVICaWQSIAoLZGVzY3JpcHRpb24YAiABKAlSC2Rlc2NyaXB0aW'
    '9uEhoKCG1pbWV0eXBlGAMgASgJUghtaW1ldHlwZRIUCgVpbWFnZRgEIAEoCVIFaW1hZ2USHgoK'
    'YXJjaGl2ZV9pZBgFIAEoCVIKYXJjaGl2ZV9pZBIeCgp0b3JyZW50X2lkGAYgASgJUgp0b3JyZW'
    '50X2lkEh4KCmNyZWF0ZWRfYXQYByABKAlSCmNyZWF0ZWRfYXQSHgoKdXBkYXRlZF9hdBgIIAEo'
    'CVIKdXBkYXRlZF9hdBImCg5rbm93bl9tZWRpYV9pZBgJIAEoCVIOa25vd25fbWVkaWFfaWQSKA'
    'oPZW5jcnlwdGlvbl9zZWVkGAogASgJUg9lbmNyeXB0aW9uX3NlZWQSEAoDdXJpGAsgASgJUgN1'
    'cmk=');

@$core.Deprecated('Use mediaSearchRequestDescriptor instead')
const MediaSearchRequest$json = {
  '1': 'MediaSearchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'mimetypes', '3': 2, '4': 3, '5': 9, '10': 'mimetypes'},
    {'1': 'adult', '3': 3, '4': 1, '5': 8, '10': 'adult'},
    {'1': 'hidden', '3': 4, '4': 1, '5': 8, '10': 'hidden'},
    {'1': 'excluded', '3': 5, '4': 3, '5': 9, '10': 'excluded'},
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 6, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `MediaSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaSearchRequestDescriptor = $convert.base64Decode(
    'ChJNZWRpYVNlYXJjaFJlcXVlc3QSFAoFcXVlcnkYASABKAlSBXF1ZXJ5EhwKCW1pbWV0eXBlcx'
    'gCIAMoCVIJbWltZXR5cGVzEhQKBWFkdWx0GAMgASgIUgVhZHVsdBIWCgZoaWRkZW4YBCABKAhS'
    'BmhpZGRlbhIaCghleGNsdWRlZBgFIAMoCVIIZXhjbHVkZWQSFwoGb2Zmc2V0GIQHIAEoBFIGb2'
    'Zmc2V0EhUKBWxpbWl0GIUHIAEoBFIFbGltaXRKBQgGEIQHSgYIhgcQ6Ac=');

@$core.Deprecated('Use mediaSearchResponseDescriptor instead')
const MediaSearchResponse$json = {
  '1': 'MediaSearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.MediaSearchRequest',
      '10': 'next'
    },
    {'1': 'items', '3': 2, '4': 3, '5': 11, '6': '.media.Media', '10': 'items'},
  ],
};

/// Descriptor for `MediaSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaSearchResponseDescriptor = $convert.base64Decode(
    'ChNNZWRpYVNlYXJjaFJlc3BvbnNlEi0KBG5leHQYASABKAsyGS5tZWRpYS5NZWRpYVNlYXJjaF'
    'JlcXVlc3RSBG5leHQSIgoFaXRlbXMYAiADKAsyDC5tZWRpYS5NZWRpYVIFaXRlbXM=');

@$core.Deprecated('Use mediaFindResponseDescriptor instead')
const MediaFindResponse$json = {
  '1': 'MediaFindResponse',
  '2': [
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
  ],
};

/// Descriptor for `MediaFindResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaFindResponseDescriptor = $convert.base64Decode(
    'ChFNZWRpYUZpbmRSZXNwb25zZRIiCgVtZWRpYRgBIAEoCzIMLm1lZGlhLk1lZGlhUgVtZWRpYQ'
    '==');

@$core.Deprecated('Use mediaUpdateRequestDescriptor instead')
const MediaUpdateRequest$json = {
  '1': 'MediaUpdateRequest',
  '2': [
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
  ],
};

/// Descriptor for `MediaUpdateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaUpdateRequestDescriptor = $convert.base64Decode(
    'ChJNZWRpYVVwZGF0ZVJlcXVlc3QSIgoFbWVkaWEYASABKAsyDC5tZWRpYS5NZWRpYVIFbWVkaW'
    'E=');

@$core.Deprecated('Use mediaUpdateResponseDescriptor instead')
const MediaUpdateResponse$json = {
  '1': 'MediaUpdateResponse',
  '2': [
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
  ],
};

/// Descriptor for `MediaUpdateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaUpdateResponseDescriptor = $convert.base64Decode(
    'ChNNZWRpYVVwZGF0ZVJlc3BvbnNlEiIKBW1lZGlhGAEgASgLMgwubWVkaWEuTWVkaWFSBW1lZG'
    'lh');

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
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
  ],
};

/// Descriptor for `MediaDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaDeleteResponseDescriptor = $convert.base64Decode(
    'ChNNZWRpYURlbGV0ZVJlc3BvbnNlEiIKBW1lZGlhGAEgASgLMgwubWVkaWEuTWVkaWFSBW1lZG'
    'lh');

@$core.Deprecated('Use mediaUploadResponseDescriptor instead')
const MediaUploadResponse$json = {
  '1': 'MediaUploadResponse',
  '2': [
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
  ],
};

/// Descriptor for `MediaUploadResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaUploadResponseDescriptor = $convert.base64Decode(
    'ChNNZWRpYVVwbG9hZFJlc3BvbnNlEiIKBW1lZGlhGAEgASgLMgwubWVkaWEuTWVkaWFSBW1lZG'
    'lh');

@$core.Deprecated('Use downloadDescriptor instead')
const Download$json = {
  '1': 'Download',
  '2': [
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
    {'1': 'bytes', '3': 2, '4': 1, '5': 4, '10': 'bytes'},
    {'1': 'downloaded', '3': 3, '4': 1, '5': 4, '10': 'downloaded'},
    {'1': 'initiated_at', '3': 4, '4': 1, '5': 9, '10': 'initiated_at'},
    {'1': 'paused_at', '3': 5, '4': 1, '5': 9, '10': 'paused_at'},
    {'1': 'peers', '3': 6, '4': 1, '5': 13, '10': 'peers'},
    {'1': 'distributing', '3': 7, '4': 1, '5': 8, '10': 'distributing'},
    {'1': 'path', '3': 8, '4': 1, '5': 9, '10': 'path'},
    {'1': 'peers_seeders', '3': 9, '4': 1, '5': 13, '10': 'peers_seeders'},
    {'1': 'peers_half', '3': 10, '4': 1, '5': 13, '10': 'peers_half'},
    {'1': 'peers_available', '3': 11, '4': 1, '5': 13, '10': 'peers_available'},
    {'1': 'completed_at', '3': 12, '4': 1, '5': 9, '10': 'completed_at'},
    {'1': 'verify_at', '3': 13, '4': 1, '5': 9, '10': 'verify_at'},
  ],
};

/// Descriptor for `Download`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadDescriptor = $convert.base64Decode(
    'CghEb3dubG9hZBIiCgVtZWRpYRgBIAEoCzIMLm1lZGlhLk1lZGlhUgVtZWRpYRIUCgVieXRlcx'
    'gCIAEoBFIFYnl0ZXMSHgoKZG93bmxvYWRlZBgDIAEoBFIKZG93bmxvYWRlZBIiCgxpbml0aWF0'
    'ZWRfYXQYBCABKAlSDGluaXRpYXRlZF9hdBIcCglwYXVzZWRfYXQYBSABKAlSCXBhdXNlZF9hdB'
    'IUCgVwZWVycxgGIAEoDVIFcGVlcnMSIgoMZGlzdHJpYnV0aW5nGAcgASgIUgxkaXN0cmlidXRp'
    'bmcSEgoEcGF0aBgIIAEoCVIEcGF0aBIkCg1wZWVyc19zZWVkZXJzGAkgASgNUg1wZWVyc19zZW'
    'VkZXJzEh4KCnBlZXJzX2hhbGYYCiABKA1SCnBlZXJzX2hhbGYSKAoPcGVlcnNfYXZhaWxhYmxl'
    'GAsgASgNUg9wZWVyc19hdmFpbGFibGUSIgoMY29tcGxldGVkX2F0GAwgASgJUgxjb21wbGV0ZW'
    'RfYXQSHAoJdmVyaWZ5X2F0GA0gASgJUgl2ZXJpZnlfYXQ=');

@$core.Deprecated('Use magnetCreateRequestDescriptor instead')
const MagnetCreateRequest$json = {
  '1': 'MagnetCreateRequest',
  '2': [
    {'1': 'uri', '3': 1, '4': 1, '5': 9, '10': 'uri'},
  ],
};

/// Descriptor for `MagnetCreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List magnetCreateRequestDescriptor = $convert
    .base64Decode('ChNNYWduZXRDcmVhdGVSZXF1ZXN0EhAKA3VyaRgBIAEoCVIDdXJp');

@$core.Deprecated('Use magnetCreateResponseDescriptor instead')
const MagnetCreateResponse$json = {
  '1': 'MagnetCreateResponse',
  '2': [
    {
      '1': 'download',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.Download',
      '10': 'download'
    },
  ],
};

/// Descriptor for `MagnetCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List magnetCreateResponseDescriptor = $convert.base64Decode(
    'ChRNYWduZXRDcmVhdGVSZXNwb25zZRIrCghkb3dubG9hZBgBIAEoCzIPLm1lZGlhLkRvd25sb2'
    'FkUghkb3dubG9hZA==');

@$core.Deprecated('Use downloadSearchRequestDescriptor instead')
const DownloadSearchRequest$json = {
  '1': 'DownloadSearchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'mimetypes', '3': 2, '4': 3, '5': 9, '10': 'mimetypes'},
    {'1': 'completed', '3': 3, '4': 1, '5': 8, '10': 'completed'},
    {'1': 'hidden', '3': 4, '4': 1, '5': 8, '10': 'hidden'},
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 5, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `DownloadSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadSearchRequestDescriptor = $convert.base64Decode(
    'ChVEb3dubG9hZFNlYXJjaFJlcXVlc3QSFAoFcXVlcnkYASABKAlSBXF1ZXJ5EhwKCW1pbWV0eX'
    'BlcxgCIAMoCVIJbWltZXR5cGVzEhwKCWNvbXBsZXRlZBgDIAEoCFIJY29tcGxldGVkEhYKBmhp'
    'ZGRlbhgEIAEoCFIGaGlkZGVuEhcKBm9mZnNldBiEByABKARSBm9mZnNldBIVCgVsaW1pdBiFBy'
    'ABKARSBWxpbWl0SgUIBRCEB0oGCIYHEOgH');

@$core.Deprecated('Use downloadSearchResponseDescriptor instead')
const DownloadSearchResponse$json = {
  '1': 'DownloadSearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.DownloadSearchRequest',
      '10': 'next'
    },
    {
      '1': 'items',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.media.Download',
      '10': 'items'
    },
  ],
};

/// Descriptor for `DownloadSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadSearchResponseDescriptor = $convert.base64Decode(
    'ChZEb3dubG9hZFNlYXJjaFJlc3BvbnNlEjAKBG5leHQYASABKAsyHC5tZWRpYS5Eb3dubG9hZF'
    'NlYXJjaFJlcXVlc3RSBG5leHQSJQoFaXRlbXMYAiADKAsyDy5tZWRpYS5Eb3dubG9hZFIFaXRl'
    'bXM=');

@$core.Deprecated('Use downloadUpdateRequestDescriptor instead')
const DownloadUpdateRequest$json = {
  '1': 'DownloadUpdateRequest',
  '2': [
    {
      '1': 'download',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.Download',
      '10': 'download'
    },
  ],
};

/// Descriptor for `DownloadUpdateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadUpdateRequestDescriptor = $convert.base64Decode(
    'ChVEb3dubG9hZFVwZGF0ZVJlcXVlc3QSKwoIZG93bmxvYWQYASABKAsyDy5tZWRpYS5Eb3dubG'
    '9hZFIIZG93bmxvYWQ=');

@$core.Deprecated('Use downloadUpdateResponseDescriptor instead')
const DownloadUpdateResponse$json = {
  '1': 'DownloadUpdateResponse',
  '2': [
    {
      '1': 'download',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.Download',
      '10': 'download'
    },
  ],
};

/// Descriptor for `DownloadUpdateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadUpdateResponseDescriptor =
    $convert.base64Decode(
        'ChZEb3dubG9hZFVwZGF0ZVJlc3BvbnNlEisKCGRvd25sb2FkGAEgASgLMg8ubWVkaWEuRG93bm'
        'xvYWRSCGRvd25sb2Fk');

@$core.Deprecated('Use downloadMetadataRequestDescriptor instead')
const DownloadMetadataRequest$json = {
  '1': 'DownloadMetadataRequest',
};

/// Descriptor for `DownloadMetadataRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadMetadataRequestDescriptor =
    $convert.base64Decode('ChdEb3dubG9hZE1ldGFkYXRhUmVxdWVzdA==');

@$core.Deprecated('Use downloadMetadataResponseDescriptor instead')
const DownloadMetadataResponse$json = {
  '1': 'DownloadMetadataResponse',
  '2': [
    {
      '1': 'download',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.Download',
      '10': 'download'
    },
  ],
};

/// Descriptor for `DownloadMetadataResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadMetadataResponseDescriptor =
    $convert.base64Decode(
        'ChhEb3dubG9hZE1ldGFkYXRhUmVzcG9uc2USKwoIZG93bmxvYWQYASABKAsyDy5tZWRpYS5Eb3'
        'dubG9hZFIIZG93bmxvYWQ=');

@$core.Deprecated('Use downloadBeginRequestDescriptor instead')
const DownloadBeginRequest$json = {
  '1': 'DownloadBeginRequest',
};

/// Descriptor for `DownloadBeginRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadBeginRequestDescriptor =
    $convert.base64Decode('ChREb3dubG9hZEJlZ2luUmVxdWVzdA==');

@$core.Deprecated('Use downloadBeginResponseDescriptor instead')
const DownloadBeginResponse$json = {
  '1': 'DownloadBeginResponse',
  '2': [
    {
      '1': 'download',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.Download',
      '10': 'download'
    },
  ],
};

/// Descriptor for `DownloadBeginResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadBeginResponseDescriptor = $convert.base64Decode(
    'ChVEb3dubG9hZEJlZ2luUmVzcG9uc2USKwoIZG93bmxvYWQYASABKAsyDy5tZWRpYS5Eb3dubG'
    '9hZFIIZG93bmxvYWQ=');

@$core.Deprecated('Use downloadPauseRequestDescriptor instead')
const DownloadPauseRequest$json = {
  '1': 'DownloadPauseRequest',
};

/// Descriptor for `DownloadPauseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadPauseRequestDescriptor =
    $convert.base64Decode('ChREb3dubG9hZFBhdXNlUmVxdWVzdA==');

@$core.Deprecated('Use downloadPauseResponseDescriptor instead')
const DownloadPauseResponse$json = {
  '1': 'DownloadPauseResponse',
  '2': [
    {
      '1': 'download',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.Download',
      '10': 'download'
    },
  ],
};

/// Descriptor for `DownloadPauseResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadPauseResponseDescriptor = $convert.base64Decode(
    'ChVEb3dubG9hZFBhdXNlUmVzcG9uc2USKwoIZG93bmxvYWQYASABKAsyDy5tZWRpYS5Eb3dubG'
    '9hZFIIZG93bmxvYWQ=');

@$core.Deprecated('Use downloadTuneRequestDescriptor instead')
const DownloadTuneRequest$json = {
  '1': 'DownloadTuneRequest',
  '2': [
    {'1': 'peers', '3': 1, '4': 3, '5': 9, '10': 'peers'},
  ],
};

/// Descriptor for `DownloadTuneRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadTuneRequestDescriptor =
    $convert.base64Decode(
        'ChNEb3dubG9hZFR1bmVSZXF1ZXN0EhQKBXBlZXJzGAEgAygJUgVwZWVycw==');

@$core.Deprecated('Use downloadTuneResponseDescriptor instead')
const DownloadTuneResponse$json = {
  '1': 'DownloadTuneResponse',
};

/// Descriptor for `DownloadTuneResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadTuneResponseDescriptor =
    $convert.base64Decode('ChREb3dubG9hZFR1bmVSZXNwb25zZQ==');

@$core.Deprecated('Use downloadDeleteRequestDescriptor instead')
const DownloadDeleteRequest$json = {
  '1': 'DownloadDeleteRequest',
};

/// Descriptor for `DownloadDeleteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadDeleteRequestDescriptor =
    $convert.base64Decode('ChVEb3dubG9hZERlbGV0ZVJlcXVlc3Q=');

@$core.Deprecated('Use downloadDeleteResponseDescriptor instead')
const DownloadDeleteResponse$json = {
  '1': 'DownloadDeleteResponse',
  '2': [
    {
      '1': 'download',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.Download',
      '10': 'download'
    },
  ],
};

/// Descriptor for `DownloadDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadDeleteResponseDescriptor =
    $convert.base64Decode(
        'ChZEb3dubG9hZERlbGV0ZVJlc3BvbnNlEisKCGRvd25sb2FkGAEgASgLMg8ubWVkaWEuRG93bm'
        'xvYWRSCGRvd25sb2Fk');

@$core.Deprecated('Use publishedDescriptor instead')
const Published$json = {
  '1': 'Published',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'mimetype', '3': 2, '4': 1, '5': 9, '10': 'mimetype'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
    {'1': 'bytes', '3': 4, '4': 1, '5': 4, '10': 'bytes'},
    {'1': 'entropy', '3': 5, '4': 1, '5': 9, '10': 'entropy'},
    {'1': 'expires_at', '3': 6, '4': 1, '5': 9, '10': 'expires_at'},
  ],
};

/// Descriptor for `Published`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishedDescriptor = $convert.base64Decode(
    'CglQdWJsaXNoZWQSDgoCaWQYASABKAlSAmlkEhoKCG1pbWV0eXBlGAIgASgJUghtaW1ldHlwZR'
    'IgCgtkZXNjcmlwdGlvbhgDIAEoCVILZGVzY3JpcHRpb24SFAoFYnl0ZXMYBCABKARSBWJ5dGVz'
    'EhgKB2VudHJvcHkYBSABKAlSB2VudHJvcHkSHgoKZXhwaXJlc19hdBgGIAEoCVIKZXhwaXJlc1'
    '9hdA==');

@$core.Deprecated('Use publishedUploadRequestDescriptor instead')
const PublishedUploadRequest$json = {
  '1': 'PublishedUploadRequest',
  '2': [
    {'1': 'entropy', '3': 1, '4': 1, '5': 9, '10': 'entropy'},
    {'1': 'mimetype', '3': 2, '4': 1, '5': 9, '10': 'mimetype'},
    {'1': 'ttl', '3': 3, '4': 1, '5': 4, '10': 'ttl'},
  ],
};

/// Descriptor for `PublishedUploadRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishedUploadRequestDescriptor =
    $convert.base64Decode(
        'ChZQdWJsaXNoZWRVcGxvYWRSZXF1ZXN0EhgKB2VudHJvcHkYASABKAlSB2VudHJvcHkSGgoIbW'
        'ltZXR5cGUYAiABKAlSCG1pbWV0eXBlEhAKA3R0bBgDIAEoBFIDdHRs');

@$core.Deprecated('Use publishedUploadResponseDescriptor instead')
const PublishedUploadResponse$json = {
  '1': 'PublishedUploadResponse',
  '2': [
    {
      '1': 'published',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.Published',
      '10': 'published'
    },
  ],
};

/// Descriptor for `PublishedUploadResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishedUploadResponseDescriptor =
    $convert.base64Decode(
        'ChdQdWJsaXNoZWRVcGxvYWRSZXNwb25zZRIuCglwdWJsaXNoZWQYASABKAsyEC5tZWRpYS5QdW'
        'JsaXNoZWRSCXB1Ymxpc2hlZA==');

@$core.Deprecated('Use metadataSyncRequestDescriptor instead')
const MetadataSyncRequest$json = {
  '1': 'MetadataSyncRequest',
  '2': [
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
  ],
};

/// Descriptor for `MetadataSyncRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metadataSyncRequestDescriptor = $convert.base64Decode(
    'ChNNZXRhZGF0YVN5bmNSZXF1ZXN0EiIKBW1lZGlhGAEgASgLMgwubWVkaWEuTWVkaWFSBW1lZG'
    'lh');

@$core.Deprecated('Use metadataSyncResponseDescriptor instead')
const MetadataSyncResponse$json = {
  '1': 'MetadataSyncResponse',
  '2': [
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
  ],
};

/// Descriptor for `MetadataSyncResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metadataSyncResponseDescriptor = $convert.base64Decode(
    'ChRNZXRhZGF0YVN5bmNSZXNwb25zZRIiCgVtZWRpYRgBIAEoCzIMLm1lZGlhLk1lZGlhUgVtZW'
    'RpYQ==');

@$core.Deprecated('Use publishedRequestDescriptor instead')
const PublishedRequest$json = {
  '1': 'PublishedRequest',
  '2': [
    {'1': 'known_media_id', '3': 1, '4': 1, '5': 9, '10': 'known_media_id'},
  ],
};

/// Descriptor for `PublishedRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishedRequestDescriptor = $convert.base64Decode(
    'ChBQdWJsaXNoZWRSZXF1ZXN0EiYKDmtub3duX21lZGlhX2lkGAEgASgJUg5rbm93bl9tZWRpYV'
    '9pZA==');

@$core.Deprecated('Use publishedResponseDescriptor instead')
const PublishedResponse$json = {
  '1': 'PublishedResponse',
  '2': [
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
    {'1': 'magnet_uri', '3': 2, '4': 1, '5': 9, '10': 'magnet_uri'},
  ],
};

/// Descriptor for `PublishedResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishedResponseDescriptor = $convert.base64Decode(
    'ChFQdWJsaXNoZWRSZXNwb25zZRIiCgVtZWRpYRgBIAEoCzIMLm1lZGlhLk1lZGlhUgVtZWRpYR'
    'IeCgptYWduZXRfdXJpGAIgASgJUgptYWduZXRfdXJp');
