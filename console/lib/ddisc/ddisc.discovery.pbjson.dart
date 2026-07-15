// This is a generated file - do not edit.
//
// Generated from ddisc.discovery.proto.

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

@$core.Deprecated('Use discoveryDescriptor instead')
const Discovery$json = {
  '1': 'Discovery',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'infohash', '3': 2, '4': 1, '5': 12, '10': 'infohash'},
    {'1': 'attempts', '3': 3, '4': 1, '5': 13, '10': 'attempts'},
    {'1': 'next_check', '3': 4, '4': 1, '5': 9, '10': 'next_check'},
    {'1': 'created_at', '3': 5, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 6, '4': 1, '5': 9, '10': 'updated_at'},
  ],
};

/// Descriptor for `Discovery`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List discoveryDescriptor = $convert.base64Decode(
    'CglEaXNjb3ZlcnkSDgoCaWQYASABKAlSAmlkEhoKCGluZm9oYXNoGAIgASgMUghpbmZvaGFzaB'
    'IaCghhdHRlbXB0cxgDIAEoDVIIYXR0ZW1wdHMSHgoKbmV4dF9jaGVjaxgEIAEoCVIKbmV4dF9j'
    'aGVjaxIeCgpjcmVhdGVkX2F0GAUgASgJUgpjcmVhdGVkX2F0Eh4KCnVwZGF0ZWRfYXQYBiABKA'
    'lSCnVwZGF0ZWRfYXQ=');

@$core.Deprecated('Use discoverySearchRequestDescriptor instead')
const DiscoverySearchRequest$json = {
  '1': 'DiscoverySearchRequest',
  '2': [
    {
      '1': 'next_check',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.DateRange',
      '10': 'next_check'
    },
    {'1': 'id', '3': 2, '4': 3, '5': 9, '10': 'id'},
    {'1': 'attempts_min', '3': 3, '4': 1, '5': 4, '10': 'attempts_min'},
    {'1': 'attempts_max', '3': 4, '4': 1, '5': 4, '10': 'attempts_max'},
    {'1': 'offset', '3': 1000, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 1001, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 5, '2': 1000},
  ],
};

/// Descriptor for `DiscoverySearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List discoverySearchRequestDescriptor = $convert.base64Decode(
    'ChZEaXNjb3ZlcnlTZWFyY2hSZXF1ZXN0Ei8KCm5leHRfY2hlY2sYASABKAsyDy5tZXRhLkRhdG'
    'VSYW5nZVIKbmV4dF9jaGVjaxIOCgJpZBgCIAMoCVICaWQSIgoMYXR0ZW1wdHNfbWluGAMgASgE'
    'UgxhdHRlbXB0c19taW4SIgoMYXR0ZW1wdHNfbWF4GAQgASgEUgxhdHRlbXB0c19tYXgSFwoGb2'
    'Zmc2V0GOgHIAEoBFIGb2Zmc2V0EhUKBWxpbWl0GOkHIAEoBFIFbGltaXRKBQgFEOgH');

@$core.Deprecated('Use discoverySearchResponseDescriptor instead')
const DiscoverySearchResponse$json = {
  '1': 'DiscoverySearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.ddisc.DiscoverySearchRequest',
      '10': 'next'
    },
    {
      '1': 'items',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.ddisc.Discovery',
      '10': 'items'
    },
  ],
};

/// Descriptor for `DiscoverySearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List discoverySearchResponseDescriptor = $convert.base64Decode(
    'ChdEaXNjb3ZlcnlTZWFyY2hSZXNwb25zZRIxCgRuZXh0GAEgASgLMh0uZGRpc2MuRGlzY292ZX'
    'J5U2VhcmNoUmVxdWVzdFIEbmV4dBImCgVpdGVtcxgCIAMoCzIQLmRkaXNjLkRpc2NvdmVyeVIF'
    'aXRlbXM=');

@$core.Deprecated('Use discoveryCreateRequestDescriptor instead')
const DiscoveryCreateRequest$json = {
  '1': 'DiscoveryCreateRequest',
  '2': [
    {
      '1': 'discovery',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.ddisc.Discovery',
      '10': 'discovery'
    },
  ],
};

/// Descriptor for `DiscoveryCreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List discoveryCreateRequestDescriptor =
    $convert.base64Decode(
        'ChZEaXNjb3ZlcnlDcmVhdGVSZXF1ZXN0Ei4KCWRpc2NvdmVyeRgBIAEoCzIQLmRkaXNjLkRpc2'
        'NvdmVyeVIJZGlzY292ZXJ5');

@$core.Deprecated('Use discoveryCreateResponseDescriptor instead')
const DiscoveryCreateResponse$json = {
  '1': 'DiscoveryCreateResponse',
  '2': [
    {
      '1': 'discovery',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.ddisc.Discovery',
      '10': 'discovery'
    },
  ],
};

/// Descriptor for `DiscoveryCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List discoveryCreateResponseDescriptor =
    $convert.base64Decode(
        'ChdEaXNjb3ZlcnlDcmVhdGVSZXNwb25zZRIuCglkaXNjb3ZlcnkYASABKAsyEC5kZGlzYy5EaX'
        'Njb3ZlcnlSCWRpc2NvdmVyeQ==');

@$core.Deprecated('Use discoveryDownloadRequestDescriptor instead')
const DiscoveryDownloadRequest$json = {
  '1': 'DiscoveryDownloadRequest',
};

/// Descriptor for `DiscoveryDownloadRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List discoveryDownloadRequestDescriptor =
    $convert.base64Decode('ChhEaXNjb3ZlcnlEb3dubG9hZFJlcXVlc3Q=');

@$core.Deprecated('Use discoveryDownloadResponseDescriptor instead')
const DiscoveryDownloadResponse$json = {
  '1': 'DiscoveryDownloadResponse',
  '2': [
    {
      '1': 'discovery',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.ddisc.Discovery',
      '10': 'discovery'
    },
  ],
};

/// Descriptor for `DiscoveryDownloadResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List discoveryDownloadResponseDescriptor =
    $convert.base64Decode(
        'ChlEaXNjb3ZlcnlEb3dubG9hZFJlc3BvbnNlEi4KCWRpc2NvdmVyeRgBIAEoCzIQLmRkaXNjLk'
        'Rpc2NvdmVyeVIJZGlzY292ZXJ5');

@$core.Deprecated('Use discoveryDeleteRequestDescriptor instead')
const DiscoveryDeleteRequest$json = {
  '1': 'DiscoveryDeleteRequest',
};

/// Descriptor for `DiscoveryDeleteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List discoveryDeleteRequestDescriptor =
    $convert.base64Decode('ChZEaXNjb3ZlcnlEZWxldGVSZXF1ZXN0');

@$core.Deprecated('Use discoveryDeleteResponseDescriptor instead')
const DiscoveryDeleteResponse$json = {
  '1': 'DiscoveryDeleteResponse',
  '2': [
    {
      '1': 'discovery',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.ddisc.Discovery',
      '10': 'discovery'
    },
  ],
};

/// Descriptor for `DiscoveryDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List discoveryDeleteResponseDescriptor =
    $convert.base64Decode(
        'ChdEaXNjb3ZlcnlEZWxldGVSZXNwb25zZRIuCglkaXNjb3ZlcnkYASABKAsyEC5kZGlzYy5EaX'
        'Njb3ZlcnlSCWRpc2NvdmVyeQ==');
