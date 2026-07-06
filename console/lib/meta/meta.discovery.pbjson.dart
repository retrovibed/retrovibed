// This is a generated file - do not edit.
//
// Generated from meta.discovery.proto.

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

@$core.Deprecated('Use discoveryDiagnosticsDescriptor instead')
const DiscoveryDiagnostics$json = {
  '1': 'DiscoveryDiagnostics',
  '2': [
    {'1': 'enabled', '3': 1, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'ratio', '3': 2, '4': 1, '5': 13, '10': 'ratio'},
    {'1': 'partitions', '3': 3, '4': 1, '5': 13, '10': 'partitions'},
    {'1': 'workloads', '3': 4, '4': 1, '5': 13, '10': 'workloads'},
    {'1': 'local_partition', '3': 5, '4': 1, '5': 9, '10': 'local_partition'},
    {'1': 'peers', '3': 6, '4': 1, '5': 4, '10': 'peers'},
    {'1': 'peers_ddisc', '3': 7, '4': 1, '5': 4, '10': 'peers_ddisc'},
    {'1': 'peers_bep51', '3': 8, '4': 1, '5': 4, '10': 'peers_bep51'},
    {'1': 'unidentified', '3': 9, '4': 1, '5': 4, '10': 'unidentified'},
    {'1': 'queued', '3': 10, '4': 1, '5': 4, '10': 'queued'},
    {'1': 'indexing', '3': 11, '4': 1, '5': 4, '10': 'indexing'},
    {'1': 'offload', '3': 12, '4': 1, '5': 4, '10': 'offload'},
    {'1': 'indexed', '3': 13, '4': 1, '5': 4, '10': 'indexed'},
  ],
};

/// Descriptor for `DiscoveryDiagnostics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List discoveryDiagnosticsDescriptor = $convert.base64Decode(
    'ChREaXNjb3ZlcnlEaWFnbm9zdGljcxIYCgdlbmFibGVkGAEgASgIUgdlbmFibGVkEhQKBXJhdG'
    'lvGAIgASgNUgVyYXRpbxIeCgpwYXJ0aXRpb25zGAMgASgNUgpwYXJ0aXRpb25zEhwKCXdvcmts'
    'b2FkcxgEIAEoDVIJd29ya2xvYWRzEigKD2xvY2FsX3BhcnRpdGlvbhgFIAEoCVIPbG9jYWxfcG'
    'FydGl0aW9uEhQKBXBlZXJzGAYgASgEUgVwZWVycxIgCgtwZWVyc19kZGlzYxgHIAEoBFILcGVl'
    'cnNfZGRpc2MSIAoLcGVlcnNfYmVwNTEYCCABKARSC3BlZXJzX2JlcDUxEiIKDHVuaWRlbnRpZm'
    'llZBgJIAEoBFIMdW5pZGVudGlmaWVkEhYKBnF1ZXVlZBgKIAEoBFIGcXVldWVkEhoKCGluZGV4'
    'aW5nGAsgASgEUghpbmRleGluZxIYCgdvZmZsb2FkGAwgASgEUgdvZmZsb2FkEhgKB2luZGV4ZW'
    'QYDSABKARSB2luZGV4ZWQ=');

@$core.Deprecated('Use discoveryMetricsResponseDescriptor instead')
const DiscoveryMetricsResponse$json = {
  '1': 'DiscoveryMetricsResponse',
  '2': [
    {
      '1': 'discovery',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.DiscoveryDiagnostics',
      '10': 'discovery'
    },
  ],
};

/// Descriptor for `DiscoveryMetricsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List discoveryMetricsResponseDescriptor =
    $convert.base64Decode(
        'ChhEaXNjb3ZlcnlNZXRyaWNzUmVzcG9uc2USOAoJZGlzY292ZXJ5GAEgASgLMhoubWV0YS5EaX'
        'Njb3ZlcnlEaWFnbm9zdGljc1IJZGlzY292ZXJ5');
