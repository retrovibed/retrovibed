// This is a generated file - do not edit.
//
// Generated from torrent.proto.

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
    {'1': 'auto_bootstrap', '3': 3, '4': 1, '5': 8, '10': 'auto_bootstrap'},
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
    {'1': 10, '2': 998},
  ],
};

/// Descriptor for `TorrentSettings`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List torrentSettingsDescriptor = $convert.base64Decode(
    'Cg9Ub3JyZW50U2V0dGluZ3MSEgoEc2VlZBgBIAEoCFIEc2VlZBIQCgNwZXgYAiABKAhSA3BleB'
    'ImCg5hdXRvX2Jvb3RzdHJhcBgDIAEoCFIOYXV0b19ib290c3RyYXASLAoRYXV0b19sb2NhdGVf'
    'bWVkaWEYBCABKAhSEWF1dG9fbG9jYXRlX21lZGlhEh4KCmZpcmV3YWxsZWQYBSABKAhSCmZpcm'
    'V3YWxsZWQSHAoJcmVzdW1hYmxlGAYgASgIUglyZXN1bWFibGUSEAoDaXA0GAcgASgJUgNpcDQS'
    'EAoDaXA2GAggASgJUgNpcDYSEgoEcG9ydBgJIAEoDVIEcG9ydBIRCgNsb2cY5gcgASgIUgNsb2'
    'cSFQoFZGVidWcY5wcgASgIUgVkZWJ1ZxIsCghkb3dubG9hZBjoByABKAsyDy50b3JyZW50cy5M'
    'aW1pdFIIZG93bmxvYWQSKAoGdXBsb2FkGOkHIAEoCzIPLnRvcnJlbnRzLkxpbWl0UgZ1cGxvYW'
    'QSKgoHaW5ib3VuZBjqByABKAsyDy50b3JyZW50cy5MaW1pdFIHaW5ib3VuZBIsCghvdXRib3Vu'
    'ZBjrByABKAsyDy50b3JyZW50cy5MaW1pdFIIb3V0Ym91bmQSJgoFcGVlcnMY7AcgASgLMg8udG'
    '9ycmVudHMuUGVlcnNSBXBlZXJzEisKEG1heGltdW1fcmVxdWVzdHMY7QcgASgEUhBtYXhpbXVt'
    'X3JlcXVlc3RzSgUIChDmBw==');

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
