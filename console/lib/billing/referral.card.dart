import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './registered.dart';
import './referral.detail.dart';

class ReferralCard extends StatelessWidget {
  final EdgeInsets? margin;
  final Function(Widget w) onPressed;
  const ReferralCard({super.key, required this.onPressed, this.margin});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    final billing = Registered.of(context);
    final count = billing.attributionCount;
    final rate = billing.attributionRate;
    final revenue = (count * rate / 100).toStringAsFixed(2);
    final tap = () => onPressed(
          ReferralDetail(margin: EdgeInsets.zero, padding: EdgeInsets.zero),
        );
    return ds.Card(
      alignment: Alignment.center,
      margin: margin ?? defaults.margin,
      help: ds.Hint(const Text("view your referral earnings")),
      onTap: tap,
      Column(
        spacing: defaults.spacing,
        children: [
          Text("Referrals", style: theme.textTheme.titleMedium),
          Text(
            '$count',
            style: theme.textTheme.headlineMedium?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          Text(
            '\$$revenue/mo',
            style: theme.textTheme.bodySmall,
          ),
        ],
      ),
    );
  }
}
