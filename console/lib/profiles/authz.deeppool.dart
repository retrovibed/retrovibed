import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/authn.dart' as authn;
import './authz.permission.row.dart';

class AuthzDeeppool extends StatelessWidget {
  const AuthzDeeppool({super.key});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final token = authn.DeeppoolAuthzCache.of(context).meta.current.token ?? meta.Token();

    return forms.Container(
      padding: EdgeInsets.symmetric(horizontal: 10),
      Wrap(
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
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 288.0),
            child: forms.Field(
              label: Text("Archive Upload"),
              input: Text(token.archiveUpload.toString()),
            ),
          ),
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 288.0),
            child: forms.Field(
              label: Text("Archive Download"),
              input: Text(token.archiveDownload.toString()),
            ),
          ),
        ],
      ),
    );
  }
}
