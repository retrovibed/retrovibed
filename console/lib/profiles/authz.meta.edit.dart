import 'package:flutter/material.dart';

import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as meta;
import './authz.permission.row.dart';

class AuthzMetaEdit extends StatelessWidget {
  final meta.Token current;
  final Function(meta.Token)? onChange;

  const AuthzMetaEdit(this.current, {super.key, this.onChange});

  static FutureBuilder<meta.Token> future(
    Future<meta.Token> pending, {
    Function(meta.Token)? onChange,
  }) {
    return ds.future(meta.Token(), pending, (snapshot) {
      return ds.ErrorScreen(
        AuthzMetaEdit(snapshot.data ?? meta.Token(), onChange: onChange),
        cause: snapshot.hasError ? ds.Error.unknown(snapshot.error!) : ds.Error.zero,
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return Wrap(
      spacing: defaults.spacing,
      runSpacing: defaults.spacing,
      alignment: WrapAlignment.start,
      children: [
        AuthzPermissionRow(
          "User Management",
          description: "Can manage user access",
          value: current.usermanagement,
          onChanged: (v) => onChange?.call(current..usermanagement = v),
        ),
        AuthzPermissionRow(
          "Library Read",
          description: "Can view library content",
          value: current.libraryRead,
          onChanged: (v) => onChange?.call(current..libraryRead = v),
        ),
        AuthzPermissionRow(
          "Library Modify",
          description: "Can modify library content",
          value: current.libraryModify,
          onChanged: (v) => onChange?.call(current..libraryModify = v),
        ),
        AuthzPermissionRow(
          "Community Modify",
          description: "Can modify community content",
          value: current.communityModify,
          onChanged: (v) => onChange?.call(current..communityModify = v),
        ),
        AuthzPermissionRow(
          "Billing Read",
          description: "Can view billing information",
          value: current.billingRead,
          onChanged: (v) => onChange?.call(current..billingRead = v),
        ),
        AuthzPermissionRow(
          "Billing Modify",
          description: "Can modify billing settings",
          value: current.billingModify,
          onChanged: (v) => onChange?.call(current..billingModify = v),
        ),
      ],
    );
  }
}
