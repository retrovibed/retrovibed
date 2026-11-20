// This is a generated file - do not edit.
//
// Generated from torrents/torrent.proto.

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

@$core.Deprecated('Use peersDescriptor instead')
const Peers$json = {
  '1': 'Peers',
  '2': [
    {'1': 'min', '3': 1, '4': 1, '5': 13, '10': 'min'},
    {'1': 'max', '3': 2, '4': 1, '5': 13, '10': 'max'},
  ],
};

/// Descriptor for `Peers`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List peersDescriptor = $convert.base64Decode(
    'CgVQZWVycxIQCgNtaW4YASABKA1SA21pbhIQCgNtYXgYAiABKA1SA21heA==');

@$core.Deprecated('Use limitDescriptor instead')
const Limit$json = {
  '1': 'Limit',
  '2': [
    {'1': 'rate', '3': 1, '4': 1, '5': 13, '10': 'rate'},
    {'1': 'burst', '3': 2, '4': 1, '5': 13, '10': 'burst'},
  ],
};

/// Descriptor for `Limit`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List limitDescriptor = $convert.base64Decode(
    'CgVMaW1pdBISCgRyYXRlGAEgASgNUgRyYXRlEhQKBWJ1cnN0GAIgASgNUgVidXJzdA==');

@$core.Deprecated('Use torrentSettingsDescriptor instead')
const TorrentSettings$json = {
  '1': 'TorrentSettings',
  '2': [
    {'1': 'seed', '3': 1, '4': 1, '5': 8, '10': 'seed'},
    {'1': 'pex', '3': 2, '4': 1, '5': 8, '10': 'pex'},
    {
      '1': 'auto_locate_media',
      '3': 4,
      '4': 1,
      '5': 8,
      '10': 'auto_locate_media'
    },
    {'1': 'firewalled', '3': 5, '4': 1, '5': 8, '10': 'firewalled'},
    {'1': 'resumable', '3': 6, '4': 1, '5': 8, '10': 'resumable'},
    {'1': 'ip4', '3': 7, '4': 1, '5': 9, '10': 'ip4'},
    {'1': 'ip6', '3': 8, '4': 1, '5': 9, '10': 'ip6'},
    {'1': 'port', '3': 9, '4': 1, '5': 13, '10': 'port'},
    {'1': 'log', '3': 998, '4': 1, '5': 8, '10': 'log'},
    {'1': 'debug', '3': 999, '4': 1, '5': 8, '10': 'debug'},
    {
      '1': 'download',
      '3': 1000,
      '4': 1,
      '5': 11,
      '6': '.torrents.Limit',
      '10': 'download'
    },
    {
      '1': 'upload',
      '3': 1001,
      '4': 1,
      '5': 11,
      '6': '.torrents.Limit',
      '10': 'upload'
    },
    {
      '1': 'inbound',
      '3': 1002,
      '4': 1,
      '5': 11,
      '6': '.torrents.Limit',
      '10': 'inbound'
    },
    {
      '1': 'outbound',
      '3': 1003,
      '4': 1,
      '5': 11,
      '6': '.torrents.Limit',
      '10': 'outbound'
    },
    {
      '1': 'peers',
      '3': 1004,
      '4': 1,
      '5': 11,
      '6': '.torrents.Peers',
      '10': 'peers'
    },
    {
      '1': 'maximum_requests',
      '3': 1005,
      '4': 1,
      '5': 4,
      '10': 'maximum_requests'
    },
  ],
  '9': [
    {'1': 3, '2': 4},
    {'1': 10, '2': 11},
    {'1': 11, '2': 998},
  ],
  '10': ['auto_bootstrap', 'dht'],
};

