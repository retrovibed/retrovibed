import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/meta.dart' as meta;

class DisabledIcon extends StatelessWidget {
  final Iterable<DateTime> dates;
  final double width;
  final double height;
  final Future<meta.Profile> Function(bool enabled)? onTap;
  const DisabledIcon({
    Key? key,
    required this.dates,
    this.width = 12.0,
    this.height = 12.0,
    this.onTap,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final DateTime earliestDate = timex.min(dates);
    final disabled = earliestDate.isBefore(timex.inf);
    final Color indicatorColor = disabled ? defaults.danger : defaults.success;

    return InkWell(
      borderRadius: defaults.borderRadius,
      onTap: this.onTap == null ? null : () => this.onTap?.call(!disabled),
      child: Container(
        width: width,
        height: height,
        decoration: BoxDecoration(
          color: indicatorColor,
          shape: BoxShape.circle,
        ),
      ),
    );
  }
}
