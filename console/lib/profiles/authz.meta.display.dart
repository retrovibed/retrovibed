import 'package:flutter/widgets.dart';

import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/authn.dart' as authn;
import './authz.permission.row.dart';

class AuthzMetaDisplay extends StatelessWidget {
  final meta.Token token;

  const AuthzMetaDisplay(this.token, {super.key});

  static Widget current() {
    return Builder(
      builder: (context) {
        return AuthzMetaDisplay(
          authn.AuthzCache.of(context).meta.current.token,
        );
      },
    );
  }

  static FutureBuilder<meta.Token> future(Future<meta.Token> pending) {
    return ds.future(meta.Token(), pending, (snapshot) {
      return ds.ErrorScreen(
        AuthzMetaDisplay(snapshot.data ?? meta.Token()),
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
          value: token.usermanagement,
        ),
        AuthzPermissionRow(
          "Library Read",
          description: "Can view library content",
          value: token.libraryRead,
        ),
        AuthzPermissionRow(
          "Library Modify",
          description: "Can modify library content",
          value: token.libraryModify,
        ),
        AuthzPermissionRow(
          "Community Modify",
          description: "Can modify community content",
          value: token.communityModify,
        ),
        AuthzPermissionRow(
          "Billing Read",
          description: "Can view billing information",
          value: token.billingRead,
        ),
        AuthzPermissionRow(
          "Billing Modify",
          description: "Can modify billing settings",
          value: token.billingModify,
        ),
      ],
    );
  }
}
