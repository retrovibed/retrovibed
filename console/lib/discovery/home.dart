import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'recent.dart';
import 'releases.dart';
import 'recommendations.dart';

class Home extends StatelessWidget {
  const Home(
    this.mimetype, {
    super.key,
    this.padding,
    this.margin,
    this.decoration,
    this.constraints,
    this.alignment,
    this.background,
    this.clipBehavior = Clip.none,
  });

  final String mimetype;
  final EdgeInsets? padding;
  final EdgeInsets? margin;
  final BoxDecoration? decoration;
  final BoxConstraints? constraints;
  final Alignment? alignment;
  final Color? background;
  final Clip clipBehavior;

  @override
  Widget build(BuildContext context) {
    return ds.build((context) {
      final defaults = ds.Defaults.of(context);
      final compact = defaults.isCompact;
      return ds.Container(
        alignment: alignment ?? Alignment.topCenter,
        background: background ?? Colors.transparent,
        padding: padding,
        margin: margin,
        decoration: decoration,
        constraints: constraints,
        clipBehavior: clipBehavior,
        Column(
          spacing: defaults.spacing,
          verticalDirection: compact ? VerticalDirection.up : VerticalDirection.down,
          children: [
            Recent(mimetype),
            if (authn.developer(context).recommendations) Recommendations(mimetype),
            if (authn.developer(context).releases) NewReleases(mimetype),
          ],
        ),
      );
    });
  }
}
