import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/profiles.dart' as profiles;

class Display extends StatefulWidget {
  final AlignmentGeometry alignment;
  final EdgeInsets? margin;
  final EdgeInsets? padding;
  const Display({
    super.key,
    this.alignment = Alignment.topLeft,
    this.margin,
    this.padding,
  });

  @override
  State<Display> createState() => _DisplayState();
}

class _DisplayState extends State<Display> {
  final TextEditingController controller = TextEditingController();
  final ValueNotifier<int> events = ValueNotifier(0);

  @override
  void dispose() {
    controller.dispose();
    events.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final theme = Theme.of(context);
    return ds.Container(
      padding: widget.padding ?? defaults.padding,
      margin: widget.margin ?? defaults.margin,
      Column(
        spacing: defaults.spacing,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text("User Management", style: theme.textTheme.titleLarge),
            ],
          ),
          profiles.ListDisplay(controller: controller, events: events),
        ],
      ),
    );
  }
}
