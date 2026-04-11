import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/forms.dart' as forms;
import './registered.dart';

class ReferralDetail extends StatelessWidget {
  final Alignment alignment;
  final EdgeInsets? margin;
  final EdgeInsets? padding;

  const ReferralDetail({
    super.key,
    this.alignment = Alignment.topLeft,
    this.margin,
    this.padding,
  });

  @override
  Widget build(BuildContext context) {
    final billing = Registered.of(context);
    final count = billing.attributionCount;
    final rate = billing.attributionRate;
    final revenue = (count * rate / 100).toStringAsFixed(2);

    return forms.Container(
      alignment: alignment,
      margin: margin,
      padding: padding,
      Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          forms.Field(
            label: Text("referred users"),
            input: Text('$count'),
          ),
          forms.Field(
            label: Text("monthly revenue"),
            input: Text('\$$revenue'),
          ),
        ],
      ),
    );
  }
}
