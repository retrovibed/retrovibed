// This is a generated file - do not edit.
//
// Generated from media/media.remote.control.proto.

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

@$core.Deprecated('Use queueDescriptor instead')
const Queue$json = {
  '1': 'Queue',
  '2': [
    {'1': 'media', '3': 1, '4': 1, '5': 11, '6': '.media.Media', '10': 'media'},
  ],
};

/// Descriptor for `Queue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queueDescriptor = $convert.base64Decode(
    'CgVRdWV1ZRIiCgVtZWRpYRgBIAEoCzIMLm1lZGlhLk1lZGlhUgVtZWRpYQ==');

@$core.Deprecated('Use dequeueDescriptor instead')
const Dequeue$json = {
  '1': 'Dequeue',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `Dequeue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dequeueDescriptor =
    $convert.base64Decode('CgdEZXF1ZXVlEg4KAmlkGAEgASgJUgJpZA==');

@$core.Deprecated('Use playPauseDescriptor instead')
const PlayPause$json = {
  '1': 'PlayPause',
  '2': [
    {'1': 'paused', '3': 1, '4': 1, '5': 8, '10': 'paused'},
  ],
};

/// Descriptor for `PlayPause`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List playPauseDescriptor =
    $convert.base64Decode('CglQbGF5UGF1c2USFgoGcGF1c2VkGAEgASgIUgZwYXVzZWQ=');

@$core.Deprecated('Use seekDescriptor instead')
const Seek$json = {
  '1': 'Seek',
  '2': [
    {'1': 'offset', '3': 1, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `Seek`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List seekDescriptor =
    $convert.base64Decode('CgRTZWVrEhYKBm9mZnNldBgBIAEoBVIGb2Zmc2V0');

@$core.Deprecated('Use volumeDescriptor instead')
const Volume$json = {
  '1': 'Volume',
  '2': [
    {'1': 'muted', '3': 1, '4': 1, '5': 8, '10': 'muted'},
    {'1': 'level', '3': 2, '4': 1, '5': 2, '10': 'level'},
  ],
};

/// Descriptor for `Volume`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List volumeDescriptor = $convert.base64Decode(
    'CgZWb2x1bWUSFAoFbXV0ZWQYASABKAhSBW11dGVkEhQKBWxldmVsGAIgASgCUgVsZXZlbA==');

@$core.Deprecated('Use syncDescriptor instead')
const Sync$json = {
  '1': 'Sync',
  '2': [
    {
      '1': 'library',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.Daemon',
      '10': 'library'
    },
    {'1': 'capacity', '3': 2, '4': 1, '5': 13, '10': 'capacity'},
    {
      '1': 'current',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.media.Media',
      '10': 'current'
    },
    {'1': 'token', '3': 4, '4': 1, '5': 9, '10': 'token'},
    {'1': 'expiration', '3': 5, '4': 1, '5': 3, '10': 'expiration'},
    {'1': 'volume', '3': 6, '4': 1, '5': 2, '10': 'volume'},
    {
      '1': 'queue',
      '3': 1000,
      '4': 3,
      '5': 11,
      '6': '.media.Media',
      '10': 'queue'
    },
  ],
  '9': [
    {'1': 7, '2': 1000},
  ],
};

/// Descriptor for `Sync`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List syncDescriptor = $convert.base64Decode(
    'CgRTeW5jEiYKB2xpYnJhcnkYASABKAsyDC5tZXRhLkRhZW1vblIHbGlicmFyeRIaCghjYXBhY2'
    'l0eRgCIAEoDVIIY2FwYWNpdHkSJgoHY3VycmVudBgDIAEoCzIMLm1lZGlhLk1lZGlhUgdjdXJy'
    'ZW50EhQKBXRva2VuGAQgASgJUgV0b2tlbhIeCgpleHBpcmF0aW9uGAUgASgDUgpleHBpcmF0aW'
    '9uEhYKBnZvbHVtZRgGIAEoAlIGdm9sdW1lEiMKBXF1ZXVlGOgHIAMoCzIMLm1lZGlhLk1lZGlh'
    'UgVxdWV1ZUoFCAcQ6Ac=');

@$core.Deprecated('Use streamDescriptor instead')
const Stream$json = {
  '1': 'Stream',
  '2': [
    {'1': 'sid', '3': 1, '4': 1, '5': 9, '10': 'sid'},
    {
      '1': 'queue',
      '3': 1000,
      '4': 1,
      '5': 11,
      '6': '.media.Queue',
      '9': 0,
      '10': 'queue'
    },
    {
      '1': 'dequeue',
      '3': 1002,
      '4': 1,
      '5': 11,
      '6': '.media.Dequeue',
      '9': 0,
      '10': 'dequeue'
    },
    {
      '1': 'playpause',
      '3': 1003,
      '4': 1,
      '5': 11,
      '6': '.media.PlayPause',
      '9': 0,
      '10': 'playpause'
    },
    {
      '1': 'seek',
      '3': 1004,
      '4': 1,
      '5': 11,
      '6': '.media.Seek',
      '9': 0,
      '10': 'seek'
    },
    {
      '1': 'sync',
      '3': 1005,
      '4': 1,
      '5': 11,
      '6': '.media.Sync',
      '9': 0,
      '10': 'sync'
    },
    {
      '1': 'volume',
      '3': 1006,
      '4': 1,
      '5': 11,
      '6': '.media.Volume',
      '9': 0,
      '10': 'volume'
    },
  ],
  '8': [
    {'1': 'Command'},
  ],
};

/// Descriptor for `Stream`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamDescriptor = $convert.base64Decode(
    'CgZTdHJlYW0SEAoDc2lkGAEgASgJUgNzaWQSJQoFcXVldWUY6AcgASgLMgwubWVkaWEuUXVldW'
    'VIAFIFcXVldWUSKwoHZGVxdWV1ZRjqByABKAsyDi5tZWRpYS5EZXF1ZXVlSABSB2RlcXVldWUS'
    'MQoJcGxheXBhdXNlGOsHIAEoCzIQLm1lZGlhLlBsYXlQYXVzZUgAUglwbGF5cGF1c2USIgoEc2'
    'VlaxjsByABKAsyCy5tZWRpYS5TZWVrSABSBHNlZWsSIgoEc3luYxjtByABKAsyCy5tZWRpYS5T'
    'eW5jSABSBHN5bmMSKAoGdm9sdW1lGO4HIAEoCzINLm1lZGlhLlZvbHVtZUgAUgZ2b2x1bWVCCQ'
    'oHQ29tbWFuZA==');
