//
//  Generated code. Do not modify.
//  source: media.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use mediaDescriptor instead')
const Media$json = {
  '1': 'Media',
  '2': [
    {'1': 'title', '3': 1, '4': 1, '5': 9, '10': 'title'},
    {'1': 'description', '3': 2, '4': 1, '5': 9, '10': 'description'},
    {'1': 'mimetype', '3': 3, '4': 1, '5': 9, '10': 'mimetype'},
    {'1': 'image', '3': 4, '4': 1, '5': 9, '10': 'image'},
  ],
};

/// Descriptor for `Media`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaDescriptor = $convert.base64Decode(
    'CgVNZWRpYRIUCgV0aXRsZRgBIAEoCVIFdGl0bGUSIAoLZGVzY3JpcHRpb24YAiABKAlSC2Rlc2'
    'NyaXB0aW9uEhoKCG1pbWV0eXBlGAMgASgJUghtaW1ldHlwZRIUCgVpbWFnZRgEIAEoCVIFaW1h'
    'Z2U=');

@$core.Deprecated('Use mediaRequestDescriptor instead')
const MediaRequest$json = {
  '1': 'MediaRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 2, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `MediaRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaRequestDescriptor = $convert.base64Decode(
    'CgxNZWRpYVJlcXVlc3QSFAoFcXVlcnkYASABKAlSBXF1ZXJ5EhcKBm9mZnNldBiEByABKARSBm'
    '9mZnNldBIVCgVsaW1pdBiFByABKARSBWxpbWl0SgUIAhCEB0oGCIYHEOgH');

@$core.Deprecated('Use mediaResponseDescriptor instead')
const MediaResponse$json = {
  '1': 'MediaResponse',
  '2': [
    {'1': 'next', '3': 1, '4': 1, '5': 11, '6': '.media.MediaRequest', '10': 'next'},
    {'1': 'items', '3': 2, '4': 3, '5': 11, '6': '.media.Media', '10': 'items'},
  ],
};

/// Descriptor for `MediaResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mediaResponseDescriptor = $convert.base64Decode(
    'Cg1NZWRpYVJlc3BvbnNlEicKBG5leHQYASABKAsyEy5tZWRpYS5NZWRpYVJlcXVlc3RSBG5leH'
    'QSIgoFaXRlbXMYAiADKAsyDC5tZWRpYS5NZWRpYVIFaXRlbXM=');

