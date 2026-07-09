import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import './edit.dart';
import './authz.meta.edit.dart';
import './typography.dart' as typo;

void _Noop(meta.Profile up) {}

typedef FnUpdate =
    Future<meta.ProfileUpdateResponse> Function(
      meta.ProfileUpdateRequest req, {
      List<httpx.Option> options,
    });

typedef FnAuthzGet =
    Future<meta.AuthzProfileResponse> Function(
      String id, {
      List<httpx.Option> options,
    });

typedef FnAuthzGrant =
    Future<meta.AuthzGrantResponse> Function(
      String id,
      meta.Token token, {
      List<httpx.Option> options,
    });

class ListRow extends StatelessWidget {
  final meta.Profile current;
  final void Function(meta.Profile upd) onChange;
  final FnUpdate apiprofileupdate;
  final FnAuthzGet apiauthzget;
  final FnAuthzGrant apiauthzgrant;

  const ListRow(
    this.current, {
    super.key,
    this.onChange = _Noop,
    this.apiprofileupdate = meta.profiles.update,
    this.apiauthzget = meta.authz.get,
    this.apiauthzgrant = meta.authz.grant,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return ds.TableRow.single(
      typo.Typography(
        current,
        onChange: (pending) {
          pending.then((v) => onChange(v));
        },
      ),
      expanded: Container(
        padding: theme.buttonTheme.padding,
        child: Wrap(
          children: [
            Edit(
              current,
              onChange: (u, _) {
                apiprofileupdate(
                  meta.ProfileUpdateRequest(profile: u),
                  options: [authn.request(authn.AuthzCache.meta(context))],
                ).then((resp) {
                  onChange(resp.profile);
                });
              },
            ),
            AuthzMetaEdit.future(
              () => httpx.withRetry(
                () => apiauthzget(current.id, options: [authn.request(authn.AuthzCache.meta(context))]).then((v) => v.token),
              ),
              onChange: (t) {
                apiauthzgrant(current.id, t, options: [authn.request(authn.AuthzCache.meta(context))]).then((_) {
                  onChange(current);
                });
              },
            ),
          ],
        ),
      ),
    );
  }
}
