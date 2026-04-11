import 'package:flutter/material.dart';
import '../theme.defaults.dart';

class Masked extends StatelessWidget {
  final Widget child;
  final Alignment alignment;
  final BorderRadius borderRadius;
  final FocusScopeNode _selffocus = FocusScopeNode(debugLabel: "screen.masked");
  final Function()? reset;
  Masked(
    this.child, {
    super.key,
    this.borderRadius = BorderRadius.zero,
    this.alignment = Alignment.topLeft,
    this.reset,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);

    return LayoutBuilder(
      builder: (context, constraints) {
        final mq = MediaQuery.of(context);
        final available = Size(
          mq.size.width - mq.padding.left - mq.padding.right,
          mq.size.height - mq.padding.top - mq.padding.bottom,
        );
        final calculated = Size(
          constraints.hasBoundedWidth ? constraints.maxWidth : available.width,
          constraints.hasBoundedHeight
              ? constraints.maxHeight
              : available.height,
        );

        return FocusScope(
          canRequestFocus: true,
          autofocus: true,
          node: _selffocus,
          child: GestureDetector(
            onTap: reset,
            child: Container(
              constraints: BoxConstraints.tight(calculated),
              decoration: BoxDecoration(
                color: Colors.transparent.withValues(alpha: defaults.opaque.a),
              ),
              child: GestureDetector(
                onTap: () {},
                child: Align(alignment: alignment, child: child),
              ),
            ),
          ),
        );
      },
    );
  }
}
