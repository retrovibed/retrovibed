import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'recent.dart';
import 'releases.dart';
import 'recommendations.dart';

class Home extends StatefulWidget {
  const Home({
    super.key,
    this.padding,
    this.margin,
    this.decoration,
    this.constraints,
    this.alignment,
    this.background,
    this.clipBehavior = Clip.none,
  });

  final EdgeInsets? padding;
  final EdgeInsets? margin;
  final BoxDecoration? decoration;
  final BoxConstraints? constraints;
  final Alignment? alignment;
  final Color? background;
  final Clip clipBehavior;

  @override
  State<Home> createState() => _HomeState();
}

class _HomeState extends State<Home> {
  @override
  Widget build(BuildContext context) {
    return ds.build((context) {
      final defaults = ds.Defaults.of(context);
      final compact = defaults.isCompact;
      return ds.Container(
        alignment: widget.alignment ?? Alignment.topCenter,
        background: widget.background ?? Colors.transparent,
        padding: widget.padding,
        margin: widget.margin,
        decoration: widget.decoration,
        constraints: widget.constraints,
        clipBehavior: widget.clipBehavior,
        Column(
          spacing: defaults.spacing,
          verticalDirection: compact ? VerticalDirection.up : VerticalDirection.down,
          children: [
            Recent(),
            Recommendations(),
            NewReleases(),
          ],
        ),
      );
    });
  }
}
