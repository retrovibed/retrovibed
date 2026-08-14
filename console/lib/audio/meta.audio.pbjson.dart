// This is a generated file - do not edit.
//
// Generated from audio/meta.audio.proto.

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

@$core.Deprecated('Use audioSinkDescriptor instead')
const AudioSink$json = {
  '1': 'AudioSink',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `AudioSink`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List audioSinkDescriptor = $convert.base64Decode(
    'CglBdWRpb1NpbmsSDgoCaWQYASABKAlSAmlkEhIKBG5hbWUYAiABKAlSBG5hbWU=');

@$core.Deprecated('Use audioSinkSearchResponseDescriptor instead')
const AudioSinkSearchResponse$json = {
  '1': 'AudioSinkSearchResponse',
  '2': [
    {
      '1': 'items',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.meta.AudioSink',
      '10': 'items'
    },
  ],
};

/// Descriptor for `AudioSinkSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List audioSinkSearchResponseDescriptor =
    $convert.base64Decode(
        'ChdBdWRpb1NpbmtTZWFyY2hSZXNwb25zZRIlCgVpdGVtcxgBIAMoCzIPLm1ldGEuQXVkaW9TaW'
        '5rUgVpdGVtcw==');

@$core.Deprecated('Use audioSinkCurrentResponseDescriptor instead')
const AudioSinkCurrentResponse$json = {
  '1': 'AudioSinkCurrentResponse',
  '2': [
    {
      '1': 'sink',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.AudioSink',
      '10': 'sink'
    },
  ],
};

/// Descriptor for `AudioSinkCurrentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List audioSinkCurrentResponseDescriptor =
    $convert.base64Decode(
        'ChhBdWRpb1NpbmtDdXJyZW50UmVzcG9uc2USIwoEc2luaxgBIAEoCzIPLm1ldGEuQXVkaW9TaW'
        '5rUgRzaW5r');

@$core.Deprecated('Use audioSinkTouchRequestDescriptor instead')
const AudioSinkTouchRequest$json = {
  '1': 'AudioSinkTouchRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `AudioSinkTouchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List audioSinkTouchRequestDescriptor = $convert
    .base64Decode('ChVBdWRpb1NpbmtUb3VjaFJlcXVlc3QSDgoCaWQYASABKAlSAmlk');

@$core.Deprecated('Use audioSinkTouchResponseDescriptor instead')
const AudioSinkTouchResponse$json = {
  '1': 'AudioSinkTouchResponse',
  '2': [
    {
      '1': 'sink',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.AudioSink',
      '10': 'sink'
    },
  ],
};

/// Descriptor for `AudioSinkTouchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List audioSinkTouchResponseDescriptor =
    $convert.base64Decode(
        'ChZBdWRpb1NpbmtUb3VjaFJlc3BvbnNlEiMKBHNpbmsYASABKAsyDy5tZXRhLkF1ZGlvU2lua1'
        'IEc2luaw==');
