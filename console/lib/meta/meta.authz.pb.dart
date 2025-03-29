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
    $core.String? issuer,
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
    ..aOS(2, _omitFieldNames ? '' : 'iss', protoName: 'issuer')
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

class AuthzGrantRequest extends $pb.GeneratedMessage {
  factory AuthzGrantRequest({
    Token? token,
  }) {
    final $result = create();
    if (token != null) {
      $result.token = token;
    }
    return $result;
  }
  AuthzGrantRequest._() : super();
  factory AuthzGrantRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AuthzGrantRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AuthzGrantRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AuthzGrantRequest clone() => AuthzGrantRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AuthzGrantRequest copyWith(void Function(AuthzGrantRequest) updates) => super.copyWith((message) => updates(message as AuthzGrantRequest)) as AuthzGrantRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzGrantRequest create() => AuthzGrantRequest._();
  AuthzGrantRequest createEmptyInstance() => create();
  static $pb.PbList<AuthzGrantRequest> createRepeated() => $pb.PbList<AuthzGrantRequest>();
  @$core.pragma('dart2js:noInline')
  static AuthzGrantRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AuthzGrantRequest>(create);
  static AuthzGrantRequest? _defaultInstance;

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

class AuthzGrantResponse extends $pb.GeneratedMessage {
  factory AuthzGrantResponse({
    Token? token,
  }) {
    final $result = create();
    if (token != null) {
      $result.token = token;
    }
    return $result;
  }
  AuthzGrantResponse._() : super();
  factory AuthzGrantResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AuthzGrantResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AuthzGrantResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AuthzGrantResponse clone() => AuthzGrantResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AuthzGrantResponse copyWith(void Function(AuthzGrantResponse) updates) => super.copyWith((message) => updates(message as AuthzGrantResponse)) as AuthzGrantResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzGrantResponse create() => AuthzGrantResponse._();
  AuthzGrantResponse createEmptyInstance() => create();
  static $pb.PbList<AuthzGrantResponse> createRepeated() => $pb.PbList<AuthzGrantResponse>();
  @$core.pragma('dart2js:noInline')
  static AuthzGrantResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AuthzGrantResponse>(create);
  static AuthzGrantResponse? _defaultInstance;

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

class AuthzRevokeRequest extends $pb.GeneratedMessage {
  factory AuthzRevokeRequest({
    Token? token,
  }) {
    final $result = create();
    if (token != null) {
      $result.token = token;
    }
    return $result;
  }
  AuthzRevokeRequest._() : super();
  factory AuthzRevokeRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AuthzRevokeRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AuthzRevokeRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AuthzRevokeRequest clone() => AuthzRevokeRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AuthzRevokeRequest copyWith(void Function(AuthzRevokeRequest) updates) => super.copyWith((message) => updates(message as AuthzRevokeRequest)) as AuthzRevokeRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzRevokeRequest create() => AuthzRevokeRequest._();
  AuthzRevokeRequest createEmptyInstance() => create();
  static $pb.PbList<AuthzRevokeRequest> createRepeated() => $pb.PbList<AuthzRevokeRequest>();
  @$core.pragma('dart2js:noInline')
  static AuthzRevokeRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AuthzRevokeRequest>(create);
  static AuthzRevokeRequest? _defaultInstance;

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

class AuthzRevokeResponse extends $pb.GeneratedMessage {
  factory AuthzRevokeResponse({
    Token? token,
  }) {
    final $result = create();
    if (token != null) {
      $result.token = token;
    }
    return $result;
  }
  AuthzRevokeResponse._() : super();
  factory AuthzRevokeResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AuthzRevokeResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AuthzRevokeResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AuthzRevokeResponse clone() => AuthzRevokeResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AuthzRevokeResponse copyWith(void Function(AuthzRevokeResponse) updates) => super.copyWith((message) => updates(message as AuthzRevokeResponse)) as AuthzRevokeResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzRevokeResponse create() => AuthzRevokeResponse._();
  AuthzRevokeResponse createEmptyInstance() => create();
  static $pb.PbList<AuthzRevokeResponse> createRepeated() => $pb.PbList<AuthzRevokeResponse>();
  @$core.pragma('dart2js:noInline')
  static AuthzRevokeResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AuthzRevokeResponse>(create);
  static AuthzRevokeResponse? _defaultInstance;

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

class AuthzProfileRequest extends $pb.GeneratedMessage {
  factory AuthzProfileRequest({
    $core.String? profileId,
  }) {
    final $result = create();
    if (profileId != null) {
      $result.profileId = profileId;
    }
    return $result;
  }
  AuthzProfileRequest._() : super();
  factory AuthzProfileRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AuthzProfileRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AuthzProfileRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'profile_id')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AuthzProfileRequest clone() => AuthzProfileRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AuthzProfileRequest copyWith(void Function(AuthzProfileRequest) updates) => super.copyWith((message) => updates(message as AuthzProfileRequest)) as AuthzProfileRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzProfileRequest create() => AuthzProfileRequest._();
  AuthzProfileRequest createEmptyInstance() => create();
  static $pb.PbList<AuthzProfileRequest> createRepeated() => $pb.PbList<AuthzProfileRequest>();
  @$core.pragma('dart2js:noInline')
  static AuthzProfileRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AuthzProfileRequest>(create);
  static AuthzProfileRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get profileId => $_getSZ(0);
  @$pb.TagNumber(1)
  set profileId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasProfileId() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfileId() => clearField(1);
}

class AuthzProfileResponse extends $pb.GeneratedMessage {
  factory AuthzProfileResponse({
    Token? token,
  }) {
    final $result = create();
    if (token != null) {
      $result.token = token;
    }
    return $result;
  }
  AuthzProfileResponse._() : super();
  factory AuthzProfileResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AuthzProfileResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AuthzProfileResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Token>(1, _omitFieldNames ? '' : 'token', subBuilder: Token.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AuthzProfileResponse clone() => AuthzProfileResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AuthzProfileResponse copyWith(void Function(AuthzProfileResponse) updates) => super.copyWith((message) => updates(message as AuthzProfileResponse)) as AuthzProfileResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthzProfileResponse create() => AuthzProfileResponse._();
  AuthzProfileResponse createEmptyInstance() => create();
  static $pb.PbList<AuthzProfileResponse> createRepeated() => $pb.PbList<AuthzProfileResponse>();
  @$core.pragma('dart2js:noInline')
  static AuthzProfileResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AuthzProfileResponse>(create);
  static AuthzProfileResponse? _defaultInstance;

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
