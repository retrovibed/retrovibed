// This is a generated file - do not edit.
//
// Generated from meta/meta.authn.proto.

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

@$core.Deprecated('Use loginOptionsDescriptor instead')
const LoginOptions$json = {
  '1': 'LoginOptions',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'account_id', '3': 2, '4': 1, '5': 9, '10': 'account_id'},
    {'1': 'profile_id', '3': 3, '4': 1, '5': 9, '10': 'profile_id'},
    {'1': 'redirect', '3': 4, '4': 1, '5': 9, '10': 'redirect'},
    {'1': 'autologin', '3': 5, '4': 1, '5': 8, '10': 'autologin'},
  ],
};

/// Descriptor for `LoginOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loginOptionsDescriptor = $convert.base64Decode(
    'CgxMb2dpbk9wdGlvbnMSDgoCaWQYASABKAlSAmlkEh4KCmFjY291bnRfaWQYAiABKAlSCmFjY2'
    '91bnRfaWQSHgoKcHJvZmlsZV9pZBgDIAEoCVIKcHJvZmlsZV9pZBIaCghyZWRpcmVjdBgEIAEo'
    'CVIIcmVkaXJlY3QSHAoJYXV0b2xvZ2luGAUgASgIUglhdXRvbG9naW4=');

@$core.Deprecated('Use loginOptionsResponseDescriptor instead')
const LoginOptionsResponse$json = {
  '1': 'LoginOptionsResponse',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 9, '10': 'token'},
  ],
};

/// Descriptor for `LoginOptionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loginOptionsResponseDescriptor =
    $convert.base64Decode(
        'ChRMb2dpbk9wdGlvbnNSZXNwb25zZRIUCgV0b2tlbhgBIAEoCVIFdG9rZW4=');

@$core.Deprecated('Use identityDescriptor instead')
const Identity$json = {
  '1': 'Identity',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'display', '3': 2, '4': 1, '5': 9, '10': 'display'},
    {'1': 'email', '3': 4, '4': 1, '5': 9, '10': 'email'},
    {'1': 'uid', '3': 7, '4': 1, '5': 9, '10': 'uid'},
    {'1': 'uidraw', '3': 8, '4': 1, '5': 9, '10': 'uidraw'},
  ],
};

/// Descriptor for `Identity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List identityDescriptor = $convert.base64Decode(
    'CghJZGVudGl0eRIOCgJpZBgBIAEoCVICaWQSGAoHZGlzcGxheRgCIAEoCVIHZGlzcGxheRIUCg'
    'VlbWFpbBgEIAEoCVIFZW1haWwSEAoDdWlkGAcgASgJUgN1aWQSFgoGdWlkcmF3GAggASgJUgZ1'
    'aWRyYXc=');

@$core.Deprecated('Use authnDescriptor instead')
const Authn$json = {
  '1': 'Authn',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 9, '10': 'token'},
    {
      '1': 'profile',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.meta.Profile',
      '10': 'profile'
    },
    {
      '1': 'account',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.meta.Account',
      '10': 'account'
    },
  ],
};

/// Descriptor for `Authn`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authnDescriptor = $convert.base64Decode(
    'CgVBdXRobhIUCgV0b2tlbhgBIAEoCVIFdG9rZW4SJwoHcHJvZmlsZRgCIAEoCzINLm1ldGEuUH'
    'JvZmlsZVIHcHJvZmlsZRInCgdhY2NvdW50GAMgASgLMg0ubWV0YS5BY2NvdW50UgdhY2NvdW50');

@$core.Deprecated('Use authedDescriptor instead')
const Authed$json = {
  '1': 'Authed',
  '2': [
    {'1': 'signup_token', '3': 1, '4': 1, '5': 9, '10': 'signup_token'},
    {
      '1': 'identity',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.meta.Identity',
      '10': 'identity'
    },
    {
      '1': 'profiles',
      '3': 4,
      '4': 3,
      '5': 11,
      '6': '.meta.Authn',
      '10': 'profiles'
    },
    {
      '1': 'options',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.meta.LoginOptions',
      '10': 'opts'
    },
  ],
};

/// Descriptor for `Authed`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authedDescriptor = $convert.base64Decode(
    'CgZBdXRoZWQSIgoMc2lnbnVwX3Rva2VuGAEgASgJUgxzaWdudXBfdG9rZW4SKgoIaWRlbnRpdH'
    'kYAyABKAsyDi5tZXRhLklkZW50aXR5UghpZGVudGl0eRInCghwcm9maWxlcxgEIAMoCzILLm1l'
    'dGEuQXV0aG5SCHByb2ZpbGVzEikKB29wdGlvbnMYBSABKAsyEi5tZXRhLkxvZ2luT3B0aW9uc1'
    'IEb3B0cw==');

@$core.Deprecated('Use sessionDescriptor instead')
const Session$json = {
  '1': 'Session',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 9, '10': 'token'},
    {
      '1': 'profile',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.meta.Profile',
      '10': 'profile'
    },
    {
      '1': 'account',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.meta.Account',
      '10': 'account'
    },
    {'1': 'expires', '3': 4, '4': 1, '5': 3, '10': 'expires'},
    {'1': 'redirect', '3': 5, '4': 1, '5': 9, '10': 'redirect'},
  ],
};

/// Descriptor for `Session`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionDescriptor = $convert.base64Decode(
    'CgdTZXNzaW9uEhQKBXRva2VuGAEgASgJUgV0b2tlbhInCgdwcm9maWxlGAIgASgLMg0ubWV0YS'
    '5Qcm9maWxlUgdwcm9maWxlEicKB2FjY291bnQYAyABKAsyDS5tZXRhLkFjY291bnRSB2FjY291'
    'bnQSGAoHZXhwaXJlcxgEIAEoA1IHZXhwaXJlcxIaCghyZWRpcmVjdBgFIAEoCVIIcmVkaXJlY3'
    'Q=');
