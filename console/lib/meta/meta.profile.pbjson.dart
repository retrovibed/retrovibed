//
//  Generated code. Do not modify.
//  source: meta.profile.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use profileDescriptor instead')
const Profile$json = {
  '1': 'Profile',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'account_id', '3': 2, '4': 1, '5': 9, '10': 'account_id'},
    {'1': 'session_watermark', '3': 3, '4': 1, '5': 9, '10': 'session_watermark'},
    {'1': 'created_at', '3': 4, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 5, '4': 1, '5': 9, '10': 'updated_at'},
    {'1': 'disabled_at', '3': 6, '4': 1, '5': 9, '10': 'disabled_at'},
    {'1': 'disabled_manually_at', '3': 7, '4': 1, '5': 9, '10': 'disabled_manually_at'},
    {'1': 'disabled_pending_approval_at', '3': 8, '4': 1, '5': 9, '10': 'disabled_pending_approval_at'},
    {'1': 'display', '3': 9, '4': 1, '5': 9, '10': 'display'},
    {'1': 'email', '3': 10, '4': 1, '5': 9, '10': 'email'},
  ],
};

/// Descriptor for `Profile`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileDescriptor = $convert.base64Decode(
    'CgdQcm9maWxlEg4KAmlkGAEgASgJUgJpZBIeCgphY2NvdW50X2lkGAIgASgJUgphY2NvdW50X2'
    'lkEiwKEXNlc3Npb25fd2F0ZXJtYXJrGAMgASgJUhFzZXNzaW9uX3dhdGVybWFyaxIeCgpjcmVh'
    'dGVkX2F0GAQgASgJUgpjcmVhdGVkX2F0Eh4KCnVwZGF0ZWRfYXQYBSABKAlSCnVwZGF0ZWRfYX'
    'QSIAoLZGlzYWJsZWRfYXQYBiABKAlSC2Rpc2FibGVkX2F0EjIKFGRpc2FibGVkX21hbnVhbGx5'
    'X2F0GAcgASgJUhRkaXNhYmxlZF9tYW51YWxseV9hdBJCChxkaXNhYmxlZF9wZW5kaW5nX2FwcH'
    'JvdmFsX2F0GAggASgJUhxkaXNhYmxlZF9wZW5kaW5nX2FwcHJvdmFsX2F0EhgKB2Rpc3BsYXkY'
    'CSABKAlSB2Rpc3BsYXkSFAoFZW1haWwYCiABKAlSBWVtYWls');

@$core.Deprecated('Use profileSearchRequestDescriptor instead')
const ProfileSearchRequest$json = {
  '1': 'ProfileSearchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'offset', '3': 2, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 3, '4': 1, '5': 4, '10': 'limit'},
    {'1': 'status', '3': 4, '4': 1, '5': 13, '10': 'status'},
  ],
};

/// Descriptor for `ProfileSearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileSearchRequestDescriptor = $convert.base64Decode(
    'ChRQcm9maWxlU2VhcmNoUmVxdWVzdBIUCgVxdWVyeRgBIAEoCVIFcXVlcnkSFgoGb2Zmc2V0GA'
    'IgASgEUgZvZmZzZXQSFAoFbGltaXQYAyABKARSBWxpbWl0EhYKBnN0YXR1cxgEIAEoDVIGc3Rh'
    'dHVz');

@$core.Deprecated('Use profileSearchResponseDescriptor instead')
const ProfileSearchResponse$json = {
  '1': 'ProfileSearchResponse',
  '2': [
    {'1': 'next', '3': 1, '4': 1, '5': 11, '6': '.meta.ProfileSearchRequest', '10': 'next'},
    {'1': 'items', '3': 2, '4': 3, '5': 11, '6': '.meta.Profile', '10': 'items'},
  ],
};

/// Descriptor for `ProfileSearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileSearchResponseDescriptor = $convert.base64Decode(
    'ChVQcm9maWxlU2VhcmNoUmVzcG9uc2USLgoEbmV4dBgBIAEoCzIaLm1ldGEuUHJvZmlsZVNlYX'
    'JjaFJlcXVlc3RSBG5leHQSIwoFaXRlbXMYAiADKAsyDS5tZXRhLlByb2ZpbGVSBWl0ZW1z');

@$core.Deprecated('Use profileCreateRequestDescriptor instead')
const ProfileCreateRequest$json = {
  '1': 'ProfileCreateRequest',
  '2': [
    {'1': 'profile', '3': 1, '4': 1, '5': 11, '6': '.meta.Profile', '10': 'profile'},
  ],
};

