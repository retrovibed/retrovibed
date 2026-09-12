// This is a generated file - do not edit.
//
// Generated from ftux/ftux.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import '../community/community.pb.dart' as $0;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class CommunitySuggestions extends $pb.GeneratedMessage {
  factory CommunitySuggestions({
    $core.Iterable<$0.Community>? community,
  }) {
    final result = CommunitySuggestions._();
    if (community != null) result.community.addAll(community);
    return result;
  }

  CommunitySuggestions._();

  factory CommunitySuggestions.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CommunitySuggestions()..mergeFromBuffer(data, registry);
  factory CommunitySuggestions.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CommunitySuggestions()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunitySuggestions',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.ftux'),
      createEmptyInstance: CommunitySuggestions.$_createMessage)
    ..pPM<$0.Community>(1000, _omitFieldNames ? '' : 'community',
        subBuilder: $0.Community.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunitySuggestions clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunitySuggestions copyWith(void Function(CommunitySuggestions) updates) =>
      super.copyWith((message) => updates(message as CommunitySuggestions))
          as CommunitySuggestions;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CommunitySuggestions() / CommunitySuggestions.new instead')
  static CommunitySuggestions create() => CommunitySuggestions._();
  static $pb.GeneratedMessage $_createMessage() => CommunitySuggestions._();
  @$core.override
  CommunitySuggestions createEmptyInstance() => CommunitySuggestions._();
  @$core.pragma('dart2js:noInline')
  static CommunitySuggestions getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunitySuggestions>(
          CommunitySuggestions.$_createMessage);
  static CommunitySuggestions? _defaultInstance;

  @$pb.TagNumber(1000)
  $pb.PbList<$0.Community> get community => $_getList(0);
}

class SubscribeCommunitiesRequest extends $pb.GeneratedMessage {
  factory SubscribeCommunitiesRequest({
    $core.Iterable<$core.String>? communityId,
  }) {
    final result = SubscribeCommunitiesRequest._();
    if (communityId != null) result.communityId.addAll(communityId);
    return result;
  }

  SubscribeCommunitiesRequest._();

  factory SubscribeCommunitiesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SubscribeCommunitiesRequest()..mergeFromBuffer(data, registry);
  factory SubscribeCommunitiesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SubscribeCommunitiesRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SubscribeCommunitiesRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.ftux'),
      createEmptyInstance: SubscribeCommunitiesRequest.$_createMessage)
    ..pPS(1, _omitFieldNames ? '' : 'community_id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeCommunitiesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeCommunitiesRequest copyWith(
          void Function(SubscribeCommunitiesRequest) updates) =>
      super.copyWith(
              (message) => updates(message as SubscribeCommunitiesRequest))
          as SubscribeCommunitiesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use SubscribeCommunitiesRequest() / SubscribeCommunitiesRequest.new instead')
  static SubscribeCommunitiesRequest create() =>
      SubscribeCommunitiesRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      SubscribeCommunitiesRequest._();
  @$core.override
  SubscribeCommunitiesRequest createEmptyInstance() =>
      SubscribeCommunitiesRequest._();
  @$core.pragma('dart2js:noInline')
  static SubscribeCommunitiesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SubscribeCommunitiesRequest>(
          SubscribeCommunitiesRequest.$_createMessage);
  static SubscribeCommunitiesRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get communityId => $_getList(0);
}

class SubscribeCommunitiesResponse extends $pb.GeneratedMessage {
  factory SubscribeCommunitiesResponse() => SubscribeCommunitiesResponse._();

  SubscribeCommunitiesResponse._();

  factory SubscribeCommunitiesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SubscribeCommunitiesResponse()..mergeFromBuffer(data, registry);
  factory SubscribeCommunitiesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SubscribeCommunitiesResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SubscribeCommunitiesResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'retrovibed.ftux'),
      createEmptyInstance: SubscribeCommunitiesResponse.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeCommunitiesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeCommunitiesResponse copyWith(
          void Function(SubscribeCommunitiesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as SubscribeCommunitiesResponse))
          as SubscribeCommunitiesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use SubscribeCommunitiesResponse() / SubscribeCommunitiesResponse.new instead')
  static SubscribeCommunitiesResponse create() =>
      SubscribeCommunitiesResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      SubscribeCommunitiesResponse._();
  @$core.override
  SubscribeCommunitiesResponse createEmptyInstance() =>
      SubscribeCommunitiesResponse._();
  @$core.pragma('dart2js:noInline')
  static SubscribeCommunitiesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SubscribeCommunitiesResponse>(
          SubscribeCommunitiesResponse.$_createMessage);
  static SubscribeCommunitiesResponse? _defaultInstance;
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
