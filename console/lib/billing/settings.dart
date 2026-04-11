import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import './api.dart' as api;
import './plan.summary.dart';
import './purchase.dart';
import './registered.dart';

class Settings extends StatefulWidget {
  final Alignment alignment;
  final EdgeInsets? margin;
  final EdgeInsets? padding;
  final Future<api.BillingPlansResponse> Function({List<httpx.Option> options}) apibillingplans;

  Settings({
    super.key,
    this.alignment = Alignment.topLeft,
    this.margin,
    this.padding,
    this.apibillingplans = api.plans,
  });

  @override
  State<Settings> createState() => _Settings();
}

class _Settings extends State<Settings> {
  Widget _cause = ds.Error.zero;
  RegisteredState? _billing;
  List<(PlanSummary, api.Plan)> _plans = [];
  PlanSummary current = free();
  (PlanSummary, api.Plan) desired = (free(), PlanSummary.plan(free()));

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void refresh() {
    final pid = Registered.of(context).current.planId;
    setState(() {
      desired =
          _plans.where((v) {
            return v.$2.stripeId == pid;
          }).firstOrNull ??
          desired;
      current = desired.$1;
    });
  }

  Future<void> _loadPlans() {
    return httpx
        .withRetry(
          () => widget.apibillingplans(
            options: [authn.Authenticated.bearer(context)],
          ),
        )
        .then((resp) {
          setState(() {
            _plans = resp.plans.map(PlanSummary.fromPlan).toList();
          });
        })
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: _reseterr);
          });
        });
  }

  void _reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  @override
  void initState() {
    super.initState();
    _billing = Registered.of(context);
    _billing?.refresh.addListener(refresh);
    _loadPlans().then((_) => refresh());
  }

  @override
  void dispose() {
    super.dispose();
    _billing?.refresh.removeListener(refresh);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final visible =
        _plans.where((p) {
          return !p.$1.hidden || p.$1.key == current.key;
        }).toList();
    return ds.ErrorScreen(
      cause: _cause,
      forms.Container(
        alignment: widget.alignment,
        margin: widget.margin,
        padding: widget.padding,
        Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            forms.Field(
              label: Text("plan"),
              input: DropdownButton(
                mouseCursor: SystemMouseCursors.click,
                borderRadius: defaults.borderRadius,
                alignment: Alignment.topLeft,
                isExpanded: true,
                value: desired.$1,
                items: [
                  for (final p in visible) DropdownMenuItem(child: p.$1.description, value: p.$1),
                ],
                onChanged: (v) {
                  setState(() {
                    desired = visible.firstWhere((x) => x.$1.id == (v ?? current).id);
                  });
                },
              ),
            ),
            desired.$1,
            Purchase(
              current: current,
              desired: desired.$2,
              onChange: (pending) {
                return pending
                    .then((v) {
                      _billing?.replace(v);
                    })
                    .catchError((cause) {
                      setState(() {
                        _cause = ds.Error.unknown(cause, onTap: _reseterr);
                      });
                    });
              },
            ),
          ],
        ),
      ),
    );
  }
}