/// Descriptor for `ProfileCreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileCreateRequestDescriptor = $convert.base64Decode(
    'ChRQcm9maWxlQ3JlYXRlUmVxdWVzdBInCgdwcm9maWxlGAEgASgLMg0ubWV0YS5Qcm9maWxlUg'
    'dwcm9maWxl');

@$core.Deprecated('Use profileCreateResponseDescriptor instead')
const ProfileCreateResponse$json = {
  '1': 'ProfileCreateResponse',
  '2': [
    {'1': 'profile', '3': 1, '4': 1, '5': 11, '6': '.meta.Profile', '10': 'profile'},
  ],
};

/// Descriptor for `ProfileCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileCreateResponseDescriptor = $convert.base64Decode(
    'ChVQcm9maWxlQ3JlYXRlUmVzcG9uc2USJwoHcHJvZmlsZRgBIAEoCzINLm1ldGEuUHJvZmlsZV'
    'IHcHJvZmlsZQ==');

@$core.Deprecated('Use profileLookupRequestDescriptor instead')
const ProfileLookupRequest$json = {
  '1': 'ProfileLookupRequest',
};

/// Descriptor for `ProfileLookupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileLookupRequestDescriptor = $convert.base64Decode(
    'ChRQcm9maWxlTG9va3VwUmVxdWVzdA==');

@$core.Deprecated('Use profileLookupResponseDescriptor instead')
const ProfileLookupResponse$json = {
  '1': 'ProfileLookupResponse',
  '2': [
    {'1': 'profile', '3': 1, '4': 1, '5': 11, '6': '.meta.Profile', '10': 'profile'},
  ],
};

/// Descriptor for `ProfileLookupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileLookupResponseDescriptor = $convert.base64Decode(
    'ChVQcm9maWxlTG9va3VwUmVzcG9uc2USJwoHcHJvZmlsZRgBIAEoCzINLm1ldGEuUHJvZmlsZV'
    'IHcHJvZmlsZQ==');

@$core.Deprecated('Use profileUpdateRequestDescriptor instead')
const ProfileUpdateRequest$json = {
  '1': 'ProfileUpdateRequest',
  '2': [
    {'1': 'profile', '3': 1, '4': 1, '5': 11, '6': '.meta.Profile', '10': 'profile'},
  ],
};

/// Descriptor for `ProfileUpdateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileUpdateRequestDescriptor = $convert.base64Decode(
    'ChRQcm9maWxlVXBkYXRlUmVxdWVzdBInCgdwcm9maWxlGAEgASgLMg0ubWV0YS5Qcm9maWxlUg'
    'dwcm9maWxl');

@$core.Deprecated('Use profileUpdateResponseDescriptor instead')
const ProfileUpdateResponse$json = {
  '1': 'ProfileUpdateResponse',
  '2': [
    {'1': 'profile', '3': 1, '4': 1, '5': 11, '6': '.meta.Profile', '10': 'profile'},
  ],
};

/// Descriptor for `ProfileUpdateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileUpdateResponseDescriptor = $convert.base64Decode(
    'ChVQcm9maWxlVXBkYXRlUmVzcG9uc2USJwoHcHJvZmlsZRgBIAEoCzINLm1ldGEuUHJvZmlsZV'
    'IHcHJvZmlsZQ==');

@$core.Deprecated('Use profileDisableRequestDescriptor instead')
const ProfileDisableRequest$json = {
  '1': 'ProfileDisableRequest',
  '2': [
    {'1': 'profile', '3': 1, '4': 1, '5': 11, '6': '.meta.Profile', '10': 'profile'},
  ],
};

/// Descriptor for `ProfileDisableRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileDisableRequestDescriptor = $convert.base64Decode(
    'ChVQcm9maWxlRGlzYWJsZVJlcXVlc3QSJwoHcHJvZmlsZRgBIAEoCzINLm1ldGEuUHJvZmlsZV'
    'IHcHJvZmlsZQ==');

@$core.Deprecated('Use profileDisableResponseDescriptor instead')
const ProfileDisableResponse$json = {
  '1': 'ProfileDisableResponse',
  '2': [
    {'1': 'profile', '3': 1, '4': 1, '5': 11, '6': '.meta.Profile', '10': 'profile'},
  ],
};

/// Descriptor for `ProfileDisableResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileDisableResponseDescriptor = $convert.base64Decode(
    'ChZQcm9maWxlRGlzYWJsZVJlc3BvbnNlEicKB3Byb2ZpbGUYASABKAsyDS5tZXRhLlByb2ZpbG'
    'VSB3Byb2ZpbGU=');

