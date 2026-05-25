import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/billing/api.dart' as billingapi;

class CancellationButton extends StatelessWidget {
  static Future<void> noopidentitydelete({List<httpx.Option> options = const []}) => Future.value();
  final Future<void> Function({List<httpx.Option> options}) apibillingdelete;
  final Future<void> Function({List<httpx.Option> options}) apiidentitydelete;
  const CancellationButton({
    super.key,
    this.apibillingdelete = billingapi.delete,
    this.apiidentitydelete = noopidentitydelete,
  });

  @override
  Widget build(BuildContext context) {
    final authzmd = authn.AuthzCache.authzmetadata(context);
    final Future<void> Function() deleteaccount =
        () => apibillingdelete(
          options: [authn.DeeppoolAuthzCache.bearer(context)],
        ).then((_) {
          authn.Login.logout(context);
        });

    final Future<void> Function() deleteidentity =
        () => apiidentitydelete(
          options: [authn.DeeppoolAuthzCache.bearer(context)],
        ).then((_) {
          authn.Login.logout(context);
        });

    return ds.LoadingIconButton.delete(
      tooltip: 'Delete Account',
      help: ds.Hint(const Text('permanently remove this account')),
      onPressed:
          () => ds.modals.asyncfn<void>(
            context,
            (completion) => ds.Confirmation.yesNo(
              content: const Text(
                'Are you sure you want to delete your account? This action cannot be undone.',
              ),
              onConfirm: () {
                final pending = authzmd.billingModify ? deleteaccount() : deleteidentity();
                pending.then((_) => completion.complete()).catchError(completion.completeError);
              },
              onCancel: () => completion.complete(),
            ),
          ),
    );
  }
}
