//
//  Generated code. Do not modify.
//  source: media.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

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
  ],
};

/// Descriptor for `Media`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaDescriptor = $convert.base64Decode(
    'CgVNZWRpYRIOCgJpZBgBIAEoCVICaWQSIAoLZGVzY3JpcHRpb24YAiABKAlSC2Rlc2NyaXB0aW'
    '9uEhoKCG1pbWV0eXBlGAMgASgJUghtaW1ldHlwZRIUCgVpbWFnZRgEIAEoCVIFaW1hZ2U=');

@$core.Deprecated('Use mediaSearchRequestDescriptor instead')
const MediaSearchRequest$json = {
  '1': 'MediaSearchRequest',
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

/// Descriptor for `MediaSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaSearchRequestDescriptor = $convert.base64Decode(
    'ChJNZWRpYVNlYXJjaFJlcXVlc3QSFAoFcXVlcnkYASABKAlSBXF1ZXJ5EhcKBm9mZnNldBiEBy'
    'ABKARSBm9mZnNldBIVCgVsaW1pdBiFByABKARSBWxpbWl0SgUIAhCEB0oGCIYHEOgH');

@$core.Deprecated('Use mediaSearchResponseDescriptor instead')
const MediaSearchResponse$json = {
  '1': 'MediaSearchResponse',
  '2': [
    {'1': 'next', '3': 1, '4': 1, '5': 11, '6': '.media.MediaSearchRequest', '10': 'next'},
    {'1': 'items', '3': 2, '4': 3, '5': 11, '6': '.media.Media', '10': 'items'},
  ],
};

/// Descriptor for `MediaSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaSearchResponseDescriptor = $convert.base64Decode(
    'ChNNZWRpYVNlYXJjaFJlc3BvbnNlEi0KBG5leHQYASABKAsyGS5tZWRpYS5NZWRpYVNlYXJjaF'
    'JlcXVlc3RSBG5leHQSIgoFaXRlbXMYAiADKAsyDC5tZWRpYS5NZWRpYVIFaXRlbXM=');

@$core.Deprecated('Use downloadDescriptor instead')
const Download$json = {
  '1': 'Download',
  '2': [
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
    {'1': 'bytes', '3': 2, '4': 1, '5': 4, '10': 'bytes'},
    {'1': 'downloaded', '3': 3, '4': 1, '5': 4, '10': 'downloaded'},
    {'1': 'initiated_at', '3': 4, '4': 1, '5': 9, '10': 'initiated_at'},
    {'1': 'paused_at', '3': 5, '4': 1, '5': 9, '10': 'paused_at'},
  ],
};

/// Descriptor for `Download`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadDescriptor = $convert.base64Decode(
    'CghEb3dubG9hZBIiCgVtZWRpYRgBIAEoCzIMLm1lZGlhLk1lZGlhUgVtZWRpYRIUCgVieXRlcx'
    'gCIAEoBFIFYnl0ZXMSHgoKZG93bmxvYWRlZBgDIAEoBFIKZG93bmxvYWRlZBIiCgxpbml0aWF0'
    'ZWRfYXQYBCABKAlSDGluaXRpYXRlZF9hdBIcCglwYXVzZWRfYXQYBSABKAlSCXBhdXNlZF9hdA'
    '==');

@$core.Deprecated('Use downloadSearchRequestDescriptor instead')
const DownloadSearchRequest$json = {
  '1': 'DownloadSearchRequest',
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

/// Descriptor for `DownloadSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadSearchRequestDescriptor = $convert.base64Decode(
    'ChVEb3dubG9hZFNlYXJjaFJlcXVlc3QSFAoFcXVlcnkYASABKAlSBXF1ZXJ5EhcKBm9mZnNldB'
    'iEByABKARSBm9mZnNldBIVCgVsaW1pdBiFByABKARSBWxpbWl0SgUIAhCEB0oGCIYHEOgH');

@$core.Deprecated('Use downloadSearchResponseDescriptor instead')
const DownloadSearchResponse$json = {
  '1': 'DownloadSearchResponse',
  '2': [
    {'1': 'next', '3': 1, '4': 1, '5': 11, '6': '.media.DownloadSearchRequest', '10': 'next'},
    {'1': 'items', '3': 2, '4': 3, '5': 11, '6': '.media.Download', '10': 'items'},
  ],
};

/// Descriptor for `DownloadSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadSearchResponseDescriptor = $convert.base64Decode(
    'ChZEb3dubG9hZFNlYXJjaFJlc3BvbnNlEjAKBG5leHQYASABKAsyHC5tZWRpYS5Eb3dubG9hZF'
    'NlYXJjaFJlcXVlc3RSBG5leHQSJQoFaXRlbXMYAiADKAsyDy5tZWRpYS5Eb3dubG9hZFIFaXRl'
    'bXM=');

@$core.Deprecated('Use downloadBeginRequestDescriptor instead')
const DownloadBeginRequest$json = {
  '1': 'DownloadBeginRequest',
};

/// Descriptor for `DownloadBeginRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadBeginRequestDescriptor = $convert.base64Decode(
    'ChREb3dubG9hZEJlZ2luUmVxdWVzdA==');

@$core.Deprecated('Use downloadBeginResponseDescriptor instead')
const DownloadBeginResponse$json = {
  '1': 'DownloadBeginResponse',
  '2': [
    {'1': 'download', '3': 1, '4': 1, '5': 11, '6': '.media.Download', '10': 'download'},
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
final $typed_data.Uint8List downloadPauseRequestDescriptor = $convert.base64Decode(
    'ChREb3dubG9hZFBhdXNlUmVxdWVzdA==');

@$core.Deprecated('Use downloadPauseResponseDescriptor instead')
const DownloadPauseResponse$json = {
  '1': 'DownloadPauseResponse',
  '2': [
    {'1': 'download', '3': 1, '4': 1, '5': 11, '6': '.media.Download', '10': 'download'},
  ],
};

/// Descriptor for `DownloadPauseResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List downloadPauseResponseDescriptor = $convert.base64Decode(
    'ChVEb3dubG9hZFBhdXNlUmVzcG9uc2USKwoIZG93bmxvYWQYASABKAsyDy5tZWRpYS5Eb3dubG'
    '9hZFIIZG93bmxvYWQ=');

