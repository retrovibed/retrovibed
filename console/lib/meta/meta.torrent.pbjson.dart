// This is a generated file - do not edit.
//
// Generated from meta/meta.torrent.proto.

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

@$core.Deprecated('Use torrentDiagnosticsDescriptor instead')
const TorrentDiagnostics$json = {
  '1': 'TorrentDiagnostics',
  '2': [
    {'1': 'total', '3': 1, '4': 1, '5': 4, '10': 'total'},
    {'1': 'seeding', '3': 2, '4': 1, '5': 4, '10': 'seeding'},
    {'1': 'bytes', '3': 3, '4': 1, '5': 4, '10': 'bytes'},
    {'1': 'downloaded', '3': 4, '4': 1, '5': 4, '10': 'downloaded'},
    {'1': 'uploaded', '3': 5, '4': 1, '5': 4, '10': 'uploaded'},
    {'1': 'peers', '3': 6, '4': 1, '5': 4, '10': 'peers'},
  ],
};

/// Descriptor for `TorrentDiagnostics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List torrentDiagnosticsDescriptor = $convert.base64Decode(
    'ChJUb3JyZW50RGlhZ25vc3RpY3MSFAoFdG90YWwYASABKARSBXRvdGFsEhgKB3NlZWRpbmcYAi'
    'ABKARSB3NlZWRpbmcSFAoFYnl0ZXMYAyABKARSBWJ5dGVzEh4KCmRvd25sb2FkZWQYBCABKARS'
    'CmRvd25sb2FkZWQSGgoIdXBsb2FkZWQYBSABKARSCHVwbG9hZGVkEhQKBXBlZXJzGAYgASgEUg'
    'VwZWVycw==');

@$core.Deprecated('Use torrentMetricsResponseDescriptor instead')
const TorrentMetricsResponse$json = {
  '1': 'TorrentMetricsResponse',
  '2': [
    {
      '1': 'torrent',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.TorrentDiagnostics',
      '10': 'torrent'
    },
  ],
};

/// Descriptor for `TorrentMetricsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List torrentMetricsResponseDescriptor =
    $convert.base64Decode(
        'ChZUb3JyZW50TWV0cmljc1Jlc3BvbnNlEjIKB3RvcnJlbnQYASABKAsyGC5tZXRhLlRvcnJlbn'
        'REaWFnbm9zdGljc1IHdG9ycmVudA==');

@$core.Deprecated('Use torrentInfoResponseDescriptor instead')
const TorrentInfoResponse$json = {
  '1': 'TorrentInfoResponse',
  '2': [
    {
      '1': 'meta',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.TorrentMeta',
      '10': 'meta'
    },
    {
      '1': 'details',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.meta.TorrentDetails',
      '10': 'details'
    },
    {
      '1': 'files',
      '3': 1000,
      '4': 3,
      '5': 11,
      '6': '.meta.TorrentFile',
      '10': 'files'
    },
  ],
};

/// Descriptor for `TorrentInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List torrentInfoResponseDescriptor = $convert.base64Decode(
    'ChNUb3JyZW50SW5mb1Jlc3BvbnNlEiUKBG1ldGEYASABKAsyES5tZXRhLlRvcnJlbnRNZXRhUg'
    'RtZXRhEi4KB2RldGFpbHMYAiABKAsyFC5tZXRhLlRvcnJlbnREZXRhaWxzUgdkZXRhaWxzEigK'
    'BWZpbGVzGOgHIAMoCzIRLm1ldGEuVG9ycmVudEZpbGVSBWZpbGVz');

@$core.Deprecated('Use torrentFileDescriptor instead')
const TorrentFile$json = {
  '1': 'TorrentFile',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'length', '3': 2, '4': 1, '5': 4, '10': 'length'},
    {'1': 'path', '3': 3, '4': 1, '5': 9, '10': 'path'},
  ],
};

/// Descriptor for `TorrentFile`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List torrentFileDescriptor = $convert.base64Decode(
    'CgtUb3JyZW50RmlsZRISCgRuYW1lGAEgASgJUgRuYW1lEhYKBmxlbmd0aBgCIAEoBFIGbGVuZ3'
    'RoEhIKBHBhdGgYAyABKAlSBHBhdGg=');

@$core.Deprecated('Use torrentMetaDescriptor instead')
const TorrentMeta$json = {
  '1': 'TorrentMeta',
  '2': [
    {'1': 'comment', '3': 1, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'encoding', '3': 2, '4': 1, '5': 9, '10': 'encoding'},
    {'1': 'created_by', '3': 3, '4': 1, '5': 9, '10': 'created_by'},
    {'1': 'creation_date', '3': 4, '4': 1, '5': 4, '10': 'creation_date'},
    {'1': 'announce_list', '3': 1000, '4': 3, '5': 9, '10': 'announce_list'},
    {'1': 'url_list', '3': 1001, '4': 3, '5': 9, '10': 'url_list'},
  ],
};

/// Descriptor for `TorrentMeta`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List torrentMetaDescriptor = $convert.base64Decode(
    'CgtUb3JyZW50TWV0YRIYCgdjb21tZW50GAEgASgJUgdjb21tZW50EhoKCGVuY29kaW5nGAIgAS'
    'gJUghlbmNvZGluZxIeCgpjcmVhdGVkX2J5GAMgASgJUgpjcmVhdGVkX2J5EiQKDWNyZWF0aW9u'
    'X2RhdGUYBCABKARSDWNyZWF0aW9uX2RhdGUSJQoNYW5ub3VuY2VfbGlzdBjoByADKAlSDWFubm'
    '91bmNlX2xpc3QSGwoIdXJsX2xpc3QY6QcgAygJUgh1cmxfbGlzdA==');

@$core.Deprecated('Use torrentDetailsDescriptor instead')
const TorrentDetails$json = {
  '1': 'TorrentDetails',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'length', '3': 2, '4': 1, '5': 4, '10': 'length'},
    {'1': 'source', '3': 3, '4': 1, '5': 9, '10': 'source'},
    {'1': 'private', '3': 4, '4': 1, '5': 8, '10': 'private'},
    {'1': 'downloaded', '3': 5, '4': 1, '5': 4, '10': 'downloaded'},
    {'1': 'uploaded', '3': 6, '4': 1, '5': 4, '10': 'uploaded'},
  ],
};

/// Descriptor for `TorrentDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List torrentDetailsDescriptor = $convert.base64Decode(
    'Cg5Ub3JyZW50RGV0YWlscxISCgRuYW1lGAEgASgJUgRuYW1lEhYKBmxlbmd0aBgCIAEoBFIGbG'
    'VuZ3RoEhYKBnNvdXJjZRgDIAEoCVIGc291cmNlEhgKB3ByaXZhdGUYBCABKAhSB3ByaXZhdGUS'
    'HgoKZG93bmxvYWRlZBgFIAEoBFIKZG93bmxvYWRlZBIaCgh1cGxvYWRlZBgGIAEoBFIIdXBsb2'
    'FkZWQ=');
