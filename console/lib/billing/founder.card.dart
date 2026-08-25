import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './registered.dart';
import './plan.summary.dart';
import './founder.signal.detail.dart';

class FounderCard extends StatelessWidget {
  final EdgeInsets? margin;
  final Function(Widget w) onPressed;
  const FounderCard({super.key, required this.onPressed, this.margin});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    final billing = Registered.of(context);
    final plan = PlanSummary.fromID(billing.plan.id);
    final tappable = plan.id == founder().id;
    final tap = () => onPressed(
      SignalGroupDetail(margin: EdgeInsets.zero, padding: EdgeInsets.zero),
    );

    return ds.Card(
      alignment: Alignment.center,
      margin: this.margin ?? defaults.margin,
      help: ds.Hint(const Text("scan the founder Signal group invite QR code")),
      onTap: tappable ? tap : null,
      Column(
        spacing: defaults.spacing,
        children: [
          Text("Founder Chat", style: theme.textTheme.titleMedium),
          Text(
            "Signal group invite for founder-plan members",
            style: theme.textTheme.bodySmall,
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }
}
