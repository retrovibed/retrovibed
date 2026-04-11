import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import './plan.summary.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'api.dart' as api;

class Purchase extends StatelessWidget {
  final PlanSummary current;
  final api.Plan desired;
  final Future<void> Function(Future<api.BillingLookupResponse>) onChange;
  final Future<void> Function(String plan, {List<httpx.Option> options}) session;
  final Future<api.BillingLookupResponse> Function({List<httpx.Option> options}) lookup;
  final Duration interval;

  const Purchase({
    super.key,
    required this.current,
    required this.desired,
    required this.onChange,
    this.session = Purchase.defaultSession,
    this.lookup = api.lookup,
    this.interval = const Duration(seconds: 5),
  });

  static Future<void> defaultSession(
    String plan, {
    List<httpx.Option> options = const [],
  }) async {
    return api.session(plan, options: options).then((v) {
      return launchUrl(Uri.parse(v.redirect)).then((ok) {
        if (!ok) throw new Exception("failed to open url");
      });
    });
  }

  Future<api.BillingLookupResponse> poll(BuildContext context) {
    return httpx.withRetry(() => lookup(options: [authn.Authenticated.bearer(context)])).then((v) {
      if (!context.mounted) return v;
      if (v.billing.planId == desired.stripeId) return v;
      return Future.delayed(interval).then((_) => poll(context));
    });
  }

  @override
  Widget build(BuildContext context) {
    return ds.LoadingButton(
      Text("upgrade"),
      disabled: current.id == desired.id,
      onPressed: () async {
        return onChange(
          httpx
              .withRetry(
                () => session(
                  desired.token,
                  options: [authn.Authenticated.bearer(context)],
                ),
              )
              .then((_) => poll(context)),
        );
      },
    );
  }
}
