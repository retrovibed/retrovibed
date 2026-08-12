import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'api.dart' as api;
import 'package:retrovibed/design.kit/stateful.dart';
import 'plan.summary.dart';

class Registered extends StatefulWidget {
  final Widget child;
  final AlignmentGeometry alignment;
  final Future<api.BillingLookupResponse> Function({List<httpx.Option> options}) lookup;
  final Future<api.BillingCreateResponse> Function({List<httpx.Option> options}) create;

  const Registered(
    this.child, {
    super.key,
    this.alignment = Alignment.center,
    this.lookup = api.lookup,
    this.create = api.create,
  });

  static RegisteredState of(BuildContext context) {
    return context.findAncestorStateOfType<RegisteredState>() ?? RegisteredState();
  }

  @override
  State<StatefulWidget> createState() => RegisteredState();
}

class RegisteredState extends State<Registered> with LoadingState {
  final ValueNotifier<api.Billing> refresh = ValueNotifier(api.Billing());
  api.Billing current = api.Billing();
  api.Plan plan = PlanSummary.plan(free());
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  int attributionCount = 0;
  int attributionRate = 0;

  void replace(api.BillingLookupResponse upd) {
    setState(() {
      current = upd.billing;
      plan = upd.plan;
    });
    refresh.value = upd.billing;
  }

  void _reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  @override
  void initState() {
    super.initState();

    httpx
        .withRetry(
          () => widget.lookup(options: [authn.Authenticated.bearer(context)]),
        )
        .then((v) {
          setState(() {
            attributionCount = v.attributionCount.toInt();
            attributionRate = v.attributionRate;
            plan = v.plan;
          });

          if (v.billing.customerId.isEmpty) {
            return httpx.withRetry(
              () => widget.create(options: [authn.Authenticated.bearer(context)]).then((v) => v.billing),
            );
          }
          return Future.value(v.billing);
        })
        .then((billing) {
          setState(() {
            current = billing;
          });
          refresh.value = billing;
        })
        .catchError((_) {
          final billing = api.Billing(planId: free().id);
          setState(() {
            current = billing;
          });
          refresh.value = billing;
        }, test: httpx.ErrorsTest.err404)
        .catchError((cause) {
          setState(() {
            _cause = ds.Errors.httpauto(cause, onTap: _reseterr);
          });
        })
        .whenComplete(() {
          setState(() {
            _loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ds.LoadingBoundary(
      key: ValueKey(
        uuidx.md5x(
          current.customerId + current.subscriptionId + current.planId,
        ),
      ),
      loading: _loading,
      ds.ErrorScreen(
        cause: _cause,
        widget.child,
      ),
    );
  }
}
