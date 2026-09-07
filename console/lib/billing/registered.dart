import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/timex.dart' as timex;
import 'api.dart' as api;
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
    return context.dependOnInheritedWidgetOfExactType<RegisteredData>()?.state ?? RegisteredState();
  }

  @override
  State<StatefulWidget> createState() => RegisteredState();
}

class RegisteredState extends State<Registered> {
  final ValueNotifier<api.Billing> refresh = ValueNotifier(api.Billing());
  api.Billing current = api.Billing(subscriptionEndedAt: timex.inf.toIso8601String());
  api.Plan plan = PlanSummary.plan(free());
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  int attributionCount = 0;
  int attributionRate = 0;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

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
    return RegisteredData(
      state: this,
      child: ds.LoadingBoundary(
        loading: _loading,
        origin: 'RegisteredState',
        ds.ErrorScreen(
          cause: _cause,
          widget.child,
        ),
      ),
    );
  }
}

// Publishes the billing state to descendants. widget.child is a stable widget
// instance, so a setState here short circuits at Element.updateChild and never
// reaches the widgets reading Registered.of(context) - the dependency this
// registers is what actually propagates a lookup landing to them.
class RegisteredData extends InheritedWidget {
  final RegisteredState state;
  final api.Billing current;
  final api.Plan plan;
  final int attributionCount;
  final int attributionRate;

  RegisteredData({required this.state, required super.child})
    : current = state.current,
      plan = state.plan,
      attributionCount = state.attributionCount,
      attributionRate = state.attributionRate;

  @override
  bool updateShouldNotify(RegisteredData old) =>
      current != old.current ||
      plan != old.plan ||
      attributionCount != old.attributionCount ||
      attributionRate != old.attributionRate;
}
