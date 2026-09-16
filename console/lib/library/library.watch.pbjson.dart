// This is a generated file - do not edit.
//
// Generated from library/library.watch.proto.

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

@$core.Deprecated('Use watchHistoryRecordDescriptor instead')
const WatchHistoryRecord$json = {
  '1': 'WatchHistoryRecord',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'media_id', '3': 2, '4': 1, '5': 9, '10': 'media_id'},
    {'1': 'duration', '3': 3, '4': 1, '5': 4, '10': 'duration'},
  ],
};

/// Descriptor for `WatchHistoryRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List watchHistoryRecordDescriptor = $convert.base64Decode(
    'ChJXYXRjaEhpc3RvcnlSZWNvcmQSDgoCaWQYASABKAlSAmlkEhoKCG1lZGlhX2lkGAIgASgJUg'
    'htZWRpYV9pZBIaCghkdXJhdGlvbhgDIAEoBFIIZHVyYXRpb24=');

@$core.Deprecated('Use watchHistoryRecordRequestDescriptor instead')
const WatchHistoryRecordRequest$json = {
  '1': 'WatchHistoryRecordRequest',
  '2': [
    {
      '1': 'record',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.library.WatchHistoryRecord',
      '10': 'record'
    },
  ],
};

/// Descriptor for `WatchHistoryRecordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List watchHistoryRecordRequestDescriptor =
    $convert.base64Decode(
        'ChlXYXRjaEhpc3RvcnlSZWNvcmRSZXF1ZXN0EjMKBnJlY29yZBgBIAEoCzIbLmxpYnJhcnkuV2'
        'F0Y2hIaXN0b3J5UmVjb3JkUgZyZWNvcmQ=');

@$core.Deprecated('Use watchHistoryRecordResponseDescriptor instead')
const WatchHistoryRecordResponse$json = {
  '1': 'WatchHistoryRecordResponse',
};

/// Descriptor for `WatchHistoryRecordResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List watchHistoryRecordResponseDescriptor =
    $convert.base64Decode('ChpXYXRjaEhpc3RvcnlSZWNvcmRSZXNwb25zZQ==');
