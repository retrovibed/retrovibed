// This is a generated file - do not edit.
//
// Generated from searchplugin.proto.

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

@$core.Deprecated('Use pluginDescriptor instead')
const Plugin$json = {
  '1': 'Plugin',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'size', '3': 3, '4': 1, '5': 4, '10': 'size'},
    {'1': 'installed_at', '3': 4, '4': 1, '5': 9, '10': 'installed_at'},
  ],
};

/// Descriptor for `Plugin`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pluginDescriptor = $convert.base64Decode(
    'CgZQbHVnaW4SDgoCaWQYASABKAlSAmlkEhIKBG5hbWUYAiABKAlSBG5hbWUSEgoEc2l6ZRgDIA'
    'EoBFIEc2l6ZRIiCgxpbnN0YWxsZWRfYXQYBCABKAlSDGluc3RhbGxlZF9hdA==');

@$core.Deprecated('Use pluginSearchRequestDescriptor instead')
const PluginSearchRequest$json = {
  '1': 'PluginSearchRequest',
  '2': [
    {'1': 'offset', '3': 1, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 2, '4': 1, '5': 4, '10': 'limit'},
  ],
};

/// Descriptor for `PluginSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pluginSearchRequestDescriptor = $convert.base64Decode(
    'ChNQbHVnaW5TZWFyY2hSZXF1ZXN0EhYKBm9mZnNldBgBIAEoBFIGb2Zmc2V0EhQKBWxpbWl0GA'
    'IgASgEUgVsaW1pdA==');

@$core.Deprecated('Use pluginSearchResponseDescriptor instead')
const PluginSearchResponse$json = {
  '1': 'PluginSearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.searchplugin.PluginSearchRequest',
      '10': 'next'
    },
    {
      '1': 'items',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.searchplugin.Plugin',
      '10': 'items'
    },
  ],
};

/// Descriptor for `PluginSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pluginSearchResponseDescriptor = $convert.base64Decode(
    'ChRQbHVnaW5TZWFyY2hSZXNwb25zZRI1CgRuZXh0GAEgASgLMiEuc2VhcmNocGx1Z2luLlBsdW'
    'dpblNlYXJjaFJlcXVlc3RSBG5leHQSKgoFaXRlbXMYAiADKAsyFC5zZWFyY2hwbHVnaW4uUGx1'
    'Z2luUgVpdGVtcw==');

@$core.Deprecated('Use pluginCreateRequestDescriptor instead')
const PluginCreateRequest$json = {
  '1': 'PluginCreateRequest',
};

/// Descriptor for `PluginCreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pluginCreateRequestDescriptor =
    $convert.base64Decode('ChNQbHVnaW5DcmVhdGVSZXF1ZXN0');

@$core.Deprecated('Use pluginCreateResponseDescriptor instead')
const PluginCreateResponse$json = {
  '1': 'PluginCreateResponse',
  '2': [
    {
      '1': 'plugin',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.searchplugin.Plugin',
      '10': 'plugin'
    },
  ],
};

/// Descriptor for `PluginCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pluginCreateResponseDescriptor = $convert.base64Decode(
    'ChRQbHVnaW5DcmVhdGVSZXNwb25zZRIsCgZwbHVnaW4YASABKAsyFC5zZWFyY2hwbHVnaW4uUG'
    'x1Z2luUgZwbHVnaW4=');

@$core.Deprecated('Use pluginFindRequestDescriptor instead')
const PluginFindRequest$json = {
  '1': 'PluginFindRequest',
};

/// Descriptor for `PluginFindRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pluginFindRequestDescriptor =
    $convert.base64Decode('ChFQbHVnaW5GaW5kUmVxdWVzdA==');

@$core.Deprecated('Use pluginFindResponseDescriptor instead')
const PluginFindResponse$json = {
  '1': 'PluginFindResponse',
  '2': [
    {
      '1': 'plugin',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.searchplugin.Plugin',
      '10': 'plugin'
    },
  ],
};

/// Descriptor for `PluginFindResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pluginFindResponseDescriptor = $convert.base64Decode(
    'ChJQbHVnaW5GaW5kUmVzcG9uc2USLAoGcGx1Z2luGAEgASgLMhQuc2VhcmNocGx1Z2luLlBsdW'
    'dpblIGcGx1Z2lu');

@$core.Deprecated('Use pluginDeleteRequestDescriptor instead')
const PluginDeleteRequest$json = {
  '1': 'PluginDeleteRequest',
};

/// Descriptor for `PluginDeleteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pluginDeleteRequestDescriptor =
    $convert.base64Decode('ChNQbHVnaW5EZWxldGVSZXF1ZXN0');

@$core.Deprecated('Use pluginDeleteResponseDescriptor instead')
const PluginDeleteResponse$json = {
  '1': 'PluginDeleteResponse',
  '2': [
    {
      '1': 'plugin',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.searchplugin.Plugin',
      '10': 'plugin'
    },
  ],
};

/// Descriptor for `PluginDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pluginDeleteResponseDescriptor = $convert.base64Decode(
    'ChRQbHVnaW5EZWxldGVSZXNwb25zZRIsCgZwbHVnaW4YASABKAsyFC5zZWFyY2hwbHVnaW4uUG'
    'x1Z2luUgZwbHVnaW4=');
