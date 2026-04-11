import 'package:flutter/material.dart';
import 'package:retrovibed/billing/api.dart' as api;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/designkit.dart' as ds;

PlanSummary free() => const PlanSummary(
  id: "plans.free",
  description: Text("free"),
  key: ValueKey("plans.free"),
  storage: Text("none"),
  bandwidth: Text("none"),
  price: Text("\$0/month"),
  mobile: Text("no"),
);

PlanSummary personal4() => const PlanSummary(
  id: "plans.personal.2025",
  description: Text("personal"),
  key: ValueKey("plans.personal.2025"),
  storage: Text("2TB included + \$0.20 each additional 128 GB"),
  bandwidth: Text("12 GB / month"),
  price: Text("\$4/month"),
  mobile: Text("no"),
  hidden: true,
);

PlanSummary personal() => const PlanSummary(
  id: "plans.personal",
  description: Text("personal"),
  key: ValueKey("plans.personal"),
  storage: Text("2TB included + \$0.20 each additional 128 GB"),
  bandwidth: Text("12 GB / month"),
  price: Text("\$6/month"),
  mobile: Text("no"),
  hidden: true,
);

PlanSummary founder() => const PlanSummary(
  id: "plans.founder",
  description: Text("founder"),
  key: ValueKey("plans.founder"),
  storage: Text("2TB included + \$0.17 each additional 128 GB"),
  bandwidth: Text("12 GB / month"),
  price: Text("\$6/month"),
  mobile: Text("yes"),
);

PlanSummary family() => const PlanSummary(
  id: "plans.family",
  description: Text("family"),
  key: ValueKey("plans.family"),
  storage: Text("4TB included + \$0.18 each additional 128 GB"),
  bandwidth: Text("no hard limits, but abuse will result in rate limits"),
  price: Text("\$12/month"),
  mobile: Text("yes"),
  hidden: true,
);

PlanSummary premium() => const PlanSummary(
  id: "plans.premium",
  description: Text("premium"),
  key: ValueKey("plans.premium"),
  storage: Text("6TB included + \$0.17 each additional 128 GB"),
  bandwidth: Text("no hard limits, but abuse will result in rate limits"),
  price: Text("\$20/month"),
  mobile: Text("yes"),
  hidden: true,
);

class PlanSummary extends StatelessWidget {
  final String id;
  final bool hidden;
  final Widget description;
  final Widget storage;
  final Widget bandwidth;
  final Widget price;
  final Widget mobile;

  const PlanSummary({
    super.key,
    required this.id,
    required this.description,
    required this.storage,
    required this.price,
    required this.mobile,
    required this.bandwidth,
    this.hidden = false,
  });

  static (PlanSummary, api.Plan) fromPlan(api.Plan plan) {
    final display = fromID(plan.id);
    return (display, plan);
  }

  static api.Plan plan(PlanSummary s) {
    return api.Plan(id: s.id);
  }

  static PlanSummary fromID(String id) {
    final _personal4 = personal4();
    final _personal = personal();
    final _founder = founder();
    final _family = family();
    final _premium = premium();
    if (id == _personal4.id) return _personal4;
    if (id == _personal.id) return _personal;
    if (id == _founder.id) return _founder;
    if (id == _family.id) return _family;
    if (id == _premium.id) return _premium;
    return free();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return forms.Container(
      Column(
        mainAxisSize: MainAxisSize.min,
        spacing: defaults.spacing,
        children: [
          forms.Field(label: Text("Price"), input: price),
          forms.Field(label: Text("storage"), input: storage),
          forms.Field(
            label: Text("bandwidth"),
            input: Tooltip(
              message: "only related to downloading of archived data, accumulates monthly with a cap of 120 TB.",
              child: bandwidth,
            ),
          ),
          forms.Field(label: Text("mobile support"), input: mobile),
        ],
      ),
    );
  }
}
