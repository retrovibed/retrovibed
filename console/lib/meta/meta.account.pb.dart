// This is a generated file - do not edit.
//
// Generated from meta/meta.account.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class Account extends $pb.GeneratedMessage {
  factory Account({
    $core.String? id,
    $core.String? description,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? disabledAt,
    $core.String? email,
    $core.String? phone,
  }) {
    final result = Account._();
    if (id != null) result.id = id;
    if (description != null) result.description = description;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (disabledAt != null) result.disabledAt = disabledAt;
    if (email != null) result.email = email;
    if (phone != null) result.phone = phone;
    return result;
  }

  Account._();

  factory Account.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Account()..mergeFromBuffer(data, registry);
  factory Account.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Account()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Account',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: Account.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'description')
    ..aOS(3, _omitFieldNames ? '' : 'created_at')
    ..aOS(4, _omitFieldNames ? '' : 'updated_at')
    ..aOS(5, _omitFieldNames ? '' : 'disabled_at')
    ..aOS(6, _omitFieldNames ? '' : 'email')
    ..aOS(7, _omitFieldNames ? '' : 'phone')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Account clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Account copyWith(void Function(Account) updates) =>
      super.copyWith((message) => updates(message as Account)) as Account;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Account() / Account.new instead')
  static Account create() => Account._();
  static $pb.GeneratedMessage $_createMessage() => Account._();
  @$core.override
  Account createEmptyInstance() => Account._();
  @$core.pragma('dart2js:noInline')
  static Account getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Account>(Account.$_createMessage);
  static Account? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(2)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(2)
  void clearDescription() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get createdAt => $_getSZ(2);
  @$pb.TagNumber(3)
  set createdAt($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCreatedAt() => $_has(2);
  @$pb.TagNumber(3)
  void clearCreatedAt() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get updatedAt => $_getSZ(3);
  @$pb.TagNumber(4)
  set updatedAt($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasUpdatedAt() => $_has(3);
  @$pb.TagNumber(4)
  void clearUpdatedAt() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get disabledAt => $_getSZ(4);
  @$pb.TagNumber(5)
  set disabledAt($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasDisabledAt() => $_has(4);
  @$pb.TagNumber(5)
  void clearDisabledAt() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get email => $_getSZ(5);
  @$pb.TagNumber(6)
  set email($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasEmail() => $_has(5);
  @$pb.TagNumber(6)
  void clearEmail() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get phone => $_getSZ(6);
  @$pb.TagNumber(7)
  set phone($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasPhone() => $_has(6);
  @$pb.TagNumber(7)
  void clearPhone() => $_clearField(7);
}

class AccountLookupRequest extends $pb.GeneratedMessage {
  factory AccountLookupRequest() => AccountLookupRequest._();

  AccountLookupRequest._();

  factory AccountLookupRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AccountLookupRequest()..mergeFromBuffer(data, registry);
  factory AccountLookupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AccountLookupRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccountLookupRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: AccountLookupRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountLookupRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountLookupRequest copyWith(void Function(AccountLookupRequest) updates) =>
      super.copyWith((message) => updates(message as AccountLookupRequest))
          as AccountLookupRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AccountLookupRequest() / AccountLookupRequest.new instead')
  static AccountLookupRequest create() => AccountLookupRequest._();
  static $pb.GeneratedMessage $_createMessage() => AccountLookupRequest._();
  @$core.override
  AccountLookupRequest createEmptyInstance() => AccountLookupRequest._();
  @$core.pragma('dart2js:noInline')
  static AccountLookupRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AccountLookupRequest>(
          AccountLookupRequest.$_createMessage);
  static AccountLookupRequest? _defaultInstance;
}

class AccountLookupResponse extends $pb.GeneratedMessage {
  factory AccountLookupResponse({
    Account? account,
  }) {
    final result = AccountLookupResponse._();
    if (account != null) result.account = account;
    return result;
  }

  AccountLookupResponse._();

  factory AccountLookupResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AccountLookupResponse()..mergeFromBuffer(data, registry);
  factory AccountLookupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AccountLookupResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccountLookupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: AccountLookupResponse.$_createMessage)
    ..aOM<Account>(1, _omitFieldNames ? '' : 'account',
        subBuilder: Account.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountLookupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountLookupResponse copyWith(
          void Function(AccountLookupResponse) updates) =>
      super.copyWith((message) => updates(message as AccountLookupResponse))
          as AccountLookupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AccountLookupResponse() / AccountLookupResponse.new instead')
  static AccountLookupResponse create() => AccountLookupResponse._();
  static $pb.GeneratedMessage $_createMessage() => AccountLookupResponse._();
  @$core.override
  AccountLookupResponse createEmptyInstance() => AccountLookupResponse._();
  @$core.pragma('dart2js:noInline')
  static AccountLookupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AccountLookupResponse>(
          AccountLookupResponse.$_createMessage);
  static AccountLookupResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Account get account => $_getN(0);
  @$pb.TagNumber(1)
  set account(Account value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasAccount() => $_has(0);
  @$pb.TagNumber(1)
  void clearAccount() => $_clearField(1);
  @$pb.TagNumber(1)
  Account ensureAccount() => $_ensure(0);
}

class AccountUpdateRequest extends $pb.GeneratedMessage {
  factory AccountUpdateRequest({
    Account? account,
  }) {
    final result = AccountUpdateRequest._();
    if (account != null) result.account = account;
    return result;
  }

  AccountUpdateRequest._();

  factory AccountUpdateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AccountUpdateRequest()..mergeFromBuffer(data, registry);
  factory AccountUpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AccountUpdateRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccountUpdateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: AccountUpdateRequest.$_createMessage)
    ..aOM<Account>(1, _omitFieldNames ? '' : 'account',
        subBuilder: Account.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountUpdateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountUpdateRequest copyWith(void Function(AccountUpdateRequest) updates) =>
      super.copyWith((message) => updates(message as AccountUpdateRequest))
          as AccountUpdateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AccountUpdateRequest() / AccountUpdateRequest.new instead')
  static AccountUpdateRequest create() => AccountUpdateRequest._();
  static $pb.GeneratedMessage $_createMessage() => AccountUpdateRequest._();
  @$core.override
  AccountUpdateRequest createEmptyInstance() => AccountUpdateRequest._();
  @$core.pragma('dart2js:noInline')
  static AccountUpdateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AccountUpdateRequest>(
          AccountUpdateRequest.$_createMessage);
  static AccountUpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Account get account => $_getN(0);
  @$pb.TagNumber(1)
  set account(Account value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasAccount() => $_has(0);
  @$pb.TagNumber(1)
  void clearAccount() => $_clearField(1);
  @$pb.TagNumber(1)
  Account ensureAccount() => $_ensure(0);
}

class AccountUpdateResponse extends $pb.GeneratedMessage {
  factory AccountUpdateResponse({
    Account? account,
  }) {
    final result = AccountUpdateResponse._();
    if (account != null) result.account = account;
    return result;
  }

  AccountUpdateResponse._();

  factory AccountUpdateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AccountUpdateResponse()..mergeFromBuffer(data, registry);
  factory AccountUpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AccountUpdateResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccountUpdateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: AccountUpdateResponse.$_createMessage)
    ..aOM<Account>(1, _omitFieldNames ? '' : 'account',
        subBuilder: Account.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountUpdateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountUpdateResponse copyWith(
          void Function(AccountUpdateResponse) updates) =>
      super.copyWith((message) => updates(message as AccountUpdateResponse))
          as AccountUpdateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AccountUpdateResponse() / AccountUpdateResponse.new instead')
  static AccountUpdateResponse create() => AccountUpdateResponse._();
  static $pb.GeneratedMessage $_createMessage() => AccountUpdateResponse._();
  @$core.override
  AccountUpdateResponse createEmptyInstance() => AccountUpdateResponse._();
  @$core.pragma('dart2js:noInline')
  static AccountUpdateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AccountUpdateResponse>(
          AccountUpdateResponse.$_createMessage);
  static AccountUpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Account get account => $_getN(0);
  @$pb.TagNumber(1)
  set account(Account value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasAccount() => $_has(0);
  @$pb.TagNumber(1)
  void clearAccount() => $_clearField(1);
  @$pb.TagNumber(1)
  Account ensureAccount() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
