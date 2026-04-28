import 'package:flutter/material.dart' as material;
import 'package:retrovibed/designkit.dart' as ds;

class Card extends material.StatelessWidget {
  final List<material.Widget> leading;
  final material.Widget child;
  final List<material.Widget> trailing;
  final material.Widget help;

  final material.Alignment? alignment;
  final material.BoxConstraints? constraints;
  final material.EdgeInsets margin;
  final material.EdgeInsets? padding;
  final material.FlexFit fit;
  final List<material.BoxShadow> tint;

  final material.GestureTapCallback? onTap;
  final material.GestureTapCallback? onDoubleTap;
  final material.GestureTapCallback? onSecondaryTap;
  final material.GestureLongPressCallback? onLongPress;

  const Card(
    this.child, {
    this.leading = const [],
    this.trailing = const [],
    this.onTap,
    this.onDoubleTap,
    this.onSecondaryTap,
    this.onLongPress,
    this.margin = material.EdgeInsets.zero,
    this.padding,
    this.tint = const [],
    this.fit = material.FlexFit.loose,
    this.alignment,
    this.constraints,
    this.help = ds.HelpScope.None,
  });

  @override
  material.Widget build(material.BuildContext context) {
    final defaults = ds.Defaults.of(context);

    final interactive = onTap != null || onDoubleTap != null || onSecondaryTap != null || onLongPress != null;
    return ds.Help(
      material.MouseRegion(
        cursor: interactive ? material.SystemMouseCursors.click : material.SystemMouseCursors.basic,
        child: material.Card(
          margin: margin,
          clipBehavior: material.Clip.antiAlias,
          child: material.InkWell(
            mouseCursor: material.MouseCursor.defer,
            onTap: onTap,
            onDoubleTap: onDoubleTap,
            onSecondaryTap: onSecondaryTap,
            onLongPress: onLongPress,
            child: material.Container(
              alignment: alignment,
              constraints: constraints,
              padding: padding ?? defaults.padding / 2,
              decoration: material.BoxDecoration(boxShadow: tint),
              child: material.Column(
                spacing: defaults.spacing / 2.5,
                mainAxisSize: material.MainAxisSize.min,
                children: [
                  ...leading,
                  material.Flexible(child: child, fit: fit),
                  ...trailing,
                ],
              ),
            ),
          ),
        ),
      ),
      help,
    );
  }
}
