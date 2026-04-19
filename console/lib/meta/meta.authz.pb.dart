// This is a generated file - do not edit.
//
// Generated from meta.authz.proto.

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

class Bearer extends $pb.GeneratedMessage {
  factory Bearer({
    $core.String? id,
    $core.String? issuer,
    $core.String? profileId,
    $core.String? sessionId,
    $fixnum.Int64? issued,
    $fixnum.Int64? expires,
    $fixnum.Int64? notBefore,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (issuer != null) result.issuer = issuer;
    if (profileId != null) result.profileId = profileId;
    if (sessionId != null) result.sessionId = sessionId;
    if (issued != null) result.issued = issued;
    if (expires != null) result.expires = expires;
    if (notBefore != null) result.notBefore = notBefore;
    return result;
  }

  Bearer._();

  factory Bearer.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Bearer.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Bearer',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'jti', protoName: 'id')
    ..aOS(2, _omitFieldNames ? '' : 'iss', protoName: 'issuer')
    ..aOS(3, _omitFieldNames ? '' : 'sub', protoName: 'profile_id')
    ..aOS(4, _omitFieldNames ? '' : 'sid', protoName: 'session_id')
    ..aInt64(5, _omitFieldNames ? '' : 'iat', protoName: 'issued')
    ..aInt64(6, _omitFieldNames ? '' : 'exp', protoName: 'expires')
    ..aInt64(7, _omitFieldNames ? '' : 'nbf', protoName: 'not_before')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Bearer clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Bearer copyWith(void Function(Bearer) updates) =>
      super.copyWith((message) => updates(message as Bearer)) as Bearer;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Bearer create() => Bearer._();
  @$core.override
  Bearer createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Bearer getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Bearer>(create);
  static Bearer? _defaultInstance;

  /// START OF STANDARD FIELDS
  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get issuer => $_getSZ(1);
  @$pb.TagNumber(2)
  set issuer($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasIssuer() => $_has(1);
  @$pb.TagNumber(2)
  void clearIssuer() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get profileId => $_getSZ(2);
  @$pb.TagNumber(3)
  set profileId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasProfileId() => $_has(2);
  @$pb.TagNumber(3)
  void clearProfileId() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get sessionId => $_getSZ(3);
  @$pb.TagNumber(4)
  set sessionId($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasSessionId() => $_has(3);
  @$pb.TagNumber(4)
  void clearSessionId() => $_clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get issued => $_getI64(4);
  @$pb.TagNumber(5)
  set issued($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasIssued() => $_has(4);
  @$pb.TagNumber(5)
  void clearIssued() => $_clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get expires => $_getI64(5);
  @$pb.TagNumber(6)
  set expires($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(6)
  $core.bool hasExpires() => $_has(5);
  @$pb.TagNumber(6)
  void clearExpires() => $_clearField(6);

  @$pb.TagNumber(7)
  $fixnum.Int64 get notBefore => $_getI64(6);
  @$pb.TagNumber(7)
  set notBefore($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(7)
  $core.bool hasNotBefore() => $_has(6);
  @$pb.TagNumber(7)
  void clearNotBefore() => $_clearField(7);
}

class Token extends $pb.GeneratedMessage {
  factory Token({
    $core.String? id,
    $core.String? issuer,
    $core.String? profileId,
    $core.String? sessionId,
    $fixnum.Int64? issued,
    $fixnum.Int64? expires,
    $fixnum.Int64? notBefore,
    $core.bool? usermanagement,
    $core.bool? billingRead,
    $core.bool? billingModify,
    $core.bool? communityModify,
    $fixnum.Int64? archiveUpload,
    $fixnum.Int64? archiveDownload,
    $core.bool? libraryRead,
    $core.bool? libraryModify,
    $core.bool? archiveSync,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (issuer != null) result.issuer = issuer;
    if (profileId != null) result.profileId = profileId;
    if (sessionId != null) result.sessionId = sessionId;
    if (issued != null) result.issued = issued;
    if (expires != null) result.expires = expires;
    if (notBefore != null) result.notBefore = notBefore;
    if (usermanagement != null) result.usermanagement = usermanagement;
    if (billingRead != null) result.billingRead = billingRead;
    if (billingModify != null) result.billingModify = billingModify;
    if (communityModify != null) result.communityModify = communityModify;
    if (archiveUpload != null) result.archiveUpload = archiveUpload;
    if (archiveDownload != null) result.archiveDownload = archiveDownload;
    if (libraryRead != null) result.libraryRead = libraryRead;
    if (libraryModify != null) result.libraryModify = libraryModify;
    if (archiveSync != null) result.archiveSync = archiveSync;
    return result;
  }

  Token._();

  factory Token.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Token.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Token',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'jti', protoName: 'id')
    ..aOS(2, _omitFieldNames ? '' : 'iss', protoName: 'issuer')
    ..aOS(3, _omitFieldNames ? '' : 'sub', protoName: 'profile_id')
    ..aOS(4, _omitFieldNames ? '' : 'sid', protoName: 'session_id')
    ..aInt64(5, _omitFieldNames ? '' : 'iat', protoName: 'issued')
    ..aInt64(6, _omitFieldNames ? '' : 'exp', protoName: 'expires')
    ..aInt64(7, _omitFieldNames ? '' : 'nbf', protoName: 'not_before')
    ..aOB(1000, _omitFieldNames ? '' : 'usermanagement')
    ..aOB(1002, _omitFieldNames ? '' : 'billing_read')
    ..aOB(1003, _omitFieldNames ? '' : 'billing_modify')
    ..aOB(1004, _omitFieldNames ? '' : 'community_modify')
    ..a<$fixnum.Int64>(
        1005, _omitFieldNames ? '' : 'archive_upload', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        1006, _omitFieldNames ? '' : 'archive_download', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOB(1007, _omitFieldNames ? '' : 'library_read')
    ..aOB(1008, _omitFieldNames ? '' : 'library_modify')
    ..aOB(1009, _omitFieldNames ? '' : 'archive_sync')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Token clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Token copyWith(void Function(Token) updates) =>
      super.copyWith((message) => updates(message as Token)) as Token;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Token create() => Token._();
  @$core.override
  Token createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Token getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Token>(create);
  static Token? _defaultInstance;

  /// START OF STANDARD FIELDS
  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get issuer => $_getSZ(1);
  @$pb.TagNumber(2)
  set issuer($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasIssuer() => $_has(1);
  @$pb.TagNumber(2)
  void clearIssuer() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get profileId => $_getSZ(2);
  @$pb.TagNumber(3)
  set profileId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasProfileId() => $_has(2);
  @$pb.TagNumber(3)
  void clearProfileId() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get sessionId => $_getSZ(3);
  @$pb.TagNumber(4)
  set sessionId($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasSessionId() => $_has(3);
  @$pb.TagNumber(4)
  void clearSessionId() => $_clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get issued => $_getI64(4);
  @$pb.TagNumber(5)
  set issued($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasIssued() => $_has(4);
  @$pb.TagNumber(5)
  void clearIssued() => $_clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get expires => $_getI64(5);
  @$pb.TagNumber(6)
  set expires($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(6)
  $core.bool hasExpires() => $_has(5);
  @$pb.TagNumber(6)
  void clearExpires() => $_clearField(6);

  @$pb.TagNumber(7)
  $fixnum.Int64 get notBefore => $_getI64(6);
  @$pb.TagNumber(7)
  set notBefore($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(7)
  $core.bool hasNotBefore() => $_has(6);
  @$pb.TagNumber(7)
  void clearNotBefore() => $_clearField(7);

  @$pb.TagNumber(1000)
  $core.bool get usermanagement => $_getBF(7);
  @$pb.TagNumber(1000)
  set usermanagement($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(1000)
  $core.bool hasUsermanagement() => $_has(7);
  @$pb.TagNumber(1000)
  void clearUsermanagement() => $_clearField(1000);

  @$pb.TagNumber(1002)
  $core.bool get billingRead => $_getBF(8);
  @$pb.TagNumber(1002)
  set billingRead($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(1002)
  $core.bool hasBillingRead() => $_has(8);
  @$pb.TagNumber(1002)
  void clearBillingRead() => $_clearField(1002);

  @$pb.TagNumber(1003)
  $core.bool get billingModify => $_getBF(9);
  @$pb.TagNumber(1003)
  set billingModify($core.bool value) => $_setBool(9, value);
  @$pb.TagNumber(1003)
  $core.bool hasBillingModify() => $_has(9);
  @$pb.TagNumber(1003)
  void clearBillingModify() => $_clearField(1003);

  @$pb.TagNumber(1004)
  $core.bool get communityModify => $_getBF(10);
  @$pb.TagNumber(1004)
  set communityModify($core.bool value) => $_setBool(10, value);
  @$pb.TagNumber(1004)
  $core.bool hasCommunityModify() => $_has(10);
  @$pb.TagNumber(1004)
  void clearCommunityModify() => $_clearField(1004);

  @$pb.TagNumber(1005)
  $fixnum.Int64 get archiveUpload => $_getI64(11);
  @$pb.TagNumber(1005)
  set archiveUpload($fixnum.Int64 value) => $_setInt64(11, value);
  @$pb.TagNumber(1005)
  $core.bool hasArchiveUpload() => $_has(11);
  @$pb.TagNumber(1005)
  void clearArchiveUpload() => $_clearField(1005);

  @$pb.TagNumber(1006)
  $fixnum.Int64 get archiveDownload => $_getI64(12);
  @$pb.TagNumber(1006)
  set archiveDownload($fixnum.Int64 value) => $_setInt64(12, value);
  @$pb.TagNumber(1006)
  $core.bool hasArchiveDownload() => $_has(12);
  @$pb.TagNumber(1006)
  void clearArchiveDownload() => $_clearField(1006);

  @$pb.TagNumber(1007)
  $core.bool get libraryRead => $_getBF(13);
  @$pb.TagNumber(1007)
  set libraryRead($core.bool value) => $_setBool(13, value);
  @$pb.TagNumber(1007)
  $core.bool hasLibraryRead() => $_has(13);
  @$pb.TagNumber(1007)
  void clearLibraryRead() => $_clearField(1007);

  @$pb.TagNumber(1008)
  $core.bool get libraryModify => $_getBF(14);
  @$pb.TagNumber(1008)
  set libraryModify($core.bool value) => $_setBool(14, value);
  @$pb.TagNumber(1008)
  $core.bool hasLibraryModify() => $_has(14);
  @$pb.TagNumber(1008)
  void clearLibraryModify() => $_clearField(1008);

  @$pb.TagNumber(1009)
  $core.bool get archiveSync => $_getBF(15);
  @$pb.TagNumber(1009)
  set archiveSync($core.bool value) => $_setBool(15, value);
  @$pb.TagNumber(1009)
  $core.bool hasArchiveSync() => $_has(15);
  @$pb.TagNumber(1009)
  void clearArchiveSync() => $_clearField(1009);
}

class AuthzRequest extends $pb.GeneratedMessage {
  factory AuthzRequest() => create();

  AuthzRequest._();

  factory AuthzRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AuthzRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AuthzRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzRequest copyWith(void Function(AuthzRequest) updates) =>
      super.copyWith((message) => updates(message as AuthzRequest))
          as AuthzRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzRequest create() => AuthzRequest._();
  @$core.override
  AuthzRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AuthzRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AuthzRequest>(create);
  static AuthzRequest? _defaultInstance;
}

class AuthzResponse extends $pb.GeneratedMessage {
  factory AuthzResponse({
    $core.String? bearer,
    Token? token,
  }) {
    final result = create();
    if (bearer != null) result.bearer = bearer;
    if (token != null) result.token = token;
    return result;
  }

  AuthzResponse._();

  factory AuthzResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AuthzResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AuthzResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'bearer')
    ..aOM<Token>(2, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzResponse copyWith(void Function(AuthzResponse) updates) =>
      super.copyWith((message) => updates(message as AuthzResponse))
          as AuthzResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzResponse create() => AuthzResponse._();
  @$core.override
  AuthzResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AuthzResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AuthzResponse>(create);
  static AuthzResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get bearer => $_getSZ(0);
  @$pb.TagNumber(1)
  set bearer($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasBearer() => $_has(0);
  @$pb.TagNumber(1)
  void clearBearer() => $_clearField(1);

  @$pb.TagNumber(2)
  Token get token => $_getN(1);
  @$pb.TagNumber(2)
  set token(Token value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasToken() => $_has(1);
  @$pb.TagNumber(2)
  void clearToken() => $_clearField(2);
  @$pb.TagNumber(2)
  Token ensureToken() => $_ensure(1);
}

class AuthzGrantRequest extends $pb.GeneratedMessage {
  factory AuthzGrantRequest({
    Token? token,
  }) {
    final result = create();
    if (token != null) result.token = token;
    return result;
  }

  AuthzGrantRequest._();

  factory AuthzGrantRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AuthzGrantRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AuthzGrantRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzGrantRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzGrantRequest copyWith(void Function(AuthzGrantRequest) updates) =>
      super.copyWith((message) => updates(message as AuthzGrantRequest))
          as AuthzGrantRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzGrantRequest create() => AuthzGrantRequest._();
  @$core.override
  AuthzGrantRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AuthzGrantRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AuthzGrantRequest>(create);
  static AuthzGrantRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Token get token => $_getN(0);
  @$pb.TagNumber(1)
  set token(Token value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => $_clearField(1);
  @$pb.TagNumber(1)
  Token ensureToken() => $_ensure(0);
}

class AuthzGrantResponse extends $pb.GeneratedMessage {
  factory AuthzGrantResponse({
    Token? token,
  }) {
    final result = create();
    if (token != null) result.token = token;
    return result;
  }

  AuthzGrantResponse._();

  factory AuthzGrantResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AuthzGrantResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AuthzGrantResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzGrantResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzGrantResponse copyWith(void Function(AuthzGrantResponse) updates) =>
      super.copyWith((message) => updates(message as AuthzGrantResponse))
          as AuthzGrantResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzGrantResponse create() => AuthzGrantResponse._();
  @$core.override
  AuthzGrantResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AuthzGrantResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AuthzGrantResponse>(create);
  static AuthzGrantResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Token get token => $_getN(0);
  @$pb.TagNumber(1)
  set token(Token value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => $_clearField(1);
  @$pb.TagNumber(1)
  Token ensureToken() => $_ensure(0);
}

class AuthzRevokeRequest extends $pb.GeneratedMessage {
  factory AuthzRevokeRequest({
    Token? token,
  }) {
    final result = create();
    if (token != null) result.token = token;
    return result;
  }

  AuthzRevokeRequest._();

  factory AuthzRevokeRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AuthzRevokeRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AuthzRevokeRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzRevokeRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzRevokeRequest copyWith(void Function(AuthzRevokeRequest) updates) =>
      super.copyWith((message) => updates(message as AuthzRevokeRequest))
          as AuthzRevokeRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzRevokeRequest create() => AuthzRevokeRequest._();
  @$core.override
  AuthzRevokeRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AuthzRevokeRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AuthzRevokeRequest>(create);
  static AuthzRevokeRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Token get token => $_getN(0);
  @$pb.TagNumber(1)
  set token(Token value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => $_clearField(1);
  @$pb.TagNumber(1)
  Token ensureToken() => $_ensure(0);
}

class AuthzRevokeResponse extends $pb.GeneratedMessage {
  factory AuthzRevokeResponse({
    Token? token,
  }) {
    final result = create();
    if (token != null) result.token = token;
    return result;
  }

  AuthzRevokeResponse._();

  factory AuthzRevokeResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AuthzRevokeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AuthzRevokeResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzRevokeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzRevokeResponse copyWith(void Function(AuthzRevokeResponse) updates) =>
      super.copyWith((message) => updates(message as AuthzRevokeResponse))
          as AuthzRevokeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzRevokeResponse create() => AuthzRevokeResponse._();
  @$core.override
  AuthzRevokeResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AuthzRevokeResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AuthzRevokeResponse>(create);
  static AuthzRevokeResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Token get token => $_getN(0);
  @$pb.TagNumber(1)
  set token(Token value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => $_clearField(1);
  @$pb.TagNumber(1)
  Token ensureToken() => $_ensure(0);
}

class AuthzProfileRequest extends $pb.GeneratedMessage {
  factory AuthzProfileRequest({
    $core.String? profileId,
  }) {
    final result = create();
    if (profileId != null) result.profileId = profileId;
    return result;
  }

  AuthzProfileRequest._();

  factory AuthzProfileRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AuthzProfileRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AuthzProfileRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'profile_id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzProfileRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzProfileRequest copyWith(void Function(AuthzProfileRequest) updates) =>
      super.copyWith((message) => updates(message as AuthzProfileRequest))
          as AuthzProfileRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzProfileRequest create() => AuthzProfileRequest._();
  @$core.override
  AuthzProfileRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AuthzProfileRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AuthzProfileRequest>(create);
  static AuthzProfileRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get profileId => $_getSZ(0);
  @$pb.TagNumber(1)
  set profileId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasProfileId() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfileId() => $_clearField(1);
}

class AuthzProfileResponse extends $pb.GeneratedMessage {
  factory AuthzProfileResponse({
    Token? token,
  }) {
    final result = create();
    if (token != null) result.token = token;
    return result;
  }

  AuthzProfileResponse._();

  factory AuthzProfileResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AuthzProfileResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AuthzProfileResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzProfileResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthzProfileResponse copyWith(void Function(AuthzProfileResponse) updates) =>
      super.copyWith((message) => updates(message as AuthzProfileResponse))
          as AuthzProfileResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzProfileResponse create() => AuthzProfileResponse._();
  @$core.override
  AuthzProfileResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AuthzProfileResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AuthzProfileResponse>(create);
  static AuthzProfileResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Token get token => $_getN(0);
  @$pb.TagNumber(1)
  set token(Token value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => $_clearField(1);
  @$pb.TagNumber(1)
  Token ensureToken() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
