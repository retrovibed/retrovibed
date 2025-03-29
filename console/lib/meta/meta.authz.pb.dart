//
//  Generated code. Do not modify.
//  source: meta.authz.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

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
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (issuer != null) {
      $result.issuer = issuer;
    }
    if (profileId != null) {
      $result.profileId = profileId;
    }
    if (sessionId != null) {
      $result.sessionId = sessionId;
    }
    if (issued != null) {
      $result.issued = issued;
    }
    if (expires != null) {
      $result.expires = expires;
    }
    if (notBefore != null) {
      $result.notBefore = notBefore;
    }
    return $result;
  }
  Bearer._() : super();
  factory Bearer.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory Bearer.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'Bearer', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'jti', protoName: 'id')
    ..aOS(2, _omitFieldNames ? '' : 'iss', protoName: 'issuer')
    ..aOS(3, _omitFieldNames ? '' : 'sub', protoName: 'profile_id')
    ..aOS(4, _omitFieldNames ? '' : 'sid', protoName: 'session_id')
    ..aInt64(5, _omitFieldNames ? '' : 'iat', protoName: 'issued')
    ..aInt64(6, _omitFieldNames ? '' : 'exp', protoName: 'expires')
    ..aInt64(7, _omitFieldNames ? '' : 'nbf', protoName: 'not_before')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  Bearer clone() => Bearer()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  Bearer copyWith(void Function(Bearer) updates) => super.copyWith((message) => updates(message as Bearer)) as Bearer;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Bearer create() => Bearer._();
  Bearer createEmptyInstance() => create();
  static $pb.PbList<Bearer> createRepeated() => $pb.PbList<Bearer>();
  @$core.pragma('dart2js:noInline')
  static Bearer getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Bearer>(create);
  static Bearer? _defaultInstance;

  /// START OF STANDARD FIELDS
  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get issuer => $_getSZ(1);
  @$pb.TagNumber(2)
  set issuer($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasIssuer() => $_has(1);
  @$pb.TagNumber(2)
  void clearIssuer() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get profileId => $_getSZ(2);
  @$pb.TagNumber(3)
  set profileId($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasProfileId() => $_has(2);
  @$pb.TagNumber(3)
  void clearProfileId() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get sessionId => $_getSZ(3);
  @$pb.TagNumber(4)
  set sessionId($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasSessionId() => $_has(3);
  @$pb.TagNumber(4)
  void clearSessionId() => clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get issued => $_getI64(4);
  @$pb.TagNumber(5)
  set issued($fixnum.Int64 v) { $_setInt64(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasIssued() => $_has(4);
  @$pb.TagNumber(5)
  void clearIssued() => clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get expires => $_getI64(5);
  @$pb.TagNumber(6)
  set expires($fixnum.Int64 v) { $_setInt64(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasExpires() => $_has(5);
  @$pb.TagNumber(6)
  void clearExpires() => clearField(6);

  @$pb.TagNumber(7)
  $fixnum.Int64 get notBefore => $_getI64(6);
  @$pb.TagNumber(7)
  set notBefore($fixnum.Int64 v) { $_setInt64(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasNotBefore() => $_has(6);
  @$pb.TagNumber(7)
  void clearNotBefore() => clearField(7);
}

class Token extends $pb.GeneratedMessage {
  factory Token({
    $core.String? id,
    $core.String? accountId,
    $core.String? profileId,
    $core.String? sessionId,
    $fixnum.Int64? issued,
    $fixnum.Int64? expires,
    $fixnum.Int64? notBefore,
    $core.bool? usermanagement,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (accountId != null) {
      $result.accountId = accountId;
    }
    if (profileId != null) {
      $result.profileId = profileId;
    }
    if (sessionId != null) {
      $result.sessionId = sessionId;
    }
    if (issued != null) {
      $result.issued = issued;
    }
    if (expires != null) {
      $result.expires = expires;
    }
    if (notBefore != null) {
      $result.notBefore = notBefore;
    }
    if (usermanagement != null) {
      $result.usermanagement = usermanagement;
    }
    return $result;
  }
  Token._() : super();
  factory Token.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory Token.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'Token', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'jti', protoName: 'id')
    ..aOS(2, _omitFieldNames ? '' : 'iss', protoName: 'account_id')
    ..aOS(3, _omitFieldNames ? '' : 'sub', protoName: 'profile_id')
    ..aOS(4, _omitFieldNames ? '' : 'sid', protoName: 'session_id')
    ..aInt64(5, _omitFieldNames ? '' : 'iat', protoName: 'issued')
    ..aInt64(6, _omitFieldNames ? '' : 'exp', protoName: 'expires')
    ..aInt64(7, _omitFieldNames ? '' : 'nbf', protoName: 'not_before')
    ..aOB(1000, _omitFieldNames ? '' : 'usermanagement')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  Token clone() => Token()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  Token copyWith(void Function(Token) updates) => super.copyWith((message) => updates(message as Token)) as Token;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Token create() => Token._();
  Token createEmptyInstance() => create();
  static $pb.PbList<Token> createRepeated() => $pb.PbList<Token>();
  @$core.pragma('dart2js:noInline')
  static Token getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Token>(create);
  static Token? _defaultInstance;

  /// START OF STANDARD FIELDS
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
  $core.String get profileId => $_getSZ(2);
  @$pb.TagNumber(3)
  set profileId($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasProfileId() => $_has(2);
  @$pb.TagNumber(3)
  void clearProfileId() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get sessionId => $_getSZ(3);
  @$pb.TagNumber(4)
  set sessionId($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasSessionId() => $_has(3);
  @$pb.TagNumber(4)
  void clearSessionId() => clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get issued => $_getI64(4);
  @$pb.TagNumber(5)
  set issued($fixnum.Int64 v) { $_setInt64(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasIssued() => $_has(4);
  @$pb.TagNumber(5)
  void clearIssued() => clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get expires => $_getI64(5);
  @$pb.TagNumber(6)
  set expires($fixnum.Int64 v) { $_setInt64(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasExpires() => $_has(5);
  @$pb.TagNumber(6)
  void clearExpires() => clearField(6);

  @$pb.TagNumber(7)
  $fixnum.Int64 get notBefore => $_getI64(6);
  @$pb.TagNumber(7)
  set notBefore($fixnum.Int64 v) { $_setInt64(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasNotBefore() => $_has(6);
  @$pb.TagNumber(7)
  void clearNotBefore() => clearField(7);

  @$pb.TagNumber(1000)
  $core.bool get usermanagement => $_getBF(7);
  @$pb.TagNumber(1000)
  set usermanagement($core.bool v) { $_setBool(7, v); }
  @$pb.TagNumber(1000)
  $core.bool hasUsermanagement() => $_has(7);
  @$pb.TagNumber(1000)
  void clearUsermanagement() => clearField(1000);
}

class AuthzRequest extends $pb.GeneratedMessage {
  factory AuthzRequest() => create();
  AuthzRequest._() : super();
  factory AuthzRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AuthzRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AuthzRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AuthzRequest clone() => AuthzRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AuthzRequest copyWith(void Function(AuthzRequest) updates) => super.copyWith((message) => updates(message as AuthzRequest)) as AuthzRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzRequest create() => AuthzRequest._();
  AuthzRequest createEmptyInstance() => create();
  static $pb.PbList<AuthzRequest> createRepeated() => $pb.PbList<AuthzRequest>();
  @$core.pragma('dart2js:noInline')
  static AuthzRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AuthzRequest>(create);
  static AuthzRequest? _defaultInstance;
}

class AuthzResponse extends $pb.GeneratedMessage {
  factory AuthzResponse({
    $core.String? bearer,
    Token? token,
  }) {
    final $result = create();
    if (bearer != null) {
      $result.bearer = bearer;
    }
    if (token != null) {
      $result.token = token;
    }
    return $result;
  }
  AuthzResponse._() : super();
  factory AuthzResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AuthzResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AuthzResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'bearer')
    ..aOM<Token>(2, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AuthzResponse clone() => AuthzResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AuthzResponse copyWith(void Function(AuthzResponse) updates) => super.copyWith((message) => updates(message as AuthzResponse)) as AuthzResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzResponse create() => AuthzResponse._();
  AuthzResponse createEmptyInstance() => create();
  static $pb.PbList<AuthzResponse> createRepeated() => $pb.PbList<AuthzResponse>();
  @$core.pragma('dart2js:noInline')
  static AuthzResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AuthzResponse>(create);
  static AuthzResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get bearer => $_getSZ(0);
  @$pb.TagNumber(1)
  set bearer($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasBearer() => $_has(0);
  @$pb.TagNumber(1)
  void clearBearer() => clearField(1);

  @$pb.TagNumber(2)
  Token get token => $_getN(1);
  @$pb.TagNumber(2)
  set token(Token v) { setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasToken() => $_has(1);
  @$pb.TagNumber(2)
  void clearToken() => clearField(2);
  @$pb.TagNumber(2)
  Token ensureToken() => $_ensure(1);
}

class GrantRequest extends $pb.GeneratedMessage {
  factory GrantRequest({
    Token? token,
  }) {
    final $result = create();
    if (token != null) {
      $result.token = token;
    }
    return $result;
  }
  GrantRequest._() : super();
  factory GrantRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GrantRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GrantRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GrantRequest clone() => GrantRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GrantRequest copyWith(void Function(GrantRequest) updates) => super.copyWith((message) => updates(message as GrantRequest)) as GrantRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GrantRequest create() => GrantRequest._();
  GrantRequest createEmptyInstance() => create();
  static $pb.PbList<GrantRequest> createRepeated() => $pb.PbList<GrantRequest>();
  @$core.pragma('dart2js:noInline')
  static GrantRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GrantRequest>(create);
  static GrantRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Token get token => $_getN(0);
  @$pb.TagNumber(1)
  set token(Token v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => clearField(1);
  @$pb.TagNumber(1)
  Token ensureToken() => $_ensure(0);
}

class GrantResponse extends $pb.GeneratedMessage {
  factory GrantResponse({
    Token? token,
  }) {
    final $result = create();
    if (token != null) {
      $result.token = token;
    }
    return $result;
  }
  GrantResponse._() : super();
  factory GrantResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GrantResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GrantResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GrantResponse clone() => GrantResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GrantResponse copyWith(void Function(GrantResponse) updates) => super.copyWith((message) => updates(message as GrantResponse)) as GrantResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GrantResponse create() => GrantResponse._();
  GrantResponse createEmptyInstance() => create();
  static $pb.PbList<GrantResponse> createRepeated() => $pb.PbList<GrantResponse>();
  @$core.pragma('dart2js:noInline')
  static GrantResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GrantResponse>(create);
  static GrantResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Token get token => $_getN(0);
  @$pb.TagNumber(1)
  set token(Token v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => clearField(1);
  @$pb.TagNumber(1)
  Token ensureToken() => $_ensure(0);
}

class RevokeRequest extends $pb.GeneratedMessage {
  factory RevokeRequest({
    Token? token,
  }) {
    final $result = create();
    if (token != null) {
      $result.token = token;
    }
    return $result;
  }
  RevokeRequest._() : super();
  factory RevokeRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory RevokeRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'RevokeRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  RevokeRequest clone() => RevokeRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  RevokeRequest copyWith(void Function(RevokeRequest) updates) => super.copyWith((message) => updates(message as RevokeRequest)) as RevokeRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RevokeRequest create() => RevokeRequest._();
  RevokeRequest createEmptyInstance() => create();
  static $pb.PbList<RevokeRequest> createRepeated() => $pb.PbList<RevokeRequest>();
  @$core.pragma('dart2js:noInline')
  static RevokeRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RevokeRequest>(create);
  static RevokeRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Token get token => $_getN(0);
  @$pb.TagNumber(1)
  set token(Token v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => clearField(1);
  @$pb.TagNumber(1)
  Token ensureToken() => $_ensure(0);
}

class RevokeResponse extends $pb.GeneratedMessage {
  factory RevokeResponse({
    Token? token,
  }) {
    final $result = create();
    if (token != null) {
      $result.token = token;
    }
    return $result;
  }
  RevokeResponse._() : super();
  factory RevokeResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory RevokeResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'RevokeResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  RevokeResponse clone() => RevokeResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  RevokeResponse copyWith(void Function(RevokeResponse) updates) => super.copyWith((message) => updates(message as RevokeResponse)) as RevokeResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RevokeResponse create() => RevokeResponse._();
  RevokeResponse createEmptyInstance() => create();
  static $pb.PbList<RevokeResponse> createRepeated() => $pb.PbList<RevokeResponse>();
  @$core.pragma('dart2js:noInline')
  static RevokeResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RevokeResponse>(create);
  static RevokeResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Token get token => $_getN(0);
  @$pb.TagNumber(1)
  set token(Token v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => clearField(1);
  @$pb.TagNumber(1)
  Token ensureToken() => $_ensure(0);
}

class ProfileRequest extends $pb.GeneratedMessage {
  factory ProfileRequest({
    $core.String? profileId,
  }) {
    final $result = create();
    if (profileId != null) {
      $result.profileId = profileId;
    }
    return $result;
  }
  ProfileRequest._() : super();
  factory ProfileRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ProfileRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProfileRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'profile_id')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ProfileRequest clone() => ProfileRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ProfileRequest copyWith(void Function(ProfileRequest) updates) => super.copyWith((message) => updates(message as ProfileRequest)) as ProfileRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileRequest create() => ProfileRequest._();
  ProfileRequest createEmptyInstance() => create();
  static $pb.PbList<ProfileRequest> createRepeated() => $pb.PbList<ProfileRequest>();
  @$core.pragma('dart2js:noInline')
  static ProfileRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileRequest>(create);
  static ProfileRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get profileId => $_getSZ(0);
  @$pb.TagNumber(1)
  set profileId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasProfileId() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfileId() => clearField(1);
}

class ProfileResponse extends $pb.GeneratedMessage {
  factory ProfileResponse({
    Token? token,
  }) {
    final $result = create();
    if (token != null) {
      $result.token = token;
    }
    return $result;
  }
  ProfileResponse._() : super();
  factory ProfileResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ProfileResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProfileResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ProfileResponse clone() => ProfileResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ProfileResponse copyWith(void Function(ProfileResponse) updates) => super.copyWith((message) => updates(message as ProfileResponse)) as ProfileResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProfileResponse create() => ProfileResponse._();
  ProfileResponse createEmptyInstance() => create();
  static $pb.PbList<ProfileResponse> createRepeated() => $pb.PbList<ProfileResponse>();
  @$core.pragma('dart2js:noInline')
  static ProfileResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileResponse>(create);
  static ProfileResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Token get token => $_getN(0);
  @$pb.TagNumber(1)
  set token(Token v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => clearField(1);
  @$pb.TagNumber(1)
  Token ensureToken() => $_ensure(0);
}


const _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
