//
//  Generated code. Do not modify.
//  source: meta.profile.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

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
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (accountId != null) {
      $result.accountId = accountId;
    }
    if (sessionWatermark != null) {
      $result.sessionWatermark = sessionWatermark;
    }
    if (createdAt != null) {
      $result.createdAt = createdAt;
    }
    if (updatedAt != null) {
      $result.updatedAt = updatedAt;
    }
    if (disabledAt != null) {
      $result.disabledAt = disabledAt;
    }
    if (disabledManuallyAt != null) {
      $result.disabledManuallyAt = disabledManuallyAt;
    }
    if (disabledPendingApprovalAt != null) {
      $result.disabledPendingApprovalAt = disabledPendingApprovalAt;
    }
    if (display != null) {
      $result.display = display;
    }
    if (email != null) {
      $result.email = email;
    }
    return $result;
  }
  Profile._() : super();
  factory Profile.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory Profile.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'Profile', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
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
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  Profile clone() => Profile()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  Profile copyWith(void Function(Profile) updates) => super.copyWith((message) => updates(message as Profile)) as Profile;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Profile create() => Profile._();
  Profile createEmptyInstance() => create();
  static $pb.PbList<Profile> createRepeated() => $pb.PbList<Profile>();
  @$core.pragma('dart2js:noInline')
  static Profile getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Profile>(create);
  static Profile? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get accountId => $_getSZ(1);
  @$pb.TagNumber(2)
  set accountId($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasAccountId() => $_has(1);
  @$pb.TagNumber(2)
  void clearAccountId() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get sessionWatermark => $_getSZ(2);
  @$pb.TagNumber(3)
  set sessionWatermark($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasSessionWatermark() => $_has(2);
  @$pb.TagNumber(3)
  void clearSessionWatermark() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get createdAt => $_getSZ(3);
  @$pb.TagNumber(4)
  set createdAt($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasCreatedAt() => $_has(3);
  @$pb.TagNumber(4)
  void clearCreatedAt() => clearField(4);

  @$pb.TagNumber(5)
  $core.String get updatedAt => $_getSZ(4);
  @$pb.TagNumber(5)
  set updatedAt($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasUpdatedAt() => $_has(4);
  @$pb.TagNumber(5)
  void clearUpdatedAt() => clearField(5);

  @$pb.TagNumber(6)
  $core.String get disabledAt => $_getSZ(5);
  @$pb.TagNumber(6)
  set disabledAt($core.String v) { $_setString(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasDisabledAt() => $_has(5);
  @$pb.TagNumber(6)
  void clearDisabledAt() => clearField(6);

  @$pb.TagNumber(7)
  $core.String get disabledManuallyAt => $_getSZ(6);
  @$pb.TagNumber(7)
  set disabledManuallyAt($core.String v) { $_setString(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasDisabledManuallyAt() => $_has(6);
  @$pb.TagNumber(7)
  void clearDisabledManuallyAt() => clearField(7);

  @$pb.TagNumber(8)
  $core.String get disabledPendingApprovalAt => $_getSZ(7);
  @$pb.TagNumber(8)
  set disabledPendingApprovalAt($core.String v) { $_setString(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasDisabledPendingApprovalAt() => $_has(7);
  @$pb.TagNumber(8)
  void clearDisabledPendingApprovalAt() => clearField(8);

  @$pb.TagNumber(9)
  $core.String get display => $_getSZ(8);
  @$pb.TagNumber(9)
  set display($core.String v) { $_setString(8, v); }
  @$pb.TagNumber(9)
  $core.bool hasDisplay() => $_has(8);
  @$pb.TagNumber(9)
  void clearDisplay() => clearField(9);

  @$pb.TagNumber(10)
  $core.String get email => $_getSZ(9);
  @$pb.TagNumber(10)
  set email($core.String v) { $_setString(9, v); }
  @$pb.TagNumber(10)
  $core.bool hasEmail() => $_has(9);
  @$pb.TagNumber(10)
  void clearEmail() => clearField(10);
}

class ProfileSearchRequest extends $pb.GeneratedMessage {
  factory ProfileSearchRequest({
    $core.String? query,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
    $core.int? status,
  }) {
    final $result = create();
    if (query != null) {
      $result.query = query;
    }
    if (offset != null) {
      $result.offset = offset;
    }
    if (limit != null) {
      $result.limit = limit;
    }
    if (status != null) {
      $result.status = status;
    }
    return $result;
  }
  ProfileSearchRequest._() : super();
  factory ProfileSearchRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ProfileSearchRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProfileSearchRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$core.int>(4, _omitFieldNames ? '' : 'status', $pb.PbFieldType.OU3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ProfileSearchRequest clone() => ProfileSearchRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ProfileSearchRequest copyWith(void Function(ProfileSearchRequest) updates) => super.copyWith((message) => updates(message as ProfileSearchRequest)) as ProfileSearchRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileSearchRequest create() => ProfileSearchRequest._();
  ProfileSearchRequest createEmptyInstance() => create();
  static $pb.PbList<ProfileSearchRequest> createRepeated() => $pb.PbList<ProfileSearchRequest>();
  @$core.pragma('dart2js:noInline')
  static ProfileSearchRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileSearchRequest>(create);
  static ProfileSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get offset => $_getI64(1);
  @$pb.TagNumber(2)
  set offset($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasOffset() => $_has(1);
  @$pb.TagNumber(2)
  void clearOffset() => clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get limit => $_getI64(2);
  @$pb.TagNumber(3)
  set limit($fixnum.Int64 v) { $_setInt64(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(3)
  void clearLimit() => clearField(3);

  @$pb.TagNumber(4)
  $core.int get status => $_getIZ(3);
  @$pb.TagNumber(4)
  set status($core.int v) { $_setUnsignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasStatus() => $_has(3);
  @$pb.TagNumber(4)
  void clearStatus() => clearField(4);
}

class ProfileSearchResponse extends $pb.GeneratedMessage {
  factory ProfileSearchResponse({
    ProfileSearchRequest? next,
    $core.Iterable<Profile>? items,
  }) {
    final $result = create();
    if (next != null) {
      $result.next = next;
    }
    if (items != null) {
      $result.items.addAll(items);
    }
    return $result;
  }
  ProfileSearchResponse._() : super();
  factory ProfileSearchResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ProfileSearchResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProfileSearchResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<ProfileSearchRequest>(1, _omitFieldNames ? '' : 'next', subBuilder: ProfileSearchRequest.create)
    ..pc<Profile>(2, _omitFieldNames ? '' : 'items', $pb.PbFieldType.PM, subBuilder: Profile.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ProfileSearchResponse clone() => ProfileSearchResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ProfileSearchResponse copyWith(void Function(ProfileSearchResponse) updates) => super.copyWith((message) => updates(message as ProfileSearchResponse)) as ProfileSearchResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileSearchResponse create() => ProfileSearchResponse._();
  ProfileSearchResponse createEmptyInstance() => create();
  static $pb.PbList<ProfileSearchResponse> createRepeated() => $pb.PbList<ProfileSearchResponse>();
  @$core.pragma('dart2js:noInline')
  static ProfileSearchResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileSearchResponse>(create);
  static ProfileSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ProfileSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(ProfileSearchRequest v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => clearField(1);
  @$pb.TagNumber(1)
  ProfileSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.List<Profile> get items => $_getList(1);
}

class ProfileCreateRequest extends $pb.GeneratedMessage {
  factory ProfileCreateRequest({
    Profile? profile,
  }) {
    final $result = create();
    if (profile != null) {
      $result.profile = profile;
    }
    return $result;
  }
  ProfileCreateRequest._() : super();
  factory ProfileCreateRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ProfileCreateRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProfileCreateRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile', subBuilder: Profile.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ProfileCreateRequest clone() => ProfileCreateRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ProfileCreateRequest copyWith(void Function(ProfileCreateRequest) updates) => super.copyWith((message) => updates(message as ProfileCreateRequest)) as ProfileCreateRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileCreateRequest create() => ProfileCreateRequest._();
  ProfileCreateRequest createEmptyInstance() => create();
  static $pb.PbList<ProfileCreateRequest> createRepeated() => $pb.PbList<ProfileCreateRequest>();
  @$core.pragma('dart2js:noInline')
  static ProfileCreateRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileCreateRequest>(create);
  static ProfileCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);
}

class ProfileCreateResponse extends $pb.GeneratedMessage {
  factory ProfileCreateResponse({
    Profile? profile,
  }) {
    final $result = create();
    if (profile != null) {
      $result.profile = profile;
    }
    return $result;
  }
  ProfileCreateResponse._() : super();
  factory ProfileCreateResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ProfileCreateResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProfileCreateResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile', subBuilder: Profile.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ProfileCreateResponse clone() => ProfileCreateResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ProfileCreateResponse copyWith(void Function(ProfileCreateResponse) updates) => super.copyWith((message) => updates(message as ProfileCreateResponse)) as ProfileCreateResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileCreateResponse create() => ProfileCreateResponse._();
  ProfileCreateResponse createEmptyInstance() => create();
  static $pb.PbList<ProfileCreateResponse> createRepeated() => $pb.PbList<ProfileCreateResponse>();
  @$core.pragma('dart2js:noInline')
  static ProfileCreateResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileCreateResponse>(create);
  static ProfileCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);
}

class ProfileLookupRequest extends $pb.GeneratedMessage {
  factory ProfileLookupRequest() => create();
  ProfileLookupRequest._() : super();
  factory ProfileLookupRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ProfileLookupRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProfileLookupRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ProfileLookupRequest clone() => ProfileLookupRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ProfileLookupRequest copyWith(void Function(ProfileLookupRequest) updates) => super.copyWith((message) => updates(message as ProfileLookupRequest)) as ProfileLookupRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileLookupRequest create() => ProfileLookupRequest._();
  ProfileLookupRequest createEmptyInstance() => create();
  static $pb.PbList<ProfileLookupRequest> createRepeated() => $pb.PbList<ProfileLookupRequest>();
  @$core.pragma('dart2js:noInline')
  static ProfileLookupRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileLookupRequest>(create);
  static ProfileLookupRequest? _defaultInstance;
}

class ProfileLookupResponse extends $pb.GeneratedMessage {
  factory ProfileLookupResponse({
    Profile? profile,
  }) {
    final $result = create();
    if (profile != null) {
      $result.profile = profile;
    }
    return $result;
  }
  ProfileLookupResponse._() : super();
  factory ProfileLookupResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ProfileLookupResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProfileLookupResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile', subBuilder: Profile.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ProfileLookupResponse clone() => ProfileLookupResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ProfileLookupResponse copyWith(void Function(ProfileLookupResponse) updates) => super.copyWith((message) => updates(message as ProfileLookupResponse)) as ProfileLookupResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileLookupResponse create() => ProfileLookupResponse._();
  ProfileLookupResponse createEmptyInstance() => create();
  static $pb.PbList<ProfileLookupResponse> createRepeated() => $pb.PbList<ProfileLookupResponse>();
  @$core.pragma('dart2js:noInline')
  static ProfileLookupResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileLookupResponse>(create);
  static ProfileLookupResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);
}

class ProfileUpdateRequest extends $pb.GeneratedMessage {
  factory ProfileUpdateRequest({
    Profile? profile,
  }) {
    final $result = create();
    if (profile != null) {
      $result.profile = profile;
    }
    return $result;
  }
  ProfileUpdateRequest._() : super();
  factory ProfileUpdateRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ProfileUpdateRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProfileUpdateRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile', subBuilder: Profile.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ProfileUpdateRequest clone() => ProfileUpdateRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ProfileUpdateRequest copyWith(void Function(ProfileUpdateRequest) updates) => super.copyWith((message) => updates(message as ProfileUpdateRequest)) as ProfileUpdateRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileUpdateRequest create() => ProfileUpdateRequest._();
  ProfileUpdateRequest createEmptyInstance() => create();
  static $pb.PbList<ProfileUpdateRequest> createRepeated() => $pb.PbList<ProfileUpdateRequest>();
  @$core.pragma('dart2js:noInline')
  static ProfileUpdateRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileUpdateRequest>(create);
  static ProfileUpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);
}

class ProfileUpdateResponse extends $pb.GeneratedMessage {
  factory ProfileUpdateResponse({
    Profile? profile,
  }) {
    final $result = create();
    if (profile != null) {
      $result.profile = profile;
    }
    return $result;
  }
  ProfileUpdateResponse._() : super();
  factory ProfileUpdateResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ProfileUpdateResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProfileUpdateResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile', subBuilder: Profile.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ProfileUpdateResponse clone() => ProfileUpdateResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ProfileUpdateResponse copyWith(void Function(ProfileUpdateResponse) updates) => super.copyWith((message) => updates(message as ProfileUpdateResponse)) as ProfileUpdateResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileUpdateResponse create() => ProfileUpdateResponse._();
  ProfileUpdateResponse createEmptyInstance() => create();
  static $pb.PbList<ProfileUpdateResponse> createRepeated() => $pb.PbList<ProfileUpdateResponse>();
  @$core.pragma('dart2js:noInline')
  static ProfileUpdateResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileUpdateResponse>(create);
  static ProfileUpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);
}

class ProfileDisableRequest extends $pb.GeneratedMessage {
  factory ProfileDisableRequest({
    Profile? profile,
  }) {
    final $result = create();
    if (profile != null) {
      $result.profile = profile;
    }
    return $result;
  }
  ProfileDisableRequest._() : super();
  factory ProfileDisableRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ProfileDisableRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProfileDisableRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile', subBuilder: Profile.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ProfileDisableRequest clone() => ProfileDisableRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ProfileDisableRequest copyWith(void Function(ProfileDisableRequest) updates) => super.copyWith((message) => updates(message as ProfileDisableRequest)) as ProfileDisableRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileDisableRequest create() => ProfileDisableRequest._();
  ProfileDisableRequest createEmptyInstance() => create();
  static $pb.PbList<ProfileDisableRequest> createRepeated() => $pb.PbList<ProfileDisableRequest>();
  @$core.pragma('dart2js:noInline')
  static ProfileDisableRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileDisableRequest>(create);
  static ProfileDisableRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);
}

class ProfileDisableResponse extends $pb.GeneratedMessage {
  factory ProfileDisableResponse({
    Profile? profile,
  }) {
    final $result = create();
    if (profile != null) {
      $result.profile = profile;
    }
    return $result;
  }
  ProfileDisableResponse._() : super();
  factory ProfileDisableResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ProfileDisableResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProfileDisableResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Profile>(1, _omitFieldNames ? '' : 'profile', subBuilder: Profile.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ProfileDisableResponse clone() => ProfileDisableResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ProfileDisableResponse copyWith(void Function(ProfileDisableResponse) updates) => super.copyWith((message) => updates(message as ProfileDisableResponse)) as ProfileDisableResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileDisableResponse create() => ProfileDisableResponse._();
  ProfileDisableResponse createEmptyInstance() => create();
  static $pb.PbList<ProfileDisableResponse> createRepeated() => $pb.PbList<ProfileDisableResponse>();
  @$core.pragma('dart2js:noInline')
  static ProfileDisableResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileDisableResponse>(create);
  static ProfileDisableResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Profile get profile => $_getN(0);
  @$pb.TagNumber(1)
  set profile(Profile v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasProfile() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfile() => clearField(1);
  @$pb.TagNumber(1)
  Profile ensureProfile() => $_ensure(0);
}


const _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
