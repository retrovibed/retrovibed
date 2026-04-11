import 'package:flutter/material.dart';
import 'package:retrovibed/timex.dart' as timex;

class TimeRange extends StatelessWidget {
  final timex.Range range;
  final String Function(timex.Range)? format;

  const TimeRange(this.range, {super.key, this.format});

  static String _defaultFormat(timex.Range r) {
    final duration = r.end.difference(r.begin);

    if (duration.inDays >= 365) {
      final years = (duration.inDays / 365).round();
      return years == 1 ? '1 Year' : '$years Years';
    } else if (duration.inDays >= 30) {
      final months = (duration.inDays / 30).round();
      return months == 1 ? '1 Month' : '$months Months';
    } else if (duration.inDays >= 7) {
      final weeks = (duration.inDays / 7).round();
      return weeks == 1 ? '1 Week' : '$weeks Weeks';
    } else {
      return duration.inDays == 1 ? '1 Day' : '${duration.inDays} Days';
    }
  }

  @override
  Widget build(BuildContext context) {
    return Text(
      (format ?? _defaultFormat)(range),
      overflow: TextOverflow.ellipsis,
    );
  }
}
