import 'package:flutter/material.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/quotas.dart' as quotas;

class ArchiveStorage extends StatelessWidget {
  final httpx.Endpoint<String, quotas.QuotaFindResponse> quota;
  const ArchiveStorage({super.key, this.quota = quotas.api.get});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final authz = authn.Authenticated.bearer(context);
    return ds.Container(
      padding: defaults.padding,
      quotas.Typography.future(
        quota(quotas.Storage.sku, options: [authz]).then((v) => v.quota),
      ),
    );
  }
}
