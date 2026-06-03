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
    $fixnum.Int64? defaultTtl,
    $core.String? defaultLanguage,
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
    if (defaultTtl != null) result.defaultTtl = defaultTtl;
    if (defaultLanguage != null) result.defaultLanguage = defaultLanguage;
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
    ..a<$fixnum.Int64>(
        16, _omitFieldNames ? '' : 'default_ttl', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(17, _omitFieldNames ? '' : 'default_language')
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

  @$pb.TagNumber(16)
  $fixnum.Int64 get defaultTtl => $_getI64(14);
  @$pb.TagNumber(16)
  set defaultTtl($fixnum.Int64 value) => $_setInt64(14, value);
  @$pb.TagNumber(16)
  $core.bool hasDefaultTtl() => $_has(14);
  @$pb.TagNumber(16)
  void clearDefaultTtl() => $_clearField(16);

  @$pb.TagNumber(17)
  $core.String get defaultLanguage => $_getSZ(15);
  @$pb.TagNumber(17)
  set defaultLanguage($core.String value) => $_setString(15, value);
  @$pb.TagNumber(17)
  $core.bool hasDefaultLanguage() => $_has(15);
  @$pb.TagNumber(17)
  void clearDefaultLanguage() => $_clearField(17);
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
