// This is a generated file - do not edit.
//
// Generated from community.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'community.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'community.pbenum.dart';

class Community extends $pb.GeneratedMessage {
  factory Community({
    $core.String? id,
    $core.String? accountId,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? mimetype,
    $core.String? domain,
    $core.String? description,
    $core.String? entropy,
    $fixnum.Int64? bytes,
    $core.String? subscribedAt,
    PublishMode? defaultPublishMode,
    $core.bool? hidden,
    $core.String? url,
    $core.bool? adult,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (accountId != null) result.accountId = accountId;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (mimetype != null) result.mimetype = mimetype;
    if (domain != null) result.domain = domain;
    if (description != null) result.description = description;
    if (entropy != null) result.entropy = entropy;
    if (bytes != null) result.bytes = bytes;
    if (subscribedAt != null) result.subscribedAt = subscribedAt;
    if (defaultPublishMode != null)
      result.defaultPublishMode = defaultPublishMode;
    if (hidden != null) result.hidden = hidden;
    if (url != null) result.url = url;
    if (adult != null) result.adult = adult;
    return result;
  }

  Community._();

  factory Community.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Community.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Community',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'account_id')
    ..aOS(4, _omitFieldNames ? '' : 'created_at')
    ..aOS(5, _omitFieldNames ? '' : 'updated_at')
    ..aOS(6, _omitFieldNames ? '' : 'mimetype')
    ..aOS(7, _omitFieldNames ? '' : 'domain')
    ..aOS(8, _omitFieldNames ? '' : 'description')
    ..aOS(9, _omitFieldNames ? '' : 'entropy')
    ..a<$fixnum.Int64>(10, _omitFieldNames ? '' : 'bytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(11, _omitFieldNames ? '' : 'subscribed_at')
    ..aE<PublishMode>(12, _omitFieldNames ? '' : 'default_publish_mode',
        enumValues: PublishMode.values)
    ..aOB(13, _omitFieldNames ? '' : 'hidden')
    ..aOS(14, _omitFieldNames ? '' : 'url')
    ..aOB(15, _omitFieldNames ? '' : 'adult')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Community clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Community copyWith(void Function(Community) updates) =>
      super.copyWith((message) => updates(message as Community)) as Community;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Community create() => Community._();
  @$core.override
  Community createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Community getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Community>(create);
  static Community? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get accountId => $_getSZ(1);
  @$pb.TagNumber(2)
  set accountId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAccountId() => $_has(1);
  @$pb.TagNumber(2)
  void clearAccountId() => $_clearField(2);

  @$pb.TagNumber(4)
  $core.String get createdAt => $_getSZ(2);
  @$pb.TagNumber(4)
  set createdAt($core.String value) => $_setString(2, value);
  @$pb.TagNumber(4)
  $core.bool hasCreatedAt() => $_has(2);
  @$pb.TagNumber(4)
  void clearCreatedAt() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get updatedAt => $_getSZ(3);
  @$pb.TagNumber(5)
  set updatedAt($core.String value) => $_setString(3, value);
  @$pb.TagNumber(5)
  $core.bool hasUpdatedAt() => $_has(3);
  @$pb.TagNumber(5)
  void clearUpdatedAt() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get mimetype => $_getSZ(4);
  @$pb.TagNumber(6)
  set mimetype($core.String value) => $_setString(4, value);
  @$pb.TagNumber(6)
  $core.bool hasMimetype() => $_has(4);
  @$pb.TagNumber(6)
  void clearMimetype() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get domain => $_getSZ(5);
  @$pb.TagNumber(7)
  set domain($core.String value) => $_setString(5, value);
  @$pb.TagNumber(7)
  $core.bool hasDomain() => $_has(5);
  @$pb.TagNumber(7)
  void clearDomain() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get description => $_getSZ(6);
  @$pb.TagNumber(8)
  set description($core.String value) => $_setString(6, value);
  @$pb.TagNumber(8)
  $core.bool hasDescription() => $_has(6);
  @$pb.TagNumber(8)
  void clearDescription() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.String get entropy => $_getSZ(7);
  @$pb.TagNumber(9)
  set entropy($core.String value) => $_setString(7, value);
  @$pb.TagNumber(9)
  $core.bool hasEntropy() => $_has(7);
  @$pb.TagNumber(9)
  void clearEntropy() => $_clearField(9);

  @$pb.TagNumber(10)
  $fixnum.Int64 get bytes => $_getI64(8);
  @$pb.TagNumber(10)
  set bytes($fixnum.Int64 value) => $_setInt64(8, value);
  @$pb.TagNumber(10)
  $core.bool hasBytes() => $_has(8);
  @$pb.TagNumber(10)
  void clearBytes() => $_clearField(10);

  @$pb.TagNumber(11)
  $core.String get subscribedAt => $_getSZ(9);
  @$pb.TagNumber(11)
  set subscribedAt($core.String value) => $_setString(9, value);
  @$pb.TagNumber(11)
  $core.bool hasSubscribedAt() => $_has(9);
  @$pb.TagNumber(11)
  void clearSubscribedAt() => $_clearField(11);

  @$pb.TagNumber(12)
  PublishMode get defaultPublishMode => $_getN(10);
  @$pb.TagNumber(12)
  set defaultPublishMode(PublishMode value) => $_setField(12, value);
  @$pb.TagNumber(12)
  $core.bool hasDefaultPublishMode() => $_has(10);
  @$pb.TagNumber(12)
  void clearDefaultPublishMode() => $_clearField(12);

  @$pb.TagNumber(13)
  $core.bool get hidden => $_getBF(11);
  @$pb.TagNumber(13)
  set hidden($core.bool value) => $_setBool(11, value);
  @$pb.TagNumber(13)
  $core.bool hasHidden() => $_has(11);
  @$pb.TagNumber(13)
  void clearHidden() => $_clearField(13);

  @$pb.TagNumber(14)
  $core.String get url => $_getSZ(12);
  @$pb.TagNumber(14)
  set url($core.String value) => $_setString(12, value);
  @$pb.TagNumber(14)
  $core.bool hasUrl() => $_has(12);
  @$pb.TagNumber(14)
  void clearUrl() => $_clearField(14);

  @$pb.TagNumber(15)
  $core.bool get adult => $_getBF(13);
  @$pb.TagNumber(15)
  set adult($core.bool value) => $_setBool(13, value);
  @$pb.TagNumber(15)
  $core.bool hasAdult() => $_has(13);
  @$pb.TagNumber(15)
  void clearAdult() => $_clearField(15);
}

class CommunitySearchRequest extends $pb.GeneratedMessage {
  factory CommunitySearchRequest({
    $core.String? query,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (query != null) result.query = query;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  CommunitySearchRequest._();

  factory CommunitySearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunitySearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunitySearchRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunitySearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunitySearchRequest copyWith(
          void Function(CommunitySearchRequest) updates) =>
      super.copyWith((message) => updates(message as CommunitySearchRequest))
          as CommunitySearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunitySearchRequest create() => CommunitySearchRequest._();
  @$core.override
  CommunitySearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunitySearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunitySearchRequest>(create);
  static CommunitySearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get offset => $_getI64(1);
  @$pb.TagNumber(2)
  set offset($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasOffset() => $_has(1);
  @$pb.TagNumber(2)
  void clearOffset() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get limit => $_getI64(2);
  @$pb.TagNumber(3)
  set limit($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(3)
  void clearLimit() => $_clearField(3);
}

class CommunitySearchResponse extends $pb.GeneratedMessage {
  factory CommunitySearchResponse({
    CommunitySearchRequest? next,
    $core.Iterable<Community>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  CommunitySearchResponse._();

  factory CommunitySearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunitySearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunitySearchResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<CommunitySearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: CommunitySearchRequest.create)
    ..pPM<Community>(2, _omitFieldNames ? '' : 'items',
        subBuilder: Community.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunitySearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunitySearchResponse copyWith(
          void Function(CommunitySearchResponse) updates) =>
      super.copyWith((message) => updates(message as CommunitySearchResponse))
          as CommunitySearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunitySearchResponse create() => CommunitySearchResponse._();
  @$core.override
  CommunitySearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunitySearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunitySearchResponse>(create);
  static CommunitySearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  CommunitySearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(CommunitySearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  CommunitySearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Community> get items => $_getList(1);
}

class CommunityCreateRequest extends $pb.GeneratedMessage {
  factory CommunityCreateRequest({
    Community? community,
  }) {
    final result = create();
    if (community != null) result.community = community;
    return result;
  }

  CommunityCreateRequest._();

  factory CommunityCreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityCreateRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<Community>(1, _omitFieldNames ? '' : 'community',
        subBuilder: Community.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityCreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityCreateRequest copyWith(
          void Function(CommunityCreateRequest) updates) =>
      super.copyWith((message) => updates(message as CommunityCreateRequest))
          as CommunityCreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityCreateRequest create() => CommunityCreateRequest._();
  @$core.override
  CommunityCreateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityCreateRequest>(create);
  static CommunityCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Community get community => $_getN(0);
  @$pb.TagNumber(1)
  set community(Community value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunity() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunity() => $_clearField(1);
  @$pb.TagNumber(1)
  Community ensureCommunity() => $_ensure(0);
}

class CommunityCreateResponse extends $pb.GeneratedMessage {
  factory CommunityCreateResponse({
    Community? community,
  }) {
    final result = create();
    if (community != null) result.community = community;
    return result;
  }

  CommunityCreateResponse._();

  factory CommunityCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityCreateResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<Community>(1, _omitFieldNames ? '' : 'community',
        subBuilder: Community.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityCreateResponse copyWith(
          void Function(CommunityCreateResponse) updates) =>
      super.copyWith((message) => updates(message as CommunityCreateResponse))
          as CommunityCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityCreateResponse create() => CommunityCreateResponse._();
  @$core.override
  CommunityCreateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityCreateResponse>(create);
  static CommunityCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Community get community => $_getN(0);
  @$pb.TagNumber(1)
  set community(Community value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunity() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunity() => $_clearField(1);
  @$pb.TagNumber(1)
  Community ensureCommunity() => $_ensure(0);
}

class CommunityFindRequest extends $pb.GeneratedMessage {
  factory CommunityFindRequest() => create();

  CommunityFindRequest._();

  factory CommunityFindRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityFindRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityFindRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityFindRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityFindRequest copyWith(void Function(CommunityFindRequest) updates) =>
      super.copyWith((message) => updates(message as CommunityFindRequest))
          as CommunityFindRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityFindRequest create() => CommunityFindRequest._();
  @$core.override
  CommunityFindRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityFindRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityFindRequest>(create);
  static CommunityFindRequest? _defaultInstance;
}

class CommunityFindResponse extends $pb.GeneratedMessage {
  factory CommunityFindResponse({
    Community? community,
  }) {
    final result = create();
    if (community != null) result.community = community;
    return result;
  }

  CommunityFindResponse._();

  factory CommunityFindResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityFindResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityFindResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<Community>(1, _omitFieldNames ? '' : 'community',
        subBuilder: Community.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityFindResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityFindResponse copyWith(
          void Function(CommunityFindResponse) updates) =>
      super.copyWith((message) => updates(message as CommunityFindResponse))
          as CommunityFindResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityFindResponse create() => CommunityFindResponse._();
  @$core.override
  CommunityFindResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityFindResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityFindResponse>(create);
  static CommunityFindResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Community get community => $_getN(0);
  @$pb.TagNumber(1)
  set community(Community value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunity() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunity() => $_clearField(1);
  @$pb.TagNumber(1)
  Community ensureCommunity() => $_ensure(0);
}

class CommunityUploadRequest extends $pb.GeneratedMessage {
  factory CommunityUploadRequest({
    Community? community,
  }) {
    final result = create();
    if (community != null) result.community = community;
    return result;
  }

  CommunityUploadRequest._();

  factory CommunityUploadRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityUploadRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityUploadRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<Community>(1, _omitFieldNames ? '' : 'community',
        subBuilder: Community.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityUploadRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityUploadRequest copyWith(
          void Function(CommunityUploadRequest) updates) =>
      super.copyWith((message) => updates(message as CommunityUploadRequest))
          as CommunityUploadRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityUploadRequest create() => CommunityUploadRequest._();
  @$core.override
  CommunityUploadRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityUploadRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityUploadRequest>(create);
  static CommunityUploadRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Community get community => $_getN(0);
  @$pb.TagNumber(1)
  set community(Community value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunity() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunity() => $_clearField(1);
  @$pb.TagNumber(1)
  Community ensureCommunity() => $_ensure(0);
}

class CommunityUploadResponse extends $pb.GeneratedMessage {
  factory CommunityUploadResponse({
    Community? community,
  }) {
    final result = create();
    if (community != null) result.community = community;
    return result;
  }

  CommunityUploadResponse._();

  factory CommunityUploadResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityUploadResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityUploadResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<Community>(1, _omitFieldNames ? '' : 'community',
        subBuilder: Community.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityUploadResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityUploadResponse copyWith(
          void Function(CommunityUploadResponse) updates) =>
      super.copyWith((message) => updates(message as CommunityUploadResponse))
          as CommunityUploadResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityUploadResponse create() => CommunityUploadResponse._();
  @$core.override
  CommunityUploadResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityUploadResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityUploadResponse>(create);
  static CommunityUploadResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Community get community => $_getN(0);
  @$pb.TagNumber(1)
  set community(Community value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunity() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunity() => $_clearField(1);
  @$pb.TagNumber(1)
  Community ensureCommunity() => $_ensure(0);
}

class CommunityDeleteRequest extends $pb.GeneratedMessage {
  factory CommunityDeleteRequest() => create();

  CommunityDeleteRequest._();

  factory CommunityDeleteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityDeleteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityDeleteRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityDeleteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityDeleteRequest copyWith(
          void Function(CommunityDeleteRequest) updates) =>
      super.copyWith((message) => updates(message as CommunityDeleteRequest))
          as CommunityDeleteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityDeleteRequest create() => CommunityDeleteRequest._();
  @$core.override
  CommunityDeleteRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityDeleteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityDeleteRequest>(create);
  static CommunityDeleteRequest? _defaultInstance;
}

class CommunityDeleteResponse extends $pb.GeneratedMessage {
  factory CommunityDeleteResponse({
    Community? community,
  }) {
    final result = create();
    if (community != null) result.community = community;
    return result;
  }

  CommunityDeleteResponse._();

  factory CommunityDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityDeleteResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<Community>(1, _omitFieldNames ? '' : 'community',
        subBuilder: Community.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityDeleteResponse copyWith(
          void Function(CommunityDeleteResponse) updates) =>
      super.copyWith((message) => updates(message as CommunityDeleteResponse))
          as CommunityDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityDeleteResponse create() => CommunityDeleteResponse._();
  @$core.override
  CommunityDeleteResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityDeleteResponse>(create);
  static CommunityDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Community get community => $_getN(0);
  @$pb.TagNumber(1)
  set community(Community value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunity() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunity() => $_clearField(1);
  @$pb.TagNumber(1)
  Community ensureCommunity() => $_ensure(0);
}

class CommunityUpdateRequest extends $pb.GeneratedMessage {
  factory CommunityUpdateRequest({
    Community? community,
  }) {
    final result = create();
    if (community != null) result.community = community;
    return result;
  }

  CommunityUpdateRequest._();

  factory CommunityUpdateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityUpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityUpdateRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<Community>(1, _omitFieldNames ? '' : 'community',
        subBuilder: Community.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityUpdateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityUpdateRequest copyWith(
          void Function(CommunityUpdateRequest) updates) =>
      super.copyWith((message) => updates(message as CommunityUpdateRequest))
          as CommunityUpdateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityUpdateRequest create() => CommunityUpdateRequest._();
  @$core.override
  CommunityUpdateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityUpdateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityUpdateRequest>(create);
  static CommunityUpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Community get community => $_getN(0);
  @$pb.TagNumber(1)
  set community(Community value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunity() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunity() => $_clearField(1);
  @$pb.TagNumber(1)
  Community ensureCommunity() => $_ensure(0);
}

class CommunityUpdateResponse extends $pb.GeneratedMessage {
  factory CommunityUpdateResponse({
    Community? community,
  }) {
    final result = create();
    if (community != null) result.community = community;
    return result;
  }

  CommunityUpdateResponse._();

  factory CommunityUpdateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityUpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityUpdateResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<Community>(1, _omitFieldNames ? '' : 'community',
        subBuilder: Community.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityUpdateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityUpdateResponse copyWith(
          void Function(CommunityUpdateResponse) updates) =>
      super.copyWith((message) => updates(message as CommunityUpdateResponse))
          as CommunityUpdateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityUpdateResponse create() => CommunityUpdateResponse._();
  @$core.override
  CommunityUpdateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityUpdateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityUpdateResponse>(create);
  static CommunityUpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Community get community => $_getN(0);
  @$pb.TagNumber(1)
  set community(Community value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunity() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunity() => $_clearField(1);
  @$pb.TagNumber(1)
  Community ensureCommunity() => $_ensure(0);
}

class PublishedContent extends $pb.GeneratedMessage {
  factory PublishedContent({
    $core.String? id,
    $core.String? communityId,
    $core.String? knownMediaId,
    $core.String? magnetUri,
    $core.String? publishedAt,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? archivedId,
    $core.String? title,
    $core.String? description,
    $core.String? mimetype,
    $core.String? libraryId,
    $core.String? oauthGoogleId,
    $core.String? encryptionSeed,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (communityId != null) result.communityId = communityId;
    if (knownMediaId != null) result.knownMediaId = knownMediaId;
    if (magnetUri != null) result.magnetUri = magnetUri;
    if (publishedAt != null) result.publishedAt = publishedAt;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (archivedId != null) result.archivedId = archivedId;
    if (title != null) result.title = title;
    if (description != null) result.description = description;
    if (mimetype != null) result.mimetype = mimetype;
    if (libraryId != null) result.libraryId = libraryId;
    if (oauthGoogleId != null) result.oauthGoogleId = oauthGoogleId;
    if (encryptionSeed != null) result.encryptionSeed = encryptionSeed;
    return result;
  }

  PublishedContent._();

  factory PublishedContent.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishedContent.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishedContent',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'community_id')
    ..aOS(3, _omitFieldNames ? '' : 'known_media_id')
    ..aOS(4, _omitFieldNames ? '' : 'magnet_uri')
    ..aOS(5, _omitFieldNames ? '' : 'published_at')
    ..aOS(6, _omitFieldNames ? '' : 'created_at')
    ..aOS(7, _omitFieldNames ? '' : 'updated_at')
    ..aOS(8, _omitFieldNames ? '' : 'archived_id')
    ..aOS(9, _omitFieldNames ? '' : 'title')
    ..aOS(10, _omitFieldNames ? '' : 'description')
    ..aOS(11, _omitFieldNames ? '' : 'mimetype')
    ..aOS(12, _omitFieldNames ? '' : 'library_id')
    ..aOS(13, _omitFieldNames ? '' : 'oauth_google_id')
    ..aOS(14, _omitFieldNames ? '' : 'encryption_seed')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContent clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContent copyWith(void Function(PublishedContent) updates) =>
      super.copyWith((message) => updates(message as PublishedContent))
          as PublishedContent;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishedContent create() => PublishedContent._();
  @$core.override
  PublishedContent createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishedContent getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishedContent>(create);
  static PublishedContent? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get communityId => $_getSZ(1);
  @$pb.TagNumber(2)
  set communityId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasCommunityId() => $_has(1);
  @$pb.TagNumber(2)
  void clearCommunityId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get knownMediaId => $_getSZ(2);
  @$pb.TagNumber(3)
  set knownMediaId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasKnownMediaId() => $_has(2);
  @$pb.TagNumber(3)
  void clearKnownMediaId() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get magnetUri => $_getSZ(3);
  @$pb.TagNumber(4)
  set magnetUri($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasMagnetUri() => $_has(3);
  @$pb.TagNumber(4)
  void clearMagnetUri() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get publishedAt => $_getSZ(4);
  @$pb.TagNumber(5)
  set publishedAt($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasPublishedAt() => $_has(4);
  @$pb.TagNumber(5)
  void clearPublishedAt() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get createdAt => $_getSZ(5);
  @$pb.TagNumber(6)
  set createdAt($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasCreatedAt() => $_has(5);
  @$pb.TagNumber(6)
  void clearCreatedAt() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get updatedAt => $_getSZ(6);
  @$pb.TagNumber(7)
  set updatedAt($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasUpdatedAt() => $_has(6);
  @$pb.TagNumber(7)
  void clearUpdatedAt() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get archivedId => $_getSZ(7);
  @$pb.TagNumber(8)
  set archivedId($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasArchivedId() => $_has(7);
  @$pb.TagNumber(8)
  void clearArchivedId() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.String get title => $_getSZ(8);
  @$pb.TagNumber(9)
  set title($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasTitle() => $_has(8);
  @$pb.TagNumber(9)
  void clearTitle() => $_clearField(9);

  @$pb.TagNumber(10)
  $core.String get description => $_getSZ(9);
  @$pb.TagNumber(10)
  set description($core.String value) => $_setString(9, value);
  @$pb.TagNumber(10)
  $core.bool hasDescription() => $_has(9);
  @$pb.TagNumber(10)
  void clearDescription() => $_clearField(10);

  @$pb.TagNumber(11)
  $core.String get mimetype => $_getSZ(10);
  @$pb.TagNumber(11)
  set mimetype($core.String value) => $_setString(10, value);
  @$pb.TagNumber(11)
  $core.bool hasMimetype() => $_has(10);
  @$pb.TagNumber(11)
  void clearMimetype() => $_clearField(11);

  @$pb.TagNumber(12)
  $core.String get libraryId => $_getSZ(11);
  @$pb.TagNumber(12)
  set libraryId($core.String value) => $_setString(11, value);
  @$pb.TagNumber(12)
  $core.bool hasLibraryId() => $_has(11);
  @$pb.TagNumber(12)
  void clearLibraryId() => $_clearField(12);

  @$pb.TagNumber(13)
  $core.String get oauthGoogleId => $_getSZ(12);
  @$pb.TagNumber(13)
  set oauthGoogleId($core.String value) => $_setString(12, value);
  @$pb.TagNumber(13)
  $core.bool hasOauthGoogleId() => $_has(12);
  @$pb.TagNumber(13)
  void clearOauthGoogleId() => $_clearField(13);

  @$pb.TagNumber(14)
  $core.String get encryptionSeed => $_getSZ(13);
  @$pb.TagNumber(14)
  set encryptionSeed($core.String value) => $_setString(13, value);
  @$pb.TagNumber(14)
  $core.bool hasEncryptionSeed() => $_has(13);
  @$pb.TagNumber(14)
  void clearEncryptionSeed() => $_clearField(14);
}

class PublishContentRequest extends $pb.GeneratedMessage {
  factory PublishContentRequest({
    PublishedContent? publishedContent,
    PublishMode? publishMode,
  }) {
    final result = create();
    if (publishedContent != null) result.publishedContent = publishedContent;
    if (publishMode != null) result.publishMode = publishMode;
    return result;
  }

  PublishContentRequest._();

  factory PublishContentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishContentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishContentRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<PublishedContent>(1, _omitFieldNames ? '' : 'published_content',
        subBuilder: PublishedContent.create)
    ..aE<PublishMode>(2, _omitFieldNames ? '' : 'publish_mode',
        enumValues: PublishMode.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishContentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishContentRequest copyWith(
          void Function(PublishContentRequest) updates) =>
      super.copyWith((message) => updates(message as PublishContentRequest))
          as PublishContentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishContentRequest create() => PublishContentRequest._();
  @$core.override
  PublishContentRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishContentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishContentRequest>(create);
  static PublishContentRequest? _defaultInstance;

  @$pb.TagNumber(1)
  PublishedContent get publishedContent => $_getN(0);
  @$pb.TagNumber(1)
  set publishedContent(PublishedContent value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPublishedContent() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublishedContent() => $_clearField(1);
  @$pb.TagNumber(1)
  PublishedContent ensurePublishedContent() => $_ensure(0);

  @$pb.TagNumber(2)
  PublishMode get publishMode => $_getN(1);
  @$pb.TagNumber(2)
  set publishMode(PublishMode value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasPublishMode() => $_has(1);
  @$pb.TagNumber(2)
  void clearPublishMode() => $_clearField(2);
}

class PublishContentResponse extends $pb.GeneratedMessage {
  factory PublishContentResponse({
    PublishedContent? publishedContent,
  }) {
    final result = create();
    if (publishedContent != null) result.publishedContent = publishedContent;
    return result;
  }

  PublishContentResponse._();

  factory PublishContentResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishContentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishContentResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<PublishedContent>(1, _omitFieldNames ? '' : 'published_content',
        subBuilder: PublishedContent.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishContentResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishContentResponse copyWith(
          void Function(PublishContentResponse) updates) =>
      super.copyWith((message) => updates(message as PublishContentResponse))
          as PublishContentResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishContentResponse create() => PublishContentResponse._();
  @$core.override
  PublishContentResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishContentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishContentResponse>(create);
  static PublishContentResponse? _defaultInstance;

  @$pb.TagNumber(1)
  PublishedContent get publishedContent => $_getN(0);
  @$pb.TagNumber(1)
  set publishedContent(PublishedContent value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPublishedContent() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublishedContent() => $_clearField(1);
  @$pb.TagNumber(1)
  PublishedContent ensurePublishedContent() => $_ensure(0);
}

class PublishedContentListRequest extends $pb.GeneratedMessage {
  factory PublishedContentListRequest({
    $core.String? communityId,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (communityId != null) result.communityId = communityId;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  PublishedContentListRequest._();

  factory PublishedContentListRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishedContentListRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishedContentListRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'community_id')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContentListRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContentListRequest copyWith(
          void Function(PublishedContentListRequest) updates) =>
      super.copyWith(
              (message) => updates(message as PublishedContentListRequest))
          as PublishedContentListRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishedContentListRequest create() =>
      PublishedContentListRequest._();
  @$core.override
  PublishedContentListRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishedContentListRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishedContentListRequest>(create);
  static PublishedContentListRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get communityId => $_getSZ(0);
  @$pb.TagNumber(1)
  set communityId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunityId() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunityId() => $_clearField(1);

  @$pb.TagNumber(900)
  $fixnum.Int64 get offset => $_getI64(1);
  @$pb.TagNumber(900)
  set offset($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(900)
  $core.bool hasOffset() => $_has(1);
  @$pb.TagNumber(900)
  void clearOffset() => $_clearField(900);

  @$pb.TagNumber(901)
  $fixnum.Int64 get limit => $_getI64(2);
  @$pb.TagNumber(901)
  set limit($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(901)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(901)
  void clearLimit() => $_clearField(901);
}

class PublishedContentListResponse extends $pb.GeneratedMessage {
  factory PublishedContentListResponse({
    Community? community,
    PublishedContentListRequest? next,
    $core.Iterable<PublishedContent>? items,
  }) {
    final result = create();
    if (community != null) result.community = community;
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  PublishedContentListResponse._();

  factory PublishedContentListResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishedContentListResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishedContentListResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<Community>(1, _omitFieldNames ? '' : 'community',
        subBuilder: Community.create)
    ..aOM<PublishedContentListRequest>(2, _omitFieldNames ? '' : 'next',
        subBuilder: PublishedContentListRequest.create)
    ..pPM<PublishedContent>(3, _omitFieldNames ? '' : 'items',
        subBuilder: PublishedContent.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContentListResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContentListResponse copyWith(
          void Function(PublishedContentListResponse) updates) =>
      super.copyWith(
              (message) => updates(message as PublishedContentListResponse))
          as PublishedContentListResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishedContentListResponse create() =>
      PublishedContentListResponse._();
  @$core.override
  PublishedContentListResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishedContentListResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishedContentListResponse>(create);
  static PublishedContentListResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Community get community => $_getN(0);
  @$pb.TagNumber(1)
  set community(Community value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunity() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunity() => $_clearField(1);
  @$pb.TagNumber(1)
  Community ensureCommunity() => $_ensure(0);

  @$pb.TagNumber(2)
  PublishedContentListRequest get next => $_getN(1);
  @$pb.TagNumber(2)
  set next(PublishedContentListRequest value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasNext() => $_has(1);
  @$pb.TagNumber(2)
  void clearNext() => $_clearField(2);
  @$pb.TagNumber(2)
  PublishedContentListRequest ensureNext() => $_ensure(1);

  @$pb.TagNumber(3)
  $pb.PbList<PublishedContent> get items => $_getList(2);
}

class CommunityMetric extends $pb.GeneratedMessage {
  factory CommunityMetric({
    $core.String? id,
    $core.String? communityId,
    $core.String? periodStart,
    $core.String? periodEnd,
    $core.int? subscribers,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (communityId != null) result.communityId = communityId;
    if (periodStart != null) result.periodStart = periodStart;
    if (periodEnd != null) result.periodEnd = periodEnd;
    if (subscribers != null) result.subscribers = subscribers;
    return result;
  }

  CommunityMetric._();

  factory CommunityMetric.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityMetric.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityMetric',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'community_id')
    ..aOS(3, _omitFieldNames ? '' : 'period_start')
    ..aOS(4, _omitFieldNames ? '' : 'period_end')
    ..aI(5, _omitFieldNames ? '' : 'subscribers',
        fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityMetric clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityMetric copyWith(void Function(CommunityMetric) updates) =>
      super.copyWith((message) => updates(message as CommunityMetric))
          as CommunityMetric;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityMetric create() => CommunityMetric._();
  @$core.override
  CommunityMetric createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityMetric getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityMetric>(create);
  static CommunityMetric? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get communityId => $_getSZ(1);
  @$pb.TagNumber(2)
  set communityId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasCommunityId() => $_has(1);
  @$pb.TagNumber(2)
  void clearCommunityId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get periodStart => $_getSZ(2);
  @$pb.TagNumber(3)
  set periodStart($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPeriodStart() => $_has(2);
  @$pb.TagNumber(3)
  void clearPeriodStart() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get periodEnd => $_getSZ(3);
  @$pb.TagNumber(4)
  set periodEnd($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasPeriodEnd() => $_has(3);
  @$pb.TagNumber(4)
  void clearPeriodEnd() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.int get subscribers => $_getIZ(4);
  @$pb.TagNumber(5)
  set subscribers($core.int value) => $_setUnsignedInt32(4, value);
  @$pb.TagNumber(5)
  $core.bool hasSubscribers() => $_has(4);
  @$pb.TagNumber(5)
  void clearSubscribers() => $_clearField(5);
}

class PublishedContentMetric extends $pb.GeneratedMessage {
  factory PublishedContentMetric({
    $core.String? id,
    $core.String? publishedContentId,
    $core.String? periodStart,
    $core.String? periodEnd,
    $core.int? archivers,
    $fixnum.Int64? bytes,
    $fixnum.Int64? revenue,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (publishedContentId != null)
      result.publishedContentId = publishedContentId;
    if (periodStart != null) result.periodStart = periodStart;
    if (periodEnd != null) result.periodEnd = periodEnd;
    if (archivers != null) result.archivers = archivers;
    if (bytes != null) result.bytes = bytes;
    if (revenue != null) result.revenue = revenue;
    return result;
  }

  PublishedContentMetric._();

  factory PublishedContentMetric.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishedContentMetric.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishedContentMetric',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'published_content_id')
    ..aOS(3, _omitFieldNames ? '' : 'period_start')
    ..aOS(4, _omitFieldNames ? '' : 'period_end')
    ..aI(5, _omitFieldNames ? '' : 'archivers', fieldType: $pb.PbFieldType.OU3)
    ..aInt64(6, _omitFieldNames ? '' : 'bytes')
    ..aInt64(7, _omitFieldNames ? '' : 'revenue')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContentMetric clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContentMetric copyWith(
          void Function(PublishedContentMetric) updates) =>
      super.copyWith((message) => updates(message as PublishedContentMetric))
          as PublishedContentMetric;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishedContentMetric create() => PublishedContentMetric._();
  @$core.override
  PublishedContentMetric createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishedContentMetric getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishedContentMetric>(create);
  static PublishedContentMetric? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get publishedContentId => $_getSZ(1);
  @$pb.TagNumber(2)
  set publishedContentId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPublishedContentId() => $_has(1);
  @$pb.TagNumber(2)
  void clearPublishedContentId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get periodStart => $_getSZ(2);
  @$pb.TagNumber(3)
  set periodStart($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPeriodStart() => $_has(2);
  @$pb.TagNumber(3)
  void clearPeriodStart() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get periodEnd => $_getSZ(3);
  @$pb.TagNumber(4)
  set periodEnd($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasPeriodEnd() => $_has(3);
  @$pb.TagNumber(4)
  void clearPeriodEnd() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.int get archivers => $_getIZ(4);
  @$pb.TagNumber(5)
  set archivers($core.int value) => $_setUnsignedInt32(4, value);
  @$pb.TagNumber(5)
  $core.bool hasArchivers() => $_has(4);
  @$pb.TagNumber(5)
  void clearArchivers() => $_clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get bytes => $_getI64(5);
  @$pb.TagNumber(6)
  set bytes($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(6)
  $core.bool hasBytes() => $_has(5);
  @$pb.TagNumber(6)
  void clearBytes() => $_clearField(6);

  @$pb.TagNumber(7)
  $fixnum.Int64 get revenue => $_getI64(6);
  @$pb.TagNumber(7)
  set revenue($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(7)
  $core.bool hasRevenue() => $_has(6);
  @$pb.TagNumber(7)
  void clearRevenue() => $_clearField(7);
}

class CommunityMetricsRequest extends $pb.GeneratedMessage {
  factory CommunityMetricsRequest({
    $core.String? communityId,
    $core.String? period,
    $core.String? startDate,
    $core.String? endDate,
  }) {
    final result = create();
    if (communityId != null) result.communityId = communityId;
    if (period != null) result.period = period;
    if (startDate != null) result.startDate = startDate;
    if (endDate != null) result.endDate = endDate;
    return result;
  }

  CommunityMetricsRequest._();

  factory CommunityMetricsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityMetricsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityMetricsRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'community_id')
    ..aOS(2, _omitFieldNames ? '' : 'period')
    ..aOS(3, _omitFieldNames ? '' : 'start_date')
    ..aOS(4, _omitFieldNames ? '' : 'end_date')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityMetricsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityMetricsRequest copyWith(
          void Function(CommunityMetricsRequest) updates) =>
      super.copyWith((message) => updates(message as CommunityMetricsRequest))
          as CommunityMetricsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityMetricsRequest create() => CommunityMetricsRequest._();
  @$core.override
  CommunityMetricsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityMetricsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityMetricsRequest>(create);
  static CommunityMetricsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get communityId => $_getSZ(0);
  @$pb.TagNumber(1)
  set communityId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunityId() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunityId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get period => $_getSZ(1);
  @$pb.TagNumber(2)
  set period($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPeriod() => $_has(1);
  @$pb.TagNumber(2)
  void clearPeriod() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get startDate => $_getSZ(2);
  @$pb.TagNumber(3)
  set startDate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasStartDate() => $_has(2);
  @$pb.TagNumber(3)
  void clearStartDate() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get endDate => $_getSZ(3);
  @$pb.TagNumber(4)
  set endDate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasEndDate() => $_has(3);
  @$pb.TagNumber(4)
  void clearEndDate() => $_clearField(4);
}

class CommunityMetricsResponse extends $pb.GeneratedMessage {
  factory CommunityMetricsResponse({
    CommunityMetric? summary,
    $core.int? totalArchivers,
    $core.Iterable<PublishedContentMetric>? items,
  }) {
    final result = create();
    if (summary != null) result.summary = summary;
    if (totalArchivers != null) result.totalArchivers = totalArchivers;
    if (items != null) result.items.addAll(items);
    return result;
  }

  CommunityMetricsResponse._();

  factory CommunityMetricsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityMetricsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityMetricsResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<CommunityMetric>(1, _omitFieldNames ? '' : 'summary',
        subBuilder: CommunityMetric.create)
    ..aI(2, _omitFieldNames ? '' : 'total_archivers')
    ..pPM<PublishedContentMetric>(3, _omitFieldNames ? '' : 'items',
        subBuilder: PublishedContentMetric.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityMetricsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityMetricsResponse copyWith(
          void Function(CommunityMetricsResponse) updates) =>
      super.copyWith((message) => updates(message as CommunityMetricsResponse))
          as CommunityMetricsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityMetricsResponse create() => CommunityMetricsResponse._();
  @$core.override
  CommunityMetricsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityMetricsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityMetricsResponse>(create);
  static CommunityMetricsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  CommunityMetric get summary => $_getN(0);
  @$pb.TagNumber(1)
  set summary(CommunityMetric value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasSummary() => $_has(0);
  @$pb.TagNumber(1)
  void clearSummary() => $_clearField(1);
  @$pb.TagNumber(1)
  CommunityMetric ensureSummary() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.int get totalArchivers => $_getIZ(1);
  @$pb.TagNumber(2)
  set totalArchivers($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTotalArchivers() => $_has(1);
  @$pb.TagNumber(2)
  void clearTotalArchivers() => $_clearField(2);

  @$pb.TagNumber(3)
  $pb.PbList<PublishedContentMetric> get items => $_getList(2);
}

class MetricsSyncRequest extends $pb.GeneratedMessage {
  factory MetricsSyncRequest({
    $core.String? communityId,
    $core.String? since,
  }) {
    final result = create();
    if (communityId != null) result.communityId = communityId;
    if (since != null) result.since = since;
    return result;
  }

  MetricsSyncRequest._();

  factory MetricsSyncRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricsSyncRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricsSyncRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'community_id')
    ..aOS(2, _omitFieldNames ? '' : 'since')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricsSyncRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricsSyncRequest copyWith(void Function(MetricsSyncRequest) updates) =>
      super.copyWith((message) => updates(message as MetricsSyncRequest))
          as MetricsSyncRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricsSyncRequest create() => MetricsSyncRequest._();
  @$core.override
  MetricsSyncRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricsSyncRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricsSyncRequest>(create);
  static MetricsSyncRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get communityId => $_getSZ(0);
  @$pb.TagNumber(1)
  set communityId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunityId() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunityId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get since => $_getSZ(1);
  @$pb.TagNumber(2)
  set since($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSince() => $_has(1);
  @$pb.TagNumber(2)
  void clearSince() => $_clearField(2);
}

class MetricsSyncResponse extends $pb.GeneratedMessage {
  factory MetricsSyncResponse({
    $core.Iterable<CommunityMetric>? communityMetrics,
    $core.Iterable<PublishedContentMetric>? contentMetrics,
    $core.String? syncedAt,
    $core.bool? complete,
  }) {
    final result = create();
    if (communityMetrics != null)
      result.communityMetrics.addAll(communityMetrics);
    if (contentMetrics != null) result.contentMetrics.addAll(contentMetrics);
    if (syncedAt != null) result.syncedAt = syncedAt;
    if (complete != null) result.complete = complete;
    return result;
  }

  MetricsSyncResponse._();

  factory MetricsSyncResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricsSyncResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricsSyncResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..pPM<CommunityMetric>(1, _omitFieldNames ? '' : 'community_metrics',
        subBuilder: CommunityMetric.create)
    ..pPM<PublishedContentMetric>(2, _omitFieldNames ? '' : 'content_metrics',
        subBuilder: PublishedContentMetric.create)
    ..aOS(3, _omitFieldNames ? '' : 'synced_at')
    ..aOB(4, _omitFieldNames ? '' : 'complete')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricsSyncResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricsSyncResponse copyWith(void Function(MetricsSyncResponse) updates) =>
      super.copyWith((message) => updates(message as MetricsSyncResponse))
          as MetricsSyncResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricsSyncResponse create() => MetricsSyncResponse._();
  @$core.override
  MetricsSyncResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricsSyncResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricsSyncResponse>(create);
  static MetricsSyncResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<CommunityMetric> get communityMetrics => $_getList(0);

  @$pb.TagNumber(2)
  $pb.PbList<PublishedContentMetric> get contentMetrics => $_getList(1);

  @$pb.TagNumber(3)
  $core.String get syncedAt => $_getSZ(2);
  @$pb.TagNumber(3)
  set syncedAt($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasSyncedAt() => $_has(2);
  @$pb.TagNumber(3)
  void clearSyncedAt() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get complete => $_getBF(3);
  @$pb.TagNumber(4)
  set complete($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasComplete() => $_has(3);
  @$pb.TagNumber(4)
  void clearComplete() => $_clearField(4);
}

class MetricsSyncProgress extends $pb.GeneratedMessage {
  factory MetricsSyncProgress({
    $core.String? status,
    $core.int? communityMetricsCount,
    $core.int? contentMetricsCount,
    $core.String? error,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (communityMetricsCount != null)
      result.communityMetricsCount = communityMetricsCount;
    if (contentMetricsCount != null)
      result.contentMetricsCount = contentMetricsCount;
    if (error != null) result.error = error;
    return result;
  }

  MetricsSyncProgress._();

  factory MetricsSyncProgress.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricsSyncProgress.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricsSyncProgress',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'status')
    ..aI(2, _omitFieldNames ? '' : 'community_metrics_count')
    ..aI(3, _omitFieldNames ? '' : 'content_metrics_count')
    ..aOS(4, _omitFieldNames ? '' : 'error')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricsSyncProgress clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricsSyncProgress copyWith(void Function(MetricsSyncProgress) updates) =>
      super.copyWith((message) => updates(message as MetricsSyncProgress))
          as MetricsSyncProgress;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricsSyncProgress create() => MetricsSyncProgress._();
  @$core.override
  MetricsSyncProgress createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricsSyncProgress getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricsSyncProgress>(create);
  static MetricsSyncProgress? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get status => $_getSZ(0);
  @$pb.TagNumber(1)
  set status($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(1)
  void clearStatus() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get communityMetricsCount => $_getIZ(1);
  @$pb.TagNumber(2)
  set communityMetricsCount($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasCommunityMetricsCount() => $_has(1);
  @$pb.TagNumber(2)
  void clearCommunityMetricsCount() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get contentMetricsCount => $_getIZ(2);
  @$pb.TagNumber(3)
  set contentMetricsCount($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasContentMetricsCount() => $_has(2);
  @$pb.TagNumber(3)
  void clearContentMetricsCount() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get error => $_getSZ(3);
  @$pb.TagNumber(4)
  set error($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasError() => $_has(3);
  @$pb.TagNumber(4)
  void clearError() => $_clearField(4);
}

class CommunitySubscribeRequest extends $pb.GeneratedMessage {
  factory CommunitySubscribeRequest() => create();

  CommunitySubscribeRequest._();

  factory CommunitySubscribeRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunitySubscribeRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunitySubscribeRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunitySubscribeRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunitySubscribeRequest copyWith(
          void Function(CommunitySubscribeRequest) updates) =>
      super.copyWith((message) => updates(message as CommunitySubscribeRequest))
          as CommunitySubscribeRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunitySubscribeRequest create() => CommunitySubscribeRequest._();
  @$core.override
  CommunitySubscribeRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunitySubscribeRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunitySubscribeRequest>(create);
  static CommunitySubscribeRequest? _defaultInstance;
}

class CommunitySubscribeResponse extends $pb.GeneratedMessage {
  factory CommunitySubscribeResponse() => create();

  CommunitySubscribeResponse._();

  factory CommunitySubscribeResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunitySubscribeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunitySubscribeResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunitySubscribeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunitySubscribeResponse copyWith(
          void Function(CommunitySubscribeResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CommunitySubscribeResponse))
          as CommunitySubscribeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunitySubscribeResponse create() => CommunitySubscribeResponse._();
  @$core.override
  CommunitySubscribeResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunitySubscribeResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunitySubscribeResponse>(create);
  static CommunitySubscribeResponse? _defaultInstance;
}

class YouTubeStatus extends $pb.GeneratedMessage {
  factory YouTubeStatus({
    $core.bool? linked,
    $core.String? id,
  }) {
    final result = create();
    if (linked != null) result.linked = linked;
    if (id != null) result.id = id;
    return result;
  }

  YouTubeStatus._();

  factory YouTubeStatus.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory YouTubeStatus.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'YouTubeStatus',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'linked')
    ..aOS(2, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  YouTubeStatus clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  YouTubeStatus copyWith(void Function(YouTubeStatus) updates) =>
      super.copyWith((message) => updates(message as YouTubeStatus))
          as YouTubeStatus;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static YouTubeStatus create() => YouTubeStatus._();
  @$core.override
  YouTubeStatus createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static YouTubeStatus getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<YouTubeStatus>(create);
  static YouTubeStatus? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get linked => $_getBF(0);
  @$pb.TagNumber(1)
  set linked($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasLinked() => $_has(0);
  @$pb.TagNumber(1)
  void clearLinked() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get id => $_getSZ(1);
  @$pb.TagNumber(2)
  set id($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasId() => $_has(1);
  @$pb.TagNumber(2)
  void clearId() => $_clearField(2);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
