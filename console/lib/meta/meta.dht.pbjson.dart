// This is a generated file - do not edit.
//
// Generated from meta.dht.proto.

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

@$core.Deprecated('Use dHTDiagnosticsDescriptor instead')
const DHTDiagnostics$json = {
  '1': 'DHTDiagnostics',
  '2': [
    {'1': 'node_id', '3': 1, '4': 1, '5': 9, '10': 'node_id'},
    {'1': 'good_nodes', '3': 2, '4': 1, '5': 5, '10': 'good_nodes'},
    {'1': 'nodes', '3': 3, '4': 1, '5': 5, '10': 'nodes'},
    {
      '1': 'outstanding_transactions',
      '3': 4,
      '4': 1,
      '5': 5,
      '10': 'outstanding_transactions'
    },
    {
      '1': 'successful_outbound_announce_peer_queries',
      '3': 5,
      '4': 1,
      '5': 3,
      '10': 'successful_outbound_announce_peer_queries'
    },
    {'1': 'bad_nodes', '3': 6, '4': 1, '5': 13, '10': 'bad_nodes'},
    {
      '1': 'outbound_queries_attempted',
      '3': 7,
      '4': 1,
      '5': 3,
      '10': 'outbound_queries_attempted'
    },
  ],
};

/// Descriptor for `DHTDiagnostics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dHTDiagnosticsDescriptor = $convert.base64Decode(
    'Cg5ESFREaWFnbm9zdGljcxIYCgdub2RlX2lkGAEgASgJUgdub2RlX2lkEh4KCmdvb2Rfbm9kZX'
    'MYAiABKAVSCmdvb2Rfbm9kZXMSFAoFbm9kZXMYAyABKAVSBW5vZGVzEjoKGG91dHN0YW5kaW5n'
    'X3RyYW5zYWN0aW9ucxgEIAEoBVIYb3V0c3RhbmRpbmdfdHJhbnNhY3Rpb25zElwKKXN1Y2Nlc3'
    'NmdWxfb3V0Ym91bmRfYW5ub3VuY2VfcGVlcl9xdWVyaWVzGAUgASgDUilzdWNjZXNzZnVsX291'
    'dGJvdW5kX2Fubm91bmNlX3BlZXJfcXVlcmllcxIcCgliYWRfbm9kZXMYBiABKA1SCWJhZF9ub2'
    'RlcxI+ChpvdXRib3VuZF9xdWVyaWVzX2F0dGVtcHRlZBgHIAEoA1Iab3V0Ym91bmRfcXVlcmll'
    'c19hdHRlbXB0ZWQ=');

@$core.Deprecated('Use dHTMetricsResponseDescriptor instead')
const DHTMetricsResponse$json = {
  '1': 'DHTMetricsResponse',
  '2': [
    {
      '1': 'dht',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.DHTDiagnostics',
      '10': 'dht'
    },
  ],
};

/// Descriptor for `DHTMetricsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dHTMetricsResponseDescriptor = $convert.base64Decode(
    'ChJESFRNZXRyaWNzUmVzcG9uc2USJgoDZGh0GAEgASgLMhQubWV0YS5ESFREaWFnbm9zdGljc1'
    'IDZGh0');
