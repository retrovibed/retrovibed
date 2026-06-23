// This is a generated file - do not edit.
//
// Generated from meta.torrent.proto.

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
