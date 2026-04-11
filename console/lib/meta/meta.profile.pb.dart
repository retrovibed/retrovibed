// This is a generated file - do not edit.
//
// Generated from meta.profile.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'meta.profile.pbenum.dart';

class Profile extends $pb.GeneratedMessage {
  factory Profile({
    $core.String? id,
    $core.String? accountId,
    $core.String? sessionWatermark,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? disabledAt,
    $core.String? disabledManuallyAt,
    $core.String? disabledPendingApprovalAt,
    $core.String? display,
    $core.String? email,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (accountId != null) result.accountId = accountId;
    if (sessionWatermark != null) result.sessionWatermark = sessionWatermark;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (disabledAt != null) result.disabledAt = disabledAt;
    if (disabledManuallyAt != null)
      result.disabledManuallyAt = disabledManuallyAt;
    if (disabledPendingApprovalAt != null)
      result.disabledPendingApprovalAt = disabledPendingApprovalAt;
    if (display != null) result.display = display;
    if (email != null) result.email = email;
    return result;
  }

  Profile._();

  factory Profile.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Profile.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Profile',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'account_id')
    ..aOS(3, _omitFieldNames ? '' : 'session_watermark')
    ..aOS(4, _omitFieldNames ? '' : 'created_at')
    ..aOS(5, _omitFieldNames ? '' : 'updated_at')
    ..aOS(6, _omitFieldNames ? '' : 'disabled_at')
    ..aOS(7, _omitFieldNames ? '' : 'disabled_manually_at')
    ..aOS(8, _omitFieldNames ? '' : 'disabled_pending_approval_at')
    ..aOS(9, _omitFieldNames ? '' : 'display')
    ..aOS(10, _omitFieldNames ? '' : 'email')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Profile clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Profile copyWith(void Function(Profile) updates) =>
      super.copyWith((message) => updates(message as Profile)) as Profile;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Profile create() => Profile._();
  @$core.override
  Profile createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Profile getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Profile>(create);
  static Profile? _defaultInstance;

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

  @$pb.TagNumber(3)
  $core.String get sessionWatermark => $_getSZ(2);
  @$pb.TagNumber(3)
  set sessionWatermark($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasSessionWatermark() => $_has(2);
  @$pb.TagNumber(3)
  void clearSessionWatermark() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get createdAt => $_getSZ(3);
  @$pb.TagNumber(4)
  set createdAt($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasCreatedAt() => $_has(3);
  @$pb.TagNumber(4)
  void clearCreatedAt() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get updatedAt => $_getSZ(4);
  @$pb.TagNumber(5)
  set updatedAt($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasUpdatedAt() => $_has(4);
  @$pb.TagNumber(5)
  void clearUpdatedAt() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get disabledAt => $_getSZ(5);
  @$pb.TagNumber(6)
  set disabledAt($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasDisabledAt() => $_has(5);
  @$pb.TagNumber(6)
  void clearDisabledAt() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get disabledManuallyAt => $_getSZ(6);
  @$pb.TagNumber(7)
  set disabledManuallyAt($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasDisabledManuallyAt() => $_has(6);
  @$pb.TagNumber(7)
  void clearDisabledManuallyAt() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get disabledPendingApprovalAt => $_getSZ(7);
  @$pb.TagNumber(8)
  set disabledPendingApprovalAt($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasDisabledPendingApprovalAt() => $_has(7);
  @$pb.TagNumber(8)
  void clearDisabledPendingApprovalAt() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.String get display => $_getSZ(8);
  @$pb.TagNumber(9)
  set display($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasDisplay() => $_has(8);
  @$pb.TagNumber(9)
  void clearDisplay() => $_clearField(9);

  @$pb.TagNumber(10)
  $core.String get email => $_getSZ(9);
  @$pb.TagNumber(10)
  set email($core.String value) => $_setString(9, value);
  @$pb.TagNumber(10)
  $core.bool hasEmail() => $_has(9);
  @$pb.TagNumber(10)
  void clearEmail() => $_clearField(10);
}

class ProfileSearchRequest extends $pb.GeneratedMessage {
  factory ProfileSearchRequest({
    $core.String? query,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
    $core.int? status,
  }) {
    final result = create();
    if (query != null) result.query = query;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    if (status != null) result.status = status;
    return result;
  }

  ProfileSearchRequest._();

  factory ProfileSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProfileSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProfileSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aI(4, _omitFieldNames ? '' : 'status', fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileSearchRequest copyWith(void Function(ProfileSearchRequest) updates) =>
      super.copyWith((message) => updates(message as ProfileSearchRequest))
          as ProfileSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileSearchRequest create() => ProfileSearchRequest._();
  @$core.override
  ProfileSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProfileSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProfileSearchRequest>(create);
  static ProfileSearchRequest? _defaultInstance;

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

  @$pb.TagNumber(4)
  $core.int get status => $_getIZ(3);
  @$pb.TagNumber(4)
  set status($core.int value) => $_setUnsignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasStatus() => $_has(3);
  @$pb.TagNumber(4)
  void clearStatus() => $_clearField(4);
}

class ProfileSearchResponse extends $pb.GeneratedMessage {
  factory ProfileSearchResponse({
    ProfileSearchRequest? next,
    $core.Iterable<Profile>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  ProfileSearchResponse._();

  factory ProfileSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProfileSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProfileSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<ProfileSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: ProfileSearchRequest.create)
    ..pPM<Profile>(2, _omitFieldNames ? '' : 'items',
        subBuilder: Profile.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileSearchResponse copyWith(
          void Function(ProfileSearchResponse) updates) =>
      super.copyWith((message) => updates(message as ProfileSearchResponse))
          as ProfileSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileSearchResponse create() => ProfileSearchResponse._();
  @$core.override
  ProfileSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProfileSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProfileSearchResponse>(create);
  static ProfileSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ProfileSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(ProfileSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  ProfileSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Profile> get items => $_getList(1);
}

class ProfileLookupRequest extends $pb.GeneratedMessage {
  factory ProfileLookupRequest() => create();

  ProfileLookupRequest._();

  factory ProfileLookupRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProfileLookupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProfileLookupRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileLookupRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileLookupRequest copyWith(void Function(ProfileLookupRequest) updates) =>
      super.copyWith((message) => updates(message as ProfileLookupRequest))
          as ProfileLookupRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileLookupRequest create() => ProfileLookupRequest._();
  @$core.override
  ProfileLookupRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProfileLookupRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProfileLookupRequest>(create);
  static ProfileLookupRequest? _defaultInstance;
}

class ProfileLookupResponse extends $pb.GeneratedMessage {
  factory ProfileLookupResponse({
    Profile? profile,
  }) {
    final result = create();
    if (profile != null) result.profile = profile;
    return result;
  }

  ProfileLookupResponse._();

  factory ProfileLookupResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProfileLookupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProfileLookupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile',
        subBuilder: Profile.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileLookupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileLookupResponse copyWith(
          void Function(ProfileLookupResponse) updates) =>
      super.copyWith((message) => updates(message as ProfileLookupResponse))
          as ProfileLookupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileLookupResponse create() => ProfileLookupResponse._();
  @$core.override
  ProfileLookupResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProfileLookupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProfileLookupResponse>(create);
  static ProfileLookupResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => $_clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);
}

class ProfileUpdateRequest extends $pb.GeneratedMessage {
  factory ProfileUpdateRequest({
    Profile? profile,
  }) {
    final result = create();
    if (profile != null) result.profile = profile;
    return result;
  }

  ProfileUpdateRequest._();

  factory ProfileUpdateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProfileUpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProfileUpdateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile',
        subBuilder: Profile.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileUpdateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileUpdateRequest copyWith(void Function(ProfileUpdateRequest) updates) =>
      super.copyWith((message) => updates(message as ProfileUpdateRequest))
          as ProfileUpdateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileUpdateRequest create() => ProfileUpdateRequest._();
  @$core.override
  ProfileUpdateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProfileUpdateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProfileUpdateRequest>(create);
  static ProfileUpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => $_clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);
}

class ProfileUpdateResponse extends $pb.GeneratedMessage {
  factory ProfileUpdateResponse({
    Profile? profile,
  }) {
    final result = create();
    if (profile != null) result.profile = profile;
    return result;
  }

  ProfileUpdateResponse._();

  factory ProfileUpdateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProfileUpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProfileUpdateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile',
        subBuilder: Profile.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileUpdateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileUpdateResponse copyWith(
          void Function(ProfileUpdateResponse) updates) =>
      super.copyWith((message) => updates(message as ProfileUpdateResponse))
          as ProfileUpdateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileUpdateResponse create() => ProfileUpdateResponse._();
  @$core.override
  ProfileUpdateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProfileUpdateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProfileUpdateResponse>(create);
  static ProfileUpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => $_clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);
}

class ProfileDisableRequest extends $pb.GeneratedMessage {
  factory ProfileDisableRequest({
    Profile? profile,
  }) {
    final result = create();
    if (profile != null) result.profile = profile;
    return result;
  }

  ProfileDisableRequest._();

  factory ProfileDisableRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProfileDisableRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProfileDisableRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile',
        subBuilder: Profile.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileDisableRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileDisableRequest copyWith(
          void Function(ProfileDisableRequest) updates) =>
      super.copyWith((message) => updates(message as ProfileDisableRequest))
          as ProfileDisableRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileDisableRequest create() => ProfileDisableRequest._();
  @$core.override
  ProfileDisableRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProfileDisableRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProfileDisableRequest>(create);
  static ProfileDisableRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => $_clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);
}

class ProfileDisableResponse extends $pb.GeneratedMessage {
  factory ProfileDisableResponse({
    Profile? profile,
  }) {
    final result = create();
    if (profile != null) result.profile = profile;
    return result;
  }

  ProfileDisableResponse._();

  factory ProfileDisableResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProfileDisableResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProfileDisableResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile',
        subBuilder: Profile.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileDisableResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileDisableResponse copyWith(
          void Function(ProfileDisableResponse) updates) =>
      super.copyWith((message) => updates(message as ProfileDisableResponse))
          as ProfileDisableResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileDisableResponse create() => ProfileDisableResponse._();
  @$core.override
  ProfileDisableResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProfileDisableResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProfileDisableResponse>(create);
  static ProfileDisableResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => $_clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);
}

class ProfileCreateRequest extends $pb.GeneratedMessage {
  factory ProfileCreateRequest({
    Profile? profile,
    $core.String? publicKey,
  }) {
    final result = create();
    if (profile != null) result.profile = profile;
    if (publicKey != null) result.publicKey = publicKey;
    return result;
  }

  ProfileCreateRequest._();

  factory ProfileCreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProfileCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProfileCreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile',
        subBuilder: Profile.create)
    ..aOS(2, _omitFieldNames ? '' : 'public_key')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileCreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileCreateRequest copyWith(void Function(ProfileCreateRequest) updates) =>
      super.copyWith((message) => updates(message as ProfileCreateRequest))
          as ProfileCreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileCreateRequest create() => ProfileCreateRequest._();
  @$core.override
  ProfileCreateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProfileCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProfileCreateRequest>(create);
  static ProfileCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => $_clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.String get publicKey => $_getSZ(1);
  @$pb.TagNumber(2)
  set publicKey($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPublicKey() => $_has(1);
  @$pb.TagNumber(2)
  void clearPublicKey() => $_clearField(2);
}

class ProfileCreateResponse extends $pb.GeneratedMessage {
  factory ProfileCreateResponse({
    Profile? profile,
  }) {
    final result = create();
    if (profile != null) result.profile = profile;
    return result;
  }

  ProfileCreateResponse._();

  factory ProfileCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProfileCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProfileCreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile',
        subBuilder: Profile.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileCreateResponse copyWith(
          void Function(ProfileCreateResponse) updates) =>
      super.copyWith((message) => updates(message as ProfileCreateResponse))
          as ProfileCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileCreateResponse create() => ProfileCreateResponse._();
  @$core.override
  ProfileCreateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProfileCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProfileCreateResponse>(create);
  static ProfileCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => $_clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