/// Descriptor for `TorrentSettings`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List torrentSettingsDescriptor = $convert.base64Decode(
    'Cg9Ub3JyZW50U2V0dGluZ3MSEgoEc2VlZBgBIAEoCFIEc2VlZBIQCgNwZXgYAiABKAhSA3BleB'
    'IsChFhdXRvX2xvY2F0ZV9tZWRpYRgEIAEoCFIRYXV0b19sb2NhdGVfbWVkaWESHgoKZmlyZXdh'
    'bGxlZBgFIAEoCFIKZmlyZXdhbGxlZBIcCglyZXN1bWFibGUYBiABKAhSCXJlc3VtYWJsZRIQCg'
    'NpcDQYByABKAlSA2lwNBIQCgNpcDYYCCABKAlSA2lwNhISCgRwb3J0GAkgASgNUgRwb3J0EhEK'
    'A2xvZxjmByABKAhSA2xvZxIVCgVkZWJ1ZxjnByABKAhSBWRlYnVnEiwKCGRvd25sb2FkGOgHIA'
    'EoCzIPLnRvcnJlbnRzLkxpbWl0Ughkb3dubG9hZBIoCgZ1cGxvYWQY6QcgASgLMg8udG9ycmVu'
    'dHMuTGltaXRSBnVwbG9hZBIqCgdpbmJvdW5kGOoHIAEoCzIPLnRvcnJlbnRzLkxpbWl0Ugdpbm'
    'JvdW5kEiwKCG91dGJvdW5kGOsHIAEoCzIPLnRvcnJlbnRzLkxpbWl0UghvdXRib3VuZBImCgVw'
    'ZWVycxjsByABKAsyDy50b3JyZW50cy5QZWVyc1IFcGVlcnMSKwoQbWF4aW11bV9yZXF1ZXN0cx'
    'jtByABKARSEG1heGltdW1fcmVxdWVzdHNKBAgDEARKBAgKEAtKBQgLEOYHUg5hdXRvX2Jvb3Rz'
    'dHJhcFIDZGh0');

@$core.Deprecated('Use dHTSettingsDescriptor instead')
const DHTSettings$json = {
  '1': 'DHTSettings',
  '2': [
    {'1': 'passive', '3': 1, '4': 1, '5': 8, '10': 'passive'},
    {'1': 'auto_bootstrap', '3': 2, '4': 1, '5': 8, '10': 'auto_bootstrap'},
    {'1': 'log', '3': 998, '4': 1, '5': 8, '10': 'log'},
    {'1': 'debug', '3': 999, '4': 1, '5': 8, '10': 'debug'},
  ],
};

/// Descriptor for `DHTSettings`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dHTSettingsDescriptor = $convert.base64Decode(
    'CgtESFRTZXR0aW5ncxIYCgdwYXNzaXZlGAEgASgIUgdwYXNzaXZlEiYKDmF1dG9fYm9vdHN0cm'
    'FwGAIgASgIUg5hdXRvX2Jvb3RzdHJhcBIRCgNsb2cY5gcgASgIUgNsb2cSFQoFZGVidWcY5wcg'
    'ASgIUgVkZWJ1Zw==');

@$core.Deprecated('Use discoverySettingsDescriptor instead')
const DiscoverySettings$json = {
  '1': 'DiscoverySettings',
  '2': [
    {'1': 'enabled', '3': 1, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'ratio', '3': 2, '4': 1, '5': 13, '10': 'ratio'},
    {'1': 'partitions', '3': 3, '4': 1, '5': 13, '10': 'partitions'},
    {'1': 'workloads', '3': 4, '4': 1, '5': 13, '10': 'workloads'},
    {'1': 'seed', '3': 5, '4': 1, '5': 9, '10': 'seed'},
    {'1': 'locate_p2p', '3': 1000, '4': 1, '5': 8, '10': 'locate_p2p'},
  ],
};

/// Descriptor for `DiscoverySettings`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List discoverySettingsDescriptor = $convert.base64Decode(
    'ChFEaXNjb3ZlcnlTZXR0aW5ncxIYCgdlbmFibGVkGAEgASgIUgdlbmFibGVkEhQKBXJhdGlvGA'
    'IgASgNUgVyYXRpbxIeCgpwYXJ0aXRpb25zGAMgASgNUgpwYXJ0aXRpb25zEhwKCXdvcmtsb2Fk'
    'cxgEIAEoDVIJd29ya2xvYWRzEhIKBHNlZWQYBSABKAlSBHNlZWQSHwoKbG9jYXRlX3AycBjoBy'
    'ABKAhSCmxvY2F0ZV9wMnA=');
