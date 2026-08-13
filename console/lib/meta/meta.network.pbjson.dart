// This is a generated file - do not edit.
//
// Generated from meta/meta.network.proto.

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

@$core.Deprecated('Use wireguardDiagnosticsDescriptor instead')
const WireguardDiagnostics$json = {
  '1': 'WireguardDiagnostics',
  '2': [
    {'1': 'peer_key', '3': 1, '4': 1, '5': 9, '10': 'peer_key'},
    {
      '1': 'keepalive_interval',
      '3': 2,
      '4': 1,
      '5': 4,
      '10': 'keepalive_interval'
    },
    {'1': 'tx_bytes', '3': 3, '4': 1, '5': 4, '10': 'tx_bytes'},
    {'1': 'rx_bytes', '3': 4, '4': 1, '5': 4, '10': 'rx_bytes'},
    {
      '1': 'last_handshake_sec',
      '3': 5,
      '4': 1,
      '5': 3,
      '10': 'last_handshake_sec'
    },
    {'1': 'status', '3': 6, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `WireguardDiagnostics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List wireguardDiagnosticsDescriptor = $convert.base64Decode(
    'ChRXaXJlZ3VhcmREaWFnbm9zdGljcxIaCghwZWVyX2tleRgBIAEoCVIIcGVlcl9rZXkSLgoSa2'
    'VlcGFsaXZlX2ludGVydmFsGAIgASgEUhJrZWVwYWxpdmVfaW50ZXJ2YWwSGgoIdHhfYnl0ZXMY'
    'AyABKARSCHR4X2J5dGVzEhoKCHJ4X2J5dGVzGAQgASgEUghyeF9ieXRlcxIuChJsYXN0X2hhbm'
    'RzaGFrZV9zZWMYBSABKANSEmxhc3RfaGFuZHNoYWtlX3NlYxIWCgZzdGF0dXMYBiABKAlSBnN0'
    'YXR1cw==');

@$core.Deprecated('Use networkInterfaceDescriptor instead')
const NetworkInterface$json = {
  '1': 'NetworkInterface',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'ip', '3': 2, '4': 1, '5': 9, '10': 'ip'},
    {'1': 'metered', '3': 3, '4': 1, '5': 8, '10': 'metered'},
  ],
};

/// Descriptor for `NetworkInterface`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List networkInterfaceDescriptor = $convert.base64Decode(
    'ChBOZXR3b3JrSW50ZXJmYWNlEhIKBG5hbWUYASABKAlSBG5hbWUSDgoCaXAYAiABKAlSAmlwEh'
    'gKB21ldGVyZWQYAyABKAhSB21ldGVyZWQ=');

@$core.Deprecated('Use networkDescriptor instead')
const Network$json = {
  '1': 'Network',
  '2': [
    {
      '1': 'interfaces',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.meta.NetworkInterface',
      '10': 'interfaces'
    },
    {'1': 'have_v4', '3': 2, '4': 1, '5': 8, '10': 'have_v4'},
    {'1': 'have_v6', '3': 3, '4': 1, '5': 8, '10': 'have_v6'},
    {
      '1': 'default_interface',
      '3': 4,
      '4': 1,
      '5': 9,
      '10': 'default_interface'
    },
  ],
};

/// Descriptor for `Network`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List networkDescriptor = $convert.base64Decode(
    'CgdOZXR3b3JrEjYKCmludGVyZmFjZXMYASADKAsyFi5tZXRhLk5ldHdvcmtJbnRlcmZhY2VSCm'
    'ludGVyZmFjZXMSGAoHaGF2ZV92NBgCIAEoCFIHaGF2ZV92NBIYCgdoYXZlX3Y2GAMgASgIUgdo'
    'YXZlX3Y2EiwKEWRlZmF1bHRfaW50ZXJmYWNlGAQgASgJUhFkZWZhdWx0X2ludGVyZmFjZQ==');

@$core.Deprecated('Use networkMetricsResponseDescriptor instead')
const NetworkMetricsResponse$json = {
  '1': 'NetworkMetricsResponse',
  '2': [
    {
      '1': 'wireguard',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.WireguardDiagnostics',
      '10': 'wireguard'
    },
    {
      '1': 'network',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.meta.Network',
      '10': 'network'
    },
  ],
};

/// Descriptor for `NetworkMetricsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List networkMetricsResponseDescriptor = $convert.base64Decode(
    'ChZOZXR3b3JrTWV0cmljc1Jlc3BvbnNlEjgKCXdpcmVndWFyZBgBIAEoCzIaLm1ldGEuV2lyZW'
    'd1YXJkRGlhZ25vc3RpY3NSCXdpcmVndWFyZBInCgduZXR3b3JrGAIgASgLMg0ubWV0YS5OZXR3'
    'b3JrUgduZXR3b3Jr');
