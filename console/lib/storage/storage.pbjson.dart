// This is a generated file - do not edit.
//
// Generated from storage.proto.

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

@$core.Deprecated('Use localDescriptor instead')
const Local$json = {
  '1': 'Local',
  '2': [
    {'1': 'reclaim', '3': 1, '4': 1, '5': 8, '10': 'reclaim'},
    {'1': 'maximum', '3': 2, '4': 1, '5': 4, '10': 'maximum'},
  ],
};

/// Descriptor for `Local`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List localDescriptor = $convert.base64Decode(
    'CgVMb2NhbBIYCgdyZWNsYWltGAEgASgIUgdyZWNsYWltEhgKB21heGltdW0YAiABKARSB21heG'
    'ltdW0=');

@$core.Deprecated('Use storageSettingsRequestDescriptor instead')
const StorageSettingsRequest$json = {
  '1': 'StorageSettingsRequest',
  '2': [
    {
      '1': 'local',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.storage.Local',
      '10': 'local'
    },
  ],
};

/// Descriptor for `StorageSettingsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List storageSettingsRequestDescriptor =
    $convert.base64Decode(
        'ChZTdG9yYWdlU2V0dGluZ3NSZXF1ZXN0EiQKBWxvY2FsGAEgASgLMg4uc3RvcmFnZS5Mb2NhbF'
        'IFbG9jYWw=');

@$core.Deprecated('Use storageSettingsResponseDescriptor instead')
const StorageSettingsResponse$json = {
  '1': 'StorageSettingsResponse',
  '2': [
    {
      '1': 'local',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.storage.Local',
      '10': 'local'
    },
  ],
};

/// Descriptor for `StorageSettingsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List storageSettingsResponseDescriptor =
    $convert.base64Decode(
        'ChdTdG9yYWdlU2V0dGluZ3NSZXNwb25zZRIkCgVsb2NhbBgBIAEoCzIOLnN0b3JhZ2UuTG9jYW'
        'xSBWxvY2Fs');
