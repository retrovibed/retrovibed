// This is a generated file - do not edit.
//
// Generated from community/community.metrics.proto.

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
