import 'dart:io';

import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './registered.dart';
import './plan.summary.dart';
import './settings.dart';
import './cancellation.dart';

class Card extends StatelessWidget {
  final EdgeInsets? margin;
  final Function(Widget w) onPressed;
  const Card({super.key, required this.onPressed, this.margin});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final billing = Registered.of(context);
    final plan = PlanSummary.fromID(billing.plan.id);
    final defaults = ds.Defaults.of(context);
    final tap =
        () => onPressed(
          Settings(margin: EdgeInsets.zero, padding: EdgeInsets.zero),
        );

    return ds.Card(
      alignment: Alignment.center,
      margin: this.margin ?? defaults.margin,
      help: ds.Hint(const Text("view your account plan and billing details")),
      onTap: (Platform.isAndroid || Platform.isIOS || Platform.isMacOS) ? null : tap,
      Column(
        spacing: defaults.spacing,
        children: [
          Text("Account", style: theme.textTheme.titleMedium),
          DefaultTextStyle(
            style:
                theme.textTheme.bodyLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                ) ??
                TextStyle(fontWeight: FontWeight.bold),
            child: plan.description,
          ),
          DefaultTextStyle(
            style: theme.textTheme.bodySmall ?? TextStyle(),
            child: plan.price,
          ),
          ds.Timestamp.iso8601(
            leading: Text("Ends: "),
            billing.current.subscriptionEndedAt,
          ),
          CancellationButton(),
        ],
      ),
    );
  }
}
